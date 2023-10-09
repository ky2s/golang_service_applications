package tables

import (
	"time"

	"gorm.io/gorm"
)

type InputForms struct {
	ID           int
	UserID       int
	UserName     string
	Phone        string
	Avatar       string
	Address      string
	Latitude     float64
	Longitude    float64
	UpdatedCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	DeletedBy    int
}

type InputFormCustomAnswers struct {
	ID           int
	FormID       int
	FormFieldID  int
	InputFormID  int
	CustomAnswer string
}

type InputFormOrganizations struct {
	ID             int
	OrganizationID int
	FormID         int
	InputFormID    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type InputFormOrganization struct {
	ID             int
	OrganizationID int
	FormID         int
	InputFormID    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type InputFormJoinOrganizations struct {
	ID               int
	UserID           int
	UserName         string
	Phone            string
	Avatar           string
	Address          string
	Latitude         float64
	Longitude        float64
	OrganizationID   int
	OrganizationName string
	UpdatedCount     int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
