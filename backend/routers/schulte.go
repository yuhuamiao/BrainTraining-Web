package routers

import (
	"braintraining/backend/dao"
	"braintraining/backend/models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

// MatrixResponse 生成矩阵的响应结构
type MatrixResponse struct {
	Matrix [][]int `json:"matrix"`
	Size   int     `json:"size"`
}

// SchulteScoreRequest 提交成绩的请求结构
type SchulteScoreRequest struct {
	IsPassed    bool    `json:"isPassed"`
	SuccessNum  int     `json:"successNum"`
	TrainingNum int     `json:"trainingNum"`
	TimeElapsed float64 `json:"timeElapsed"`
}

// SchulteMatrix 生成舒尔特矩阵
// @Summary 舒尔特表格
// @Description 为舒尔特表格提供随机数字矩阵
// @Tags 舒尔特矩阵
// @Accept application/json
// @Produce application/json
// @Param size query int true "矩阵大小（如 5 表示 5x5 矩阵）" default(5) minimum(3) maximum(10)
// @Success 200 {object} MatrixResponse "成功返回矩阵数据"
// @Router /api/v1/schulte/matrix [get]
func SchulteMatrix(c *gin.Context) {
	sizestr := c.DefaultQuery("size", "5")
	size, err := strconv.Atoi(sizestr)

	if err != nil || size < 3 || size > 10 {
		c.JSON(400, gin.H{
			"code":  400,
			"error": "无效矩阵参数",
		})
		return
	}

	c.JSON(200, MatrixResponse{
		Matrix: generateMatrix(size),
		Size:   size,
	})
}

// SubmitScore 提交舒尔特成绩
// @Summary 提交成绩
// @Description 保存用户舒尔特表格游戏成绩
// @Tags 舒尔特矩阵
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token" default(Bearer <your_token>)
// @Param score body SchulteScoreRequest true "成绩数据"
// @Success 200 {object} map[string]interface{} "保存成功响应"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/schulte/scores [post]
func SubmitScore(c *gin.Context, dao *dao.SchulteDAO) {
	//dao := c.MustGet("dao").(*dao.SchulteDAO)

	var score SchulteScoreRequest
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(400, ErrorResponse{Error: "InvalidRequest", Message: "无效的请求参数"})
		return
	}
	passed := score.SuccessNum == score.TrainingNum
	if score.TrainingNum != 25 || score.SuccessNum < 0 || score.SuccessNum > score.TrainingNum || score.IsPassed != passed || score.TimeElapsed < 0 || score.TimeElapsed > 30.5 {
		c.JSON(400, ErrorResponse{Error: "InvalidScore", Message: "成绩数据不合法"})
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(401, ErrorResponse{Error: "Unauthorized", Message: "用户身份无效"})
		return
	}

	// TODO: 存储到数据库
	scoreSql := &models.SchulteScore{
		UserID:      userID,
		IsPassed:    score.IsPassed,
		SuccessNum:  score.SuccessNum,
		TrainingNum: score.TrainingNum,
		TimeElapsed: score.TimeElapsed,
	}

	if err := dao.CreateScore(scoreSql); err != nil {
		c.JSON(500, ErrorResponse{Error: "DatabaseError", Message: "保存成绩失败"})
		return
	}

	success(c, gin.H{"isPassed": score.IsPassed, "timeElapsed": score.TimeElapsed})
}

//// SchulteRoutes 注册舒尔特相关路由
//func SchulteRoutes(router *gin.Engine, dao *dao.SchulteDAO) {
//	router.Use(func(c *gin.Context) {
//		c.Set("dao", dao)
//		c.Next()
//	})
//	router.GET("/api/v1/schulte/matrix", SchulteMatrix)
//	router.POST("/api/v1/schulte/scores", SubmitScore)
//}

func generateMatrix(size int) [][]int {
	Matrix := make([][]int, size) //创建一个 size * size 大小的矩阵切片

	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	// 使用 Source 创建一个 Rand 实例
	r := rand.New(source) //通过时间设立种子获取随机数生成的源，再用源生成 Rand 实例，以此达到真正随机的效果

	num := r.Perm(size * size) //创建一个 size * size 大小的包含 0 到 size * size - 1 的随机排列的整数切片
	for i := 0; i < size; i++ {
		Matrix[i] = num[i*size : (i+1)*size] //一行填入 size 数的 num 切片
		for j := 0; j < size; j++ {
			Matrix[i][j]++ //因为生成的随机数是 0 到 size * size - 1 的，需要每个再加一
		}
	}
	return Matrix
}
