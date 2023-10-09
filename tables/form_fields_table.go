package tables

import (
	"time"

	"gorm.io/gorm"
)

type FormFields struct {
	gorm.Model
	ID                 int
	ParentID           int `gorm:"default:null"`
	FormID             int
	FieldTypeID        int `gorm:"default:null"`
	Label              string
	Description        string
	Option             string
	ConditionType      string `gorm:"default:null"`
	UpperlowerCaseType string `gorm:"default:null"`
	IsMultiple         bool   `gorm:"default:false"`
	IsRequired         bool   `gorm:"default:false"`
	IsSection          bool   `gorm:"default:false"`
	IsCountryPhoneCode bool   `gorm:"default:false"`
	SectionColor       string
	SortOrder          int
	IsCondition        bool `gorm:"default:false"`
	TagLocIcon         string
	TagLocColor        string
	// ProvinceID         int    `gorm:"default:null"`
	// CityID             int    `gorm:"default:null"`
	// DistrictID         int    `gorm:"default:null"`
	// SubDistrictID      int    `gorm:"default:null"`
	// AddressType  string `gorm:"default:null"`
	// CurrencyType int    `gorm:"default:null"`
	// Currency     string `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type FormFieldPics struct {
	ID          int
	FormFieldID int
	Pic         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
