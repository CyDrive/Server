package master

import (
	"net/http"
	"strings"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func SetRequestId(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_id", uuid.NewString())
		c.Next()
	}
}

func LoginAuth(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiUrl := strings.Trim(c.Request.URL.Path, "/")
		if apiUrl == "login" || apiUrl == "register" {
			c.Next()
			return
		}

		userSession := sessions.DefaultMany(c, "account")
		account := userSession.Get("userStruct")
		expire := userSession.Get("expire")
		if account == nil || expire == nil {
			c.AbortWithStatusJSON(http.StatusOK, models.Response{
				StatusCode: consts.StatusCode_AuthError,
				Message:    "not login",
			})
			return
		}

		if !expire.(time.Time).After(time.Now()) {
			c.AbortWithStatusJSON(http.StatusOK, models.Response{
				StatusCode: consts.StatusCode_AuthError,
				Message:    "timeout, login again",
			})
			userSession.Clear()
			return
		}

		// Flush expire time
		userSession.Set("expire", time.Now().Add(time.Hour*12))
		if err := userSession.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, models.Response{
				StatusCode: consts.StatusCode_SessionError,
				Message:    err.Error(),
			})
			return
		}

		// Store account struct into context
		c.Set("account", account)
	}
}

func Log(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			account   *models.Account = nil
			requestId string          = ""
		)

		if accountI, ok := c.Get("account"); ok {
			account = accountI.(*models.Account)
		}

		if requestIdI, ok := c.Get("request_id"); ok {
			requestId = requestIdI.(string)
		}

		log.Infof("request_id=%+v (%s request_url=%s): headers=%+v from account=%+v, client_ip=%+v",
			requestId,
			c.Request.Method,
			c.Request.URL.String(),
			c.Request.Header,
			account,
			c.Request.RemoteAddr)

		c.Next()

		var headers http.Header = nil
		if c.Request.Response != nil {
			headers = c.Request.Response.Header
		}
		log.Infof("response to request_id=%+v: headers=%+v",
			requestId,
			headers)
	}
}
