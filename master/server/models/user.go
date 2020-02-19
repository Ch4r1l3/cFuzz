package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// swagger:model
type User struct {
	ID uint64 `gorm:"primary_key" json:"id"`
	// example:123
	Username string `json:"username" sql:"type:varchar(255) NOT NULL UNIQUE"`

	Password string `json:"-"`
	Salt     string `json:"-"`

	// example:true
	IsAdmin bool `json:"isAdmin"`
}

type CustomClaims struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"isAdmin"`
	jwt.StandardClaims
}
