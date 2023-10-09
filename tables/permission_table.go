package tables

import (
	"time"

	"gorm.io/gorm"
)

type Permissions struct {
	ID        int
	Name      string
	Slug      string
	HttpPath  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type FormUserPermission struct {
	ID           int
	FormUserID   int
	PermissionID int
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type FormUserPermissionJoin struct {
	ID             int
	FormUserID     int
	PermissionID   int
	PermissionName string
	Status         bool
}

type UserOrganizationPermission struct {
	ID                 int
	UserOrganizationID int `json:"user_organization_id"`
	PermissionID       int
	IsChecked          bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type UserOrganizationPermissionJoin struct {
	ID                 int
	UserOrganizationID int
	PermissionID       int
	PermissionName     string
	IsChecked          bool
}
