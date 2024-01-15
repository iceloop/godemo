package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的路由器
	r := gin.Default()

	// 处理根路由 ("/") 的 GET 请求
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "欢迎使用 Gin Web 框架!",
		})
	})

	// 在 8080 端口启动服务
	r.Run(":8080")
}
