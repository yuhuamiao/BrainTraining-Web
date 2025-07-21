package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusHandler struct {
	dao *dao.BusDAO
}

func NewBusHandler(dao *dao.BusDAO) *BusHandler {
	return &BusHandler{dao: dao}
}

func (h *BusHandler) BusRoutes(r *gin.Engine) {
	r.POST("/api/v1/bus/scores", h.SubmitScore)
}

// SubmitScore 提交成绩
// @Summary 提交公交车训练成绩
// @Description 保存用户公交车训练成绩
// @Tags 公交车人数
// @Accept json
// @Produce json
// @Param score body BusScoreRequest true "成绩数据"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/bus/scores [post]
func (h *BusHandler) SubmitScore(c *gin.Context) {
	var req BusScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidRequest",
			Message: "无效的请求参数",
			Details: gin.H{"error": err.Error()},
		})
		return
	}

	accuracy := float64(req.SuccessNum) / float64(3)

	record := &models.BusRecord{
		UserID:      req.UserID,
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

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data: gin.H{
			"accuracy": accuracy,
			"level":    req.Level,
		},
	})
}

type BusScoreRequest struct {
	UserID      string  `json:"userId" binding:"required"`
	SuccessNum  int     `json:"successNum" binding:"required"`
	Level       string  `json:"level" binding:"required"`
	Accuracy    float64 `json:"accuracy" binding:"required"`
	TrainingNum int     `json:"trainingNum" binding:"required"`
}
