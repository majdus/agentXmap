package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondSuccess sends a success response with data
func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// RespondCreated sends a created response with data
func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// RespondError sends an error response with a specific status code
func RespondError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, Response{
		Success: false,
		Error: gin.H{
			"message": message,
			"details": details,
		},
	})
}
