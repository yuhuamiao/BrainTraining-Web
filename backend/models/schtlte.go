// 舒尔特方格数据模型
package models

import "gorm.io/gorm"

type SchulteScore struct {
	gorm.Model
	UserID      string  `json:"-" gorm:"size:36;not null;index"`
	IsPassed    bool    `json:"isPassed" gorm:"not null"`
	SuccessNum  int     `json:"successNum" gorm:"not null"`
	TrainingNum int     `json:"trainingNum" gorm:"not null"`
	TimeElapsed float64 `json:"timeElapsed" gorm:"not null"`
}
