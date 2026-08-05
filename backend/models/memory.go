package models

import "gorm.io/gorm"

type MemoryRecord struct {
	gorm.Model
	UserID      string  `json:"-" gorm:"size:36;not null;index"`
	SuccessNum  int     `json:"successNum" gorm:"not null"`
	Level       string  `json:"level" gorm:"size:10;not null"`
	Accuracy    float64 `json:"accuracy" gorm:"not null"`
	TrainingNum int     `json:"trainingNum" gorm:"not null"`
}

//type MemoryGameConfig struct {
//	gorm.Model
//	GridSize     int     `gorm:"not null;default:30"`  // 方格总数(5x6=30)
//	InitialCells int     `gorm:"not null;default:5"`   // 初始显示单元格数
//	DisplayTime  float64 `gorm:"not null;default:3.0"` // 显示时间(秒)
//	AnswerTime   float64 `gorm:"not null;default:5.0"` // 答题时间(秒)
//	Difficulty   string  `gorm:"size:10;not null"`     // 难度级别
//}
