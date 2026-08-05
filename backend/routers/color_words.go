package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		parsed, err := url.Parse(origin)
		if err != nil {
			return false
		}
		if strings.EqualFold(parsed.Host, r.Host) {
			return true
		}
		allowedOrigin := strings.TrimRight(strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGIN")), "/")
		return allowedOrigin != "" && strings.EqualFold(strings.TrimRight(origin, "/"), allowedOrigin)
	},
}

// 这里新加了一个名为 ColorWordHandler 的结构体及其构造函数，用于处理 HTTP 请求并调用 DAO 层进行数据操作
type ColorWordHandler struct {
	dao *dao.ColorWordDAO
}

func NewColorWordHandler(dao *dao.ColorWordDAO) *ColorWordHandler {
	return &ColorWordHandler{dao: dao}
}

//func (h *ColorWordHandler) ColorWordsRouter(r *gin.Engine) {
//	r.GET("/api/v1/color_words/matrix", h.ColorWordsMatrix)
//	r.POST("/api/v1/color_words/scores", h.SubmitScore)
//	r.GET("/api/v1/color_words/color", h.ColorStream)
//
//}

// ColorWordsMatrix 生成颜色文字矩阵
// @Summary 生成颜色文字
// @Description 生成一组文字和颜色不匹配的条目
// @Tags 多色文字
// @Accept json
// @Produce json
// @Param difficulty query string false "难度级别" Enums(easy,medium,hard) default(easy)
// @Success 200
// @Failure 400
// @Router /api/v1/color_words/matrix [get]
func (h *ColorWordHandler) ColorWordsMatrix(c *gin.Context) {
	difficulty := c.DefaultQuery("difficulty", "easy")
	if !validDifficulty(difficulty) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "InvalidDifficulty", Message: "难度必须是 easy、medium 或 hard"})
		return
	}
	matrix := generateColorMatrix(difficulty)
	c.JSON(200, gin.H{
		"matrix":     matrix,
		"difficulty": difficulty,
	})
}

// SubmitScore 提交成绩
// @Summary 提交颜色文字成绩
// @Description 保存用户颜色文字训练成绩
// @Tags 多色文字
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param score body TrainingScoreRequest true "成绩数据"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/color_words/scores [post]
func (h *ColorWordHandler) SubmitScore(c *gin.Context) {
	var req TrainingScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
			Details: gin.H{"error": err.Error()},
		})
		return
	}
	if !validateTrainingScore(c, req) {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}

	record := &models.ColorWordRecord{
		UserID:      userID,
		SuccessNum:  req.SuccessNum,
		Level:       req.Level,
		Accuracy:    req.Accuracy, //准确度
		TrainingNum: req.TrainingNum,
	}

	if err := h.dao.CreateRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DatabaseError",
			Message: "保存成绩失败",
		})
		return
	}

	success(c, gin.H{"accuracy": req.Accuracy, "level": req.Level})
}

// ColorStream WebSocket颜色流
// @Summary 颜色流推送
// @Description 通过WebSocket推送颜色变化
// @Tags 多色文字
// @Accept json
// @Produce json
// @Router /api/v1/color_words/color [get]
func (h *ColorWordHandler) ColorStream(c *gin.Context) {
	if c.Request.Header.Get("Upgrade") != "websocket" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "不是有效的WebSocket请求，缺少 Upgrade: websocket 头",
		})
		return
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 记录详细错误日志
		log.Printf("WebSocket升级失败: %v", err)

		// 检查是否是因为请求头不匹配导致的错误
		if _, ok := err.(websocket.HandshakeError); ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "WebSocketHandshakeError",
				Message: "WebSocket握手失败，请确保请求包含正确的头信息",
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "WebSocketError",
				Message: "无法建立WebSocket连接",
			})
		}
		return
	}

	defer conn.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			color := generateRandomColor()
			if err := conn.WriteJSON(gin.H{
				"color": color,
				"time":  time.Now().Unix(),
			}); err != nil {
				return
			}
		}
	}
}

// 辅助函数
func generateColorMatrix(difficulty string) [][]gin.H {
	// 根据难度生成不同复杂度的矩阵
	var count int
	switch difficulty {
	case "easy":
		count = 5
	case "medium":
		count = 6
	case "hard":
		count = 7
	default:
		count = 5
	}

	// 正确：应创建 count 个外层元素，每个元素再初始化内层切片
	matrix := make([][]gin.H, count)
	for i := range matrix {
		matrix[i] = make([]gin.H, count) // 初始化内层切片
	}
	colorNames := []string{"红", "蓝", "绿", "黄", "黑"}

	for i := 0; i < count; i++ {
		for j := 0; j < count; j++ {
			wordIdx := rand.Intn(len(colorNames))

			matrix[i][j] = gin.H{
				"word": colorNames[wordIdx],
			}

		}
	}

	return matrix
}

func generateRandomColor() string {
	colors := []string{"red", "blue", "green", "yellow", "black"}
	return colors[rand.Intn(len(colors))]
}

type ColorWordMatrixResponse struct {
	Matrix     [][]gin.H `json:"matrix"`
	Difficulty string    `json:"difficulty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details gin.H  `json:"details,omitempty"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    gin.H  `json:"data,omitempty"`
}
