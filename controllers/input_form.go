package controllers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"snapin-form/config"
	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/360EntSecGroup-Skylar/excelize"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// interface
type InputFormController interface {
	FieldSaveData(c *gin.Context)
	FieldUpdateData(c *gin.Context)
	FieldDestroyData(c *gin.Context)
	DataFormDetail(c *gin.Context)
	DataFormDelete(c *gin.Context)
	DataFormDetailDownload(c *gin.Context)
	DataFormDetail2Download(c *gin.Context) // existing used report
	DataFormDetail3Download(c *gin.Context)
	DataFormDetail4Download(c *gin.Context) // multy akses latest
	DataFormDetailDownloadCSV(c *gin.Context)
	DataFormDetail4DownloadCSV(c *gin.Context) // multy akses
	DataFormDetailGrafic(c *gin.Context)
	DataFormDetailGraficByOrganization(c *gin.Context)
	DataFormResponList(c *gin.Context)
	DataFormRespondenList(c *gin.Context)
	ReportFormRespondenList(c *gin.Context)
	ReportFormRespondenList2(c *gin.Context) // admin apps(report)
	DataFormResponMapList(c *gin.Context)    // map report
	DataFormResponMapPostList(c *gin.Context)
	// DataFormPerUserMapList(c *gin.Context)    // map report

	DataFormCompanyResponList(c *gin.Context) // lihat data with company

	GenerateFormUserOrg(c *gin.Context)
	GenerateInputFormOrg(c *gin.Context)
	GenerateInputFormOrgNoCopy(c *gin.Context)
	GenerateInputFormOrgLatest(c *gin.Context) //new
}

type inputFormController struct {
	inputForm models.InputFormModels
	formMod   models.FormModels
	formField models.FormFieldModels
	helper    helpers.Helper
	userMod   models.UserModels
	subsMod   models.SubsModels
	compMod   models.CompaniesModels
	attMod    models.AttendanceModels
	settModel models.SettingModels
	shrtMod   models.ShortenUrlModels
	conf      config.Configurations
}

func NewInputFormController(inputMod models.InputFormModels, formField models.FormFieldModels, frmMod models.FormModels, help helpers.Helper, userModel models.UserModels, subsModel models.SubsModels, compModel models.CompaniesModels, attModel models.AttendanceModels, settModel models.SettingModels, shorModel models.ShortenUrlModels, configs config.Configurations) InputFormController {
	return &inputFormController{
		inputForm: inputMod,
		formField: formField,
		formMod:   frmMod,
		helper:    help,
		userMod:   userModel,
		subsMod:   subsModel,
		compMod:   compModel,
		attMod:    attModel,
		settModel: settModel,
		shrtMod:   shorModel,
		conf:      configs,
	}
}

func (ctr *inputFormController) FieldSaveData(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	myCompanyID := 0
	if len(claims) >= 5 {
		myCompanyID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("companyID :::", myCompanyID)
	}

	fmt.Println("myCompanyID :::", myCompanyID)

	var reqData objects.FormData
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

	if len(reqData.FieldData) > 0 {
		reqData.UserID = userID

		// get form user organizations ID
		getFormUserOrg, _ := ctr.formMod.GetFormUserToOrganizationRow(tables.JoinFormUsers{FormID: reqData.FormID, UserID: userID}, "")

		// get form Company
		getFormOrg, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: reqData.FormID})

		// cek organization dari form_organization
		// next nya getFormUserOrg.OrganizationID tidak terdeteksi karena relasi antara TIM & FORM dimana user terhubung melalui tim
		// jdi harus dicari organization ID dari team_form
		if getFormUserOrg.OrganizationID <= 0 {
			getFormUserOrg.OrganizationID = getFormOrg.OrganizationID
		}

		// cheking quota zero
		// InsertInjuryPlan
		checkQuota, _ := ctr.subsMod.GetPlanRow(objects.SubsPlan{OrganizationID: getFormUserOrg.OrganizationID})
		if checkQuota.IsBlocked == true {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Your data respon isn't saved",
			})
			return
		}

		// insert input form /submit data --------------------------------------------------------------
		_, err = ctr.inputForm.InsertFormDataWithOrganization(reqData, getFormUserOrg.OrganizationID)
		if err != nil {
			fmt.Println("InsertFormData", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if checkQuota.QuotaCurrent == 0 {

			var whreSett tables.Settings
			whreSett.Code = "injury_plan_quota"
			settingInjuryQuota, _ := ctr.settModel.GetSettingRow(whreSett)
			quotaValue, _ := strconv.Atoi(settingInjuryQuota.Value)

			var whrPlan objects.SubsPlan
			whrPlan.OrganizationID = getFormUserOrg.OrganizationID
			checkPlan, _ := ctr.subsMod.GetPlanRow(whrPlan)

			var whrInjPlan objects.InjuryPlan
			whrInjPlan.OrganizationSubscriptionPlanID = checkPlan.ID
			checkInjuryPlan, _ := ctr.subsMod.GetInjuryPlanRows(whrInjPlan)

			fmt.Println("injury cek ------------------------------>", checkInjuryPlan, checkPlan.ID, len(checkInjuryPlan))

			// ke 2 kali nya tidak di auto add Quota
			if quotaValue > 0 && len(checkInjuryPlan) == 0 {
				// var dataPostInjury tables.SubsPlan
				// dataPostInjury.QuotaCurrent = quotaValue
				// dataPostInjury.QuotaTotal = quotaValue
				// addnjuryQuota, injuryRes, err := ctr.subsMod.UpdatePlan(getFormOrg.OrganizationID, dataPostInjury)
				// if err != nil {
				// 	fmt.Println("UpdatePlan dataPostInjury", err)
				// 	c.JSON(http.StatusBadRequest, gin.H{
				// 		"error": err,
				// 	})
				// 	return
				// }

				// if addnjuryQuota == true {
				var injuryData tables.InjuryPlan
				injuryData.OrganizationSubscriptionPlanID = checkPlan.ID
				injuryData.Quota = quotaValue
				ctr.subsMod.InsertInjuryPlan(injuryData)
				// }
			}
		}

		//insert submiting quota respon
		var dataPost tables.SubsPlan
		dataPost.RespondentCurrent = 1
		_, _, err = ctr.subsMod.UpdatePlanCurrent(getFormOrg.OrganizationID, dataPost)
		if err != nil {
			fmt.Println("UpdatePlanCurrent", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Submission Sukses",
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data input is empty",
		})
		return
	}
}

