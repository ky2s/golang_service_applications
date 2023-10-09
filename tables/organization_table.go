package tables

import (
	"time"

	"gorm.io/gorm"
)

type Organizations struct {
	ID          int
	Code        string
	Name        string
	Description string
	Address1    string
	Phone1      string
	Phone2      string
	Email       string
	IsDefault   bool
	ProfilePic  string
	CreatedBy   int `gorm:"default:null"`
	UpdatedBy   int `gorm:"default:null"`
	DeletedBy   int `gorm:"default:null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type UserOrganizations struct {
	ID             int
	UserID         int
	OrganizationID int
	IsDefault      bool
	CreatedBy      int
	UpdatedBy      int
	DeletedBy      int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// type UserOrganizations_ struct {
// 	ID              int
// 	UserID          int
// 	Phone           string
// 	UserEmail       string
// 	OrganizationID  int
// 	Code            string
// 	Name            string `json:"name"`
// 	Description     string
// 	Address1        string
// 	Phone1          string
// 	Phone2          string
// 	Email           string
// 	RoleID          int
// 	ProfilePic      string
// 	ContactName     string
// 	ContactPhone    string
// 	ContactPosition string
// 	IsDefault       bool
// 	CreatedBy       int
// 	UpdatedBy       int
// 	DeletedBy       int
// }

type UserOrganizationRoles struct {
	ID                 int
	UserOrganizationID int
	RoleID             int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type SelectUserOrganizations struct {
	ID              int
	UserID          int
	Phone           string
	UserEmail       string
	OrganizationID  int
	Code            string
	Name            string `json:"name"`
	Description     string
	Address1        string
	Phone1          string
	Phone2          string
	Email           string
	ProfilePic      string
	ContactName     string
	ContactPhone    string
	ContactPosition string
}
