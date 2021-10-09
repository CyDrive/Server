package master

import (
	"net/http"
	"strings"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginAuth(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Trim(c.Request.URL.Path, "/") == "login" {
			c.Next()
			return
		}

		userSession := sessions.DefaultMany(c, "user")
		user := userSession.Get("userStruct")
		expire := userSession.Get("expire")
		if user == nil || expire == nil {
			c.AbortWithStatusJSON(http.StatusOK, model.Response{
				StatusCode: consts.StatusCode_AuthError,
				Message:    "not login",
				Data:       nil,
			})
			return
		}

		if !expire.(time.Time).After(time.Now()) {
			c.AbortWithStatusJSON(http.StatusOK, model.Response{
				StatusCode: consts.StatusCode_AuthError,
				Message:    "timeout, login again",
				Data:       nil,
			})
			userSession.Clear()
			return
		}

		// Flush expire time
		userSession.Set("expire", time.Now().Add(time.Hour*12))
		if err := userSession.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, model.Response{
				StatusCode: consts.StatusCode_SessionError,
				Message:    err.Error(),
				Data:       nil,
			})
			return
		}

		// Store user struct into context
		c.Set("user", user)
	}
}