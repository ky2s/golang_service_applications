package objects

import (
	"time"

	"gorm.io/gorm"
)

type Forms struct {
	ProjectID             int          `json:"group_id" form:"id" `
	ID                    int          `json:"id" form:"id" `
	FormStatusID          int          `json:"form_status_id" form:"form_status_id"`
	FormStatus            string       `json:"form_status" form:"form_status"`
	Name                  string       `json:"name" form:"name" binding:"required"`
	Description           string       `json:"description" form:"description"`
	Notes                 string       `json:"notes" form:"notes"`
	PeriodStartDate       string       `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate         string       `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic            string       `json:"profile_pic" form:"profile_pic"`
	UserID                int          `json:"user_id,omitempty" form:"user_id"`
	TotalForms            int          `json:"total_forms" `
	TotalResponden        int          `json:"total_responden"  gorm:"default:0"`
	TotalRespondenActive  int          `json:"total_responden_active"  gorm:"default:0"`
	TotalRespon           int          `json:"total_respon"  gorm:"default:0"`
	TotalAdmin            int          `json:"total_admin"  gorm:"default:0"`
	TotalPerformance      int          `json:"total_performance"  gorm:"default:0"`
	TotalPerformanceFloat float64      `json:"total_performance_float"  `
	CreatedBy             int          `json:"created_by" form:"created_by"`
	CreatedByName         string       `json:"created_by_name" form:"created_by_name"`
	CreatedByEmail        string       `json:"created_by_email" form:"created_by_email"`
	UpdatedByName         string       `json:"updated_by_name" form:"updated_by_name"`
	LastUpdate            string       `json:"last_updated" form:"last_updated"`
	AttendanceIn          string       `json:"attendance_in" form:"attendance_in"`
	AttendanceOut         string       `json:"attendance_out" form:"attendance_out"`
	SubmissionTarget      int          `json:"submission_target" form:"submission_target"`
	IsAttendanceRequired  bool         `json:"is_attendance_required" form:"is_attendance_required"`
	PeriodeRange          int          `json:"periode_range" form:"periode_range"`
	ArchivedAt            time.Time    `json:"archive_at" form:"archive_at" gorm:"default:null"`
	Author                string       `json:"author" `
	AuthorName            string       `json:"author_name" `
	PerformanceColor      string       `json:"performance_color" `
	ShareUrl              string       `json:"group_url" `
	TotalDeletedData      int          `json:"total_deleted_data"`
	TotalUpdatedData      int          `json:"total_updated_data"`
	LastSubmission        string       `json:"last_submission" `
	IsEditResponden       bool         `json:"is_edit_responden" `
	IsContainTagloc       bool         `json:"is_contain_tagloc" `
	AttendanceOverdate    bool         `json:"attendance_overdate"`
	AttendanceOverdateAt  time.Time    `json:"attendance_overdate_at" gorm:"default:null"`
	IsAttendanceRadius    bool         `json:"is_attendance_radius"`
	FormShared            string       `json:"form_shared"`
	Type                  string       `json:"type"`
	FormFields            []FormFields `json:"form_fields"`
}

