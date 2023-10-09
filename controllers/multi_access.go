package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"snapin-form/config"
	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"
	"strings"

	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

// interface
type MultiAccessController interface {
	InviteAdmin(c *gin.Context)
	SelectCompany(c *gin.Context)
}

type multiAccessController struct {
	helper         helpers.Helper
	multiAccessMod models.MultiAccessModels
	formMod        models.FormModels
	compMod        models.CompaniesModels
	userMod        models.UserModels
	conf           config.Configurations
	permissMod     models.PermissionModels
	pgErr          *pgconn.PgError
}

func NewMultiAccessController(h helpers.Helper, multiAccessModel models.MultiAccessModels, formModels models.FormModels, companiesModels models.CompaniesModels, userModels models.UserModels, configs config.Configurations, permissModel models.PermissionModels) MultiAccessController {
	return &multiAccessController{
		helper:         h,
		multiAccessMod: multiAccessModel,
		formMod:        formModels,
		compMod:        companiesModels,
		userMod:        userModels,
		conf:           configs,
		permissMod:     permissModel,
	}
}

func (ctr *multiAccessController) InviteAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	// roleID, _ := strconv.Atoi(claims["role_id"].(string))
	companyID := 0
	if len(claims) >= 5 {
		companyID, _ = strconv.Atoi(claims["organization_id"].(string))
	}

	var reqData objects.AdminInvites
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	getFormData, _ := ctr.formMod.GetFormRow(tables.Forms{ID: reqData.FormID})

	if len(reqData.AdminInvites) >= 1 {

		var whr tables.Organizations
		whr.ID = companyID
		getCompany, _ := ctr.compMod.GetCompaniesRow(whr)

		var sendMail bool
		for i := 0; i < len(reqData.AdminInvites); i++ {

			var whrUser tables.Users
			whrUser.Email = reqData.AdminInvites[i].Email
			getUser, _ := ctr.userMod.GetUserRow(whrUser)

			var whrUC objects.UserOrganizations
			whrUC.UserID = getUser.ID
			whrUC.OrganizationID = companyID
			checkUserComp, _ := ctr.compMod.GetUserCompaniesRow(whrUC, "")

			var whreUserEmail tables.Users
			whreUserEmail.Email = reqData.AdminInvites[i].Email
			whreUserEmail.RoleID = 7
			checkEmailGuest, _ := ctr.userMod.GetUserRow(whreUserEmail)
			// fmt.Println(checkEmailGuest)
			// os.Exit(0)
			if checkUserComp.ID <= 0 && checkEmailGuest.ID <= 0 && getUser.ID >= 1 && getUser.RoleID != 7 && getUser.RoleID != 1 {
				if reqData.AdminInvites[i].Email != "" {
					originInviteLink := ctr.conf.LINK_EXTERNAL + "/admin_invite?formid=" + strconv.Itoa(reqData.FormID) + "&senderid=" + strconv.Itoa(userID) + "&senderorgid=" + strconv.Itoa(companyID) + "&receiverid=" + strconv.Itoa(getUser.ID) + "&type=" + reqData.AdminInvites[i].Type

					to := reqData.AdminInvites[i].Email
					cc := ""
					subject := " Konfirmasi Penerimaan Form Multi Akses"
					message := `
					<br><img src="https://srv-asset-snapinnew.oss-ap-southeast-5.aliyuncs.com/dev/form_data_file/Logo-Project-Snap-In-Baru_6391d448be637c3444d244da.png" width="120">
					<br>
					<br>Halo ` + getUser.Name + `, 
					<br>
					<br>` + getCompany.Name + ` membagikan form <lable style="text-transform:capitalize;font-weight: bold;">` + getFormData.Name + `</lable> untuk anda sebagai <lable style="text-transform:capitalize;font-weight: bold;">` + reqData.AdminInvites[i].Type + `</lable>
					<br>
					<br>Klik link dibawah ini untuk konfirmasi penerimaan form
					<br><a href="` + originInviteLink + `"> ` + originInviteLink + ` </a>
					<br>
					<br>
					<br>
					<br>Perlu bantuan?
					<br>Client Service Snap-in : 0813 8585 7575 (WhatsApp Chat Only)
					<br>
					<br>Jam Kerja
					<br>Senin-Jumat : 09.00 s/d 18.00
					<br>Sabtu : 09.00 s/d 15.00
					`

					send, err := ctr.helper.SendGoMail(to, cc, subject, message)
					if err != nil {
						fmt.Println("err SendMail : ", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"error":   err,
							"message": "err SendGoMail 46546",
						})
						return
					}

					sendMail = send

				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Email is requred",
					})
					return
				}
			} else if checkEmailGuest.ID >= 1 {
				var postData tables.FormUsers
				postData.UserID = checkEmailGuest.ID
				postData.FormID = reqData.FormID
				postData.Type = "guest"
				postData.FormUserStatusID = 1

				updateUserForm, err := ctr.formMod.ConnectFormUser(postData)
				if err != nil {
					fmt.Println(err)

					if errors.As(err, &ctr.pgErr) {
						fmt.Println("ctr.pgErr.Code--", ctr.pgErr.Code)
						if ctr.pgErr.Code == "23505" { //code duplicate
							c.JSON(http.StatusOK, gin.H{
								"message": "Email has register before",
								"status":  true,
							})
							return
						} else {
							c.JSON(http.StatusBadRequest, gin.H{
								"status": true,
								"error":  err,
							})
							return
						}
					}

				}

				getPermiss, err := ctr.permissMod.GetPermissionRows(tables.Permissions{HttpPath: "/form"}, "")
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"status": false,
						"error":  err.Error(),
					})
					return
				}

				if len(getPermiss) > 0 {

					for i := 0; i < len(getPermiss); i++ {
						var postData tables.FormUserPermission
						postData.FormUserID = updateUserForm.ID
						postData.PermissionID = getPermiss[i].ID
						_, err = ctr.permissMod.InsertFormUserPermission(postData)
						if err != nil {
							if errors.As(err, &ctr.pgErr) {
								if ctr.pgErr.Code != "23505" { //code duplicate

									c.JSON(http.StatusBadRequest, gin.H{
										"status":  false,
										"error":   err.Error(),
										"message": "Error: InsertFormUserPermission 89758435",
									})
									return
								}
							}

						}
					}

				}

				originInviteLink := ctr.conf.LINK_EXTERNAL + "/guest?formid=" + strconv.Itoa(reqData.FormID) + "&senderid=" + strconv.Itoa(userID) + "&senderorgid=" + strconv.Itoa(companyID) + "&receiverid=" + strconv.Itoa(checkEmailGuest.ID) + "&type=" + reqData.AdminInvites[i].Type
				to := reqData.AdminInvites[i].Email
				cc := ""
				subject := " Konfirmasi Penerimaan Form Multi Akses"
				message := `
					<br><img src="https://srv-asset-snapinnew.oss-ap-southeast-5.aliyuncs.com/dev/form_data_file/Logo-Project-Snap-In-Baru_6391d448be637c3444d244da.png" width="120">
					<br>
					<br>Halo ` + getUser.Name + `, 
					<br>
					<br>` + getCompany.Name + ` membagikan form <lable style="text-transform:capitalize;font-weight: bold;">` + getFormData.Name + `</lable> untuk anda sebagai <lable style="text-transform:capitalize;font-weight: bold;">` + reqData.AdminInvites[i].Type + `</lable>
					<br>
					<br>Klik link dibawah ini untuk konfirmasi penerimaan form
					<br><a href="` + originInviteLink + `"> ` + originInviteLink + ` </a>
					<br>
					<br>
					<br>
					<br>Perlu bantuan?
					<br>Client Service Snap-in : 0813 8585 7575 (WhatsApp Chat Only)
					<br>
					<br>Jam Kerja
					<br>Senin-Jumat : 09.00 s/d 18.00
					<br>Sabtu : 09.00 s/d 15.00
					`

				send, err := ctr.helper.SendGoMail(to, cc, subject, message)
				if err != nil {
					fmt.Println("err SendMail : ", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"error":   err.Error(),
						"message": "err SendGoMail 46546",
					})
					return
				}

				sendMail = send
			} else if getUser.RoleID == 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":   "Akun Tidak bisa diinvite",
					"error": err,
				})
				return
			} else {
				fmt.Println("else :", reqData.AdminInvites[i].Email)

				if reqData.AdminInvites[i].Email != "" {

					// check invite privilage
					if reqData.AdminInvites[i].Type != "viewer" {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"message": "Email belum terdaftar di Snap-In. Anda bisa mengubahnya ke Viewer",
						})
						return
					}

					hash, err := helpers.Hash("123456")
					if err != nil {
						fmt.Println(err)
						return
					}

					hashEmail, err := helpers.Hash(reqData.AdminInvites[i].Email)
					if err != nil {
						fmt.Println(err)
						return
					}
					userName := strings.Split(reqData.AdminInvites[i].Email, "@")

					var userData tables.UsersMA
					userData.Name = userName[0]
					userData.Phone = ""
					userData.Email = reqData.AdminInvites[i].Email
					userData.RoleID = 7 //Guest

					lowName := strings.Replace(strings.ToLower(userName[0]), " ", "", -1)
					encryptCode := lowName[0:4] + helpers.EncodeToString(4)
					userData.EncryptCode = encryptCode

					if ctr.conf.ENV_TYPE != "dev" {
						userData.RememberToken = hashEmail
					}

					if ctr.conf.ENV_TYPE == "dev" {
						userData.IsEmailVerified = true
					}

					userData.Password = hash
					userData.IsEmailVerified = true
					res, err := ctr.userMod.InsertUser(userData)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}
					originInviteLink := ctr.conf.LINK_EXTERNAL + "/guest?formid=" + strconv.Itoa(reqData.FormID) + "&senderid=" + strconv.Itoa(userID) + "&senderorgid=" + strconv.Itoa(companyID) + "&receiverid=" + strconv.Itoa(res.ID) + "&type=" + reqData.AdminInvites[i].Type

					if res.ID >= 1 {

						// auto connect to organization
						// var dataPost objects.UserOrganizations
						// dataPost.UserID = res.ID
						// dataPost.OrganizationID = companyID
						// dataPost.IsDefault = true

						// _, err = ctr.compMod.ConnectedUserCompanies(dataPost)
						// if err != nil {
						// 	fmt.Println(err)
						// 	if ctr.pgErr.Code == "23505" {

						// 		c.JSON(http.StatusBadRequest, gin.H{
						// 			"status":  false,
						// 			"message": "Email has join in organization or another organization before",
						// 			"data":    nil,
						// 		})
						// 		return

						// 	} else {
						// 		c.JSON(http.StatusBadRequest, gin.H{
						// 			"status": false,
						// 			"error":  err.Error(),
						// 		})
						// 		return
						// 	}
						// }

						// connect to form
						var postData tables.FormUsers
						postData.UserID = res.ID
						postData.FormID = reqData.FormID
						postData.Type = "guest"
						postData.FormUserStatusID = 1

						updateUserForm, err := ctr.formMod.ConnectFormUser(postData)
						if err != nil {
							fmt.Println(err)

							if errors.As(err, &ctr.pgErr) {
								fmt.Println("ctr.pgErr.Code--", ctr.pgErr.Code)
								if ctr.pgErr.Code == "23505" { //code duplicate
									c.JSON(http.StatusOK, gin.H{
										"message": "Email has register before",
										"status":  true,
									})
									return
								} else {
									c.JSON(http.StatusBadRequest, gin.H{
										"status": true,
										"error":  err,
									})
									return
								}
							}
						}

						// connect to permission
						var whre tables.Permissions
						whre.HttpPath = "/form"

						getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
						if err != nil {
							fmt.Println(err)
							c.JSON(http.StatusBadRequest, gin.H{
								"error": err,
							})
							return
						}
						if len(getPermiss) > 0 {

							for i := 0; i < len(getPermiss); i++ {
								var postData tables.FormUserPermission
								postData.FormUserID = updateUserForm.ID
								postData.PermissionID = getPermiss[i].ID
								_, err = ctr.permissMod.InsertFormUserPermission(postData)
								if err != nil {
									if errors.As(err, &ctr.pgErr) {
										if ctr.pgErr.Code == "23505" { //code duplicate

											c.JSON(http.StatusOK, gin.H{
												"status":  true,
												"message": "Data member has join in Admin available",
												"data":    nil,
											})
											return

										} else {
											c.JSON(http.StatusBadRequest, gin.H{
												"error": err,
											})
											return
										}
									}

								}
							}

						}

					} else {
						fmt.Println("err: gagal register")
						c.JSON(http.StatusBadRequest, gin.H{
							"msg":   "Failed register",
							"error": err,
						})
						return
					}

					to := reqData.AdminInvites[i].Email
					cc := ""
					subject := " Konfirmasi Penerimaan Form Multi Akses"
					message := `
					<br><img src="https://srv-asset-snapinnew.oss-ap-southeast-5.aliyuncs.com/dev/form_data_file/Logo-Project-Snap-In-Baru_6391d448be637c3444d244da.png" width="120">
					<br>
					<br>Halo ` + getUser.Name + `, 
					<br>
					<br>` + getCompany.Name + ` membagikan form <lable style="text-transform:capitalize;font-weight: bold;">` + getFormData.Name + `</lable> untuk anda sebagai <lable style="text-transform:capitalize;font-weight: bold;">` + reqData.AdminInvites[i].Type + `</lable>
					<br>
					<br>Klik link dibawah ini untuk konfirmasi penerimaan form
					<br><a href="` + originInviteLink + `"> ` + originInviteLink + ` </a>
					<br>
					<br>
					<br>
					<br>Perlu bantuan?
					<br>Client Service Snap-in : 0813 8585 7575 (WhatsApp Chat Only)
					<br>
					<br>Jam Kerja
					<br>Senin-Jumat : 09.00 s/d 18.00
					<br>Sabtu : 09.00 s/d 15.00
					`

					send, err := ctr.helper.SendGoMail(to, cc, subject, message)
					if err != nil {
						fmt.Println("err SendMail : ", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"error":   err.Error(),
							"message": "err SendGoMail 46546",
						})
						return
					}

					sendMail = send

					var whreU tables.Users
					whreStr := "email = '" + reqData.AdminInvites[i].Email + "' AND role_id in (7)"
					checkAdmin, err := ctr.userMod.GetUserWhereRow(whreU, whreStr)
					// fmt.Println(checkAdmin.ID)
					// os.Exit(0)

					var dataPost objects.UserOrganizations
					dataPost.UserID = checkAdmin.ID
					dataPost.OrganizationID = companyID
					dataPost.IsDefault = true

					userComp, err := ctr.compMod.ConnectedUserCompanies(dataPost)
					var whre tables.Permissions
					whre.HttpPath = "/admin-global"

					// getUserOrg, err := ctr.userMod.GetUserOrganization(userID)
					getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
					if len(getPermiss) > 0 {
						fmt.Println(len(getPermiss))
						for i := 0; i < len(getPermiss); i++ {
							var postData tables.UserOrganizationPermission
							postData.UserOrganizationID = userComp.ID
							postData.PermissionID = getPermiss[i].ID
							_, err = ctr.permissMod.InsertUserOrganizationPermission(postData)
							if err != nil {
								if errors.As(err, &ctr.pgErr) {
									if ctr.pgErr.Code == "23505" { //code duplicate

										c.JSON(http.StatusOK, gin.H{
											"status":  true,
											"message": "Data member has join in Admin available",
											"data":    nil,
										})
										return

									} else {
										c.JSON(http.StatusBadRequest, gin.H{
											"error": err,
										})
										return
									}
								}

							}
						}
						c.JSON(http.StatusOK, gin.H{
							"status":  true,
							"message": "Success sign data to form",
							"data":    nil,
						})
						return
					}

				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Email is required",
					})
					return
				}
			}
		}

		if sendMail {
			fmt.Println("send mail success")

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Success invites new company data",
			})
			return
		} else {
			fmt.Println("send mail failed")

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed invites new company data",
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed: email input is required",
		})
		return
	}

}