func (ctr *inputFormController) FieldUpdateData(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	submissionID, err := strconv.Atoi(c.Param("dataid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	var reqData objects.FormData
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

	if len(reqData.FieldData) > 0 {
		reqData.UserID = userID
		_, err = ctr.inputForm.UpdateFormData(submissionID, reqData)
		if err != nil {
			fmt.Println("InsertFormData", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		var whreFrm tables.FormOrganizations
		whreFrm.FormID = reqData.FormID
		getFormOrg, _ := ctr.formMod.GetFormOrganization(whreFrm)

		// cheking zero quota
		// InsertInjuryPlan
		var whr objects.SubsPlan
		whr.OrganizationID = getFormOrg.OrganizationID
		checkQuota, _ := ctr.subsMod.GetPlanRow(whr)

		fmt.Println("checkQuota.QuotaCurrent>>>>>", checkQuota.QuotaCurrent)
		if checkQuota.IsBlocked == true {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Your data respon isn't saved",
			})
			return
		}

		if checkQuota.QuotaCurrent == 0 {

			var whreSett tables.Settings
			whreSett.Code = "injury_plan_quota"
			settingInjuryQuota, _ := ctr.settModel.GetSettingRow(whreSett)
			quotaValue, _ := strconv.Atoi(settingInjuryQuota.Value)

			var whrPlan objects.SubsPlan
			whrPlan.OrganizationID = getFormOrg.OrganizationID
			checkPlan, _ := ctr.subsMod.GetPlanRow(whrPlan)

			var whrInjPlan objects.InjuryPlan
			whrInjPlan.OrganizationSubscriptionPlanID = checkPlan.ID
			checkInjuryPlan, _ := ctr.subsMod.GetInjuryPlanRows(whrInjPlan)

			fmt.Println("injury cek ------------------------------>", checkInjuryPlan, checkPlan.ID, len(checkInjuryPlan))

			// ke 2 kali nya tidak di auto add Quota
			if quotaValue > 0 && len(checkInjuryPlan) == 0 {
				// var dataPostInjury tables.SubsPlan
				// dataPostInjury.QuotaCurrent = quotaValue
				// dataPostInjury.QuotaTotal = quotaValue
				// addnjuryQuota, injuryRes, err := ctr.subsMod.UpdatePlan(getFormOrg.OrganizationID, dataPostInjury)
				// if err != nil {
				// 	fmt.Println("UpdatePlan dataPostInjury", err)
				// 	c.JSON(http.StatusBadRequest, gin.H{
				// 		"error": err,
				// 	})
				// 	return
				// }

				// if addnjuryQuota == true {
				var injuryData tables.InjuryPlan
				injuryData.OrganizationSubscriptionPlanID = checkPlan.ID
				injuryData.Quota = quotaValue
				ctr.subsMod.InsertInjuryPlan(injuryData)
				// }
			}
		}

		// getCountUpdt, _ := ctr.inputForm.GetUpdatedCount(reqData.FormID, submissionID)

		var updCnt objects.UpdCnt
		updCnt.UpdatedCount = 1
		_, _, err := ctr.inputForm.UpdatedCount(reqData.FormID, submissionID, updCnt)
		if err != nil {
			fmt.Println("UpdatePlanCurrent", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		//insert submiting quota respon
		var dataPost tables.SubsPlan
		dataPost.RespondentCurrent = 1
		_, _, err = ctr.subsMod.UpdatePlanCurrent(getFormOrg.OrganizationID, dataPost)
		if err != nil {
			fmt.Println("UpdatePlanCurrent", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Ubah submission sukses",
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data input is empty",
		})
		return
	}
}

func (ctr *inputFormController) FieldDestroyData(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	submissionID, err := strconv.Atoi(c.Param("dataid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	var reqData objects.FormData
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

	if submissionID > 0 {
		reqData.UserID = userID
		ctr.inputForm.DeleteFormData(submissionID, reqData)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data input has deleted successfully",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data input ID is required",
		})
		return
	}
}

func (ctr *inputFormController) DataFormDetail(c *gin.Context) {
	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate

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

func (ctr *inputFormController) DataFormDetailGrafic(c *gin.Context) {

	ID := c.Param("formid")
	periode := c.Param("periode")
	year := c.Param("year")
	if year == "" {
		year = time.Now().Format("2006")
	}

	graficType := c.Request.URL.Query().Get("type")

	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		formID := result.ID

		if formID > 0 {
			//get total responden
			// var whereFU tables.JoinFormUsers
			// whereFU.FormID = formID
			// whereFU.Type = "respondent"

			var bufferActiveUser bytes.Buffer
			var whreStrAU string

			if periode == "daily" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
			} else if periode == "monthly" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")
			} else if periode == "yearly" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy') = TO_CHAR(NOW()::date, 'yyyy') ")
			}

			if year != "" {
				bufferActiveUser.WriteString(" AND TO_CHAR(if.created_at::date, 'yyyy') = '" + year + "'")
			}
			whreStrAU = bufferActiveUser.String()
			var whreActive tables.InputForms
			getRespondens, err := ctr.inputForm.GetActiveUserInputForm(formID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//get total respon
			var bufferRes bytes.Buffer
			var whereInForm tables.InputForms

			var whreStrIF string
			if periode == "daily" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
			} else if periode == "monthly" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")
			} else if periode == "yearly" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy') = TO_CHAR(NOW()::date, 'yyyy') ")
			}

			if year != "" {
				bufferRes.WriteString(" AND TO_CHAR(if.created_at::date, 'yyyy') = '" + year + "'")
			}
			whreStrIF = bufferRes.String()
			getDataRespons, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whreStrIF, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//get hours

			var getDataHours []objects.GraficDataHours
			var getDataPeriode []objects.GraficDataPeriod

			if graficType == "submission" {
				var bufferHrs bytes.Buffer
				var whreStrHrs string
				if periode != "" {
					bufferHrs.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
				}
				if year != "" {
					bufferHrs.WriteString(" AND TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'")
				}
				whreStrHrs = bufferHrs.String()

				getDataHours, err = ctr.inputForm.GetDataHours(formID, whreStrHrs)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if periode == "daily" {
					whreDailySub := " where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeDays(formID, whreDailySub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "monthly" {
					whreMonthlySub := " where TO_CHAR(f.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') "
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeMonthly(formID, whreMonthlySub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "yearly" {
					whreYearSub := " where TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'"
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeYearly(formID, whreYearSub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
				}

			} else if graficType == "responden" {

				var bufferHrs bytes.Buffer
				var whreStrHrs string
				if periode != "" {
					bufferHrs.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
				}
				if year != "" {
					bufferHrs.WriteString(" AND TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'")
				}
				whreStrHrs = bufferHrs.String()

				getDataHours, err = ctr.inputForm.GetDataHoursUserResp(formID, whreStrHrs)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if periode == "daily" {
					whreDailyResp := " where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeDaysResp(formID, whreDailyResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "monthly" {

					whreMonthlyResp := " where TO_CHAR(f.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') "
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeMonthlyResp(formID, whreMonthlyResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "yearly" {
					whreYearResp := " where TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'"
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeYearlyResp(formID, whreYearResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				}
			}

			// var hoursResult []objects.GraficData
			// if len(getHours) > 0 {
			// 	for i := 0; i < len(getHours); i++ {
			// 		var ec objects.GraficData
			// 		ec.Field = getHours[i].Hours
			// 		ec.Value = "0"

			// 		hoursResult = append(hoursResult, ec)
			// 	}
			// }

			var res objects.FormGraficRes
			res.FormID = formID
			res.TotalResponden = len(getRespondens)
			res.TotalRespon = len(getDataRespons)
			res.ActiveHours = getDataHours
			res.ActivePeriod = getDataPeriode

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form id is not found",
			"data":    nil,
		})
		return
	}
}

func (ctr *inputFormController) DataFormDetailGraficByOrganization(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	ID := c.Param("formid")
	periode := c.Param("periode")
	year := c.Param("year")
	if year == "" {
		year = time.Now().Format("2006")
	}

	graficType := c.Request.URL.Query().Get("type")
	companyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, _ := strconv.Atoi(companyID)

	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		formID := result.ID

		if formID > 0 {
			//get total responden
			// var whereFU tables.JoinFormUsers
			// whereFU.FormID = formID
			// whereFU.Type = "respondent"

			var bufferActiveUser bytes.Buffer
			var whreStrAU string

			if periode == "daily" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
			} else if periode == "monthly" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")
			} else if periode == "yearly" {
				bufferActiveUser.WriteString(" TO_CHAR(if.created_at::date, 'yyyy') = TO_CHAR(NOW()::date, 'yyyy') ")
			}

			if year != "" {
				bufferActiveUser.WriteString(" AND TO_CHAR(if.created_at::date, 'yyyy') = '" + year + "'")
			}

			checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
			if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
				// company ID by TOKEN
				bufferActiveUser.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
			} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
				// select company option (form sharing only)
				bufferActiveUser.WriteString(" AND ifo.organization_id= " + strconv.Itoa(iCompanyID))
			}

			whreStrAU = bufferActiveUser.String()
			getRespondens, err := ctr.inputForm.GetActiveUserInputForm(formID, tables.InputForms{}, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"status": false,
					"error":  err.Error(),
				})
				return
			}

			//get total respon
			var bufferRes bytes.Buffer
			var whereInForm tables.InputForms

			var whreStrIF string
			if periode == "daily" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
			} else if periode == "monthly" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")
			} else if periode == "yearly" {
				bufferRes.WriteString(" TO_CHAR(if.created_at::date, 'yyyy') = TO_CHAR(NOW()::date, 'yyyy') ")
			}

			if year != "" {
				bufferRes.WriteString(" AND TO_CHAR(if.created_at::date, 'yyyy') = '" + year + "'")
			}

			if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
				bufferRes.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
			} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {

				// select company option (fomr sahring only)
				bufferRes.WriteString(" AND ifo.organization_id= " + companyID)
			}

			whreStrIF = bufferRes.String()
			getDataRespons, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whreStrIF, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//get hours

			var getDataHours []objects.GraficDataHours
			var getDataPeriode []objects.GraficDataPeriod

			if graficType == "submission" {
				var bufferHrs bytes.Buffer
				if periode != "" {
					bufferHrs.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
				}
				if year != "" {
					bufferHrs.WriteString(" AND TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'")
				}

				if organizationID >= 1 && iCompanyID <= 0 {
					bufferHrs.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
				} else if organizationID >= 1 && iCompanyID >= 1 {

					// select company option (fomr sahring only)
					bufferHrs.WriteString(" AND ifo.organization_id= " + companyID)
				}

				whreStrHrs := bufferHrs.String()

				getDataHours, err = ctr.inputForm.GetDataHours(formID, whreStrHrs)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if periode == "daily" {
					var whrDaily bytes.Buffer
					// whreDailySub := " where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
					whrDaily.WriteString("where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')")
					if organizationID >= 1 {
						whrDaily.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					}

					if organizationID >= 1 && iCompanyID <= 0 {
						whrDaily.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrDaily.WriteString(" AND ifo.organization_id= " + companyID)
					}

					whreDailySub := whrDaily.String()
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeDays(formID, whreDailySub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "monthly" {
					var whrMonthly bytes.Buffer
					whrMonthly.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")

					if organizationID >= 1 && iCompanyID <= 0 {
						whrMonthly.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrMonthly.WriteString(" AND ifo.organization_id= " + companyID)
					}

					whreMonthlySub := whrMonthly.String()
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeMonthly(formID, whreMonthlySub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "yearly" {

					var whrYearly bytes.Buffer
					whrYearly.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "' ")

					if organizationID >= 1 && iCompanyID <= 0 {
						whrYearly.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrYearly.WriteString(" AND ifo.organization_id= " + companyID)
					}

					whreYearSub := whrYearly.String()
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeYearly(formID, whreYearSub)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
				}

			} else if graficType == "responden" {

				var bufferHrs bytes.Buffer
				if periode != "" {
					bufferHrs.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")
				}
				if year != "" {
					bufferHrs.WriteString(" AND TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "'")
				}

				if organizationID >= 1 && iCompanyID <= 0 {
					bufferHrs.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
				} else if organizationID >= 1 && iCompanyID >= 1 {

					// select company option (fomr sahring only)
					bufferHrs.WriteString(" AND ifo.organization_id= " + companyID)
				}

				whreStrHrs := bufferHrs.String()

				getDataHours, err = ctr.inputForm.GetDataHoursUserResp(formID, whreStrHrs)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if periode == "daily" {

					var whrDailyResp bytes.Buffer
					whrDailyResp.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') ")

					if organizationID >= 1 && iCompanyID <= 0 {
						whrDailyResp.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrDailyResp.WriteString(" AND ifo.organization_id= " + companyID)
					}

					whreDailyResp := whrDailyResp.String()
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeDaysResp(formID, whreDailyResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "monthly" {

					// whreMonthlyResp := " where TO_CHAR(f.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') "
					var whrMonthlyResp bytes.Buffer
					whrMonthlyResp.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy-mm') = TO_CHAR(NOW()::date, 'yyyy-mm') ")

					if organizationID >= 1 && iCompanyID <= 0 {
						whrMonthlyResp.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrMonthlyResp.WriteString(" AND ifo.organization_id= " + companyID)
					}
					whreMonthlyResp := whrMonthlyResp.String()
					getDataPeriode, err = ctr.inputForm.GetDataPeriodeMonthlyResp(formID, whreMonthlyResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "yearly" {
					var whrYearlyResp bytes.Buffer
					whrYearlyResp.WriteString(" where TO_CHAR(f.created_at::date, 'yyyy') = '" + year + "' ")

					if organizationID >= 1 && iCompanyID <= 0 {
						whrYearlyResp.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && iCompanyID >= 1 {

						// select company option (fomr sahring only)
						whrYearlyResp.WriteString(" AND ifo.organization_id= " + companyID)
					}
					whreYearResp := whrYearlyResp.String()

					getDataPeriode, err = ctr.inputForm.GetDataPeriodeYearlyResp(formID, whreYearResp)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

				}
			}

			// var hoursResult []objects.GraficData
			// if len(getHours) > 0 {
			// 	for i := 0; i < len(getHours); i++ {
			// 		var ec objects.GraficData
			// 		ec.Field = getHours[i].Hours
			// 		ec.Value = "0"

			// 		hoursResult = append(hoursResult, ec)
			// 	}
			// }

			var res objects.FormGraficRes
			res.FormID = formID
			res.TotalResponden = len(getRespondens)
			res.TotalRespon = len(getDataRespons)
			res.ActiveHours = getDataHours
			res.ActivePeriod = getDataPeriode

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form id is not found",
			"data":    nil,
		})
		return
	}
}

func (ctr *inputFormController) DataFormResponList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	// getAuthorData, _ := ctr.userMod.GetUserRow(tables.Users{ID: userID})

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	companyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, _ := strconv.Atoi(companyID)
	isDeletedHide := c.Request.URL.Query().Get("is_deleted_hide")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--87sf79dsf--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"model":  "GetFormFieldNotParentRows ",
				"status": false,
				"error":  err.Error(),
			})
			return
		}

		if len(results) > 0 {

			var fieldLabels []objects.InputFormFields
			// row header ---------------------------------------------------------------------------
			for i := 0; i < len(results); i++ {

				var each objects.InputFormFields
				each.ID = results[i].ID
				each.Label = results[i].Label
				each.FieldTypeName = results[i].FieldTypeName

				fmt.Println("results[i].FieldTypeID ::------->", results[i].FieldTypeID, results[i].FieldTypeName, results[i].Label)
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					fmt.Println("len(tabDataRowHeader.TabDataHeader) :::", results[i].ID, len(tabDataRowHeader.TabDataHeader))
					var fieldLabel []objects.TabValueHeader
					for j := 0; j < len(tabDataRowHeader.TabDataHeader); j++ {
						var ech objects.TabValueHeader
						ech.Label = tabDataRowHeader.TabDataHeader[j].Value

						fmt.Println("-------------tab--------------->", tabDataRowHeader.TabDataHeader[j].Value)
						fieldLabel = append(fieldLabel, ech)
					}

					each.TabHeader = fieldLabel
				}

				fieldLabels = append(fieldLabels, each)
			}

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			var whereInForm tables.InputForms

			var whrInputForm bytes.Buffer
			var whreStr = ``

			// whrInputForm.WriteString(" if.deleted_at is null or if.deleted_at is not null")

			// get form detail
			getForm, err := ctr.formMod.GetFormRow(tables.Forms{ID: formID})
			if err != nil {
				fmt.Println("GetFormRow---------->>")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			// if iOptionCompanyID >= 1 {
			// 	whrInputForm.WriteString(" ifo.organization_id=" + optionCompanyID)
			// } else if organizationID != getForm.OrganizationID {
			// 	whrInputForm.WriteString(" ifo.organization_id=" + strconv.Itoa(organizationID))
			// } else {
			// 	whrInputForm.WriteString(" ifo.organization_id is not null")
			// }

			checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
			if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
				// company ID by TOKEN
				whrInputForm.WriteString(" ifo.organization_id= " + claims["organization_id"].(string))
			} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
				// select company option (form sharing only)
				whrInputForm.WriteString(" ifo.organization_id= " + strconv.Itoa(iCompanyID))
			}

			if searchKeyWord != "" {
				whrInputForm.WriteString(" u.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			if periodeStart != "" && periodeEnd == "" {
				whrInputForm.WriteString(" to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
			}

			if periodeEnd != "" && periodeStart == "" {
				whrInputForm.WriteString(" to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
			}

			if periodeStart != "" && periodeEnd != "" {
				whrInputForm.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
			}

			whreStr = whrInputForm.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getData, err := ctr.inputForm.GetInputFormUnscopedRows(formID, whereInForm, whreStr, paging)
			if err != nil {
				fmt.Println("err GetFormFieldRows--adae232--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"model":  "GetFormFieldRows",
					"status": false,
					"error":  err.Error(),
				})
				return
			}
			fmt.Println("isDeletedHide >>>>>>>>>>>>>>>>", isDeletedHide)
			getAllData, err := ctr.inputForm.GetInputFormUnscopedRows(formID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows-3uju5o63---", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"model":  "GetFormFieldRows all",
					"status": false,
					"error":  err.Error(),
				})
				return
			}
			count := 0
			var dataRows []objects.InputFieldData
			for i := 0; i < len(getData); i++ {

				deletedValue, _ := getData[i].DeletedAt.Value()
				if deletedValue == nil {

					var each objects.InputFieldData
					each.ID = getData[i].ID
					each.UserID = getData[i].UserID
					each.UserName = getData[i].UserName
					each.Phone = getData[i].Phone

					if getData[i].UpdatedCount == 0 {
						each.StatusData = "Terkirim"
					} else {
						each.StatusData = "Diubah " + strconv.Itoa(getData[i].UpdatedCount) + " kali\n" + getData[i].UpdatedAt.Format("2006-01-02 15:04")
					}
					each.CreatedAt = getData[i].CreatedAt.Format("2006-01-02 15:04")
					each.UpdatedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					// rows data --------------------------------------------
					var resultData []objects.InputData
					for j := 0; j < len(results); j++ {

						fieldID := strconv.Itoa(results[j].ID)
						fieldTypeID := results[j].FieldTypeID

						var ec objects.InputData
						ec.FieldID = fieldID
						ec.FieldTypeID = fieldTypeID

						//data user here -------------------------------------------------------------
						var fields tables.InputForms
						fields.ID = getData[i].ID
						fieldStrings := "coalesce(f" + fieldID + ",'') as f"

						inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, "")
						if err != nil {
							fmt.Println("err GetFormFieldRows--jgjh755W--", err)
							c.JSON(http.StatusBadRequest, gin.H{
								"model":  "GetInputDataRows",
								"status": false,
								"error":  err.Error(),
							})
							return
						}

						ec.Value = inputData[0][0]

						if fieldTypeID == 22 {
							var fieldData []objects.TabValueAnswer
							obj := inputData[0][0]
							json.Unmarshal([]byte(obj), &fieldData)

							ec.TabValue = fieldData
						}

						fmt.Println("+++++++++++++>>", inputData[0][0])
						if fieldTypeID == 10 || fieldTypeID == 18 || fieldTypeID == 20 {
							var value = inputData[0][0]

							if value != "" && len(value) < 35 {
								words := strings.Split(value, "/")
								lastWord := words[len(words)-1]

								// if i == 5 {
								// 	time.Sleep(3 * time.Second)
								// }

								// fmt.Println("lastWord ===============>", count, " --- ", fieldID, " :: ", lastWord)

								linkReal, err := ctr.shrtMod.GetLinkReal(lastWord)
								if err != nil {
									// Handle the error.
									fmt.Println(err)
									c.JSON(http.StatusBadRequest, gin.H{
										"model":  "GetLinkReal",
										"status": false,
										"error":  err.Error(),
									})
									return
								}
								count += 1
								// if count == 59 {
								// 	time.Sleep(3 * time.Second)
								// } else {
								// 	time.Sleep(500 * time.Millisecond)
								// }

								fmt.Println("linkReal ------------------------->", linkReal, count)
								ec.Value = linkReal.Data.URL
							}

							// ec.Value = value

						}

						resultData = append(resultData, ec)
					}

					each.InputData = resultData
					dataRows = append(dataRows, each)
				} else if deletedValue != nil && isDeletedHide == "false" {

					getDeleteUserData, _ := ctr.userMod.GetUserRow(tables.Users{ID: getData[i].DeletedBy})

					var each objects.InputFieldData
					each.ID = getData[i].ID
					each.StatusData = "Dihapus"
					each.Note = "Data telah dihapus oleh (" + getDeleteUserData.Name + ")\n" + getData[i].DeletedAt.Time.Format("2006-01-02 15:04")
					each.DeletedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					// rows data --------------------------------------------
					var resultData []objects.InputData
					for j := 0; j < len(results); j++ {

						fieldID := strconv.Itoa(results[j].ID)
						fieldTypeID := results[j].FieldTypeID

						var ec objects.InputData
						ec.FieldID = fieldID
						ec.FieldTypeID = fieldTypeID
						ec.Value = "Data dihapus"

						resultData = append(resultData, ec)
					}

					each.InputData = resultData
					dataRows = append(dataRows, each)
				}
			}

			// form detail data ------------------------------------------------------------------------------
			var res objects.InputFormRes
			res.FormID = formID
			res.FormName = getForm.Name
			res.FormDescription = getForm.Description
			res.PeriodStartDate = getForm.PeriodStartDate
			res.PeriodEndDate = getForm.PeriodEndDate
			res.FieldLabel = fieldLabels
			res.FieldData = dataRows

			totalPage := 0
			if limit > 0 {

				totalPage = len(getAllData) / limit
				if (len(getAllData) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var pagingRes objects.DataRows
			pagingRes.TotalRows = len(getAllData)
			pagingRes.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  pagingRes,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *inputFormController) DataFormCompanyResponList(c *gin.Context) {

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}

	selectedCompanyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err.Error(),
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--34543k5jkj--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"model":  "GetFormFieldNotParentRows ",
				"status": false,
				"error":  err.Error(),
			})
			return
		}

		if len(results) > 0 {

			var fieldLabels []objects.InputFormFields

			// row header ---------------------------------------------------------------------------
			for i := 0; i < len(results); i++ {

				var each objects.InputFormFields
				each.ID = results[i].ID
				each.Label = results[i].Label
				each.FieldTypeName = results[i].FieldTypeName

				fmt.Println("results[i].FieldTypeID ::------->", results[i].FieldTypeID, results[i].FieldTypeName, results[i].Label)
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					fmt.Println("len(tabDataRowHeader.TabDataHeader) :::", results[i].ID, len(tabDataRowHeader.TabDataHeader))
					var fieldLabel []objects.TabValueHeader
					for j := 0; j < len(tabDataRowHeader.TabDataHeader); j++ {
						var ech objects.TabValueHeader
						ech.Label = tabDataRowHeader.TabDataHeader[j].Value

						fmt.Println("-------------tab--------------->", tabDataRowHeader.TabDataHeader[j].Value)
						fieldLabel = append(fieldLabel, ech)
					}

					each.TabHeader = fieldLabel
				}

				fieldLabels = append(fieldLabels, each)
			}

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			var whereInForm tables.InputFormJoinOrganizations

			var buffer bytes.Buffer
			var whreStr = ``
			if searchKeyWord != "" {
				buffer.WriteString(" u.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			if periodeStart != "" && periodeEnd == "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
			}

			if periodeEnd != "" && periodeStart == "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
			}

			if periodeStart != "" && periodeEnd != "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
			}

			if selectedCompanyID >= 1 {
				buffer.WriteString(" AND ifo.organization_id =" + c.Param("company_id"))
			} else {
				buffer.WriteString(" AND ifo.organization_id is null")
			}

			whreStr = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whreStr, paging)
			if err != nil {
				fmt.Println("err GetFormFieldRows--kj46jk--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"model":  "GetFormFieldRows",
					"status": false,
					"error":  err.Error(),
				})
				return
			}

			getAllData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows--846i46j--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"model":  "GetFormFieldRows all",
					"status": false,
					"error":  err.Error(),
				})
				return
			}
			// count := 0
			var dataRows []objects.InputFieldData
			for i := 0; i < len(getData); i++ {

				var each objects.InputFieldData
				each.ID = getData[i].ID
				each.UserID = getData[i].UserID
				each.UserName = getData[i].UserName
				each.Phone = getData[i].Phone
				each.CreatedAt = getData[i].CreatedAt.Format("2006-01-02 15:04")

				// rows data --------------------------------------------
				var resultData []objects.InputData
				for j := 0; j < len(results); j++ {

					fieldID := strconv.Itoa(results[j].ID)
					fieldTypeID := results[j].FieldTypeID

					var ec objects.InputData
					ec.FieldID = fieldID
					ec.FieldTypeID = fieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[i].ID
					fieldStrings := "coalesce(f" + fieldID + ",'') as f"

					var whrOrg string
					if selectedCompanyID >= 1 {
						whrOrg = "ifa.organization_id =" + c.Param("company_id")
					} else {
						whrOrg = "ifa.organization_id is null"

					}

					inputData, err := ctr.inputForm.GetInputDataOrganizationRows(formID, fieldStrings, fields, whrOrg)
					if err != nil {
						fmt.Println("err GetFormFieldRows--k34k5jm43l--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"model":  "GetInputDataRows",
							"status": false,
							"error":  err.Error(),
						})
						return
					}

					ec.Value = inputData[0][0]

					if fieldTypeID == 22 {
						var fieldData []objects.TabValueAnswer
						obj := inputData[0][0]
						json.Unmarshal([]byte(obj), &fieldData)

						ec.TabValue = fieldData
					}

					/*
						fmt.Println("+++++++++++++>>", inputData[0][0])
						if fieldTypeID == 10 || fieldTypeID == 18 || fieldTypeID == 20 {
							var value = inputData[0][0]

							if value != "" && len(value) < 35 {
								words := strings.Split(value, "/")
								lastWord := words[len(words)-1]

								// if i == 5 {
								// 	time.Sleep(3 * time.Second)
								// }

								// fmt.Println("lastWord ===============>", count, " --- ", fieldID, " :: ", lastWord)

								linkReal, err := ctr.shrtMod.GetLinkReal(lastWord)
								if err != nil {
									// Handle the error.
									fmt.Println(err)
									c.JSON(http.StatusBadRequest, gin.H{
										"model":  "GetLinkReal",
										"status": false,
										"error":  err.Error(),
									})
									return
								}
								count += 1
								// if count == 59 {
								// 	time.Sleep(3 * time.Second)
								// } else {
								// 	time.Sleep(500 * time.Millisecond)
								// }

								fmt.Println("linkReal ------------------------->", linkReal, count)
								ec.Value = linkReal.Data.URL
							}

							// ec.Value = value

						}*/

					resultData = append(resultData, ec)
				}

				each.InputData = resultData
				dataRows = append(dataRows, each)
			}

			// form detail data ------------------------------------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, err := ctr.formMod.GetFormRow(fieldForm)

			var res objects.InputFormRes
			res.FormID = formID
			res.FormName = getForm.Name
			res.FormDescription = getForm.Description
			res.PeriodStartDate = getForm.PeriodStartDate
			res.PeriodEndDate = getForm.PeriodEndDate
			res.FieldLabel = fieldLabels
			res.FieldData = dataRows

			totalPage := 0
			if limit > 0 {

				totalPage = len(getAllData) / limit
				if (len(getAllData) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var pagingRes objects.DataRows
			pagingRes.TotalRows = len(getAllData)
			pagingRes.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  pagingRes,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *inputFormController) DataFormDetailDownload(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFieldRows--ekjtl56kj--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [697]string{"F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx

			xlsx := excelize.NewFile() // xlsx
			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

			xlsx.SetCellValue(sheet1Name, "A1", "NAME") // xlsx
			xlsx.SetCellValue(sheet1Name, "B1", "DATE") // xlsx
			xlsx.SetCellValue(sheet1Name, "C1", "HP")   // xlsx
			xlsx.SetCellValue(sheet1Name, "D1", "LAT")  // xlsx
			xlsx.SetCellValue(sheet1Name, "E1", "LNG")  // xlsx

			// row header --------------------------
			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {

				fmt.Println(":: ID ----", results[i].ID)
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						// fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
						xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[nextIndex1+1]+"1", results[i].Label)
						nextIndex1++
					} else {
						// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[index]+"1", results[i].Label)
						index++
					}
				}

			}

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			if periodeStart != "" && periodeEnd == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' "
			}

			if periodeEnd != "" && periodeStart == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  "
			}

			if periodeStart != "" && periodeEnd != "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' "
			}

			var whereInForm tables.InputForms

			getData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows--l3k5lk5l--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			no := 1
			for a := 0; a < len(getData); a++ {

				// formDataDate := getData[a].CreatedAt.Format("2006-01-02")
				// var whr tables.Attendances
				// whrStr := "to_char(attendances.created_at, 'yyyy-mm-dd') =  '" + formDataDate + "'"
				// getAtt, _ := ctr.attMod.GetAttendanceObjRow(whr, whrStr)

				// var lat float64
				// var long float64

				// if getAtt.ID > 0 {
				// 	lat = getAtt.Latitude
				// 	long = getAtt.Longitude
				// }

				row := strconv.Itoa(no + 1)
				xlsx.SetCellValue(sheet1Name, "A"+row, getData[a].UserName)  // xlsx
				xlsx.SetCellValue(sheet1Name, "B"+row, getData[a].CreatedAt) // xlsx
				xlsx.SetCellValue(sheet1Name, "C"+row, getData[a].Phone)     // xlsx
				xlsx.SetCellValue(sheet1Name, "D"+row, getData[a].Latitude)  // xlsx
				xlsx.SetCellValue(sheet1Name, "E"+row, getData[a].Longitude) // xlsx

				// rows data --------------------------------------------
				indexData := 0
				nextIndexData := 0

				fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
				for j := 0; j < len(results); j++ {

					fieldID := strconv.Itoa(results[j].ID)
					fieldTypeID := results[j].FieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[a].ID

					var bufferInputData bytes.Buffer
					fieldStrings := "coalesce(f" + fieldID + ",'') as f"

					if periodeStart != "" && periodeEnd == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					fmt.Println(periodeStart, periodeEnd)
					whereStr := bufferInputData.String()
					inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, whereStr)
					if err != nil {
						fmt.Println("err GetFormFieldRows--3k4l3lll;--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					if len(inputData) > 0 {
						obj := inputData[0][0]

						if fieldTypeID != 22 { // (22 = tabulasi)

							xlsx.SetCellValue(sheet1Name, col[indexData]+row, obj)

						} else {

							var fieldData []objects.TabValueAnswer
							json.Unmarshal([]byte(obj), &fieldData)

							nextIndexDataCol := 0
							for m := 0; m < len(fieldData); m++ {

								xlsx.SetCellValue(sheet1Name, col[m+indexData]+row, fieldData[m].Answer)
								nextIndexDataCol = m + indexData

							}
							indexData = nextIndexDataCol
						}
						indexData++
					}

					// if len(inputData) > 0 {
					//

					// 	// if fieldTypeID == 22 {
					// 	// 	var fieldData []objects.TabValueAnswer
					// 	// 	json.Unmarshal([]byte(obj), &fieldData)

					// 	// 	nextIndexDataCol := 0
					// 	// 	for m := 0; m < len(fieldData); m++ {

					// 	// 		xlsx.SetCellValue(sheet1Name, col[m]+row, fieldData[m].Answer)
					// 	// 		nextIndexDataCol = indexData + m

					// 	// 		fmt.Println("SetCellValue, --------IF-----inpputdata-------->> ::", col[m]+row, fieldData[m].Answer)
					// 	// 	}
					// 	// 	nextIndexData = nextIndexDataCol

					// 	// } else {

					// 	// 	if nextIndexData > 0 {

					// 	// 		fmt.Println("formID, --------IF-----2222-------->> ::", col[nextIndexData+1]+row, obj)
					// 	// 		xlsx.SetCellValue(sheet1Name, col[nextIndexData+1]+row, obj)
					// 	// 	} else {
					// 	// 		fmt.Println("formID, --------ELSE-----2222-------->> ::", col[indexData]+row, obj)
					// 	// 		xlsx.SetCellValue(sheet1Name, col[indexData]+row, obj)

					// 	// 	}
					// 	// 	nextIndexData++
					// 	// 	indexData++
					// 	// }
					// 	fmt.Println(fieldTypeID, nextIndexData, ":::", indexData, row, col[indexData]+row, obj)
					// 	xlsx.SetCellValue(sheet1Name, col[indexData]+row, obj)
					// } else {
					// 	c.JSON(http.StatusBadRequest, gin.H{
					// 		"status":  false,
					// 		"message": "Data is not available",
					// 		"data":    nil,
					// 	})
					// 	return
					// }

				}

				no++
			}

			// err1 := xlsx.AutoFilter(sheet1Name, "A1", "C1", "")
			// if err1 != nil {
			// 	log.Fatal("ERROR", err1.Error())
			// }

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}

}

