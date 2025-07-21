package models

import "gorm.io/gorm"

type BusRecord struct {
	gorm.Model
	UserID      string  `gorm:"size:36;not null;index"` // 用户ID
	SuccessNum  int     `gorm:"check:success_num"`
	Level       string  `gorm:"size:10;check:level IN ('easy','medium','hard')"`
	Accuracy    float64 `gorm:"not null"` // 准确率(SuccessNum/TotalNum)
	TrainingNum int     `gorm:"check:training_num BETWEEN 1 AND 6"`
}
