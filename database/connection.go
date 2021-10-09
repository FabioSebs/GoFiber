package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() {
	_, err := gorm.Open(mysql.Open("fab@/go_fiber"), &gorm.Config{})

	if err != nil {
		panic("could not connect to database")
	}
}