func (ctr *inputFormController) DataFormDetail2Download(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--l234klk34lk--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [696]string{"G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx

			xlsx := excelize.NewFile() // xlsx
			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

			xlsx.SetCellValue(sheet1Name, "A1", "No.")  // xlsx
			xlsx.SetCellValue(sheet1Name, "B1", "NAME") // xlsx
			xlsx.SetCellValue(sheet1Name, "C1", "DATE") // xlsx
			xlsx.SetCellValue(sheet1Name, "D1", "HP")   // xlsx
			xlsx.SetCellValue(sheet1Name, "E1", "LAT")  // xlsx
			xlsx.SetCellValue(sheet1Name, "F1", "LNG")  // xlsx

			// row header --------------------------
			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {

				fmt.Println(":: ID ----", results[i].ID)
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						// fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
						xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[nextIndex1+1]+"1", results[i].Label)
						nextIndex1++
					} else {
						// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[index]+"1", results[i].Label)
						index++
					}
				}

			}
			// end row header ----------------------

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			if periodeStart != "" && periodeEnd == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' "
			}

			if periodeEnd != "" && periodeStart == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  "
			}

			if periodeStart != "" && periodeEnd != "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' "
			}

			var whereInForm tables.InputForms

			getData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows--k5jlm435kjl35j--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			no := 1
			for a := 0; a < len(getData); a++ {

				row := strconv.Itoa(no + 1)
				xlsx.SetCellValue(sheet1Name, "A"+row, no)                   // xlsx
				xlsx.SetCellValue(sheet1Name, "B"+row, getData[a].UserName)  // xlsx
				xlsx.SetCellValue(sheet1Name, "C"+row, getData[a].CreatedAt) // xlsx
				xlsx.SetCellValue(sheet1Name, "D"+row, getData[a].Phone)     // xlsx
				xlsx.SetCellValue(sheet1Name, "E"+row, getData[a].Latitude)  // xlsx
				xlsx.SetCellValue(sheet1Name, "F"+row, getData[a].Longitude) // xlsx

				// rows data --------------------------------------------
				indexData := 0
				nextIndexData := 0

				fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
				for j := 0; j < len(results); j++ {

					fieldID := strconv.Itoa(results[j].ID)
					fieldTypeID := results[j].FieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[a].ID

					var bufferInputData bytes.Buffer
					fieldStrings := "coalesce(f" + fieldID + ",'') as f"

					if periodeStart != "" && periodeEnd == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					fmt.Println(periodeStart, periodeEnd)
					whereStr := bufferInputData.String()
					inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, whereStr)
					if err != nil {
						fmt.Println("err GetFormFieldRows--q,eq,em--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					if len(inputData) > 0 {
						obj := inputData[0][0]

						if fieldTypeID != 22 { // (22 = tabulasi)
							rowData := strings.Replace(obj, "\n", " ", 100)

							xlsx.SetCellValue(sheet1Name, col[indexData]+row, rowData)
							indexData++
						} else {

							var fieldData []objects.TabValueAnswer
							json.Unmarshal([]byte(obj), &fieldData)

							nextIndexDataCol := 0
							if len(fieldData) > 0 {
								for m := 0; m < len(fieldData); m++ {

									rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

									xlsx.SetCellValue(sheet1Name, col[m+indexData]+row, rowData)
									nextIndexDataCol = m + indexData

								}
								indexData = nextIndexDataCol
								indexData++
							}
						}

					}
				}

				no++
			}
			// end row data --------------------------------------------------

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			formName = strings.Replace(formName, "/", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}

}

