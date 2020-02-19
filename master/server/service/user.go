package service

import (
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"github.com/Ch4r1l3/cFuzz/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

func IsUserExistsByUsername(username string) bool {
	return IsObjectExistsCustom(&models.User{}, []string{"username = ?"}, []interface{}{username})
}

func IsUserExistsByID(id uint64) bool {
	return IsObjectExistsByID(&models.User{}, id)
}

func CreateUser(username string, password string, isAdmin bool) error {
	if IsUserExistsByUsername(username) {
		return errors.New("username cannot be the same")
	}
	salt, err := utils.RandomString(12)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Password: utils.GetEncryptPassword(password, salt),
		Salt:     salt,
		IsAdmin:  isAdmin,
	}
	return InsertObject(&user)
}

func VerifyUser(username string, password string) (bool, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	encPassword := utils.GetEncryptPassword(password, user.Salt)
	return encPassword == user.Password, nil
}

func CreateToken(id uint64, isAdmin bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := models.CustomClaims{
		int64(id),
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cfuzz",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.ServerConf.SigningKey))
}

func ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ServerConf.SigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("Token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token not active yet")
			} else {
				return nil, errors.New("Couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Couldn't handle this token")
}

func GetObjectsByUserID(objs interface{}, userID uint64) error {
	return models.DB.Order("id").Where("user_id = ?", userID).Find(objs).Error
}

func DeleteObjectsByUserID(objs interface{}, userID uint64) error {
	return models.DB.Where("user_id = ?", userID).Delete(objs).Error
}

func GetCountByUserID(objs interface{}, userID uint64) (int, error) {
	var count int
	err := models.DB.Model(objs).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := models.DB.Model(models.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetNormalUserCombine(offset, limit int, name string) ([]models.User, int, error) {
	var users []models.User
	var count int
	var err error
	if name != "" {
		count, err = GetObjectCombinCustom(&users, offset, limit, "", []string{"is_admin = ?", "username like ?"}, []interface{}{false, "%" + name + "%"})
	} else {
		count, err = GetObjectCombinCustom(&users, offset, limit, "", []string{"is_admin = ?"}, []interface{}{false})
	}
	return users, count, err
}

func DeleteUserByID(id uint64) error {
	DeleteStorageItemCustom("user_id = ?", id)
	var tasks []models.Task
	if GetObjectsByUserID(&tasks, id) == nil {
		for _, t := range tasks {
			DeleteTask(t.ID)
		}
	}
	return DeleteObjectByID(&models.User{}, id)
}
