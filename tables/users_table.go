package tables

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID              int
	Name            string
	Tipe            string
	Phone           string
	Email           string
	Password        string
	RememberToken   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Deleted         string
	Apps            string
	RoleID          int
	EncryptCode     string
	IsEmailVerified bool `gorm:"default:false"`
}

type UsersMA struct {
	ID              int
	Name            string
	Phone           string `gorm:"type:varchar(13);default:null" json:"phone"`
	Email           string `gorm:"type:varchar(255);default:null" json:"email"`
	Password        string `gorm:"default:null"`
	Avatar          string
	DateOfBirth     string `sql:"default:null" gorm:"default:null"`
	GenderID        int    `gorm:"default:null"`
	RememberToken   string
	IsEmailVerified bool `gorm:"default:false"`
	isPhoneVerified bool `gorm:"default:false"`
	Status          bool `gorm:"default:true"`
	EncryptCode     string
	RoleID          int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserData struct {
	ID              int
	Name            string
	Phone           string `gorm:"type:varchar(13);default:null" json:"phone"`
	Email           string `gorm:"type:varchar(255);default:null" json:"email"`
	Password        string
	Avatar          string
	DateOfBirth     string `sql:"default:CURRENT_TIMESTAMP"`
	GenderID        int
	GenderName      string
	RememberToken   string
	IsEmailVerified bool `gorm:"default:false"`
	isPhoneVerified bool `gorm:"default:false"`
	Status          bool `gorm:"default:true"`
	EncryptCode     string
	RoleID          int
	CompanyName     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
