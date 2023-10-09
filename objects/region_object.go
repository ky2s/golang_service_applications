package objects

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"omitempty"`
	UpdatedAt time.Time      `json:"omitempty"`
	DeletedAt gorm.DeletedAt `json:"omitempty" gorm:"index"`
}

type Cities struct {
	ID         int            `json:"id"`
	ProvinceID int            `json:"province_id"`
	Name       string         `json:"name"`
	Status     bool           `json:"status"`
	CreatedAt  time.Time      `json:"omitempty"`
	UpdatedAt  time.Time      `json:"omitempty"`
	DeletedAt  gorm.DeletedAt `json:"omitempty" gorm:"index"`
}

type Districts struct {
	ID        int            `json:"id"`
	CityID    int            `json:"city_id"`
	Name      string         `json:"name"`
	Status    bool           `json:"status"`
	CreatedAt time.Time      `json:"omitempty"`
	UpdatedAt time.Time      `json:"omitempty"`
	DeletedAt gorm.DeletedAt `json:"omitempty" gorm:"index"`
}

type SubDistricts struct {
	ID         int            `json:"id"`
	DistrictID int            `json:"district_id"`
	Name       string         `json:"name"`
	PostalCode int            `json:"postal_code"`
	Status     bool           `json:"status"`
	CreatedAt  time.Time      `json:"omitempty"`
	UpdatedAt  time.Time      `json:"omitempty"`
	DeletedAt  gorm.DeletedAt `json:"omitempty" gorm:"index"`
}

type Radius struct {
	LocationID int     `json:"location_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Distance   float64 `json:"distance"`
	IsRadius   bool    `json:"is_radius"`
}

type RadiusSlow struct {
	Distance float64 `json:"distance"`
	IsRadius bool    `json:"is_radius"`
}
