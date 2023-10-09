package tables

import (
	"time"

	"gorm.io/gorm"
)

type SubsPlan struct {
	ID                 int
	OrganizationID     int
	SubscriptionPlanID int
	RespondentCurrent  int
	BlashCurrent       int
	QuotaCurrent       int
	QuotaTotal         int
	IsBlocked          bool
	TotalPeriodDays    int
	CreatedBy          int `gorm:"default:null"`
	UpdatedBy          int `gorm:"default:null"`
	DeletedBy          int `gorm:"default:null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type InjuryPlan struct {
	ID                             int
	OrganizationSubscriptionPlanID int
	Quota                          int
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	DeletedAt                      gorm.DeletedAt `gorm:"index"`
}
