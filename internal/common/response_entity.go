package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseEntity struct {
	Code    int            `json:"code"`    // custom business code
	Message string         `json:"message"` // message for client
	Data    map[string]any `json:"data"`    // actual payload
}

func Success() *ResponseEntity {
	resp := &ResponseEntity{
		Code:    20000,
		Message: "success",
		Data:    make(map[string]interface{}),
	}
	return resp
}

// link append data
func (r *ResponseEntity) Append(key string, value interface{}) *ResponseEntity {
	r.Data[key] = value
	return r
}

// linq end write response
func (r *ResponseEntity) End(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func Fail(c *gin.Context, code int, msg string) {
	resp := &ResponseEntity{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
