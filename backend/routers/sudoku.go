package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SudokuHandler struct {
	dao *dao.SudokuDAO
}

type SudokuScoreRequest struct {
	Level       string  `json:"level"`
	IsPassed    bool    `json:"isPassed"`
	Mistakes    int     `json:"mistakes"`
	TimeElapsed float64 `json:"timeElapsed"`
}

func NewSudokuHandler(sudokuDAO *dao.SudokuDAO) *SudokuHandler {
	return &SudokuHandler{dao: sudokuDAO}
}

// SubmitScore 提交数独成绩
// @Summary 提交数独成绩
// @Description 保存当前用户的数独完成情况
// @Tags 数独挑战
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param score body SudokuScoreRequest true "成绩数据"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /api/v1/sudoku/scores [post]
func (h *SudokuHandler) SubmitScore(c *gin.Context) {
	var req SudokuScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil || !validDifficulty(req.Level) || req.Mistakes < 0 || req.TimeElapsed < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "InvalidScore", Message: "成绩数据不合法"})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}
	record := &models.SudokuRecord{
		UserID:      userID,
		Level:       req.Level,
		IsPassed:    req.IsPassed,
		Mistakes:    req.Mistakes,
		TimeElapsed: req.TimeElapsed,
	}
	if err := h.dao.CreateRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "DatabaseError", Message: "保存成绩失败"})
		return
	}
	success(c, gin.H{"isPassed": req.IsPassed, "timeElapsed": req.TimeElapsed})
}
