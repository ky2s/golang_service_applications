package tables

import (
	"time"

	"gorm.io/gorm"
)

type Forms struct {
	gorm.Model
	ID                   int
	FormStatusID         int `gorm:"default:1"`
	Name                 string
	Description          string
	Notes                string
	PeriodStartDate      string `gorm:"default:null"`
	PeriodEndDate        string `gorm:"default:null"`
	SubmissionTargetUser int
	ProfilePic           string
	IsPublish            bool
	IsAttendanceRequired bool
	// OrganizationID       int
	EncryptCode          string
	ShareUrl             string
	AttendanceOverdateAt time.Time `gorm:"default:null"`
	CreatedBy            int       `gorm:"default:null"`
	UpdatedBy            int       `gorm:"default:null"`
	DeletedBy            int       `gorm:"default:null"`
	IsAttendanceRadius   bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"default:null"`
	ArchivedAt           time.Time      `gorm:"default:null"`
}

type FormOut struct {
	gorm.Model
	ID                   int
	FormStatusID         int    `gorm:"default:1"`
	FormStatus           string `json:"omitempty"`
	Name                 string
	Description          string
	Notes                string
	PeriodStartDate      string `gorm:"default:null"`
	PeriodEndDate        string `gorm:"default:null"`
	SubmissionTargetUser int
	ProfilePic           string
	IsPublish            bool
	IsAttendanceRequired bool
	EncryptCode          string
	ShareUrl             string
	AttendanceOverdateAt time.Time `gorm:"default:null"`
	UserID               int
	UserPhone            string
	OrganizationID       int
	OrganizationName     string
	CreatedBy            int `json:"createdBy" gorm:"default:null"`
	UpdatedBy            int `json:"cpdatesBy" gorm:"default:null"`
	DeletedBy            int `json:"deletedBy" gorm:"default:null"`
	IsAttendanceRadius   bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	ArchivedAt           time.Time      `gorm:"default:null"`
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
	AccessType               string
	FormShared               string
	FormShareCount           int
	FormExternalCompanyName  string
	FormExternalCompanyImage string
	IsAttendanceRadius       bool
}

type FormOrganizationsJoin struct {
	gorm.Model
	ID                   int
	FormStatusID         int `gorm:"default:1"`
	Name                 string
	Description          string
	Notes                string
	PeriodStartDate      string `gorm:"default:null"`
	PeriodEndDate        string `gorm:"default:null"`
	SubmissionTargetUser int
	ProfilePic           string
	IsPublish            bool
	IsAttendanceRequired bool
	EncryptCode          string
	Type                 string `json:"type"`
	OrganizationID       int
	CreatedBy            int `json:"createdBy" gorm:"default:null"`
	UpdatedBy            int `json:"cpdatesBy" gorm:"default:null"`
	DeletedBy            int `json:"deletedBy" gorm:"default:null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	ArchivedAt           time.Time      `gorm:"default:null"`
}

type FormUsers struct {
	gorm.Model
	ID               int `sql:"type:int(11);primary key"`
	FormID           int
	UserID           int
	FormUserStatusID int    `gorm:"default:1"`
	Type             string `gorm:"default:respondent"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type JoinFormUsers struct {
	gorm.Model
	ID                   int `sql:"type:int(11);primary key"`
	FormID               int
	UserID               int
	UserName             string
	UserImage            string
	Email                string
	Phone                string
	Name                 string
	Description          string
	Notes                string
	ProfilePic           string
	PeriodStartDate      string
	PeriodEndDate        string
	FormStatusID         int
	FormStatus           string
	FormUserStatusID     int
	FormUserStatusName   string
	Type                 string
	AttendanceIn         string
	AttendanceOut        string
	IsAttendanceRequired bool
	IsAttendanceRadius   bool
	IsAttendanceOverdate bool
	OrganizationID       int
	CreatedAt            time.Time
}

type FormFieldTempAssets struct {
	ID          int
	FormID      int
	FormFieldID int
	TempAsset   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type InputFormUsers struct {
	gorm.Model
	ID        int
	UserID    int
	UserName  string
	UserPhone string
	Avatar    string
}

type Date struct {
	Date string
}

type Months struct {
	Month string
}

type TotalDate struct {
	Date  string
	Value string
}

type FormFieldSection struct {
	ID          int
	Name        string
	Description string
	Color       string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time      `gorm:"default:null"`
	DeletedAt   gorm.DeletedAt `gorm:"default:null"`
}

type FormOrganizations struct {
	ID             int
	FormID         int
	OrganizationID int
	CreatedAt      time.Time
	UpdatedAt      time.Time      `gorm:"default:null"`
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`
}

type FormPeriodRange struct {
	gorm.Model
	PeriodRange int `json:"period_range"`
}

type AddLog struct {
	ID           int
	UserID       int
	PermissionID int
	FormID       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FormOtps struct {
	ID        int
	FormID    int
	UserID    int
	OtpCode   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type UserFormOrganizations struct {
	ID             int
	FormID         int
	OrganizationID int
	UserID         int
	CreatedAt      time.Time
	UpdatedAt      time.Time      `gorm:"default:null"`
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`
}

type UserOrganizationInvites struct {
	ID             int
	UserID         int
	OrganizationID int
	IsQuotaSharing bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type FormOrganizationInvites struct {
	ID             int
	FormID         int
	OrganizationID int
	IsQuotaSharing bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type UserOrganizationInviteDetail struct {
	ID                       int
	UserID                   int
	OrganizationID           int
	IsQuotaSharing           bool
	OrganizationName         string
	OrganizationContactName  string
	OrganizationContactPhone string
	OrganizationProfilePic   string
}

type JoinFormCompanies struct {
	ID                       int    `sql:"type:int(11);primary key" json:"id"`
	FormID                   int    `json:"form_id"`
	OrganizationID           int    `json:"organization_id"`
	IsQuotaSharing           bool   `json:"is_quota_sharing"`
	OrganizationName         string `json:"organization_name"`
	Type                     string `json:"type"`
	OrganizationContactName  string `json:"organization_contact_name"`
	OrganizationContactPhone string `json:"organization_contact_phone"`
	OrganizationProfilePic   string `json:"organization_profile_pic"`
}

type FormUserOrganizations struct {
	ID             int
	FormUserID     int
	OrganizationID int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type FormAttendanceLocations struct {
	ID         int
	FormID     int
	Name       string
	Location   string
	Geometry   string
	IsCheckIn  bool
	IsCheckOut bool
	Radius     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
