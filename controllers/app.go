package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"snapin-form/config"
	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"

	"github.com/gin-gonic/gin"
)

// interface
type AppController interface {
	Home(c *gin.Context)
	FormList(c *gin.Context)
	FormPerformance(c *gin.Context)
	ProjectList(c *gin.Context)
	SubmissionForm(c *gin.Context)
	SubmissionFormData(c *gin.Context)
	SubmissionFormDataField(c *gin.Context)
	SubmissionEditRequest(c *gin.Context)
	SubmissionFormOTPChecking(c *gin.Context)

	//admin apps
	HomeAdmin(c *gin.Context)
	HomeAdminContent(c *gin.Context)
	SubmissionDetailUser(c *gin.Context)
}

type appController struct {
	formMod      models.FormModels
	formFieldMod models.FormFieldModels
	formOtpMod   models.FormOtpModels
	ftMod        models.FieldTypeModels
	ruleMod      models.RuleModels
	helper       helpers.Helper
	userMod      models.UserModels
	projectMod   models.ProjectModels
	inputForm    models.InputFormModels
	attendMod    models.AttendanceModels
	compMod      models.CompaniesModels
	subsMod      models.SubsModels
	conf         config.Configurations
	settingMod   models.SettingModels
	pgErr        *pgconn.PgError
	shrtMod      models.ShortenUrlModels
}

func NewAppController(formModel models.FormModels, formFieldModel models.FormFieldModels, ftModel models.FieldTypeModels, ruleMod models.RuleModels, helper helpers.Helper, projectModel models.ProjectModels, userModel models.UserModels, inputForm models.InputFormModels, attendModel models.AttendanceModels, compModel models.CompaniesModels, formOtpModel models.FormOtpModels, settingModels models.SettingModels, configs config.Configurations, subsModel models.SubsModels, shorModel models.ShortenUrlModels) AppController {
	return &appController{
		formMod:      formModel,
		ftMod:        ftModel,
		formFieldMod: formFieldModel,
		ruleMod:      ruleMod,
		helper:       helper,
		projectMod:   projectModel,
		inputForm:    inputForm,
		userMod:      userModel,
		attendMod:    attendModel,
		compMod:      compModel,
		formOtpMod:   formOtpModel,
		subsMod:      subsModel,
		conf:         configs,
		settingMod:   settingModels,
		shrtMod:      shorModel,
	}
}

