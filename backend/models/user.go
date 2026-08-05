package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      string `json:"userId" gorm:"size:36;not null;uniqueIndex"`
	Username    string `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Password    string `json:"-" gorm:"size:255;not null"`
	Avatar      string `json:"avatar" gorm:"size:255"`
	LastLoginAt int64  `json:"lastLogin"`
}

// 加密密码 (在创建/更新用户前调用)
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
