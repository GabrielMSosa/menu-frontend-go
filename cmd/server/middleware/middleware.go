package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMIddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timerq := time.Now()
		ctx.Next()
		durati := time.Since(timerq).Milliseconds()
		method := ctx.Request.Method
		code := ctx.Writer.Status()
		url := ctx.Request.URL.String()
		client := ctx.ClientIP()
		log.Printf("Method: %s | URL: %s |Duration:%d ms | Code %d | IP: %s", method, url, durati, code, client)
	}
}
