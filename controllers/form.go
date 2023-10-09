package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jackc/pgconn"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// interface
type FormController interface {
	FormGroup1List(c *gin.Context)
	FormGroup2List(c *gin.Context)
	FormGroup3List(c *gin.Context)
	FormGroup4List(c *gin.Context)
	FormList(c *gin.Context)
	FormListArchive(c *gin.Context)
	FormListArchiveLast(c *gin.Context)
	FormCreate(c *gin.Context)
	FormCreateLocation(c *gin.Context)
	FormUpdate(c *gin.Context)
	FormUpdateLocation(c *gin.Context)
	FormDestroy(c *gin.Context)
	FormDestroyLocation(c *gin.Context)
	FormDuplicate(c *gin.Context)
	FormDuplicate2(c *gin.Context)
	FormDuplicate3(c *gin.Context)
	FormUserList(c *gin.Context)
	FormUserCreate(c *gin.Context)
	FormUserUpdate(c *gin.Context)
	FormUserStatusUpdate(c *gin.Context)
	FormUserConnect(c *gin.Context)
	FormUserDisconnect(c *gin.Context)
	FieldList(c *gin.Context)
	FieldGroupList(c *gin.Context)
	Field2GroupList(c *gin.Context)
	FieldGroup3List(c *gin.Context)
	FieldCreate(c *gin.Context)
	FieldUpdate(c *gin.Context)
	FieldDestroy(c *gin.Context)
	FieldTypeList(c *gin.Context)
	FieldGroupCreate(c *gin.Context)
	FieldGroupUpdate(c *gin.Context)
	FieldSectionCreate(c *gin.Context)
	FieldSectionUpdate(c *gin.Context)
	ConditionRuleList(c *gin.Context)
	FieldConditionSave(c *gin.Context)
	FieldSaveImage(c *gin.Context)
	FormShare(c *gin.Context)
	FormGetShare(c *gin.Context)
	FormUpdateStatus(c *gin.Context)
	FieldSortOrderSave(c *gin.Context)
	FormUserAdminFrm(c *gin.Context)
	FormUserAdminFrmConnect(c *gin.Context)
	FormUserAdminFrmDisconnect(c *gin.Context)
	FormUserAdminCheck(c *gin.Context)
	FormAttendanceRequired(c *gin.Context)
	UserGetFormList(c *gin.Context)
	UpdateAdminPermission(c *gin.Context)
	GetFillingType(c *gin.Context)

	// admin
	AdminFormList(c *gin.Context)
	AdminFormListNew(c *gin.Context)
	AdminListPermission(c *gin.Context)
	ListAdminEks(c *gin.Context)
	DeleteAdminEks(c *gin.Context)
	AddAdminPermisMan(c *gin.Context)
	AddAdminPermisOto(c *gin.Context)

	// company form
	FormCompanyConnect(c *gin.Context)
	FormCompanyDisconnect(c *gin.Context)
	FormToCompanyList(c *gin.Context)
	FormCompanyNotInList(c *gin.Context)
	FormCompanyUpdateQuota(c *gin.Context)
	FormMultyAccessList(c *gin.Context)
	FormCompanySharingList(c *gin.Context)
	FormCompanySharingDelete(c *gin.Context)
	FormExternalSharingList(c *gin.Context)
	FormExternalSharingTotal(c *gin.Context)
	FilterFormToCompanyList(c *gin.Context)
	CheckForm(c *gin.Context)

	FormDataExport(c *gin.Context) //export

	//invite
	// InviteAdmin(c *gin.Context)

	//List Form Template
	ListFormTemplate(c *gin.Context)
	ListProject(c *gin.Context)
	FormTemplate(c *gin.Context)
}

type formController struct {
	formMod       models.FormModels
	formFieldMod  models.FormFieldModels
	ftMod         models.FieldTypeModels
	ruleMod       models.RuleModels
	helper        helpers.Helper
	userMod       models.UserModels
	compMod       models.CompaniesModels
	inputForm     models.InputFormModels
	pgErr         *pgconn.PgError
	permissMod    models.PermissionModels
	projectMod    models.ProjectModels
	attendanceMod models.AttendanceModels
}

func NewFormController(formModel models.FormModels, formFieldModel models.FormFieldModels, ftModel models.FieldTypeModels, ruleMod models.RuleModels, helper helpers.Helper, userModel models.UserModels, compModel models.CompaniesModels, inputFormModel models.InputFormModels, permissModel models.PermissionModels, projModel models.ProjectModels, attendanceModel models.AttendanceModels) FormController {
	return &formController{
		formMod:       formModel,
		ftMod:         ftModel,
		formFieldMod:  formFieldModel,
		ruleMod:       ruleMod,
		helper:        helper,
		userMod:       userModel,
		compMod:       compModel,
		inputForm:     inputFormModel,
		permissMod:    permissModel,
		projectMod:    projModel,
		attendanceMod: attendanceModel,
	}
}

