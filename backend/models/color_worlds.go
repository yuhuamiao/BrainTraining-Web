package models

import "gorm.io/gorm"

type ColorWordRecord struct {
	gorm.Model
	UserID      string `gorm:"size:36;not null;index"` // 用户ID
	SuccessNum  int    `gorm:"check:success_num"`
	TrainingNum int    `gorm:"check:training_num BETWEEN 1 AND 6"`
}

type ColorWordConfig struct {
	gorm.Model
	ActiveColors []string `gorm:"serializer:json"` // 可用的颜色列表
}
