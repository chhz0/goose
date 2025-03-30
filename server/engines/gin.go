package engines

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GinEnginWrapper struct {
	ginEngine *gin.Engine
}

func (g *GinEnginWrapper) Handler() http.Handler {
	return g.ginEngine
}

func Gin(g *gin.Engine) *GinEnginWrapper {
	return &GinEnginWrapper{ginEngine: g}
}

/*
todo: 中间件：
	1. 认证和授权
	2. 日志记录
	3. 错误处理
	4. 性能监视
*/

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// todo: log
		ctx.Next()
	}
}

func NoCahce() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "Tue, 01 Jan 1970 00:00:00 GMT")
		ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		ctx.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "OPTIONS" {
			ctx.Next()
		} else {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			ctx.Header("Access-Control-Max-Age", "3600")
			ctx.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			ctx.Header("Content-Type", "application/json")
			ctx.AbortWithStatus(200)
		}
	}
}

func RequestID(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.Request.Header.Get(key)

		if requestID == "" {
			requestID = uuid.New().String()
		}
		valCtx := context.WithValue(ctx.Request.Context(), key, requestID)
		ctx.Request = ctx.Request.WithContext(valCtx)

		ctx.Writer.Header().Set(key, requestID)

		ctx.Next()
	}
}

// todo: 统一返回体 && 错误处理
