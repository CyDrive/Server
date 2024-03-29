package network

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

type DataTask struct {
	// filled when the server deliver task id
	Id           types.TaskId
	FileInfo     *models.FileInfo // note: the FilePath here is not relative to the account data folder
	StartAt      time.Time
	Type         consts.DataTaskType
	HasDoneBytes int64

	// filled when client connects to the server
	Conn           *net.TCPConn
	File           types.FileHandle
	LastAccessTime int64

	// callbacks
	OnConnect func()
	OnStart   func()
	OnEnd     func()
	OnError   func()
}

type FileTransferor struct {
	taskMap *sync.Map
	idGen   *utils.IdGenerator
}

func NewFileTransferor() *FileTransferor {
	idGen := utils.NewIdGenerator()
	return &FileTransferor{
		taskMap: &sync.Map{},
		idGen:   idGen,
	}
}

func (ft *FileTransferor) Listen() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: consts.FileTransferorListenPort})
	if err != nil {
		panic(err)
	}

	go ft.GcMaintenance()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("accept tcp connection error: %+v", err)
		}

		go ft.ProcessConn(conn)
	}
}

func (ft *FileTransferor) CreateTask(fileInfo *models.FileInfo, file types.FileHandle, taskType consts.DataTaskType, doneBytes int64) *DataTask {
	taskId := ft.idGen.NextAndRef()
	// host, _, _ := net.SplitHostPort(clientIp)
	task := &DataTask{
		Id:           taskId,
		FileInfo:     fileInfo,
		StartAt:      time.Now(),
		Type:         taskType,
		HasDoneBytes: doneBytes,

		Conn:           nil,
		File:           file,
		LastAccessTime: time.Now().Unix(),
	}

	log.Infof("create new task: %+v", task)
	ft.taskMap.Store(taskId, task)

	return task
}

func (ft *FileTransferor) GetTask(taskId types.TaskId) *DataTask {
	taskI, ok := ft.taskMap.Load(taskId)
	if !ok {
		return nil
	}

	return taskI.(*DataTask)
}

func (ft *FileTransferor) ProcessConn(conn *net.TCPConn) {
	var taskId int32
	err := binary.Read(conn, binary.LittleEndian, &taskId)
	if err != nil {
		log.Errorf("read task id error: %+v", err)
		return
	}

	log.Infof("connection from: %+v, taskId=%+v", conn.RemoteAddr(), taskId)

	taskI, ok := ft.taskMap.Load(taskId)
	if !ok {
		log.Errorf("task not exist, taskId=%+v", taskId)
		return
	}
	task := taskI.(*DataTask)

	task.Conn = conn

	// Connection established from now

	if task.OnConnect != nil {
		task.OnConnect()
	}

	switch task.Type {
	case consts.DataTaskType_Download:
		go ft.DownloadHandle(task)

	case consts.DataTaskType_Upload:
		go ft.UploadHandle(task)
	}
}

func (ft *FileTransferor) DownloadHandle(task *DataTask) {
	var err error

	if _, err = task.File.Seek(task.HasDoneBytes, io.SeekStart); err != nil {
		log.Errorf("file seeks to %+v error: %+v", task.HasDoneBytes, err)
	}

	for {
		log.Infof("reading data from file=%s...", task.FileInfo.FilePath)
		n, err := io.Copy(task.Conn, task.File)
		if err != nil {
			log.Errorf("failed to copy data from file to conn, filePath=%s, err=%v", task.FileInfo.FilePath, err)
			break
		}
		log.Infof("copy %v bytes from file=%s", n, task.FileInfo.FilePath)

		task.HasDoneBytes += n
		task.LastAccessTime = time.Now().Unix()
		if task.HasDoneBytes >= task.FileInfo.Size {
			log.Infof("task finished: task=%+v", task)
			break
		}
	}

	if task.OnEnd != nil {
		task.OnEnd()
	}

	ft.deleteTask(task.Id)
}

func (ft *FileTransferor) UploadHandle(task *DataTask) {
	if err := task.File.Truncate(task.HasDoneBytes); err != nil {
		log.Errorf("failed to truncated file, err=%+v, task=%+v", err, task)
		return
	}

	for {
		log.Infof("reading data to file=%s...", task.FileInfo.FilePath)
		read, err := io.Copy(task.File, task.Conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("conn has been closed")
			} else {
				log.Errorf("read conn fail ed: err=%+v", err)
			}

			break
		}
		log.Infof("copy %v bytes to file=%s", read, task.FileInfo.FilePath)

		task.HasDoneBytes += read
		task.LastAccessTime = time.Now().Unix()
		if task.HasDoneBytes >= task.FileInfo.Size {
			log.Infof("task finished: %+v", task)
			break
		}
	}

	if task.OnEnd != nil {
		task.OnEnd()
	}

	ft.deleteTask(task.Id)
}

func (ft *FileTransferor) GcMaintenance() {
	for {
		tasksShouldBeDeleted := []*DataTask{}
		ft.taskMap.Range(func(key, value interface{}) bool {
			task := value.(*DataTask)

			// No response for a long time
			if time.Now().Unix()-atomic.LoadInt64(&task.LastAccessTime) >= consts.DataTaskExpireTime {
				tasksShouldBeDeleted = append(tasksShouldBeDeleted, task)
			}

			return true
		})

		if len(tasksShouldBeDeleted) > 0 {
			log.Infof("task should be dropped: %+v", tasksShouldBeDeleted)
			for _, task := range tasksShouldBeDeleted {
				ft.deleteTask(task.Id)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (ft *FileTransferor) deleteTask(taskId int32) {
	if taskI, ok := ft.taskMap.LoadAndDelete(taskId); ok {
		task := taskI.(*DataTask)
		if task.Conn != nil {
			task.Conn.Close()
		}
		ft.idGen.UnRef(taskId)
	}
}
