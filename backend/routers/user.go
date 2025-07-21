package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
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
// @Param userId query string true "用户ID"
// @Success 200
// @Failure 400
// @Router /api/v1/user/scores [get]
func (h *UserHandler) GetUserScores(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "MissingParameter",
			Message: "缺少userId参数",
		})
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

	c.JSON(http.StatusOK, scores)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户基本信息
// @Description 获取用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param userId query string true "用户ID"
// @Success 200
// @Failure 400
// @Router /api/v1/user/users [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "MissingParameter",
			Message: "缺少userId参数",
		})
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
// @Router /api/v1/user/change [post]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
		})
		return
	}

	user := &models.User{
		UserID:   req.UserID,
		Username: req.Username,
		//Email:    req.Email,
	}

	if err := h.dao.UpdateUser(user); err != nil {
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
// @Param userId formData string true "用户ID"
// @Param avatar formData file true "头像文件"
// @Success 200
// @Failure 400
// @Router /api/v1/user/photo [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.PostForm("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "MissingParameter",
			Message: "缺少userId参数",
		})
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
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidFileType",
			Message: "只支持jpg、jpeg和png格式的头像",
		})
		return
	}

	// 生成保存路径 (实际项目中应该上传到云存储)
	filename := "avatars/" + userID + ext
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SaveFileError",
			Message: "保存头像文件失败",
		})
		return
	}

	// 更新数据库中的头像URL (这里简化处理，实际应该使用完整URL)
	avatarURL := "/static/" + filename
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
	UserID   string `json:"userId" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserScoresResponse struct {
	Schulte   []models.SchulteScore    `json:"schulte"`
	ColorWord []models.ColorWordRecord `json:"color_word"`
	Memory    []models.MemoryRecord    `json:"memory"`
	Bus       []models.BusRecord       `json:"bus"`
}
