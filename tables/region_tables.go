package tables

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	gorm.Model
	ID        int
	Name      string
	Status    bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Cities struct {
	gorm.Model
	ID         int
	ProvinceID int
	Name       string
	Status     bool `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Districts struct {
	gorm.Model
	ID        int
	CityID    int
	Name      string
	Status    bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SubDistricts struct {
	gorm.Model
	ID         int
	DistrictID int
	Name       string
	PostalCode int
	Status     bool `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