type MergeForms struct {
	ProjectID                int          `json:"group_id" form:"id" `
	ID                       int          `json:"id" form:"id" `
	FormStatusID             int          `json:"form_status_id" form:"form_status_id"`
	FormStatus               string       `json:"form_status" form:"form_status"`
	Name                     string       `json:"name" form:"name" binding:"required"`
	Description              string       `json:"description" form:"description"`
	Notes                    string       `json:"notes" form:"notes"`
	PeriodStartDate          string       `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate            string       `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic               string       `json:"profile_pic" form:"profile_pic"`
	UserID                   int          `json:"user_id,omitempty" form:"user_id"`
	TotalForms               int          `json:"total_forms" `
	TotalResponden           int          `json:"total_responden"  gorm:"default:0"`
	TotalRespondenActive     int          `json:"total_responden_active"  gorm:"default:0"`
	TotalRespon              int          `json:"total_respon"  gorm:"default:0"`
	TotalAdmin               int          `json:"total_admin"  gorm:"default:0"`
	TotalPerformance         int          `json:"total_performance"  gorm:"default:0"`
	TotalPerformanceFloat    float64      `json:"total_performance_float"  `
	CreatedBy                int          `json:"created_by" form:"created_by"`
	CreatedByName            string       `json:"created_by_name" form:"created_by_name"`
	CreatedByEmail           string       `json:"created_by_email" form:"created_by_email"`
	UpdatedByName            string       `json:"updated_by_name" form:"updated_by_name"`
	LastUpdate               string       `json:"last_updated" form:"last_updated"`
	AttendanceIn             string       `json:"attendance_in" form:"attendance_in"`
	AttendanceOut            string       `json:"attendance_out" form:"attendance_out"`
	SubmissionTarget         int          `json:"submission_target" form:"submission_target"`
	IsAttendanceRequired     bool         `json:"is_attendance_required" form:"is_attendance_required"`
	PeriodeRange             int          `json:"periode_range" form:"periode_range"`
	ArchivedAt               time.Time    `json:"archive_at" form:"archive_at" gorm:"default:null"`
	Author                   string       `json:"author" `
	AuthorName               string       `json:"author_name" `
	PerformanceColor         string       `json:"performance_color" `
	ShareUrl                 string       `json:"group_url" `
	TotalDeletedData         int          `json:"total_deleted_data"`
	TotalUpdatedData         int          `json:"total_updated_data"`
	LastSubmission           string       `json:"last_submission" `
	IsEditResponden          bool         `json:"is_edit_responden" `
	IsContainTagloc          bool         `json:"is_contain_tagloc" `
	AttendanceOverdate       bool         `json:"attendance_overdate"`
	AttendanceOverdateAt     time.Time    `json:"attendance_overdate_at" gorm:"default:null"`
	IsAttendanceRadius       bool         `json:"is_attendance_radius"`
	NameLocation             string       `json:"name_location" form:"name_location" `
	Location                 string       `json:"location" form:"location" `
	Latitude                 float64      `json:"latitude" form:"latitude" `
	Longitude                float64      `json:"longitude" form:"longitude" `
	IsCheckIn                bool         `json:"Is_check_in" form:"Is_check_in"`
	IsCheckOut               bool         `json:"is_check_out" form:"is_check_out"`
	Radius                   int          `json:"radius" form:"radius"`
	SharingSaldo             string       `json:"sharing_saldo"`
	StatusAdmin              string       `json:"status_admin"`
	FormShared               string       `json:"form_shared"`
	FormSharedCount          string       `json:"form_shared_count"`
	FormExternalCompanyImage string       `json:"form_external_company_image"`
	FormExternalCompanyName  string       `json:"form_external_company_name"`
	Type                     string       `json:"type"`
	AccessType               string       `json:"access_type"`
	IsQuotaSharing           bool         `json:"is_quota_sharing"`
	FormFields               []FormFields `json:"form_fields"`
}

type GraficDataHours struct {
	ID         int    `json:"id"`
	FieldHours string `json:"field_hours"`
	Value      int    `json:"value"`
}

type GraficDataPeriod struct {
	ID          int    `json:"id"`
	FieldPeriod string `json:"field_period"`
	Value       int    `json:"value"`
}

type FormGraficRes struct {
	FormID         int                `json:"form_id"`
	TotalResponden int                `json:"total_responden"`
	TotalRespon    int                `json:"total_respon"`
	ActiveHours    []GraficDataHours  `json:"active_hours"`
	ActivePeriod   []GraficDataPeriod `json:"active_period"`
	ActiveDays     []GraficDataPeriod `json:"active_days"`
}

type FormResponse struct {
	ID int `json:"id" form:"id" `
}

type FormUsers struct {
	gorm.Model
	ID              int          `json:"id,omitempty"`
	FormID          int          `json:"form_id" form:"form_id" binding:"required"`
	UserID          int          `json:"user_id,omitempty" form:"user_id" binding:"required"`
	Status          bool         `json:"status" form:"status" gorm:"default:true"`
	Name            string       `json:"name" form:"name"`
	Description     string       `json:"description" form:"description"`
	Notes           string       `json:"notes" form:"notes"`
	ProfilePic      string       `json:"profile_pic" form:"profile_pic"`
	PeriodStartDate string       `json:"period_start_date" form:"period_start_date"`
	PeriodEndDate   string       `json:"period_end_date" form:"period_end_date"`
	FormStatusID    int          `json:"form_status_id" form:"form_status_id"`
	Type            string       `json:"type" form:"type"`
	FormFields      []FormFields `json:"form_fields"`
}

