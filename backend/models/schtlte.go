// 舒尔特方格数据模型
package models

import "gorm.io/gorm"

type SchulteScore struct {
	gorm.Model
	UserID   string `gorm:"size:36;not null;index"`
	IsPassed bool   `gorm:"not null"`
	//Level       string  `gorm:"size:10;check:level IN ('easy','medium','hard')"`
	TrainingNum int     `gorm:"check:training_num BETWEEN 1 AND 6"`
	TimeElapsed float64 `gorm:"type:DECIMAL(5,2)"` //已用时间
}
