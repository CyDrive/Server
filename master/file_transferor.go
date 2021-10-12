package master

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
	log "github.com/sirupsen/logrus"
)

type DataTaskType int32

const (
	DataTaskType_Download DataTaskType = iota
	DataTaskType_Upload

	DataTaskExpireTime int64 = 30 * 60
)

type DataTask struct {
	// filled when the server deliver task id
	Id           int32
	ClientIp     string
	FileInfo     *models.FileInfo
	Account      *models.Account
	StartAt      time.Time
	Type         DataTaskType
	HasDoneBytes int64

	// filled when client connects to the server
	Conn          *net.TCPConn
	LastAcessTime int64
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

		log.Infof("connection from: %+v", conn.RemoteAddr())

		go ft.ProcessConn(conn)
	}
}

func (ft *FileTransferor) CreateTask(clientIp string, fileInfo *models.FileInfo, account *models.Account, taskType DataTaskType, doneBytes int64) int32 {
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
	if task.Type == DataTaskType_Download {
		go ft.DownloadHandle(task)
	} else {
		go ft.UploadHandle(task)
	}
}

func (ft *FileTransferor) DownloadHandle(task *DataTask) {
	path := strings.Join([]string{utils.GetAccountDataDir(task.Account),
		task.FileInfo.FilePath}, "/")
	file, err := GetEnv().Open(path)
	if err != nil {
		log.Errorf("open file %+v error: %+v", task.FileInfo.FilePath, err)
		// todo: notify account by message channel
		return
	}
	defer file.Close()

	if _, err = file.Seek(task.HasDoneBytes, io.SeekStart); err != nil {
		log.Errorf("file seeks to %+v error: %+v", task.HasDoneBytes, err)
	}

	for {
		written, err := io.Copy(task.Conn, file)
		if err != nil {
			if err == io.EOF {
				log.Infof("conn has been closed")
			} else {
				log.Errorf("write conn failed: err=%+v", err)
			}
			break
		}

		task.HasDoneBytes += written
		if task.HasDoneBytes >= task.FileInfo.Size {
			log.Infof("task finished")
			break
		}
	}

	ft.deleteTask(task.Id)
}

func (ft *FileTransferor) UploadHandle(task *DataTask) {
	filePath := filepath.Join(utils.GetAccountDataDir(task.Account), task.FileInfo.FilePath)

	file, err := GetEnv().OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("open file %+v error: %+v", filePath, err)
		// todo: notify account by message channel
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
			if time.Now().Unix()-atomic.LoadInt64(&task.LastAcessTime) >= DataTaskExpireTime {
				tasksShouldBeDeleted = append(tasksShouldBeDeleted, task)
			}

			return true
		})

		log.Infof("task should be dropped: %+v", tasksShouldBeDeleted)
		for _, task := range tasksShouldBeDeleted {
			ft.deleteTask(task.Id)
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
