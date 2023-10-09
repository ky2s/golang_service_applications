package tables

import (
	"time"

	"gorm.io/gorm"
)

type ConditionRules struct {
	gorm.Model
	ID                int
	Name              string
	Type              string
	Code              string
	NameTextcontentID int
}

type Translations struct {
	gorm.Model
	ID            int
	TextContentID int
	LanguageID    int
	Translations  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type FormFieldConditionRules struct {
	gorm.Model
	ID                     int
	FormFieldID            int
	ConditionRuleID        int
	Value1                 string
	Value2                 string
	ErrMsg                 string `gorm:"default:null"`
	ConditionParentFieldID int    `gorm:"default:null"`
	ConditionAllRight      bool   `gorm:"default:null"`
	TabMaxOnePerLine       bool   `gorm:"default:null"`
	TabEachLineRequire     bool   `gorm:"default:null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type SelectFormFieldConditionRules struct {
	gorm.Model
	ID                     int
	FormFieldID            int
	ParentID               int `gorm:"default:null"`
	FormID                 int
	FieldTypeID            int
	FieldTypeName          string
	Label                  string
	Description            string
	Option                 string
	ConditionType          string `gorm:"default:null"`
	UpperlowerCaseType     string `gorm:"default:null"`
	IsMultiple             bool
	IsRequired             bool
	IsSection              bool
	SectionColor           string
	ConditionRuleID        int
	Value1                 string
	Value2                 string
	ErrMsg                 string
	ConditionParentFieldID int
	ConditionAllRight      bool `gorm:"default:null"`
	TabMaxOnePerLine       bool `gorm:"default:null"`
	TabEachLineRequire     bool `gorm:"default:null"`
	IsCondition            bool
	IsCountryPhoneCode     bool
	SortOrder              int
	Image                  string
	TagLocColor            string
	TagLocIcon             string
	ProvinceID             int    `gorm:"default:null"`
	CityID                 int    `gorm:"default:null"`
	DistrictID             int    `gorm:"default:null"`
	SubDistrictID          int    `gorm:"default:null"`
	AddressType            string `gorm:"default:null"`
	CurrencyType           int    `gorm:"default:null"`
	Currency               string `gorm:"default:null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}
