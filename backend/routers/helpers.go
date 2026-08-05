package routers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TrainingScoreRequest struct {
	SuccessNum  int     `json:"successNum"`
	Level       string  `json:"level"`
	Accuracy    float64 `json:"accuracy"`
	TrainingNum int     `json:"trainingNum"`
}

func currentUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get("userID")
	userID, ok := value.(string)
	return userID, exists && ok && userID != ""
}

func validateTrainingScore(c *gin.Context, req TrainingScoreRequest) bool {
	if !validDifficulty(req.Level) || req.TrainingNum < 1 || req.SuccessNum < 0 || req.SuccessNum > req.TrainingNum || math.IsNaN(req.Accuracy) || req.Accuracy < 0 || req.Accuracy > 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "InvalidScore",
			Message: "成绩数据不合法",
		})
		return false
	}
	return true
}

func validDifficulty(level string) bool {
	return level == "easy" || level == "medium" || level == "hard"
}

func success(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, SuccessResponse{Success: true, Data: data})
}
