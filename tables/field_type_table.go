package tables

import (
	"time"

	"gorm.io/gorm"
)

type FieldTypes struct {
	ID                int `sql:"type:int(11);primary key"`
	Code              string
	RealVarType       string
	IsGroup           bool
	IsMedia           bool
	IsChildField      bool
	NameTextcontentId int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

type FieldTypeTrans struct {
	ID                int `sql:"type:int(11);primary key"`
	Code              string
	Name              string
	RealVarType       string
	IsGroup           bool
	IsMedia           bool
	IsChildField      bool
	NameTextcontentId int
	Info              string
}
