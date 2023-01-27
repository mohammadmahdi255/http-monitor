package tables

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Urls     []URL  `gorm:"foreignkey:user_id"`
}

// HashPassword generates a hashed string from 'pass'
// returns error if 'pass' is empty
func HashPassword(pass string) (string, error) {

	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hash), err
}
