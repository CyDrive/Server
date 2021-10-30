package network

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/envs"
	"github.com/CyDrive/models"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

type DataTask struct {
	// filled when the server deliver task id
	Id           types.TaskId
	ClientIp     string
	FileInfo     *models.FileInfo
	Account      *models.Account
	StartAt      time.Time
	Type         types.DataTaskType
	HasDoneBytes int64

	// filled when client connects to the server
	Conn          *net.TCPConn
	LastAcessTime int64

	// filled when the task starts
	FileHandle envs.FileHandle
}

type FileTransferor struct {
	taskMap *sync.Map
	idGen   *utils.IdGenerator
	env     envs.Env
}

func NewFileTransferor(env envs.Env) *FileTransferor {
	idGen := utils.NewIdGenerator()
	return &FileTransferor{
		taskMap: &sync.Map{},
		idGen:   idGen,
		env:     env,
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

		log.Infof("connection from: %+v", conn.RemoteAddr())

		go ft.ProcessConn(conn)
	}
}

func (ft *FileTransferor) CreateTask(clientIp string, fileInfo *models.FileInfo, account *models.Account, taskType types.DataTaskType, doneBytes int64) int32 {
	taskId := ft.idGen.NextAndRef()
	// host, _, _ := net.SplitHostPort(clientIp)
	task := &DataTask{
		Id:           taskId,
		ClientIp:     clientIp,
		FileInfo:     fileInfo,
		Account:      account,
		StartAt:      time.Now(),
		Type:         taskType,
		HasDoneBytes: doneBytes,

		Conn:          nil,
		LastAcessTime: time.Now().Unix(),
	}

	log.Infof("create new task: %+v", task)
	ft.taskMap.Store(taskId, task)

	return taskId
}

func (ft *FileTransferor) GetTask(taskId types.TaskId) *DataTask {
	taskI, ok := ft.taskMap.Load(taskId)
	if !ok {
		return nil
	}

	return taskI.(*DataTask)
}

func (ft *FileTransferor) ProcessConn(conn *net.TCPConn) {
	bufReader := bufio.NewReader(conn)
	var taskId int32
	err := binary.Read(bufReader, binary.LittleEndian, &taskId)
	if err != nil {
		log.Errorf("read task id error: %+v", err)
		return
	}

	taskI, ok := ft.taskMap.Load(taskId)
	if !ok {
		log.Errorf("task not exist, taskId=%+v", taskId)
		return
	}
	task := taskI.(*DataTask)

	// validate
	tcpHost, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	tcpIp := net.ParseIP(tcpHost)
	taskClientIp := net.ParseIP(task.ClientIp)

	if !(tcpIp.Equal(taskClientIp) ||
		tcpIp.IsLoopback() && taskClientIp.IsLoopback()) {
		log.Warnf("IPs not match, tcpIp=%b, taskClientIp=%b", tcpIp, taskClientIp)
		conn.Write([]byte("please DO NOT try to steal data"))
		conn.Close()
		return
	}

	task.Conn = conn

	switch task.Type {
	case consts.DataTaskType_Download:
		go ft.DownloadHandle(task)

	case consts.DataTaskType_Upload:
		go ft.UploadHandle(task)
	}
}

func (ft *FileTransferor) DownloadHandle(task *DataTask) {
	var err error

	path := utils.AccountFilePath(task.Account, task.FileInfo.FilePath)
	task.FileHandle, err = ft.env.Open(path)

	if remoteFileHandle, ok := task.FileHandle.(*envs.RemoteFile); ok && remoteFileHandle.CallOnStart != nil {
		remoteFileHandle.CallOnStart(task.Id)
	}

	if err != nil {
		log.Errorf("open file %+v error: %+v", task.FileInfo.FilePath, err)
		// todo: notify account by message channel
		return
	}
	defer task.FileHandle.Close()

	if _, err = task.FileHandle.Seek(task.HasDoneBytes, io.SeekStart); err != nil {
		log.Errorf("file seeks to %+v error: %+v", task.HasDoneBytes, err)
	}

	for {
		written, err := io.Copy(task.Conn, task.FileHandle)
		if err != nil {
			if err == io.EOF {
				log.Infof("conn has been closed")
			} else {
				log.Errorf("write conn failed: err=%+v", err)
			}
			break
		}

		task.HasDoneBytes += written
		task.LastAcessTime = time.Now().Unix()
		if task.HasDoneBytes >= task.FileInfo.Size {
			log.Infof("task finished: task=%+v", task)
			break
		}
	}

	ft.deleteTask(task.Id)
}

func (ft *FileTransferor) UploadHandle(task *DataTask) {
	filePath := utils.AccountFilePath(task.Account, task.FileInfo.FilePath)

	file, err := ft.env.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Errorf("open file %+v error: %+v", filePath, err)
		// todo: notify account by message channel
		return
	}
	if err = file.Truncate(task.HasDoneBytes); err != nil {
		log.Errorf("failed to truncated file, err=%+v, task=%+v", err, task)
		return
	}

	defer file.Close()

	for {
		read, err := io.Copy(file, task.Conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("conn has been closed")
			} else {
				log.Errorf("read conn failed: err=%+v", err)
			}

			break
		}

		task.HasDoneBytes += read
		task.LastAcessTime = time.Now().Unix()
		if task.HasDoneBytes >= task.FileInfo.Size {
			log.Infof("task finished: %+v", task)
			break
		}
	}

	ft.deleteTask(task.Id)
}

func (ft *FileTransferor) GcMaintenance() {
	for {
		tasksShouldBeDeleted := []*DataTask{}
		ft.taskMap.Range(func(key, value interface{}) bool {
			task := value.(*DataTask)

			// No response for a long time
			if time.Now().Unix()-atomic.LoadInt64(&task.LastAcessTime) >= consts.DataTaskExpireTime {
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
