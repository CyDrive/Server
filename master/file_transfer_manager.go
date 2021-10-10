package master

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/model"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

type TaskType int32

const (
	DownloadTaskType TaskType = iota
	UploadTaskType
)

var (
	ftm = NewFileTransferManager()
)

type Task struct {
	// filled when the server deliver task id
	Id        int32
	FileInfo  *model.FileInfo
	User      *model.Account
	Expire    time.Duration
	StartAt   time.Time
	Type      TaskType
	DoneBytes int64

	// filled when client connects to the server
	Conn *net.TCPConn
}

type FileTransferManager struct {
	taskMap *sync.Map
	idGen   *utils.IdGenerator
}

func NewFileTransferManager() *FileTransferManager {
	idGen := utils.NewIdGenerator()
	return &FileTransferManager{
		taskMap: &sync.Map{},
		idGen:   idGen,
	}
}

func (ftm *FileTransferManager) Listen() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: consts.FtmListenPort})
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("accept tcp connection error: %+v", err)
		}

		go func(conn *net.TCPConn) {
			bufReader := bufio.NewReader(conn)
			var taskId int64
			err := binary.Read(bufReader, binary.LittleEndian, &taskId)
			if err != nil {
				log.Errorf("read task id error: %+v", err)
				return
			}

			taskI, ok := ftm.taskMap.Load(taskId)
			if !ok {
				log.Errorf("task not exist, taskId=%+v", taskId)
				return
			}
			task := taskI.(*Task)
			task.Conn = conn
			if task.Type == DownloadTaskType {
				go ftm.DownloadHandle(task)
			} else {
				go ftm.UploadHandle(task)
			}
		}(conn)
	}
}

func (ftm *FileTransferManager) AddTask(fileInfo *model.FileInfo, user *model.Account, taskType TaskType, doneBytes int64) int32 {
	taskId := ftm.idGen.NextAndRef()
	ftm.taskMap.Store(taskId, &Task{
		Id:        taskId,
		FileInfo:  fileInfo,
		User:      user,
		Expire:    24 * time.Hour,
		StartAt:   time.Now(),
		Type:      taskType,
		DoneBytes: doneBytes,
	})

	return taskId
}

func (ftm *FileTransferManager) DownloadHandle(task *Task) {
	file, err := GetEnv().Open(task.FileInfo.FilePath)
	if err != nil {
		log.Errorf("open file %+v error: %+v", task.FileInfo.FilePath, err)
		// todo: notify user by message channel
		return
	}
	defer file.Close()

	if _, err = file.Seek(task.DoneBytes, io.SeekStart); err != nil {
		log.Errorf("file seeks to %+v error: %+v", task.DoneBytes, err)
	}

	totalBytes := task.DoneBytes
	for {
		written, err := io.Copy(task.Conn, file)
		if err != nil {
			if err != io.EOF {
				log.Errorf("upload failed: err=%+v", err)
			} else {
				log.Infof("upload task finish")
			}

			break
		}

		totalBytes += written
		if totalBytes >= task.FileInfo.Size {
			log.Infof("upload task finish")
			break
		}
	}

	ftm.dropTask(task)
}

func (ftm *FileTransferManager) UploadHandle(task *Task) {
	filePath := filepath.Join(task.User.DataDir, task.FileInfo.FilePath)

	file, err := GetEnv().OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("open file %+v error: %+v", filePath, err)
		// todo: notify user by message channel
		return
	}
	defer file.Close()

	totalBytes := task.DoneBytes
	for {
		written, err := io.Copy(file, task.Conn)
		if err != nil {
			if err != io.EOF {
				log.Errorf("upload failed: err=%+v", err)
			} else {
				log.Infof("upload task finish")
			}

			break
		}

		totalBytes += written
		if totalBytes >= task.FileInfo.Size {
			log.Infof("upload task finish")
			break
		}
	}

	ftm.dropTask(task)
}

func (ftm *FileTransferManager) dropTask(task *Task) {
	task.Conn.Close()
	ftm.taskMap.Delete(task.Id)
	ftm.idGen.UnRef(task.Id)
}
