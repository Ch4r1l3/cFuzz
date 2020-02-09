package models

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/pbkdf2"
	"time"
)

// swagger:model
type User struct {
	ID uint64 `gorm:"primary_key" json:"id"`
	// example:123
	Username string `json:"username"`

	Password string `json:"-"`
	Salt     string `json:"-"`

	// example:true
	IsAdmin bool `json:"isAdmin"`
}

func GetEncryptPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 15000, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(dk)
}

func CreateUser(username string, password string, isAdmin bool) error {
	var count int
	if err := DB.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count >= 1 {
		return errors.New("username cannot be the same")
	}
	salt, err := utils.RandomString(12)
	if err != nil {
		return err
	}

	user := User{
		Username: username,
		Password: GetEncryptPassword(password, salt),
		Salt:     salt,
		IsAdmin:  isAdmin,
	}
	return InsertObject(&user)
}

func VerifyUser(username string, password string) (bool, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	encPassword := GetEncryptPassword(password, user.Salt)
	return encPassword == user.Password, nil
}

type CustomClaims struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"isAdmin"`
	jwt.StandardClaims
}

func CreateToken(id uint64, isAdmin bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := CustomClaims{
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

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Couldn't handle this token")
}

func GetObjectsByUserID(objs interface{}, userID uint64) error {
	return DB.Order("id").Where("user_id = ?", userID).Find(objs).Error
}

func GetCountByUserID(objs interface{}, userID uint64) (int, error) {
	var count int
	err := DB.Model(objs).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Model(User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUsernameExists(username string) bool {
	return IsObjectExistsCustom(&User{}, []string{"username = ?"}, []interface{}{username})
}

func GetNormalUserCombine(offset, limit int, name string) ([]User, int, error) {
	var users []User
	var count int
	var err error
	if name != "" {
		count, err = GetObjectCombinCustom(&users, offset, limit, "", []string{"is_admin = ?", "username like ?"}, []interface{}{false, "%" + name + "%"})
	} else {
		count, err = GetObjectCombinCustom(&users, offset, limit, "", []string{"is_admin = ?"}, []interface{}{false})
	}
	return users, count, err
}
