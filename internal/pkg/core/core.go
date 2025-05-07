package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onexstack/fastgo/internal/pkg/errorsx"
)

// ErrorResponse 定义了错误响应的结构体.
// 用于API 请求中发生错误时返回统一的格式化错误信息.
type ErrorResponse struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		errx := errorsx.FromError(err) // 提取错误详细信息
		c.JSON(errx.Code, ErrorResponse{
			Reason:  errx.Reason,
			Message: errx.Message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
