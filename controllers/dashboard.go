package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// interface
type HomeController interface {
	Home(c *gin.Context)
}

type homeController struct {
	userMod   models.UserModels
	formMod   models.FormModels
	helper    helpers.Helper
	inputForm models.InputFormModels
	compMod   models.CompaniesModels
	subMod    models.SubsModels
}

func NewHomeController(userModels models.UserModels, formModel models.FormModels, helper helpers.Helper, inputModels models.InputFormModels, compModels models.CompaniesModels, subsModel models.SubsModels) HomeController {
	return &homeController{
		userMod:   userModels,
		formMod:   formModel,
		helper:    helper,
		inputForm: inputModels,
		compMod:   compModels,
		subMod:    subsModel,
	}
}

func (ctr *homeController) Home(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println(userID, organizationID)
	}
	var getComp tables.Organizations
	projectID := c.Param("groupid")
	// iProjectID, _ := strconv.Atoi(projectID)
	//get my form
	var getForms []tables.FormAll
	var getTotalSubmission []tables.FormAll
	var getTotalSubmissionFormEksternal []tables.FormAll
	// var getFormsEksternal []tables.FormAll
	if roleID == 1 { // 1 is owner

		// SUPER ADMIN here
		var fields tables.FormOrganizationsJoin
		fields.OrganizationID = organizationID

		var buffer bytes.Buffer
		whereString := ""
		buffer.WriteString(" forms.form_status_id in (1)")
		if projectID != "" {
			buffer.WriteString(" AND forms.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + projectID + ")")
		}

		whereString = buffer.String()

		results, err := ctr.formMod.GetFormOwnerRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getForms = results

		//EKSTERNAL
		results, err = ctr.formMod.GetFormEksternalOwnerRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getTotalSubmissionFormEksternal = results

	} else {
		var buffer bytes.Buffer
		var fields tables.Forms
		fields.FormStatusID = 1

		whereString := ""
		buffer.WriteString(" forms.form_status_id in (1) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ") AND fo.organization_id = " + strconv.Itoa(organizationID) + "")

		if projectID != "" {
			buffer.WriteString(" AND forms.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + projectID + ")")
		}

		whereString = buffer.String()

		results, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		getForms = results

	}
	if roleID == 1 {
		// SUPER ADMIN here
		var whre tables.FormOrganizationsJoin
		whre.OrganizationID = organizationID

		var buffer bytes.Buffer
		whereString := ""
		buffer.WriteString(" forms.form_status_id in (1,2,3)")
		if projectID != "" {
			buffer.WriteString(" AND forms.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + projectID + ")")
		}

		whereString = buffer.String()

		results, err := ctr.formMod.GetFormOwnerRows(whre, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getTotalSubmission = results

	} else {
		fmt.Println("dash/home", roleID)

		var buffer bytes.Buffer
		var fields tables.Forms

		whereString := ""
		buffer.WriteString(" forms.form_status_id in (1,2,3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ") AND fo.organization_id = " + strconv.Itoa(organizationID) + "")

		if projectID != "" {
			buffer.WriteString(" AND forms.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + projectID + ")")
		}

		whereString = buffer.String()

		results, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getTotalSubmission = results

		//EKSTERNAL
		var bufferr bytes.Buffer
		var fieldds tables.Forms
		whereStr := ""

		bufferr.WriteString(" f.form_status_id not in (3) AND f.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereStr = bufferr.String()

		var paging objects.Paging
		getFormsEksternal, err := ctr.formMod.GetListFormEksternal(fieldds, whereStr, paging, userID, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		getTotalSubmissionFormEksternal = getFormsEksternal
	}
	// -----------------
	// var wherefld tables.Forms
	// wherefld.CreatedBy = userID
	// wherefld.FormStatusID = 1

	// var buffer bytes.Buffer
	// // buffer.WriteString("forms.form_status_id not in (3) ")
	// if projectID != "" {
	// 	buffer.WriteString("forms.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + projectID + ")")
	// }

	// whreStr := buffer.String()
	// getForms, err := ctr.formMod.GetFormWhreRows(wherefld, whreStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }

	fmt.Println("len(getForms) :::", len(getForms))
	// fmt.Println(getForms[1])
	// os.Exit(0)

	totalSubmission := 0
	totalSubmissionEksternal := 0
	totalActiveRespondens := 0
	totalActiveResponden := 0

	var res objects.Home
	if len(getTotalSubmissionFormEksternal) > 0 {

		for i := 0; i < len(getTotalSubmissionFormEksternal); i++ {
			var whereInForm tables.InputForms

			whreStr := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""
			getDataRespons, err := ctr.inputForm.GetInputFormRows(getTotalSubmissionFormEksternal[i].ID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			totalSubmissionEksternal += len(getDataRespons)

			whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND fuo.organization_id = " + strconv.Itoa(organizationID) + ""
			var whreActive tables.InputForms
			getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(getTotalSubmissionFormEksternal[i].ID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			totalActiveRespondens += len(getActiveRespondens)

		}
	}
	if len(getTotalSubmission) > 0 {

		for i := 0; i < len(getTotalSubmission); i++ {

			//count respon all form
			var whereInForm tables.InputForms
			// whreStr := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

			getDataRespons, err := ctr.inputForm.GetInputFormRows(getTotalSubmission[i].ID, whereInForm, "", objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			totalSubmission += len(getDataRespons)
			// fmt.Println(getDataRespons)

		}
	}

	if len(getForms) > 0 {

		// totalRespon := 0
		// totalResponden := 0

		for i := 0; i < len(getForms); i++ {

			//count respon all form
			// var whereInForm tables.InputForms
			// getDataRespons, err := ctr.inputForm.GetInputFormRows(getForms[i].ID, whereInForm, "", objects.Paging{})
			// if err != nil {
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"error": err,
			// 	})
			// 	return
			// }
			// totalRespon += len(getDataRespons)

			// total responden all
			// var whereFU tables.FormUsers
			// whereFU.FormID = getForms[i].ID
			// fuRows, err := ctr.formMod.GetFormUserRows(whereFU)
			// if err != nil {
			// 	fmt.Println("err: GetFormUserRows", err)
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"error": err,
			// 	})
			// 	return
			// }
			// fmt.Println("getForms[i].ID---", getForms[i].ID)
			// totalResponden += len(fuRows)

			//total active responden
			whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			var whreActive tables.InputForms
			getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(getForms[i].ID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			totalActiveResponden += len(getActiveRespondens)

		}

		//total responden uniq all
		// fuRows, err := ctr.formMod.GetFormUserUniqRows(userID, iProjectID)
		// if err != nil {
		// 	fmt.Println("err: GetFormUserRows", err)
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		// res.TotalRespon = totalRespon
	}

	//get data user account
	var fields tables.Users
	fields.ID = userID
	getUser, err := ctr.userMod.GetUserRow(fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	getComp, _ = ctr.compMod.GetCompaniesRow(tables.Organizations{ID: organizationID})
	isUserCompanyActive := false

	// if roleID == 1 { // owner
	var whrSubs objects.SubsPlan
	whrSubs.OrganizationID = getComp.ID
	checkSubsPlan, _ := ctr.subMod.GetPlanRow(whrSubs)
	if checkSubsPlan.IsBlocked == false {
		isUserCompanyActive = true
	}
	// } else {
	// 	isUserCompanyActive = false
	// }
	fmt.Println("checkSubsPlan.IsBlocked ::", checkSubsPlan.IsBlocked)
	//check profile data
	checkProfile := false
	if getUser.Phone != "" {
		checkProfile = true
	}

	res.UserID = userID
	res.UserName = getUser.Name
	res.UserAvatar = getUser.Avatar
	res.CompanyName = getComp.Name
	res.IsUserCompanyActive = isUserCompanyActive
	res.IsProfileComplete = checkProfile
	res.TotalRespon = totalSubmission + totalSubmissionEksternal
	res.TotalForm = len(getForms) + len(getTotalSubmissionFormEksternal)
	res.TotalResponden = totalActiveResponden + totalActiveRespondens

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    res,
	})
	return

}

func (ctr *homeController) Home_last(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	//get data user account
	var fields tables.Users
	fields.ID = userID
	getUser, err := ctr.userMod.GetUserRow(fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var res objects.Home
	res.UserID = userID
	res.UserName = getUser.Name
	res.UserAvatar = getUser.Avatar
	res.CompanyName = helpers.Substr(getUser.CompanyName, 0, 10)
	res.TotalForm = 0
	res.TotalResponden = 0
	res.TotalRespon = 0

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    res,
	})
	return

}
