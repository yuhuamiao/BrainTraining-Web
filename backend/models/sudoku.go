package models

import "gorm.io/gorm"

type SudokuRecord struct {
	gorm.Model
	UserID      string  `json:"-" gorm:"size:36;not null;index"`
	Level       string  `json:"level" gorm:"size:10;not null"`
	IsPassed    bool    `json:"isPassed" gorm:"not null"`
	Mistakes    int     `json:"mistakes" gorm:"not null"`
	TimeElapsed float64 `json:"timeElapsed" gorm:"not null"`
}
