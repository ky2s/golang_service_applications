package tables

import (
	"time"

	"gorm.io/gorm"
)

type Settings struct {
	ID        int
	AppTypeID int
	Code      string
	Name      string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