func (ctr *formController) FormGroup1List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			fmt.Println("err : GetFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			fmt.Println("err : GetFormPeriodeRangeRow")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			fmt.Println("err : GetInputFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		var resultgetProject []tables.FormAll
		if roleID == 1 { // 1 is owner

			var searchForm bytes.Buffer
			var whereName string
			// SUPER ADMIN here
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			whereGroupStr := ``
			whereString := "AND forms.form_status_id not in (3)"
			if searchKeyWord != "" {
				searchForm.WriteString(" where name ilike '%" + searchKeyWord + "%'")
				whereName = searchForm.String()
			}

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeSuperAdminNew1(fields, whereName, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormMergeSuperAdminNew1(fields, whereName, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			fmt.Println("result :::::", len(result), len(getFormsAll))

			whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormMergeSuperAdminNew1(fields, whereName, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

			getProject, err := ctr.formMod.GetProjectSuperAdminNew(fields, whereName, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetProject = getProject

		} else {
			var buffer bytes.Buffer
			var whreGroup bytes.Buffer
			var whre bytes.Buffer
			var searchForm bytes.Buffer
			var whereName string
			var whereString string
			var whereGroupStr string
			var whereStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				searchForm.WriteString(" where name ilike '%" + searchKeyWord + "%'")
				whereName = searchForm.String()
			}

			whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
			whereGroupStr = whreGroup.String()

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			whre.WriteString(" f.form_status_id not in (3) AND f.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereStr = whre.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeAdminNew1(fields, whereName, whereString, whereGroupStr, whereStr, userID, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms

			//get all data
			getFormsAll, err := ctr.formMod.GetFormMergeAdminNew1(fields, whereName, whereString, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormMergeAdminNew1(fields, whereName, whereAll, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

			getProject, err := ctr.formMod.GetProjectAdminNew(fields, whereName, whereAll, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetProject = getProject

		}

		if len(result) > 0 {

			var res []objects.MergeForms
			for i := 0; i < len(result); i++ {

				if result[i].ProjectID > 0 {

					//get total form in project/groups
					var whreFrm tables.ProjectForms
					var whrStr string
					whreFrm.ProjectID = result[i].ProjectID

					if roleID != 1 {
						whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
					}
					getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

					var each objects.MergeForms
					each.ProjectID = result[i].ProjectID
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.TotalForms = len(getFormsInGroup)
					each.FormShared = result[i].FormShared
					res = append(res, each)

				} else {
					fmt.Println("-------------------------->")
					// get total responden form internal
					var whereFU tables.JoinFormUsers
					whereFU.FormID = result[i].ID
					whereFU.Type = "respondent"
					whreStr := ""

					getRespondenIntenal, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// get total responden form external
					var where tables.JoinFormUsers
					where.FormID = result[i].ID
					where.Type = "respondent"
					whreString := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

					getRespondenExternal, err := ctr.formMod.GetFormUserRows(where, whreString)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total active responden form internal
					whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					var whreActive tables.InputForms
					getActiveRespondensInternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total active responden form external
					whreStrDate := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND ifo.organization_id = " + strconv.Itoa(organizationID) + ""
					// var whreActive tables.InputForms
					getActiveRespondensExternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrDate)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon Internal
					var whereInForm tables.InputForms
					getDataResponFormInternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon External
					whreStr = "ifo.organization_id = " + strconv.Itoa(organizationID) + ""
					getDataResponFormExternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStr, objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total performance Internal
					var totalPerformInternal int
					var totalPerformFloatInternal float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloatInternal = float64(len(getDataResponFormInternal)) / float64(result[i].SubmissionTargetUser)
						totalPerformInternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatInternal, 'f', 0, 64))
					}

					//total performance External
					var totalPerformExternal int
					var totalPerformFloatExternal float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloatExternal = float64(len(getDataResponFormExternal)) / float64(result[i].SubmissionTargetUser)
						totalPerformExternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatExternal, 'f', 0, 64))
					}

					// get permission admin
					isPermission := false
					if roleID > 1 {
						var whr tables.FormUserPermissionJoin
						whr.PermissionID = 6 //(6 is edit responden)
						whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
						getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}
						isPermission = getPermission.Status
					}

					if userID == result[i].CreatedBy {
						isPermission = true
					}

					if roleID == 1 {
						isPermission = true
					}

					var each objects.MergeForms
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.ProfilePic = result[i].ProfilePic
					each.FormStatusID = result[i].FormStatusID
					each.FormStatus = result[i].FormStatus
					each.Notes = result[i].Notes
					each.PeriodStartDate = result[i].PeriodStartDate
					each.PeriodEndDate = result[i].PeriodEndDate
					each.CreatedBy = result[i].CreatedBy
					each.CreatedByName = result[i].CreatedByName
					each.CreatedByEmail = result[i].CreatedByEmail
					each.IsAttendanceRequired = result[i].IsAttendanceRequired
					each.UpdatedByName = ""
					each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
					each.SubmissionTarget = result[i].SubmissionTargetUser
					each.PeriodeRange = 0
					each.IsEditResponden = isPermission
					each.Type = result[i].Type
					if result[i].Type == "internal" {
						each.TotalRespon = len(getDataResponFormInternal)
						each.TotalPerformance = totalPerformInternal
						each.TotalPerformanceFloat = totalPerformFloatInternal
						each.TotalResponden = len(getRespondenIntenal)
						each.TotalRespondenActive = len(getActiveRespondensInternal)
					} else if result[i].Type == "external" {
						each.TotalRespon = len(getDataResponFormExternal)
						each.TotalPerformance = totalPerformExternal
						each.TotalPerformanceFloat = totalPerformFloatExternal
						each.TotalResponden = len(getRespondenExternal)
						each.TotalRespondenActive = len(getActiveRespondensExternal)
					}
					// if result[i].FormShareCount > 0 {
					// 	each.FormShared = "out"
					// } else {
					// 	each.FormShared = result[i].FormShared
					// }
					each.FormShared = result[i].FormShared
					each.FormSharedCount = "Dibagikan ke " + strconv.Itoa(result[i].FormShareCount) + " Organisasi"
					each.FormExternalCompanyName = result[i].FormExternalCompanyName
					each.FormExternalCompanyImage = result[i].FormExternalCompanyImage

					each.IsQuotaSharing = result[i].IsQuotaSharing
					each.AccessType = result[i].AccessType

					res = append(res, each)
				}
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll) - len(resultgetProject)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormGroup2List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			fmt.Println("err : GetFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			fmt.Println("err : GetFormPeriodeRangeRow")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			fmt.Println("err : GetInputFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			// SUPER ADMIN here
			var buffer bytes.Buffer
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			whereGroupStr := ``
			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" AND forms.name ilike '%" + searchKeyWord + "%'  ")
				whereGroupStr = " AND t.name ilike '%" + searchKeyWord + "%'"
			}

			buffer.WriteString(" AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			fmt.Println("result :::::", len(result), len(getFormsAll))

			whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

		} else {
			var buffer bytes.Buffer
			var whreGroup bytes.Buffer
			var whereString string
			var whereGroupStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%'   ")
				whreGroup.WriteString(" AND t.name ilike '%" + searchKeyWord + "%'")
			}

			whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
			whereGroupStr = whreGroup.String()

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms

			//get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll
		}

		if len(result) > 0 {

			var res []objects.Forms
			for i := 0; i < len(result); i++ {

				if result[i].ProjectID > 0 {

					//get total form in project/groups
					var whreFrm tables.ProjectForms
					var whrStr string
					whreFrm.ProjectID = result[i].ProjectID

					if roleID != 1 {
						whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
					}
					getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

					var each objects.Forms
					each.ProjectID = result[i].ProjectID
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.TotalForms = len(getFormsInGroup)
					res = append(res, each)

				} else {

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

					//total active responden
					whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					var whreActive tables.InputForms
					getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon
					var whereInForm tables.InputForms
					getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total performance
					var totalPerform int
					var totalPerformFloat float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
						totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
					}

					// get permission admin
					isPermission := false
					if roleID > 1 {
						var whr tables.FormUserPermissionJoin
						whr.PermissionID = 6 //(6 is edit responden)
						whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
						getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}
						isPermission = getPermission.Status
					}

					if userID == result[i].CreatedBy {
						isPermission = true
					}

					if roleID == 1 {
						isPermission = true
					}

					var each objects.Forms
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.ProfilePic = result[i].ProfilePic
					each.FormStatusID = result[i].FormStatusID
					each.FormStatus = result[i].FormStatus
					each.Notes = result[i].Notes
					each.PeriodStartDate = result[i].PeriodStartDate
					each.PeriodEndDate = result[i].PeriodEndDate
					each.CreatedBy = result[i].CreatedBy
					each.CreatedByName = result[i].CreatedByName
					each.CreatedByEmail = result[i].CreatedByEmail
					each.TotalResponden = len(getResponden)
					each.TotalRespondenActive = len(getActiveRespondens)
					each.TotalRespon = len(getDataRespons)
					each.TotalPerformance = totalPerform
					each.TotalPerformanceFloat = totalPerformFloat
					each.IsAttendanceRequired = result[i].IsAttendanceRequired
					each.UpdatedByName = ""
					each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
					each.SubmissionTarget = result[i].SubmissionTargetUser
					each.PeriodeRange = 0
					each.IsEditResponden = isPermission

					res = append(res, each)
				}
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormGroup3List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			fmt.Println("err : GetFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			fmt.Println("err : GetFormPeriodeRangeRow")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			fmt.Println("err : GetInputFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		var resultgetProject []tables.FormAll
		if roleID == 1 { // 1 is owner

			var searchForm bytes.Buffer
			var whereName string
			// SUPER ADMIN here
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			whereGroupStr := ``
			whereString := "AND forms.form_status_id not in (3)"
			if searchKeyWord != "" {
				searchForm.WriteString(" where name ilike '%" + searchKeyWord + "%'")
				whereName = searchForm.String()
			}

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeSuperAdminNew(fields, whereName, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormMergeSuperAdminNew(fields, whereName, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			fmt.Println("result :::::", len(result), len(getFormsAll))

			whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormMergeSuperAdminNew(fields, whereName, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

			getProject, err := ctr.formMod.GetProjectSuperAdminNew(fields, whereName, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetProject = getProject

		} else {
			var buffer bytes.Buffer
			var whreGroup bytes.Buffer
			var whre bytes.Buffer
			var searchForm bytes.Buffer
			var whereName string
			var whereString string
			var whereGroupStr string
			var whereStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				searchForm.WriteString(" where name ilike '%" + searchKeyWord + "%'")
				whereName = searchForm.String()
			}

			whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
			whereGroupStr = whreGroup.String()

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			whre.WriteString(" f.form_status_id not in (3) AND f.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereStr = whre.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeAdminNew(fields, whereName, whereString, whereGroupStr, whereStr, userID, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms

			//get all data
			getFormsAll, err := ctr.formMod.GetFormMergeAdminNew(fields, whereName, whereString, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormMergeAdminNew(fields, whereName, whereAll, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

			getProject, err := ctr.formMod.GetProjectAdminNew(fields, whereName, whereAll, whereGroupStr, whereStr, userID, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetProject = getProject

		}

		if len(result) > 0 {

			var res []objects.MergeForms
			for i := 0; i < len(result); i++ {

				if result[i].ProjectID > 0 {

					//get total form in project/groups
					var whreFrm tables.ProjectForms
					var whrStr string
					whreFrm.ProjectID = result[i].ProjectID

					if roleID != 1 {
						whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
					}
					getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

					var each objects.MergeForms
					each.ProjectID = result[i].ProjectID
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.TotalForms = len(getFormsInGroup)
					each.FormShared = result[i].FormShared
					res = append(res, each)

				} else {
					fmt.Println("-------------------------->")
					// get total responden form internal
					var whereFU tables.JoinFormUsers
					whereFU.FormID = result[i].ID
					whereFU.Type = "respondent"
					whreStr := ""

					getRespondenIntenal, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// get total responden form external
					var where tables.JoinFormUsers
					where.FormID = result[i].ID
					where.Type = "respondent"
					whreString := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

					getRespondenExternal, err := ctr.formMod.GetFormUserRows(where, whreString)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total active responden form internal
					whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					var whreActive tables.InputForms
					getActiveRespondensInternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					whreStrDate := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND ifo.organization_id = " + strconv.Itoa(organizationID) + ""
					// var whreActive tables.InputForms
					getActiveRespondensExternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrDate)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon Internal
					var whereInForm tables.InputForms
					getDataResponFormInternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon External
					whreStr = "ifo.organization_id = " + strconv.Itoa(organizationID) + ""
					getDataResponFormExternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStr, objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total performance Internal
					var totalPerformInternal int
					var totalPerformFloatInternal float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloatInternal = float64(len(getDataResponFormInternal)) / float64(result[i].SubmissionTargetUser)
						totalPerformInternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatInternal, 'f', 0, 64))
					}

					//total performance External
					var totalPerformExternal int
					var totalPerformFloatExternal float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloatExternal = float64(len(getDataResponFormExternal)) / float64(result[i].SubmissionTargetUser)
						totalPerformExternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatExternal, 'f', 0, 64))
					}

					// get permission admin
					isPermission := false
					if roleID > 1 {
						var whr tables.FormUserPermissionJoin
						whr.PermissionID = 6 //(6 is edit responden)
						whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
						getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}
						isPermission = getPermission.Status
					}

					if userID == result[i].CreatedBy {
						isPermission = true
					}

					if roleID == 1 {
						isPermission = true
					}

					var each objects.MergeForms
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.ProfilePic = result[i].ProfilePic
					each.FormStatusID = result[i].FormStatusID
					each.FormStatus = result[i].FormStatus
					each.Notes = result[i].Notes
					each.PeriodStartDate = result[i].PeriodStartDate
					each.PeriodEndDate = result[i].PeriodEndDate
					each.CreatedBy = result[i].CreatedBy
					each.CreatedByName = result[i].CreatedByName
					each.CreatedByEmail = result[i].CreatedByEmail
					each.IsAttendanceRequired = result[i].IsAttendanceRequired
					each.UpdatedByName = ""
					each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
					each.SubmissionTarget = result[i].SubmissionTargetUser
					each.PeriodeRange = 0
					each.IsEditResponden = isPermission
					each.Type = result[i].Type
					if result[i].Type == "internal" {
						each.TotalRespon = len(getDataResponFormInternal)
						each.TotalPerformance = totalPerformInternal
						each.TotalPerformanceFloat = totalPerformFloatInternal
						each.TotalResponden = len(getRespondenIntenal)
						each.TotalRespondenActive = len(getActiveRespondensInternal)
					} else if result[i].Type == "external" {
						each.TotalRespon = len(getDataResponFormExternal)
						each.TotalPerformance = totalPerformExternal
						each.TotalPerformanceFloat = totalPerformFloatExternal
						each.TotalResponden = len(getRespondenExternal)
						each.TotalRespondenActive = len(getActiveRespondensExternal)
					}
					// if result[i].FormShareCount > 0 {
					// 	each.FormShared = "out"
					// } else {
					// 	each.FormShared = result[i].FormShared
					// }
					each.FormShared = result[i].FormShared
					each.FormSharedCount = "Dibagikan ke " + strconv.Itoa(result[i].FormShareCount) + " Organisasi"
					each.FormExternalCompanyName = result[i].FormExternalCompanyName
					each.FormExternalCompanyImage = result[i].FormExternalCompanyImage

					each.SharingSaldo = "Ya"
					each.StatusAdmin = "Editor"

					res = append(res, each)
				}
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll) - len(resultgetProject)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormGroup4List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			fmt.Println("err : GetFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			fmt.Println("err : GetFormPeriodeRangeRow")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			fmt.Println("err : GetInputFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			// SUPER ADMIN here
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			whereGroupStr := ``
			whereString := ""
			if searchKeyWord != "" {
				whereGroupStr = " AND t.name ilike '%" + searchKeyWord + "%'"
			}

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			fmt.Println("result :::::", len(result), len(getFormsAll))

			// whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(+ "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, "", whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

		} else {
			var buffer bytes.Buffer
			var whreGroup bytes.Buffer
			var whereString string
			var whereGroupStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%'   ")
				whreGroup.WriteString(" AND t.name ilike '%" + searchKeyWord + "%'")
			}

			whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
			whereGroupStr = whreGroup.String()

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms

			//get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll
		}

		if len(result) > 0 {

			var res []objects.MergeForms
			for i := 0; i < len(result); i++ {

				if result[i].ProjectID > 0 {

					//get total form in project/groups
					var whreFrm tables.ProjectForms
					var whrStr string
					whreFrm.ProjectID = result[i].ProjectID

					if roleID != 1 {
						whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
					}
					getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

					var each objects.MergeForms
					each.ProjectID = result[i].ProjectID
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.TotalForms = len(getFormsInGroup)
					res = append(res, each)

				} else {
					fmt.Println("-------------------------->")
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

					//total active responden
					whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					var whreActive tables.InputForms
					getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon
					var whereInForm tables.InputForms
					getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total performance
					var totalPerform int
					var totalPerformFloat float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
						totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
					}

					// get permission admin
					isPermission := false
					if roleID > 1 {
						var whr tables.FormUserPermissionJoin
						whr.PermissionID = 6 //(6 is edit responden)
						whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
						getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}
						isPermission = getPermission.Status
					}

					if userID == result[i].CreatedBy {
						isPermission = true
					}

					if roleID == 1 {
						isPermission = true
					}

					var each objects.MergeForms
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.ProfilePic = result[i].ProfilePic
					each.FormStatusID = result[i].FormStatusID
					each.FormStatus = result[i].FormStatus
					each.Notes = result[i].Notes
					each.PeriodStartDate = result[i].PeriodStartDate
					each.PeriodEndDate = result[i].PeriodEndDate
					each.CreatedBy = result[i].CreatedBy
					each.CreatedByName = result[i].CreatedByName
					each.CreatedByEmail = result[i].CreatedByEmail
					each.TotalResponden = len(getResponden)
					each.TotalRespondenActive = len(getActiveRespondens)
					each.TotalRespon = len(getDataRespons)
					each.TotalPerformance = totalPerform
					each.TotalPerformanceFloat = totalPerformFloat
					each.IsAttendanceRequired = result[i].IsAttendanceRequired
					each.UpdatedByName = ""
					each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
					each.SubmissionTarget = result[i].SubmissionTargetUser
					each.PeriodeRange = 0
					each.IsEditResponden = isPermission

					each.FormShared = result[i].FormShared
					if result[i].FormShareCount >= 1 {
						each.FormSharedCount = "Dibagikan ke " + strconv.Itoa(result[i].FormShareCount) + " Akun"
					}
					each.FormExternalCompanyName = result[i].FormExternalCompanyName
					each.FormExternalCompanyImage = result[i].FormExternalCompanyImage

					each.SharingSaldo = "Ya"
					each.StatusAdmin = "Editor"

					res = append(res, each)
				}
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormGroup5List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		result, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			fmt.Println("err : GetFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			fmt.Println("err : GetFormPeriodeRangeRow")

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			fmt.Println("err : GetInputFormRow")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			// SUPER ADMIN here
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			whereGroupStr := ``
			whereString := ""
			if searchKeyWord != "" {
				whereGroupStr = " AND t.name ilike '%" + searchKeyWord + "%'"
			}

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			fmt.Println("result :::::", len(result), len(getFormsAll))

			// whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(+ "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, "", whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

		} else {
			var buffer bytes.Buffer
			var whreGroup bytes.Buffer
			var whereString string
			var whereGroupStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%'   ")
				whreGroup.WriteString(" AND t.name ilike '%" + searchKeyWord + "%'")
			}

			whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
			whereGroupStr = whreGroup.String()

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms

			//get all data
			getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
			getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll
		}

		if len(result) > 0 {

			var res []objects.MergeForms
			for i := 0; i < len(result); i++ {

				if result[i].ProjectID > 0 {

					//get total form in project/groups
					var whreFrm tables.ProjectForms
					var whrStr string
					whreFrm.ProjectID = result[i].ProjectID

					if roleID != 1 {
						whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
					}
					getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

					var each objects.MergeForms
					each.ProjectID = result[i].ProjectID
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.TotalForms = len(getFormsInGroup)
					res = append(res, each)

				} else {
					fmt.Println("-------------------------->")
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

					//total active responden
					whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
					var whreActive tables.InputForms
					getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
					if err != nil {
						fmt.Println("err: GetFormUserRows", err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					// total respon
					var whereInForm tables.InputForms
					getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					//total performance
					var totalPerform int
					var totalPerformFloat float64
					if result[i].SubmissionTargetUser > 0 {
						totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
						totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
					}

					// get permission admin
					isPermission := false
					if roleID > 1 {
						var whr tables.FormUserPermissionJoin
						whr.PermissionID = 6 //(6 is edit responden)
						whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
						getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
							return
						}
						isPermission = getPermission.Status
					}

					if userID == result[i].CreatedBy {
						isPermission = true
					}

					if roleID == 1 {
						isPermission = true
					}

					var each objects.MergeForms
					each.ID = result[i].ID
					each.Name = result[i].Name
					each.Description = result[i].Description
					each.ProfilePic = result[i].ProfilePic
					each.FormStatusID = result[i].FormStatusID
					each.FormStatus = result[i].FormStatus
					each.Notes = result[i].Notes
					each.PeriodStartDate = result[i].PeriodStartDate
					each.PeriodEndDate = result[i].PeriodEndDate
					each.CreatedBy = result[i].CreatedBy
					each.CreatedByName = result[i].CreatedByName
					each.CreatedByEmail = result[i].CreatedByEmail
					each.TotalResponden = len(getResponden)
					each.TotalRespondenActive = len(getActiveRespondens)
					each.TotalRespon = len(getDataRespons)
					each.TotalPerformance = totalPerform
					each.TotalPerformanceFloat = totalPerformFloat
					each.IsAttendanceRequired = result[i].IsAttendanceRequired
					each.UpdatedByName = ""
					each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
					each.SubmissionTarget = result[i].SubmissionTargetUser
					each.PeriodeRange = 0
					each.IsEditResponden = isPermission

					each.FormShared = result[i].FormShared
					if result[i].FormShareCount >= 1 {
						each.FormSharedCount = "Dibagikan ke " + strconv.Itoa(result[i].FormShareCount) + " Akun"
					}
					each.FormExternalCompanyName = result[i].FormExternalCompanyName
					each.FormExternalCompanyImage = result[i].FormExternalCompanyImage

					each.SharingSaldo = "Ya"
					each.StatusAdmin = "Editor"

					res = append(res, each)
				}
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID >= 1 {
		fmt.Println("roleid = 1---------->>", formID)

		// var fields tables.Forms
		// fields.ID = formID
		// result, err := ctr.formMod.GetFormRow(fields)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		// CEK USER FORM PRIVILAGE --------------------------------

		var result tables.FormOut
		var err error

		/*
			if roleID == 1 {

				fmt.Println("roleid = 1---------->>")

				var fields tables.UserFormOrganizations
				fields.FormID = formID
				whStr := "o.created_by=" + claims["id"].(string)
				checkPrivilage, err := ctr.formMod.GetUserFormOrganization(fields, whStr)
				if err != nil {
					fmt.Println("GetUserFormOrganization---------->>")
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if checkPrivilage.ID <= 0 {
					fmt.Println("checkPrivilage---------->>")
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"error":   err,
						"message": "Form tidak tersedia atau Anda tidak mempunyai privilage ke Form",
					})
					return
				}

				var fieldForm tables.Forms
				fieldForm.ID = formID
				result, err = ctr.formMod.GetFormRow(fieldForm)
				if err != nil {
					fmt.Println("GetFormRow---------->>")
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

			} else {
				fmt.Println("roleid = else---------->>", formID)
				// admin creator check
				var whrFrm tables.Forms
				whrFrm.CreatedBy = userID
				whrFrm.ID = formID
				checkAuthor, _ := ctr.formMod.GetFormRow(whrFrm)

				var whrField tables.FormUsers
				whrField.ID = formID
				whrStr := "fu.user_id=" + claims["id"].(string)
				checkFormUser, err := ctr.formMod.GetDetailFormUserRow(whrField, whrStr)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{
						"status": false,
						"error":  err.Error(),
					})
					return
				}

				if checkFormUser.ID <= 0 && checkAuthor.ID == 0 {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Form tidak tersedia atau Anda tidak mempunyai privilage ke Form",
					})
					return
				}

				var fieldForm tables.Forms
				fieldForm.ID = formID
				result, err = ctr.formMod.GetFormRow(fieldForm)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status": false,
						"error":  err.Error(),
					})
					return
				}

			}*/

		result, err = ctr.formMod.GetFormRow(tables.Forms{ID: formID})
		if err != nil {
			fmt.Println("GetFormRow---------->>")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// END CHECK PRIVILAGE --------------------

		if result.ID > 0 {

			getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(tables.Forms{ID: formID})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// get last submit
			var whrData tables.InputForms
			getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var isContainTagloc = false
			var whrFields tables.FormFields
			whrFields.FormID = result.ID
			whrFields.FieldTypeID = 5
			checkContainTagloc, _ := ctr.formFieldMod.GetFormFieldNotParentRows(whrFields, "")
			if len(checkContainTagloc) > 0 {
				isContainTagloc = true
			}

			var res objects.FormDetail
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.FormStatus = result.FormStatus
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm
			res.IsContainTagloc = isContainTagloc
			res.AttendanceOverdateAt = result.AttendanceOverdateAt
			res.OrganizationID = result.OrganizationID
			res.OrganizationName = result.OrganizationName
			res.IsAttendanceRadius = result.IsAttendanceRadius

			// filter organization is hide / unhide
			checkFormShare, _ := ctr.formMod.GetFormCompanyInviteRow(tables.JoinFormCompanies{FormID: formID}, "")
			res.IsShowFilterOrganization = false
			if organizationID == result.OrganizationID && checkFormShare.ID >= 1 {
				res.IsShowFilterOrganization = true
			}

			if organizationID == result.OrganizationID {
				res.IsShowTabOrganization = true
			}

			if result.AttendanceOverdateAt.IsZero() == false {
				res.AttendanceOverdate = true
			}

			if res.AttendanceIn == "" && res.AttendanceOut == "" {
				fmt.Println("------- in 1")
				res.IsButton = "check_in"

				// cek chekout yg kosong sebelum hari ini (kurang dari sama dengan today) dari batas tanggal overdate At
				checkCheckoutBefore, err := ctr.attendanceMod.GetLastAttendanceOverdate(res.ID, userID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"error":   err,
						"message": "Data is available",
					})
					return
				}

				if res.AttendanceOverdate == true && checkCheckoutBefore.ID > 0 && checkCheckoutBefore.AttendanceOut == "" {
					fmt.Println("------- sub in 1 :::", res.ID, userID, ":::", checkCheckoutBefore.ID, "-----", checkCheckoutBefore, checkCheckoutBefore.AttendanceOut)

					// jika absen sebelumnya belum checkout maka button akan terus checkout
					res.IsButton = "check_out"
				}

			} else if res.AttendanceIn != "" && res.AttendanceOut == "" {
				fmt.Println("------- in 2")
				res.IsButton = "check_out"

			} else if res.AttendanceIn != "" && res.AttendanceOut != "" {
				fmt.Println("------- in 3")
				res.IsButton = "finish"

			} else if res.IsAttendanceRequired == false {
				res.IsButton = ""
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

		// cek admin user
		// var whre0 tables.UserOrganizationRoles

		// whrestr1 := `uo.user_id =` + strconv.Itoa(userID)
		// checkAdmin, _ := ctr.compMod.GetUserCompaniyToRole(whre0, whrestr1)

		// var whre2 tables.Organizations
		// whre2.CreatedBy = userID
		// checkOwner, _ := ctr.compMod.GetCompaniesRow(whre2)

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			var whrComp tables.Organizations
			whrComp.CreatedBy = userID
			whrComp.IsDefault = true
			getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)

			// SUPER ADMIN here
			var buffer bytes.Buffer
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = getComp.ID

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			// getForms, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
			// if err != nil {
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"error": err,
			// 	})
			// 	return
			// }

			getForms, err := ctr.formMod.GetFormOwnerRows(fields, whereString, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormOwnerRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := " forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

			getAll, err := ctr.formMod.GetFormOwnerRows(fields, whereAll, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

		} else {
			var buffer bytes.Buffer
			var fields tables.Forms

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
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
			//get all data
			getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

			getAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereAll, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll
		}

		if len(result) > 0 {

			var res []objects.Forms
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

				//total active responden
				whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total performance
				var totalPerform int
				var totalPerformFloat float64
				if result[i].SubmissionTargetUser > 0 {
					totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
					totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
				}

				// get permission admin
				isPermission := false
				if roleID > 1 {
					var whr tables.FormUserPermissionJoin
					whr.PermissionID = 6 //(6 is edit responden)
					whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
					getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
					isPermission = getPermission.Status
				}

				if userID == result[i].CreatedBy {
					isPermission = true
				}

				if roleID == 1 {
					isPermission = true
				}

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedBy = result[i].CreatedBy
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespondenActive = len(getActiveRespondens)
				each.TotalRespon = len(getDataRespons)
				each.TotalPerformance = totalPerform
				each.TotalPerformanceFloat = totalPerformFloat
				each.IsAttendanceRequired = result[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = result[i].SubmissionTargetUser
				each.PeriodeRange = 0
				each.IsEditResponden = isPermission
				each.IsAttendanceRadius = result[i].IsAttendanceRadius

				res = append(res, each)
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormListArchiveLast(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	formID, _ := strconv.Atoi(c.Param("id"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

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

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		// get last submit
		var whrData tables.InputForms
		getSubmission, err := ctr.inputForm.GetInputFormRow(formID, whrData, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			lastUpdate := result.UpdatedAt.Format("2006-01-02 15:04")

			lastSubm := ""
			if getSubmission.ID > 0 {
				lastSubm = getSubmission.CreatedAt.Format("2006-01-02 15:04")
			}

			var res objects.Forms
			res.ID = result.ID
			res.Name = result.Name
			res.Description = result.Description
			res.FormStatusID = result.FormStatusID
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate
			res.IsAttendanceRequired = result.IsAttendanceRequired
			res.SubmissionTarget = result.SubmissionTargetUser
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = result.FormStatusID
			res.ShareUrl = result.ShareUrl
			res.LastSubmission = lastSubm

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		var resultgetAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			var whrComp tables.Organizations
			whrComp.CreatedBy = userID
			whrComp.IsDefault = true
			getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)

			// SUPER ADMIN here
			var buffer bytes.Buffer
			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = getComp.ID

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id in (3) AND forms.deleted_at is null")
			whereString = buffer.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			// getForms, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
			// if err != nil {
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"error": err,
			// 	})
			// 	return
			// }

			getForms, err := ctr.formMod.GetFormOwnerRows(fields, whereString, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormOwnerRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := " forms.form_status_id in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + ")) AND forms.deleted_at is null"

			getAll, err := ctr.formMod.GetFormOwnerRows(fields, whereAll, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll

		} else {
			var buffer bytes.Buffer
			var fields tables.Forms

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id in (3) AND forms.deleted_at is null AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
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
			//get all data
			getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

			whereAll := "forms.form_status_id in (3) AND forms.deleted_at is null AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

			getAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereAll, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultgetAll = getAll
		}

		if len(result) > 0 {

			var res []objects.Forms
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

				//total active responden
				whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total performance
				var totalPerform int
				var totalPerformFloat float64
				if result[i].SubmissionTargetUser > 0 {
					totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
					totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
				}

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespondenActive = len(getActiveRespondens)
				each.TotalRespon = len(getDataRespons)
				each.TotalPerformance = totalPerform
				each.IsAttendanceRequired = result[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = result[i].SubmissionTargetUser
				each.PeriodeRange = 0

				res = append(res, each)
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

			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			var detail objects.DataRowsDetail
			detail.AllRows = len(resultgetAll)

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
				"detail":  detail,
			})
			return
		}
	}
}

func (ctr *formController) FormListArchive(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if userID > 0 {

		var fields tables.Forms
		fields.CreatedBy = userID
		fields.FormStatusID = 3

		whereString := ""
		if searchKeyWord != "" {
			whereString = "forms.name ilike '%" + searchKeyWord + "%'"
		}

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		result, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var res []objects.Forms
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
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespon = len(getDataRespons)
				each.ArchivedAt = result[i].ArchivedAt

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

func (ctr *formController) FormCreate(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	// roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	var reqData objects.Forms
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

	var postData tables.Forms
	postData.Name = reqData.Name
	postData.Description = reqData.Description
	postData.Notes = reqData.Notes
	postData.PeriodStartDate = reqData.PeriodStartDate
	postData.ProfilePic = reqData.ProfilePic
	postData.EncryptCode = helpers.EncodeToString(6)
	postData.CreatedBy = userID
	postData.SubmissionTargetUser = reqData.SubmissionTarget
	postData.IsAttendanceRequired = reqData.IsAttendanceRequired
	postData.ShareUrl = reqData.ShareUrl

	if reqData.PeriodEndDate != "" {
		postData.PeriodEndDate = reqData.PeriodEndDate
	}

	if reqData.IsAttendanceRequired == true && reqData.AttendanceOverdate == true {
		postData.AttendanceOverdateAt = helpers.DateNow()
	}

	if reqData.IsAttendanceRequired == true {
		postData.IsAttendanceRadius = reqData.IsAttendanceRadius
	}

	res, err := ctr.formMod.InsertForm(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctr.helper.AddLogBook(userID, 16, res.ID)

	//insert form organization
	var postData2 tables.FormOrganizations
	postData2.FormID = res.ID
	postData2.OrganizationID = organizationID
	fOrg, err := ctr.formMod.InsertFormOrganization(postData2)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if fOrg.ID > 0 {
		var obj objects.FormResponse
		obj.ID = res.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created form",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed created form, please try again",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormCreateLocation(c *gin.Context) {

	var reqData objects.ObjectFormAttendanceLocations
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

	var postData objects.ObjectFormAttendanceLocations
	postData.FormID = reqData.FormID
	postData.Name = reqData.Name
	postData.Location = reqData.Location
	postData.Longitude = reqData.Longitude
	postData.Latitude = reqData.Latitude
	postData.IsCheckIn = reqData.IsCheckIn
	postData.IsCheckOut = reqData.IsCheckOut
	postData.Radius = reqData.Radius

	Ifal, err := ctr.formMod.InsertFormAttendanceLocation(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created form location",
		"data":    Ifal,
	})
	return
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (ctr *formController) FormUpdate(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
	fmt.Println("helpers.DateNow() ::::: ", reqData.AttendanceOverdateAt, reqData.IsAttendanceRequired, helpers.DateNow())

	var sendData tables.Forms
	sendData.Name = reqData.Name
	sendData.Description = reqData.Description
	sendData.Notes = reqData.Notes
	sendData.PeriodStartDate = reqData.PeriodStartDate
	sendData.PeriodEndDate = reqData.PeriodEndDate
	sendData.UpdatedBy = userID
	sendData.SubmissionTargetUser = reqData.SubmissionTarget
	sendData.IsAttendanceRequired = reqData.IsAttendanceRequired
	sendData.ShareUrl = reqData.ShareUrl
	sendData.IsAttendanceRadius = reqData.IsAttendanceRadius

	if reqData.PeriodEndDate == "" {
		fmt.Println("masukk -----", sendData.PeriodEndDate)
	}

	sendData.ProfilePic = ""
	if reqData.ProfilePic != "" {
		sendData.ProfilePic = reqData.ProfilePic
	}

	if reqData.IsAttendanceRequired == true && reqData.AttendanceOverdate == true {
		sendData.AttendanceOverdateAt = helpers.DateNow()
	}

	res, err := ctr.formMod.UpdateForm(id, sendData)
	if err != nil {
		fmt.Println("InsertUser", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res == true {
		var obj objects.ProjectRes
		obj.ID = id

		ctr.helper.AddLogBook(userID, 8, id)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success update data",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Failed update data",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormUpdateLocation(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	var reqData objects.ObjectFormAttendanceLocations
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

	// var geo string
	var postData3 objects.ObjectFormAttendanceLocations
	postData3.ID = id
	postData3.Name = reqData.Name
	postData3.Location = reqData.Location
	// geo = `ST_SetSRID(ST_MakePoint(` + FloatToString(reqData.Longitude) + `, ` + FloatToString(reqData.Latitude) + `), 4326)`
	// postData3.Geometry = geo
	postData3.Longitude = reqData.Longitude
	postData3.Latitude = reqData.Latitude
	postData3.IsCheckIn = reqData.IsCheckIn
	postData3.IsCheckOut = reqData.IsCheckOut
	postData3.Radius = reqData.Radius

	res, err := ctr.formMod.UpdateFormAttendanceLocation(id, postData3)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res == true {
		var obj objects.ProjectRes
		obj.ID = id

		ctr.helper.AddLogBook(userID, 8, id)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success update data",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Failed update data",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormDestroy(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	_, err = ctr.formMod.DeleteForm(id, userID)
	if err != nil {
		fmt.Println("DeleteProject", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	ctr.helper.AddLogBook(userID, 3, id)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success deleted data",
		"data":    nil,
	})
	return
}

func (ctr *formController) FormDestroyLocation(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	_, err = ctr.formMod.DeleteFormLocation(id)
	if err != nil {
		fmt.Println("DeleteProject", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	ctr.helper.AddLogBook(userID, 3, id)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success deleted data",
		"data":    nil,
	})
	return
}

// new function
func (ctr *formController) FormUserList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println(userID, organizationID)
	}

	formID, _ := strconv.Atoi(c.Param("formid"))
	param := c.Request.URL.Query().Get("type")
	searchKeyWord := c.Request.URL.Query().Get("search")

	if formID > 0 {
		//cekform id
		var frm tables.Forms
		frm.ID = formID
		getForm, err := ctr.formMod.GetFormRow(frm)
		if err != nil {
			fmt.Println("err: GetFormRow", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if getForm.ID == 0 {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is wrong",
				"error":   err,
			})
			return
		}

		formUserType := "respondent"
		if param != "" {
			formUserType = param
		}

		var fields tables.JoinFormUsers
		fields.FormID = formID
		fields.Type = formUserType

		var whreFormUserOrg bytes.Buffer
		if organizationID >= 1 {
			whreFormUserOrg.WriteString(" fuo.organization_id =" + strconv.Itoa(organizationID))
		} else {
			whreFormUserOrg.WriteString(" fuo.organization_id is null")
		}

		// whereString := " u.id in (select uo.user_id from usr.user_organizations uo where uo.organization_id= " + claims["organization_id"].(string) + ")"
		if searchKeyWord != "" {
			// whereString = " AND u.name ilike '%" + searchKeyWord + "%'  "
			whreFormUserOrg.WriteString("  AND u.name ilike '%" + searchKeyWord + "%'   ")
		}

		fuRows, err := ctr.formMod.GetFormUserToOrganizationRows(fields, whreFormUserOrg.String())
		if err != nil {
			fmt.Println("err: GetFormUserRows", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		allData, _ := ctr.formMod.GetFormUserToOrganizationRows(fields, whreFormUserOrg.String())

		// static row woner
		var res []objects.UserMember
		if formUserType == "admin" {
			//get author
			var whreForm tables.Forms
			whreForm.ID = getForm.ID
			getForm, _ := ctr.formMod.GetFormRow(whreForm)

			// static author ----------------------------------------------------------------
			var whre0 tables.Users
			whre0.ID = getForm.CreatedBy
			staticAuthor, _ := ctr.userMod.GetUserRows(whre0)
			fmt.Println("len(staticAuthor) ::::", len(staticAuthor), getForm.CreatedBy)
			for j := 0; j < len(staticAuthor); j++ {

				var eachStr1 objects.UserMember
				eachStr1.ID = staticAuthor[j].ID
				eachStr1.Email = staticAuthor[j].Email
				eachStr1.Phone = staticAuthor[j].Phone
				eachStr1.Name = staticAuthor[j].Name
				eachStr1.StatusID = 1
				eachStr1.StatusName = "Active"

				// permission
				var permissRes = []objects.FormUserPermissionJoin{
					{ID: 0, PermissionID: 6, PermissionName: "Edit Responden", Status: true},
					{ID: 0, PermissionID: 7, PermissionName: "View", Status: true},
					{ID: 0, PermissionID: 8, PermissionName: "Edit Form", Status: true},
					{ID: 0, PermissionID: 9, PermissionName: "Download", Status: true},
				}

				eachStr1.Permissions = permissRes

				res = append(res, eachStr1)
			}
		}

		var detail objects.DataRowsDetail
		detail.AllRows = len(allData)

		if len(fuRows) > 0 {

			// data user
			for i := 0; i < len(fuRows); i++ {
				var each objects.UserMember
				each.ID = fuRows[i].UserID
				each.Email = fuRows[i].Email
				each.Phone = fuRows[i].Phone
				each.Name = fuRows[i].UserName
				each.StatusID = fuRows[i].FormUserStatusID
				each.StatusName = fuRows[i].FormUserStatusName

				if formUserType == "admin" {

					var permissRes []objects.FormUserPermissionJoin

					var whreStr2 tables.FormUserPermissionJoin
					whreStr2.FormUserID = fuRows[i].ID
					getFormUserPerm, _ := ctr.permissMod.GetFormUserPermissionRows(whreStr2, "")
					for j := 0; j < len(getFormUserPerm); j++ {
						var eachStr2 objects.FormUserPermissionJoin

						eachStr2.ID = getFormUserPerm[j].ID
						eachStr2.PermissionID = getFormUserPerm[j].PermissionID
						eachStr2.PermissionName = getFormUserPerm[j].PermissionName
						eachStr2.Status = getFormUserPerm[j].Status

						permissRes = append(permissRes, eachStr2)
					}
					each.Permissions = permissRes
				}
				res = append(res, each)
			}

		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Form user data is available",
			"data":    res,
			"detail":  detail,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data form ID is not registered",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) UserGetFormList(c *gin.Context) {

	ID := c.Param("userid")
	userID, _ := strconv.Atoi(ID)

	formUserType := "respondent"

	var fields tables.JoinFormUsers
	fields.UserID = userID
	fields.Type = formUserType
	whreStr := ""

	fuRows, err := ctr.formMod.GetFormUserRows(fields, whreStr)
	if err != nil {
		fmt.Println("err: GetFormUserRows", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	if len(fuRows) > 0 {

		var res []objects.UserGetFormList

		// data user admin
		for i := 0; i < len(fuRows); i++ {
			var each objects.UserGetFormList
			each.ID = fuRows[i].FormID
			each.Name = fuRows[i].Name
			each.Description = fuRows[i].Description
			each.TotalRespon = 0
			each.FormStatus = true

			res = append(res, each)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Form user data is available",
			"data":    res,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User data is not available",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FormUserList__old(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	ID := c.Param("formid")
	iID, _ := strconv.Atoi(ID)

	if iID > 0 {

		var fields tables.FormUsers
		fields.UserID = userID
		fields.FormID = iID
		fields.Type = "respondent"
		result, err := ctr.formMod.GetFormUserRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {

			var res objects.FormUsers
			res.FormID = result.FormID
			res.UserID = result.UserID
			res.Name = result.Name
			res.Description = result.Description
			res.Notes = result.Notes
			res.ProfilePic = result.ProfilePic
			res.PeriodStartDate = result.PeriodStartDate
			res.PeriodEndDate = result.PeriodEndDate

			var FFfields tables.FormFields
			FFfields.FormID = result.FormID
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
		var fields tables.JoinFormUsers
		fields.UserID = userID
		fields.Type = "respondent"
		whreStr := ""

		result, err := ctr.formMod.GetFormUserRows(fields, whreStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var res []objects.FormUsers
			for i := 0; i < len(result); i++ {
				var each objects.FormUsers
				each.FormID = result[i].FormID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.Notes = result[i].Notes
				each.ProfilePic = result[i].ProfilePic
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate

				res = append(res, each)
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}

}

func (ctr *formController) FormUserCreate(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.Forms
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

	// fileName := "file_name.jpg"
	if reqData.ProfilePic != "" {
		//upload file
		// fileName = ctr.helper.UploadFileToOSS(reqData.ProfilePic, "form_file", "form")
	}

	var postData tables.Forms
	postData.Name = reqData.Name
	postData.Description = reqData.Description
	postData.Notes = reqData.Notes
	postData.PeriodStartDate = reqData.PeriodStartDate
	postData.PeriodEndDate = reqData.PeriodEndDate
	postData.ProfilePic = reqData.ProfilePic
	postData.EncryptCode = helpers.EncodeToString(6)

	res, err := ctr.formMod.InsertForm(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// ini akan dihapus
	var postData2 tables.FormUsers
	postData2.FormID = res.ID
	postData2.UserID = userID

	_, err = ctr.formMod.ConnectFormUser(postData2)
	if err != nil {
		fmt.Println("InsertFormUser", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var obj objects.FormResponse
	obj.ID = res.ID

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created form",
		"data":    obj,
	})
	return
}

func (ctr *formController) FormUserUpdate(c *gin.Context) {

	ID := c.Param("formid")
	if ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Field ID in URL is required",
		})
		return
	}
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
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

	// fileName := "file_name.jpg"
	if reqData.ProfilePic != "" {
		//upload file
		// fileName = ctr.helper.UploadFileToOSS(reqData.ProfilePic, "form_file", "form")
	}

	var postData tables.Forms
	postData.Name = reqData.Name
	postData.Description = reqData.Description
	postData.Notes = reqData.Notes
	postData.PeriodStartDate = reqData.PeriodStartDate
	postData.PeriodEndDate = reqData.PeriodEndDate
	postData.ProfilePic = reqData.ProfilePic
	postData.EncryptCode = helpers.EncodeToString(6)

	res, err := ctr.formMod.UpdateForm(formID, postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success updated form",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed updated form",
		})
		return
	}
}

func (ctr *formController) FormUserStatusUpdate(c *gin.Context) {

	ID := c.Param("formid")
	if ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Field ID in URL is required",
		})
		return
	}
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var reqData objects.FormUserStatus
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

	var postData tables.FormUsers
	postData.FormUserStatusID = reqData.UserStatusID
	fmt.Println("err", formID, reqData.UserID, postData)
	res, err := ctr.formMod.UpdateFormUser(reqData.UserID, formID, postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success updated user form status",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed updated user form",
		})
		return
	}
}

func (ctr *formController) FormUserConnect(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	myCompanyID := 0
	if len(claims) >= 5 {
		myCompanyID, _ = strconv.Atoi(claims["organization_id"].(string))
	}
	fmt.Println("my companyID :::", myCompanyID)

	var reqData objects.FormUsers
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
	var postData objects.InputFormUserOrganizations
	postData.FormID = reqData.FormID
	postData.UserID = reqData.UserID
	postData.Type = "respondent"
	postData.OrganizationID = myCompanyID

	res, err := ctr.formMod.ConnectFormUserOrg(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var obj objects.FormResponse
	obj.ID = res.ID

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created form",
		"data":    obj,
	})
	return
}

func (ctr *formController) FormUserDisconnect(c *gin.Context) {

	var reqData objects.FormUsers
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

	_, err = ctr.formMod.DeleteFormUserOrg(reqData.UserID, reqData.FormID)
	if err != nil {
		fmt.Println("DeleteFormUser", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success disconnect user from form",
		"data":    nil,
	})
	return
}

func (ctr *formController) FieldCreate(c *gin.Context) {

	var reqData objects.FormFields
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

	var dataField tables.FormFields
	// if reqData.CityID == 0 && reqData.DistrictID == 0 && reqData.SubDistrictID == 0 {
	// 	dataField.AddressType = reqData.AddressType
	// 	dataField.ProvinceID = reqData.ProvinceID
	// }
	// if reqData.DistrictID == 0 && reqData.SubDistrictID == 0 {
	// 	dataField.AddressType = reqData.AddressType
	// 	dataField.ProvinceID = reqData.ProvinceID
	// 	dataField.CityID = reqData.CityID
	// }
	// if reqData.SubDistrictID == 0 {
	// 	dataField.AddressType = reqData.AddressType
	// 	dataField.ProvinceID = reqData.ProvinceID
	// 	dataField.CityID = reqData.CityID
	// 	dataField.DistrictID = reqData.DistrictID
	// }

	dataField.ParentID = reqData.ParentID
	dataField.FormID = reqData.FormID
	dataField.FieldTypeID = reqData.FieldTypeID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.Option = reqData.Option
	dataField.ConditionType = reqData.ConditionType
	dataField.UpperlowerCaseType = reqData.UpperlowerCaseType
	dataField.IsMultiple = reqData.IsMultiple
	dataField.IsRequired = reqData.IsRequired
	dataField.SortOrder = reqData.SortOrder
	dataField.TagLocIcon = reqData.TagLocIcon
	dataField.TagLocColor = reqData.TagLocColor
	// dataField.AddressType = reqData.AddressType
	// dataField.ProvinceID = reqData.ProvinceID
	// dataField.CityID = reqData.CityID
	// dataField.DistrictID = reqData.DistrictID
	// dataField.SubDistrictID = reqData.SubDistrictID
	// dataField.CurrencyType = reqData.CurrencyType
	// dataField.Currency = reqData.Currency
	dataField.IsCountryPhoneCode = reqData.IsCountryPhoneCode

	res, err := ctr.formFieldMod.InsertFormField(dataField)
	if err != nil {
		fmt.Println("InsertFormField", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// insert rule
	if res.ID > 0 && reqData.ConditionRulesID > 0 {

		var dataRule tables.FormFieldConditionRules
		dataRule.FormFieldID = res.ID
		dataRule.ConditionRuleID = reqData.ConditionRulesID
		dataRule.Value1 = reqData.ConditionRuleValue1
		dataRule.Value2 = reqData.ConditionRuleValue2
		dataRule.ErrMsg = reqData.ConditionRuleMsg
		dataRule.TabMaxOnePerLine = reqData.TabMaxOnePerLine
		dataRule.TabEachLineRequire = reqData.TabEachLineRequire

		_, err = ctr.ruleMod.InsertFormFieldRule(dataRule)
		if err != nil {
			fmt.Println("InsertFormFieldRule", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else if res.ID > 1 && reqData.ConditionRulesID == 0 {

		_, err = ctr.ruleMod.DeleteFormFieldRulePrimary(res.ID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldRule", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	if res.ID > 0 && reqData.Image != "" {

		// fileName := ctr.helper.UploadFileToOSS(reqData.Image, "form_file_pic", "form_field")

		var postData tables.FormFieldPics
		postData.FormFieldID = res.ID
		postData.Pic = reqData.Image
		_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
		if err != nil {
			fmt.Println("InsertFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else if res.ID > 0 && reqData.Image == "" {

		_, err = ctr.formFieldMod.DeleteFormFieldPic(res.ID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	if res.ID > 0 {
		var obj objects.FormResponse
		obj.ID = res.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created form field",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed created form field",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FieldUpdate(c *gin.Context) {

	ID := c.Param("fieldid")
	if ID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Field ID in URL is required",
		})
		return
	}
	fieldID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var reqData objects.FormFields
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

	var dataField tables.FormFields
	dataField.ParentID = reqData.ParentID
	dataField.FormID = reqData.FormID
	dataField.FieldTypeID = reqData.FieldTypeID
	dataField.FormID = reqData.FormID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.Option = reqData.Option
	dataField.ConditionType = reqData.ConditionType
	dataField.UpperlowerCaseType = reqData.UpperlowerCaseType
	dataField.IsMultiple = reqData.IsMultiple
	dataField.IsRequired = reqData.IsRequired
	dataField.SortOrder = reqData.SortOrder
	// dataField.AddressType = reqData.AddressType
	// dataField.ProvinceID = reqData.ProvinceID
	// dataField.CityID = reqData.CityID
	// dataField.DistrictID = reqData.DistrictID
	// dataField.SubDistrictID = reqData.SubDistrictID
	// dataField.CurrencyType = reqData.CurrencyType
	// dataField.Currency = reqData.Currency
	// dataField.IsCountryPhoneCode = reqData.IsCountryPhoneCode

	res, err := ctr.formFieldMod.UpdateFormField(fieldID, dataField)
	if err != nil {
		fmt.Println("InsertFormField", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed update field",
		})
		return
	}

	// insert rule / single validation
	if fieldID > 0 && reqData.ConditionRulesID > 0 {

		var fields tables.FormFieldConditionRules
		fields.FormFieldID = fieldID
		whereString := "condition_parent_field_id is null"
		getRule, err := ctr.ruleMod.GetFormFieldRuleRow(fields, whereString)

		if getRule.ID > 0 {
			var dataRule tables.FormFieldConditionRules
			dataRule.FormFieldID = fieldID
			dataRule.ConditionRuleID = reqData.ConditionRulesID
			dataRule.Value1 = reqData.ConditionRuleValue1
			dataRule.Value2 = reqData.ConditionRuleValue2
			dataRule.ErrMsg = reqData.ConditionRuleMsg
			dataRule.TabMaxOnePerLine = reqData.TabMaxOnePerLine
			dataRule.TabEachLineRequire = reqData.TabEachLineRequire
			// dataRule.ConditionAllRight = reqData.C

			_, err = ctr.ruleMod.UpdateFormFieldRule(getRule.ID, dataRule)
			if err != nil {
				fmt.Println("InsertFormFieldRule", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		} else {
			var dataRule tables.FormFieldConditionRules
			dataRule.FormFieldID = fieldID
			dataRule.ConditionRuleID = reqData.ConditionRulesID
			dataRule.Value1 = reqData.ConditionRuleValue1
			dataRule.Value2 = reqData.ConditionRuleValue2
			dataRule.ErrMsg = reqData.ConditionRuleMsg
			dataRule.TabMaxOnePerLine = reqData.TabMaxOnePerLine
			dataRule.TabEachLineRequire = reqData.TabEachLineRequire

			_, err = ctr.ruleMod.InsertFormFieldRule(dataRule)
			if err != nil {
				fmt.Println("InsertFormFieldRule", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}
	} else if fieldID > 0 && reqData.ConditionRulesID == 0 {

		_, err = ctr.ruleMod.DeleteFormFieldRulePrimary(fieldID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldRule", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

	}

	if fieldID > 0 && reqData.Image != "" {

		// fileName := ctr.helper.UploadFileToOSS(reqData.Image, "form_file_pic", "form_field")
		var fields tables.FormFieldPics
		fields.FormFieldID = fieldID
		getFieldPic, _ := ctr.formFieldMod.GetFormFieldPicRow(fields)
		if err != nil {
			newErr := errors.New("record not found")
			if err != newErr {
				fmt.Println("GetFormFieldPicRow ---", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}

		var postData tables.FormFieldPics
		postData.FormFieldID = fieldID
		postData.Pic = reqData.Image
		if getFieldPic.ID > 0 {
			_, err = ctr.formFieldMod.UpdateFormFieldPic(getFieldPic.ID, postData)
			if err != nil {
				fmt.Println("UpdateFormFieldPic", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		} else {
			_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
			if err != nil {
				fmt.Println("InsertFormFieldPic", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}

	} else if fieldID > 0 && reqData.Image == "" {

		_, err = ctr.formFieldMod.DeleteFormFieldPic(fieldID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	var obj objects.FormResponse
	obj.ID = fieldID

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success created form",
		"data":    obj,
	})
	return
}

func (ctr *formController) FieldGroupCreate(c *gin.Context) {

	var reqData objects.FormFieldGroup
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

	var dataField tables.FormFields
	dataField.FormID = reqData.FormID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.SortOrder = reqData.SortOrder

	res, err := ctr.formFieldMod.InsertFormField(dataField)
	if err != nil {
		fmt.Println("InsertFormField group", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed save form field",
			"message": err,
		})
		return
	}

	if res.ID > 0 {
		var obj objects.FormResponse
		obj.ID = res.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created group field",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed created group field",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FieldGroupUpdate(c *gin.Context) {

	ID := c.Param("fieldid")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Field ID in URL is required",
		})
		return
	}
	fieldID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var reqData objects.FormFieldGroup
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

	var dataField tables.FormFields
	dataField.FormID = reqData.FormID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.SortOrder = reqData.SortOrder

	res, err := ctr.formFieldMod.UpdateFormField(fieldID, dataField)
	if err != nil {
		fmt.Println("InsertFormField group", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed save form field",
			"message": err,
		})
		return
	}

	if res {

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created group field",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed created group field",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FieldSectionCreate(c *gin.Context) {

	var reqData objects.FormFieldSection
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

	var dataField tables.FormFields
	dataField.ParentID = reqData.ParentID
	dataField.FormID = reqData.FormID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.IsSection = true
	dataField.SectionColor = reqData.SectionColor
	dataField.SortOrder = reqData.SortOrder

	res, err := ctr.formFieldMod.InsertFormField(dataField)
	if err != nil {
		fmt.Println("InsertFormField", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res.ID > 0 && reqData.Image != "" {

		// fileName := ctr.helper.UploadFileToOSS(reqData.Image, "form_file_pic", "form_field")

		var postData tables.FormFieldPics
		postData.FormFieldID = res.ID
		postData.Pic = reqData.Image
		_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
		if err != nil {
			fmt.Println("InsertFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else if res.ID > 0 && reqData.Image == "" {

		_, err = ctr.formFieldMod.DeleteFormFieldPic(res.ID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	if res.ID > 0 {
		var obj objects.FormResponse
		obj.ID = res.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created field section",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed created field section",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FieldSectionUpdate(c *gin.Context) {

	ID := c.Param("fieldid")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Field ID in URL is required",
		})
		return
	}
	fieldID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	var reqData objects.FormFieldSection
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

	var dataField tables.FormFields
	dataField.ParentID = reqData.ParentID
	dataField.FormID = reqData.FormID
	dataField.Label = reqData.Label
	dataField.Description = reqData.Description
	dataField.IsSection = true
	dataField.SectionColor = reqData.SectionColor
	dataField.SortOrder = reqData.SortOrder

	res, err := ctr.formFieldMod.UpdateFormField(fieldID, dataField)
	if err != nil {
		fmt.Println("UpdateFormField", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if res == true && reqData.Image != "" {

		// fileName := ctr.helper.UploadFileToOSS(reqData.Image, "form_file_pic", "form_field")
		var whrPic tables.FormFieldPics
		whrPic.FormFieldID = fieldID
		chekPic, _ := ctr.formFieldMod.GetFormFieldPicRow(whrPic)

		if chekPic.ID > 0 {
			var updateData tables.FormFieldPics
			updateData.Pic = reqData.Image
			_, err := ctr.formFieldMod.UpdateFormFieldPic(chekPic.ID, updateData)
			if err != nil {
				fmt.Println("InsertFormFieldPic", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		} else {
			var postData tables.FormFieldPics
			postData.FormFieldID = fieldID
			postData.Pic = reqData.Image
			_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
			if err != nil {
				fmt.Println("InsertFormFieldPic", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}
	} else if res == true && reqData.Image == "" {

		_, err = ctr.formFieldMod.DeleteFormFieldPic(fieldID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldPic", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	if res {
		var obj objects.FormResponse
		obj.ID = fieldID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success updated field section",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed update field section",
			"data":    nil,
		})
		return
	}
}

// field tanpa group
func (ctr *formController) FieldList(c *gin.Context) {

	ID := c.Param("formid")
	fieldTypeID := c.Request.URL.Query().Get("field_type_id")

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

		if fieldTypeID != "" {
			fields.FieldTypeID, err = strconv.Atoi(fieldTypeID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": false,
					"error":  err,
				})
				return
			}
		}

		result, err := ctr.formFieldMod.GetFormFieldNotParentRows(fields, "")
		fmt.Println(len(result))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {
			var res []objects.FormFields
			for i := 0; i < len(result); i++ {

				var each objects.FormFields
				each.ID = result[i].ID
				each.Label = result[i].Label
				each.FieldTypeID = result[i].FieldTypeID
				each.FormID = result[i].FormID
				each.Description = result[i].Description
				each.IsMultiple = result[i].IsMultiple
				each.IsRequired = result[i].IsRequired
				each.Option = result[i].Option
				each.ConditionRulesID = result[i].ConditionRuleID
				each.ConditionRuleValue1 = result[i].Value1
				each.ConditionRuleValue2 = result[i].Value2
				each.ConditionRuleMsg = result[i].ErrMsg
				each.ConditionParentFieldID = result[i].ConditionParentFieldID
				each.IsSection = result[i].IsSection
				each.SectionColor = result[i].SectionColor
				each.AddressType = result[i].AddressType
				each.ProvinceID = result[i].ProvinceID
				each.CityID = result[i].CityID
				each.DistrictID = result[i].DistrictID
				each.SubDistrictID = result[i].SubDistrictID
				each.CurrencyType = result[i].CurrencyType
				each.Currency = result[i].Currency
				each.IsCountryPhoneCode = result[i].IsCountryPhoneCode

				var fields tables.FormFieldPics
				fields.FormFieldID = result[i].ID
				pic, _ := ctr.formFieldMod.GetFormFieldPicRow(fields)

				each.Image = pic.Pic

				if result[i].IsCondition {
					var fcs tables.FormFieldConditionRules
					fcs.FormFieldID = result[i].ID
					stringFields := "condition_parent_field_id is not null"
					getRules, err := ctr.ruleMod.GetFormFieldRuleRows(fcs, stringFields)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}

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
					each.Conditions = condRules
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
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data is not available",
		})
		return
	}
}

// field dgn group
func (ctr *formController) FieldGroupList(c *gin.Context) {

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

		//groups
		var fields tables.FormFields
		fields.FormID = formID

		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		// var fields tables.FormFields
		// fields.FormID = iID
		// result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		if len(result) > 0 {

			var res []objects.FormFieldGroups
			for i := 0; i < len(result); i++ {

				var each objects.FormFieldGroups
				each.ID = result[i].ID
				each.Label = result[i].Label
				each.Description = result[i].Description

				//form fileds
				var getFieldres []objects.FormFields
				var whre tables.FormFields
				whre.FormID = formID
				whre.ParentID = result[i].ID
				getFields, _ := ctr.formFieldMod.GetFormFieldRows(whre)

				if len(getFields) > 0 {

					for j := 0; j < len(getFields); j++ {
						fmt.Println("getFields[j].ID ------------", getFields[j].ID)
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
						// eachj.SortOrder = getFields[j].SortOrder
						eachj.Image = getFields[j].Image
						eachj.ConditionRulesID = getFields[j].ConditionRuleID
						eachj.ConditionRuleValue1 = getFields[j].Value1
						eachj.ConditionRuleValue2 = getFields[j].Value2
						eachj.ConditionRuleMsg = getFields[j].ErrMsg
						eachj.ConditionParentFieldID = getFields[j].ConditionParentFieldID
						// eachj.TabMaxOnePerLine = getFields[j].TabMaxOnePerLine
						// eachj.TabEachLineRequire = getFields[j].TabEachLineRequire
						// eachj.QrCode = getFields[j].QrCode
						// eachj.Conditions = getFields[j].Conditions

						getFieldres = append(getFieldres, eachj)
					}
				}
				each.FormFields = getFieldres
				//group condition
				if result[i].IsCondition {
					var fcs tables.FormFieldConditionRules
					fcs.FormFieldID = result[i].ID
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
					each.Conditions = condRules
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

// field dgn group 2nd
func (ctr *formController) Field2GroupList(c *gin.Context) {

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

		// get groups
		var fields tables.FormFields
		fields.FormID = formID

		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(result) > 0 {

			var res []objects.FormFieldGroups
			for i := 0; i < len(result); i++ {

				var each objects.FormFieldGroups
				each.ID = 0
				each.Label = ""
				each.Description = ""
				each.SortOrder = 0

				if result[i].FieldTypeID == 0 && result[i].IsSection == false {
					each.ID = result[i].ID
					each.Label = result[i].Label
					each.Description = result[i].Description
					each.SortOrder = result[i].SortOrder

					//form fileds
					var getFieldres []objects.FormFields
					var whre tables.FormFields
					whre.ParentID = result[i].ID
					whre.FormID = formID
					getFields, _ := ctr.formFieldMod.GetFormFieldRows(whre)

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
							// eachj.TabMaxOnePerLine = getFields[j].TabMaxOnePerLine
							// eachj.TabEachLineRequire = getFields[j].TabEachLineRequire
							// eachj.QrCode = getFields[j].QrCode

							//group condition
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
							getFieldres = append(getFieldres, eachj)
						}
					}
					each.FormFields = getFieldres

				} else {

					// tidak memiliki group
					each.ID = 0

					each.Label = ""
					each.Description = ""
					each.SortOrder = 0

					var getFieldres []objects.FormFields
					var whre tables.FormFields
					whre.ID = result[i].ID
					getFields, _ := ctr.formFieldMod.GetFormFieldRows(whre)

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

							//group condition
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

							getFieldres = append(getFieldres, eachj)
						}
					}
					each.FormFields = getFieldres

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

// latest func
func (ctr *formController) FieldGroup3List(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	var formID int
	ID := c.Param("formid")
	if ID != "" {
		ID, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Println("InsertFormFieldRule", err)
			c.JSON(http.StatusNoContent, gin.H{
				"error": err,
			})
			return
		}

		formID = ID
	}

	formCode := c.Request.URL.Query().Get("code")
	if formCode != "" {

		// cek user privilage by form user
		var whrFrmUser tables.FormUsers
		whrFrmUser.UserID = userID
		whrFrmUser.FormID = formID
		checkFormUser, _ := ctr.formMod.GetFormUserRow(whrFrmUser)
		if checkFormUser.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Your account does not have privilage in this form.",
			})
			return
		}

		var whrFrm tables.Forms
		whrFrm.EncryptCode = formCode
		getFormByCode, err := ctr.formMod.GetFormRow(whrFrm)
		if err != nil {
			fmt.Println("get form by code", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err,
			})
			return
		}

		if getFormByCode.EncryptCode != "" {
			formID = getFormByCode.ID
		}
	}

	if formID > 0 {
		fmt.Println(formID)
		// get groups
		var fields tables.FormFields
		fields.FormID = formID
		// fields.FieldTypeID = -2

		result, err := ctr.formFieldMod.GetFormFieldRows(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println("field grou 3 - len(result) :::--- ", len(result))

		if len(result) > 0 {

			var res []objects.FormField3Groups
			for i := 0; i < len(result); i++ {

				if result[i].FieldTypeID == 0 && result[i].IsSection == false {

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

					//form fileds
					var getFieldres []objects.FormFields
					var whre tables.FormFields
					whre.ParentID = result[i].ID
					whre.FormID = formID
					getFields, _ := ctr.formFieldMod.GetFormFieldRows(whre)

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
							// eachj.TabMaxOnePerLine = getFields[j].TabMaxOnePerLine
							// eachj.TabEachLineRequire = getFields[j].TabEachLineRequire
							// eachj.QrCode = getFields[j].QrCode

							//group condition
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
							getFieldres = append(getFieldres, eachj)
						}
					}

					each.FormFields = getFieldres

					res = append(res, each)

				} else {

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

						res = append(res, each)
					}
				}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is empty",
		})
		return
	}
}

func (ctr *formController) FieldDestroy(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println("controller : FieldDestroy", userID)

	id, err := strconv.Atoi(c.Param("fieldid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err,
		})
		return
	}

	_, err = ctr.formFieldMod.DeleteFormField(id)
	if err != nil {
		fmt.Println("DeleteFormUser", err)
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

func (ctr *formController) FieldTypeList(c *gin.Context) {

	ID := c.Param("id")
	iID, _ := strconv.Atoi(ID)

	if iID > 0 {

		var fields tables.FieldTypes
		fields.ID = iID

		result, err := ctr.ftMod.GetFieldTypeRow(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if result.ID > 0 {
			var res objects.FieldTypes
			res.ID = result.ID
			res.Code = result.Code
			res.Name = result.Name

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

		var fields tables.FieldTypes
		results, err := ctr.ftMod.GetFieldTypeRows(fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.FieldTypes
			for i := 0; i < len(results); i++ {
				var each objects.FieldTypes
				each.ID = results[i].ID
				each.Code = results[i].Code
				each.Name = results[i].Name
				each.Info = results[i].Info

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

func (ctr *formController) ConditionRuleList(c *gin.Context) {

	ruleType := c.Param("type")

	var fields tables.ConditionRules
	langID := 1

	if ruleType == "text" {
		fields.Type = "text"
	} else if ruleType == "numeric" {
		fields.Type = "numeric"
	} else if ruleType == "length" {
		fields.Type = "text_length"
	} else if ruleType == "choice" {
		fields.Type = "number_choice"
	} else if ruleType == "tab" {
		fields.Type = "tab_howtoinput"
	}

	result, err := ctr.ruleMod.ConditionRuleList(langID, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if len(result) > 0 {

		var res []objects.ConditionRulesRes
		for i := 0; i < len(result); i++ {
			var each objects.ConditionRulesRes
			each.ID = result[i].ID
			each.Name = result[i].Name
			each.Code = result[i].Code

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

func (ctr *formController) FieldSaveImage(c *gin.Context) {

	var reqData objects.FieldFile
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

	// upload file
	fileName := ""
	if reqData.File != "" {
		fileName = ctr.helper.UploadFileToOSS(reqData.File, "file_data", reqData.FileType)
	}

	var postData tables.FormFieldTempAssets
	postData.FormID = reqData.FormID
	postData.FormFieldID = reqData.FieldID
	postData.TempAsset = fileName

	check, res, err := ctr.formMod.InsertFieldFile(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	if check == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed upload file",
			"data":    nil,
		})
		return
	}

	var obj objects.FieldFileRes
	obj.FormFileID = res.ID
	obj.UrlFile = fileName

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success uploaded file",
		"data":    obj,
	})
	return
}

func (ctr *formController) FieldConditionSave(c *gin.Context) {

	var reqData objects.FieldConditionSave
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

	if len(reqData.Conditions) > 0 {

		var cr tables.FormFieldConditionRules
		cr.FormFieldID = reqData.FormFieldID
		cr.ConditionAllRight = reqData.ConditionAllRight

		//udpate field is conddition
		var dataUpdate tables.FormFields
		dataUpdate.IsCondition = true
		_, err := ctr.formFieldMod.UpdateFormOnly(reqData.FormFieldID, dataUpdate)
		if err != nil {
			fmt.Println("Err :UpdateFormField", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		// clear old data
		_, err = ctr.ruleMod.DeleteFormFieldRule(reqData.FormFieldID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldRule", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		// var save tables.FormFieldConditionRules
		for i := 0; i < len(reqData.Conditions); i++ {

			//insert new rules
			cr.ConditionParentFieldID = reqData.Conditions[i].ParentFieldID
			cr.ConditionRuleID = reqData.Conditions[i].ConditionRuleID
			cr.Value1 = reqData.Conditions[i].ConditionRuleValue1
			cr.Value2 = reqData.Conditions[i].ConditionRuleValue2

			_, err = ctr.ruleMod.InsertFormFieldRule(cr)
			if err != nil {
				fmt.Println(" Err : InsertFormFieldRule", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created condition form",
			"data":    nil,
		})
		return
	} else {
		// clear old data
		_, err = ctr.ruleMod.DeleteFormFieldRule(reqData.FormFieldID)
		if err != nil {
			fmt.Println("Err :DeleteFormFieldRule", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success deleted condition rule",
			"data":    nil,
		})
		return
	}

	// c.JSON(http.StatusBadRequest, gin.H{
	// 	"status":  false,
	// 	"message": "Failed saved condition rule",
	// 	"data":    nil,
	// })
	// return

}

func (ctr *formController) FormShare(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	fmt.Println(authorID)

	fID, _ := strconv.Atoi(c.Param("formid"))

	needValidate := c.Request.URL.Query().Get("need_validate")

	var whre tables.Forms
	whre.ID = fID
	getForm, err := ctr.formMod.GetFormRow(whre)
	if err != nil {
		fmt.Println("Err :DeleteFormFieldRule", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if getForm.ID > 0 {
		var res objects.ResFormShare
		res.Code = getForm.EncryptCode
		res.SenderOrganizationID = organizationID
		// scheme := "http"
		// if c.Request.TLS != nil {
		// 	scheme = "https"
		// }
		// res.Url = scheme + "://" + c.Request.Host + "/?code=" + getForm.EncryptCode + "&nval=" + needValidate
		fmt.Println(needValidate)

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

func (ctr *formController) FormGetShare(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	respondenID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.SendShareCode
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

	if reqData.ShareCode != "" {
		var whre tables.Forms
		whre.EncryptCode = reqData.ShareCode
		getForm, err := ctr.formMod.GetFormRow(whre)
		if err != nil {
			fmt.Println("Err :GetFormRow", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		fmt.Println("getForm :::", getForm)

		if getForm.ID >= 1 {

			var whereComp tables.FormOrganizations
			whereComp.FormID = getForm.ID
			getComp, err := ctr.formMod.GetFormOrganization(whereComp)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			// join to company
			var dataPost objects.UserOrganizations
			dataPost.UserID = respondenID
			dataPost.OrganizationID = getComp.OrganizationID

			userComp, err := ctr.compMod.ConnectedUserCompanies(dataPost)
			if err != nil {
				fmt.Println(err)
				// cek duplicate key
				if errors.As(err, &ctr.pgErr) {

					fmt.Println("err", ctr.pgErr.Code)
					if ctr.pgErr.Code == "23505" {

						// auto join to form
						formStatusID := 1
						if reqData.NeedValidate == true {
							formStatusID = 2 //(2 = non active, 3 = neww validate)
						}
						var dataFU tables.FormUsers
						dataFU.FormID = getForm.ID
						dataFU.UserID = respondenID
						dataFU.FormUserStatusID = formStatusID // non aktif/perlu validasi
						resultFUO1, err := ctr.formMod.ConnectFormUser(dataFU)
						if err != nil {

							// cek duplicate key
							if errors.As(err, &ctr.pgErr) {
								fmt.Println("err", ctr.pgErr.Code)
								if ctr.pgErr.Code == "23505" {

									var shareCode objects.ShareCodeResp
									shareCode.FormID = getForm.ID
									shareCode.IsFirstConnected = false

									c.JSON(http.StatusOK, gin.H{
										"status":  true,
										"message": "Anda telah tergabung ke dalam form " + getForm.Name,
										"data":    shareCode,
									})
									return
								}
							}

							c.JSON(http.StatusBadRequest, gin.H{
								"error": err,
							})
							return
						}

						// insert to form user organization
						var postFUO tables.FormUserOrganizations
						postFUO.FormUserID = resultFUO1.ID
						postFUO.OrganizationID = getComp.OrganizationID
						_, err = ctr.formMod.GenerateFormUserOrg(postFUO)
						if err != nil {
							c.JSON(http.StatusBadRequest, gin.H{
								"error":   err,
								"message": err.Error(),
								"status":  http.StatusBadRequest,
							})
							return
						}

						var shareCode objects.ShareCodeResp
						shareCode.FormID = getForm.ID
						shareCode.IsFirstConnected = true
						//Kamu telah bergabung dengan perusahaan dan saat ini tergabung dengan form
						c.JSON(http.StatusOK, gin.H{
							"status":  true,
							"message": "Anda telah tergabung ke dalam form " + getForm.Name,
							"data":    shareCode,
						})
						return
					}
				}

				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			if userComp.ID >= 1 {
				var saveData tables.UserOrganizationRoles
				saveData.UserOrganizationID = userComp.ID
				saveData.RoleID = 5 //(5 = responden)

				_, err = ctr.compMod.InsertUserCompaniyToRole(saveData)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// auto join to form
				formStatusID := 1
				if reqData.NeedValidate == true {
					formStatusID = 2 //(2 = non active, 3 = neww validate)
				}
				var dataFU tables.FormUsers
				dataFU.FormID = getForm.ID
				dataFU.UserID = respondenID
				dataFU.FormUserStatusID = formStatusID // non aktif/perlu validasi
				resultFUO, err := ctr.formMod.ConnectFormUser(dataFU)
				if err != nil {

					// cek duplicate key
					if errors.As(err, &ctr.pgErr) {
						fmt.Println("err", ctr.pgErr.Code)
						if ctr.pgErr.Code == "23505" {
							var shareCode objects.ShareCodeResp
							shareCode.FormID = getForm.ID
							shareCode.IsFirstConnected = false

							c.JSON(http.StatusOK, gin.H{
								"status":  true,
								"message": "Anda telah tergabung ke dalam form " + getForm.Name + ".",
								"data":    shareCode,
							})
							return
						}
					}
				}

				// insert to form user orgnization
				var postFUO tables.FormUserOrganizations
				postFUO.FormUserID = resultFUO.ID
				postFUO.OrganizationID = getComp.OrganizationID
				_, err = ctr.formMod.GenerateFormUserOrg(postFUO)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   err,
						"message": err.Error(),
						"status":  http.StatusBadRequest,
					})
					return
				}

				fmt.Println("getComp.OrganizationID ::", resultFUO.ID, getComp.OrganizationID)

				var shareCode objects.ShareCodeResp
				shareCode.FormID = getForm.ID
				shareCode.IsFirstConnected = false
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Anda tergabung ke dalam form " + getForm.Name,
					"data":    shareCode,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Gagal bergabung ke dalam form" + getForm.Name,
					"data":    nil,
				})
				return
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form code is not found",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form code is empty or wrong",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FormUpdateStatus(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))

	fID, _ := strconv.Atoi(c.Param("formid"))

	var reqData objects.FormStatus
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

	if fID > 0 {

		var data tables.Forms
		data.UpdatedBy = authorID
		data.FormStatusID = reqData.StatusID

		formStatus := ""
		if reqData.StatusID == 1 {
			formStatus = "ACTIVED"
		} else if reqData.StatusID == 2 {
			formStatus = "INACTIVED"
		} else if reqData.StatusID == 3 {
			formStatus = "ARCHIVED"
			data.ArchivedAt = helpers.DateNow()
		}

		saveUpdate, err := ctr.formMod.UpdateFormStatus(fID, data)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if saveUpdate {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Form update is successful",
				"data":    "Form is " + formStatus,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed form update",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "form ID is required",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FieldSortOrderSave(c *gin.Context) {

	fID, _ := strconv.Atoi(c.Param("formid"))

	var reqData objects.FormSortOrder
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println("err -----", err)
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
	if fID > 0 {
		// var whre tables.FormFields
		// whre.FormID = fID
		// results, err := ctr.formFieldMod.GetFormFieldRows(whre)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		results := reqData.SortOrders
		// fmt.Println("zzzz", len(results))
		if len(results) > 0 {
			// update field
			saveUpdate := false
			for i := 0; i < len(results); i++ {

				var data tables.FormFields
				data.SortOrder = i + 1
				saveUpdate, err = ctr.formFieldMod.UpdateFormFieldSort(results[i].ID, data)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
			}

			if saveUpdate {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Form update is successful",
					"data":    nil,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed form update",
					"data":    nil,
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Field is empty",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormUserAdminFrm(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	myOrganizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	fmt.Println(myOrganizationID)

	formID, _ := strconv.Atoi(c.Param("formid"))
	// os.Exit(0)
	if formID > 0 {

		//check author
		var wh tables.Forms
		wh.ID = formID
		getForm, _ := ctr.formMod.GetFormRow(wh)
		fmt.Println("roleID", roleID)
		// if author is init form sel permission

		var res []objects.FormUserPermissionRes

		if (getForm.CreatedBy > 0 && getForm.CreatedBy == authorID) || roleID == 1 { // ID 1 (role admin auto all permission)
			fmt.Println("cond 1 ::", formID, " +++++ ", getForm.CreatedBy, "==", authorID)
			var permissRes = []objects.FormUserPermissionJoin{
				{ID: 0, PermissionID: 6, PermissionName: "Edit Responden", Status: true},
				{ID: 0, PermissionID: 7, PermissionName: "View", Status: true},
				{ID: 0, PermissionID: 8, PermissionName: "Edit Form", Status: true},
				{ID: 0, PermissionID: 9, PermissionName: "Download", Status: true},
			}

			// check form ini adalah form company dri user admin ini atau bukan
			checkFormToOwner, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID}) // if

			for i := 0; i < len(permissRes); i++ {

				var each objects.FormUserPermissionRes
				each.ID = permissRes[i].ID
				each.FormUserID = permissRes[i].FormUserID
				each.PermissionID = permissRes[i].PermissionID
				each.PermissionName = permissRes[i].PermissionName
				each.Status = permissRes[i].Status
				if permissRes[i].PermissionID == 8 && checkFormToOwner.OrganizationID != myOrganizationID {
					each.Status = false
				}

				res = append(res, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available.",
				"data":    res,
			})
			return
		}

		var whre tables.FormUserPermissionJoin
		shreStr2 := `fu.form_id = ` + strconv.Itoa(formID) + ` and fu.user_id = ` + strconv.Itoa(authorID)
		getFormuserPermiss, err := ctr.permissMod.GetFormUserPermissionRows(whre, shreStr2)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		// fmt.Println("len(getFormuserPermiss)", len(getFormuserPermiss), strconv.Itoa(formID), strconv.Itoa(authorID))

		// check form ini adalah form company dri user admin ini atau bukan
		checkFormOrganization, _ := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID, OrganizationID: myOrganizationID}) // if

		// var res []objects.FormUserPermissionRes
		for i := 0; i < len(getFormuserPermiss); i++ {

			if getFormuserPermiss[i].PermissionID == 7 {
				getFormuserPermiss[i].Status = true
			}

			var each objects.FormUserPermissionRes
			each.ID = getFormuserPermiss[i].ID
			each.FormUserID = getFormuserPermiss[i].FormUserID
			each.PermissionID = getFormuserPermiss[i].PermissionID
			each.PermissionName = getFormuserPermiss[i].PermissionName
			each.Status = getFormuserPermiss[i].Status
			if getFormuserPermiss[i].PermissionID == 8 && checkFormOrganization.OrganizationID <= 0 {
				each.Status = false
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
			"message": "Failed author ID is not found",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FormUserAdminFrmConnect(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	myCompanyID := 0
	if len(claims) >= 5 {
		myCompanyID, _ = strconv.Atoi(claims["organization_id"].(string))
	}
	fmt.Println("my companyID :::", myCompanyID)

	var reqData objects.FormUsers
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
	formUserID := 0

	var postData objects.InputFormUserOrganizations
	postData.FormID = reqData.FormID
	postData.UserID = reqData.UserID
	postData.Type = "admin"
	postData.OrganizationID = myCompanyID

	updateCF, err := ctr.formMod.ConnectFormUserOrg(postData)
	if err != nil {
		fmt.Println(err)

		if errors.As(err, &ctr.pgErr) {
			fmt.Println("ctr.pgErr.Code--", ctr.pgErr.Code)
			if ctr.pgErr.Code == "23505" { //code duplicate

				var whr3 tables.FormUsers
				whr3.FormID = reqData.FormID
				whr3.UserID = reqData.UserID
				whr3.Type = "admin"
				getFormUser, err := ctr.formMod.GetFormUserRow(whr3)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}
				formUserID = getFormUser.ID

			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
		}
	}

	if updateCF.ID > 0 {
		formUserID = updateCF.ID
	}

	var whre tables.Permissions
	whre.HttpPath = "/form"
	getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
	if len(getPermiss) > 0 {

		// send mail
		// var whreUser tables.Users
		// whreUser.ID = reqData.UserID
		// getUser, err := ctr.userMod.GetUserRow(whreUser)
		// to := []string{carttoorders[i].CustomersEmail}
		// cc := []string{}
		// subject := "Snap-in | Akun anda di inveite sebagai Admin"
		// message := `Halo [` + customerName + `], terima kasih telah berbelanja di [` + customerName + `], untuk [` + ItemsCount + `]produk, sejumlah Rp. ` + grandtotal + `. Unduh aplikasi Semarket disini [` + url + `]`

		// helpers.sendMail(to, cc, subject, message)

		fmt.Println(len(getPermiss))
		for i := 0; i < len(getPermiss); i++ {
			var postData tables.FormUserPermission
			postData.FormUserID = formUserID
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
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success sign data to form",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed sign data to form",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FormUserAdminFrmDisconnect(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	authorID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println(authorID)

	var reqData objects.FormUsers
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

	var whre tables.FormUsers
	whre.UserID = reqData.UserID
	whre.FormID = reqData.FormID
	whre.Type = "admin"
	getFormUser, _ := ctr.formMod.GetFormUserRow(whre)

	//delete form user permission
	remove, _ := ctr.permissMod.DeleteFormUserPermission(getFormUser.ID, "")

	if remove {
		_, err := ctr.formMod.DeleteFormUserOrg(reqData.UserID, reqData.FormID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success deleted data in form user",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed deleted data in form user",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) FormUserAdminCheck(c *gin.Context) {

	// claims := jwt.ExtractClaims(c)
	// userID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.FormUserPermissions
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

	if len(reqData.Data) > 0 {
		for i := 0; i < len(reqData.Data); i++ {
			var whreP tables.FormUserPermission
			whreP.Status = reqData.Data[i].Status
			_, err := ctr.permissMod.UpdateFormUserPermission(reqData.Data[i].ID, whreP)
			if err != nil {
				fmt.Println("ERR: UpdateFormUserPermission in loop ", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success update user permission",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed update user permission, please try again",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormAttendanceRequired(c *gin.Context) {

	fID, _ := strconv.Atoi(c.Param("formid"))

	var reqData objects.FormAttendanceRequired
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

	if fID > 0 {

		formStatus := "not required"
		if reqData.IsAttendanceRequired {
			formStatus = "ATTENDANCE REQUIRED "
		}
		var data tables.Forms
		data.IsAttendanceRequired = reqData.IsAttendanceRequired
		saveUpdate, err := ctr.formMod.UpdateForm(fID, data)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if saveUpdate {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Form update is successful",
				"data":    "Form is " + formStatus,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed form update",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "form ID is required",
			"data":    nil,
		})
		return
	}

}

func (ctr *formController) AdminFormList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	formID, _ := strconv.Atoi(c.Param("formid"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		form, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if form.ID > 0 {
			updatedAt := form.UpdatedAt
			lastUpdate := updatedAt.Format("2006-02-01 15:04")

			// get total responden
			var whereFU tables.JoinFormUsers
			whereFU.FormID = form.ID
			whereFU.Type = "respondent"

			getResponden, err := ctr.formMod.GetFormUserRows(whereFU, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// get total admin
			var whrFu tables.JoinFormUsers
			whrFu.FormID = form.ID
			whrFu.Type = "admin"

			getAdmins, err := ctr.formMod.GetFormUserRows(whrFu, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// total respon
			var whereInForm tables.InputForms
			getDataRespons, err := ctr.inputForm.GetInputFormRows(form.ID, whereInForm, "", objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			var totalUpdated int
			var whreStr = "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(formID, whereInForm, whreStr)
			getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedDataWithDate(formID, whreStr)
			for i := 0; i < len(getDataSubmissionUpdated); i++ {
				totalUpdated = getDataSubmissionUpdated[i].UpdatedCount
			}

			var res objects.Forms
			res.ID = form.ID
			res.Name = form.Name
			res.Description = form.Description
			res.FormStatusID = form.FormStatusID
			res.Notes = form.Notes
			res.ProfilePic = form.ProfilePic
			res.PeriodStartDate = form.PeriodStartDate
			res.PeriodEndDate = form.PeriodEndDate
			res.IsAttendanceRequired = form.IsAttendanceRequired
			res.SubmissionTarget = form.SubmissionTargetUser * len(getResponden)
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = form.FormStatusID
			res.TotalRespon = len(getDataRespons)
			res.TotalAdmin = len(getAdmins) + 1
			res.TotalResponden = len(getResponden)
			res.ShareUrl = form.ShareUrl
			res.TotalDeletedData = len(getDataSubmissionDeleted)
			res.TotalUpdatedData = totalUpdated

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			// SUPER ADMIN here
			var whrComp tables.Organizations
			whrComp.CreatedBy = userID
			whrComp.IsDefault = true
			getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)

			var buffer bytes.Buffer
			var fields tables.FormOrganizationsJoin
			// fields.FormStatusID = 1
			fields.OrganizationID = getComp.ID

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id = 1")
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

			// get all data
			getFormsAll, err := ctr.formMod.GetFormOwnerRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

		} else {

			var buffer bytes.Buffer
			var fields tables.Forms

			whereString := ""
			if searchKeyWord != "" {
				buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
			}

			buffer.WriteString(" forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")")
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
			//get all data
			getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll
		}

		var res []objects.Forms
		if len(result) > 0 {

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

				//total active responden
				whreStrToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrToday)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStrToday, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total performance
				// inputFormData, err := ctr.userMod.InputFormUserPaging(result[i].ID, "", objects.Paging{})
				// if err != nil {
				// 	c.JSON(http.StatusInternalServerError, gin.H{
				// 		"error": err,
				// 	})
				// 	return
				// }

				var totalPerform int
				var totalPerformFloat float64
				if result[i].SubmissionTargetUser > 0 && len(getResponden) > 0 {
					totalPerformFloat = ((float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)) * 100) / float64(len(getResponden))
					totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))

					if totalPerform >= 100 {
						totalPerform = 100
						totalPerformFloat = 100
					}

				} else {
					totalPerform = 0
				}
				// "#CB3939" //red
				// "#D6C31A" //yellow
				// "#398037;" // green

				// 0-40% : merah
				// 41-60% : kuning
				// 61-100% : hijau

				// performance color
				performColor := "#CB3939"
				if totalPerform >= 41 && totalPerform <= 60 {
					performColor = "#D6C31A"
				} else if totalPerform >= 61 {
					performColor = "#398037"
				}
				strTotalPerform := strconv.FormatFloat(totalPerformFloat, 'f', 1, 64)

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespondenActive = len(getActiveRespondens)
				each.TotalRespon = len(getDataRespons)
				each.TotalPerformance = totalPerform
				each.TotalPerformanceFloat, _ = strconv.ParseFloat(strTotalPerform, 1)
				each.IsAttendanceRequired = result[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = result[i].SubmissionTargetUser
				each.PeriodeRange = 0
				each.Author = result[i].Author
				each.PerformanceColor = performColor
				each.ShareUrl = result[i].ShareUrl

				res = append(res, each)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
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

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
			"paging":  paging,
		})
		return

	}
}

func (ctr *formController) AdminFormListNew(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	formID, _ := strconv.Atoi(c.Param("formid"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {
		var fields tables.Forms
		fields.ID = formID

		form, err := ctr.formMod.GetFormRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getPeriodeRange, err := ctr.formMod.GetFormPeriodeRangeRow(fields)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if form.ID > 0 {
			updatedAt := form.UpdatedAt
			lastUpdate := updatedAt.Format("2006-02-01 15:04")

			// get total responden
			var whereFU tables.JoinFormUsers
			whereFU.FormID = form.ID
			whereFU.Type = "respondent"

			getResponden, err := ctr.formMod.GetFormUserRows(whereFU, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// get total admin
			var whrFu tables.JoinFormUsers
			whrFu.FormID = form.ID
			whrFu.Type = "admin"

			getAdmins, err := ctr.formMod.GetFormUserRows(whrFu, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// total respon
			var whereInForm tables.InputForms
			getDataRespons, err := ctr.inputForm.GetInputFormRows(form.ID, whereInForm, "", objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			var totalUpdated int
			var whreStr = "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			getDataSubmissionDeleted, _ := ctr.inputForm.GetDeletedData(formID, whereInForm, whreStr)
			getDataSubmissionUpdated, _ := ctr.inputForm.GetUpdatedDataWithDate(formID, whreStr)
			for i := 0; i < len(getDataSubmissionUpdated); i++ {
				totalUpdated = getDataSubmissionUpdated[i].UpdatedCount
			}

			var res objects.Forms
			res.ID = form.ID
			res.Name = form.Name
			res.Description = form.Description
			res.FormStatusID = form.FormStatusID
			res.Notes = form.Notes
			res.ProfilePic = form.ProfilePic
			res.PeriodStartDate = form.PeriodStartDate
			res.PeriodEndDate = form.PeriodEndDate
			res.IsAttendanceRequired = form.IsAttendanceRequired
			res.SubmissionTarget = form.SubmissionTargetUser * len(getResponden)
			res.UpdatedByName = ""
			res.LastUpdate = lastUpdate
			res.PeriodeRange = getPeriodeRange.PeriodRange
			res.FormStatusID = form.FormStatusID
			res.TotalRespon = len(getDataRespons)
			res.TotalAdmin = len(getAdmins) + 1
			res.TotalResponden = len(getResponden)
			res.ShareUrl = form.ShareUrl
			res.TotalDeletedData = len(getDataSubmissionDeleted)
			res.TotalUpdatedData = totalUpdated

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

		var result []tables.FormAll
		var resultAll []tables.FormAll
		if roleID == 1 { // 1 is owner

			var buffer bytes.Buffer
			var fields tables.FormOrganizationsJoin
			// SUPER ADMIN here
			fields.OrganizationID = organizationID

			whereName := ``
			whereString := "AND forms.form_status_id in (1)"
			if searchKeyWord != "" {
				buffer.WriteString(" name ilike '%" + searchKeyWord + "%'")
				// buffer.WriteString(" forms.form_status_id = 1")
				whereString = buffer.String()
			}

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeSuperAdminApps(fields, whereName, whereString, paging)
			// getForms, err := ctr.formMod.GetFormOwnerRows(fields, whereString, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			result = getForms

			// get all data
			getFormsAll, err := ctr.formMod.GetFormMergeSuperAdminApps(fields, whereName, whereString, paging)
			// getFormsAll, err := ctr.formMod.GetFormOwnerRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll

		} else {

			var searchForm bytes.Buffer
			var buffer bytes.Buffer
			var whre bytes.Buffer

			var whereName string
			var whereString string
			var whereStr string

			var fields tables.FormOrganizationsJoin
			fields.OrganizationID = organizationID

			if searchKeyWord != "" {
				searchForm.WriteString(" where name ilike '%" + searchKeyWord + "%'")
				whereName = searchForm.String()
			}

			buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")")
			whereString = buffer.String()

			whre.WriteString(" f.form_status_id not in (3)")
			whereStr = whre.String()

			var paging objects.Paging
			paging.Page = page
			paging.Limit = limit
			paging.SortBy = sortBy
			paging.Sort = sort

			getForms, err := ctr.formMod.GetFormMergeAdminApps(fields, whereName, whereString, whereStr, userID, paging)
			// getForms, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			result = getForms
			//get all data
			getFormsAll, err := ctr.formMod.GetFormMergeAdminApps(fields, whereName, whereString, whereStr, userID, paging)
			// getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			resultAll = getFormsAll
		}

		var res []objects.Forms
		if len(result) > 0 {

			for i := 0; i < len(result); i++ {

				// get total responden form internal
				var whereFU tables.JoinFormUsers
				whereFU.FormID = result[i].ID
				whereFU.Type = "respondent"
				whreStr := ""

				getRespondenInternal, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// get total responden form external
				var where tables.JoinFormUsers
				where.FormID = result[i].ID
				where.Type = "respondent"
				whreString := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

				getRespondenExternal, err := ctr.formMod.GetFormUserRows(where, whreString)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total active responden form internal
				whreStrToday := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondensInternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrToday)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total active responden form external
				whreStrDate := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND ifo.organization_id = " + strconv.Itoa(organizationID) + ""
				// var whreActive tables.InputForms
				getActiveRespondensExternal, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrDate)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon Internal
				var whereInForm tables.InputForms
				getDataResponFormInternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon External
				whreStr = "ifo.organization_id = " + strconv.Itoa(organizationID) + ""
				getDataResponFormExternal, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStr, objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total performance
				// inputFormData, err := ctr.userMod.InputFormUserPaging(result[i].ID, "", objects.Paging{})
				// if err != nil {
				// 	c.JSON(http.StatusInternalServerError, gin.H{
				// 		"error": err,
				// 	})
				// 	return
				// }

				// var totalPerform int
				// var totalPerformFloat float64
				// if result[i].SubmissionTargetUser > 0 && len(getResponden) > 0 {
				// 	totalPerformFloat = ((float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)) * 100) / float64(len(getResponden))
				// 	totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))

				// 	if totalPerform >= 100 {
				// 		totalPerform = 100
				// 		totalPerformFloat = 100
				// 	}

				// } else {
				// 	totalPerform = 0
				// }

				//total performance Internal
				var totalPerformInternal int
				var totalPerformFloatInternal float64
				if result[i].SubmissionTargetUser > 0 && len(getRespondenInternal) > 0 {
					totalPerformFloatInternal = ((float64(len(getDataResponFormInternal)) / float64(result[i].SubmissionTargetUser)) * 100) / float64(len(getRespondenInternal))
					totalPerformInternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatInternal, 'f', 0, 64))

					if totalPerformInternal >= 100 {
						totalPerformInternal = 100
						totalPerformFloatInternal = 100
					}
				} else {
					totalPerformInternal = 0
				}

				var totalPerformExternal int
				var totalPerformFloatExternal float64
				if result[i].SubmissionTargetUser > 0 && len(getRespondenExternal) > 0 {
					totalPerformFloatExternal = ((float64(len(getDataResponFormExternal)) / float64(result[i].SubmissionTargetUser)) * 100) / float64(len(getRespondenExternal))
					totalPerformExternal, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloatExternal, 'f', 0, 64))

					if totalPerformExternal >= 100 {
						totalPerformExternal = 100
						totalPerformFloatExternal = 100
					}
				} else {
					totalPerformExternal = 0
				}
				// "#CB3939" //red
				// "#D6C31A" //yellow
				// "#398037;" // green

				// 0-40% : merah
				// 41-60% : kuning
				// 61-100% : hijau

				// performance color
				performColorInternal := "#CB3939"
				if totalPerformInternal >= 41 && totalPerformInternal <= 60 {
					performColorInternal = "#D6C31A"
				} else if totalPerformInternal >= 61 {
					performColorInternal = "#398037"
				}
				strtotalPerformInternal := strconv.FormatFloat(totalPerformFloatInternal, 'f', 1, 64)

				performColorExternal := "#CB3939"
				if totalPerformExternal >= 41 && totalPerformExternal <= 60 {
					performColorExternal = "#D6C31A"
				} else if totalPerformExternal >= 61 {
					performColorExternal = "#398037"
				}
				strtotalPerformExternal := strconv.FormatFloat(totalPerformFloatExternal, 'f', 1, 64)

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.IsAttendanceRequired = result[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = result[i].SubmissionTargetUser
				each.PeriodeRange = 0
				each.Author = result[i].Author
				each.ShareUrl = result[i].ShareUrl
				each.Type = result[i].Type
				if result[i].Type == "internal" {
					each.TotalResponden = len(getDataResponFormInternal)
					each.TotalRespondenActive = len(getActiveRespondensInternal)
					each.TotalRespon = len(getDataResponFormInternal)
					each.TotalPerformance = totalPerformInternal
					each.TotalPerformanceFloat, _ = strconv.ParseFloat(strtotalPerformInternal, 1)
					each.PerformanceColor = performColorInternal
				} else if result[i].Type == "external" {
					each.TotalResponden = len(getDataResponFormExternal)
					each.TotalRespondenActive = len(getActiveRespondensExternal)
					each.TotalRespon = len(getDataResponFormExternal)
					each.TotalPerformance = totalPerformExternal
					each.TotalPerformanceFloat, _ = strconv.ParseFloat(strtotalPerformExternal, 1)
					each.PerformanceColor = performColorExternal
				}

				res = append(res, each)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
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

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
			"paging":  paging,
		})
		return

	}
}

func (ctr *formController) FormDuplicate(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  false,
		"message": "Failed created form, please try again",
		"data":    nil,
	})
	return
}

func (ctr *formController) FormDuplicate2(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	var reqData objects.FormDuplicate
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

	if reqData.FormID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}

	var companyID int
	if roleID == 1 {
		var whre tables.Organizations
		whre.CreatedBy = userID
		whre.IsDefault = true
		getComp, _ := ctr.compMod.GetCompaniesRow(whre)
		companyID = getComp.ID
	} else {
		var whrUserComp objects.UserOrganizations
		whrUserComp.UserID = userID
		getUComp, _ := ctr.compMod.GetUserCompaniesRow(whrUserComp, "")
		companyID = getUComp.OrganizationID
	}

	//get Data form
	var whrFrm tables.Forms
	whrFrm.ID = reqData.FormID
	getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

	//check name
	var whrFrm2 tables.Forms
	whreStr := "fo.organization_id = " + strconv.Itoa(companyID) + " AND forms.name ilike '%" + getFormData.Name + "%'"
	getFormsByname, _ := ctr.formMod.GetFormWhreRows(whrFrm2, whreStr)
	countFile := len(getFormsByname) + 1

	var postData tables.Forms
	postData.Name = getFormData.Name + "(" + strconv.Itoa(countFile) + ")"
	postData.Description = getFormData.Description
	postData.Notes = getFormData.Notes
	postData.PeriodStartDate = getFormData.PeriodStartDate
	postData.ProfilePic = getFormData.ProfilePic
	postData.EncryptCode = helpers.EncodeToString(6)
	postData.CreatedBy = userID
	postData.SubmissionTargetUser = getFormData.SubmissionTargetUser
	postData.IsAttendanceRequired = getFormData.IsAttendanceRequired
	postData.ShareUrl = getFormData.ShareUrl
	postData.CreatedBy = userID

	if getFormData.PeriodEndDate != "" {
		postData.PeriodEndDate = getFormData.PeriodEndDate
	}

	newForm, err := ctr.formMod.InsertForm(postData)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// ctr.helper.AddLogBook(userID, 16, res.ID)

	//insert form organization
	var postData2 tables.FormOrganizations
	postData2.FormID = newForm.ID
	postData2.OrganizationID = companyID
	fOrg, err := ctr.formMod.InsertFormOrganization(postData2)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	//get form field existing
	var whrField tables.FormFields
	whrField.FormID = reqData.FormID
	// whrField.ParentID = 0
	getExistingForm, _ := ctr.formFieldMod.GetFormFieldNotParentRows(whrField, "")

	if len(getExistingForm) > 0 {
		for k := 0; k < len(getExistingForm); k++ {

			//cek field rule yg existing (yg mw di duplicate : utk dicari sort order nya)
			var newParentField22 tables.SelectFormFieldConditionRules
			if getExistingForm[k].ParentID > 0 {
				var whrPrntField tables.FormFields
				whrPrntField.ID = getExistingForm[k].ParentID
				getExistingPrntField, _ := ctr.formFieldMod.GetFormFieldRow(whrPrntField)

				var whrNewPrntField22 tables.FormFields
				whrNewPrntField22.FormID = newForm.ID
				whrNewPrntField22.SortOrder = getExistingPrntField.SortOrder
				newParentField22, _ = ctr.formFieldMod.GetFormFieldRow(whrNewPrntField22)
			}

			// insert field
			var dataField tables.FormFields
			dataField.ParentID = newParentField22.ID
			dataField.FormID = newForm.ID
			dataField.FieldTypeID = getExistingForm[k].FieldTypeID
			dataField.Label = getExistingForm[k].Label
			dataField.Description = getExistingForm[k].Description
			dataField.Option = getExistingForm[k].Option
			dataField.ConditionType = getExistingForm[k].ConditionType
			dataField.UpperlowerCaseType = getExistingForm[k].UpperlowerCaseType
			dataField.IsMultiple = getExistingForm[k].IsMultiple
			dataField.IsCondition = getExistingForm[k].IsCondition
			dataField.IsRequired = getExistingForm[k].IsRequired
			dataField.SortOrder = getExistingForm[k].SortOrder
			dataField.IsSection = getExistingForm[k].IsSection
			dataField.SectionColor = getExistingForm[k].SectionColor
			dataField.TagLocColor = getExistingForm[k].TagLocColor
			dataField.TagLocIcon = getExistingForm[k].TagLocIcon

			newField, err := ctr.formFieldMod.InsertFormField(dataField)
			if err != nil {
				fmt.Println("InsertFormField", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			// insert rule
			if newField.ID > 0 && getExistingForm[k].ConditionRuleID > 0 {

				var dataRule tables.FormFieldConditionRules
				dataRule.FormFieldID = newField.ID
				dataRule.ConditionRuleID = getExistingForm[k].ConditionRuleID
				dataRule.Value1 = getExistingForm[k].Value1
				dataRule.Value2 = getExistingForm[k].Value2
				dataRule.ErrMsg = getExistingForm[k].ErrMsg
				dataRule.TabMaxOnePerLine = getExistingForm[k].TabMaxOnePerLine
				dataRule.TabEachLineRequire = getExistingForm[k].TabEachLineRequire

				_, err = ctr.ruleMod.InsertFormFieldRule(dataRule)
				if err != nil {
					fmt.Println("InsertFormFieldRule", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
			} else if newField.ID > 1 && getExistingForm[k].ConditionRuleID == 0 {

				_, err = ctr.ruleMod.DeleteFormFieldRulePrimary(newField.ID)
				if err != nil {
					fmt.Println("Err :DeleteFormFieldRule", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
			}

			if newField.ID > 0 && getExistingForm[k].Image != "" {

				var postData tables.FormFieldPics
				postData.FormFieldID = newField.ID
				postData.Pic = getExistingForm[k].Image
				_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
				if err != nil {
					fmt.Println("InsertFormFieldPic", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
			} else if newField.ID > 0 && getExistingForm[k].Image == "" {

				_, err = ctr.formFieldMod.DeleteFormFieldPic(newField.ID)
				if err != nil {
					fmt.Println("Err :DeleteFormFieldPic", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
			}

			// insert condition rule
			if getExistingForm[k].IsCondition == true {
				var whreRules tables.FormFieldConditionRules
				whreRules.FormFieldID = getExistingForm[k].ID
				getConditionRules, err := ctr.ruleMod.GetFormFieldRuleRows(whreRules, "")

				if len(getConditionRules) > 0 {

					// clear old data tidak berlaku di duplicate
					// _, err = ctr.ruleMod.DeleteFormFieldRule(getExistingForm[k].ID)
					// if err != nil {
					// 	fmt.Println("Err :DeleteFormFieldRule", err)
					// 	c.JSON(http.StatusInternalServerError, gin.H{
					// 		"error": err,
					// 	})
					// 	return
					// }

					// fmt.Println("#getConditionRules ------------------------------------------------------", getExistingForm[k].ID, "---", newField.ID, len(getConditionRules))

					// var save tables.FormFieldConditionRules
					var cr tables.FormFieldConditionRules
					cr.FormFieldID = newField.ID
					cr.ConditionAllRight = getExistingForm[k].ConditionAllRight
					for j := 0; j < len(getConditionRules); j++ {

						//cek field rule yg existing (yg mw di duplicate : utk mencari sort order nya)
						if getConditionRules[j].ConditionParentFieldID > 0 {
							var whrParentField tables.FormFields
							whrParentField.ID = getConditionRules[j].ConditionParentFieldID
							getExistingParentField, _ := ctr.formFieldMod.GetFormFieldRow(whrParentField)

							// fmt.Println("existingParentField >>>>>>>", getExistingParentField.SortOrder)
							// fmt.Println("getExistingParentField.ID >>>>>>>", getExistingParentField.ID, getExistingParentField.Label)

							var whrNewParentField tables.FormFields
							whrNewParentField.FormID = newForm.ID
							whrNewParentField.ParentID = 0
							whrNewParentField.SortOrder = getExistingParentField.SortOrder
							newParentField, _ := ctr.formFieldMod.GetFormFieldRow(whrNewParentField)

							fmt.Println("new GetFormFieldRow =================", newForm.ID, "--new field:", newField.ID, "--", getExistingParentField.SortOrder, "result ::", newParentField.ID)
							//insert new rules
							cr.ConditionParentFieldID = newParentField.ID
							cr.ConditionRuleID = getConditionRules[j].ConditionRuleID
							cr.Value1 = getConditionRules[j].Value1
							cr.Value2 = getConditionRules[j].Value2

							_, err = ctr.ruleMod.InsertFormFieldRule(cr)
							if err != nil {
								fmt.Println(" Err : InsertFormFieldRule", err)
								c.JSON(http.StatusInternalServerError, gin.H{
									"error": err,
								})
								return
							}
						}
					}
					getConditionRules = []tables.FormFieldConditionRules{}
				}
			}
		}

		if fOrg.ID > 0 {
			var obj objects.FormResponse
			obj.ID = newForm.ID

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Success created form",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed created form, please try again",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed duplicate form, field of form is empty",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormDuplicate3(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	// roleID, _ := strconv.Atoi(claims["role_id"].(string))
	companyID, _ := strconv.Atoi(claims["organization_id"].(string))

	var reqData objects.FormDuplicate
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

	fmt.Println("companyID ::", companyID)

	if reqData.FormID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}

	generateDuplicate, formDataResult, err := MagicFormDuplicate(ctr, companyID, userID, reqData.FormID, reqData.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if generateDuplicate {
		var obj objects.FormData
		obj.FormID = formDataResult.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created form",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Success created form",
		})
		return
	}
}

func MagicFormDuplicate(ctr *formController, companyID, userID int, formID int, projectID int) (bool, tables.Forms, error) {
	//get Data form
	var whrFrm tables.Forms
	whrFrm.ID = formID
	getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

	//check name
	var whrFrm2 tables.Forms
	whreStr := "fo.organization_id = " + strconv.Itoa(companyID) + " AND forms.name ilike '%" + getFormData.Name + "%'"
	getFormsByname, _ := ctr.formMod.GetFormWhreRows(whrFrm2, whreStr)
	countFile := len(getFormsByname) + 1

	var postData tables.Forms
	postData.Name = getFormData.Name + "(" + strconv.Itoa(countFile) + ")"
	postData.Description = getFormData.Description
	postData.Notes = getFormData.Notes
	postData.PeriodStartDate = getFormData.PeriodStartDate
	postData.ProfilePic = getFormData.ProfilePic
	postData.EncryptCode = helpers.EncodeToString(6)
	postData.CreatedBy = userID
	postData.SubmissionTargetUser = getFormData.SubmissionTargetUser
	postData.IsAttendanceRequired = getFormData.IsAttendanceRequired
	postData.ShareUrl = getFormData.ShareUrl
	postData.CreatedBy = userID

	if getFormData.PeriodEndDate != "" {
		postData.PeriodEndDate = getFormData.PeriodEndDate
	}

	if getFormData.IsAttendanceRequired == true && getFormData.AttendanceOverdateAt.IsZero() == false {
		postData.AttendanceOverdateAt = helpers.DateNow()
	}

	if getFormData.IsAttendanceRequired == true {
		postData.IsAttendanceRadius = getFormData.IsAttendanceRadius
	}

	newForm, err := ctr.formMod.InsertForm(postData)
	if err != nil {
		fmt.Println("InsertForm", err)

		return false, tables.Forms{}, err
	}

	// duplicate location radius
	getFormAttLocations, err := ctr.attendanceMod.GetFormAttendanceLocationRows(objects.ObjectFormAttendanceLocations{FormID: getFormData.ID})
	if err != nil {
		fmt.Println("InsertForm", err)

		return false, tables.Forms{}, err
	}
	if len(getFormAttLocations) >= 1 {
		for i := 0; i < len(getFormAttLocations); i++ {
			var postAttData objects.ObjectFormAttendanceLocations
			postAttData.FormID = newForm.ID
			postAttData.Name = getFormAttLocations[i].Name
			postAttData.Location = getFormAttLocations[i].Location
			postAttData.Latitude = getFormAttLocations[i].Latitude
			postAttData.Longitude = getFormAttLocations[i].Longitude
			postAttData.IsCheckIn = getFormAttLocations[i].IsCheckIn
			postAttData.IsCheckOut = getFormAttLocations[i].IsCheckOut
			postAttData.Radius = getFormAttLocations[i].Radius

			ctr.formMod.InsertFormAttendanceLocation(postAttData)
		}
	}

	// check form in project
	if projectID != 0 {
		// check form in project
		var whrProjectFrm tables.ProjectForms
		whrProjectFrm.FormID = getFormData.ID
		getPorjectForm, _ := ctr.projectMod.GetProjectForms(whrProjectFrm, "")
		if len(getPorjectForm) > 0 {

			var postPF tables.ProjectForms
			postPF.FormID = newForm.ID
			postPF.ProjectID = getPorjectForm[0].ProjectID
			ctr.projectMod.InsertProjectForm(postPF)
		}
	}

	// store log form activity
	// ctr.helper.AddLogBook(userID, 16, res.ID)

	//insert form organization
	var postData2 tables.FormOrganizations
	postData2.FormID = newForm.ID
	postData2.OrganizationID = companyID
	fOrg, err := ctr.formMod.InsertFormOrganization(postData2)
	if err != nil {
		fmt.Println("InsertForm", err)
		return false, tables.Forms{}, err
	}

	//get form field existing
	getExistingForm, _ := ctr.formFieldMod.GetFormFieldNotParentRows(tables.FormFields{FormID: formID}, "")

	if len(getExistingForm) > 0 {
		for k := 0; k < len(getExistingForm); k++ {

			//cek field rule yg existing (yg mw di duplicate : utk dicari sort order nya)
			var newParentField22 tables.SelectFormFieldConditionRules
			if getExistingForm[k].ParentID > 0 {
				var whrPrntField tables.FormFields
				whrPrntField.ID = getExistingForm[k].ParentID
				getExistingPrntField, _ := ctr.formFieldMod.GetFormFieldRow(whrPrntField)

				var whrNewPrntField22 tables.FormFields
				whrNewPrntField22.FormID = newForm.ID
				whrNewPrntField22.SortOrder = getExistingPrntField.SortOrder
				newParentField22, _ = ctr.formFieldMod.GetFormFieldRow(whrNewPrntField22)
			}

			// insert field
			var dataField tables.FormFields
			dataField.ParentID = newParentField22.ID
			dataField.FormID = newForm.ID
			dataField.FieldTypeID = getExistingForm[k].FieldTypeID
			dataField.Label = getExistingForm[k].Label
			dataField.Description = getExistingForm[k].Description
			dataField.Option = getExistingForm[k].Option
			dataField.ConditionType = getExistingForm[k].ConditionType
			dataField.UpperlowerCaseType = getExistingForm[k].UpperlowerCaseType
			dataField.IsMultiple = getExistingForm[k].IsMultiple
			dataField.IsCondition = getExistingForm[k].IsCondition
			dataField.IsRequired = getExistingForm[k].IsRequired
			dataField.SortOrder = getExistingForm[k].SortOrder
			dataField.IsSection = getExistingForm[k].IsSection
			dataField.SectionColor = getExistingForm[k].SectionColor
			dataField.TagLocColor = getExistingForm[k].TagLocColor
			dataField.TagLocIcon = getExistingForm[k].TagLocIcon

			newField, err := ctr.formFieldMod.InsertFormField(dataField)
			if err != nil {
				fmt.Println("InsertFormField", err)
				return false, tables.Forms{}, err
			}

			// insert rule
			if newField.ID > 0 && getExistingForm[k].ConditionRuleID > 0 {

				var dataRule tables.FormFieldConditionRules
				dataRule.FormFieldID = newField.ID
				dataRule.ConditionRuleID = getExistingForm[k].ConditionRuleID
				dataRule.Value1 = getExistingForm[k].Value1
				dataRule.Value2 = getExistingForm[k].Value2
				dataRule.ErrMsg = getExistingForm[k].ErrMsg
				dataRule.TabMaxOnePerLine = getExistingForm[k].TabMaxOnePerLine
				dataRule.TabEachLineRequire = getExistingForm[k].TabEachLineRequire

				_, err = ctr.ruleMod.InsertFormFieldRule(dataRule)
				if err != nil {
					fmt.Println("InsertFormFieldRule", err)
					return false, tables.Forms{}, err
				}
			} else if newField.ID > 1 && getExistingForm[k].ConditionRuleID == 0 {

				_, err = ctr.ruleMod.DeleteFormFieldRulePrimary(newField.ID)
				if err != nil {
					fmt.Println("Err :DeleteFormFieldRule", err)
					return false, tables.Forms{}, err
				}
			}

			if newField.ID > 0 && getExistingForm[k].Image != "" {

				var postData tables.FormFieldPics
				postData.FormFieldID = newField.ID
				postData.Pic = getExistingForm[k].Image
				_, err := ctr.formFieldMod.InsertFormFieldPic(postData)
				if err != nil {
					fmt.Println("InsertFormFieldPic", err)
					return false, tables.Forms{}, err
				}
			} else if newField.ID > 0 && getExistingForm[k].Image == "" {

				_, err = ctr.formFieldMod.DeleteFormFieldPic(newField.ID)
				if err != nil {
					fmt.Println("Err :DeleteFormFieldPic", err)
					return false, tables.Forms{}, err
				}
			}

			// insert condition rule
			if getExistingForm[k].IsCondition == true {
				var whreRules tables.FormFieldConditionRules
				whreRules.FormFieldID = getExistingForm[k].ID
				getConditionRules, err := ctr.ruleMod.GetFormFieldRuleRows(whreRules, "")

				if len(getConditionRules) > 0 {

					// clear old data tidak berlaku di duplicate
					// _, err = ctr.ruleMod.DeleteFormFieldRule(getExistingForm[k].ID)
					// if err != nil {
					// 	fmt.Println("Err :DeleteFormFieldRule", err)
					// 	c.JSON(http.StatusInternalServerError, gin.H{
					// 		"error": err,
					// 	})
					// 	return
					// }

					// fmt.Println("#getConditionRules ------------------------------------------------------", getExistingForm[k].ID, "---", newField.ID, len(getConditionRules))

					// var save tables.FormFieldConditionRules
					var cr tables.FormFieldConditionRules
					cr.FormFieldID = newField.ID
					cr.ConditionAllRight = getExistingForm[k].ConditionAllRight
					for j := 0; j < len(getConditionRules); j++ {

						//cek field rule yg existing (yg mw di duplicate : utk mencari sort order nya)
						if getConditionRules[j].ConditionParentFieldID > 0 {
							var whrParentField tables.FormFields
							whrParentField.ID = getConditionRules[j].ConditionParentFieldID
							getExistingParentField, _ := ctr.formFieldMod.GetFormFieldRow(whrParentField)

							// fmt.Println("existingParentField >>>>>>>", getExistingParentField.SortOrder)
							// fmt.Println("getExistingParentField.ID >>>>>>>", getExistingParentField.ID, getExistingParentField.Label)

							var whrNewParentField tables.FormFields
							whrNewParentField.FormID = newForm.ID
							whrNewParentField.ParentID = 0
							whrNewParentField.SortOrder = getExistingParentField.SortOrder
							newParentField, _ := ctr.formFieldMod.GetFormFieldRow(whrNewParentField)

							//insert new rules
							cr.ConditionParentFieldID = newParentField.ID
							cr.ConditionRuleID = getConditionRules[j].ConditionRuleID
							cr.Value1 = getConditionRules[j].Value1
							cr.Value2 = getConditionRules[j].Value2

							_, err = ctr.ruleMod.InsertFormFieldRule(cr)
							if err != nil {
								fmt.Println(" Err : InsertFormFieldRule", err)
								return false, tables.Forms{}, err
							}
						}
					}
					getConditionRules = []tables.FormFieldConditionRules{}
				}
			}
		}

		if fOrg.ID > 0 {
			var formDataRes tables.Forms
			formDataRes.ID = newForm.ID

			return true, formDataRes, nil
		} else {
			err := errors.New("Failed generate Company of form")
			return false, tables.Forms{}, err
		}
	} else {
		err := errors.New("Failed duplicate of form ID")
		return false, tables.Forms{}, err
	}

}

func (ctr *formController) FormTemplate(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	var reqData objects.FormTemplate
	if err := c.BindJSON(&reqData); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var companyID int
	if roleID == 1 {
		var whre tables.Organizations
		whre.CreatedBy = userID
		whre.IsDefault = true
		getComp, _ := ctr.compMod.GetCompaniesRow(whre)
		companyID = getComp.ID
	} else {
		var whrUserComp objects.UserOrganizations
		whrUserComp.UserID = userID
		getUComp, _ := ctr.compMod.GetUserCompaniesRow(whrUserComp, "")
		companyID = getUComp.OrganizationID
	}

	// fmt.Println(companyID)
	// os.Exit(0)
	generateDuplicate, formDataResult, err := MagicFormDuplicate(ctr, companyID, userID, reqData.FormID, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	if generateDuplicate {
		var obj objects.FormData
		obj.FormID = formDataResult.ID

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success created form",
			"data":    obj,
		})
		return
	}

	return
}

func (ctr *formController) FormCompanyConnect(c *gin.Context) {

	var reqData objects.FormCompany
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

	//check the company ID has in your connecting
	var whrConnect tables.UserOrganizationInvites
	whrConnect.OrganizationID = reqData.OrganizationID
	checkCompanyConnect, err := ctr.formMod.GetUserCompaniesListInvitedRows(whrConnect, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	if len(checkCompanyConnect) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Company belum terhubung ke organisasi Anda, invite company terkait terlebih dahulu ",
		})
		return
	}

	var postData tables.FormOrganizationInvites
	postData.FormID = reqData.FormID
	postData.OrganizationID = reqData.OrganizationID
	_, err = ctr.formMod.InsertFormCompanyInvites(postData)
	if err != nil {

		if errors.As(err, &ctr.pgErr) {

			// fmt.Println("err", ctr.pgErr.Code)
			if ctr.pgErr.Code == "23505" {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Organization ID tersebut sudah masuk ke dalam Form tersebut",
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success connecting company inviting to form",
	})
	return
}

func (ctr *formController) FormCompanyDisconnect(c *gin.Context) {

	var reqData objects.FormCompany
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

	//check the company ID has in your connecting
	var whrConnect tables.UserOrganizationInvites
	whrConnect.OrganizationID = reqData.OrganizationID
	checkCompanyConnect, err := ctr.formMod.GetUserCompaniesListInvitedRows(whrConnect, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
	}
	if len(checkCompanyConnect) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Company belum terhubung ke organisasi Anda, invite company terkait terlebih dahulu ",
		})
		return
	}

	_, err = ctr.formMod.DeleteFormCompanyInvites(reqData.FormID, reqData.OrganizationID)
	if err != nil {
		fmt.Println("InsertForm", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success delete company inviting from form",
	})
	return
}

func (ctr *formController) FormToCompanyList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println(userID)

	formID, err := strconv.Atoi(c.Request.URL.Query().Get("form_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")

	if formID > 0 {

		//cekform id
		var frm tables.Forms
		frm.ID = formID
		getForm, err := ctr.formMod.GetFormRow(frm)
		if err != nil {
			fmt.Println("err: GetFormRow", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if getForm.ID == 0 {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is wrong",
				"error":   err,
			})
			return
		}

		var fields tables.JoinFormCompanies
		fields.FormID = formID

		whereString := ""
		if searchKeyWord != "" {
			whereString = " o.name ilike '%" + searchKeyWord + "%'  "
		}
		// whereType := "admin"
		// getFormOrganizationInvite, err := ctr.formMod.GetFormOrganizationInvite(formID)

		fuRows, err := ctr.formMod.GetFormCompanyInviteNew(fields, whereString)
		if err != nil {
			fmt.Println("err: GetFormCompanyInviteRows", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var results []objects.FormToCompanyList

		if len(fuRows) >= 1 {

			for i := 0; i < len(fuRows); i++ {
				var each objects.FormToCompanyList
				// capitalizedStr := strings.Title(fuRows[i].Type)

				each.ID = fuRows[i].ID
				each.FormID = fuRows[i].FormID
				each.OrganizationID = fuRows[i].OrganizationID
				each.OrganizationName = fuRows[i].OrganizationName
				each.OrganizationContactName = fuRows[i].OrganizationContactName
				each.OrganizationContactPhone = fuRows[i].OrganizationContactPhone
				each.OrganizationContactStatus = fuRows[i].Type
				each.OrganizationProfilePic = fuRows[i].OrganizationProfilePic
				each.IsQuotaSharing = fuRows[i].IsQuotaSharing

				results = append(results, each)
			}
		}

		allData, _ := ctr.formMod.GetFormCompanyInviteRows(fields, "")

		// static row owner
		var detail objects.DataRowsDetail
		detail.AllRows = len(allData)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Form user data is available",
			"data":    results,
			"detail":  detail,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormCompanyNotInList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println(userID)

	formID, err := strconv.Atoi(c.Request.URL.Query().Get("form_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")

	if formID > 0 {

		//cekform id
		var frm tables.Forms
		frm.ID = formID
		getForm, err := ctr.formMod.GetFormRow(frm)
		if err != nil {
			fmt.Println("err: GetFormRow", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if getForm.ID == 0 {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is wrong",
				"error":   err,
			})
			return
		}

		var fields tables.UserOrganizationInvites
		fields.UserID = userID

		whereString := " o.id not in (select foi.organization_id from frm.form_organization_invites foi where foi.organization_id=" + c.Request.URL.Query().Get("form_id") + ")"
		if searchKeyWord != "" {
			whereString = " AND o.name ilike '%" + searchKeyWord + "%'  "
		}

		fuRows, err := ctr.formMod.GetUserCompaniesListInvitedRows(fields, whereString)
		if err != nil {
			fmt.Println("err: GetFormCompanyRows", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		allData, _ := ctr.formMod.GetUserCompaniesListInvitedRows(fields, "")

		// static row owner
		var detail objects.DataRowsDetail
		detail.AllRows = len(allData)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Form user data is available",
			"data":    fuRows,
			"detail":  detail,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormCompanyUpdateQuota(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println(userID)

	frmCompanyID, err := strconv.Atoi(c.Param("frm_company_id"))
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	isQuotaSharing, _ := strconv.ParseBool(c.Request.URL.Query().Get("is_quota_sharing"))

	if frmCompanyID > 0 {
		fmt.Println("isQuotaSharing ::", isQuotaSharing)
		var whreUpdate tables.FormOrganizationInvites
		whreUpdate.ID = frmCompanyID
		whreUpdate.IsQuotaSharing = isQuotaSharing
		_, err := ctr.formMod.UpdateCompanyInviteForm(frmCompanyID, whreUpdate)
		if err != nil {
			fmt.Println("err: GetFormCompanyRows", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Company form success updated",
		})
		return

	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data company form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormMultyAccessList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	var result []tables.FormAll
	var resultAll []tables.FormAll
	var resultgetAll []tables.FormAll
	if roleID == 1 { // 1 is owner

		var whrComp tables.Organizations
		whrComp.CreatedBy = userID
		whrComp.IsDefault = true
		getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)

		// SUPER ADMIN here
		var buffer bytes.Buffer
		var fields tables.FormOrganizationsJoin
		fields.OrganizationID = getComp.ID

		whereGroupStr := ``
		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" AND forms.name ilike '%" + searchKeyWord + "%'  ")
			whereGroupStr = " AND t.name ilike '%" + searchKeyWord + "%'"
		}

		buffer.WriteString(" AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		result = getForms

		// get all data
		getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultAll = getFormsAll

		fmt.Println("result :::::", len(result), len(getFormsAll))

		whereAll := " AND forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
		getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll

	} else {
		var whrUserComp objects.UserOrganizations
		whrUserComp.UserID = userID
		getUComp, _ := ctr.compMod.GetUserCompaniesRow(whrUserComp, "")

		var buffer bytes.Buffer
		var whreGroup bytes.Buffer
		var whereString string
		var whereGroupStr string

		var fields tables.FormOrganizationsJoin
		fields.OrganizationID = getUComp.OrganizationID

		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%'   ")
			whreGroup.WriteString(" AND t.name ilike '%" + searchKeyWord + "%'")
			// whereGroupStr = " AND t.name ilike '%" + searchKeyWord + "%'"
		}

		whreGroup.WriteString(" AND (t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) OR t.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select f.id from frm.forms f where created_by = " + strconv.Itoa(userID) + ")))")
		whereGroupStr = whreGroup.String()

		buffer.WriteString(" AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		getForms, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		result = getForms

		//get all data
		getFormsAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereString, whereGroupStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultAll = getFormsAll

		whereAll := "AND forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
		getAll, err := ctr.formMod.GetFormUnionProjectRows(fields, whereAll, whereGroupStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll
	}

	if len(result) > 0 {

		var res []objects.Forms
		for i := 0; i < len(result); i++ {

			if result[i].ProjectID > 0 {

				//get total form in project/groups
				var whreFrm tables.ProjectForms
				var whrStr string
				whreFrm.ProjectID = result[i].ProjectID

				if roleID != 1 {
					whrStr = " f.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR f.id in (select f.id from frm.forms f where f.created_by= " + strconv.Itoa(userID) + ")"
				}
				getFormsInGroup, _ := ctr.projectMod.GetProjectForms(whreFrm, whrStr)

				var each objects.Forms
				each.ProjectID = result[i].ProjectID
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.TotalForms = len(getFormsInGroup)
				res = append(res, each)

			} else {

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

				//total active responden
				whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
				var whreActive tables.InputForms
				getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
				if err != nil {
					fmt.Println("err: GetFormUserRows", err)
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				// total respon
				var whereInForm tables.InputForms
				getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				//total performance
				var totalPerform int
				var totalPerformFloat float64
				if result[i].SubmissionTargetUser > 0 {
					totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
					totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
				}

				// get permission admin
				isPermission := false
				if roleID > 1 {
					var whr tables.FormUserPermissionJoin
					whr.PermissionID = 6 //(6 is edit responden)
					whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
					getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error": err,
						})
						return
					}
					isPermission = getPermission.Status
				}

				if userID == result[i].CreatedBy {
					isPermission = true
				}

				if roleID == 1 {
					isPermission = true
				}

				var each objects.Forms
				each.ID = result[i].ID
				each.Name = result[i].Name
				each.Description = result[i].Description
				each.ProfilePic = result[i].ProfilePic
				each.FormStatusID = result[i].FormStatusID
				each.FormStatus = result[i].FormStatus
				each.Notes = result[i].Notes
				each.PeriodStartDate = result[i].PeriodStartDate
				each.PeriodEndDate = result[i].PeriodEndDate
				each.CreatedBy = result[i].CreatedBy
				each.CreatedByName = result[i].CreatedByName
				each.CreatedByEmail = result[i].CreatedByEmail
				each.TotalResponden = len(getResponden)
				each.TotalRespondenActive = len(getActiveRespondens)
				each.TotalRespon = len(getDataRespons)
				each.TotalPerformance = totalPerform
				each.TotalPerformanceFloat = totalPerformFloat
				each.IsAttendanceRequired = result[i].IsAttendanceRequired
				each.UpdatedByName = ""
				each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
				each.SubmissionTarget = result[i].SubmissionTargetUser
				each.PeriodeRange = 0
				each.IsEditResponden = isPermission

				res = append(res, each)
			}
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

		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
			"paging":  paging,
			"detail":  detail,
		})
		return
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

}

func (ctr *formController) UpdateAdminPermission(c *gin.Context) {

	// claims := jwt.ExtractClaims(c)
	// userID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.UserOrganizationPermissions
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

	if len(reqData.Data) > 0 {
		for i := 0; i < len(reqData.Data); i++ {
			var whreP tables.UserOrganizationPermission
			whreP.IsChecked = reqData.Data[i].IsChecked
			_, err := ctr.permissMod.UpdateUserOrganizationPermission(reqData.Data[i].ID, whreP)
			if err != nil {
				fmt.Println("ERR: UpdateUserOrganizationPermission in loop ", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success update user permission",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed update user permission, please try again",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) AdminListPermission(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	userID, _ := strconv.Atoi(claims["id"].(string))
	fmt.Println("ini userId :", userID)
	isOwner := true

	companyID := 0

	if len(claims) >= 5 {
		companyID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("companyID :::", companyID)
	}
	fmt.Println("roleID :::", roleID)

	if companyID > 0 {
		if roleID == 1 {
			isOwner = true
			var res []objects.UserOrgPermissionRes

			var whre tables.Permissions
			whre.HttpPath = "/admin-global"
			getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			for i := 0; i < len(getPermiss); i++ {
				var each objects.UserOrgPermissionRes
				each.ID = 0
				each.UserOrganizationID = 0
				each.PermissionID = getPermiss[i].ID
				each.PermissionName = getPermiss[i].Name
				each.IsChecked = true

				res = append(res, each)
			}
			var resData objects.GlobalAdmin
			resData.IsOwner = isOwner
			resData.UserOrgPermissionRes = res

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    resData,
			})
			return
		} else {
			isOwner = false
			var res []objects.UserOrgPermissionRes

			getUserOrgPermiss, err := ctr.permissMod.GetUserOrganizationPermissionRow(userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			for i := 0; i < len(getUserOrgPermiss); i++ {

				var each objects.UserOrgPermissionRes
				each.ID = getUserOrgPermiss[i].ID
				each.UserOrganizationID = getUserOrgPermiss[i].UserOrganizationID
				each.PermissionID = getUserOrgPermiss[i].PermissionID
				each.PermissionName = getUserOrgPermiss[i].PermissionName
				each.IsChecked = getUserOrgPermiss[i].IsChecked

				res = append(res, each)
			}

			var resData objects.GlobalAdmin
			resData.IsOwner = isOwner
			resData.UserOrgPermissionRes = res

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    resData,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed company ID is not found",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) ListAdminEks(c *gin.Context) {
	userSearch := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	claims := jwt.ExtractClaims(c)
	companyID := 0

	// var getTotalFormEksternal []objects.FormData
	var numrows int

	if len(claims) >= 5 {
		companyID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println("companyID :::", companyID)
	}

	if companyID > 0 {

		whreStr := ""
		var buffer bytes.Buffer
		if userSearch != "" {
			buffer.WriteString(" AND u.name ilike '%" + userSearch + "%' AND u2.name ilike '%" + userSearch + "%'")
		}

		whreStr = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort
		var res []objects.AdminEks

		getAdminEks, err := ctr.userMod.GetListAdminEks(companyID, whreStr, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		getAdminEksAll, err := ctr.userMod.GetListAdminEks(companyID, whreStr, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(getAdminEks) > 0 {
			for i := 0; i < len(getAdminEks); i++ {
				getTotalFormEksternal, err := ctr.userMod.GetTotalFormAdminEksternal(getAdminEks[i].UserID, companyID)
				fmt.Println(getTotalFormEksternal)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				var each objects.AdminEks
				// each.ID = getAdminEks[i].ID
				each.UserID = getAdminEks[i].UserID
				each.OrganizationID = getAdminEks[i].OrganizationID
				each.Name = getAdminEks[i].Name
				each.Email = getAdminEks[i].Email
				each.Phone = getAdminEks[i].Phone
				each.OrganizationName = getAdminEks[i].OrganizationName
				each.TotalForm = len(getTotalFormEksternal)

				res = append(res, each)
			}
			totalPage := 0
			if limit > 0 {
				totalPage = len(getAdminEksAll) / limit
				if (len(getAdminEksAll) % limit) > 0 {
					totalPage = totalPage + numrows
				}
			}

			var paging objects.DataRows
			paging.TotalRows = len(getAdminEksAll) + numrows
			paging.TotalPages = totalPage

			var detail objects.DataRowsDetail
			detail.AllRows = len(getAdminEks) + numrows

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
				"detail":  detail,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is empty",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed company ID is not found",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) DeleteAdminEks(c *gin.Context) {
	// formID, _ := strconv.Atoi(c.Param("id"))
	adminEksID, _ := strconv.Atoi(c.Param("admineksid"))
	// if adminEksID == ' ' {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Field ID in URL is required",
	// 	})
	// 	return
	// }
	// fmt.Println(adminEksID)
	// os.Exit(0)
	getAllID, _ := ctr.formMod.GetAllIDForDelete(adminEksID)

	var res objects.DeleteAdminEksObj
	res.ID = adminEksID
	res.UserID = getAllID.UserID
	res.OrganizationID = getAllID.OrganizationID
	res.UserOrganizationID = getAllID.UserOrganizationID
	res.UserOrganizationRolesID = getAllID.UserOrganizationRolesID
	res.FormOrganizationInvitesID = getAllID.FormOrganizationInvitesID

	_, err := ctr.formMod.DeleteAdminEks(res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success Delete",
	})
	return
}

func (ctr *formController) FormCompanySharingList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID := 0
	if len(claims) >= 5 {
		organizationID, _ = strconv.Atoi(claims["organization_id"].(string))
		fmt.Println(userID, organizationID)
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	var result []tables.FormAll
	var resultAll []tables.FormAll
	var resultgetAll []tables.FormAll

	if roleID == 1 { // 1 is owner

		var whrComp tables.Organizations
		whrComp.CreatedBy = userID
		whrComp.IsDefault = true
		getComp, _ := ctr.compMod.GetCompaniesRow(whrComp)

		// SUPER ADMIN here
		var buffer bytes.Buffer
		var fields tables.FormOrganizationsJoin
		fields.OrganizationID = getComp.ID

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		// getForms, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		getForms, err := ctr.formMod.GetFormOtherCompanyRows(fields, whereString, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		result = getForms

		// get all data
		getFormsAll, err := ctr.formMod.GetFormOtherCompanyRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultAll = getFormsAll

		whereAll := " forms.form_status_id not in (3) AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

		getAll, err := ctr.formMod.GetFormOtherCompanyRows(fields, whereAll, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll

	} else {
		var buffer bytes.Buffer
		var fields tables.Forms

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" f.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" f.form_status_id not in (3) AND f.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		getForms, err := ctr.formMod.GetListFormEksternal(fields, whereString, paging, userID, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		result = getForms
		//get all data
		getFormsAll, err := ctr.formMod.GetListFormEksternal(fields, whereString, objects.Paging{}, userID, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultAll = getFormsAll

		whereAll := "f.form_status_id not in (3) AND f.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

		getAll, err := ctr.formMod.GetListFormEksternal(fields, whereAll, objects.Paging{}, userID, organizationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll
	}

	if len(result) > 0 {

		var res []objects.FormSharing
		for i := 0; i < len(result); i++ {

			// get total responden
			var whereFU tables.JoinFormUsers
			whereFU.FormID = result[i].ID
			whereFU.Type = "respondent"
			whreStr := "fuo.organization_id = " + strconv.Itoa(organizationID) + ""

			getResponden, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//total active responden
			whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd') AND fuo.organization_id = " + strconv.Itoa(organizationID) + ""
			var whreActive tables.InputForms
			getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// total respon
			var whereInForm tables.InputForms

			whreStr = "fuo.organization_id = " + strconv.Itoa(organizationID) + ""
			getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, whreStr, objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			var data tables.Organizations
			data.ID = organizationID
			getCompaniesRow, err := ctr.compMod.GetCompaniesRow(data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//total performance
			var totalPerform int
			var totalPerformFloat float64
			if result[i].SubmissionTargetUser > 0 {
				totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
				totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
			}

			// get permission admin
			isPermission := false
			if roleID > 1 {
				var whr tables.FormUserPermissionJoin
				whr.PermissionID = 6 //(6 is edit responden)
				whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(result[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
				getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
				isPermission = getPermission.Status
			}

			if userID == result[i].CreatedBy {
				isPermission = true
			}

			if roleID == 1 {
				isPermission = true
			}

			var each objects.FormSharing
			each.ID = result[i].ID
			each.Name = result[i].Name
			each.Description = result[i].Description
			each.ProfilePic = result[i].ProfilePic
			each.FormStatusID = result[i].FormStatusID
			each.FormStatus = result[i].FormStatus
			each.Notes = result[i].Notes
			each.PeriodStartDate = result[i].PeriodStartDate
			each.PeriodEndDate = result[i].PeriodEndDate
			each.CreatedBy = result[i].CreatedBy
			each.CreatedByName = result[i].CreatedByName
			each.CreatedByEmail = result[i].CreatedByEmail
			each.TotalResponden = len(getResponden)
			each.TotalRespondenActive = len(getActiveRespondens)
			each.TotalRespon = len(getDataRespons)
			each.TotalPerformance = totalPerform
			each.TotalPerformanceFloat = totalPerformFloat
			each.IsAttendanceRequired = result[i].IsAttendanceRequired
			each.UpdatedByName = ""
			each.LastUpdate = result[i].UpdatedAt.Format("2006-02-01 15:04")
			each.SubmissionTarget = result[i].SubmissionTargetUser
			each.PeriodeRange = 0
			each.IsEditResponden = isPermission

			sharingSaldo := "Tidak"
			if result[i].IsQuotaSharing == true {
				sharingSaldo = "Ya"
			}
			each.SharingSaldo = sharingSaldo
			each.StatusAdmin = "Editor"
			each.OrganizationName = getCompaniesRow.Name
			each.OrganizationProfilePic = getCompaniesRow.ProfilePic

			res = append(res, each)
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

		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
			"paging":  paging,
			"detail":  detail,
		})
		return
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

}

func (ctr *formController) FormExternalSharingList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	var results []tables.FormAll
	var resultAll []tables.FormAll
	var resultgetAll []tables.FormAll
	var err error

	// 6 is role GUEST
	if roleID == 7 {
		var buffer bytes.Buffer
		var fields tables.Forms

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
		whereString = buffer.String()

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		results, err = ctr.formMod.GetFormNotInProjectRows(fields, whereString, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		//get all data
		getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		resultAll = getFormsAll
		whereAll := "forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"
		getAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereAll, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data role ID is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

	if len(results) > 0 {

		var res []objects.FormSharing
		for i := 0; i < len(results); i++ {

			// get total responden
			var whereFU tables.JoinFormUsers
			whereFU.FormID = results[i].ID
			whereFU.Type = "respondent"
			whreStr := ""

			getResponden, err := ctr.formMod.GetFormUserRows(whereFU, whreStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//total active responden
			whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			var whreActive tables.InputForms
			getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(results[i].ID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// total respon
			var whereInForm tables.InputForms
			getDataRespons, err := ctr.inputForm.GetInputFormRows(results[i].ID, whereInForm, "", objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//total performance
			var totalPerform int
			var totalPerformFloat float64
			if results[i].SubmissionTargetUser > 0 {
				totalPerformFloat = float64(len(getDataRespons)) / float64(results[i].SubmissionTargetUser)
				totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
			}

			// get permission admin
			isPermission := false
			if roleID > 1 {
				var whr tables.FormUserPermissionJoin
				whr.PermissionID = 6 //(6 is edit responden)
				whrStr := "form_user_id in (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(results[i].ID) + " AND fu.user_id=" + strconv.Itoa(userID) + " )"
				getPermission, err := ctr.permissMod.GetFormUserPermissionRow(whr, whrStr)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
					return
				}
				isPermission = getPermission.Status
			}

			if userID == results[i].CreatedBy {
				isPermission = true
			}

			if roleID == 1 {
				isPermission = true
			}

			getOrg, _ := ctr.compMod.GetFormOrganizationRow(tables.FormOrganizations{FormID: results[i].ID})

			var each objects.FormSharing
			each.ID = results[i].ID
			each.Name = results[i].Name
			each.Description = results[i].Description
			each.ProfilePic = results[i].ProfilePic
			each.FormStatusID = results[i].FormStatusID
			each.FormStatus = results[i].FormStatus
			each.Notes = results[i].Notes
			each.PeriodStartDate = results[i].PeriodStartDate
			each.PeriodEndDate = results[i].PeriodEndDate
			each.CreatedBy = results[i].CreatedBy
			each.CreatedByName = results[i].CreatedByName
			each.CreatedByEmail = results[i].CreatedByEmail
			each.TotalResponden = len(getResponden)
			each.TotalRespondenActive = len(getActiveRespondens)
			each.TotalRespon = len(getDataRespons)
			each.TotalPerformance = totalPerform
			each.TotalPerformanceFloat = totalPerformFloat
			each.IsAttendanceRequired = results[i].IsAttendanceRequired
			each.UpdatedByName = ""
			each.LastUpdate = results[i].UpdatedAt.Format("2006-02-01 15:04")
			each.SubmissionTarget = results[i].SubmissionTargetUser
			each.PeriodeRange = 0
			each.IsEditResponden = isPermission
			each.OrganizationName = getOrg.OrganizationName
			each.OrganizationProfilePic = getOrg.OrganizationProfilePic
			each.SharingSaldo = "-"
			each.StatusAdmin = "Viewer"

			res = append(res, each)
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

		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
			"paging":  paging,
			"detail":  detail,
		})
		return
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

}

func (ctr *formController) FormExternalSharingTotal(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	var result []tables.FormAll
	// var resultAll []tables.FormAll
	var resultgetAll []tables.FormAll

	// 7 is role GUEST
	if roleID == 7 {
		var buffer bytes.Buffer
		var fields tables.Forms

		whereString := ""
		if searchKeyWord != "" {
			buffer.WriteString(" forms.name ilike '%" + searchKeyWord + "%' AND  ")
		}

		buffer.WriteString(" forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))")
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
		//get all data
		// getFormsAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereString, objects.Paging{})
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }
		// resultAll = getFormsAll

		whereAll := "forms.form_status_id not in (3) AND (forms.id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ") OR forms.created_by = " + strconv.Itoa(userID) + ")AND forms.id not in (select pf.form_id from frm.project_forms pf where pf.project_id in (select p.id from frm.projects p where p.created_by = " + strconv.Itoa(userID) + "))"

		getAll, err := ctr.formMod.GetFormNotInProjectRows(fields, whereAll, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		resultgetAll = getAll
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data role ID is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

	if len(result) > 0 {

		var totalform int
		var totalrespondenactive int
		var submission int
		var res []objects.FormSharing
		for i := 0; i < len(result); i++ {

			//total active responden
			whreStrAU := "TO_CHAR(if.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')"
			var whreActive tables.InputForms
			getActiveRespondens, err := ctr.inputForm.GetActiveUserInputForm(result[i].ID, whreActive, whreStrAU)
			if err != nil {
				fmt.Println("err: GetFormUserRows", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			// total respon
			var whereInForm tables.InputForms
			getDataRespons, err := ctr.inputForm.GetInputFormRows(result[i].ID, whereInForm, "", objects.Paging{})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			//total performance
			var totalPerform int
			var totalPerformFloat float64
			if result[i].SubmissionTargetUser > 0 {
				totalPerformFloat = float64(len(getDataRespons)) / float64(result[i].SubmissionTargetUser)
				totalPerform, _ = strconv.Atoi(strconv.FormatFloat(totalPerformFloat, 'f', 0, 64))
			}

			fmt.Println(totalPerform, result[i].SubmissionTargetUser)
			var each objects.FormSharing

			totalform += 1
			totalrespondenactive += len(getActiveRespondens)
			submission += len(getDataRespons)
			res = append(res, each)
		}

		getDataUser, _ := ctr.userMod.GetUserRow(tables.Users{ID: userID})

		var total objects.Home
		total.UserID = userID
		total.UserName = getDataUser.Name
		total.UserAvatar = getDataUser.Avatar
		total.CompanyName = ""
		total.IsUserCompanyActive = true
		total.IsProfileComplete = true
		total.TotalForm = totalform
		total.TotalResponden = totalrespondenactive
		total.TotalRespon = submission

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data user is available",
			"data":    total,
		})
		return
	} else {
		var detail objects.DataRowsDetail
		detail.AllRows = len(resultgetAll)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data is not available",
			"data":    nil,
			"detail":  detail,
		})
		return
	}

}

func (ctr *formController) AddAdminPermisMan(c *gin.Context) {
	var reqData tables.UserOrganizationPermission
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

	// fmt.Println(reqData.UserOrganizationID)
	// os.Exit(0)
	var whre tables.Permissions
	whre.HttpPath = "/admin-global"

	getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
	if len(getPermiss) > 0 {
		fmt.Println(len(getPermiss))
		for i := 0; i < len(getPermiss); i++ {
			var postData tables.UserOrganizationPermission
			postData.UserOrganizationID = reqData.UserOrganizationID
			postData.PermissionID = getPermiss[i].ID
			_, err = ctr.permissMod.InsertAttendanceOrganization(postData)
			if err != nil {
				if errors.As(err, &ctr.pgErr) {
					if ctr.pgErr.Code == "23505" || ctr.pgErr.Code == "23503" { //code duplicate

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
}

func (ctr *formController) AddAdminPermisOto(c *gin.Context) {
	var fields objects.IDAdminEks
	getMissingID, _ := ctr.permissMod.GetMissingID(fields)

	if len(getMissingID) > 0 {

		var whre tables.Permissions
		whre.HttpPath = "/admin-global"

		getPermiss, err := ctr.permissMod.GetPermissionRows(whre, "")
		for i := 0; i < len(getMissingID); i++ {
			for j := 0; j < len(getPermiss); j++ {

				var postData tables.UserOrganizationPermission
				postData.UserOrganizationID = getMissingID[i].ID
				postData.PermissionID = getPermiss[j].ID
				postData.IsChecked = true
				// fmt.Println(postData)
				// os.Exit(0)
				_, err = ctr.permissMod.InsertAttendanceOrganization(postData)
				if err != nil {
					if errors.As(err, &ctr.pgErr) {
						if ctr.pgErr.Code == "23505" || ctr.pgErr.Code == "23503" { //code duplicate

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
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success sign data to form",
			"data":    nil,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data Kosong brader",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FilterFormToCompanyList(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	// organizationID, _ := strconv.Atoi(claims["organization_id"].(string))
	fmt.Println(userID)

	formID, err := strconv.Atoi(c.Request.URL.Query().Get("form_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")

	if formID > 0 {
		// company data
		// getCompany, err := ctr.compMod.GetCompaniesRow(tables.Organizations{ID: organizationID})
		getFormComp, err := ctr.formMod.GetFormOrganization(tables.FormOrganizations{FormID: formID})

		// cekform id
		var frm tables.Forms
		frm.ID = formID
		getForm, err := ctr.formMod.GetFormRow(frm)
		if err != nil {
			fmt.Println("err: GetFormRow", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if getForm.ID == 0 {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Form ID is wrong",
				"error":   err,
			})
			return
		}

		whereString := ""
		if searchKeyWord != "" {
			whereString = " o.name ilike '%" + searchKeyWord + "%'  "
		}

		fuRows, err := ctr.formMod.GetFormCompanyInviteRows(tables.JoinFormCompanies{FormID: formID}, whereString)
		if err != nil {
			fmt.Println("err: GetFormCompanyInviteRows", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var results []objects.FilterFormCompanyList

		// static list my company
		for i := 0; i < 1; i++ {
			var each objects.FilterFormCompanyList
			each.ID = getFormComp.OrganizationID
			each.OrganizationName = getFormComp.OrganizationName

			results = append(results, each)
		}

		if len(fuRows) >= 1 {

			for i := 0; i < len(fuRows); i++ {
				var each objects.FilterFormCompanyList
				each.ID = fuRows[i].OrganizationID
				each.OrganizationName = fuRows[i].OrganizationName

				results = append(results, each)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Form company list data is available",
			"data":    results,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *formController) FormCompanySharingDelete(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, _ := strconv.Atoi(c.Param("formid"))

	if formID >= 0 && userID >= 0 {
		deletedFormUser, err := ctr.formMod.DeleteFormUserOrg(userID, formID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"status":  false,
				"message": "Error: Data user is not deleted",
			})
			return
		}

		if deletedFormUser {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Success: Data user has deleted",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed: Data user is not deleted",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data form ID is required",
		})
		return
	}
}

func (ctr *formController) GetFillingType(c *gin.Context) {
	var where objects.FillingType
	getFillingType, err := ctr.formMod.GetFillingType(where)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"status":  false,
			"message": "Error: Data user is not deleted",
		})
		return
	}

	var res []objects.FillingType

	if len(getFillingType) > 0 {

		for i := 0; i < len(getFillingType); i++ {

			var each objects.FillingType
			each.ID = getFillingType[i].ID
			each.Name = getFillingType[i].Name
			each.Status = getFillingType[i].Status

			res = append(res, each)

		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
		})
		return
	}

}

func (ctr *formController) FormDataExport(c *gin.Context) {

	getForms, _ := ctr.formMod.GetFormByOrganization(objects.HistoryBalanceSaldo{OrganizationID: 30})

	results, err := ctr.formMod.GetDate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(results) >= 1 {
		fmt.Println("len(results)", len(results))
		xlsx := excelize.NewFile()
		sheetName := "Export-Data"
		xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)

		xlsx.SetCellValue(sheetName, "A1", "Tanggal")
		xlsx.SetCellValue(sheetName, "B1", "Total Submission")
		xlsx.SetCellValue(sheetName, "C1", "Total Responden")
		xlsx.SetCellValue(sheetName, "D1", "Total Form")

		row := 2

		for i := 0; i < len(results); i++ {
			xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(row), results[i].Date)

			// resForm, err := ctr.formMod.GetAllForm()
			// if err != nil {
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"error": err.Error(),
			// 	})
			// 	return
			// }

			// fmt.Println("len(getForms)", len(getForms))
			totalData := 0
			totalUsers := 0
			totalFormActive := 0
			if len(getForms) >= 1 {

				for j := 0; j < len(getForms); j++ {

					whr := " to_char(if.created_at, 'yyyy-mm-dd') = '" + results[i].Date + "'"
					getData, _ := ctr.inputForm.GetInputDataRows(getForms[j].FormID, "", tables.InputForms{}, whr)
					getUsers, _ := ctr.inputForm.GetActiveUserInputFormOld(getForms[j].FormID, tables.InputForms{}, whr)

					totalData += len(getData)
					totalUsers += len(getUsers)
					if len(getData) >= 1 {

						totalFormActive += 1
					}
				}
			}
			xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(row), totalData)
			xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(row), totalUsers)
			xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(row), totalFormActive)
			row++
		}

		today := time.Now()
		dateFormat := today.Format("02012006-1504")

		fileName := "Export-Snapin-" + dateFormat
		fileGroup := "general_download"
		fileExtention := "xlsx"
		fileLocation := "file/" + fileName + "." + fileExtention // local path file location

		err2 := xlsx.SaveAs(fileLocation)
		if err2 != nil {
			fmt.Println(err2)
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
			"message": "Data export is available",
			"data":    obj,
		})
		return
	}

	/*
		if len(results) > 0 {
			xlsx := excelize.NewFile()
			sheetName := "Export-Absen-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)

			style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#e4e4e4"],"pattern":1}}`)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadGateway, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}
			xlsx.SetCellStyle(sheetName, "A1", "K1", style)
			xlsx.SetColWidth(sheetName, "A", "K", 20)

			xlsx.SetCellValue(sheetName, "A1", "Nama")
			xlsx.SetCellValue(sheetName, "B1", "Kontak")
			xlsx.SetCellValue(sheetName, "C1", "Selfie Absen Masuk")
			xlsx.SetCellValue(sheetName, "D1", "Tanggal Sistem Absen Masuk")
			xlsx.SetCellValue(sheetName, "E1", "Waktu Absen Masuk")
			xlsx.SetCellValue(sheetName, "F1", "Lokasi Absen Masuk")

			xlsx.SetCellValue(sheetName, "G1", "Selfie Absen Keluar")
			xlsx.SetCellValue(sheetName, "H1", "Tanggal Sistem Absen Keluar")
			xlsx.SetCellValue(sheetName, "I1", "Waktu Absen Keluar")
			xlsx.SetCellValue(sheetName, "J1", "Lokasi Absen Keluar")
			xlsx.SetCellValue(sheetName, "K1", "Durasi Kerja")

			row := 2
			for i := 0; i < len(results); i++ {

				updatedAt := ""
				if results[i].AttendanceOut != "" {
					updatedAt = results[i].UpdatedAt
				}

				xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(row), results[i].UserName)
				xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(row), results[i].UserPhone)
				xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(row), results[i].FacePicIn)
				xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(row), results[i].CreatedAt)
				xlsx.SetCellValue(sheetName, "E"+strconv.Itoa(row), results[i].AttendanceIn)
				xlsx.SetCellValue(sheetName, "F"+strconv.Itoa(row), results[i].AddressIn)
				xlsx.SetCellValue(sheetName, "G"+strconv.Itoa(row), results[i].FacePicOut)
				xlsx.SetCellValue(sheetName, "H"+strconv.Itoa(row), updatedAt)
				xlsx.SetCellValue(sheetName, "I"+strconv.Itoa(row), results[i].AttendanceOut)
				xlsx.SetCellValue(sheetName, "J"+strconv.Itoa(row), results[i].AddressOut)
				xlsx.SetCellValue(sheetName, "K"+strconv.Itoa(row), results[i].Duration)
				row++
			}

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			fileName := "Absensi-Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "attendance_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(err2)
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
				"message": "Data export is available",
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
	*/
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is  available",
		"data":    results,
	})
	return
}
func (ctr *formController) CheckForm(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))
	roleID, _ := strconv.Atoi(claims["role_id"].(string))
	organizationID, _ := strconv.Atoi(claims["organization_id"].(string))

	fmt.Println(userID, roleID, organizationID)
}

func (ctr *formController) ListFormTemplate(c *gin.Context) {

	UserID := 2546
	ProjectID, _ := strconv.Atoi(c.Request.URL.Query().Get("jkid"))

	if ProjectID <= 0 {

		getFormTemplate, _ := ctr.formMod.GetFormTemplate(UserID)
		var res []objects.FormTemplateNew

		var each objects.FormTemplateNew

		staticValue := objects.FormTemplateNew{
			FormID:      0,
			ProjectID:   0,
			Name:        "Form Kosong",
			ProfilePic:  "https://srv-asset-snapinnew.oss-ap-southeast-5.aliyuncs.com/dev/forms/form_image_64e3027ade9ab40db59494ec.jpg",
			Description: "",
			Total:       0,
		}
		res = append(res, staticValue)
		for i := 0; i < len(getFormTemplate); i++ {

			var fields tables.FormFields
			fields.FormID = getFormTemplate[i].FormID
			result, _ := ctr.formFieldMod.GetFormFieldNotParentRows(fields, "")

			each.FormID = getFormTemplate[i].FormID
			each.ProjectID = getFormTemplate[i].ProjectID
			each.Name = getFormTemplate[i].Name
			each.ProfilePic = getFormTemplate[i].ProfilePic
			each.Description = getFormTemplate[i].Description
			each.Total = len(result)

			res = append(res, each)

		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    res,
		})
		return

	} else {
		getFormTemplateByProjectID, _ := ctr.formMod.GetFormTemplateByProjectID(UserID, ProjectID)
		var res []objects.FormTemplateNew

		if len(getFormTemplateByProjectID) > 0 {
			var each objects.FormTemplateNew
			for i := 0; i < len(getFormTemplateByProjectID); i++ {

				var fields tables.FormFields
				fields.FormID = getFormTemplateByProjectID[i].FormID
				result, _ := ctr.formFieldMod.GetFormFieldNotParentRows(fields, "")

				each.FormID = getFormTemplateByProjectID[i].FormID
				each.ProjectID = getFormTemplateByProjectID[i].ProjectID
				each.Name = getFormTemplateByProjectID[i].Name
				each.ProfilePic = getFormTemplateByProjectID[i].ProfilePic
				each.Description = getFormTemplateByProjectID[i].Description
				each.Total = len(result)

				res = append(res, each)

			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is not available",
				"data":    []string{},
			})
			return
		}

	}

}

func (ctr *formController) ListProject(c *gin.Context) {

	UserID := 2546

	getProject, _ := ctr.formMod.GetProject(UserID)
	var res []objects.Projects

	var each objects.Projects
	for i := 0; i < len(getProject); i++ {

		each.ID = getProject[i].ID
		each.Name = getProject[i].Name
		each.Description = getProject[i].Description

		res = append(res, each)

	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data is available",
		"data":    res,
	})
	return

}
