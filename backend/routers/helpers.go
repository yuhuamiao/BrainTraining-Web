package routers

import (
	"math"
	"net/http"
	"strings"

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

func normalizeUsername(value string) (string, bool) {
	username := strings.TrimSpace(value)
	length := len([]rune(username))
	return username, length >= 3 && length <= 50
}

func validPassword(password string) bool {
	return len([]rune(password)) >= 6 && len([]byte(password)) <= 72
}

func validateTrainingScore(c *gin.Context, req TrainingScoreRequest) bool {
	expectedAccuracy := 0.0
	if req.TrainingNum > 0 {
		expectedAccuracy = float64(req.SuccessNum) / float64(req.TrainingNum)
	}
	if !validDifficulty(req.Level) || req.TrainingNum < 1 || req.SuccessNum < 0 || req.SuccessNum > req.TrainingNum || math.IsNaN(req.Accuracy) || req.Accuracy < 0 || req.Accuracy > 1 || math.Abs(req.Accuracy-expectedAccuracy) > 0.0001 {
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