func (ctr *inputFormController) DataFormDetail3Download(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("organizationID :::", organizationID)
	}

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")
	selectedCompanyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, _ := strconv.Atoi(selectedCompanyID)

	fmt.Println("iCompanyID", iCompanyID)

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--jrkemwr234--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": true,
				"error":  err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [696]string{"G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx
			xlsx := excelize.NewFile()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     // xlsx

			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

			// row header --------------------------
			xlsx.SetCellValue(sheet1Name, "A1", "No.")  // xlsx
			xlsx.SetCellValue(sheet1Name, "B1", "NAME") // xlsx
			xlsx.SetCellValue(sheet1Name, "C1", "DATE") // xlsx
			xlsx.SetCellValue(sheet1Name, "D1", "HP")   // xlsx
			xlsx.SetCellValue(sheet1Name, "E1", "LAT")  // xlsx
			xlsx.SetCellValue(sheet1Name, "F1", "LNG")  // xlsx

			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {

				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						// fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
						xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[nextIndex1+1]+"1", results[i].Label)
						nextIndex1++
					} else {
						// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[index]+"1", results[i].Label)
						index++
					}
				}
			}

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			var bufferInputForm bytes.Buffer

			// selecting company id
			checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
			if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
				bufferInputForm.WriteString(" ifo.organization_id= " + claims["organization_id"].(string))

			} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
				bufferInputForm.WriteString(" ifo.organization_id= " + strconv.Itoa(iCompanyID))

			}

			if periodeStart != "" && periodeEnd == "" {
				bufferInputForm.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
			}

			if periodeEnd != "" && periodeStart == "" {
				bufferInputForm.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
			}

			if periodeStart != "" && periodeEnd != "" {
				bufferInputForm.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
			}
			whereData = bufferInputForm.String()

			var whereInForm tables.InputFormJoinOrganizations
			getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows--kjkj6k7j--", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			no := 1
			for a := 0; a < len(getData); a++ {

				row := strconv.Itoa(no + 1)
				xlsx.SetCellValue(sheet1Name, "A"+row, no)                   // xlsx
				xlsx.SetCellValue(sheet1Name, "B"+row, getData[a].UserName)  // xlsx
				xlsx.SetCellValue(sheet1Name, "C"+row, getData[a].CreatedAt) // xlsx
				xlsx.SetCellValue(sheet1Name, "D"+row, getData[a].Phone)     // xlsx
				xlsx.SetCellValue(sheet1Name, "E"+row, getData[a].Latitude)  // xlsx
				xlsx.SetCellValue(sheet1Name, "F"+row, getData[a].Longitude) // xlsx

				// rows data --------------------------------------------
				indexData := 0
				nextIndexData := 0

				fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
				for j := 0; j < len(results); j++ {

					fieldID := strconv.Itoa(results[j].ID)
					fieldTypeID := results[j].FieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[a].ID

					var bufferInputData bytes.Buffer
					selectFieldStr := "coalesce(f" + fieldID + ",'') as f"

					if periodeStart != "" && periodeEnd == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}

					whereStr := bufferInputData.String()
					inputData, err := ctr.inputForm.GetInputDataOrganizationRows(formID, selectFieldStr, fields, whereStr)
					if err != nil {
						fmt.Println("err GetFormFieldRows--k2j3lqmlkl--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					if len(inputData) >= 1 {
						obj := inputData[0][0]

						if fieldTypeID != 22 { // (22 = tabulasi)
							rowData := strings.Replace(obj, "\n", " ", 100)

							xlsx.SetCellValue(sheet1Name, col[indexData]+row, rowData)
							indexData++
						} else {

							var fieldData []objects.TabValueAnswer
							json.Unmarshal([]byte(obj), &fieldData)

							nextIndexDataCol := 0
							if len(fieldData) > 0 {
								for m := 0; m < len(fieldData); m++ {

									rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

									xlsx.SetCellValue(sheet1Name, col[m+indexData]+row, rowData)
									nextIndexDataCol = m + indexData

								}
								indexData = nextIndexDataCol
								indexData++
							}
						}

					}
				}

				no++
			}
			// end row data --------------

			// check multycompany here -----------------------------------------------------------------------------------------
			var wf tables.InputFormOrganizations
			wf.FormID = formID

			if iCompanyID >= 1 {
				wf.OrganizationID = iCompanyID

				// delete sheet index 1
				// xlsx.DeleteSheet(sheet1Name)
				fmt.Println("xlsx.DeleteSheet(sheet1Name) ------------------------------->", sheet1Name)
			}

			checkOrgInputForm, err := ctr.inputForm.GetOrganizationInputForm(wf, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			fmt.Println("checkOrgInputForm ::=======>", len(checkOrgInputForm), selectedCompanyID)

			if len(checkOrgInputForm) >= 1 {

				for o := 0; o < len(checkOrgInputForm); o++ {

					sheet4Company := "Export-Form-Data-" + checkOrgInputForm[o].OrganizationName
					xlsx.NewSheet(sheet4Company)

					// row header
					xlsx.SetCellValue(sheet4Company, "A1", "No.")  // xlsx
					xlsx.SetCellValue(sheet4Company, "B1", "NAME") // xlsx
					xlsx.SetCellValue(sheet4Company, "C1", "DATE") // xlsx
					xlsx.SetCellValue(sheet4Company, "D1", "HP")   // xlsx
					xlsx.SetCellValue(sheet4Company, "E1", "LAT")  // xlsx
					xlsx.SetCellValue(sheet4Company, "F1", "LNG")  // xlsx

					index := 0
					nextIndex1 := 0
					for i := 0; i < len(results); i++ {

						fmt.Println(":: ID ----", results[i].ID)
						if results[i].FieldTypeID == 22 {

							var tabDataRowHeader objects.TabDataRowHeader
							objHeader := results[i].Option
							json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

							if nextIndex1 > 0 {
								index = nextIndex1 + 1
							}
							nextIndexCol := 0
							for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

								// fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
								xlsx.SetCellValue(sheet4Company, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
								nextIndexCol = index + k
							}
							nextIndex1 = nextIndexCol

						} else {
							if nextIndex1 > 0 {
								// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[nextIndex1+1]+"1", results[i].Label)
								nextIndex1++
							} else {
								// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[index]+"1", results[i].Label)
								index++
							}
						}

					}
					// end row header

					// row data --------------------
					whereData := ""
					var bufferInputFormOrg bytes.Buffer

					checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
						// company ID by TOKEN
						bufferInputFormOrg.WriteString(" ifo.organization_id =" + claims["organization_id"].(string))
					} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
						// select company option (form sharing only)
						bufferInputFormOrg.WriteString(" ifo.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
					}

					if periodeStart != "" && periodeEnd == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					whereData = bufferInputFormOrg.String()

					var whereInForm tables.InputFormJoinOrganizations
					// whereInForm.OrganizationID, _ = strconv.Atoi(selectedCompanyID)
					getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whereData, objects.Paging{})
					if err != nil {
						fmt.Println("err GetFormFieldRows--q,wmeqw,em--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					no := 1
					for a := 0; a < len(getData); a++ {

						row := strconv.Itoa(no + 1)
						xlsx.SetCellValue(sheet4Company, "A"+row, no)                   // xlsx
						xlsx.SetCellValue(sheet4Company, "B"+row, getData[a].UserName)  // xlsx
						xlsx.SetCellValue(sheet4Company, "C"+row, getData[a].CreatedAt) // xlsx
						xlsx.SetCellValue(sheet4Company, "D"+row, getData[a].Phone)     // xlsx
						xlsx.SetCellValue(sheet4Company, "E"+row, getData[a].Latitude)  // xlsx
						xlsx.SetCellValue(sheet4Company, "F"+row, getData[a].Longitude) // xlsx

						// rows data --------------------------------------------
						indexData := 0
						nextIndexData := 0

						fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
						for j := 0; j < len(results); j++ {

							fieldID := strconv.Itoa(results[j].ID)
							fieldTypeID := results[j].FieldTypeID

							//data user here -------------------------------------------------------------
							var fields tables.InputForms
							fields.ID = getData[a].ID

							var bufferInputData bytes.Buffer
							fieldStrings := "coalesce(f" + fieldID + ",'') as f"

							bufferInputData.WriteString(" ifo.organization_id= " + strconv.Itoa(checkOrgInputForm[o].OrganizationID))

							if periodeStart != "" && periodeEnd == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
							}

							if periodeEnd != "" && periodeStart == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
							}

							if periodeStart != "" && periodeEnd != "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
							}

							// if periodeStart != "" || periodeEnd != "" {

							// 	bufferInputData.WriteString(" AND ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// } else {
							// 	bufferInputData.WriteString(" ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// }

							whereStr := bufferInputData.String()
							inputData, err := ctr.inputForm.GetInputDataOrganizationRows(formID, fieldStrings, fields, whereStr)
							if err != nil {
								fmt.Println("err GetFormFieldRows--jk3k24hkj5--", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							if len(inputData) > 0 {
								obj := inputData[0][0]

								if fieldTypeID != 22 { // (22 = tabulasi)
									rowData := strings.Replace(obj, "\n", " ", 100)

									xlsx.SetCellValue(sheet4Company, col[indexData]+row, rowData)
									indexData++
								} else {

									var fieldData []objects.TabValueAnswer
									json.Unmarshal([]byte(obj), &fieldData)

									nextIndexDataCol := 0
									if len(fieldData) > 0 {
										for m := 0; m < len(fieldData); m++ {

											rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

											xlsx.SetCellValue(sheet4Company, col[m+indexData]+row, rowData)
											nextIndexDataCol = m + indexData

										}
										indexData = nextIndexDataCol
										indexData++
									}
								}

							}
						}

						no++
					}
					// end row data --------------------------------------------------
				} // for organization
			}

			// delete sheet index 1
			if iCompanyID >= 1 {
				xlsx.DeleteSheet(sheet1Name)
				fmt.Println("xlsx.DeleteSheet(sheet1Name) -----2-------------------------->", sheet1Name)
			}
			// END MUlty company --------------------------------------------

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			formName = strings.Replace(formName, "/", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data field is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}

}

// multy organization
func (ctr *inputFormController) DataFormDetail4DownloadBackup(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("organizationID :::", organizationID)
	}

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")
	selectedCompanyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, _ := strconv.Atoi(selectedCompanyID)

	fmt.Println("iCompanyID", iCompanyID)

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--23jjk2jh4kj--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": true,
				"error":  err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [696]string{"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx
			xlsx := excelize.NewFile()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                // xlsx

			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
			// end row data --------------

			// check multycompany here -----------------------------------------------------------------------------------------
			var wf tables.InputFormOrganizations
			wf.FormID = formID

			if iCompanyID >= 1 {
				wf.OrganizationID = iCompanyID

				// delete sheet index 1
				// xlsx.DeleteSheet(sheet1Name)
				// fmt.Println("xlsx.DeleteSheet(sheet1Name) ------------------------------->", sheet1Name)
			}

			checkOrgInputForm, err := ctr.inputForm.GetOrganizationInputForm(wf, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if len(checkOrgInputForm) >= 1 {

				for o := 0; o < len(checkOrgInputForm); o++ {

					sheet4Company := "Export-Form-Data-" + checkOrgInputForm[o].OrganizationName
					xlsx.NewSheet(sheet4Company)

					// row header
					xlsx.SetCellValue(sheet4Company, "A1", "No.")         // xlsx
					xlsx.SetCellValue(sheet4Company, "B1", "STATUS DATA") // xlsx
					xlsx.SetCellValue(sheet4Company, "C1", "NAME")        // xlsx
					xlsx.SetCellValue(sheet4Company, "D1", "DATE")        // xlsx
					xlsx.SetCellValue(sheet4Company, "E1", "HP")          // xlsx
					xlsx.SetCellValue(sheet4Company, "F1", "LAT")         // xlsx
					xlsx.SetCellValue(sheet4Company, "G1", "LNG")         // xlsx

					index := 0
					nextIndex1 := 0
					for i := 0; i < len(results); i++ {

						fmt.Println(":: ID ----", results[i].ID)
						if results[i].FieldTypeID == 22 {

							var tabDataRowHeader objects.TabDataRowHeader
							objHeader := results[i].Option
							json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

							if nextIndex1 > 0 {
								index = nextIndex1 + 1
							}
							nextIndexCol := 0
							for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

								// fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
								xlsx.SetCellValue(sheet4Company, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
								nextIndexCol = index + k
							}
							nextIndex1 = nextIndexCol

						} else {
							if nextIndex1 > 0 {
								// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[nextIndex1+1]+"1", results[i].Label)
								nextIndex1++
							} else {
								// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[index]+"1", results[i].Label)
								index++
							}
						}

					}

					style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#e4e4e4"],"pattern":1}}`)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadGateway, gin.H{
							"status":  false,
							"message": err.Error(),
						})
						return
					}
					xlsx.SetCellStyle(sheet4Company, "A1", col[nextIndex1+1]+"1", style)
					xlsx.SetColWidth(sheet4Company, "B", col[index], 15)
					// end row header

					// row data --------------------
					whereData := ""
					var bufferInputFormOrg bytes.Buffer

					// checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					// if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
					// 	// company ID by TOKEN
					// 	bufferInputFormOrg.WriteString(" ifo.organization_id =" + claims["organization_id"].(string))
					// } else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
					// 	// select company option (form sharing only)
					// 	bufferInputFormOrg.WriteString(" ifo.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
					// }

					bufferInputFormOrg.WriteString(" ifo.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))

					if periodeStart != "" && periodeEnd == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					whereData = bufferInputFormOrg.String()

					var whereInForm tables.InputFormJoinOrganizations
					// whereInForm.OrganizationID, _ = strconv.Atoi(selectedCompanyID)
					getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whereData, objects.Paging{})
					if err != nil {
						fmt.Println("err GetFormFieldRows--m42,23m4m,--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					no := 1
					for a := 0; a < len(getData); a++ {

						// if a <= 10 {

						var statusData string
						if getData[a].UpdatedCount == 0 {
							statusData = "Terkirim"
						} else {
							statusData = "Diubah " + strconv.Itoa(getData[a].UpdatedCount) + " kali (" + getData[a].UpdatedAt.Format("2006-01-02 15:04") + ")"
						}

						row := strconv.Itoa(no + 1)
						xlsx.SetCellValue(sheet4Company, "A"+row, no)                   // xlsx
						xlsx.SetCellValue(sheet4Company, "B"+row, statusData)           // xlsx
						xlsx.SetCellValue(sheet4Company, "C"+row, getData[a].UserName)  // xlsx
						xlsx.SetCellValue(sheet4Company, "D"+row, getData[a].CreatedAt) // xlsx
						xlsx.SetCellValue(sheet4Company, "E"+row, getData[a].Phone)     // xlsx
						xlsx.SetCellValue(sheet4Company, "F"+row, getData[a].Latitude)  // xlsx
						xlsx.SetCellValue(sheet4Company, "G"+row, getData[a].Longitude) // xlsx

						// rows data --------------------------------------------
						indexData := 0
						nextIndexData := 0

						fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
						for j := 0; j < len(results); j++ {

							fieldID := strconv.Itoa(results[j].ID)
							fieldTypeID := results[j].FieldTypeID

							//data user here -------------------------------------------------------------
							var fields tables.InputForms
							fields.ID = getData[a].ID

							var bufferInputData bytes.Buffer
							fieldStrings := "coalesce(f" + fieldID + ",'') as f"

							bufferInputData.WriteString(" ifo.organization_id= " + strconv.Itoa(checkOrgInputForm[o].OrganizationID))

							if periodeStart != "" && periodeEnd == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
							}

							if periodeEnd != "" && periodeStart == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
							}

							if periodeStart != "" && periodeEnd != "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
							}

							// if periodeStart != "" || periodeEnd != "" {

							// 	bufferInputData.WriteString(" AND ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// } else {
							// 	bufferInputData.WriteString(" ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// }

							whereStr := bufferInputData.String()
							inputData, err := ctr.inputForm.GetInputDataOrganizationRows(formID, fieldStrings, fields, whereStr)
							if err != nil {
								fmt.Println("err GetFormFieldRows--wlerjlwekrj--", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							if len(inputData) > 0 {
								obj := inputData[0][0]

								if fieldTypeID != 22 { // (22 = tabulasi)
									rowData := strings.Replace(obj, "\n", " ", 100)

									if fieldTypeID == 3 || fieldTypeID == 21 {

										var dataOption objects.DataOption
										json.Unmarshal([]byte(rowData), &dataOption)

										if len(dataOption.Data) >= 1 {

											optVal := make([]string, len(dataOption.Data))
											for k, v := range dataOption.Data {
												optVal[k] = v.Value
											}

											rowData = strings.Join(optVal, ",")
										}

										// os.Exit(0)
									}
									xlsx.SetCellValue(sheet4Company, col[indexData]+row, rowData)

									fmt.Println("-----reguler XXX-----", col[indexData]+row, "--index-", indexData, row, rowData)

									// indexData++
								}
								// } else {

								// 	nextIndexDataCol := 0

								// 	var tabDataRowHeader objects.TabDataRowHeader
								// 	objHeader := results[j].Option
								// 	json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

								// 	for p := 0; p < len(tabDataRowHeader.TabDataHeader); p++ {

								// 		fmt.Println("-----tabulasi-----", p)
								// 		var fieldData []objects.TabValueAnswer
								// 		json.Unmarshal([]byte(obj), &fieldData)

								// 		if len(fieldData) >= 1 {

								// 			for m := 0; m < len(fieldData); m++ {

								// 				if p == m {
								// 					fmt.Println("-----tabulasi XXX-----", p, "----m--", m, "--index-", indexData, col[m]+row)
								// 					rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

								// 					xlsx.SetCellValue(sheet4Company, col[m]+row, rowData)
								// 					nextIndexDataCol = m + indexData
								// 				}

								// 			}

								// 		}
								// 	}
								// 	indexData = nextIndexDataCol
								// 	// indexData++
								// 	// os.Exit(0)

								// 	// -------------------------

								// 	// var tabDataRowHeader objects.TabDataRowHeader
								// 	// objHeader := results[j].Option
								// 	// json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

								// 	// for p := 0; p < len(tabDataRowHeader.TabDataHeader); p++ {

								// 	// 	for m := 0; m < len(fieldData); m++ {

								// 	// 		if p == m {

								// 	// 			rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

								// 	// 			xlsx.SetCellValue(sheet4Company, col[m+indexData]+row, rowData)
								// 	// 			nextIndexDataCol = m + indexData
								// 	// 		}

								// 	// 	}
								// 	// }
								// }

								indexData++

							}
						}
						// }

						no++
					}
					// end row data --------------------------------------------------
				} // for organization
			}

			// delete sheet index 1
			// if iCompanyID >= 1 {
			xlsx.DeleteSheet(sheet1Name)
			// fmt.Println("xlsx.DeleteSheet(sheet1Name) -----2-------------------------->", sheet1Name)
			// }
			// END MUlty company --------------------------------------------

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			formName = strings.Replace(formName, "/", "-", 100)
			formName = strings.Replace(formName, "+", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed data generate Excel",
				})
				return
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data field is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}

}

