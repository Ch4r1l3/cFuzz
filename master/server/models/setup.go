package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
)

var DB *gorm.DB

type DeleteAble interface {
	Delete() error
}

type UserIDAble interface {
	GetUserID() uint64
}

func InsertObject(obj interface{}) error {
	return DB.Create(obj).Error
}

func GetObjects(objs interface{}) error {
	return DB.Order("id").Find(objs).Error
}

func GetCount(objs interface{}) (int, error) {
	var count int
	err := DB.Model(objs).Count(&count).Error
	return count, err
}

func GetObjectCombinCustom(objs interface{}, offset int, limit int, name string, queries []string, values []interface{}) (int, error) {
	var count int
	var err error
	query := ""
	for _, v := range queries {
		query += v + " AND "
	}
	if name != "" {
		err = DB.Model(objs).Where(query+"name LIKE ?", append(values, "%"+name+"%")...).Count(&count).Error
	} else {
		err = DB.Model(objs).Where(query+"1=1", values...).Count(&count).Error
	}
	if err != nil {
		return 0, err
	}
	t := DB.Order("id")
	if name != "" {
		t = t.Where(query+"name LIKE ?", append(values, "%"+name+"%")...)
	} else {
		t = t.Where(query+"1=1", values...)
	}
	if limit >= 0 && offset >= 0 {
		t = t.Offset(offset).Limit(limit)
	}
	return count, t.Find(objs).Error
}

func GetObjectCombine(objs interface{}, offset int, limit int, name string, userID uint64, isAdmin bool) (int, error) {
	if isAdmin {
		return GetObjectCombinCustom(objs, offset, limit, name, nil, nil)
	} else {
		return GetObjectCombinCustom(objs, offset, limit, name, []string{"user_id = ?"}, []interface{}{userID})
	}
}

func GetObjectByID(obj interface{}, id uint64) error {
	return DB.Where("id = ?", id).First(obj).Error
}

func DeleteObjectByID(obj interface{}, id uint64) error {
	return DB.Where("id = ?", id).Delete(obj).Error
}

func IsObjectExistsCustom(objs interface{}, queries []string, values []interface{}) bool {
	query := ""
	for i, v := range queries {
		query += v
		if i != len(queries)-1 {
			query += " AND "
		}
	}
	return !DB.Where(query, values...).First(objs).RecordNotFound()
}

func IsObjectExistsByID(obj interface{}, id uint64) bool {
	return IsObjectExistsCustom(obj, []string{"id = ?"}, []interface{}{id})
}

func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", "./master.db")
	if err != nil {
		log.Fatal(err)
	}
	DB.Exec("PRAGMA foreign_keys = ON")
	DB.SingularTable(true)
	DB.AutoMigrate(&Deployment{}, &Task{}, &StorageItem{}, &TaskEnvironment{}, &TaskArgument{}, &TaskCrash{}, &TaskFuzzResult{}, &TaskFuzzResultStat{}, &User{})

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
		if err = CreateUser(username, string(password), true); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("\ncreate success")
		}
	}
}
