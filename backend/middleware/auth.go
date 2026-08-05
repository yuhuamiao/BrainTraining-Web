package middleware

import (
	"braintraining/backend/dao"
	"braintraining/backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userDAO *dao.UserDAO) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "缺少认证token",
			})
			return
		}

		// 2. 验证token格式: Bearer <token>
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "InvalidTokenFormat",
				"message": "token格式应为: Bearer <token>",
			})
			return
		}

		token := parts[1]

		// 3. 解析JWT token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "InvalidToken",
				"message": "无效的token: " + err.Error(),
			})
			return
		}

		// 4. 验证用户是否存在
		user, err := userDAO.GetUserByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "UserNotFound",
				"message": "用户不存在",
			})
			return
		}

		// 5. 将用户信息存入上下文
		c.Set("user", user)
		c.Set("userID", user.UserID)
		c.Next()
	}
}
