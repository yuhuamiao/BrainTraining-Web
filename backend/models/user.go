package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   string `gorm:"size:36;not null;uniqueIndex"` // 用户唯一标识
	Username string `gorm:"size:50;not null"`             // 用户名
	//Email       string `gorm:"size:100;not null"`            // 邮箱
	Password    string `gorm:"size:255;not null"` // 密码(加密存储)
	Avatar      string `gorm:"size:255"`          // 头像URL
	LastLoginAt int64  // 最后登录时间(unix timestamp)
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