type FormFields struct {
	ID                     int                       `json:"id"`
	ParentID               int                       `json:"group_id" form:"group_id"`
	FieldTypeID            int                       `json:"field_type_id" form:"field_type_id" binding:"required"`
	FieldTypeName          string                    `json:"field_type_name" form:"field_type_name" `
	FormID                 int                       `json:"form_id" form:"form_id" binding:"required"`
	Label                  string                    `json:"label" form:"label" `
	Description            string                    `json:"description" form:"description"`
	Option                 string                    `json:"option" form:"option"`
	ConditionType          string                    `json:"condition_type" form:"condition_type"`
	UpperlowerCaseType     string                    `json:"upperlower_case_type" form:"upperlower_case_type"`
	IsMultiple             bool                      `json:"is_multiple" form:"is_multiple"`
	IsRequired             bool                      `json:"is_required" form:"is_required"`
	IsSection              bool                      `json:"is_section" form:"is_section"`
	SectionColor           string                    `json:"section_color" form:"section_color"`
	SortOrder              int                       `json:"sort_order" form:"sort_order"`
	Image                  string                    `json:"image" form:"image"`
	ConditionRulesID       int                       `json:"rule_id" form:"rule_id"`
	ConditionRuleValue1    string                    `json:"rule_value_1" form:"rule_value_1"`
	ConditionRuleValue2    string                    `json:"rule_value_2" form:"rule_value_2"`
	ConditionRuleMsg       string                    `json:"rule_msg" form:"rule_msg"`
	ConditionParentFieldID int                       `json:"condition_parent_field_id" form:"condition_parent_field_id"`
	TabMaxOnePerLine       bool                      `json:"tab_max_one_per_line" form:"tab_max_one_per_line"`
	TabEachLineRequire     bool                      `json:"tab_each_line_require" form:"tab_each_line_require"`
	QrCode                 string                    `json:"qr_code" form:"qr_code"`
	Conditions             []FormFieldConditionRules `json:"conditions"`
	FieldData              string                    `json:"field_data" form:"field_data"`
	TagLocColor            string                    `json:"tag_loc_color" form:"tag_loc_color"`
	TagLocIcon             string                    `json:"tag_loc_icon" form:"tag_loc_icon"`
	ProvinceID             int                       `json:"province_id" form:"province_id"`
	CityID                 int                       `json:"city_id" form:"city_id"`
	DistrictID             int                       `json:"district_id" form:"district_id"`
	SubDistrictID          int                       `json:"sub_district_id" form:"sub_district_id"`
	AddressType            string                    `json:"address_type" form:"address_type"`
	CurrencyType           int                       `json:"currency_type" form:"currency_type"`
	Currency               string                    `json:"currency" form:"currency"`
	CustomOption           string                    `json:"custom_option" `
	IsCountryPhoneCode     bool                      `json:"is_country_phone_code" form:"is_country_phone_code"`
}

type FormFieldGroups struct {
	ID          int                       `json:"group_id"`
	Label       string                    `json:"label" form:"label" `
	Description string                    `json:"description" form:"description"`
	SortOrder   int                       `json:"sort_order" form:"sort_order"`
	FormFields  []FormFields              `json:"fields" `
	Conditions  []FormFieldConditionRules `json:"conditions"`
}

type FormField3Groups struct {
	ID                     int                       `json:"id"`
	ParentID               int                       `json:"group_id" form:"group_id"`
	IsGroup                bool                      `json:"is_group" form:"is_group"`
	FieldTypeID            int                       `json:"field_type_id" form:"field_type_id" binding:"required"`
	FormID                 int                       `json:"form_id" form:"form_id" binding:"required"`
	Label                  string                    `json:"label" form:"label" `
	Description            string                    `json:"description" form:"description"`
	Option                 string                    `json:"option" form:"option"`
	ConditionType          string                    `json:"condition_type" form:"condition_type"`
	UpperlowerCaseType     string                    `json:"upperlower_case_type" form:"upperlower_case_type"`
	IsMultiple             bool                      `json:"is_multiple" form:"is_multiple"`
	IsRequired             bool                      `json:"is_required" form:"is_required"`
	IsSection              bool                      `json:"is_section" form:"is_section"`
	SectionColor           string                    `json:"section_color" form:"section_color"`
	SortOrder              int                       `json:"sort_order" form:"sort_order"`
	Image                  string                    `json:"image" form:"image"`
	ConditionRulesID       int                       `json:"rule_id" form:"rule_id"`
	ConditionRuleValue1    string                    `json:"rule_value_1" form:"rule_value_1"`
	ConditionRuleValue2    string                    `json:"rule_value_2" form:"rule_value_2"`
	ConditionRuleMsg       string                    `json:"rule_msg" form:"rule_msg"`
	ConditionParentFieldID int                       `json:"condition_parent_field_id" form:"condition_parent_field_id"`
	TabMaxOnePerLine       bool                      `json:"tab_max_one_per_line" form:"tab_max_one_per_line"`
	TabEachLineRequire     bool                      `json:"tab_each_line_require" form:"tab_each_line_require"`
	QrCode                 string                    `json:"qr_code" form:"qr_code"`
	FormFields             []FormFields              `json:"fields" `
	Conditions             []FormFieldConditionRules `json:"conditions"`
	CustomOption           string                    `json:"custom_option" `
	FieldData              string                    `json:"field_data" `
}

