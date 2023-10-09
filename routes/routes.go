package routes

import (
	"log"
	"snapin-form/config"
	"snapin-form/controllers"
	"snapin-form/helpers"
	m "snapin-form/middleware"
	"snapin-form/models"

	"gorm.io/gorm"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
)

// var (
// 	varconf        config.Configurations =
// 	db             *gorm.DB                   = config.ConnectDB(varconf)
// 	userModels     models.UserModels          = models.NewUserModels(db)
// 	userController controllers.UserController = controllers.NewUserController(userModels)
// )

func SetupRoutes(conf config.Configurations) *gin.Engine {

	r := gin.Default()

	var varconf config.Configurations = conf

	var db *gorm.DB = config.ConnectDB(varconf)

	var formModels models.FormModels = models.NewFormModels(db)
	var helpers helpers.Helper = helpers.NewHelper(varconf, db)
	var ftModels models.FieldTypeModels = models.NewFieldTypeModels(db)
	var formFieldModels models.FormFieldModels = models.NewFormFieldModels(db)
	var ruleModels models.RuleModels = models.NewRuleModels(db)
	var projectModels models.ProjectModels = models.NewProjectModels(db)
	var inputFormModels models.InputFormModels = models.NewInputFormModels(db)
	var userModels models.UserModels = models.NewUserModels(db)
	var compModels models.CompaniesModels = models.NewCompaniesModels(db)
	var subsModels models.SubsModels = models.NewSubsModels(db)
	var permissModels models.PermissionModels = models.NewPermissionModels(db)
	var attendModels models.AttendanceModels = models.NewAttendanceModels(db)
	var settingModels models.SettingModels = models.NewSettingModels(db)
	var formOtpModes models.FormOtpModels = models.NewFormOtpModels(db)
	var shortenModels models.ShortenUrlModels = models.NewShortenUrlModels(db)
	var regionModels models.RegionModels = models.NewRegionModels(db)
	var multiAccessModels models.MultiAccessModels = models.NewMultiAccessModels(db)

	var homeController controllers.HomeController = controllers.NewHomeController(userModels, formModels, helpers, inputFormModels, compModels, subsModels)
	var formController controllers.FormController = controllers.NewFormController(formModels, formFieldModels, ftModels, ruleModels, helpers, userModels, compModels, inputFormModels, permissModels, projectModels, attendModels)
	var projectController controllers.ProjectController = controllers.NewProjectController(projectModels, formModels, inputFormModels, compModels, permissModels)
	var subsController controllers.SubscriptionController = controllers.NewSubsController(formModels, inputFormModels)
	var fileController controllers.FileController = controllers.NewFileController(helpers)
	var inputFormCtr controllers.InputFormController = controllers.NewInputFormController(inputFormModels, formFieldModels, formModels, helpers, userModels, subsModels, compModels, attendModels, settingModels, shortenModels, varconf)
	var appController controllers.AppController = controllers.NewAppController(formModels, formFieldModels, ftModels, ruleModels, helpers, projectModels, userModels, inputFormModels, attendModels, compModels, formOtpModes, settingModels, varconf, subsModels, shortenModels)
	var attController controllers.AttendanceController = controllers.NewAttController(formModels, userModels, attendModels, helpers, shortenModels)
	var regionController controllers.RegionController = controllers.NewRegionController(regionModels)
	var multiAccessController controllers.MultiAccessController = controllers.NewMultiAccessController(helpers, multiAccessModels, formModels, compModels, userModels, varconf, permissModels)

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,content-type,authorization, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Set("db", db)
		c.Next()
	})

	authMiddleware := m.SetupMiddleware(db)

	// When you use jwt.New(), the function is already automatically called for checking,// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// rules
	r.GET("/rules/:type", formController.ConditionRuleList)

	// file save global fitur
	r.POST("/file/save", fileController.FieldSaveData)
	r.POST("/file/save_video", fileController.FieldSaveDataVideo)
	r.POST("/file/save_any", fileController.FieldSaveDataVideo)
	r.POST("/attendance/insert_existing_data", attController.AutoInsertData)
	r.GET("/export", formController.FormDataExport)

	dash := r.Group("/dash")
	dash.Use(authMiddleware.MiddlewareFunc())
	{
		dash.GET("/home", homeController.Home)
		dash.GET("/home/:groupid", homeController.Home)

		dash.GET("/home_non_account", formController.FormExternalSharingTotal) // (GUEST) list formuser extenal tanpa company ID di token

	}

	form := r.Group("/form")
	form.Use(authMiddleware.MiddlewareFunc())
	{
		// forms
		form.GET("/group1list", formController.FormGroup1List)
		form.GET("/group2list", formController.FormGroup2List)
		form.GET("/group3list", formController.FormGroup3List) // latest
		form.GET("/list", formController.FormList)
		form.GET("/:id", formController.FormList)
		form.POST("/create", formController.FormCreate)
		form.PUT("/update/:id", formController.FormUpdate)
		form.DELETE("/destroy/:id", formController.FormDestroy)
		form.GET("/archive", formController.FormListArchiveLast)
		form.POST("/duplicate", formController.FormDuplicate3)

		form.POST("/createlocation", formController.FormCreateLocation)
		form.PUT("/updatelocation/:id", formController.FormUpdateLocation)
		form.DELETE("/destroylocation/:id", formController.FormDestroyLocation)

		// form question list
		form.GET("/fieldtypes", formController.FieldTypeList)
		form.GET("/fields/:formid", formController.FieldList)
		// form.GET("/fieldgroups/:formid", formController.FieldGroupList)
		// form.GET("/fieldgroups/:formid", formController.Field2GroupList)
		form.GET("/fieldgroups", formController.FieldGroup3List)         //list questin apps current
		form.GET("/fieldgroups/:formid", formController.FieldGroup3List) //list questin apps current
		form.POST("/fieldsave", formController.FieldCreate)
		form.PUT("/fieldupdate/:fieldid", formController.FieldUpdate)
		form.DELETE("/fielddestroy/:fieldid", formController.FieldDestroy)
		form.POST("/fieldconditionsave", formController.FieldConditionSave)
		form.POST("/fieldgroupsave", formController.FieldGroupCreate)
		form.PUT("/fieldgroupupdate/:fieldid", formController.FieldGroupUpdate)
		form.POST("/fieldsectionsave", formController.FieldSectionCreate)
		form.PUT("/fieldsectionupdate/:fieldid", formController.FieldSectionUpdate)
		form.GET("/filling_type", formController.GetFillingType)

		// form save data
		form.POST("/savedata", inputFormCtr.FieldSaveData)
		form.PUT("/updatedata/:dataid", inputFormCtr.FieldUpdateData)
		form.DELETE("/destroydata/:dataid", inputFormCtr.FieldDestroyData)
		form.POST("/saveimage", formController.FieldSaveImage)

		// form user
		form.GET("/formuserlist", formController.FormUserList)
		form.GET("/formuserlist/:formid", formController.FormUserList)
		form.POST("/formusercreate", formController.FormUserCreate)
		form.PUT("/formuserupdate/:formid", formController.FormUserUpdate)
		form.PUT("/formuserstatus/:formid", formController.FormUserStatusUpdate)

		// form to user
		form.POST("/formuserconnect", formController.FormUserConnect)
		form.DELETE("/formuserdisconnect", formController.FormUserDisconnect)

		form.GET("/share/:formid", formController.FormShare)
		form.PUT("/updatestatus/:formid", formController.FormUpdateStatus)
		form.PUT("/updateattendance/:formid", formController.FormAttendanceRequired)
		form.PUT("/sortordersave/:formid", formController.FieldSortOrderSave)

		// user to role form_admin
		form.POST("/useradminfrmconnect", formController.FormUserAdminFrmConnect)
		form.POST("/useradminfrmdisconnect", formController.FormUserAdminFrmDisconnect)
		form.PUT("/useradminpermissionsave", formController.FormUserAdminCheck)

		// form user permission
		form.GET("/permissions/:formid", formController.FormUserAdminFrm)
		form.GET("/attendance/:formid", attController.FormAttendanceList)
		form.GET("/download_attendance/:formid", attController.FormAttendanceListExport)
		form.GET("/download_attendance_csv/:formid", attController.FormAttendanceListExportCSV)
		form.POST("/download_attendance", attController.FormAttendanceListExportPost)
		form.GET("/profile_forms/:userid", formController.UserGetFormList)
		form.GET("/map_attendance/:formid", attController.FormAttendanceMapsList2)

		// Form Template
		form.GET("/get_form_template", formController.ListFormTemplate)
		form.GET("/get_project", formController.ListProject)
		form.POST("/template", formController.FormTemplate)

		// form company invite
		form.POST("/formcompanyconnect", formController.FormCompanyConnect)
		form.DELETE("/formcompanydisconnect", formController.FormCompanyDisconnect)
		form.GET("/formcompanylist", formController.FormToCompanyList)
		form.GET("/formcompanynotinlist", formController.FormCompanyNotInList)
		form.PUT("/formcompanyupdatequota/:frm_company_id", formController.FormCompanyUpdateQuota)
		form.GET("/filterformcompanylist", formController.FilterFormToCompanyList)

		form.GET("/formcompanysharinglist", formController.FormCompanySharingList)                // tab form multy access
		form.DELETE("/formcompanysharingdelete/:formid", formController.FormCompanySharingDelete) // delete form multy access
		form.GET("/formexternalsharinglist", formController.FormExternalSharingList)              // (GUEST) list formuser extenal tanpa company ID di token
		form.GET("/history_balance_saldo", subsController.HistoryBalanceSaldo)
		form.GET("/history_balance_saldo/:form_id", subsController.HistoryBalanceSaldo)
		form.GET("/history_balance_saldo_by_date", subsController.HistoryBalanceSaldoByDate)
		form.GET("/check_form", formController.CheckForm)

	}

	//formdata
	formData := r.Group("/formdata")
	formData.Use(authMiddleware.MiddlewareFunc())
	{
		formData.GET("/:formid", inputFormCtr.DataFormDetail)
		formData.DELETE("/delete", inputFormCtr.DataFormDelete)
		formData.GET("/download/:formid", inputFormCtr.DataFormDetail4Download)
		formData.GET("/downloadv1/:formid", inputFormCtr.DataFormDetail2Download)
		formData.GET("/download_csv/:formid", inputFormCtr.DataFormDetail4DownloadCSV)
		formData.GET("/respon/:formid", inputFormCtr.DataFormResponList)
		formData.GET("/responden/:formid/:periode/:year", inputFormCtr.DataFormRespondenList)
		// formData.GET("/grafic/:formid/:periode/:year", inputFormCtr.DataFormDetailGrafic)
		formData.GET("/grafic/:formid/:periode/:year", inputFormCtr.DataFormDetailGraficByOrganization)
		formData.GET("/map/:formid", inputFormCtr.DataFormResponMapList)
		formData.POST("/map", inputFormCtr.DataFormResponMapPostList) // existing (POST method)
	}

	// data
	data := r.Group("/data")
	data.Use(authMiddleware.MiddlewareFunc())
	{
		// data.POST("/data/savedata", inputFormCtr.FieldSaveData)
		// r.GET("/data/detail/:formid", inputFormCtr.DetailDataList)
		// generate data form usr organizations
		data.POST("/generateformuser", inputFormCtr.GenerateFormUserOrg)
		data.POST("/generateinputform", inputFormCtr.GenerateInputFormOrg)
		data.POST("/generateinputformnocopy", inputFormCtr.GenerateInputFormOrgNoCopy)
		data.POST("/generateinputformlatest", inputFormCtr.GenerateInputFormOrgLatest)

	}

	// group/projects
	project := r.Group("/project")
	project.Use(authMiddleware.MiddlewareFunc())
	{
		project.GET("/list", projectController.ProjectList)
		project.GET("/:id", projectController.ProjectList)
		project.POST("/create", projectController.ProjectCreate)
		project.PUT("/update/:id", projectController.ProjectUpdate)
		project.DELETE("/destroy/:id", projectController.ProjectDestroy)

		// project form
		project.POST("/projectformcreate", projectController.ProjectFormCreate)
		project.DELETE("/projectformdestroy", projectController.ProjectFormDestroy)
	}

	app := r.Group("/app")
	app.Use(authMiddleware.MiddlewareFunc())
	{
		app.GET("/home", appController.Home)
		app.GET("/project/list", appController.ProjectList)
		app.GET("/form/list", appController.FormList)
		app.GET("/form/performance", appController.FormPerformance)
		app.GET("/form/performance/:formid", appController.FormPerformance)

		app.POST("/sendsharecode", formController.FormGetShare)

		// submission data
		app.GET("/submission/form", appController.SubmissionForm)
		app.GET("/submission/form_data/:formid", appController.SubmissionFormData)                           // pertnyaan only
		app.GET("/submission/form_data_field/:formid/:field_data_id", appController.SubmissionFormDataField) // pertanyaan & jawaban per user
		app.POST("/submission_edit_request", appController.SubmissionEditRequest)
		app.POST("/submission_otp_request_check", appController.SubmissionFormOTPChecking)
		app.GET("/submission_user_data", appController.SubmissionDetailUser)

		//absen submit dan location
		app.POST("/attendance", attController.InsertAttendance)
		app.POST("/warn_attendance", attController.InsertWarningAttendance)
		app.POST("/attendance_offline_mode", attController.InsertWarning2Attendance)
		app.GET("/list_location_attendance/:formid", attController.GetListLocationAttendance)
		app.GET("/list_location_attendance_dash/:formid", attController.GetListLocationAttendance)

		//admin app
		app.GET("/home_admin", appController.HomeAdmin)
	}

	admin := r.Group("/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		admin.GET("/home", appController.HomeAdminContent)
		admin.GET("/forms", formController.AdminFormList)
		admin.GET("/form", formController.AdminFormListNew) //new
		admin.GET("/forms/:formid", formController.AdminFormList)
		admin.GET("/form_report_responden/:formid", inputFormCtr.ReportFormRespondenList2)
		admin.PUT("/update_admin_permission", formController.UpdateAdminPermission)
		admin.GET("/list_admin_permission", formController.AdminListPermission)
		admin.GET("/global_permission", formController.AdminListPermission) // permission admin global
		admin.GET("/list_admin_eksternal", formController.ListAdminEks)
		admin.DELETE("/delete_admin_eksternal/:admineksid", formController.DeleteAdminEks)
		admin.POST("/add_admin_permission_manual", formController.AddAdminPermisMan)
		admin.POST("/add_admin_permission_otomatic", formController.AddAdminPermisOto)
		admin.POST("/invite_admin", multiAccessController.InviteAdmin)
		admin.POST("/select_company", multiAccessController.SelectCompany)

	}

	region := r.Group("/region")
	region.Use(authMiddleware.MiddlewareFunc())
	{
		region.GET("/province", regionController.GetProvince)
		region.GET("/cities", regionController.GetCity)
		region.GET("/cities/:provinceid", regionController.GetCity)
		region.GET("/district", regionController.GetDistrict)
		region.GET("/district/:cityid", regionController.GetDistrict)
		region.GET("/sub_district", regionController.GetSubDistrict)
		region.GET("/sub_district/:districtid", regionController.GetSubDistrict)

		region.POST("/getradius", regionController.GetRadius)
	}

	return r
}
