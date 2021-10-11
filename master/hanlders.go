package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/model"
	"github.com/CyDrive/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Account handlers
func RegisterHandle(c *gin.Context) {
	email, ok := c.GetPostForm("email")
	if !ok {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no account email",
		})
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no password",
		})
		return
	}

	account := &model.Account{
		Email:    email,
		Password: password,
	}

	if name, ok := c.GetPostForm("name"); ok {
		account.Name = name
	}
	if cap, ok := c.GetPostForm("cap"); ok {
		var err error
		account.Cap, err = strconv.ParseInt(cap, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, model.Response{
				StatusCode: consts.StatusCode_InvalidParameters,
				Message:    "invalid parameter: cap",
			})
			return
		}
	}

	err := GetAccountStore().AddAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "register account error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "account created",
	})
}

func LoginHandle(c *gin.Context) {
	email, ok := c.GetPostForm("email")
	if !ok {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no account email",
		})
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no password",
		})
		return
	}

	account, err := GetAccountStore().GetAccountByEmail(email)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no such account",
		})
		return
	}
	if utils.PasswordHash(account.Password) != password {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "account name or password not correct",
		})
		return
	}

	userSession := sessions.DefaultMany(c, "account")

	userSession.Set("userStruct", &account)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err = userSession.Save()
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	safeUser := utils.PackSafeAccount(account)
	userBytes, err := json.Marshal(safeUser)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "Welcome to CyDrive!",
		Data:       string(userBytes),
	})
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*model.Account)

	path := c.Param("path")

	path = strings.Trim(path, "/")
	absPath := strings.Join([]string{account.DataDir, path}, "/")

	fileList, err := GetEnv().ReadDir(absPath)
	for i := range fileList {
		fileList[i].FilePath = strings.ReplaceAll(fileList[i].FilePath, "\\", "/")
		fileList[i].FilePath = strings.TrimPrefix(fileList[i].FilePath, account.DataDir)
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	fileListJson, err := json.Marshal(fileList)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "list done",
		Data:       string(fileListJson),
	})
}

func GetFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*model.Account)

	filePath := c.Param("path")
	filePath = strings.Trim(filePath, "/")
	absFilePath := filepath.Join(account.DataDir, filePath)

	fileInfo, err := GetEnv().Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	fileInfoBytes, err := json.Marshal(fileInfo)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "get file info done",
		Data:       string(fileInfoBytes),
	})
}

// func PutFileInfoHandle(c *gin.Context) {
// 	userI, _ := c.Get("account")
// 	account := userI.(*model.User)

// 	filePath := c.Param("path")
// 	filePath = strings.Trim(filePath, "/")
// 	absFilePath := filepath.Join(account.DataDir, filePath)

// 	_, err := GetEnv().Stat(absFilePath)
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Resp{
// 			StatusCode:  consts.StatusCode_IoError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	len := c.Request.ContentLength
// 	fileInfoJson := make([]byte, len)
// 	c.Request.Body.Read(fileInfoJson)

// 	fileInfo := model.FileInfo{}
// 	if err := json.Unmarshal(fileInfoJson, &fileInfo); err != nil {
// 		c.JSON(http.StatusOK, model.Resp{
// 			StatusCode:  consts.StatusCode_InternalError,
// 			Message: "error when parsing file info",
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	openFile, err := GetEnv().OpenFile(absFilePath, os.O_RDWR, os.FileMode(fileInfo.FileMode))
// 	if err != nil {
// 		c.JSON(http.StatusOK, model.Resp{
// 			StatusCode:  consts.StatusCode_IoError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})
// 		return
// 	}
// 	defer openFile.Close()

// 	if err = GetEnv().Chtimes(absFilePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
// 		c.JSON(http.StatusOK, model.Resp{
// 			StatusCode:  consts.StatusCode_InternalError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})

// 		return
// 	}

// 	c.JSON(http.StatusOK, model.Resp{
// 		StatusCode:  consts.StatusCode_Ok,
// 		Message: "put file info done",
// 		Data:    nil,
// 	})
// }

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*model.Account)

	// relative path
	filePath := c.Param("path")

	// absolute filepath
	filePath = strings.Join([]string{account.DataDir, filePath}, "/")
	fileInfo, _ := GetEnv().Stat(filePath)

	if fileInfo.IsDir {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    "not a file",
		})
		return
	}

	// range
	var begin, _ int64 = 0, fileInfo.Size - 1
	bytesRange := c.GetHeader("Range")
	if len(bytesRange) > 0 {
		begin, _ = utils.UnpackRange(bytesRange)
	}

	taskId := GetFileTransferManager().AddTask(&fileInfo, account, DownloadTaskType, begin)

	uFileInfo := fileInfo
	uFileInfo.FilePath, _ = filepath.Rel(account.DataDir, uFileInfo.FilePath)
	uFileInfo.FilePath = strings.ReplaceAll(uFileInfo.FilePath, "\\", "/")

	resp := model.DownloadResponse{
		NodeAddr: config.IpAddr,
		TaskId:   taskId,
		FileInfo: &uFileInfo,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "download task created",
		Data:       string(respBytes),
	})
}

func UploadHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*model.Account)

	filePath := c.Param("path")

	filePath = strings.Join([]string{account.DataDir, filePath}, "/")
	fileDir := filepath.Dir(filePath)
	if err := GetEnv().MkdirAll(fileDir, 0666); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	var req model.UploadRequest

	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    fmt.Sprintf("read request body error: %+v", err),
		})
		return
	}
	if len(jsonBytes) == 0 || json.Unmarshal(jsonBytes, &req) != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "need file info",
		})
		return
	}

	fileInfo := req.FileInfo

	// Check file size
	if fileInfo.Size > consts.FileSizeLimit {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_FileTooLarge,
			Message:    "file is too large",
		})
		return
	}

	// Check account storage capability
	if fileInfo.Size > account.Cap {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_FileTooLarge,
			Message: fmt.Sprintf("no enough capability, free storage: %vMiB, and size of the file: %vMiB",
				(account.Cap-account.Usage)>>20, fileInfo.Size>>20), // Convert Byte to MB
		})
		return
	}

	// Change the modified time
	if err = GetEnv().Chtimes(filePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	taskId := ftm.AddTask(fileInfo, account, UploadTaskType, fileInfo.Size)

	resp := model.UploadResponse{
		NodeAddr: config.IpAddr,
		TaskId:   taskId,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "upload task created",
		Data:       string(respBytes),
	})
}
