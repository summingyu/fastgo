package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/onexstack/fastgo/internal/pkg/contextx"
	"github.com/onexstack/fastgo/internal/pkg/core"
	"github.com/onexstack/fastgo/internal/pkg/errorsx"
	"github.com/onexstack/fastgo/pkg/token"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, nil, errorsx.ErrTokenInvalid)
			c.Abort()
			return
		}
		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
