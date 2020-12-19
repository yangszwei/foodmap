package delivery

import (
	"foodmap/internal/infra/errors"
	"github.com/gin-gonic/gin"
)

// Error create error response for API requests
func Error(err error) gin.H {
	if e, ok := err.(errors.Error); ok {
		return gin.H{
			"code":    e.Name,
			"message": e.Message,
		}
	}
	if e, ok := err.(errors.ValidationError); ok {
		return gin.H{
			"code":   "invalid query",
			"fields": e.Field,
			"type":   e.Tag,
		}
	}
	return gin.H{
		"code":  "native",
		"error": err.Error(),
	}
}