func (ctr *multiAccessController) SelectCompany(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	companyID := 0
	if len(claims) >= 5 {
		companyID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("companyID :::", companyID, userID, roleID)
	}

	var reqData objects.SelectOrganization
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	checkFormAlreadyHave, _ := ctr.multiAccessMod.CheckFormAlreadyHave(reqData.FormID, reqData.UserReceiverID, reqData.OrganizationReceiverID)

	if len(checkFormAlreadyHave) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Admin already in this Form with same Organization",
		})
		return

	} else {
		//Get Sender Role
		getSender, _ := ctr.multiAccessMod.GetSenderRole(reqData.UserSenderID)
		getReceiver, _ := ctr.multiAccessMod.GetSenderRole(reqData.UserReceiverID)

		var postData objects.SelectOrganization
		postData.FormID = reqData.FormID
		postData.IsQuotaSharing = false
		postData.AccessType = reqData.AccessType
		postData.UserSenderID = reqData.UserSenderID
		postData.UserReceiverID = reqData.UserReceiverID
		postData.UserSenderTypeID = getSender.RoleID
		postData.UserReceiverTypeID = getReceiver.RoleID
		postData.OrganizationSenderID = reqData.OrganizationSenderID
		postData.OrganizationReceiverID = reqData.OrganizationReceiverID
		_, err = ctr.multiAccessMod.InsertFormToUserInvites(postData)
		if err != nil {
			fmt.Println("InsertForm", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err,
				"message": "err ConnectFormUserOrg 44645",
			})
			return
		}

		var dataInput objects.InputFormUserOrganizations
		dataInput.FormID = reqData.FormID
		dataInput.UserID = reqData.UserReceiverID
		dataInput.OrganizationID = reqData.OrganizationReceiverID
		dataInput.Type = "admin"

		checkConnectToFormUser, _ := ctr.multiAccessMod.CheckUserAlreadyConnect(dataInput)
		if len(checkConnectToFormUser) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "User Already in this Form",
			})
			return

		} else {
			//insert to form_users & form_user_organizations
			insertToFormUser, err := ctr.formMod.ConnectFormUserOrg(dataInput)
			if err != nil {
				fmt.Println("InsertFormUser", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			// fmt.Println(insertToFormUser)
			var whre tables.Permissions
			whre.HttpPath = "/form"
			getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"error":   err,
					"message": "err GetPermissionRows 65465",
				})
				return
			}
			if len(getPermiss) >= 1 {

				for i := 0; i < len(getPermiss); i++ {
					var postData tables.FormUserPermission
					postData.FormUserID = insertToFormUser.ID
					postData.PermissionID = getPermiss[i].ID
					if reqData.AccessType == "editor" {
						if getPermiss[i].ID == 8 {
							postData.Status = false
						} else {
							postData.Status = true

						}
					} else {
						if getPermiss[i].ID == 7 {
							postData.Status = true
						} else {
							postData.Status = false

						}
					}
					_, _ = ctr.permissMod.InsertFormUserPermission(postData)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success connected user to companies",
		})
		return
	}
}
