package middleware

import (
	"hyper-pen-service/utils"

	"github.com/kataras/iris/v12"
)

// AuthRequired 验证用户是否已登录
func AuthRequired(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "Authorization header is required"})
		return
	}

	// 检查Bearer token格式
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "Invalid token format"})
		return
	}

	tokenString := authHeader[7:]
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "Invalid token"})
		return
	}

	// 将用户ID存储到上下文中
	ctx.Values().Set("userID", claims.UserID)
	ctx.Next()
}
