package master

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/managers"
	"github.com/CyDrive/master/store"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Account services
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

	err = utils.GetJsonDecoder().Unmarshal(reqBytes, &req)
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

	safeAccountBytes, _ := utils.GetJsonEncoder().Marshal(utils.PackSafeAccount(account))

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

	err = utils.GetJsonDecoder().Unmarshal(reqBytes, &req)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    "error when unmarshal request body: " + err.Error(),
		})
		return
	}
	log.Infof("req=%+v", req)

	account, err := GetAccountStore().GetAccountByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "no such account",
		})
		return
	}
	log.Infof("account=%+v", account)

	if account.Password != req.Password {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    "account name or password not correct",
		})
		return
	}

	userSession := sessions.DefaultMany(c, "account")

	userSession.Set("userStruct", account)
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

	safeAccountBytes, err := utils.GetJsonEncoder().Marshal(safeAccount)
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

func GetAccountInfo(c *gin.Context) {
	accountI, _ := c.Get("account")
	account := accountI.(*models.Account)

	safeAccount := utils.PackSafeAccount(account)

	utils.ResponseData(c, consts.StatusCode_Ok, "get account info done", safeAccount)
}

// Storage services
func ListHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	path := c.Param("path")
	path = strings.Trim(path, "/")
	absPath := utils.AccountFilePath(account, path)
	absPath = strings.Trim(absPath, "/")

	fileList, err := GetEnv().ReadDir(absPath)
	for i := range fileList {
		fileList[i].FilePath = strings.ReplaceAll(fileList[i].FilePath, "\\", "/")
		fileList[i].FilePath = strings.TrimPrefix(fileList[i].FilePath, utils.GetAccountDataDir(account.Id))
	}
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	fileListJson, err := utils.GetJsonEncoder().Marshal(&models.FileInfoList{
		FileInfoList: fileList,
	})
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

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	// relative path
	filePath := c.Param("path")

	// absolute filepath
	filePath = utils.AccountFilePath(account, filePath)
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

	file, err := GetEnv().Open(fileInfo.FilePath)
	if err != nil {
		log.Errorf("failed to open file %s, err=%+v", fileInfo.FilePath, err)
		utils.Response(c, consts.StatusCode_IoError, err.Error())
		return
	}

	uFileInfo := fileInfo
	uFileInfo.FilePath, _ = filepath.Rel(utils.GetAccountDataDir(account.Id), uFileInfo.FilePath)
	uFileInfo.FilePath = strings.ReplaceAll(uFileInfo.FilePath, "\\", "/")

	log.Infof("clientIp=%+v", c.ClientIP())
	task := GetFileTransferor().CreateTask(uFileInfo, file, consts.DataTaskType_Download, begin)
	resp := models.DownloadResponse{
		NodeAddr: config.IpAddr + consts.FileTransferorListenPortStr,
		TaskId:   task.Id,
		FileInfo: uFileInfo,
	}
	log.Infof("downloadResp=%+v", resp)
	respBytes, err := utils.GetJsonEncoder().Marshal(&resp)
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
	filePath = utils.AccountFilePath(account, filePath)

	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    fmt.Sprintf("read request body error: %+v", err),
		})
		return
	}

	var req models.UploadRequest
	if err = utils.GetJsonDecoder().Unmarshal(jsonBytes, &req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	fileInfo := req.FileInfo

	fileDir := filePath
	if !fileInfo.IsDir {
		fileDir = filepath.Dir(fileDir)
	}
	if err := GetEnv().MkdirAll(fileDir, 0666); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}
	if fileInfo.IsDir {
		utils.Response(c, consts.StatusCode_Ok, "mkdir done")
		return
	}

	// Check account storage capability
	if account.Usage+fileInfo.Size > account.Cap {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_FileTooLarge,
			Message: fmt.Sprintf("no enough cap, free storage: %vMiB, and size of the file: %vMiB",
				(account.Cap-account.Usage)>>20, fileInfo.Size>>20), // Convert Byte to MB
		})
		return
	}

	flags := os.O_CREATE | os.O_WRONLY
	if req.ShouldTruncate {
		flags |= os.O_TRUNC
	}
	GetEnv().SetFileInfo(filePath, fileInfo)
	file, err := GetEnv().OpenFile(filePath, flags, 0666)
	if err != nil {
		utils.Response(c, consts.StatusCode_IoError, fmt.Sprintf("failed to open file %s, err=%v", filePath, err))
		return
	}
	task := GetFileTransferor().CreateTask(fileInfo, file, consts.DataTaskType_Upload, fileInfo.Size)

	GetAccountStore().AddUsage(account.Email, fileInfo.Size)

	offset := int64(0)

	if !req.ShouldTruncate {
		existFileInfo, err := GetEnv().Stat(filePath)
		if err != nil {
			offset = existFileInfo.Size
		}
	}

	resp := models.UploadResponse{
		NodeAddr: config.IpAddr + consts.FileTransferorListenPortStr,
		TaskId:   task.Id,
		Offset:   offset,
	}

	respBytes, err := utils.GetJsonEncoder().Marshal(&resp)
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

