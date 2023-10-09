package objects

import (
	"time"

	"gorm.io/gorm"
)

type Organizations struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Address1    string `json:"address"`
	Phone1      string `json:"phone"`
	Email       string `json:"email"`
	CreatedBy   int    `json:"created_by"`
	UpdatedBy   int    `json:"updated_by"`
	DeletedBy   int    `json:"deleted_by"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type UserOrganizations struct {
	ID             int
	UserID         int
	Phone          string
	UserEmail      string
	OrganizationID int
	Code           string
	Name           string `json:"name"`
	Description    string
	Address1       string
	Phone1         string
	Phone2         string
	Email          string
	CreatedBy      int
	UpdatedBy      int
	DeletedBy      int
	IsDefault      bool
}

type UserOrganizationRoles struct {
	ID                 int    `json:"id"`
	UserOrganizationID int    `json:"user_organization_id"`
	OrganizationID     int    `json:"organization_id"`
	OrganizationName   string `json:"organization_name"`
	RoleID             int    `json:"role_id"`
	RoleName           string `json:"role_name"`
	UserID             int    `json:"user_id"`
	UserName           string `json:"user_name"`
}
