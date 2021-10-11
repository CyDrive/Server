package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Account handlers
func RegisterHandle(c *gin.Context) {
	var req models.RegisterRequest
	reqBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    "error when read request body: " + err.Error(),
		})
		return
	}

	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "error when unmarshal request body: " + err.Error(),
		})
		return
	}

	if _, err = mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InvalidParameters,
			Message:    "invalid email address: " + req.Email,
		})
		return
	}

	account := &models.Account{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Cap:      req.Cap,
	}

	err = GetAccountStore().AddAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "register account error: " + err.Error(),
		})
		return
	}

	safeAccountBytes, _ := json.Marshal(utils.PackSafeAccount(account))

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "account created",
		Data:       string(safeAccountBytes),
	})
}

func LoginHandle(c *gin.Context) {
	var req models.LoginRequest
	reqBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    "error when read request body: " + err.Error(),
		})
		return
	}

	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "error when unmarshal request body: " + err.Error(),
		})
		return
	}

	account, err := GetAccountStore().GetAccountByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no such account",
		})
		return
	}
	if account.Password != req.Password {
		c.JSON(http.StatusOK, models.Response{
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
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	safeAccount := utils.PackSafeAccount(account)
	safeAccountBytes, err := json.Marshal(safeAccount)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "Welcome to CyDrive!",
		Data:       string(safeAccountBytes),
	})
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	path := c.Param("path")

	path = strings.Trim(path, "/")
	absPath := strings.Join([]string{account.DataDir, path}, "/")

	fileList, err := GetEnv().ReadDir(absPath)
	for i := range fileList {
		fileList[i].FilePath = strings.ReplaceAll(fileList[i].FilePath, "\\", "/")
		fileList[i].FilePath = strings.TrimPrefix(fileList[i].FilePath, account.DataDir)
	}
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	fileListJson, err := json.Marshal(fileList)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "list done",
		Data:       string(fileListJson),
	})
}

func GetFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	filePath := c.Param("path")
	filePath = strings.Trim(filePath, "/")
	absFilePath := filepath.Join(account.DataDir, filePath)

	fileInfo, err := GetEnv().Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	fileInfoBytes, err := json.Marshal(fileInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "get file info done",
		Data:       string(fileInfoBytes),
	})
}

// func PutFileInfoHandle(c *gin.Context) {
// 	userI, _ := c.Get("account")
// 	account := userI.(*models.User)

// 	filePath := c.Param("path")
// 	filePath = strings.Trim(filePath, "/")
// 	absFilePath := filepath.Join(account.DataDir, filePath)

// 	_, err := GetEnv().Stat(absFilePath)
// 	if err != nil {
// 		c.JSON(http.StatusOK, models.Resp{
// 			StatusCode:  consts.StatusCode_IoError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	len := c.Request.ContentLength
// 	fileInfoJson := make([]byte, len)
// 	c.Request.Body.Read(fileInfoJson)

// 	fileInfo := models.FileInfo{}
// 	if err := json.Unmarshal(fileInfoJson, &fileInfo); err != nil {
// 		c.JSON(http.StatusOK, models.Resp{
// 			StatusCode:  consts.StatusCode_InternalError,
// 			Message: "error when parsing file info",
// 			Data:    nil,
// 		})
// 		return
// 	}

// 	openFile, err := GetEnv().OpenFile(absFilePath, os.O_RDWR, os.FileMode(fileInfo.FileMode))
// 	if err != nil {
// 		c.JSON(http.StatusOK, models.Resp{
// 			StatusCode:  consts.StatusCode_IoError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})
// 		return
// 	}
// 	defer openFile.Close()

// 	if err = GetEnv().Chtimes(absFilePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
// 		c.JSON(http.StatusOK, models.Resp{
// 			StatusCode:  consts.StatusCode_InternalError,
// 			Message: err.Error(),
// 			Data:    nil,
// 		})

// 		return
// 	}

// 	c.JSON(http.StatusOK, models.Resp{
// 		StatusCode:  consts.StatusCode_Ok,
// 		Message: "put file info done",
// 		Data:    nil,
// 	})
// }

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	// relative path
	filePath := c.Param("path")

	// absolute filepath
	filePath = strings.Join([]string{account.DataDir, filePath}, "/")
	fileInfo, _ := GetEnv().Stat(filePath)

	if fileInfo.IsDir {
		c.JSON(http.StatusOK, models.Response{
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

	resp := models.DownloadResponse{
		NodeAddr: config.IpAddr,
		TaskId:   taskId,
		FileInfo: &uFileInfo,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "download task created",
		Data:       string(respBytes),
	})
}

func UploadHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	filePath := c.Param("path")

	filePath = strings.Join([]string{account.DataDir, filePath}, "/")
	fileDir := filepath.Dir(filePath)
	if err := GetEnv().MkdirAll(fileDir, 0666); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	var req models.UploadRequest

	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    fmt.Sprintf("read request body error: %+v", err),
		})
		return
	}
	if len(jsonBytes) == 0 || json.Unmarshal(jsonBytes, &req) != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "need file info",
		})
		return
	}

	fileInfo := req.FileInfo

	// Check file size
	if fileInfo.Size > consts.FileSizeLimit {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_FileTooLarge,
			Message:    "file is too large",
		})
		return
	}

	// Check account storage capability
	if fileInfo.Size > account.Cap {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_FileTooLarge,
			Message: fmt.Sprintf("no enough capability, free storage: %vMiB, and size of the file: %vMiB",
				(account.Cap-account.Usage)>>20, fileInfo.Size>>20), // Convert Byte to MB
		})
		return
	}

	// Change the modified time
	if err = GetEnv().Chtimes(filePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	taskId := ftm.AddTask(fileInfo, account, UploadTaskType, fileInfo.Size)

	resp := models.UploadResponse{
		NodeAddr: config.IpAddr,
		TaskId:   taskId,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "upload task created",
		Data:       string(respBytes),
	})
}
