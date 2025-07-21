package models

import "gorm.io/gorm"

type MemoryRecord struct {
	gorm.Model
	UserID      string  `gorm:"size:36;not null;index"` // 用户ID
	SuccessNum  int     `gorm:"check:success_num"`
	Level       string  `gorm:"size:10;check:level IN ('easy','medium','hard')"`
	Accuracy    float64 `gorm:"not null"` // 准确率(SuccessNum/TotalNum)
	TrainingNum int     `gorm:"check:training_num BETWEEN 1 AND 6"`
}

//type MemoryGameConfig struct {
//	gorm.Model
//	GridSize     int     `gorm:"not null;default:30"`  // 方格总数(5x6=30)
//	InitialCells int     `gorm:"not null;default:5"`   // 初始显示单元格数
//	DisplayTime  float64 `gorm:"not null;default:3.0"` // 显示时间(秒)
//	AnswerTime   float64 `gorm:"not null;default:5.0"` // 答题时间(秒)
//	Difficulty   string  `gorm:"size:10;not null"`     // 难度级别
//}