type FormFieldGroup struct {
	FormID      int    `json:"form_id" form:"form_id" binding:"required"`
	Label       string `json:"label" form:"label" `
	Description string `json:"description" form:"description"`
	SortOrder   int    `json:"sort_order" form:"sort_order"`
}

type FormFieldSection struct {
	ID           int                       `json:"id"`
	ParentID     int                       `json:"group_id" form:"group_id"`
	FormID       int                       `json:"form_id" form:"form_id" binding:"required"`
	Label        string                    `json:"label" form:"label" `
	Description  string                    `json:"description" form:"description"`
	IsSection    bool                      `json:"is_section" form:"is_section" `
	SectionColor string                    `json:"color" form:"color" `
	SortOrder    int                       `json:"sort_order" form:"sort_order"`
	Image        string                    `json:"image" form:"image"`
	Conditions   []FormFieldConditionRules `json:"conditions"`
}

type ConditionRulesRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type FormData struct {
	Timestamp string  `json:"timestamp" `
	FormID    int     `json:"form_id" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address" `
	UserID    int
	// IsPending bool        `json:"is_pending" gorm:"default:false"`
	FieldData []FieldData `json:"field_data" form:"field_data" `
}

type UpdCnt struct {
	UpdatedCount int `json:"updated_count"`
}

type InputFormData struct {
	ID int `json:"id"`
}

type FieldData struct {
	FieldID      int    `json:"field_id" form:"field_id"`
	Answer       string `json:"answer" form:"answer"`
	CustomOption string `json:"custom_option" form:"custom_option"`
}

type FieldFile struct {
	FormID   int    `json:"form_id" form:"form_id" `
	FieldID  int    `json:"field_id" form:"field_id" `
	File     string `json:"file" form:"file" `
	FileType string `json:"file_type" form:"file_type" `
}

type FieldFileRes struct {
	UrlFile    string `json:"url_file" form:"url_file"`
	FormFileID int    `json:"form_file_id" form:"form_file_id"`
}

type File struct {
	FileType string `json:"file_type" form:"file_type" binding:"required"`
	FileName string `json:"file_name" form:"file_name" binding:"required"`
	File     string `json:"file" form:"file" binding:"required"`
}

type FileRes struct {
	File string `json:"file" form:"file"`
}

type OssResponse struct {
	FileName string `json:"file_name" form:"file"`
	FileUrl  string `json:"file_url" form:"file"`
}

type FieldConditionSave struct {
	FormFieldID       int         `json:"field_id" form:"field_id" binding:"required"`
	ConditionAllRight bool        `json:"condition_all_right" form:"condition_all_right"`
	Conditions        []Condition `json:"conditions" form:"conditions" `
}

type Condition struct {
	ParentFieldID       int    `json:"parent_field_id" form:"parent_field_id" binding:"required"`
	ConditionRuleID     int    `json:"rule_id" form:"rule_id" binding:"required"`
	ConditionRuleValue1 string `json:"rule_value_1" form:"rule_value_1" binding:"required"`
	ConditionRuleValue2 string `json:"rule_value_2" form:"rule_value_2"`
}

type SendShareCode struct {
	ShareCode            string `json:"share_code" binding:"required"`
	NeedValidate         bool   `json:"need_validation"`
	SenderOrganizationID int    `json:"sender_organization_id"`
}

type ResFormShare struct {
	Code                 string `json:"code"`
	Url                  string `json:"url"`
	SenderOrganizationID int    `json:"sender_organization_id"`
}

type FormStatus struct {
	StatusID int  `json:"status_id"`
	Status   bool `json:"status"`
}

type FormAttendanceRequired struct {
	IsAttendanceRequired bool `json:"is_attendance_required"`
}

type FormUserStatus struct {
	UserID       int `json:"user_id" binding:"required"`
	UserStatusID int `json:"user_status_id" binding:"required"`
}

type FormSortOrder struct {
	GroupID    int          `json:"group_id" form:"group_id"`
	SortOrders []FormFields `json:"sort_order" form:"sort_order"`
}

type FormUserPermissions struct {
	Data []PermissionCheck `json:"data" form:"data"`
}

type PermissionCheck struct {
	ID     int  `json:"id" form:"id"`
	Status bool `json:"status" form:"status"`
}

