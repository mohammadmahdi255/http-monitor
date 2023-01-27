package tables

import (
	"time"
)

type URL struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UserId      uint   `gorm:"unique_index:index_addr_user"` // for preventing url duplication for a single user
	Address     string `gorm:"unique_index:index_addr_user"`
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"`
}
