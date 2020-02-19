package service

import (
	"github.com/Ch4r1l3/cFuzz/master/server/models"
)

func insertObject(obj interface{}) error {
	return models.DB.Create(obj).Error
}

func getObjects(objs interface{}) error {
	return models.DB.Order("id").Find(objs).Error
}

func getObjectCombinCustom(objs interface{}, offset int, limit int, name string, queries []string, values []interface{}) (int, error) {
	var count int
	var err error
	query := ""
	for _, v := range queries {
		query += v + " AND "
	}
	if name != "" {
		err = models.DB.Model(objs).Where(query+"name LIKE ?", append(values, "%"+name+"%")...).Count(&count).Error
	} else {
		err = models.DB.Model(objs).Where(query+"1=1", values...).Count(&count).Error
	}
	if err != nil {
		return 0, err
	}
	t := models.DB.Order("id")
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
		return getObjectCombinCustom(objs, offset, limit, name, nil, nil)
	} else {
		return getObjectCombinCustom(objs, offset, limit, name, []string{"user_id = ?"}, []interface{}{userID})
	}
}

func GetObjectByID(obj interface{}, id uint64) error {
	return models.DB.Where("id = ?", id).First(obj).Error
}

func DeleteObjectByID(obj interface{}, id uint64) error {
	return models.DB.Where("id = ?", id).Delete(obj).Error
}

func IsObjectExistsCustom(objs interface{}, queries []string, values []interface{}) bool {
	query := ""
	for i, v := range queries {
		query += v
		if i != len(queries)-1 {
			query += " AND "
		}
	}
	return !models.DB.Where(query, values...).First(objs).RecordNotFound()
}

func IsObjectExistsByID(obj interface{}, id uint64) bool {
	return IsObjectExistsCustom(obj, []string{"id = ?"}, []interface{}{id})
}

func SaveObject(obj interface{}) error {
	return models.DB.Save(obj).Error
}

func UpdateObject(obj interface{}, data map[string]interface{}) error {
	return models.DB.Model(obj).Updates(data).Error
}
