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

//func (h *BusHandler) BusRoutes(r *gin.Engine) {
//	r.POST("/api/v1/bus/scores", h.SubmitScore)
//}

// SubmitScore 提交成绩
// @Summary 提交公交车训练成绩
// @Description 保存用户公交车训练成绩
// @Tags 公交车人数
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param score body TrainingScoreRequest true "成绩数据"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/bus/scores [post]
func (h *BusHandler) SubmitScore(c *gin.Context) {
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

	record := &models.BusRecord{
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