type FormUserPermissionRes struct {
	ID             int    `json:"id" form:"id"`
	FormUserID     int    `json:"form_user_id" form:"form_user_id"`
	PermissionID   int    `json:"permission_id" form:"permission_id"`
	PermissionName string `json:"permission_name" form:"permission_name"`
	Status         bool   `json:"status" form:"status"`
}

type UserOrgPermissionRes struct {
	ID                 int    `json:"id" form:"id"`
	UserOrganizationID int    `json:"user_organization_id" form:"user_organization_id"`
	PermissionID       int    `json:"permission_id" form:"permission_id"`
	PermissionName     string `json:"permission_name" form:"permission_name"`
	IsChecked          bool   `json:"is_checked" form:"is_checked"`
}

type GlobalAdmin struct {
	IsOwner              bool                   `json:"is_owner"`
	UserOrgPermissionRes []UserOrgPermissionRes `json:"user_organization_permissions"`
}

type FormPerformance struct {
	ID               int               `json:"id" form:"id" `
	Name             string            `json:"name" form:"name" binding:"required"`
	PeriodStartDate  string            `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate    string            `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic       string            `json:"profile_pic" form:"profile_pic"`
	LastUpdate       string            `json:"last_update" form:"last_update"`
	TotalRespon      int               `json:"total_respon" form:"total_respon" gorm:"default:0"`
	AllTotalRespon   string            `json:"all_total_respon" form:"all_total_respon" gorm:"default:0"`
	TotalTarget      int               `json:"total_target" form:"total_target" gorm:"default:0"`
	TotalAverage     float64           `json:"total_average" form:"total_average" gorm:"default:0"`
	TotalUpdateData  int               `json:"total_update_data" form:"total_update_data" gorm:"default:0"`
	TotalDeletedData int               `json:"total_deleted_data" form:"total_deleted_data" gorm:"default:null"`
	InfoTarget       string            `json:"info_target" form:"info_target"`
	InfoAverage      string            `json:"info_average" form:"info_average"`
	InfoTargetColor  string            `json:"info_target_color" form:"info_target_color"`
	InfoAverageColor string            `json:"info_average_color" form:"info_average_color"`
	ProgressData     []ProcessData     `json:"process_data"`
	HoursData        []GraficDataHours `json:"hours_data"`
}

type ProcessData struct {
	ProcessName string `json:"process_name"`
	Status      bool   `json:"status"`
}

// type SubmissionForm struct {
// 	FormID     int    `json:"form_id" form:"form_id" binding:"required"`
// 	FormName   string `json:"form_name" form:"form_name" `
// 	ProfilePic string `json:"profile_pic" form:"profile_pic"`
// 	DescStatus string `json:"desc_status" form:"desc_status"`
// 	SendStatus bool   `json:"send_status" form:"send_status"`
// }

type SubmissionForm struct {
	FormID          int    `json:"form_id" form:"form_id" binding:"required"`
	FormName        string `json:"form_name" form:"form_name" `
	ProfilePic      string `json:"profile_pic" form:"profile_pic"`
	TotalSubmission int    `json:"total_submission" form:"total_submission"`
}

type UserGetFormList struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	FormID      int    `json:"form_id"`
	Name        string `json:"name_id"`
	Description string `json:"description"`
	TotalRespon int    `json:"total_respon"`
	FormStatus  bool   `json:"form_status"`
}

type Paging struct {
	Page   int
	Limit  int
	SortBy string
	Sort   string
}

type FormDuplicate struct {
	FormID    int `json:"form_id"`
	ProjectID int `json:"project_id"`
}

type FormDataMap struct {
	FormID         int             `json:"form_id" binding:"required"`
	CompanyID      int             `json:"company_id"`
	StartDate      string          `json:"start_date"`
	EndDate        string          `json:"end_date"`
	FilterFieldIDs []FilterFieldID `json:"filter_field_ids" form:"filter_field_ids" `
}

type FilterFieldID struct {
	FieldID int `json:"field_id"`
}

type AppListForms struct {
	ID                   int    `json:"id" form:"id" `
	FormStatusID         int    `json:"form_status_id" form:"form_status_id"`
	FormStatus           string `json:"form_status" form:"form_status"`
	Name                 string `json:"name" form:"name" binding:"required"`
	Description          string `json:"description" form:"description"`
	Notes                string `json:"notes" form:"notes"`
	PeriodStartDate      string `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate        string `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic           string `json:"profile_pic" form:"profile_pic"`
	IsAttendanceRequired bool   `json:"is_attendance_required" form:"is_attendance_required"`
	IsAttendanceRadius   bool   `json:"is_attendance_radius" form:"is_attendance_radius"`
	AttendanceOverdate   bool   `json:"attendance_over_date"`
	AttendanceIn         string `json:"attendance_in" form:"attendance_in"`
	AttendanceOut        string `json:"attendance_out" form:"attendance_out"`
	IsButton             string `json:"is_button"`
	TotalResponden       int    `json:"total_responden"  gorm:"default:0"`
	TotalRespon          int    `json:"total_respon"  gorm:"default:0"`
	IsActivePeriod       bool   `json:"is_active_period"`
}

