package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UserHandler struct {
	dao *dao.UserDAO
}

func NewUserHandler(dao *dao.UserDAO) *UserHandler {
	return &UserHandler{dao: dao}
}

//func (h *UserHandler) UserRoutes(r *gin.Engine) {
//	userGroup := r.Group("/api/v1/user")
//	{
//		userGroup.GET("/scores", h.GetUserScores)
//		userGroup.GET("/users", h.GetUserInfo)
//		userGroup.POST("/change", h.UpdateUserInfo)
//		userGroup.POST("/photo", h.UploadAvatar)
//	}
//}

// GetUserScores 获取用户成绩
// @Summary 获取用户所有训练成绩
// @Description 获取用户在各个训练中的成绩记录
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Success 200
// @Failure 400
// @Router /api/v1/user/scores [get]
func (h *UserHandler) GetUserScores(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}

	scores, err := h.dao.GetUserScores(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "获取成绩失败",
		})
		return
	}

	c.JSON(http.StatusOK, normalizeScores(scores))
}

// GetUserInfo 获取用户信息
// @Summary 获取用户基本信息
// @Description 获取用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Success 200
// @Failure 400
// @Router /api/v1/user/me [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}

	user, err := h.dao.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "获取用户信息失败",
		})
		return
	}

	// 不返回密码等敏感信息
	response := gin.H{
		"userId":   user.UserID,
		"username": user.Username,
		//"email":     user.Email,
		"avatar":    user.Avatar,
		"lastLogin": user.LastLoginAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的基本信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param user body UserUpdateRequest true "用户更新数据"
// @Success 200
// @Failure 400
// @Router /api/v1/user/me [patch]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
		})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}
	username, valid := normalizeUsername(req.Username)
	if !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "InvalidUsername", Message: "用户名长度应为 3 到 50 个字符"})
		return
	}
	if existing, err := h.dao.GetUserByUsername(username); err == nil && existing.UserID != userID {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "UsernameExists", Message: "用户名已存在"})
		return
	}

	if err := h.dao.UpdateUsername(userID, username); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "更新用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "用户信息更新成功",
	})
}

// UploadAvatar 上传头像
// @Summary 上传用户头像
// @Description 上传并更新用户头像
// @Tags 用户
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param avatar formData file true "头像文件"
// @Success 200
// @Failure 400
// @Router /api/v1/user/photo [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidFile",
			Message: "获取头像文件失败",
		})
		return
	}

	// 限制文件大小 (2MB)
	if file.Size > 2<<20 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "FileTooLarge",
			Message: "头像文件不能超过2MB",
		})
		return
	}

	// 限制文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidFileType",
			Message: "只支持jpg、jpeg和png格式的头像",
		})
		return
	}

	// 生成保存路径 (实际项目中应该上传到云存储)
	if err := os.MkdirAll("avatars", 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "SaveFileError", Message: "创建头像目录失败"})
		return
	}
	storedName := userID + ext
	filename := filepath.Join("avatars", storedName)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SaveFileError",
			Message: "保存头像文件失败",
		})
		return
	}

	// 更新数据库中的头像URL (这里简化处理，实际应该使用完整URL)
	avatarURL := "/uploads/avatars/" + storedName
	if err := h.dao.UpdateAvatar(userID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "更新头像URL失败",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "头像上传成功",
		Data:    gin.H{"avatarUrl": avatarURL},
	})
}

// 请求和响应结构体
type UserUpdateRequest struct {
	Username string `json:"username" binding:"required"`
}

type ScoreSummary struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	SuccessNum  int       `json:"successNum,omitempty"`
	TrainingNum int       `json:"trainingNum,omitempty"`
	Accuracy    *float64  `json:"accuracy,omitempty"`
	Level       string    `json:"level,omitempty"`
	IsPassed    bool      `json:"isPassed"`
	TimeElapsed float64   `json:"timeElapsed,omitempty"`
	Mistakes    int       `json:"mistakes,omitempty"`
}

func normalizeScores(raw map[string]interface{}) gin.H {
	result := gin.H{
		"schulte":   []ScoreSummary{},
		"colorWord": []ScoreSummary{},
		"memory":    []ScoreSummary{},
		"bus":       []ScoreSummary{},
		"sudoku":    []ScoreSummary{},
	}

	if records, ok := raw["schulte"].([]models.SchulteScore); ok {
		items := make([]ScoreSummary, 0, len(records))
		for _, record := range records {
			items = append(items, ScoreSummary{ID: record.ID, CreatedAt: record.CreatedAt, SuccessNum: record.SuccessNum, TrainingNum: record.TrainingNum, IsPassed: record.IsPassed, TimeElapsed: record.TimeElapsed})
		}
		result["schulte"] = items
	}
	if records, ok := raw["color_word"].([]models.ColorWordRecord); ok {
		result["colorWord"] = summarizeTrainingRecords(records)
	}
	if records, ok := raw["memory"].([]models.MemoryRecord); ok {
		result["memory"] = summarizeTrainingRecords(records)
	}
	if records, ok := raw["bus"].([]models.BusRecord); ok {
		result["bus"] = summarizeTrainingRecords(records)
	}
	if records, ok := raw["sudoku"].([]models.SudokuRecord); ok {
		items := make([]ScoreSummary, 0, len(records))
		for _, record := range records {
			items = append(items, ScoreSummary{ID: record.ID, CreatedAt: record.CreatedAt, Level: record.Level, IsPassed: record.IsPassed, Mistakes: record.Mistakes, TimeElapsed: record.TimeElapsed})
		}
		result["sudoku"] = items
	}
	return result
}

func summarizeTrainingRecords[T models.ColorWordRecord | models.MemoryRecord | models.BusRecord](records []T) []ScoreSummary {
	items := make([]ScoreSummary, 0, len(records))
	for _, item := range records {
		switch record := any(item).(type) {
		case models.ColorWordRecord:
			items = append(items, ScoreSummary{ID: record.ID, CreatedAt: record.CreatedAt, SuccessNum: record.SuccessNum, TrainingNum: record.TrainingNum, Accuracy: float64Pointer(record.Accuracy), Level: record.Level})
		case models.MemoryRecord:
			items = append(items, ScoreSummary{ID: record.ID, CreatedAt: record.CreatedAt, SuccessNum: record.SuccessNum, TrainingNum: record.TrainingNum, Accuracy: float64Pointer(record.Accuracy), Level: record.Level})
		case models.BusRecord:
			items = append(items, ScoreSummary{ID: record.ID, CreatedAt: record.CreatedAt, SuccessNum: record.SuccessNum, TrainingNum: record.TrainingNum, Accuracy: float64Pointer(record.Accuracy), Level: record.Level})
		}
	}
	return items
}

func float64Pointer(value float64) *float64 {
	return &value
}
