package utils

import (
	"github.com/gin-gonic/gin"
)

type RespJson struct {
	statusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
	Msg        string `json:"message,omitempty"`
}

func JSONR(c *gin.Context, code int, dat interface{}, message interface{}) (werror error) {
	var (
		wcode int
		data  interface{}
		msg   interface{}
	)
	wcode = code
	data = dat
	msg = message

	var body interface{}

	switch msg.(type) {
	case string:
		body = gin.H{"code": wcode, "data": data, "msg": msg.(string)}
	case error:
		body = gin.H{"code": wcode, "data": data, "msg": msg.(error).Error()}
	default:
		body = gin.H{"code": wcode, "data": data, "msg": ""}
	}
	c.JSON(wcode, body)
	return
}
