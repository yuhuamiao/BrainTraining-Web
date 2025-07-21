package dao

import (
	"braintraining/backend/models"
	"gorm.io/gorm"
)

type MemoryDAO struct {
	db *gorm.DB
}

func NewMemoryDAO(db *gorm.DB) *MemoryDAO {
	return &MemoryDAO{db: db}
}

func (d *MemoryDAO) CreateRecord(record *models.MemoryRecord) error {
	return d.db.Create(record).Error
}

func (d *MemoryDAO) GetUserRecords(userID string, limit int) ([]models.MemoryRecord, error) {
	var records []models.MemoryRecord
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).
		Error
	return records, err
}

//func (d *MemoryDAO) GetConfigByDifficulty(difficulty string) (*models.MemoryGameConfig, error) {
//	var config models.MemoryGameConfig
//	err := d.db.Where("difficulty = ?", difficulty).First(&config).Error
//	if err == gorm.ErrRecordNotFound {
//		// 返回默认配置
//		return &models.MemoryGameConfig{
//			GridSize:     30,
//			InitialCells: 5,
//			DisplayTime:  3.0,
//			AnswerTime:   5.0,
//			Difficulty:   difficulty,
//		}, nil
//	}
//	return &config, err
//}
//
//func (d *MemoryDAO) UpdateConfig(config *models.MemoryGameConfig) error {
//	return d.db.Save(config).Error
//}

func (d *MemoryDAO) GetUserStats(userID string) (*models.MemoryRecord, error) {
	var stats models.MemoryRecord
	err := d.db.Model(&models.MemoryRecord{}).
		Select("AVG(accuracy) as accuracy, AVG(response_time) as response_time, MAX(level) as level").
		Where("user_id = ?", userID).
		Scan(&stats).
		Error
	return &stats, err
}
