package objects

import (
	"time"

	"gorm.io/gorm"
)

type UserLanguges struct {
	ID         int
	UserID     int
	LanguageID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
