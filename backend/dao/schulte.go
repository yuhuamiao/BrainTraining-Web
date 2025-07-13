// 舒尔特方格数据访问层
package dao

import (
	"braintraining/backend/models"
	//"errors"
	"gorm.io/gorm"
)

type SchulteDAO struct {
	db *gorm.DB
}

func NewSchulteDAO(db *gorm.DB) *SchulteDAO {
	return &SchulteDAO{db: db}
}

func (d *SchulteDAO) CreateScore(score *models.SchulteScore) error {
	return d.db.Create(score).Error
}