func (ctr *inputFormController) DataFormDetail4Download(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("organizationID :::", organizationID)
	}

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")
	selectedCompanyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, _ := strconv.Atoi(selectedCompanyID)

	fmt.Println("iCompanyID", iCompanyID)

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--23jjk2jh4kj--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": true,
				"error":  err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [696]string{"H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx
			xlsx := excelize.NewFile()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                // xlsx

			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
			// end row data --------------

			// check multycompany here -----------------------------------------------------------------------------------------
			var wf tables.InputFormOrganizations
			wf.FormID = formID

			if iCompanyID >= 1 {
				wf.OrganizationID = iCompanyID

				// delete sheet index 1
				// xlsx.DeleteSheet(sheet1Name)
				// fmt.Println("xlsx.DeleteSheet(sheet1Name) ------------------------------->", sheet1Name)
			}

			checkOrgInputForm, err := ctr.inputForm.GetOrganizationInputForm(wf, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if len(checkOrgInputForm) >= 1 {

				for o := 0; o < len(checkOrgInputForm); o++ {

					sheet4Company := "Export-Form-Data-" + checkOrgInputForm[o].OrganizationName
					xlsx.NewSheet(sheet4Company)

					// row header
					xlsx.SetCellValue(sheet4Company, "A1", "No.")         // xlsx
					xlsx.SetCellValue(sheet4Company, "B1", "STATUS DATA") // xlsx
					xlsx.SetCellValue(sheet4Company, "C1", "NAME")        // xlsx
					xlsx.SetCellValue(sheet4Company, "D1", "DATE")        // xlsx
					xlsx.SetCellValue(sheet4Company, "E1", "HP")          // xlsx
					xlsx.SetCellValue(sheet4Company, "F1", "LAT")         // xlsx
					xlsx.SetCellValue(sheet4Company, "G1", "LNG")         // xlsx

					index := 0
					nextIndex1 := 0
					for i := 0; i < len(results); i++ {

						fmt.Println(":: ID ----", results[i].ID)
						if results[i].FieldTypeID == 22 {

							var tabDataRowHeader objects.TabDataRowHeader
							objHeader := results[i].Option
							json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

							if nextIndex1 > 0 {
								index = nextIndex1 + 1
							}
							nextIndexCol := 0
							for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

								xlsx.SetCellValue(sheet4Company, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
								fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", "-- label : "+results[i].Label, " | ", tabDataRowHeader.TabDataHeader[k].Value)
								nextIndexCol = index + k
							}
							nextIndex1 = nextIndexCol

						} else {
							if nextIndex1 > 0 {
								// fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[nextIndex1+1]+"1", results[i].Label)
								nextIndex1++
							} else {
								// fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
								xlsx.SetCellValue(sheet4Company, col[index]+"1", results[i].Label)
								index++
							}
						}

					}

					style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#e4e4e4"],"pattern":1}}`)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadGateway, gin.H{
							"status":  false,
							"message": err.Error(),
						})
						return
					}
					xlsx.SetCellStyle(sheet4Company, "A1", col[nextIndex1+1]+"1", style)
					xlsx.SetColWidth(sheet4Company, "B", col[index], 15)
					// end row header

					// row data --------------------
					whereData := ""
					var bufferInputFormOrg bytes.Buffer

					// checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					// if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
					// 	// company ID by TOKEN
					// 	bufferInputFormOrg.WriteString(" ifo.organization_id =" + claims["organization_id"].(string))
					// } else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
					// 	// select company option (form sharing only)
					// 	bufferInputFormOrg.WriteString(" ifo.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
					// }

					bufferInputFormOrg.WriteString(" ifo.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))

					if periodeStart != "" && periodeEnd == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputFormOrg.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					whereData = bufferInputFormOrg.String()

					var whereInForm tables.InputFormJoinOrganizations
					// whereInForm.OrganizationID, _ = strconv.Atoi(selectedCompanyID)
					getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, whereInForm, whereData, objects.Paging{})
					if err != nil {
						fmt.Println("err GetFormFieldRows--m42,23m4m,--", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// check empty data
					if len(getData) <= 0 {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"message": "Data field is not available",
							"data":    nil,
						})
						return
					}

					no := 1
					for a := 0; a < len(getData); a++ {

						// if a <= 1 {

						var statusData string
						if getData[a].UpdatedCount == 0 {
							statusData = "Terkirim"
						} else {
							statusData = "Diubah " + strconv.Itoa(getData[a].UpdatedCount) + " kali (" + getData[a].UpdatedAt.Format("2006-01-02 15:04") + ")"
						}

						row := strconv.Itoa(no + 1)
						xlsx.SetCellValue(sheet4Company, "A"+row, no)                   // xlsx
						xlsx.SetCellValue(sheet4Company, "B"+row, statusData)           // xlsx
						xlsx.SetCellValue(sheet4Company, "C"+row, getData[a].UserName)  // xlsx
						xlsx.SetCellValue(sheet4Company, "D"+row, getData[a].CreatedAt) // xlsx
						xlsx.SetCellValue(sheet4Company, "E"+row, getData[a].Phone)     // xlsx
						xlsx.SetCellValue(sheet4Company, "F"+row, getData[a].Latitude)  // xlsx
						xlsx.SetCellValue(sheet4Company, "G"+row, getData[a].Longitude) // xlsx

						// rows data --------------------------------------------
						indexData := 0
						nextIndexData := 0

						fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
						for j := 0; j < len(results); j++ {

							fieldID := strconv.Itoa(results[j].ID)
							fieldTypeID := results[j].FieldTypeID

							//data user here -------------------------------------------------------------
							var fields tables.InputForms
							fields.ID = getData[a].ID

							var bufferInputData bytes.Buffer
							fieldStrings := "coalesce(f" + fieldID + ",'') as f"

							bufferInputData.WriteString(" ifo.organization_id= " + strconv.Itoa(checkOrgInputForm[o].OrganizationID))

							if periodeStart != "" && periodeEnd == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
							}

							if periodeEnd != "" && periodeStart == "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
							}

							if periodeStart != "" && periodeEnd != "" {
								bufferInputData.WriteString(" AND to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
							}

							// if periodeStart != "" || periodeEnd != "" {

							// 	bufferInputData.WriteString(" AND ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// } else {
							// 	bufferInputData.WriteString(" ifa.organization_id =" + strconv.Itoa(checkOrgInputForm[o].OrganizationID))
							// }

							whereStr := bufferInputData.String()
							inputData, err := ctr.inputForm.GetInputDataOrganizationRows(formID, fieldStrings, fields, whereStr)
							if err != nil {
								fmt.Println("err GetFormFieldRows--wlerjlwekrj--", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							if len(inputData) > 0 {
								obj := inputData[0][0]

								if fieldTypeID != 22 { // (22 = tabulasi)
									rowData := strings.Replace(obj, "\n", " ", 100)

									if fieldTypeID == 3 || fieldTypeID == 21 {

										var dataOption objects.DataOption
										json.Unmarshal([]byte(rowData), &dataOption)

										if len(dataOption.Data) >= 1 {

											optVal := make([]string, len(dataOption.Data))
											for k, v := range dataOption.Data {
												optVal[k] = v.Value
											}

											rowData = strings.Join(optVal, ",")
										}

										// os.Exit(0)
									}
									xlsx.SetCellValue(sheet4Company, col[indexData]+row, rowData)

									fmt.Println("-----reguler XXX-----", col[indexData]+row, "--index-", indexData, row, rowData)

									// indexData++

								} else {

									nextIndexDataCol := 0

									var tabDataRowHeader objects.TabDataRowHeader
									objHeader := results[j].Option
									json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

									// fmt.Println("-----TabDataHeader-----", tabDataRowHeader.TabDataHeader)

									for p := 0; p < len(tabDataRowHeader.TabDataHeader); p++ {

										// fmt.Println("-----tabulasi-----", p)
										var fieldData []objects.TabValueAnswer
										json.Unmarshal([]byte(obj), &fieldData)

										if len(fieldData) >= 1 {

											for m := 0; m < len(fieldData); m++ {

												if p == m {
													rowData := strings.Replace(fieldData[m].Answer, "\n", " ", 100)

													xlsx.SetCellValue(sheet4Company, col[indexData]+row, rowData)
													// fmt.Println("-----tabulasi XXX-----", p, "----m--", m, "--index-", indexData, col[indexData]+row, "data ::", rowData)
													nextIndexDataCol = m + indexData
													indexData++
												}

											}

										}
									}
									fmt.Println("nextIndexDataCol ::", nextIndexDataCol)
									// indexData = nextIndexDataCol
									indexData--
								}

								indexData++
							}
						}

						//} //if loop 3

						no++
					}
					// end row data --------------------------------------------------
				} // for organization
			}
			// os.Exit(0)

			// delete sheet index 1
			// if iCompanyID >= 1 {
			xlsx.DeleteSheet(sheet1Name)
			// fmt.Println("xlsx.DeleteSheet(sheet1Name) -----2-------------------------->", sheet1Name)
			// }
			// END MUlty company --------------------------------------------

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			formName = strings.Replace(formName, "/", "-", 100)
			formName = strings.Replace(formName, "+", "-", 100)
			formName = strings.Replace(formName, ",", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed data generate Excel",
				})
				return
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data field is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}

}

// multyorganization
func (ctr *inputFormController) DataFormDetail4DownloadCSV(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("organizationID :::", organizationID)
	}

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows--32l4kl2k4--", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header --------------------------
			var header = []string{"NAME", "COMPANY", "STATUS DATA", "DATE", "HP", "LAT", "LONG"}

			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {
				var fieldActive string
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						fmt.Println(":: -------index kolom tabulasi ::: ", index+k, tabDataRowHeader.TabDataHeader[k].Value)
						// xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						fieldActive = results[i].Label + " | " + tabDataRowHeader.TabDataHeader[k].Value
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						fmt.Println("::---else 1-------", nextIndex1+1, "-----", results[i].Label)

						fieldActive = results[i].Label
						nextIndex1++
					} else {
						fmt.Println("::-----else 2-----", index, "-----", results[i].Label)

						fieldActive = results[i].Label
						index++
					}
				}

				header = append(header, fieldActive)

			}

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			if periodeStart != "" && periodeEnd == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' "
			}

			if periodeEnd != "" && periodeStart == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  "
			}

			if periodeStart != "" && periodeEnd != "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' "
			}

			getData, err := ctr.inputForm.GetInputFormOrganizationRows(formID, tables.InputFormJoinOrganizations{}, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows--,m4234l2l4kmk--", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			var exportData = make([][]string, len(getData)+1)
			exportData[0] = header

			if len(getData) > 0 {
				no := 1
				for a := 0; a < len(getData); a++ {

					// row := strconv.Itoa(no + 1)
					timeStr := getData[a].CreatedAt.Format("2006-01-02")
					lat := strconv.FormatFloat(getData[a].Latitude, 'f', 50, 64)
					long := strconv.FormatFloat(getData[a].Longitude, 'f', 50, 64)

					var statusData string
					if getData[a].UpdatedCount == 0 {
						statusData = "Terkirim"
					} else {
						statusData = "Diubah " + strconv.Itoa(getData[a].UpdatedCount) + " kali"
					}

					var rowData = []string{getData[a].UserName, getData[a].OrganizationName, statusData, timeStr, getData[a].Phone, lat, long}
					// fmt.Println("formating :: ==============", getData[a].CreatedAt, timeStr)
					// col data --------------------------------------------
					indexData := 0
					for j := 0; j < len(results); j++ {

						var rowColActive string

						fieldID := strconv.Itoa(results[j].ID)
						fieldTypeID := results[j].FieldTypeID

						//data user here -------------------------------------------------------------
						var fields tables.InputForms
						fields.ID = getData[a].ID

						var bufferInputData bytes.Buffer
						fieldStrings := "coalesce(f" + fieldID + ",'') as f"

						if periodeStart != "" && periodeEnd == "" {
							bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
						}

						if periodeEnd != "" && periodeStart == "" {
							bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
						}

						if periodeStart != "" && periodeEnd != "" {
							bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
						}

						fmt.Println(periodeStart, periodeEnd)
						whereStr := bufferInputData.String()
						inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, whereStr)
						if err != nil {
							fmt.Println("err GetFormFieldRows--n423k43m--", err)
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}

						if len(inputData) >= 1 {
							obj := inputData[0][0]

							if fieldTypeID != 22 { // (22 = tabulasi)

								if fieldTypeID == 3 || fieldTypeID == 21 {

									var dataOption objects.DataOption
									json.Unmarshal([]byte(obj), &dataOption)

									if len(dataOption.Data) >= 1 {
										//
										optVal := make([]string, len(dataOption.Data))
										for k, v := range dataOption.Data {
											optVal[k] = v.Value
										}

										obj = strings.Join(optVal, ",")
									}

								}
								rowColActive = strings.Replace(obj, "\n", " ", 50)

							} else {

								var tabDataRowHeader objects.TabDataRowHeader
								objHeader := results[j].Option
								json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

								for p := 0; p < len(tabDataRowHeader.TabDataHeader); p++ {

									fmt.Println("tabDataRowHeader.TabDataHeader  :: --------------", tabDataRowHeader.TabDataHeader[p].Value)

									// var fieldData []objects.TabValueAnswer
									// json.Unmarshal([]byte(obj), &fieldData)

									nextIndexDataCol := 0
									// fmt.Println("len  :: --------------", len(fieldData), results[j].Label)

									// for m := 0; m < len(fieldData); m++ {

									// 	rowColActive = fieldData[m].Answer
									// 	nextIndexDataCol = m + indexData

									// 	fmt.Println("TAB  :: --------------", results[j].Label)

									// 	rowData = append(rowData, rowColActive)

									// }
									indexData = nextIndexDataCol
								}
								fmt.Println("fieldData[m].Answer  :: --------------", results[j].Label)

								if ctr.conf.ENV_TYPE == "dev" {
									if results[j].ID == 9340 {
										fmt.Println("fieldData[m].Answer 9340 :: --------------", results[j].Label)
										// os.Exit(0)
									}
								}
							}
							indexData++
						}

						rowData = append(rowData, rowColActive)
					}

					exportData[no] = rowData
					no++

				}

				// if ctr.conf.ENV_TYPE == "dev" {
				// 	os.Exit(0)
				// }

				// CONFIG file --------------------------------------------------
				var fieldForm tables.Forms
				fieldForm.ID = formID
				getForm, _ := ctr.formMod.GetFormRow(fieldForm)

				today := time.Now()
				dateFormat := today.Format("02012006-150405")

				formName := strings.Replace(getForm.Name, " ", "-", 100)
				formName = strings.Replace(formName, "/", "-", 100)
				fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
				fileGroup := "form_download"
				fileExtention := "csv"
				fileLocation := "file/" + fileName + "." + fileExtention // local path file location

				f, err := os.Create(fileLocation)
				// defer f.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				w := csv.NewWriter(f)
				// w.UseCRLF = false
				err = w.WriteAll(exportData) // calls Flush internally
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
				fmt.Println(fileGroup)
				var obj objects.FileRes
				obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": err,
						"data":    nil,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Data is available",
					"data":    obj,
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
		})
		return
	}
}

