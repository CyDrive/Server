package handles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/env"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
)

func LoginHandle(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no user name",
			Data:    nil,
		})
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no password",
			Data:    nil,
		})
		return
	}

	user := userStore.GetUserByName(username)
	if user == nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no such user",
			Data:    nil,
		})
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "user name or password not correct",
			Data:    nil,
		})
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	userJson, err := json.Marshal(user.SafeUser)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "Welcome to CyDrive!",
		Data:    string(userJson),
	})
}