type FormDetail struct {
	ProjectID                int          `json:"group_id" form:"id" `
	ID                       int          `json:"id" form:"id" `
	FormStatusID             int          `json:"form_status_id" form:"form_status_id"`
	FormStatus               string       `json:"form_status" form:"form_status"`
	Name                     string       `json:"name" form:"name" binding:"required"`
	Description              string       `json:"description" form:"description"`
	Notes                    string       `json:"notes" form:"notes"`
	PeriodStartDate          string       `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate            string       `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic               string       `json:"profile_pic" form:"profile_pic"`
	UserID                   int          `json:"user_id,omitempty" form:"user_id"`
	TotalForms               int          `json:"total_forms" `
	TotalResponden           int          `json:"total_responden"  gorm:"default:0"`
	TotalRespondenActive     int          `json:"total_responden_active"  gorm:"default:0"`
	TotalRespon              int          `json:"total_respon"  gorm:"default:0"`
	TotalAdmin               int          `json:"total_admin"  gorm:"default:0"`
	TotalPerformance         int          `json:"total_performance"  gorm:"default:0"`
	TotalPerformanceFloat    float64      `json:"total_performance_float"  `
	CreatedBy                int          `json:"created_by" form:"created_by"`
	CreatedByName            string       `json:"created_by_name" form:"created_by_name"`
	CreatedByEmail           string       `json:"created_by_email" form:"created_by_email"`
	UpdatedByName            string       `json:"updated_by_name" form:"updated_by_name"`
	LastUpdate               string       `json:"last_updated" form:"last_updated"`
	AttendanceIn             string       `json:"attendance_in" form:"attendance_in"`
	AttendanceOut            string       `json:"attendance_out" form:"attendance_out"`
	SubmissionTarget         int          `json:"submission_target" form:"submission_target"`
	IsAttendanceRequired     bool         `json:"is_attendance_required" form:"is_attendance_required"`
	PeriodeRange             int          `json:"periode_range" form:"periode_range"`
	ArchivedAt               time.Time    `json:"archive_at" form:"archive_at" gorm:"default:null"`
	Author                   string       `json:"author" `
	AuthorName               string       `json:"author_name" `
	PerformanceColor         string       `json:"performance_color" `
	ShareUrl                 string       `json:"group_url" `
	LastSubmission           string       `json:"last_submission" `
	IsEditResponden          bool         `json:"is_edit_responden" `
	IsContainTagloc          bool         `json:"is_contain_tagloc" `
	AttendanceOverdate       bool         `json:"attendance_overdate"`
	AttendanceOverdateAt     time.Time    `json:"attendance_overdate_at" gorm:"default:null"`
	IsButton                 string       `json:"is_button"`
	OrganizationID           int          `json:"organization_id"`
	OrganizationName         string       `json:"organization_name"`
	IsShowFilterOrganization bool         `json:"is_show_filter_organization"`
	IsShowTabOrganization    bool         `json:"is_show_tab_organization"`
	IsAttendanceRadius       bool         `json:"is_attendance_radius"`
	FormFields               []FormFields `json:"form_fields"`
}

type FormOtp struct {
	ID     int  `json:"id"`
	FormID int  `json:"form_id"`
	UserID int  `json:"user_id"`
	Status bool `json:"status"`
}

type SubmissionFormOtp struct {
	FormID  int    `json:"form_id" binding:"required"`
	OtpCode string `json:"otp_code" binding:"required"`
}

type ShareCodeResp struct {
	FormID           int  `json:"form_id"`
	IsFirstConnected bool `json:"is_firsttime_connected"`
}

type FormTemplate struct {
	FormID      int    `json:"template_id" binding:"required"`
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfilePic  string `json:"profile_pic"`
}

type FormTemplateNew struct {
	FormID      int    `json:"form_id"`
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfilePic  string `json:"profile_pic"`
	Total       int    `json:"total"`
}

type FormCompany struct {
	FormID         int `json:"form_id"`
	OrganizationID int `json:"organization_id"`
}

