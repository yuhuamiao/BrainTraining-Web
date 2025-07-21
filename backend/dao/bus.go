package dao

import (
	"braintraining/backend/models"
	"gorm.io/gorm"
)

type BusDAO struct {
	db *gorm.DB
}

func NewBusDAO(db *gorm.DB) *BusDAO {
	return &BusDAO{db: db}
}

func (d *BusDAO) CreateRecord(record *models.BusRecord) error {
	return d.db.Create(record).Error
}

func (d *BusDAO) GetUserRecords(userID string, limit int) ([]models.BusRecord, error) {
	var records []models.BusRecord
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).
		Error
	return records, err
}

func (d *BusDAO) GetUserStats(userID string) (*models.BusRecord, error) {
	var stats models.BusRecord
	err := d.db.Model(&models.BusRecord{}).
		Select("AVG(accuracy) as accuracy, AVG(response_time) as response_time, MAX(level) as level").
		Where("user_id = ?", userID).
		Scan(&stats).
		Error
	return &stats, err
}
