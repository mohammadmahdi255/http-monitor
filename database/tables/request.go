package tables

import (
	"time"
)

type Request struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UrlId     uint
	Result    int
}
