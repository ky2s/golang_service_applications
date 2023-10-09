package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// interface
type ProjectController interface {
	ProjectList(c *gin.Context)
	ProjectCreate(c *gin.Context)
	ProjectUpdate(c *gin.Context)
	ProjectDestroy(c *gin.Context)
	ProjectFormCreate(c *gin.Context)
	ProjectFormDestroy(c *gin.Context)
}

type projectController struct {
	projectMod models.ProjectModels
	formMod    models.FormModels
	inputForm  models.InputFormModels
	compMod    models.CompaniesModels
	permissMod models.PermissionModels
}

func NewProjectController(projectModel models.ProjectModels, formModel models.FormModels, inputFormModel models.InputFormModels, compModel models.CompaniesModels, permModel models.PermissionModels) ProjectController {
	return &projectController{
		projectMod: projectModel,
		formMod:    formModel,
		inputForm:  inputFormModel,
		compMod:    compModel,
		permissMod: permModel,
	}
}

func (ctr *projectController) ProjectList(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	// organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println(userID, organizationID)
	}

	projectID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if projectID >= 1 {

		var fields tables.Projects
		fields.ID = projectID

		result, err := ctr.projectMod.GetProjectRow(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		//get form
		var whrProjectForm bytes.Buffer
		var whereString string
		// var whrStr string
		// var wherefrm tables.ProjectForms
		// wherefrm.ProjectID = result.ID

		// whrProjectForm.WriteString(" form_status_id in (1,2)  ")

		// if roleID != 1 {
		// 	whrProjectForm.WriteString(" AND f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")")
		// }
		// whrStr = whrProjectForm.String()

		// getFroms, err := ctr.projectMod.GetProjectForms(wherefrm, whrStr)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }
		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort
		if roleID == 1 {
			whrProjectForm.WriteString(" AND forms.form_status_id not in (3) AND project_id = " + strconv.Itoa(result.ID) + "")

		} else {
			whrProjectForm.WriteString("AND forms.form_status_id not in (3) AND project_id = " + strconv.Itoa(result.ID) + " AND ( forms.id in ( select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")")
		}
		if searchKeyWord != "" {
			whrProjectForm.WriteString(" AND  forms.name ilike '%" + searchKeyWord + "%'")
			// whereName = searchForm.String()
		}
		whereString = whrProjectForm.String()

		getFroms, _ := ctr.projectMod.GetFormInProject(whereString, organizationID, paging)
		fmt.Println("Role ID 1 ==>", getFroms)

		// fmt.Println(getFroms)
		// os.Exit(0)
		//get forms
		var formList []objects.MergeForms
		if len(getFroms) > 0 {

			for i := 0; i < len(getFroms); i++ {

				// checkFormIn, _ := ctr.projectMod.CheckFormIn(getFroms[i].ID, organizationID)
				// checkFormOut, _ := ctr.projectMod.CheckFormOut(getFroms[i].ID, organizationID)

				// fmt.Println(len(checkFormIn))
				// fmt.Println(len(checkFormOut))

				// get total responden
				var whereFU tables.JoinFormUsers
				whereFU.FormID = getFroms[i].ID
				whereFU.Type = "respondent"
				whreStr := ""

				getRespondenIntenal, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
				// get total responden form external
				var where tables.JoinFormUsers
				where.FormID = getFroms[i].ID
				where.Type = "respondent"
				whreString := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

				getRespondenExternal, err := ctr.formMod.GetFormUserRows(where, whreString)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon Internal
				var whereInForm tables.InputForms
				getDataResponFormInternal, err := ctr.inputForm.GetInputFormRows(getFroms[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon External
				whreStr = "ifo.organization_id = " + strconv.Itoa(organizationID) + ""
				getDataResponFormExternal, err := ctr.inputForm.GetInputFormRows(getFroms[i].ID, whereInForm, whreStr, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total active responden form internal
				whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondensInternal, err := ctr.inputForm.GetActiveUserInputForm(getFroms[i].ID, whreActive, whreStrAU)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				whreStrDate := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND ifo.organization_id = " + strconv.Itoa(organizationID) + ""
				// var whreActive tables.InputForms
				getActiveRespondensExternal, err := ctr.inputForm.GetActiveUserInputForm(getFroms[i].ID, whreActive, whreStrDate)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// get permission admin
				isPermission := false
				if roleID > 1 {
					var whr tables.FormUserPermissionJoin
					whr.PermissionID = 6 //(6 is edit responden)
					whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(getFroms[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
					getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
					isPermission = getPermission.Status
				}

				if userID == getFroms[i].CreatedBy {
					isPermission = true
				}

				if roleID == 1 {
					isPermission = true
				}

				var each objects.MergeForms
				each.ID = getFroms[i].ID
				each.Name = getFroms[i].Name
				each.Description = getFroms[i].Description
				each.ProfilePic = getFroms[i].ProfilePic
				each.FormStatusID = getFroms[i].FormStatusID
				each.FormStatus = getFroms[i].FormStatus
				each.Notes = getFroms[i].Notes
				each.PeriodStartDate = getFroms[i].PeriodStartDate
				each.PeriodEndDate = getFroms[i].PeriodEndDate
				each.CreatedBy = getFroms[i].CreatedBy
				each.CreatedByName = getFroms[i].CreatedByName
				each.CreatedByEmail = getFroms[i].CreatedByEmail
				each.IsAttendanceRequired = getFroms[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = getFroms[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = getFroms[i].SubmissionTargetUser
				each.PeriodeRange = 0
				each.IsEditResponden = isPermission
				each.Type = getFroms[i].Type
				if getFroms[i].Type == "internal" {
					each.TotalRespon = len(getDataResponFormInternal)
					each.TotalResponden = len(getRespondenIntenal)
					each.TotalRespondenActive = len(getActiveRespondensInternal)
				} else if getFroms[i].Type == "external" {
					each.TotalRespon = len(getDataResponFormExternal)
					each.TotalResponden = len(getRespondenExternal)
					each.TotalRespondenActive = len(getActiveRespondensExternal)
				}
				each.FormShared = getFroms[i].FormShared
				each.FormSharedCount = "Dibagikan ke " + strconv.Itoa(getFroms[i].FormShareCount) + " Organisasi"
				each.FormExternalCompanyName = getFroms[i].FormExternalCompanyName
				each.FormExternalCompanyImage = getFroms[i].FormExternalCompanyImage

				each.SharingSaldo = "Ya"
				each.StatusAdmin = "Editor"

				formList = append(formList, each)
			}
			// res.FormList = formList
		}

		var res objects.ProjectListRes
		res.ID = result.ID
		res.Name = result.Name
		res.Description = result.Description
		// res.AllRows = len(getFroms)

		var page objects.DataRows
		page.TotalRows = len(getFroms)
		page.TotalPages = 0

		var detail objects.DataRowsDetail
		detail.AllRows = len(getFroms)

		if len(formList) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"detail":      detail,
				"data_detail": res,
				"data_paging": page,
				"message":     "Data is available",
				"data":        formList,
				"status":      true,
			})
			return
		} else if len(formList) <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"detail":      detail,
				"data_detail": res,
				"data_paging": page,
				"message":     "Data is not available",
				"data":        formList,
				"status":      false,
			})
			return
		} else if res.ID > 0 {
			c.JSON(http.StatusOK, gin.H{
				"detail":      detail,
				"data_detail": res,
				"data_paging": page,
				"message":     "Data is available",
				"data":        nil,
				"status":      true,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"data_detail": res,
				"data_paging": page,
				"message":     "Data is not available",
				"data":        nil,
				"status":      false,
			})
			return
		}
	} else {

		// if roleID == 1 {

		// } else {
		var fields tables.Projects
		fields.CreatedBy = userID
		result, err := ctr.projectMod.GetProjectRows(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var res []objects.ProjectListRes
			for i := 0; i < len(result); i++ {

				var wherefrm tables.ProjectForms
				wherefrm.ProjectID = result[i].ID
				getFroms, err := ctr.projectMod.GetProjectForms(wherefrm, "")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var each objects.ProjectListRes
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.FormCount = len(getFroms)

				res = append(res, each)
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
		// }
	}
}

func (ctr *projectController) ProjectList__(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))

	ID := c.Param("id")
	iID, _ := strconv.Atoi(ID)

	if iID > 0 {

		var fields tables.Projects
		fields.ID = iID

		result, err := ctr.projectMod.GetProjectRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println("result.ID", result.ID)
		//get form
		var wherefrm tables.ProjectForms
		wherefrm.ProjectID = result.ID
		getFroms, err := ctr.projectMod.GetProjectForms(wherefrm, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var res objects.ProjectListRes
		res.ID = result.ID
		res.Name = result.Name
		res.Description = result.Description
		res.FormCount = len(getFroms)

		//get forms
		if len(getFroms) > 0 {
			var formList []objects.Forms
			for i := 0; i < len(getFroms); i++ {

				// get total responden
				var whereFU tables.JoinFormUsers
				whereFU.FormID = getFroms[i].FormID
				whereFU.Type = "respondent"
				whreStr := ""

				getResponden, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(getFroms[i].FormID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var each objects.Forms
				each.ID = getFroms[i].FormID
				each.FormStatusID = getFroms[i].FormStatusID
				each.FormStatus = getFroms[i].FormStatus
				each.Name = getFroms[i].Name
				each.Description = getFroms[i].Description
				each.ProfilePic = getFroms[i].ProfilePic
				each.CreatedByName = getFroms[i].CreatedByName
				each.CreatedByEmail = getFroms[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespon = len(getDataRespons)

				formList = append(formList, each)
			}
			res.FormList = formList
		}

		if result.ID > 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {

		var fields tables.Projects
		fields.CreatedBy = authorID
		result, err := ctr.projectMod.GetProjectRows(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var res []objects.ProjectListRes
			for i := 0; i < len(result); i++ {

				var wherefrm tables.ProjectForms
				wherefrm.ProjectID = result[i].ID
				getFroms, err := ctr.projectMod.GetProjectForms(wherefrm, "")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var each objects.ProjectListRes
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.FormCount = len(getFroms)

				res = append(res, each)
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}
}

func (ctr *projectController) ProjectCreate(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	var reqData objects.Projects
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

	var postData tables.Projects
	postData.Name = reqData.Name
	postData.Description = reqData.Description
	postData.CreatedBy = authorID

	var compID int
	if roleID == 1 {
		var whrComp tables.Organizations
		whrComp.CreatedBy = authorID
		whrComp.IsDefault = true
		getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)
		compID = getComp.ID
	} else {
		var uComp objects.UserOrganizations
		uComp.UserID = authorID
		getUserComp, _ := ctr.compMod.GetUserCompaniesRow(uComp, "")
		compID = getUserComp.OrganizationID
	}

	res, err := ctr.projectMod.InsertProject(postData, compID)
	if err != nil {
		fmt.Println("InsertProject", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var obj objects.ProjectRes
	obj.ID = res.ID

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created new data",
		"data":    obj,
	})
	return
}

func (ctr *projectController) ProjectUpdate(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	var reqData objects.Forms
	err = c.ShouldBindJSON(&reqData)
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

	var sendData tables.Projects
	sendData.Name = reqData.Name
	sendData.Description = reqData.Description
	sendData.UpdatedBy = authorID

	res, err := ctr.projectMod.UpdateProject(id, sendData)
	if err != nil {
		fmt.Println("UpdateProject", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res == true {
		var obj objects.ProjectRes
		obj.ID = id

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success update data",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed update data",
			"data":    nil,
		})
		return
	}
}

func (ctr *projectController) ProjectDestroy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	_, err = ctr.projectMod.DeleteProject(id)
	if err != nil {
		fmt.Println("DeleteProject", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success deleted data",
		"data":    nil,
	})
	return
}

func (ctr *projectController) ProjectFormCreate(c *gin.Context) {

	var reqData objects.ProjectForm
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

	var postData tables.ProjectForms
	postData.ProjectID = reqData.ProjectID
	postData.FormID = reqData.FormID

	res, err := ctr.projectMod.InsertProjectForm(postData)
	if err != nil {
		fmt.Println("InsertProjectForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var obj objects.ProjectRes
	obj.ID = res.ID

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created new data",
		"data":    obj,
	})
	return
}

func (ctr *projectController) ProjectFormDestroy(c *gin.Context) {

	var reqData objects.ProjectForm
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

	var fields tables.ProjectForms
	fields.ProjectID = reqData.ProjectID
	fields.FormID = reqData.FormID

	_, err = ctr.projectMod.DeleteProjectForm(fields)
	if err != nil {
		fmt.Println("DeleteProject", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success deleted data",
		"data":    nil,
	})
	return
}
