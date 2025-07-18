package dao

import (
	"braintraining/backend/models"
	"gorm.io/gorm"
)

type ColorWordDAO struct { //这是一个数据访问对象的实现，*gorm.DB：指向 GORM 数据库连接对象的指针。
	db *gorm.DB
}

func NewColorWordDAO(db *gorm.DB) *ColorWordDAO { //建立一个新的 dao 对象
	return &ColorWordDAO{db: db}
}

func (d *ColorWordDAO) CreateRecord(record *models.ColorWordRecord) error { //实现添加数据的方法
	return d.db.Create(record).Error
}

func (d *ColorWordDAO) GetUserRecords(userID string, limit int) ([]models.ColorWordRecord, error) { //实现一个获得用户数据的方法
	var records []models.ColorWordRecord
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).
		Error
	return records, err
}

func (d *ColorWordDAO) GetColorConfig() (*models.ColorWordConfig, error) {
	var config models.ColorWordConfig
	err := d.db.FirstOrCreate(&config).Error
	return &config, err
}

func (d *ColorWordDAO) UpdateColorConfig(colors []string) error {
	return d.db.Model(&models.ColorWordConfig{}).
		Update("active_colors", colors).
		Error
}