type DeleteAdminEksObj struct {
	ID                        int `json:"id"`
	UserID                    int `json:"user_id"`
	OrganizationID            int `json:"organization_id"`
	UserOrganizationID        int `json:"user_organization_id"`
	UserOrganizationRolesID   int `json:"user_organization_roles_id"`
	FormOrganizationInvitesID int `json:"form_organization_invites_id"`
}

type FormToCompanyList struct {
	ID                        int    `json:"id"`
	FormID                    int    `json:"form_id"`
	OrganizationID            int    `json:"organization_id"`
	IsQuotaSharing            bool   `json:"is_quota_sharing"`
	OrganizationName          string `json:"organization_name"`
	OrganizationContactName   string `json:"organization_contact_name"`
	OrganizationContactPhone  string `json:"organization_contact_phone"`
	OrganizationProfilePic    string `json:"organization_profile"`
	OrganizationContactStatus string `json:"organization_contact_status"`
}

type FormSharing struct {
	ID                     int       `json:"id" form:"id" `
	FormStatusID           int       `json:"form_status_id" form:"form_status_id"`
	FormStatus             string    `json:"form_status" form:"form_status"`
	Name                   string    `json:"name" form:"name" binding:"required"`
	Description            string    `json:"description" form:"description"`
	Notes                  string    `json:"notes" form:"notes"`
	PeriodStartDate        string    `json:"period_start_date" form:"period_start_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	PeriodEndDate          string    `json:"period_end_date" form:"period_end_date" gorm:"default:null" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	ProfilePic             string    `json:"profile_pic" form:"profile_pic"`
	UserID                 int       `json:"user_id,omitempty" form:"user_id"`
	TotalForms             int       `json:"total_forms" `
	TotalResponden         int       `json:"total_responden"  gorm:"default:0"`
	TotalRespondenActive   int       `json:"total_responden_active"  gorm:"default:0"`
	TotalRespon            int       `json:"total_respon"  gorm:"default:0"`
	TotalAdmin             int       `json:"total_admin"  gorm:"default:0"`
	TotalPerformance       int       `json:"total_performance"  gorm:"default:0"`
	TotalPerformanceFloat  float64   `json:"total_performance_float"  `
	CreatedBy              int       `json:"created_by" form:"created_by"`
	CreatedByName          string    `json:"created_by_name" form:"created_by_name"`
	CreatedByEmail         string    `json:"created_by_email" form:"created_by_email"`
	UpdatedByName          string    `json:"updated_by_name" form:"updated_by_name"`
	LastUpdate             string    `json:"last_updated" form:"last_updated"`
	AttendanceIn           string    `json:"attendance_in" form:"attendance_in"`
	AttendanceOut          string    `json:"attendance_out" form:"attendance_out"`
	SubmissionTarget       int       `json:"submission_target" form:"submission_target"`
	IsAttendanceRequired   bool      `json:"is_attendance_required" form:"is_attendance_required"`
	PeriodeRange           int       `json:"periode_range" form:"periode_range"`
	ArchivedAt             time.Time `json:"archive_at" form:"archive_at" gorm:"default:null"`
	Author                 string    `json:"author" `
	AuthorName             string    `json:"author_name" `
	PerformanceColor       string    `json:"performance_color" `
	ShareUrl               string    `json:"group_url" `
	LastSubmission         string    `json:"last_submission" `
	IsEditResponden        bool      `json:"is_edit_responden" `
	IsContainTagloc        bool      `json:"is_contain_tagloc" `
	AttendanceOverdate     bool      `json:"attendance_overdate"`
	AttendanceOverdateAt   time.Time `json:"attendance_overdate_at" gorm:"default:null"`
	SharingSaldo           string    `json:"sharing_saldo"`
	StatusAdmin            string    `json:"status_admin"`
	OrganizationName       string    `json:"organization_name"`
	OrganizationProfilePic string    `json:"organization_profile_pic"`
}

type FormDataDelete struct {
	FormID      int           `json:"form_id" binding:"required"`
	FormDataIDs []FormDataIDs `json:"form_data_ids" `
}

type FormDataTotal struct {
	TotalForm            int  `json:"total_form" `
	TotalRespondenActive int  `json:"total_responden_active"  gorm:"default:0"`
	Submission           int  `json:"submission" form:"submission_target"`
	IsProfileComplete    bool `json:"is_profile_complete" `
}

type FormDataIDs struct {
	ID int
}

type InputFormUserOrganizations struct {
	ID               int
	FormID           int
	UserID           int
	OrganizationID   int
	FormUserStatusID int
	Type             string
}

type FormUserOrganizations struct {
	ID             int
	OrganizationID int
	FormUserID     int
}

