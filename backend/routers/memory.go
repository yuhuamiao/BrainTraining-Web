package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type MemoryHandler struct {
	dao *dao.MemoryDAO
}

func NewMemoryHandler(dao *dao.MemoryDAO) *MemoryHandler {
	return &MemoryHandler{
		dao: dao,
	}
}

//func (h *MemoryHandler) MemoryRoutes(r *gin.Engine) {
//	r.GET("/api/v1/memory/matrix", h.MemoryMatrix)
//	r.POST("/api/v1/memory/scores", h.SubmitScore)
//}

// MemoryMatrix 生成记忆矩阵
// @Summary 生成记忆训练矩阵
// @Description 生成一组需要记忆的方格位置
// @Tags 瞬间记忆
// @Accept json
// @Produce json
// @Param difficulty query string false "难度级别" Enums(easy,medium,hard) default(easy)
// @Success 200 {object} MemoryMatrixResponse
// @Failure 400
// @Router /api/v1/memory/matrix [get]
func (h *MemoryHandler) MemoryMatrix(c *gin.Context) {
	difficulty := c.DefaultQuery("difficulty", "easy")
	if !validDifficulty(difficulty) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "InvalidDifficulty", Message: "难度必须是 easy、medium 或 hard"})
		return
	}
	config := h.GetConfigByDifficulty(difficulty)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, ErrorResponse{
	//		Error:   "ConfigError",
	//		Message: "获取配置失败",
	//	})
	//	return
	//}

	// 生成随机方格位置
	positions := generatePositions(config.GridSize, config.InitialCells)

	c.JSON(http.StatusOK, MemoryMatrixResponse{
		Positions:   positions,
		Difficulty:  difficulty,
		TotalCells:  config.GridSize,
		DisplayTime: config.DisplayTime,
		AnswerTime:  config.AnswerTime,
	})
}

// SubmitScore 提交成绩
// @Summary 提交记忆训练成绩
// @Description 保存用户记忆训练成绩
// @Tags 瞬间记忆
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param score body TrainingScoreRequest true "成绩数据"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/memory/scores [post]
func (h *MemoryHandler) SubmitScore(c *gin.Context) {
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

	record := &models.MemoryRecord{
		UserID:      userID,
		SuccessNum:  req.SuccessNum,
		Level:       req.Level,
		Accuracy:    req.Accuracy,
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

//// GetConfig 获取配置
//// @Summary 获取记忆训练配置
//// @Description 获取记忆训练的各项配置参数
//// @Tags 瞬间记忆
//// @Accept json
//// @Produce json
//// @Param difficulty query string false "难度级别" Enums(easy,medium,hard) default(easy)
//// @Success 200 {object} MemoryConfigResponse
//// @Router /api/v1/memory/config [get]
//func (h *MemoryHandler) GetConfig(c *gin.Context) {
//	difficulty := c.DefaultQuery("difficulty", "easy")
//	config, err := h.dao.GetConfigByDifficulty(difficulty)
//	if err != nil {
//		c.JSON(http.StatusOK, MemoryConfigResponse{
//			GridSize:     30,
//			InitialCells: 5,
//			DisplayTime:  3.0,
//			AnswerTime:   5.0,
//			Difficulty:   difficulty,
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, MemoryConfigResponse{
//		GridSize:     config.GridSize,
//		InitialCells: config.InitialCells,
//		DisplayTime:  config.DisplayTime,
//		AnswerTime:   config.AnswerTime,
//		Difficulty:   config.Difficulty,
//	})
//}

// 辅助函数
func generatePositions(gridSize, InitialCells int) []int {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	// 使用 Source 创建一个 Rand 实例
	r := rand.New(source) //通过时间设立种子获取随机数生成的源，再用源生成 Rand 实例，以此达到真正随机的效果

	positions := r.Perm(gridSize)[:InitialCells]
	return positions
}

func (h *MemoryHandler) GetConfigByDifficulty(difficulty string) *MemoryConfig {
	var config MemoryConfig

	switch difficulty {
	case "easy":
		config = MemoryConfig{
			GridSize:     30,
			InitialCells: 5,
			DisplayTime:  3.0, // 显示3秒
			AnswerTime:   8.0,
		}
	case "medium":
		config = MemoryConfig{
			GridSize:     42,
			InitialCells: 7,
			DisplayTime:  2.5,
			AnswerTime:   7.0,
		}
	case "hard":
		config = MemoryConfig{
			GridSize:     56,
			InitialCells: 9,
			DisplayTime:  2.0,
			AnswerTime:   6.0,
		}
	default:
		config = MemoryConfig{
			GridSize:     30,
			InitialCells: 5,
			DisplayTime:  3.0,
			AnswerTime:   8.0,
		}
	}
	return &config
}

// 配置结构体
type MemoryConfig struct {
	GridSize     int     // 总方格数
	InitialCells int     // 初始显示单元格数
	DisplayTime  float64 // 显示时间(秒)
	AnswerTime   float64 // 答题时间(秒)
}

// 请求和响应结构体
type MemoryMatrixResponse struct {
	Positions   []int   `json:"positions"`   //返回的要亮的 15 个数字
	DisplayTime float64 `json:"displayTime"` //返回 展示时间
	AnswerTime  float64 `json:"answerTime"`  //返回 回答时间
	Difficulty  string  `json:"difficulty"`  //返回 等级难度
	TotalCells  int     `json:"totalCells"`  //返回一共所有的单元格
}

//type MemoryConfigResponse struct {
//	GridSize     int     `json:"gridSize"`
//	InitialCells int     `json:"initialCells"`
//	DisplayTime  float64 `json:"displayTime"`
//	AnswerTime   float64 `json:"answerTime"`
//	Difficulty   string  `json:"difficulty"`
//}
