package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"bradobrei/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		log.Printf(
			"[HTTP] %s %s -> %d (%s) ip=%s",
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start).Round(time.Millisecond),
			c.ClientIP(),
		)
	}
}

func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Printf("[ERROR] %s %s -> %v", c.Request.Method, c.Request.URL.Path, err)
		}
	}
}

func RecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %v\n%s", r, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "internal_server_error",
					Code:    500,
					Message: "Внутренняя ошибка сервера. Изменения откатены.",
				})
			}
		}()
		c.Next()
	}
}