type FilterFormCompanyList struct {
	ID               int    `json:"id"`
	OrganizationName string `json:"organization_name"`
}

type FormOrganizations struct {
	ID                     int
	FormID                 int
	OrganizationID         int
	OrganizationName       string
	OrganizationProfilePic string
}

type ObjectInputFormUsers struct {
	gorm.Model
	ID            int
	UserID        int
	UserName      string
	UserPhone     string
	Avatar        string
	Organizations string `json:"organization_name"`
	SubmitDate    string
}

type HistoryBalanceSaldo struct {
	FormID                 int    `json:"form_id"`
	OrganizationID         int    `json:"organization_id"`
	FormName               string `json:"form_name"`
	FormImage              string `json:"form_image"`
	OrganizationName       string `json:"organization_name"`
	TotalDeletedSubmission int    `json:"total_deleted_submission"`
	TotalUpdatedSubmission int    `json:"total_updated_submission"`
	TotalSubmission        int    `json:"total_submission"`
	TotalBlast             int    `json:"total_blast"`
	TotalUsageBalance      int    `json:"total_usage_balance"`
}

type HistoryBalanceSaldoWithDate struct {
	ID                     int    `json:"id"`
	DateDB                 string `json:"date_db"`
	Date                   string `json:"date"`
	TopupSubmission        int    `json:"topup_submission"`
	TotalDeletedSubmission int    `json:"total_deleted_submission"`
	TotalUpdatedSubmission int    `json:"total_updated_submission"`
	TotalSubmission        int    `json:"total_submission"`
	TotalBlast             int    `json:"total_blast"`
	TotalUsageBalance      int    `json:"total_usage_balance"`
}

type HistoryBalanceSaldoPerForm struct {
	FormID                 int               `json:"form_id"`
	OrganizationID         int               `json:"organization_id"`
	FormName               string            `json:"form_name"`
	OrganizationName       string            `json:"organization_name"`
	TotalDeletedSubmission int               `json:"total_deleted_submission"`
	TotalUpdatedSubmission int               `json:"total_updated_submission"`
	TotalSubmission        int               `json:"total_submission"`
	TotalUsageBalance      int               `json:"total_usage_balance"`
	DataUser               []ReportResponden `json:"data_user"`
}

type BlastInfoData struct {
	UserID    int `json:"user_id"`
	FormID    int `json:"form_id"`
	CreatedAt time.Time
}

type ObjectFormAttendanceLocations struct {
	ID         int     `json:"id"`
	FormID     int     `json:"form_id"`
	Name       string  `json:"name"`
	Location   string  `json:"location"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	IsCheckIn  bool    `json:"is_check_in"`
	IsCheckOut bool    `json:"is_check_out"`
	Radius     int     `json:"radius"`
}

type ObjectFormAttendanceLocation struct {
	ID         int     `json:"id"`
	FormID     int     `json:"form_id"`
	Name       string  `json:"name"`
	Location   string  `json:"location"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	IsCheckIn  bool    `json:"is_check_in"`
	IsCheckOut bool    `json:"is_check_out"`
}

type FillingType struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type ExportOutput struct {
	TotalSubmission int
	totalResponden  int
	TotalForm       int
}

type TopupHistory struct {
	ID             int    `json:"id"`
	OrganizationID int    `json:"organization_id"`
	Quota          int    `json:"quota"`
	Price          string `json:"price"`
	PeriodDays     int    `json:"period_days"`
	CreatedAt      string `json:"created_at"`
}

type FormAll struct {
	gorm.Model
	ProjectID                int
	ID                       int
	FormStatusID             int `gorm:"default:1"`
	FormStatus               string
	Name                     string
	Description              string
	Notes                    string
	PeriodStartDate          string `gorm:"default:null"`
	PeriodEndDate            string `gorm:"default:null"`
	ProfilePic               string
	IsPublish                bool
	IsAttendanceRequired     bool
	EncryptCode              string
	CreatedBy                int `json:"createdBy" gorm:"default:null"`
	CreatedByName            string
	CreatedByEmail           string
	Type                     string
	SubmissionTargetUser     int
	Author                   string
	ShareUrl                 string
	UpdatedBy                int `json:"cpdatesBy" gorm:"default:null"`
	DeletedBy                int `json:"deletedBy" gorm:"default:null"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt `gorm:"index"`
	ArchivedAt               time.Time      `gorm:"default:null"`
	IsQuotaSharing           bool
	FormShared               string
	FormShareCount           int
	FormExternalCompanyName  string
	FormExternalCompanyImage string
}