func (ctr *inputFormController) DataFormDetailDownloadCSV(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldNotParentRows(fields, "")
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header --------------------------
			var header = []string{"NAME", "DATE", "HP", "LAT", "LONG"}

			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {
				var fieldActive string
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						fmt.Println(":: -------index kolom tabulasi ::: ", index+k, tabDataRowHeader.TabDataHeader[k].Value)
						// xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						fieldActive = results[i].Label + " | " + tabDataRowHeader.TabDataHeader[k].Value
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						fmt.Println("::---else 1-------", nextIndex1+1, "-----", results[i].Label)
						// xlsx.SetCellValue(sheet1Name, col[nextIndex1+1]+"1", results[i].Label)
						fieldActive = results[i].Label
						nextIndex1++
					} else {
						fmt.Println("::-----else 2-----", index, "-----", results[i].Label)
						// xlsx.SetCellValue(sheet1Name, col[index]+"1", results[i].Label)
						fieldActive = results[i].Label
						index++
					}
				}

				header = append(header, fieldActive)

			}

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			if periodeStart != "" && periodeEnd == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' "
			}

			if periodeEnd != "" && periodeStart == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  "
			}

			if periodeStart != "" && periodeEnd != "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' "
			}

			var whereInForm tables.InputForms

			getData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			var exportData = make([][]string, len(getData)+1)
			exportData[0] = header

			if len(getData) > 0 {
				no := 1
				for a := 0; a < len(getData); a++ {

					// row := strconv.Itoa(no + 1)
					timeStr := getData[a].CreatedAt.Format("2006-01-02")
					lat := strconv.FormatFloat(getData[a].Latitude, 'f', 50, 64)
					long := strconv.FormatFloat(getData[a].Longitude, 'f', 50, 64)
					var rowData = []string{getData[a].UserName, timeStr, getData[a].Phone, lat, long}
					// fmt.Println("formating :: ==============", getData[a].CreatedAt, timeStr)
					// col data --------------------------------------------
					indexData := 0
					for j := 0; j < len(results); j++ {

						var rowColActive string

						fieldID := strconv.Itoa(results[j].ID)
						fieldTypeID := results[j].FieldTypeID

						//data user here -------------------------------------------------------------
						var fields tables.InputForms
						fields.ID = getData[a].ID

						var bufferInputData bytes.Buffer
						fieldStrings := "coalesce(f" + fieldID + ",'') as f"

						if periodeStart != "" && periodeEnd == "" {
							bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
						}

						if periodeEnd != "" && periodeStart == "" {
							bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
						}

						if periodeStart != "" && periodeEnd != "" {
							bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
						}

						fmt.Println(periodeStart, periodeEnd)
						whereStr := bufferInputData.String()
						inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, whereStr)
						if err != nil {
							fmt.Println("err GetFormFieldRows----", err)
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}

						if len(inputData) > 0 {
							obj := inputData[0][0]

							if fieldTypeID != 22 { // (22 = tabulasi)
								rowColActive = strings.Replace(obj, "\n", " ", 50)

							} else {

								var fieldData []objects.TabValueAnswer
								json.Unmarshal([]byte(obj), &fieldData)

								nextIndexDataCol := 0
								for m := 0; m < len(fieldData); m++ {

									rowColActive = fieldData[m].Answer
									nextIndexDataCol = m + indexData

								}
								indexData = nextIndexDataCol
							}
							indexData++
						}

						rowData = append(rowData, rowColActive)
					}

					exportData[no] = rowData
					no++

				}

				// CONFIG file --------------------------------------------------
				var fieldForm tables.Forms
				fieldForm.ID = formID
				getForm, _ := ctr.formMod.GetFormRow(fieldForm)

				today := time.Now()
				dateFormat := today.Format("02012006-150405")

				formName := strings.Replace(getForm.Name, " ", "-", 100)
				formName = strings.Replace(formName, "/", "-", 100)
				fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
				fileGroup := "form_download"
				fileExtention := "csv"
				fileLocation := "file/" + fileName + "." + fileExtention // local path file location

				f, err := os.Create(fileLocation)
				// defer f.Close()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				w := csv.NewWriter(f)
				// w.UseCRLF = false
				err = w.WriteAll(exportData) // calls Flush internally
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
				fmt.Println(fileGroup)
				var obj objects.FileRes
				obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": err,
						"data":    nil,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Data is available",
					"data":    obj,
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
}

func (ctr *inputFormController) DataFormDetailDownloadCSV__(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2 // -2 filetype is not null
		results, err := ctr.formField.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header
			col := [697]string{"F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ", "DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ", "EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ", "FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ", "GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ", "HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ", "IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ", "JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ", "KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ", "LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ", "MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ", "OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ", "QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ", "RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ", "SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ", "TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ", "UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ", "VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ", "XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ", "YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ"} // xlsx

			xlsx := excelize.NewFile() // xlsx
			sheet1Name := "Export-Form-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

			xlsx.SetCellValue(sheet1Name, "A1", "NAME") // xlsx
			xlsx.SetCellValue(sheet1Name, "B1", "DATE") // xlsx
			xlsx.SetCellValue(sheet1Name, "C1", "HP")   // xlsx
			xlsx.SetCellValue(sheet1Name, "D1", "LAT")  // xlsx
			xlsx.SetCellValue(sheet1Name, "E1", "LNG")  // xlsx

			// row header --------------------------
			index := 0
			nextIndex1 := 0
			for i := 0; i < len(results); i++ {

				fmt.Println(":: ID ----", results[i].ID)
				if results[i].FieldTypeID == 22 {

					var tabDataRowHeader objects.TabDataRowHeader
					objHeader := results[i].Option
					json.Unmarshal([]byte(objHeader), &tabDataRowHeader)

					if nextIndex1 > 0 {
						index = nextIndex1 + 1
					}
					nextIndexCol := 0
					for k := 0; k < len(tabDataRowHeader.TabDataHeader); k++ {

						fmt.Println(":: -------index kolom tabulasi ::: ", index+k, "-------", col[index+k]+"1", tabDataRowHeader.TabDataHeader[k].Value)
						xlsx.SetCellValue(sheet1Name, col[index+k]+"1", results[i].Label+" | "+tabDataRowHeader.TabDataHeader[k].Value)
						nextIndexCol = index + k
					}
					nextIndex1 = nextIndexCol

				} else {
					if nextIndex1 > 0 {
						fmt.Println("::---else 1-------", nextIndex1+1, "-----", col[nextIndex1+1]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[nextIndex1+1]+"1", results[i].Label)
						nextIndex1++
					} else {
						fmt.Println("::-----else 2-----", index, "-----", col[index]+"1", results[i].Label)
						xlsx.SetCellValue(sheet1Name, col[index]+"1", results[i].Label)
						index++
					}
				}

			}

			// row data -------------------------------------------------------------------------------------------------------
			whereData := ""
			if periodeStart != "" && periodeEnd == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' "
			}

			if periodeEnd != "" && periodeStart == "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  "
			}

			if periodeStart != "" && periodeEnd != "" {
				whereData = " to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' "
			}

			var whereInForm tables.InputForms

			getData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whereData, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			no := 1
			for a := 0; a < len(getData); a++ {

				row := strconv.Itoa(no + 1)
				xlsx.SetCellValue(sheet1Name, "A"+row, getData[a].UserName)  // xlsx
				xlsx.SetCellValue(sheet1Name, "B"+row, getData[a].CreatedAt) // xlsx
				xlsx.SetCellValue(sheet1Name, "C"+row, getData[a].Phone)     // xlsx
				xlsx.SetCellValue(sheet1Name, "D"+row, getData[a].Latitude)  // xlsx
				xlsx.SetCellValue(sheet1Name, "E"+row, getData[a].Longitude) // xlsx

				// rows data --------------------------------------------
				indexData := 0
				nextIndexData := 0

				fmt.Println(periodeStart, periodeEnd, nextIndexData, indexData)
				for j := 0; j < len(results); j++ {

					fieldID := strconv.Itoa(results[j].ID)
					fieldTypeID := results[j].FieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[a].ID

					var bufferInputData bytes.Buffer
					fieldStrings := "coalesce(f" + fieldID + ",'') as f"

					if periodeStart != "" && periodeEnd == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
					}

					if periodeEnd != "" && periodeStart == "" {
						bufferInputData.WriteString("to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
					}

					if periodeStart != "" && periodeEnd != "" {
						bufferInputData.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
					}
					fmt.Println(periodeStart, periodeEnd)
					whereStr := bufferInputData.String()
					inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, whereStr)
					if err != nil {
						fmt.Println("err GetFormFieldRows----", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					if len(inputData) > 0 {
						obj := inputData[0][0]

						if fieldTypeID != 22 { // (22 = tabulasi)

							xlsx.SetCellValue(sheet1Name, col[indexData]+row, obj)

						} else {

							var fieldData []objects.TabValueAnswer
							json.Unmarshal([]byte(obj), &fieldData)

							nextIndexDataCol := 0
							for m := 0; m < len(fieldData); m++ {

								xlsx.SetCellValue(sheet1Name, col[m+indexData]+row, fieldData[m].Answer)
								nextIndexDataCol = m + indexData

							}
							indexData = nextIndexDataCol
						}
						indexData++
					}

				}

				no++
			}

			// err1 := xlsx.AutoFilter(sheet1Name, "A1", "C1", "")
			// if err1 != nil {
			// 	log.Fatal("ERROR", err1.Error())
			// }

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			fileName := "Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "form_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(fileGroup, err2)
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err2,
					"data":    nil,
				})
				return
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}

}

func (ctr *inputFormController) DataFormRespondenList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	periode := c.Param("periode")
	year := c.Param("year")
	if year == "" {
		year = time.Now().Format("2006")
	}

	companyID := c.Request.URL.Query().Get("company_id")
	iCompanyID, err := strconv.Atoi(companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Company ID is required",
			"status":  false,
		})
		return
	}

	if formID > 0 {

		checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
		whre := " AND fu.type='respondent' AND fuo.organization_id=" + claims["organization_id"].(string)
		if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
			// company ID by TOKEN
			whre = " AND fu.type='respondent' AND fuo.organization_id=" + claims["organization_id"].(string)
		} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
			// select company option (form sharing only)
			whre = " AND fu.type='respondent' AND fuo.organization_id=" + strconv.Itoa(iCompanyID)
		}

		fmt.Println("whre ----->", whre)
		fmt.Println("IDs ----->", organizationID, checkFormCompany.OrganizationID, iCompanyID)

		formUserData, err := ctr.userMod.InputFormUserLeftJoin(formID, whre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println("len(formUserData) :::", len(formUserData))
		if len(formUserData) > 0 {

			// get date column ----------------------------------------------------------
			var resHeader []objects.InputFormFields

			if periode == "daily" {
				strWhre := "where to_char(d,'dd') = to_char(now(), 'dd')"
				getDates, err := ctr.inputForm.GetDates(strWhre)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				for i := 0; i < len(getDates); i++ {
					var each objects.InputFormFields
					each.ID = i + 1
					each.Label = getDates[i].Date

					resHeader = append(resHeader, each)
				}
			} else if periode == "monthly" {
				strWhre := ""
				getDates, err := ctr.inputForm.GetDates(strWhre)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				for i := 0; i < len(getDates); i++ {
					var each objects.InputFormFields
					each.ID = i + 1
					each.Label = getDates[i].Date

					resHeader = append(resHeader, each)
				}
			} else if periode == "yearly" {
				strWhre := ""
				getDates, err := ctr.inputForm.GetMonths(strWhre)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				for i := 0; i < len(getDates); i++ {
					var each objects.InputFormFields
					each.ID = i + 1
					each.Label = getDates[i].Month

					resHeader = append(resHeader, each)
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Parameter periode is wrong",
					"data":    nil,
				})
				return
			}

			// get rows data ----------------------------------------------------------------

			var rows []objects.Rows
			for i := 0; i < len(formUserData); i++ {

				whre := ""
				var getData []tables.TotalDate
				if periode == "daily" {

					fmt.Println("daillly :::", formUserData[i].UserID)
					// whre = "and to_char(d, 'yyyy-mm-dd') = to_char(now(), 'yyyy-mm-dd') "

					var bufferIF bytes.Buffer
					bufferIF.WriteString(" AND TO_CHAR(d, 'yyyy-mm-dd') = TO_CHAR(now(), 'yyyy-mm-dd') ")
					if year != "" {
						bufferIF.WriteString(" AND TO_CHAR(d::date, 'yyyy') = '" + year + "'")
					}

					checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
						// company ID by TOKEN
						bufferIF.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
						// select company option (form sharing only)
						bufferIF.WriteString(" AND ifo.organization_id= " + strconv.Itoa(iCompanyID))
					}
					whre = bufferIF.String()

					getData, err = ctr.inputForm.GetTotalDate(formUserData[i].UserID, formID, whre)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

				} else if periode == "monthly" {

					var bufferIF bytes.Buffer
					if year != "" {
						bufferIF.WriteString(" AND TO_CHAR(d::date, 'yyyy') = '" + year + "'")
					}

					checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
						// company ID by TOKEN
						bufferIF.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
						// select company option (form sharing only)
						bufferIF.WriteString(" AND ifo.organization_id= " + strconv.Itoa(iCompanyID))
					}
					whre = bufferIF.String()

					getData, err = ctr.inputForm.GetTotalDateMonthly(formUserData[i].UserID, formID, whre)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
				}

				if periode == "yearly" {
					var bufferIF bytes.Buffer
					if year != "" {
						bufferIF.WriteString(" AND TO_CHAR(d::date, 'yyyy') = '" + year + "'")
					}
					checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
					if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
						// company ID by TOKEN
						bufferIF.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
					} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
						// select company option (form sharing only)
						bufferIF.WriteString(" AND ifo.organization_id= " + strconv.Itoa(iCompanyID))
					}
					whre = bufferIF.String()

					getData, err = ctr.inputForm.GetTotalMonth(formUserData[i].UserID, formID, whre)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
				}

				var inputData []objects.InputData
				if len(getData) > 0 {

					for j := 0; j < len(getData); j++ {

						var each objects.InputData
						each.ID = j + 1
						each.Date = getData[j].Date
						each.Value = getData[j].Value

						inputData = append(inputData, each)
					}
				} else {

					for j := 0; j < 1; j++ {

						var each objects.InputData
						inputData = append(inputData, each)
					}
				}

				var each objects.Rows
				each.UserID = formUserData[i].UserID
				each.UserName = formUserData[i].UserName
				each.UserPhone = formUserData[i].UserPhone
				each.Organizations = formUserData[i].Organizations
				each.SubmitDate = formUserData[i].SubmitDate
				each.InputData = inputData

				rows = append(rows, each)
				// }
			}

			fmt.Println("ROWSSS ========", len(rows))

			var res objects.InputFormUsers
			res.FormID = formID
			res.FieldHeader = resHeader
			res.FieldData = rows

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true, // ini sengaja true request FE jgn diubah
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"error":   err,
		})
		return
	}
}

