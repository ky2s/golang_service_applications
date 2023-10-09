package objects

import (
	"time"

	"gorm.io/gorm"
)

type InputFormFields struct {
	ID            int              `json:"id"`
	Label         string           `json:"label" form:"label" `
	FieldTypeName string           `json:"field_type" form:"field_type" `
	TabHeader     []TabValueHeader `json:"tab_header"`
}
type TabValueHeader struct {
	Label string `json:"label"`
}

type TabDataRowHeader struct {
	TabDataHeader []TabDataHeader `json:"tab_header"`
}

type TabDataHeader struct {
	Value string `json:"value"`
}

type InputData struct {
	ID          int              `json:"id"`
	FieldID     string           `json:"field_id,omitempty"  `
	FieldTypeID int              `json:"field_type_id,omitempty"`
	Days        string           `json:"days"`
	Date        string           `json:"date"`
	Time        string           `json:"time"`
	Value       string           `json:"value"`
	TabValue    []TabValueAnswer `json:"tab_value"`
}

type TabValueAnswer struct {
	Answer string `json:"answer"`
}

type InputFieldData struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	UserName   string      `json:"user_name"`
	Phone      string      `json:"phone"`
	StatusData string      `json:"status_data"`
	Note       string      `json:"note"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
	DeletedAt  string      `json:"deleted_at"`
	InputData  []InputData `json:"input_data"`
}

type InputFormRes struct {
	FormID            int               `json:"form_id" `
	FormName          string            `json:"form_name" `
	FormDescription   string            `json:"form_description" `
	PeriodStartDate   string            `json:"period_start_date" `
	PeriodEndDate     string            `json:"period_end_date" `
	AuthorPhoneNumber string            `json:"author_phone_number" `
	CreatedAt         string            `json:"created_at" `
	UpdatedAt         string            `json:"updated_at" `
	FieldLabel        []InputFormFields `json:"field_label" `
	FieldData         []InputFieldData  `json:"field_data" `
}

type InputForms struct {
	ID        int       `json:"id" `
	UserID    int       `json:"user_i" `
	UserName  string    `json:"user_name" `
	Phone     string    `json:"phone" `
	CreatedAt time.Time `json:"created_at" `
	UpdatedAt time.Time `json:"updated_at" `
}

type Rows struct {
	ID            int         `json:"id,omitempty" `
	UserID        int         `json:"user_id" `
	UserName      string      `json:"user_name" `
	UserPhone     string      `json:"user_phone" `
	Organizations string      `json:"organization_name" `
	SubmitDate    string      `json:"submit_date" `
	InputData     []InputData `json:"input_data"`
}

type InputFormUsers struct {
	FormID      int               `json:"form_id,omitempty" `
	UserId      int               `json:"user_id,omitempty" `
	UserName    string            `json:"user_name,omitempty" `
	FieldHeader []InputFormFields `json:"field_header" `
	FieldData   []Rows            `json:"field_data" `
}

type Date struct {
	Date string `json:"date" `
}

type Hours struct {
	Hours string `json:"hours" `
}

type InputFormDetail struct {
	FormID           int                `json:"form_id"`
	FormName         string             `json:"form_name"`
	FormDescription  string             `json:"form_description"`
	PeriodStartDate  string             `json:"period_start_date"`
	PeriodEndDate    string             `json:"period_end_date"`
	FieldLabelOnData []FormField3Groups `json:"field_label_on_data"`
}

type LabelOnData struct {
	FieldTypeID   int    `json:"field_type_id" form:"field_type_id" `
	FieldTypeName string `json:"field_type" form:"field_type" `
	Label         string `json:"label" form:"label" `
	FieldData     string `json:"field_data"`
}

type ReportResponden struct {
	ID                    int     `json:"id,omitempty" `
	UserID                int     `json:"user_id" `
	UserName              string  `json:"user_name" `
	UserPhone             string  `json:"user_phone" `
	Submission            int     `json:"submission" `
	TotalActiveSubmission string  `json:"total_active_submission" `
	TotalUpdateData       int     `json:"total_update_data" `
	TotalDeletedData      int     `json:"total_deleted_data" `
	TotalAverage          float64 `json:"total_average" form:"total_average" gorm:"default:0"`
	SubmissionTargetUser  int     `json:"submission_target_user" `
	Performance           int     `json:"performance" `
	PerformanceFloat      float64 `json:"performance_float" `
	LastSubmission        string  `json:"last_submission" `
	LastSubmissionDate    string  `json:"last_submission_date" `
	SubmissionColor       string  `json:"submission_color" `
	PerformanceColor      string  `json:"performance_color" `
	Address               string  `json:"address" `
	Latitude              float64 `json:"latitude" `
	Longitude             float64 `json:"longitude" `
	Avatar                string  `json:"avatar" `
}

type HistorySaldoUse struct {
	ID                     int    `json:"id,omitempty" `
	UserID                 int    `json:"user_id" `
	UserName               string `json:"user_name" `
	UserPhone              string `json:"user_phone" `
	UserImage              string `json:"user_image" `
	TotalSubmission        int    `json:"total_submission" `
	TotalUpdatedSubmission int    `json:"total_updated_submission" `
	TotalDeletedSubmission int    `json:"total_deleted_submission" `
	TotalBlast             int    `json:"total_blast"`
	TotalUsageBalance      int    `json:"total_usage_balance" `
}

type InputFieldDataMap struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	UserName  string         `json:"user_name"`
	Phone     string         `json:"phone"`
	Avatar    string         `json:"avatar"`
	CreatedAt string         `json:"created_at"`
	InputData []InputDataMap `json:"input_data"`
}
type InputDataMap struct {
	FieldID     string  `json:"field_id,omitempty"  `
	FieldTypeID int     `json:"field_type_id,omitempty"`
	TagLoc      string  `json:"tag_loc"`
	TagLocColor string  `json:"tag_loc_color"`
	TagLocIcon  string  `json:"tag_loc_icon"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Value       string  `json:"value"`
}

type InputFormMapRes struct {
	FormID          int                 `json:"form_id" `
	FormName        string              `json:"form_name" `
	FormDescription string              `json:"form_description" `
	PeriodStartDate string              `json:"period_start_date" `
	PeriodEndDate   string              `json:"period_end_date" `
	TaglocScopes    TaglocScopes        `json:"tagloc_scopes" `
	FieldLabel      []InputFormFields   `json:"field_label" `
	FieldData       []InputFieldDataMap `json:"field_data" `
}

type TaglocScopes struct {
	Latitude1  float64 `json:"latitude_1" gorm:"default:0.0"`
	Longitude1 float64 `json:"longitude_1" gorm:"default:0.0"`
	Latitude2  float64 `json:"latitude_2" gorm:"default:0.0"`
	Longitude2 float64 `json:"longitude_2" gorm:"default:0.0"`
}

type InputFormOrganizations struct {
	ID               int
	OrganizationID   int
	OrganizationName string
	FormID           int
	InputFormID      int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type DataOption struct {
	Data []ValueOption `json:'data'`
}

type ValueOption struct {
	Value string `json:'value'`
}
