package tables

import (
	"time"

	"gorm.io/gorm"
)

type Projects struct {
	ID          int
	Name        string
	Description string
	ParentID    int `gorm:"default:null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CreatedBy   int            `gorm:"default:null"`
	UpdatedBy   int            `gorm:"default:null"`
	DeletedBy   int            `gorm:"default:null"`
}

type ProjectForms struct {
	ID        int
	ProjectID int `form:"project_id" binding:"required"`
	FormID    int `form:"form_id" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProjectFormsJoin struct {
	ID             int
	ProjectID      int `form:"project_id" binding:"required"`
	FormStatusID   int
	FormStatus     string
	FormID         int `form:"form_id" binding:"required"`
	Name           string
	Description    string
	ProfilePic     string
	CreatedBy      int
	CreatedByName  string
	CreatedByEmail string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type ProjectOrganizations struct {
	ID             int
	ProjectID      int
	OrganizationID int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