func (ctr *inputFormController) ReportFormRespondenList(c *gin.Context) {

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {

		whreStr := ""
		if searchKeyWord != "" {
			whreStr = " WHERE u.name ilike '%" + searchKeyWord + "%'  "
		}

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.Sort = sort
		paging.SortBy = sortBy

		formUserData, err := ctr.userMod.InputFormUserPaging(formID, whreStr, paging)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		resultAll, err := ctr.userMod.InputFormUserPaging(formID, whreStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		var whreForm tables.Forms
		whreForm.ID = formID
		formData, err := ctr.formMod.GetFormRow(whreForm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(formUserData) > 0 {

			// get rows data ----------------------------------------------------------------
			var lastSubmitDate string
			var rows []objects.ReportResponden
			for i := 0; i < len(formUserData); i++ {

				//get total responden
				var whereInForm tables.InputForms
				whereInForm.UserID = formUserData[i].UserID
				whereToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
				getDataRespons, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whereToday, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				submissionColor := ""

				// if len(getDataRespons) < 50 {
				// 	submissionColor = "#F9B3B3"
				// }

				// performance per user
				var perform int
				var performFloat float64
				performColor := ""
				if formData.SubmissionTargetUser > 0 {
					performFloat = (float64(len(getDataRespons)) / float64(formData.SubmissionTargetUser)) * 100
					perform, _ = strconv.Atoi(strconv.FormatFloat(performFloat, 'f', 0, 64))

					if perform >= 100 {
						perform = 100
						performFloat = 100.0
					}

					if perform < 50 {
						performColor = "#F9B3B3"
					}

				}

				// get last submit
				var lastSubmit string
				var lastAddress string
				var lastLatitude float64
				var lastLongitude float64

				if len(getDataRespons) > 0 {
					for i := 0; i < len(getDataRespons); i++ {
						fmt.Println("created ::", getDataRespons[i].CreatedAt)
						if i == 0 {
							lastSubmitDate = getDataRespons[i].CreatedAt.Format("2006-02-01 15:04")
							now := time.Now()

							diff := now.Sub(getDataRespons[i].CreatedAt)

							if diff.Minutes() > 0 {
								lastSubmit = strconv.FormatFloat(diff.Minutes(), 'f', 0, 64) + " menit yang lalu"
							}

							lastAddress = getDataRespons[i].Address
							lastLatitude = getDataRespons[i].Latitude
							lastLongitude = getDataRespons[i].Longitude
						}
					}
				}

				if formUserData[i].UserID == 11 {
					fmt.Println("perform ---------11------", perform, len(getDataRespons), submissionColor, performColor)
				}

				var each objects.ReportResponden
				each.UserID = formUserData[i].UserID
				each.UserName = formUserData[i].UserName
				each.UserPhone = formUserData[i].UserPhone
				each.Avatar = formUserData[i].Avatar
				each.Submission = len(getDataRespons)
				each.Performance = perform
				each.PerformanceFloat = performFloat
				each.LastSubmission = lastSubmit
				each.LastSubmissionDate = lastSubmitDate
				each.SubmissionColor = performColor
				each.PerformanceColor = performColor
				each.Address = lastAddress
				each.Latitude = lastLatitude
				each.Longitude = lastLongitude

				rows = append(rows, each)

			}

			totalPage := 0
			if limit > 0 {
				totalPage = len(resultAll) / limit
				if (len(resultAll) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var paging objects.DataRows
			paging.TotalRows = len(resultAll)
			paging.TotalPages = totalPage

			var dataDetail objects.DataDetail
			dataDetail.LastSubmissionDate = lastSubmitDate

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    rows,
				"paging":  paging,
				"detail":  dataDetail,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true, // ini sengaja true request FE jgn diubah
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"error":   err,
		})
		return
	}
}

func (ctr *inputFormController) ReportFormRespondenList2(c *gin.Context) {

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("strconv.Atoi", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")
	totalDataUpdated := 0
	activeSubmission := 0

	if formID > 0 {

		whreStr := ""
		if searchKeyWord != "" {
			whreStr = " WHERE u.name ilike '%" + searchKeyWord + "%'  "
		}

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.Sort = sort
		paging.SortBy = sortBy

		formUserData, err := ctr.inputForm.GetReportFormRespondenUnionTeam(formID, tables.InputForms{}, whreStr, paging)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		var pagingAll objects.Paging
		pagingAll.Sort = "desc"
		pagingAll.SortBy = "last_submission_date"
		resultAll, err := ctr.inputForm.GetReportFormRespondenUnionTeam(formID, tables.InputForms{}, whreStr, pagingAll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		var wFormUser tables.JoinFormUsers
		wFormUser.FormID = formID
		getUserForm, _ := ctr.formMod.GetFormUserRows(wFormUser, "")

		var whereAll tables.InputForms
		whereToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
		getAllDataSubmission, err := ctr.inputForm.GetInputFormRows(formID, whereAll, whereToday, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":   "GetInputFormRows all",
				"error": err,
			})
			return
		}

		totalAverage := 0.0
		if len(getAllDataSubmission) > 0 {
			totalAverage = (float64(len(getAllDataSubmission)) / float64(len(getUserForm)))
		}
		if len(formUserData) > 0 {
			var values []objects.ReportResponden
			for i := 0; i < len(formUserData); i++ {

				var whereInForm tables.InputForms
				whereInForm.UserID = formUserData[i].UserID
				whereToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
				dataDeleted, err := ctr.inputForm.GetDeletedData(formID, whereInForm, whereToday)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":   "GetDeletedData",
						"error": err,
					})
					return
				}

				whereDate := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
				dataUpdated, err := ctr.inputForm.GetUpdatedData(formID, formUserData[i].UserID, whereDate)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":   "GetUpdatedData",
						"error": err,
					})
					return
				}
				for i := 0; i < len(dataUpdated); i++ {
					totalDataUpdated = dataUpdated[i].UpdatedCount
				}
				var each objects.ReportResponden
				activeSubmission = formUserData[i].Submission - len(dataDeleted)

				each.UserID = formUserData[i].UserID
				each.UserName = formUserData[i].UserName
				each.UserPhone = formUserData[i].UserPhone
				each.Submission = formUserData[i].Submission
				each.TotalAverage = totalAverage
				each.Performance = formUserData[i].Performance
				each.PerformanceFloat = formUserData[i].PerformanceFloat
				each.SubmissionTargetUser = formUserData[i].SubmissionTargetUser
				each.LastSubmission = formUserData[i].LastSubmission
				each.LastSubmissionDate = formUserData[i].LastSubmissionDate
				each.SubmissionColor = formUserData[i].SubmissionColor
				each.PerformanceColor = formUserData[i].PerformanceColor
				each.Address = formUserData[i].Address
				each.Latitude = formUserData[i].Latitude
				each.Longitude = formUserData[i].Longitude
				each.Avatar = formUserData[i].Avatar
				each.TotalUpdateData = totalDataUpdated
				each.TotalDeletedData = len(dataDeleted)
				each.TotalActiveSubmission = strconv.Itoa(activeSubmission) + "/" + strconv.Itoa(formUserData[i].Submission)

				values = append(values, each)
			}
			totalPage := 0
			if limit > 0 {
				totalPage = len(resultAll) / limit
				if (len(resultAll) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var paging objects.DataRows
			paging.TotalRows = len(resultAll)
			paging.TotalPages = totalPage

			var dataDetail objects.DataDetail
			dataDetail.LastSubmissionDate = resultAll[len(resultAll)-1].LastSubmissionDate

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    values,
				// "data":    formUserData,
				"paging": paging,
				"detail": dataDetail,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true, // ini sengaja true request FE jgn diubah
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"error":   err,
		})
		return
	}
}

