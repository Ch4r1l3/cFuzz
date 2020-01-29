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
	return DB.Order("id").Find(objs).Error
}

func GetCount(objs interface{}) (int, error) {
	var count int
	err := DB.Find(objs).Count(&count).Error
	return count, err
}

func GetObjectsPagination(objs interface{}, offset int, limit int) error {
	return DB.Order("id").Offset(offset).Limit(limit).Find(objs).Error
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
	DB.AutoMigrate(&Deployment{}, &Task{}, &StorageItem{}, &TaskEnvironment{}, &TaskArgument{}, &TaskCrash{}, &TaskFuzzResult{}, &TaskFuzzResultStat{})
}
