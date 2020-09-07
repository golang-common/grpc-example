package http_utils

import (
	"github.com/gin-gonic/gin"
)

type ResponseJSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Object  interface{} `json:"data"`
}

type ResponseObject struct {
	Data interface{} `json:"data"`

	TotalCount int `json:"total_count,omitempty"`
}

func NormalResponse(c *gin.Context, errCode int, errMessage string, object interface{}) {
	var r ResponseJSON
	if object == nil {
		r = ResponseJSON{Code: errCode, Message: errMessage, Object: nil}
		c.JSON(200, r)
	} else {
		switch obj := object.(type) {
		case error:
			r = ResponseJSON{Code: errCode, Message: errMessage, Object: obj.Error()}
		default:
			r = ResponseJSON{Code: errCode, Message: errMessage, Object: obj}
		}
		c.JSON(200, r)
	}
}