func (ctr *inputFormController) DataFormResponMapList(c *gin.Context) {

	formID, err := strconv.Atoi(c.Param("formid"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	periodeStart := c.Request.URL.Query().Get("periode_start")
	periodeEnd := c.Request.URL.Query().Get("periode_end")

	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	filterFieldIDs := c.Request.URL.Query().Get("filter_field_ids")

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = 5

		if filterFieldIDs != "" {
			fields.ID = 5
		}

		result, err := ctr.formField.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var fieldLabels []objects.InputFormFields

			// row header ---------------------------------------------------------------------------
			for i := 0; i < len(result); i++ {

				var each objects.InputFormFields

				if result[i].FieldTypeID == 5 {

					each.ID = result[i].ID
					each.Label = result[i].Label
					each.FieldTypeName = result[i].FieldTypeName

					// each.TabHeader = fieldLabel
				}

				fieldLabels = append(fieldLabels, each)
			}

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			var whereInForm tables.InputForms

			var buffer bytes.Buffer
			var whreStr = ``
			if searchKeyWord != "" {
				buffer.WriteString(" u.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			if periodeStart != "" && periodeEnd == "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') >= '" + periodeStart + "' ")
			}

			if periodeEnd != "" && periodeStart == "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') <= '" + periodeEnd + "'  ")
			}

			if periodeStart != "" && periodeEnd != "" {
				buffer.WriteString(" to_char(if.created_at,'yyyy-mm-dd') BETWEEN '" + periodeStart + "' AND '" + periodeEnd + "' ")
			}

			whreStr = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whreStr, paging)
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			getAllData, err := ctr.inputForm.GetInputFormRows(formID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			var dataRows []objects.InputFieldDataMap
			for i := 0; i < len(getData); i++ {

				var each objects.InputFieldDataMap
				each.ID = getData[i].ID
				each.UserID = getData[i].UserID
				each.UserName = getData[i].UserName
				each.Phone = getData[i].Phone
				each.CreatedAt = getData[i].CreatedAt.Format("2006-01-02 15:04")

				// cols data --------------------------------------------
				var resultData []objects.InputDataMap
				for j := 0; j < len(result); j++ {

					fieldID := strconv.Itoa(result[j].ID)
					fieldTypeID := result[j].FieldTypeID

					var ec objects.InputDataMap
					ec.FieldID = fieldID
					ec.FieldTypeID = fieldTypeID

					//data user here -------------------------------------------------------------
					var fields tables.InputForms
					fields.ID = getData[i].ID
					fieldStrings := "coalesce(f" + fieldID + ",'') as f"

					inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, "")
					if err != nil {
						fmt.Println("err GetFormFieldRows----", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					ec.Value = inputData[0][0]

					resultData = append(resultData, ec)
				}

				each.InputData = resultData
				dataRows = append(dataRows, each)
			}

			// form detail data ------------------------------------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, err := ctr.formMod.GetFormRow(fieldForm)

			var res objects.InputFormMapRes
			res.FormID = formID
			res.FormName = getForm.Name
			res.FormDescription = getForm.Description
			res.PeriodStartDate = getForm.PeriodStartDate
			res.PeriodEndDate = getForm.PeriodEndDate
			res.FieldLabel = fieldLabels
			res.FieldData = dataRows

			totalPage := 0
			if limit > 0 {
				totalPage = len(getAllData) / limit
				if (len(getAllData) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var pagingRes objects.DataRows
			pagingRes.TotalRows = len(getAllData)
			pagingRes.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  pagingRes,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *inputFormController) DataFormResponMapPostList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	fmt.Println("userID ::", userID)

	var reqData objects.FormDataMap
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

	formID := reqData.FormID
	iCompanyID := reqData.CompanyID

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = 5

		whrStr1 := ``
		if len(reqData.FilterFieldIDs) > 0 {

			fieldIDs := make([]string, len(reqData.FilterFieldIDs))
			for i, v := range reqData.FilterFieldIDs {

				fieldIDs[i] = strconv.Itoa(v.FieldID)
			}

			whrStr1 = fmt.Sprintf("AND form_fields.id in (%s)", strings.Join(fieldIDs, ","))

		}

		results, err := ctr.formField.GetFormFieldNotParentRows(fields, whrStr1)
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println("------", len(results))
		if len(results) > 0 {

			var fieldLabels []objects.InputFormFields

			// row header ---------------------------------------------------------------------------
			for i := 0; i < len(results); i++ {

				var each objects.InputFormFields

				if results[i].FieldTypeID == 5 {

					each.ID = results[i].ID
					each.Label = results[i].Label
					each.FieldTypeName = results[i].FieldTypeName

					// each.TabHeader = fieldLabel
				}

				fieldLabels = append(fieldLabels, each)
			}

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			// var whereInForm tables.InputForms
			// if reqData.StartDate != "" && reqData.EndDate != "" {
			// 	whreStr = " to_char(if.created_at, 'yyyy-mm-dd') BETWEEN '" + reqData.StartDate + "' AND '" + reqData.EndDate + "'"
			// }

			checkFormCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})
			var whreInputFrm bytes.Buffer
			if reqData.StartDate != "" && reqData.EndDate != "" {
				whreInputFrm.WriteString("  to_char(if.created_at, 'yyyy-mm-dd') BETWEEN '" + reqData.StartDate + "' AND '" + reqData.EndDate + "'")
			}

			if organizationID >= 1 && organizationID != checkFormCompany.OrganizationID && iCompanyID <= 0 {
				// company ID by TOKEN
				whreInputFrm.WriteString(" AND ifo.organization_id= " + claims["organization_id"].(string))
			} else if organizationID >= 1 && organizationID == checkFormCompany.OrganizationID && iCompanyID >= 1 {
				// select company option (form sharing only)
				whreInputFrm.WriteString(" AND ifo.organization_id= " + strconv.Itoa(iCompanyID))
			}
			whreStr := whreInputFrm.String()

			getData, err := ctr.inputForm.GetInputFormRows(formID, tables.InputForms{}, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			var dataRows []objects.InputFieldDataMap
			if len(getData) >= 1 {
				for k := 0; k < len(getData); k++ {

					var each objects.InputFieldDataMap
					each.ID = getData[k].ID
					each.UserID = getData[k].UserID
					each.UserName = getData[k].UserName
					each.Phone = getData[k].Phone
					each.Avatar = getData[k].Avatar
					each.CreatedAt = getData[k].CreatedAt.Format("2006-01-02 15:04")

					// cols data --------------------------------------------
					var resultData []objects.InputDataMap
					if len(results) >= 1 {

						for j := 0; j < len(results); j++ {

							fieldID := strconv.Itoa(results[j].ID)
							fieldTypeID := results[j].FieldTypeID

							var ec objects.InputDataMap
							ec.FieldID = fieldID
							ec.FieldTypeID = fieldTypeID
							ec.TagLocIcon = results[j].TagLocIcon
							ec.TagLocColor = results[j].TagLocColor

							//data user here -------------------------------------------------------------
							var fields tables.InputForms
							fields.ID = getData[k].ID
							fieldStrings := "coalesce(f" + fieldID + ",'') as f"

							inputData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, fields, "")
							if err != nil {
								fmt.Println("err GetFormFieldRows----", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							ec.Value = inputData[0][0]

							// latLong := strings.Split(inputData[0][0], ",")
							regex, err := regexp.Compile(`\[.*?\]`)
							if err != nil {
								fmt.Println("err GetFormFieldRows----", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							res := regex.FindAllString(inputData[0][0], 1)
							if len(res) >= 1 {
								latlong := strings.Replace(res[0], "[", "", -1)
								latlong2 := strings.Replace(latlong, "]", "", -1)

								latLongSplit := strings.Split(latlong2, ",")
								lat, err := strconv.ParseFloat(strings.TrimSpace(latLongSplit[0]), 10)
								if err != nil {
									fmt.Println("err GetFormFieldRows----", err)
									c.JSON(http.StatusBadRequest, gin.H{
										"error": err,
									})
									return
								}
								long, err := strconv.ParseFloat(strings.TrimSpace(latLongSplit[1]), 10)
								if err != nil {
									fmt.Println("err GetFormFieldRows----", err)
									c.JSON(http.StatusBadRequest, gin.H{
										"error": err,
									})
									return
								}

								ec.TagLoc = latlong2
								ec.Latitude = lat
								ec.Longitude = long
							}
							resultData = append(resultData, ec)
						}
					}

					each.InputData = resultData
					dataRows = append(dataRows, each)
				}
			}

			// form detail data ------------------------------------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, err := ctr.formMod.GetFormRow(fieldForm)

			var res objects.InputFormMapRes
			res.FormID = formID
			res.FormName = getForm.Name
			res.FormDescription = getForm.Description
			res.PeriodStartDate = getForm.PeriodStartDate
			res.PeriodEndDate = getForm.PeriodEndDate
			res.FieldLabel = fieldLabels
			res.FieldData = dataRows

			// totalPage := 0
			// if limit > 0 {
			// 	totalPage = len(getAllData) / limit
			// 	if (len(getAllData) % limit) > 0 {
			// 		totalPage = totalPage + 1
			// 	}
			// }

			var pagingRes objects.DataRows
			pagingRes.TotalRows = len(getData)
			// pagingRes.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  pagingRes,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

// save data with join company
func (ctr *inputFormController) FieldSaveDataWithCompany(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	companyID := 0
	if len(claims) >= 5 {
		companyID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("companyID :::", companyID)
	}

	fmt.Println("companyID :::", companyID)

	var reqData objects.FormData
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

	if len(reqData.FieldData) > 0 {
		reqData.UserID = userID
		_, err = ctr.inputForm.InsertFormDataWithOrganization(reqData, companyID)
		if err != nil {
			fmt.Println("InsertFormData", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		// check form organization invite
		var whrCompany tables.JoinFormCompanies
		whrCompany.FormID = reqData.FormID
		whrCompany.OrganizationID = companyID
		checkOrgInvite, err := ctr.formMod.GetFormCompanyInviteRow(whrCompany, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		// check uota & the sharing quota is true
		var quotaCurrent = 0
		var organizationID = 0
		if checkOrgInvite.ID >= 1 && checkOrgInvite.IsQuotaSharing == true {
			// quota other company (sharing form quota)
			// InsertInjuryPlan
			var whr objects.SubsPlan
			whr.OrganizationID = checkOrgInvite.OrganizationID
			checkQuota, _ := ctr.subsMod.GetPlanRow(whr)

			fmt.Println("checkQuota.QuotaCurrent>>>>>", checkQuota.QuotaCurrent)
			if checkQuota.IsBlocked == true {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Your data respon isn't saved",
				})
				return
			}
			organizationID = checkOrgInvite.OrganizationID
			quotaCurrent = checkQuota.QuotaCurrent

		} else {
			// quota company milik sendiri (author form)
			var whreFrm tables.FormOrganizations
			whreFrm.FormID = reqData.FormID
			getFormOrg, _ := ctr.formMod.GetFormOrganization(whreFrm)

			// cheking quota zero
			// InsertInjuryPlan
			var whr objects.SubsPlan
			whr.OrganizationID = getFormOrg.OrganizationID
			checkQuota, _ := ctr.subsMod.GetPlanRow(whr)

			fmt.Println("checkQuota.QuotaCurrent>>>>>", checkQuota.QuotaCurrent)
			if checkQuota.IsBlocked == true {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Your data respon isn't saved",
				})
				return
			}
			organizationID = getFormOrg.OrganizationID
			quotaCurrent = checkQuota.QuotaCurrent

		}

		if quotaCurrent == 0 {

			var whreSett tables.Settings
			whreSett.Code = "injury_plan_quota"
			settingInjuryQuota, _ := ctr.settModel.GetSettingRow(whreSett)
			quotaValue, _ := strconv.Atoi(settingInjuryQuota.Value)

			var whrPlan objects.SubsPlan
			whrPlan.OrganizationID = organizationID
			checkPlan, _ := ctr.subsMod.GetPlanRow(whrPlan)

			var whrInjPlan objects.InjuryPlan
			whrInjPlan.OrganizationSubscriptionPlanID = checkPlan.ID
			checkInjuryPlan, _ := ctr.subsMod.GetInjuryPlanRows(whrInjPlan)

			fmt.Println("injury cek ------------------------------>", checkInjuryPlan, checkPlan.ID, len(checkInjuryPlan))

			// ke 2 kali nya tidak di auto add Quota
			if quotaValue > 0 && len(checkInjuryPlan) == 0 {
				// var dataPostInjury tables.SubsPlan
				// dataPostInjury.QuotaCurrent = quotaValue
				// dataPostInjury.QuotaTotal = quotaValue
				// addnjuryQuota, injuryRes, err := ctr.subsMod.UpdatePlan(getFormOrg.OrganizationID, dataPostInjury)
				// if err != nil {
				// 	fmt.Println("UpdatePlan dataPostInjury", err)
				// 	c.JSON(http.StatusBadRequest, gin.H{
				// 		"error": err,
				// 	})
				// 	return
				// }

				// if addnjuryQuota == true {
				var injuryData tables.InjuryPlan
				injuryData.OrganizationSubscriptionPlanID = checkPlan.ID
				injuryData.Quota = quotaValue
				ctr.subsMod.InsertInjuryPlan(injuryData)
				// }
			}
		}

		//insert submiting quota respon
		var dataPost tables.SubsPlan
		dataPost.RespondentCurrent = 1
		_, _, err = ctr.subsMod.UpdatePlanCurrent(organizationID, dataPost)
		if err != nil {
			fmt.Println("UpdatePlanCurrent", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Submission Sukses",
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data input is empty",
		})
		return
	}

}

func (ctr *inputFormController) DataFormDelete(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.FormDataDelete
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println("error val", err.Error())
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status": true,
			"error":  errorMessages,
		})
		return
	}

	formID := reqData.FormID

	if formID >= 1 {

		if len(reqData.FormDataIDs) >= 1 {
			deleteData := false
			for i := 0; i < len(reqData.FormDataIDs); i++ {
				deleteData, err = ctr.inputForm.DeleteFormData(reqData.FormDataIDs[i].ID, objects.FormData{FormID: formID, UserID: userID})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status": true,
						"error":  err,
					})
					return
				}
			}

			if deleteData == true {

				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Success delete submission data",
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed delete data",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data IDs is not availabe",
			})
			return
		}

	}
}

func (ctr *inputFormController) GenerateFormUserOrg(c *gin.Context) {

	getAllFormUser, err := ctr.formMod.GetFormUserToFormOrgRows(tables.JoinFormUsers{}, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data IDs is not availabe",
		})
		return
	}
	if len(getAllFormUser) >= 1 {

		for i := 0; i < len(getAllFormUser); i++ {

			var postData tables.FormUserOrganizations
			postData.FormUserID = getAllFormUser[i].ID
			postData.OrganizationID = getAllFormUser[i].OrganizationID
			ctr.formMod.GenerateFormUserOrg(postData)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Generate Form User Org finish",
		})
		return
	}
}

func (ctr *inputFormController) GenerateInputFormOrg(c *gin.Context) {

	getAllForms, err := ctr.formMod.GetFormRows(tables.Forms{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data IDs is not availabe",
		})
		return
	}
	if len(getAllForms) >= 1 {

		var exportData = make([][]string, len(getAllForms)+1)

		for i := 0; i < len(getAllForms); i++ {
			fmt.Println("getAllForms[i].ID ::", getAllForms[i].ID)

			selectInputData, _ := ctr.inputForm.GetInputFormRows(getAllForms[i].ID, tables.InputForms{}, "", objects.Paging{})
			// fmt.Println("selectInputData >>", len(getAllForms), getAllForms[i].ID, len(selectInputData))
			var rowData = []string{strconv.Itoa(getAllForms[i].ID), strconv.Itoa(len(selectInputData))}

			exportData[i] = rowData

			if len(selectInputData) >= 1 {
				for j := 0; j < len(selectInputData); j++ {

					var whr tables.JoinFormUsers
					whr.FormID = getAllForms[i].ID
					whr.UserID = selectInputData[j].UserID
					getCompany, err := ctr.formMod.GetFormUserToOrganizationRow(whr, "")
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"message": "Data GetFormUserToOrganizationRow is not availabe",
						})
						return
					}

					if getCompany.OrganizationID >= 1 {
						var postData tables.InputFormOrganizations
						postData.InputFormID = selectInputData[j].ID
						postData.OrganizationID = getCompany.OrganizationID
						postData.FormID = getAllForms[i].ID
						ctr.inputForm.InsertInputFormOrgData(postData)
					}
					// os.Exit(0)
				}
			}
		}

		fileLocation := "file/export_input_data.csv" // local path file location

		f, err := os.Create(fileLocation)
		// defer f.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		w := csv.NewWriter(f)
		w.UseCRLF = false
		err = w.WriteAll(exportData) // calls Flush internally
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data InsertInputFormOrgData success",
		})
		return

	}
}

func (ctr *inputFormController) GenerateInputFormOrgNoCopy(c *gin.Context) {

	getAllForms, err := ctr.formMod.GetAllFormRows(tables.Forms{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data IDs is not availabe",
		})
		return
	}
	if len(getAllForms) >= 1 {

		var exportData = make([][]string, len(getAllForms)+1)

		for i := 0; i < len(getAllForms); i++ {
			fmt.Println("getAllForms[i].ID ::", getAllForms[i].ID)

			selectInputData, _ := ctr.inputForm.GetInputFormRows(getAllForms[i].ID, tables.InputForms{}, "", objects.Paging{})
			// fmt.Println("selectInputData >>", len(getAllForms), getAllForms[i].ID, len(selectInputData))
			var rowData = []string{strconv.Itoa(getAllForms[i].ID), strconv.Itoa(len(selectInputData))}

			exportData[i] = rowData

			if len(selectInputData) >= 1 {
				for j := 0; j < len(selectInputData); j++ {

					var whr tables.JoinFormUsers
					whr.FormID = getAllForms[i].ID
					whr.UserID = selectInputData[j].UserID
					getCompany, err := ctr.formMod.GetFormUserToOrganizationRow(whr, "")
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"message": "Data GetFormUserToOrganizationRow is not availabe",
						})
						return
					}

					if getCompany.OrganizationID >= 1 {
						var postData tables.InputFormOrganizations
						postData.InputFormID = selectInputData[j].ID
						postData.OrganizationID = getCompany.OrganizationID
						postData.FormID = getAllForms[i].ID
						ctr.inputForm.InsertInputFormOrgData(postData)

					} else if getCompany.OrganizationID <= 0 {

						getTeam, _ := ctr.attMod.GetTeamByRespondent(selectInputData[j].UserID)

						if len(getTeam) > 0 {
							for i := 0; i < len(getTeam); i++ {
								getFormTeam, _ := ctr.attMod.GetFormTeam(getAllForms[i].ID, getTeam[i].TeamID)

								if getFormTeam.OrganizationID >= 1 {

									var postDataTeam tables.InputFormOrganizations
									postDataTeam.InputFormID = selectInputData[j].ID
									postDataTeam.OrganizationID = getFormTeam.OrganizationID
									postDataTeam.FormID = getAllForms[i].ID
									ctr.inputForm.InsertInputFormOrgData(postDataTeam)
								}
							}
						}

					}
					// os.Exit(0)
				}
			}
		}

		fileLocation := "file/export_input_data.csv" // local path file location

		f, err := os.Create(fileLocation)
		// defer f.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		w := csv.NewWriter(f)
		w.UseCRLF = false
		err = w.WriteAll(exportData) // calls Flush internally
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data InsertInputFormOrgData success",
		})
		return

	}
}

func (ctr *inputFormController) GenerateInputFormOrgLatest(c *gin.Context) {

	getAllForms, err := ctr.formMod.GetAllFormRows(tables.Forms{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data IDs is not availabe",
		})
		return
	}
	if len(getAllForms) >= 1 {

		for i := 0; i < len(getAllForms); i++ {

			getCompany, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: getAllForms[i].ID})

			if getCompany.OrganizationID >= 1 {
				selectInputData, _ := ctr.inputForm.GetInputFormRows(getAllForms[i].ID, tables.InputForms{}, "", objects.Paging{})

				if len(selectInputData) >= 1 {
					for j := 0; j < len(selectInputData); j++ {

						var postData tables.InputFormOrganizations
						postData.InputFormID = selectInputData[j].ID
						postData.OrganizationID = getCompany.OrganizationID
						postData.FormID = getAllForms[i].ID
						ctr.inputForm.InsertInputFormOrgData(postData)

					}
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data InsertInputFormOrgData success",
		})
		return

	}
}

func makeAPICall() error {
	// Create a context with no timeout
	ctx := context.Background()

	req, err := http.NewRequest("POST", "http://localhost:8088/data/generateinputform", nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Process the response
	// ...

	return nil
}