func DeleteHandle(c *gin.Context) {
	userI, _ := c.Get("account")
	account := userI.(*models.Account)

	filePath := c.Param("path")
	filePath = utils.AccountFilePath(account, filePath)

	fileInfo, err := GetEnv().Stat(filePath)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    err.Error(),
		})
		return
	}

	err = GetEnv().RemoveAll(filePath)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}
	GetAccountStore().AddUsage(account.Email, -fileInfo.Size)

	resp := models.DeleteResponse{
		FileInfo: fileInfo,
	}
	respBytes, err := utils.GetJsonEncoder().Marshal(&resp)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "deleted",
		Data:       string(respBytes),
	})
}

// Message service
var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ConnectMessageServiceHandle(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errMsg := fmt.Sprintf("failed to establish websocket connection, err=%+v", err)
		log.Errorf(errMsg)
		utils.Response(c, consts.StatusCode_InternalError, errMsg)
		return
	}

	accountI, _ := c.Get("account")
	account := accountI.(*models.Account)

	deviceId := c.Query("device_id")
	if err != nil || deviceId == "" {
		utils.Response(c, consts.StatusCode_InvalidParameters, "invalid device_id")
		return
	}

	hub := GetMessageManager().GetHub(account.Id)
	hub.Register(
		managers.NewMessageConn(hub, deviceId, conn))
}

// queries: int64 last_time: ms since epoch, int32 count, string device_id
func GetMessageHandle(c *gin.Context) {
	accountI, _ := c.Get("account")
	account := accountI.(*models.Account)

	timeStr := c.Query("last_time")
	countStr := c.Query("count")
	deviceId := c.Query("device_id")

	var (
		count    int64 = 10
		lastTime int64 = time.Now().Unix()
		err      error
	)

	if len(timeStr) > 0 {
		lastTime, err = strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			utils.Response(c,
				consts.StatusCode_InvalidParameters,
				fmt.Sprintf("failed to parse parameter time=%+v, err=%+v", timeStr, err))
			return
		}
	}

	if len(countStr) > 0 {
		count, err = strconv.ParseInt(countStr, 10, 32)
		if err != nil {
			utils.Response(c,
				consts.StatusCode_InvalidParameters,
				fmt.Sprintf("failed to parse parameter count=%+v, err=%+v", countStr, err))
			return
		}
	}

	messages := GetMessageManager().
		GetMessageStore().
		GetMessagesByTime(account.Id, deviceId,
			int32(count),
			time.Unix(0, 0).Add(time.Duration(lastTime)*time.Millisecond))

	resp := models.GetMessageResponse{
		Messages: messages,
	}

	utils.ResponseData(c,
		consts.StatusCode_Ok,
		"get message ok",
		&resp)
}

func ShareHandle(c *gin.Context) {
	accountI, _ := c.Get("account")
	account := accountI.(*models.Account)

	var req models.ShareRequest

	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_IoError,
			Message:    fmt.Sprintf("read request body error: %+v", err),
		})
		return
	}

	if err = utils.GetJsonDecoder().Unmarshal(jsonBytes, &req); err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	var share_link = &store.ShareLink{
		FilePath:        req.FilePath,
		From:            account.GetId(),
		Password:        req.Password,
		LeftAccessCount: req.AccessCount,
		Expire:          req.Expire,
		CreatedAt:       time.Now(),
	}
	err = GetShareStore().CreateShareLink(share_link, req.To...)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	resp := models.ShareResponse{
		Link: share_link.Uri,
	}

	respBytes, err := utils.GetJsonEncoder().Marshal(&resp)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_InternalError,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		StatusCode: consts.StatusCode_Ok,
		Message:    "share link created",
		Data:       string(respBytes),
	})

}

func GetSharedFileHandle(c *gin.Context) {
	accountI, _ := c.Get("account")
	account := accountI.(*models.Account)

	uri := c.Param("uri")
	password := c.GetHeader("password")
	filePath, err := GetShareStore().CheckPermission(uri, account.GetId(), password)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: consts.StatusCode_AuthError,
			Message:    err.Error(),
		})
		return
	}
	// absolute filepath
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

	uFileInfo := fileInfo
	file, _ := GetEnv().Open(fileInfo.FilePath)
	uFileInfo.FilePath, _ = filepath.Rel(utils.GetAccountDataDir(account.Id), uFileInfo.FilePath)
	uFileInfo.FilePath = strings.ReplaceAll(uFileInfo.FilePath, "\\", "/")

	log.Infof("clientIp=%+v", c.ClientIP())
	task := GetFileTransferor().CreateTask(uFileInfo, file, consts.DataTaskType_Download, begin)
	resp := models.DownloadResponse{
		NodeAddr: config.IpAddr + consts.FileTransferorListenPortStr,
		TaskId:   task.Id,
		FileInfo: uFileInfo,
	}
	log.Infof("downloadResp=%+v", resp)
	respBytes, err := utils.GetJsonEncoder().Marshal(&resp)
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
