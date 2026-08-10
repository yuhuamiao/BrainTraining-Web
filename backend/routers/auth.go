// 登录注册部分路由
package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"braintraining/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"

	"errors"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserDAO *dao.UserDAO
}

func NewAuthHandler(userDAO *dao.UserDAO) *AuthHandler {
	return &AuthHandler{UserDAO: userDAO}
}

//func (h *AuthHandler) AuthRoutes(r *gin.Engine) {
//	r.POST("/api/v1/login", h.Login)
//	r.POST("/api/v1/register", h.Register)
//}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	// Username 去除首尾空白后为 3 到 50 个字符
	Username string `json:"username" binding:"required" minLength:"3" maxLength:"50"`
	// Password 至少 6 个字符，UTF-8 编码后不能超过 72 字节
	Password string `json:"password" binding:"required" minLength:"6"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "注册信息"
// @Success 200
// @Failure 400
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
			Details: gin.H{"error": err.Error()},
		})
		return
	}

	var valid bool
	req.Username, valid = normalizeUsername(req.Username)
	if !valid {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidUsername",
			Message: "用户名长度应为 3 到 50 个字符",
		})
		return
	}
	if !validPassword(req.Password) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidPassword",
			Message: "密码至少需要 6 个字符，且不能超过 72 字节",
		})
		return
	}
	if _, err := h.UserDAO.GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "UsernameExists",
			Message: "用户名已存在",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "检查用户名失败",
		})
		return
	}

	// 创建用户
	user := &models.User{
		UserID:   generateUserID(), // 生成唯一ID
		Username: req.Username,
		//Email:    req.Email,
		Password: req.Password, // 会在HashPassword中加密
	}

	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "PasswordError",
			Message: "密码处理失败",
		})
		return
	}

	if err := h.UserDAO.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "创建用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "注册成功",
		Data: gin.H{
			"userId": user.UserID,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取认证token
// @Tags 认证
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "登录凭证"
// @Success 200
// @Failure 400
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "InvalidUsername", Message: "用户名不能为空"})
		return
	}

	// 获取用户
	user, err := h.UserDAO.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AuthFailed",
			Message: "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "AuthFailed",
			Message: "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now().Unix()
	if err := h.UserDAO.UpdateLastLogin(user.UserID, user.LastLoginAt); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "更新登录状态失败",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "TokenError",
			Message: "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
		User: gin.H{
			"userId":   user.UserID,
			"username": user.Username,
			//"email":    user.Email,
			"avatar": user.Avatar,
		},
	})
}

// 辅助函数
func generateUserID() string {
	// 实际项目中应该使用更好的ID生成方式
	return uuid.New().String()
}

// 响应结构体
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	User    gin.H  `json:"user"`
}
