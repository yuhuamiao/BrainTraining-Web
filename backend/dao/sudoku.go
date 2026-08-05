package dao

import (
	"braintraining/backend/models"

	"gorm.io/gorm"
)

type SudokuDAO struct {
	db *gorm.DB
}

func NewSudokuDAO(db *gorm.DB) *SudokuDAO {
	return &SudokuDAO{db: db}
}

func (d *SudokuDAO) CreateRecord(record *models.SudokuRecord) error {
	return d.db.Create(record).Error
}
