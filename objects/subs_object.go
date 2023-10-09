package objects

import (
	"time"

	"gorm.io/gorm"
)

type SubsPlan struct {
	ID                 int
	OrganizationID     int
	SubscriptionPlanID int
	RespondentCurrent  int `gorm:"default:0"`
	CreatedBy          int `gorm:"default:null"`
	UpdatedBy          int `gorm:"default:null"`
	DeletedBy          int `gorm:"default:null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type SubsPlanRes struct {
	SubsPlanID         int `json:"subs_plan_id"`
	SubsPlanTotalValue int `json:"subs_plan_total"`
	SubsPlanCurrent    int `json:"subs_plan_current"`
}

type InjuryPlan struct {
	ID                             int
	OrganizationSubscriptionPlanID int
	Quota                          int
}

type SubsPlanPeriod struct {
	ID                 int
	OrganizationID     int
	SubscriptionPlanID int
	QuotaCurrent       int
	TotalPeriodDays    int
	PeriodEndDate      string
	PeriodRemain       int
}
