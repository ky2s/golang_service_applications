package controllers

import (
	"fmt"
	"net/http"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"
	"sort"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type SubscriptionController interface {
	HistoryBalanceSaldo(c *gin.Context)
	HistoryBalanceSaldoByDate(c *gin.Context)
}

type subscribtionController struct {
	formMod   models.FormModels
	inputForm models.InputFormModels
}

func NewSubsController(formModel models.FormModels, inputFormModel models.InputFormModels) SubscriptionController {
	return &subscribtionController{
		formMod:   formModel,
		inputForm: inputFormModel,
	}
}

func (ctr *subscribtionController) HistoryBalanceSaldo(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	// roleID, _ := strconv.Atoi(claims["role_id"].(string))
	formID, _ := strconv.Atoi(c.Param("form_id"))
	date := c.Request.URL.Query().Get("date")
	searchKeyWord := c.Request.URL.Query().Get("search")

	if formID == 0 {
		whereGroupStr := ""
		if searchKeyWord != "" {
			whereGroupStr = "f.name ilike '%" + searchKeyWord + "%'"
		} else {
			whereGroupStr = ""
		}

		var whereComp objects.HistoryBalanceSaldo
		whereComp.OrganizationID = organizationID
		getFormByOrgID, err := ctr.formMod.GetFormsOrganizationWithFilter(whereComp, whereGroupStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"status":  false,
				"message": "Error: Data user is not deleted",
			})
			return
		}
		if len(getFormByOrgID) > 0 {
			var totalSubmission int
			var totalDeleted int
			var totalUpdated int
			var totalBlast int

			var res []objects.HistoryBalanceSaldo
			for i := 0; i < len(getFormByOrgID); i++ {
				var whereInForm tables.InputForms
				var whreStr = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + date + "'"
				var whreStrg = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + date + "' AND if.deleted_at IS NULL"
				var whreString = "TO_CHAR(notification_histories.created_at, 'YYYY-MM-DD') = '" + date + "'"

				getDataAllSubmission, _ := ctr.inputForm.GetInputFormUnscopedRows(getFormByOrgID[i].FormID, whereInForm, whreStrg, objects.Paging{})
				getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(getFormByOrgID[i].FormID, whereInForm, whreStr)
				getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedDataWithDate(getFormByOrgID[i].FormID, whreStr)
				getDataBlastInfo, _ := ctr.formMod.GetBlastInfoData(getFormByOrgID[i].FormID, whreString)

				for i := 0; i < len(getDataSubmissionUpdated); i++ {
					totalUpdated = getDataSubmissionUpdated[i].UpdatedCount
				}
				totalSubmission = len(getDataAllSubmission)
				totalDeleted = len(getDataSubmissionDeleted)
				totalBlast = len(getDataBlastInfo)

				each := objects.HistoryBalanceSaldo{
					FormID:                 getFormByOrgID[i].FormID,
					OrganizationID:         getFormByOrgID[i].OrganizationID,
					FormName:               getFormByOrgID[i].FormName,
					FormImage:              getFormByOrgID[i].FormImage,
					OrganizationName:       getFormByOrgID[i].OrganizationName,
					TotalDeletedSubmission: totalDeleted,
					TotalUpdatedSubmission: totalUpdated,
					TotalSubmission:        totalSubmission,
					TotalBlast:             totalBlast,
				}
				each.TotalUsageBalance = each.TotalUpdatedSubmission + each.TotalSubmission + each.TotalBlast + each.TotalDeletedSubmission
				if each.TotalUsageBalance != 0 {
					res = append(res, each)
				}
			}
			sort.Slice(res, func(i, j int) bool {
				return res[i].TotalUsageBalance > res[j].TotalUsageBalance
			})
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return

		}
	} else {
		var whereForm tables.JoinFormUsers
		whereForm.FormID = formID
		whereStr := ""
		if searchKeyWord != "" {
			whereStr = "u.name ilike '%" + searchKeyWord + "%'"
		} else {
			whereStr = ""
		}
		getUserByForm, err := ctr.formMod.GetFormUserRows(whereForm, whereStr)
		// fmt.Println("getUserByForm", getUserByForm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"status":  false,
				"message": "Error: Data user is not deleted",
			})
			return
		}
		var resHeader []objects.HistorySaldoUse
		if len(getUserByForm) > 0 {
			for i := 0; i < len(getUserByForm); i++ {

				var whereInForm tables.InputForms
				whereInForm.UserID = getUserByForm[i].UserID
				var whreStr = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + date + "'"
				var whreStrg = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + date + "' AND if.deleted_at IS NULL"
				var whreString = "TO_CHAR(notification_histories.created_at, 'YYYY-MM-DD') = '" + date + "'"

				var totalUpdated int
				getDataAllSubmission, _ := ctr.inputForm.GetInputFormUnscopedRows(formID, whereInForm, whreStrg, objects.Paging{})
				getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(formID, whereInForm, whreStr)
				getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedData(formID, getUserByForm[i].UserID, whreStr)
				getDataBlastInfo, _ := ctr.formMod.GetBlastInfoDataUsers(formID, getUserByForm[i].UserID, whreString)

				for i := 0; i < len(getDataSubmissionUpdated); i++ {
					totalUpdated = getDataSubmissionUpdated[i].UpdatedCount
				}
				// fmt.Println(getUserByForm[i].UserID)
				each := objects.HistorySaldoUse{
					UserID:                 getUserByForm[i].UserID,
					UserName:               getUserByForm[i].UserName,
					UserPhone:              getUserByForm[i].Phone,
					UserImage:              getUserByForm[i].UserImage,
					TotalSubmission:        len(getDataAllSubmission),
					TotalDeletedSubmission: len(getDataSubmissionDeleted),
					TotalUpdatedSubmission: totalUpdated,
					TotalBlast:             len(getDataBlastInfo),
				}
				each.TotalUsageBalance = each.TotalSubmission + each.TotalUpdatedSubmission + each.TotalBlast + each.TotalDeletedSubmission
				if each.TotalUsageBalance != 0 {
					resHeader = append(resHeader, each)
				}
				// var each objects.HistorySaldoUse
				// each.UserID = getUserByForm[i].UserID
				// each.UserName = getUserByForm[i].UserName
				// each.UserPhone = getUserByForm[i].Phone
				// each.UserImage = getUserByForm[i].UserImage
				// each.TotalSubmission = len(getDataAllSubmission)
				// each.TotalDeletedSubmission = len(getDataSubmissionDeleted)
				// each.TotalUpdatedSubmission = totalUpdated
				// each.TotalBlast = len(getDataBlastInfo)
				// each.TotalUsageBalance = each.TotalSubmission + each.TotalUpdatedSubmission + each.TotalBlast + each.TotalDeletedSubmission

				// resHeader = append(resHeader, each)
			}
			sort.Slice(resHeader, func(i, j int) bool {
				return resHeader[i].TotalUsageBalance > resHeader[j].TotalUsageBalance
			})
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    resHeader,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true,
				"message": "Data is not available",
				"data":    make([]objects.HistorySaldoUse, 0),
			})
			return
		}
	}
}

