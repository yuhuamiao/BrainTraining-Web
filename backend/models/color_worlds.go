package models

import "gorm.io/gorm"

type ColorWordRecord struct {
	gorm.Model
	UserID      string  `json:"-" gorm:"size:36;not null;index"`
	SuccessNum  int     `json:"successNum" gorm:"not null"`
	Level       string  `json:"level" gorm:"size:10;not null"`
	Accuracy    float64 `json:"accuracy" gorm:"not null"`
	TrainingNum int     `json:"trainingNum" gorm:"not null"`
}

type ColorWordConfig struct {
	gorm.Model
	ActiveColors []string `gorm:"serializer:json"` // 可用的颜色列表
}
