package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
)

var DB *gorm.DB

func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", "./master.db")
	if err != nil {
		log.Fatal(err)
	}
	DB.Exec("PRAGMA foreign_keys = ON")
	DB.SingularTable(true)
	DB.AutoMigrate(&Image{}, &Task{}, &StorageItem{}, &TaskEnvironment{}, &TaskArgument{}, &TaskCrash{}, &TaskFuzzResult{}, &TaskFuzzResultStat{}, &User{})

	// check if admin exist, if not, create one
	var user User
	err = DB.Where("is_admin = true").First(&user).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Fatal(err)
	} else if err != nil {
		var username string
		fmt.Print("admin username: ")
		fmt.Scan(&username)
		fmt.Print("admin password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if len(password) < 6 {
			log.Fatal("password should longer than 6")
		}
		if len(password) > 18 {
			log.Fatal("password should shorter than 18")
		}
		salt, err := utils.RandomString(12)
		if err != nil {
			log.Fatal(err)
		}
		err = DB.Create(&User{
			Username: username,
			Password: utils.GetEncryptPassword(fmt.Sprintf("%x", sha256.Sum256(password)), salt),
			Salt:     salt,
			IsAdmin:  true,
		}).Error
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\ncreate success")
	}
}
