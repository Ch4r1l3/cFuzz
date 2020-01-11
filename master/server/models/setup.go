package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var DB *gorm.DB

func InsertObject(obj interface{}) error {
	return DB.Create(obj).Error
}

func GetObjects(objs interface{}) error {
	return DB.Find(objs).Error
}

func GetObjectByID(obj interface{}, id uint64) error {
	return DB.Where("id = ?", id).First(obj).Error
}

func DeleteObjectByID(obj interface{}, id uint64) error {
	return DB.Where("id = ?", id).Delete(obj).Error
}

func IsObjectExistsByID(obj interface{}, id uint64) bool {
	return !DB.Where("id = ?", id).First(obj).RecordNotFound()
}

func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", "./master.db")
	if err != nil {
		log.Fatal(err)
	}
	DB.Exec("PRAGMA foreign_keys = ON")
	DB.SingularTable(true)
	DB.AutoMigrate(&Deployment{}, &Task{}, &TaskTarget{}, &TaskCorpus{}, &TaskEnvironment{}, &TaskArgument{}, &TaskCrash{}, &TaskFuzzResult{}, &TaskFuzzResultStat{}, &Fuzzer{})
}
