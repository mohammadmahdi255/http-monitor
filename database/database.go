package database

import (
	"fmt"
	"github.com/mohammadmahdi255/http-monitor/database/tables"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strings"
)

// New Setup initializes a database
func New(databaseName string) *gorm.DB {
	if !strings.HasSuffix(databaseName, ".db") {
		databaseName = databaseName + ".db"
	}

	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		fmt.Println("Error in creating Store file : ", err)
		return nil
	}

	err = db.AutoMigrate(&tables.User{}, &tables.Request{}, &tables.URL{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
