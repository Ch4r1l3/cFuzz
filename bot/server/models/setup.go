package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var DB *gorm.DB

func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", "./bot.db")
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&Fuzzer{}, &Task{}, &TaskArgument{}, &TaskEnvironment{})
}
