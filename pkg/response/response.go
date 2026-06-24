// Package response defines the API's shared JSON response envelope.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

func Fail(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Body{Code: httpCode, Message: message})
}