func (ctr *appController) Home(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	//get my form
	var fields tables.FormUsers
	fields.UserID = userAppsID
	getForms, err := ctr.formMod.GetFormUserUnionTeamRows(fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	totalRespon := 0
	averageAtt := 0
	var averageAttFloat float64
	if len(getForms) > 0 {

		totalFormRequired := 0
		totalUserAttToday := 0
		for i := 0; i < len(getForms); i++ {

			//count respon all form
			var whereInForm tables.InputForms
			whereInForm.UserID = userAppsID
			whereIFToday := " TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "

			getDataRespons, err := ctr.inputForm.GetInputFormRows(getForms[i].FormID, whereInForm, whereIFToday, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			totalRespon += len(getDataRespons)

			// cek he/she is has attendance
			var whreAtt tables.Attendances
			whreAtt.UserID = userAppsID
			whreAtt.FormID = getForms[i].FormID
			whereToday := " TO_CHAR(attendances.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "

			getAtt, _ := ctr.attendMod.GetAttendanceRow(whreAtt, whereToday)
			if getAtt.ID > 0 && getAtt.GeometryOut != "" {
				totalUserAttToday += 1
			}

			// get form attendance required
			if getForms[i].IsAttendanceRequired {
				totalFormRequired += 1
			}
		}

		if totalFormRequired > 0 {
			avgF := (float64(totalUserAttToday) / float64(totalFormRequired)) * 100
			averageAtt, _ = strconv.Atoi(strconv.FormatFloat(avgF, 'f', 0, 64))

			strPerform := strconv.FormatFloat(avgF, 'f', 1, 64)
			averageAttFloat, _ = strconv.ParseFloat(strPerform, 1)
		}

	}

	// get user data
	var whreUser tables.Users
	whreUser.ID = userAppsID
	getUser, err := ctr.userMod.GetUserRow(whreUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var res objects.HomeApps
	res.UserID = getUser.ID
	res.UserName = getUser.Name
	res.UserAvatar = getUser.Avatar
	res.CompanyName = getUser.CompanyName
	res.TotalForm = len(getForms)
	res.TotalRespon = totalRespon
	res.TotalAttendance = averageAtt
	res.TotalAttendanceFloat = averageAttFloat

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    res,
	})
	return
}

func (ctr *appController) FormList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	respondenID, _ := strconv.Atoi(claims["id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

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

		if result.ID >= 1 {

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate

			var FFfields tables.FormFields
			FFfields.FormID = result.ID
			getFields, err := ctr.formFieldMod.GetFormFieldRows(FFfields)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			var getFormFields []objects.FormFields
			for i := 0; i < len(getFields); i++ {
				var each objects.FormFields
				each.ID = getFields[i].ID
				each.FieldTypeID = getFields[i].FieldTypeID
				each.Label = getFields[i].Label
				each.Description = getFields[i].Description
				each.Option = getFields[i].Option
				each.ConditionType = getFields[i].ConditionType
				each.UpperlowerCaseType = getFields[i].UpperlowerCaseType
				each.IsMultiple = getFields[i].IsMultiple
				each.IsRequired = getFields[i].IsRequired

				getFormFields = append(getFormFields, each)
			}
			res.FormFields = getFormFields

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

		var fields tables.FormUsers
		fields.UserID = respondenID
		results, err := ctr.formMod.GetFormUserUnionTeamRows(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.AppListForms
			for i := 0; i < len(results); i++ {

				//get total responden
				var whereInForm tables.InputForms
				whereInForm.UserID = respondenID
				whereToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
				getDataRespons, err := ctr.inputForm.GetInputFormRows(results[i].FormID, whereInForm, whereToday, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}

				var each objects.AppListForms
				each.ID = results[i].FormID
				each.Name = results[i].Name
				each.Description = results[i].Description
				each.ProfilePic = results[i].ProfilePic
				each.PeriodStartDate = results[i].PeriodStartDate
				each.PeriodEndDate = results[i].PeriodEndDate
				each.TotalResponden = 0
				each.TotalRespon = len(getDataRespons)
				each.FormStatusID = results[i].FormStatusID
				each.FormStatus = results[i].FormStatus
				each.AttendanceIn = results[i].AttendanceIn
				each.AttendanceOut = results[i].AttendanceOut
				each.IsAttendanceRequired = results[i].IsAttendanceRequired
				each.AttendanceOverdate = results[i].IsAttendanceOverdate
				each.IsAttendanceRadius = results[i].IsAttendanceRadius

				if results[i].AttendanceIn == "" && results[i].AttendanceOut == "" {
					fmt.Println("------- in 1")

					each.IsButton = "check_in"

					// cek chekout yg kosong sebelum hari ini (kurang dari sama dengan today) dari batas tanggal overdate At
					checkCheckoutBefore, err := ctr.attendMod.GetLastAttendanceOverdate(results[i].FormID, respondenID)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  false,
							"error":   err,
							"message": "Data is available",
						})
						return
					}

					if results[i].IsAttendanceOverdate == true && checkCheckoutBefore.ID > 0 && checkCheckoutBefore.AttendanceOut == "" {
						fmt.Println("------- sub in 1 :::", results[i].FormID, respondenID, "---:::", checkCheckoutBefore.AttendanceOut)

						// jika absen sebelumnya belum checkout maka button akan terus checkout
						each.IsButton = "check_out"
					}

				} else if results[i].AttendanceIn != "" && results[i].AttendanceOut == "" {
					fmt.Println("------- in 2")
					each.IsButton = "check_out"

				} else if results[i].AttendanceIn != "" && results[i].AttendanceOut != "" {
					fmt.Println("------- in 3")

					each.IsButton = "finish"

				} else if results[i].IsAttendanceRequired == false {
					each.IsButton = ""
				}

				each.IsActivePeriod = false
				checkCompanyActivePeriod, _ := ctr.subsMod.GetPlanPeriodRow(objects.SubsPlan{OrganizationID: results[i].OrganizationID})
				if checkCompanyActivePeriod.TotalPeriodDays >= 1 && checkCompanyActivePeriod.PeriodRemain >= 1 {
					each.IsActivePeriod = true

					// if quota habis
					if checkCompanyActivePeriod.QuotaCurrent <= 0 {
						each.IsActivePeriod = false
					}
				}

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

func (ctr *appController) FormPerformance(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	ID := c.Param("formid")
	formID, _ := strconv.Atoi(ID)
	totalDataUpdated := 0

	if formID > 0 {
		var whre tables.Forms
		whre.ID = formID
		results, err := ctr.formMod.GetFormRows(whre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "GetFormUserRows",
				"error": err,
			})
			return
		}
		var whereInForm tables.InputForms
		whereInForm.UserID = userID
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
		dataUpdated, err := ctr.inputForm.GetUpdatedData(formID, userID, whereDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "GetUpdatedData",
				"error": err,
			})
			return
		}
		fmt.Println(dataUpdated)
		for i := 0; i < len(dataUpdated); i++ {
			totalDataUpdated = dataUpdated[i].UpdatedCount
		}
		if len(results) > 0 {
			var res []objects.FormPerformance
			for i := 0; i < len(results); i++ {
				//get total respon
				var whereInForm tables.InputForms
				whereInForm.UserID = userID
				whereToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
				getDataRespons, err := ctr.inputForm.GetInputFormRows(results[i].ID, whereInForm, whereToday, objects.Paging{})
				getDataResponsUnscopes, err := ctr.inputForm.GetInputFormUnscopedRows(results[i].ID, whereInForm, whereToday, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":   "GetInputFormRows",
						"error": err,
					})
					return
				}

				// total average
				var wFormUser tables.JoinFormUsers
				wFormUser.FormID = results[i].ID
				getUserForm, _ := ctr.formMod.GetFormUserRows(wFormUser, "")

				var whereAll tables.InputForms
				getAllDataSubmission, err := ctr.inputForm.GetInputFormRows(results[i].ID, whereAll, whereToday, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg":   "GetInputFormRows all",
						"error": err,
					})
					return
				}

				totalAverage := 0.0
				if len(getDataRespons) > 0 {
					fmt.Println("len(getDataRespons)", len(getDataRespons), float64(len(getAllDataSubmission)), float64(len(getDataRespons)))
					totalAverage = (float64(len(getAllDataSubmission)) / float64(len(getUserForm)))
				}

				//last update
				lastUpdate := ""
				if len(getDataRespons) > 0 {
					lastUpdate = (getDataRespons[0].CreatedAt).String()
				}

				//info target
				infoTargetColor := "#FF0000"
				infoTarget := "Belum mencapai target"
				if len(getDataRespons) >= results[i].SubmissionTargetUser {
					infoTarget = "Sudah mencapai target"
					infoTargetColor = "#000000"
				}

				// info Average
				infoAverageColor := "#FF0000"
				infoAverage := "Kinerjamu di bawah rata-rata tim"
				if float64(len(getDataRespons)) > totalAverage {
					infoAverage = "Kinerjamu di atas rata-rata tim"
					infoAverageColor = "#000000"
				} else if float64(len(getDataRespons)) == totalAverage {
					infoAverage = "Kinerjamu sama dengan rata-rata tim"
					infoAverageColor = "#000000"
				}

				// process absen
				proData := []objects.ProcessData{}
				if results[i].IsAttendanceRequired {

					// check in today
					var whreIn tables.Attendances
					whreIn.FormID = results[i].ID
					whreStrIn := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					getAttIn, _ := ctr.attendMod.GetAttendanceRow(whreIn, whreStrIn)
					checkIn := false
					if getAttIn.ID > 0 {
						checkIn = true
					}

					// checkout today
					var whreOut tables.Attendances
					whreOut.FormID = results[i].ID
					whreStrOut := "TO_CHAR(attendances.attendance_out::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					getAttOut, _ := ctr.attendMod.GetAttendanceRow(whreOut, whreStrOut)
					checkOut := false
					if getAttOut.ID > 0 {
						checkOut = true
					}

					// submission form today
					hasSubmission := false
					if lastUpdate != "" {
						hasSubmission = true
					}

					proData = []objects.ProcessData{
						{ProcessName: "Check-in", Status: checkIn},
						{ProcessName: "Isi Form", Status: hasSubmission},
						{ProcessName: "Check-out", Status: checkOut},
					}
				}

				var each objects.FormPerformance
				each.ID = results[i].ID
				each.Name = results[i].Name
				each.ProfilePic = results[i].ProfilePic
				each.PeriodStartDate = results[i].PeriodStartDate
				each.PeriodEndDate = results[i].PeriodEndDate
				each.LastUpdate = lastUpdate
				each.TotalRespon = len(getDataRespons)
				each.TotalTarget = results[i].SubmissionTargetUser
				each.TotalAverage = totalAverage
				each.InfoTarget = infoTarget
				each.InfoAverage = infoAverage
				each.InfoTargetColor = infoTargetColor
				each.InfoAverageColor = infoAverageColor
				each.ProgressData = proData
				each.TotalDeletedData = len(dataDeleted)
				each.TotalUpdateData = totalDataUpdated
				each.AllTotalRespon = strconv.Itoa(len(getDataRespons)) + "/" + strconv.Itoa(len(getDataResponsUnscopes))

				var whreStr = `where f.user_id = ` + strconv.Itoa(userID) + ` AND TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')`
				getDataHours, err := ctr.inputForm.GetDataHours(results[i].ID, whreStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
				each.HoursData = getDataHours

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

	} else {
		results, err := ctr.formMod.GetFormUserRespondenRows(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":   "GetFormUserRows",
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.FormPerformance
			for i := 0; i < len(results); i++ {

				var whreForm tables.Forms
				whreForm.ID = results[i].FormID
				getForm, _ := ctr.formMod.GetFormRow(whreForm)

				var each objects.FormPerformance
				each.ID = getForm.ID
				each.Name = getForm.Name
				each.ProfilePic = getForm.ProfilePic
				each.PeriodStartDate = getForm.PeriodStartDate
				each.PeriodEndDate = getForm.PeriodEndDate

				//get total respon
				var whereInForm tables.InputForms
				whereInForm.UserID = userID
				getDataRespons, err := ctr.inputForm.GetInputFormRows(getForm.ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
				each.TotalRespon = len(getDataRespons)
				each.TotalTarget = 100
				each.TotalAverage = 0

				var proData = []objects.ProcessData{
					{ProcessName: "Check-in", Status: true},
					{ProcessName: "Isi Form", Status: false},
					{ProcessName: "Check-out", Status: false},
				}
				each.ProgressData = proData

				var whreStr = `where f.user_id = ` + strconv.Itoa(userID)
				getDataHours, err := ctr.inputForm.GetDataHours(getForm.ID, whreStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
				each.HoursData = getDataHours

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

func (ctr *appController) ProjectList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	ID := c.Param("id")
	iID, _ := strconv.Atoi(ID)

	if iID > 0 {

		var fields tables.Forms
		fields.ID = iID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {

			fmt.Println("result.ProfilePic", result.ProfilePic)
			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate

			var FFfields tables.FormFields
			FFfields.FormID = result.ID
			getFields, err := ctr.formFieldMod.GetFormFieldRows(FFfields)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			var getFormFields []objects.FormFields
			for i := 0; i < len(getFields); i++ {
				var each objects.FormFields
				each.ID = getFields[i].ID
				each.FieldTypeID = getFields[i].FieldTypeID
				each.Label = getFields[i].Label
				each.Description = getFields[i].Description
				each.Option = getFields[i].Option
				each.ConditionType = getFields[i].ConditionType
				each.UpperlowerCaseType = getFields[i].UpperlowerCaseType
				each.IsMultiple = getFields[i].IsMultiple
				each.IsRequired = getFields[i].IsRequired

				getFormFields = append(getFormFields, each)
			}
			res.FormFields = getFormFields

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

		results, err := ctr.projectMod.GetProjectInForms(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.Projects
			for i := 0; i < len(results); i++ {
				var each objects.Projects
				each.ID = results[i].ID
				each.Name = results[i].Name
				each.Description = results[i].Description

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

func (ctr *appController) SubmissionForm(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	//get my form
	var wherefld tables.FormUsers
	wherefld.UserID = userAppsID
	getForms, err := ctr.formMod.GetFormUserUnionTeamRows(wherefld)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "GetFormUserRows",
			"error": err,
		})
		return
	}

	if len(getForms) > 0 {

		var res []objects.SubmissionForm
		for i := 0; i < len(getForms); i++ {
			var ech objects.SubmissionForm
			ech.FormID = getForms[i].FormID
			ech.FormName = getForms[i].Name
			ech.ProfilePic = getForms[i].ProfilePic

			var whreIF tables.InputForms
			whreIF.UserID = userAppsID

			whreStr := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') "
			getSubm, _ := ctr.inputForm.GetInputFormRows(getForms[i].FormID, whreIF, whreStr, objects.Paging{})
			ech.TotalSubmission = len(getSubm)

			res = append(res, ech)
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

func (ctr *appController) SubmissionFormData(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2
		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFiel_dRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			currentDate := time.Now().Format("2006-01-02")

			var whereInForm tables.InputForms
			whereInForm.UserID = userAppsID

			var buffer bytes.Buffer
			var whreStr = ``
			buffer.WriteString("to_char(if.created_at, 'yyyy-mm-dd') =  '" + currentDate + "'")
			whreStr = buffer.String()

			getData, err := ctr.inputForm.GetInputFormUnscopedRows(formID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"status": false,
					"error":  err.Error(),
				})
				return
			}

			fmt.Println("getData :::", getData)

			var dataRows []objects.InputFieldData
			for i := 0; i < len(getData); i++ {

				deletedValue, _ := getData[i].DeletedAt.Value()
				if deletedValue == nil {
					var each objects.InputFieldData
					each.ID = getData[i].ID
					each.UserID = getData[i].UserID
					each.UserName = getData[i].UserName
					if getData[i].UpdatedCount == 0 {
						each.StatusData = ""
					} else {
						each.StatusData = strconv.Itoa(getData[i].UpdatedCount) + "x edit"
					}
					each.CreatedAt = getData[i].CreatedAt.Format("2006-01-02 15:04")
					each.UpdatedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					dataRows = append(dataRows, each)
				} else {

					userDeleting, _ := ctr.userMod.GetUserRow(tables.Users{ID: getData[i].DeletedBy})

					var each objects.InputFieldData
					each.ID = getData[i].ID
					each.StatusData = "Dihapus"
					each.Note = "Data telah dihapus oleh " + userDeleting.Name
					each.DeletedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					dataRows = append(dataRows, each)
				}
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
			res.AuthorPhoneNumber = getForm.UserPhone
			res.CreatedAt = getForm.CreatedAt.Format("2006-01-02 15:04")
			res.UpdatedAt = getForm.UpdatedAt.Format("2006-01-02 15:04")
			res.FieldData = dataRows

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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
		})
		return

	}
}

func (ctr *appController) SubmissionFormDataField__(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// userAppsID, _ := strconv.Atoi(claims["id"].(string))
	// userAppsID = 17
	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	sfieldDataID := c.Param("field_data_id")
	fieldDataID, err := strconv.Atoi(sfieldDataID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2
		getQuests, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		fmt.Println("len(getQuests) .........", len(getQuests))
		if len(getQuests) > 0 {

			//get question
			var quest []objects.FormFields

			for i := 0; i < len(getQuests); i++ {

				//get answer
				var whereInForm tables.InputForms
				whereInForm.ID = fieldDataID
				// whereInForm.UserID = userAppsID
				// var buffer bytes.Buffer
				// whreStr = buffer.String()

				fieldStrings := "coalesce(f" + strconv.Itoa(getQuests[i].ID) + ",'') as f"

				getData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, whereInForm, "")
				if err != nil {
					fmt.Println("err GetFormFieldRows----", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var eachj objects.FormFields
				// each.ID = getQuests[i].ID
				// each.FieldTypeID = getQuests[i].FieldTypeID
				// each.FieldTypeName = getQuests[i].FieldTypeName
				// each.Label = getQuests[i].Label
				// each.FieldData = getData[0][0]

				eachj.ID = getQuests[i].ID
				eachj.ParentID = getQuests[i].ParentID
				eachj.FieldTypeID = getQuests[i].FieldTypeID
				eachj.FieldTypeName = getQuests[i].FieldTypeName
				eachj.FormID = getQuests[i].FormID
				eachj.Label = getQuests[i].Label
				eachj.Description = getQuests[i].Description
				eachj.Option = getQuests[i].Option
				eachj.ConditionType = getQuests[i].ConditionType
				eachj.UpperlowerCaseType = getQuests[i].UpperlowerCaseType
				eachj.IsMultiple = getQuests[i].IsMultiple
				eachj.IsRequired = getQuests[i].IsRequired
				eachj.IsSection = getQuests[i].IsSection
				eachj.SectionColor = getQuests[i].SectionColor
				eachj.SortOrder = getQuests[i].SortOrder
				eachj.Image = getQuests[i].Image
				eachj.ConditionRulesID = getQuests[i].ConditionRuleID
				eachj.ConditionRuleValue1 = getQuests[i].Value1
				eachj.ConditionRuleValue2 = getQuests[i].Value2
				eachj.ConditionRuleMsg = getQuests[i].ErrMsg
				eachj.ConditionParentFieldID = getQuests[i].ConditionParentFieldID

				//group condition
				if getQuests[i].IsCondition {
					var fcs tables.FormFieldConditionRules
					fcs.FormFieldID = getQuests[i].ID
					stringFields := "condition_parent_field_id is not null"
					getRules, _ := ctr.ruleMod.GetFormFieldRuleRows(fcs, stringFields)

					var condRules []objects.FormFieldConditionRules
					for i := 0; i < len(getRules); i++ {

						var each objects.FormFieldConditionRules

						if getRules[i].ConditionParentFieldID > 0 {

							each.FormFieldID = getRules[i].FormFieldID
							each.ConditionParentFieldID = getRules[i].ConditionParentFieldID
							each.ConditionRuleID = getRules[i].ConditionRuleID
							each.Value1 = getRules[i].Value1
							each.Value2 = getRules[i].Value2
							each.ConditionAllRight = getRules[i].ConditionAllRight
							each.ErrMsg = getRules[i].ErrMsg

						}
						condRules = append(condRules, each)

					}
					eachj.Conditions = condRules

				}
				eachj.FieldData = getData[0][0]

				quest = append(quest, eachj)
				// row data -------------------------------------------------------------------------------------------------------
			}

			// form detail data ------------------------------------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			var res objects.InputFormDetail
			res.FormID = formID
			res.FormName = getForm.Name
			res.FormDescription = getForm.Description
			res.PeriodStartDate = getForm.PeriodStartDate
			res.PeriodEndDate = getForm.PeriodEndDate
			// res.FieldLabelOnData = quest

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

func (ctr *appController) SubmissionFormDataField(c *gin.Context) {
	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	sfieldDataID := c.Param("field_data_id")
	fieldDataID, err := strconv.Atoi(sfieldDataID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var fields tables.FormFields
		fields.FormID = formID
		// fields.FieldTypeID = -2
		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFieldRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var res []objects.FormField3Groups
		if len(result) > 0 {

			for i := 0; i < len(result); i++ {
				if result[i].FieldTypeID == 0 && result[i].IsSection == false {
					fmt.Println("fieldDataID:: ------------------------>", fieldDataID)

					var each objects.FormField3Groups
					// its field group
					each.ID = result[i].ID
					each.IsGroup = true
					each.Label = result[i].Label
					each.Description = result[i].Description
					each.SortOrder = result[i].SortOrder
					each.Image = result[i].Image

					if result[i].IsCondition {
						var fcs tables.FormFieldConditionRules
						fcs.FormFieldID = result[i].ID
						stringFields := "condition_parent_field_id is not null"
						getRules, err := ctr.ruleMod.GetFormFieldRuleRows(fcs, stringFields)
						if err != nil {
							fmt.Println("------ERR :::", err)
							return
						}

						var condRules []objects.FormFieldConditionRules
						for k := 0; k < len(getRules); k++ {

							var eachK objects.FormFieldConditionRules
							if getRules[k].ConditionParentFieldID > 0 {

								eachK.FormFieldID = getRules[k].FormFieldID
								eachK.ConditionParentFieldID = getRules[k].ConditionParentFieldID
								eachK.ConditionRuleID = getRules[k].ConditionRuleID
								eachK.Value1 = getRules[k].Value1
								eachK.Value2 = getRules[k].Value2
								eachK.ConditionAllRight = getRules[k].ConditionAllRight
								eachK.ErrMsg = getRules[k].ErrMsg

							}
							condRules = append(condRules, eachK)

						}

						each.Conditions = condRules
					}

					//form fileds
					var getFieldres []objects.FormFields
					var whre tables.FormFields
					whre.ParentID = result[i].ID
					whre.FormID = formID
					getFields, err := ctr.formFieldMod.GetFormFieldRows(whre)
					if err != nil {
						fmt.Println("------ ERR GetFormFieldRows:::", err)
						return
					}

					if len(getFields) > 0 {

						for j := 0; j < len(getFields); j++ {

							var eachj objects.FormFields
							eachj.ID = getFields[j].ID
							eachj.ParentID = getFields[j].ParentID
							eachj.FieldTypeID = getFields[j].FieldTypeID
							eachj.FormID = getFields[j].FormID
							eachj.Label = getFields[j].Label
							eachj.Description = getFields[j].Description
							eachj.Option = getFields[j].Option
							eachj.ConditionType = getFields[j].ConditionType
							eachj.UpperlowerCaseType = getFields[j].UpperlowerCaseType
							eachj.IsMultiple = getFields[j].IsMultiple
							eachj.IsRequired = getFields[j].IsRequired
							eachj.IsSection = getFields[j].IsSection
							eachj.SectionColor = getFields[j].SectionColor
							eachj.SortOrder = getFields[j].SortOrder
							eachj.Image = getFields[j].Image
							eachj.ConditionRulesID = getFields[j].ConditionRuleID
							eachj.ConditionRuleValue1 = getFields[j].Value1
							eachj.ConditionRuleValue2 = getFields[j].Value2
							eachj.ConditionRuleMsg = getFields[j].ErrMsg
							eachj.ConditionParentFieldID = getFields[j].ConditionParentFieldID

							//group condition ----------
							if getFields[j].IsCondition {
								var fcs tables.FormFieldConditionRules
								fcs.FormFieldID = getFields[j].ID
								stringFields := "condition_parent_field_id is not null"
								getRules, _ := ctr.ruleMod.GetFormFieldRuleRows(fcs, stringFields)

								var condRules []objects.FormFieldConditionRules
								for i := 0; i < len(getRules); i++ {

									var each objects.FormFieldConditionRules

									if getRules[i].ConditionParentFieldID > 0 {

										each.FormFieldID = getRules[i].FormFieldID
										each.ConditionParentFieldID = getRules[i].ConditionParentFieldID
										each.ConditionRuleID = getRules[i].ConditionRuleID
										each.Value1 = getRules[i].Value1
										each.Value2 = getRules[i].Value2
										each.ConditionAllRight = getRules[i].ConditionAllRight
										each.ErrMsg = getRules[i].ErrMsg

									}
									condRules = append(condRules, each)

								}
								eachj.Conditions = condRules
							}

							//get answer---------------
							var whereInForm tables.InputForms
							whereInForm.ID = fieldDataID
							fieldStrings := "coalesce(f" + strconv.Itoa(getFields[j].ID) + ",'') as f"

							getData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, whereInForm, "")
							if err != nil {
								fmt.Println("err GetFormFieldRows----", err)
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							if len(getData) <= 0 {
								if err != nil {
									fmt.Println("err len(getData)----", err)
									c.JSON(http.StatusBadRequest, gin.H{
										"message": "Data tidak tersedia atau field data ID salah",
										"error":   err,
										"status":  false,
									})
									return
								}
							}
							eachj.FieldData = getData[0][0]

							// custom option
							getCustomOption, err := ctr.inputForm.GetInputFormCustomAnswerRow(formID, getFields[j].ID, fieldDataID)
							eachj.CustomOption = getCustomOption.CustomAnswer

							getFieldres = append(getFieldres, eachj)
						}
					}

					each.FormFields = getFieldres

					res = append(res, each)

				} else {
					fmt.Println("else:: ------------------------>", fieldDataID)
					// tidak memiliki group
					// var getFieldres []objects.FormFields
					var whre tables.FormFields
					whre.ID = result[i].ID
					getFields, _ := ctr.formFieldMod.GetFormFieldRow(whre)

					if getFields.ID > 0 && getFields.ParentID == 0 {
						var each objects.FormField3Groups

						each.ID = getFields.ID
						each.ParentID = getFields.ParentID
						each.FieldTypeID = getFields.FieldTypeID
						each.FormID = getFields.FormID
						each.Label = getFields.Label
						each.Description = getFields.Description
						each.Option = getFields.Option
						each.ConditionType = getFields.ConditionType
						each.UpperlowerCaseType = getFields.UpperlowerCaseType
						each.IsMultiple = getFields.IsMultiple
						each.IsRequired = getFields.IsRequired
						each.IsSection = getFields.IsSection
						each.SectionColor = getFields.SectionColor
						each.SortOrder = getFields.SortOrder
						each.Image = getFields.Image
						each.ConditionRulesID = getFields.ConditionRuleID
						each.ConditionRuleValue1 = getFields.Value1
						each.ConditionRuleValue2 = getFields.Value2
						each.ConditionRuleMsg = getFields.ErrMsg
						each.ConditionParentFieldID = getFields.ConditionParentFieldID

						//group condition
						if getFields.IsCondition {
							var fcs tables.FormFieldConditionRules
							fcs.FormFieldID = getFields.ID
							stringFields := "condition_parent_field_id is not null"
							getRules, _ := ctr.ruleMod.GetFormFieldRuleRows(fcs, stringFields)

							var condRules []objects.FormFieldConditionRules
							for k := 0; k < len(getRules); k++ {

								var eachK objects.FormFieldConditionRules
								if getRules[k].ConditionParentFieldID > 0 {

									eachK.FormFieldID = getRules[k].FormFieldID
									eachK.ConditionParentFieldID = getRules[k].ConditionParentFieldID
									eachK.ConditionRuleID = getRules[k].ConditionRuleID
									eachK.Value1 = getRules[k].Value1
									eachK.Value2 = getRules[k].Value2
									eachK.ConditionAllRight = getRules[k].ConditionAllRight
									eachK.ErrMsg = getRules[k].ErrMsg

								}
								condRules = append(condRules, eachK)

							}
							each.Conditions = condRules
						}

						//get answer---------------
						var whereInForm tables.InputForms
						whereInForm.ID = fieldDataID
						fieldStrings := "coalesce(f" + strconv.Itoa(getFields.ID) + ",'') as f, id, user_id"

						getData, err := ctr.inputForm.GetInputDataRows(formID, fieldStrings, whereInForm, "")
						if err != nil {
							fmt.Println("err GetFormFieldRows----", err)
							c.JSON(http.StatusBadRequest, gin.H{
								"error": err,
							})
							return
						}
						each.FieldData = getData[0][0]

						// decrypt shortlink ----------------------------------------------
						if getFields.FieldTypeID == 10 || getFields.FieldTypeID == 18 || getFields.FieldTypeID == 20 {
							var value = getData[0][0]

							if value != "" && len(value) < 35 {
								words := strings.Split(value, "/")
								lastWord := words[len(words)-1]

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

								each.FieldData = linkReal.Data.URL
							}
						}

						fmt.Println("getData ::::", getData[0])

						// custom option
						getCustomOption, err := ctr.inputForm.GetInputFormCustomAnswerRow(formID, getFields.ID, fieldDataID)
						each.CustomOption = getCustomOption.CustomAnswer
						res = append(res, each)
					}
				}
			}
		}

		// form detail data ------------------------------------------------------------------------------
		var fieldForm tables.Forms
		fieldForm.ID = formID
		getForm, _ := ctr.formMod.GetFormRow(fieldForm)

		var resData objects.InputFormDetail
		resData.FormID = formID
		resData.FormName = getForm.Name
		resData.FormDescription = getForm.Description
		resData.PeriodStartDate = getForm.PeriodStartDate
		resData.PeriodEndDate = getForm.PeriodEndDate
		resData.FieldLabelOnData = res

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    resData,
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

func (ctr *appController) HomeAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	// get user data
	var fields tables.Users
	fields.ID = userAppsID
	getUser, err := ctr.userMod.GetUserRow(fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var res objects.HomeAdminApps
	res.UserID = getUser.ID
	res.UserName = getUser.Name
	res.UserAvatar = getUser.Avatar
	res.CompanyName = getUser.CompanyName

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    res,
	})
	return
}

func (ctr *appController) HomeAdminContent(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))
	roleAppsID, _ := strconv.Atoi(claims["role_id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	// get user data
	var whreUser tables.Users
	whreUser.ID = userAppsID
	getUser, err := ctr.userMod.GetUserRow(whreUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	// get comp
	var getComp tables.Organizations
	compName := getUser.CompanyName
	if roleAppsID == 1 {
		var whrComp tables.Organizations
		whrComp.CreatedBy = userAppsID
		whrComp.IsDefault = true
		getComp, _ = ctr.compMod.GetCompaniesRow(whrComp)

		compName = getComp.Name
	} else {
		var whrOrg objects.UserOrganizations
		whrOrg.UserID = userAppsID
		getUserComp, _ := ctr.compMod.GetUserCompaniesRow(whrOrg, "")

		compName = getUserComp.Name
	}

	var result []tables.FormAll
	if roleAppsID == 1 { // 1 is owner
		// SUPER ADMIN here
		var buffer bytes.Buffer
		var fields tables.FormOrganizationsJoin
		// fields.CreatedBy = userAppsID
		// fields.FormStatusID = 1
		fields.OrganizationID = getComp.ID

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" forms.form_status_id = 1")

		// buffer.WriteString(" forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userAppsID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		getForms, err := ctr.formMod.GetFormOwnerRows(fields, whereString, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		result = getForms

	} else {

		var buffer bytes.Buffer
		var fields tables.Forms
		fields.FormStatusID = 1

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userAppsID) + ") OR forms.created_by = " + strconv.Itoa(userAppsID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userAppsID) + "))")

		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		getForms, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		result = getForms
	}

	var totalPerform float64
	var totalRespons int

	if len(result) > 0 {
		var total float64

		for i := 0; i < len(result); i++ {

			// get total responden
			var whereFU tables.JoinFormUsers
			whereFU.FormID = result[i].ID
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
			whreStrToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			var whereInForm tables.InputForms
			getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStrToday, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			totalRespons += len(getDataRespons)

			//total performance
			var totalPerformance float64
			if result[i].SubmissionTargetUser > 0 && len(getResponden) > 0 {
				totalPerformance = ((float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)) * 100) / float64(len(getResponden))

			} else {
				totalPerformance = 0
			}

			total = total + totalPerformance
		}
		totalPerform = total

	}

	// totalPerorm := 0
	totalDataPerormFloat := totalPerform / float64(len(result))
	totalDataPerorm, _ := strconv.Atoi(strconv.FormatFloat(totalDataPerormFloat, 'f', 0, 64))
	if totalDataPerormFloat >= 100 {
		totalDataPerormFloat = 100
		totalDataPerorm = 100
	}
	strPerform := strconv.FormatFloat(totalDataPerormFloat, 'f', 1, 64)
	lastTotal, _ := strconv.ParseFloat(strPerform, 1)

	var adminRes objects.HomeAdminApps
	adminRes.UserID = getUser.ID
	adminRes.UserName = getUser.Name
	adminRes.UserAvatar = getUser.Avatar
	adminRes.CompanyName = compName
	adminRes.TotalPerformance = totalDataPerorm
	adminRes.TotalPerformanceFloat = lastTotal
	adminRes.TotalRespon = totalRespons

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    adminRes,
	})
	return

}

func (ctr *appController) SubmissionEditRequest(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.SubmissionForm
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

	if reqData.FormID < 1 {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  false,
			"message": "Data form ID is required",
		})
		return
	}

	// generate otp
	otpcode := helpers.EncodeToString(4)
	if ctr.conf.ENV_TYPE == "dev" {
		otpcode = "1234"
	}

	var post tables.FormOtps
	post.FormID = reqData.FormID
	post.UserID = userAppsID
	post.OtpCode = otpcode
	_, err = ctr.formOtpMod.InsertFormOtp(post)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	// get data respnse
	getRespondenData, _ := ctr.userMod.GetUserRow(tables.Users{ID: userAppsID})

	// send form author ------------------------------
	var whrFrm tables.Forms
	whrFrm.ID = reqData.FormID
	getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

	var whrUser tables.Users
	whrUser.ID = getFormData.CreatedBy
	getAuthor, _ := ctr.userMod.GetUserRow(whrUser)

	if getAuthor.ID >= 1 && ctr.conf.ENV_TYPE != "dev" {
		msg := "*" + getRespondenData.Name + "* mengajukan perubahan data submision pada form *" + getFormData.Name + "*. Berikut kode OTP untuk mengijinkan perubahan data " + otpcode
		sendWA := helpers.SendWA(getAuthor.Phone, msg)

		fmt.Println("sendWA author ---->", sendWA)
	}

	// send otp to admin form -----------------------
	var whr tables.JoinFormUsers
	whr.FormID = reqData.FormID
	whr.Type = "admin"
	getAdminInForm, err := ctr.formMod.GetFormUserRows(whr, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(getAdminInForm) >= 1 {

		if ctr.conf.ENV_TYPE != "dev" {

			for i := 0; i < len(getAdminInForm); i++ {

				msg := "*" + getRespondenData.Name + "* mengajukan perubahan data submision pada form *" + getFormData.Name + "*. Berikut kode OTP untuk mengijinkan perubahan data " + otpcode
				sendWA := helpers.SendWA(getAdminInForm[i].Phone, msg)

				fmt.Println("sendWA ---->", sendWA)
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "OTP sent successfully!",
	})
	return
}

func (ctr *appController) SubmissionFormOTPChecking(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.SubmissionFormOtp
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
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

	var whrOtp tables.FormOtps
	whrOtp.FormID = reqData.FormID
	whrOtp.UserID = userAppsID
	whrOtp.OtpCode = reqData.OtpCode
	checkOtp, _ := ctr.formOtpMod.GetFormOtp(whrOtp)

	if checkOtp.ID > 0 {
		//delete history otp
		var whrFormOtp tables.FormOtps
		whrFormOtp.FormID = reqData.FormID
		whrFormOtp.UserID = userAppsID
		_, err := ctr.formOtpMod.DeleteWhrFormOtp(whrFormOtp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "OTP you entered is correct!",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "OTP you entered is wrong!",
		})
		return
	}
}

func (ctr *appController) SubmissionDetailUser(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userAppsID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println(userAppsID)

	fID := c.Request.URL.Query().Get("form_id")
	formID, err := strconv.Atoi(fID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	uID := c.Request.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(uID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID >= 1 && userID >= 1 {

		var fields tables.FormFields
		fields.FormID = formID
		fields.FieldTypeID = -2
		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			fmt.Println("err GetFormFiel_dRows----", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			// row data -------------------------------------------------------------------------------------------------------
			// var fieldStrings []string
			currentDate := time.Now().Format("2006-01-02")

			var whereInForm tables.InputForms
			whereInForm.UserID = userID

			var buffer bytes.Buffer
			var whreStr = ``
			buffer.WriteString("to_char(if.created_at, 'yyyy-mm-dd') =  '" + currentDate + "'")
			whreStr = buffer.String()

			getData, err := ctr.inputForm.GetInputFormUnscopedRows(formID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				fmt.Println("err GetFormFieldRows----", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

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
						each.StatusData = ""
					} else {
						each.StatusData = strconv.Itoa(getData[i].UpdatedCount) + "x edit"
					}
					each.CreatedAt = getData[i].CreatedAt.Format("2006-01-02 15:04")
					each.UpdatedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					dataRows = append(dataRows, each)
				} else {
					userDeleting, _ := ctr.userMod.GetUserRow(tables.Users{ID: getData[i].DeletedBy})

					var each objects.InputFieldData
					each.ID = getData[i].ID
					each.StatusData = "Dihapus"
					each.Note = "Data telah dihapus oleh " + userDeleting.Name
					each.DeletedAt = getData[i].UpdatedAt.Format("2006-01-02 15:04")

					dataRows = append(dataRows, each)
				}
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
			res.FieldData = dataRows

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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
		})
		return

	}
}
