package tables

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Urls      []URL  `gorm:"foreignkey:user_id"`
}
