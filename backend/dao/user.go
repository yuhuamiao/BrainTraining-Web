package dao

import (
	"braintraining/backend/models"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// 新增方法
func (d *UserDAO) CreateUser(user *models.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := d.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// 用户基本信息操作
func (d *UserDAO) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := d.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (d *UserDAO) UpdateUser(user *models.User) error {
	return d.db.Save(user).Error
}

func (d *UserDAO) UpdateUsername(userID, username string) error {
	return d.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("username", username).
		Error
}

func (d *UserDAO) UpdateLastLogin(userID string, timestamp int64) error {
	return d.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("last_login_at", timestamp).
		Error
}

func (d *UserDAO) UpdateAvatar(userID, avatarURL string) error {
	return d.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("avatar", avatarURL).
		Error
}

// 用户成绩查询
func (d *UserDAO) GetUserScores(userID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 舒尔特表格成绩
	var schulteScores []models.SchulteScore
	if err := d.db.Where("user_id = ?", userID).Find(&schulteScores).Error; err != nil {
		return nil, err
	}
	result["schulte"] = schulteScores

	// 多色文字成绩
	var colorWordScores []models.ColorWordRecord
	if err := d.db.Where("user_id = ?", userID).Find(&colorWordScores).Error; err != nil {
		return nil, err
	}
	result["color_word"] = colorWordScores

	// 瞬间记忆成绩
	var memoryScores []models.MemoryRecord
	if err := d.db.Where("user_id = ?", userID).Find(&memoryScores).Error; err != nil {
		return nil, err
	}
	result["memory"] = memoryScores

	// 公交人数成绩
	var busScores []models.BusRecord
	if err := d.db.Where("user_id = ?", userID).Find(&busScores).Error; err != nil {
		return nil, err
	}
	result["bus"] = busScores

	var sudokuScores []models.SudokuRecord
	if err := d.db.Where("user_id = ?", userID).Find(&sudokuScores).Error; err != nil {
		return nil, err
	}
	result["sudoku"] = sudokuScores

	return result, nil
}

// 获取用户最近成绩
func (d *UserDAO) GetRecentScores(userID string, limit int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 每种训练获取最近的几条记录
	var recent []map[string]interface{}

	// 舒尔特表格
	var schulte []models.SchulteScore
	d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&schulte)
	recent = append(recent, map[string]interface{}{"type": "schulte", "data": schulte})

	// 多色文字
	var colorWord []models.ColorWordRecord
	d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&colorWord)
	recent = append(recent, map[string]interface{}{"type": "color_word", "data": colorWord})

	// 瞬间记忆
	var memory []models.MemoryRecord
	d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&memory)
	recent = append(recent, map[string]interface{}{"type": "memory", "data": memory})

	// 公交人数
	var bus []models.BusRecord
	d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&bus)
	recent = append(recent, map[string]interface{}{"type": "bus", "data": bus})

	result["recent"] = recent
	return result, nil
}
