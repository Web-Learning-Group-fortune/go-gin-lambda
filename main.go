package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// 初始化 Gin 路由
	r := gin.Default()

	// 定义 GET API
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from AWS Lambda with Go Gin!",
		})
	})

	// 定义 POST API
	r.POST("/echo", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}
		c.JSON(200, gin.H{
			"received": requestBody,
		})
	})

	// 将 Gin 路由转换为 Lambda 处理器
	ginLambda = ginadapter.New(r)
}

func main() {
	// 启动 Lambda
	lambda.Start(ginLambda.Proxy)
}