func (ctr *subscribtionController) HistoryBalanceSaldoByDate(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	// roleID, _ := strconv.Atoi(claims["role_id"].(string))

	month := c.Request.URL.Query().Get("month")
	year := c.Request.URL.Query().Get("year")

	currentTime := time.Now()
	currentMonth := int(currentTime.Month())
	currentYear := int(currentTime.Year())
	monthStr := strconv.Itoa(currentMonth)
	yearStr := strconv.Itoa(currentYear)

	fmt.Println("Bulan sekarang =>", strconv.Itoa(currentMonth), "Filter by Bulan =", month)

	if monthStr == month && yearStr == year {
		strWhre := ""
		getDates, err := ctr.inputForm.GetDatesNew(strWhre)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		var resHeader []objects.HistoryBalanceSaldoWithDate
		var totalSubs int
		var totalSubsDeleted int
		var totalSubsUpdated int
		var totalBlastInfo int
		var totalTopupSubs int

		for i := 0; i < len(getDates); i++ {
			dateString := getDates[i].Date
			layout := "02 Jan 2006"

			date, err := time.Parse(layout, dateString)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return
			}
			dateDb := date.Format("2006-01-02")

			var whereComp objects.HistoryBalanceSaldo
			whereComp.OrganizationID = organizationID
			getFormByOrgID, err := ctr.formMod.GetFormsOrganization(whereComp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   err.Error(),
					"status":  false,
					"message": "Error: Data user is not deleted",
				})
				return
			}

			fmt.Println("len(getFormByOrgID)  -->", organizationID, len(getFormByOrgID))

			if len(getFormByOrgID) > 0 {
				var totalTopupSubmission int
				var totalSubmission int
				var totalDeleted int
				var totalUpdated int
				var totalBlast int

				for i := 0; i < len(getFormByOrgID); i++ {

					var whereInForm tables.InputForms
					var whreStr = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"
					var whreStrg = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + dateDb + "' AND if.deleted_at IS NULL"
					var whreString = "TO_CHAR(notification_histories.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"
					var whreDate = "TO_CHAR(organization_topup_histories.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"

					getDataAllSubmission, _ := ctr.inputForm.GetInputFormUnscopedRows(getFormByOrgID[i].FormID, whereInForm, whreStrg, objects.Paging{})
					getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(getFormByOrgID[i].FormID, whereInForm, whreStr)
					getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedDataWithDate(getFormByOrgID[i].FormID, whreStr)
					getDataBlastInfo, _ := ctr.formMod.GetBlastInfoData(getFormByOrgID[i].FormID, whreString)
					getHistoryTopupByDate, _ := ctr.formMod.GetHistoryTopupByDate(organizationID, whreDate)
					// fmt.Println(getHistoryTopupByDate)

					for i := 0; i < len(getHistoryTopupByDate); i++ {
						totalTopupSubmission = getHistoryTopupByDate[i].Quota
					}
					for i := 0; i < len(getDataSubmissionUpdated); i++ {
						totalUpdated += getDataSubmissionUpdated[i].UpdatedCount
					}
					totalSubmission += len(getDataAllSubmission)
					totalDeleted += len(getDataSubmissionDeleted)
					totalBlast += len(getDataBlastInfo)

					fmt.Println("param ---->", getFormByOrgID[i].FormID)
					fmt.Println("totalSubmission ---->", totalSubmission)

					totalSubs = totalSubmission
					totalSubsDeleted = totalDeleted
					totalSubsUpdated = totalUpdated
					totalBlastInfo = totalBlast
					totalTopupSubs = totalTopupSubmission
				}
			}

			var each objects.HistoryBalanceSaldoWithDate
			each.ID = i + 1
			each.Date = getDates[i].Date
			each.DateDB = dateDb
			each.TopupSubmission = totalTopupSubs
			each.TotalSubmission = totalSubs
			each.TotalDeletedSubmission = totalSubsDeleted
			each.TotalUpdatedSubmission = totalSubsUpdated
			each.TotalBlast = totalBlastInfo
			each.TotalUsageBalance = each.TotalSubmission + each.TotalUpdatedSubmission + each.TotalBlast + each.TotalDeletedSubmission
			fmt.Println(dateDb, "------------------>", totalSubs, totalSubsUpdated, totalSubsDeleted, totalBlastInfo, totalSubsDeleted)

			resHeader = append(resHeader, each)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    resHeader,
		})
		return
	} else {
		strWhre := ""
		getDates, err := ctr.inputForm.GetDatesWithFilter(month, year, strWhre)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var resHeader []objects.HistoryBalanceSaldoWithDate
		var totalSubs int
		var totalSubsDeleted int
		var totalSubsUpdated int
		var totalBlastInfo int
		var totalTopupSubs int

		for i := 0; i < len(getDates); i++ {
			dateString := getDates[i].Date
			layout := "02 Jan 2006"

			date, err := time.Parse(layout, dateString)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return
			}
			dateDb := date.Format("2006-01-02")

			var whereComp objects.HistoryBalanceSaldo
			whereComp.OrganizationID = organizationID
			getFormByOrgID, err := ctr.formMod.GetFormsOrganization(whereComp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   err.Error(),
					"status":  false,
					"message": "Error: Data user is not deleted",
				})
				return
			}

			if len(getFormByOrgID) > 0 {
				var totalTopupSubmission int
				var totalSubmission int
				var totalDeleted int
				var totalUpdated int
				var totalBlast int

				for i := 0; i < len(getFormByOrgID); i++ {

					var whereInForm tables.InputForms
					var whreStr = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"
					var whreStrg = "TO_CHAR(if.created_at, 'YYYY-MM-DD') = '" + dateDb + "' AND if.deleted_at IS NULL"
					var whreString = "TO_CHAR(notification_histories.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"
					var whreDate = "TO_CHAR(organization_topup_histories.created_at, 'YYYY-MM-DD') = '" + dateDb + "'"

					getDataAllSubmission, _ := ctr.inputForm.GetInputFormUnscopedRows(getFormByOrgID[i].FormID, whereInForm, whreStrg, objects.Paging{})
					getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(getFormByOrgID[i].FormID, whereInForm, whreStr)
					getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedDataWithDate(getFormByOrgID[i].FormID, whreStr)
					getDataBlastInfo, _ := ctr.formMod.GetBlastInfoData(getFormByOrgID[i].FormID, whreString)
					getHistoryTopupByDate, _ := ctr.formMod.GetHistoryTopupByDate(organizationID, whreDate)

					for i := 0; i < len(getHistoryTopupByDate); i++ {
						totalTopupSubmission = getHistoryTopupByDate[i].Quota
					}
					for i := 0; i < len(getDataSubmissionUpdated); i++ {
						totalUpdated += getDataSubmissionUpdated[i].UpdatedCount
					}
					totalSubmission += len(getDataAllSubmission)
					totalDeleted += len(getDataSubmissionDeleted)
					totalBlast += len(getDataBlastInfo)

					totalSubs = totalSubmission
					totalSubsDeleted = totalDeleted
					totalSubsUpdated = totalUpdated
					totalBlastInfo = totalBlast
					totalTopupSubs = totalTopupSubmission

				}
			}

			var each objects.HistoryBalanceSaldoWithDate
			each.ID = i + 1
			each.Date = getDates[i].Date
			each.DateDB = dateDb
			each.TopupSubmission = totalTopupSubs
			each.TotalSubmission = totalSubs
			each.TotalDeletedSubmission = totalSubsDeleted
			each.TotalUpdatedSubmission = totalSubsUpdated
			each.TotalBlast = totalBlastInfo
			each.TotalUsageBalance = each.TotalSubmission + each.TotalUpdatedSubmission + each.TotalBlast + each.TotalDeletedSubmission

			resHeader = append(resHeader, each)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    resHeader,
		})
		return
	}
}
