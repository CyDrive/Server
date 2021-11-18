package envs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CyDrive/consts"
	"github.com/CyDrive/master/managers"
	"github.com/CyDrive/models"
	"github.com/CyDrive/network"
	"github.com/CyDrive/types"
	"github.com/CyDrive/utils"
)

type Env interface {
	Open(name string) (types.FileHandle, error)                                 // for read
	OpenFile(name string, flag int, perm os.FileMode) (types.FileHandle, error) // for write
	RemoveAll(path string) error
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(dirname string) ([]*models.FileInfo, error)
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Stat(name string) (*models.FileInfo, error)
}

type LocalEnv struct{}

func NewLocalEnv() *LocalEnv {
	return &LocalEnv{}
}

func (env *LocalEnv) Open(name string) (types.FileHandle, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewLocalFile(file, name), nil
}

func (env *LocalEnv) OpenFile(name string, flag int, perm os.FileMode) (types.FileHandle, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return NewLocalFile(file, name), nil
}

func (env *LocalEnv) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (env *LocalEnv) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (env *LocalEnv) ReadDir(dirname string) ([]*models.FileInfo, error) {
	innerList, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	fileInfoList := []*models.FileInfo{}
	for _, info := range innerList {
		fileInfoList = append(fileInfoList,
			utils.NewFileInfo(info, filepath.Join(dirname, info.Name())))
	}

	return fileInfoList, nil
}

func (env *LocalEnv) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (env *LocalEnv) Stat(name string) (*models.FileInfo, error) {
	inner, err := os.Stat(name)
	if err != nil {
		return &models.FileInfo{}, err
	}

	return utils.NewFileInfo(inner, name), nil
}

type RemoteEnv struct {
	nodeManager    *managers.NodeManager
	fileTransferor *network.FileTransferor
	metaMap        *sync.Map // map: filePath -> *FileInfo or []string
}

func NewRemoteEnv(nodeManager *managers.NodeManager, fileTransferor *network.FileTransferor) *RemoteEnv {
	return &RemoteEnv{
		nodeManager:    nodeManager,
		fileTransferor: fileTransferor,
		metaMap:        &sync.Map{},
	}
}

func (env *RemoteEnv) Open(name string) (types.FileHandle, error) {
	fileInfo, ok := env.getFileInfo(name)
	if !ok {
		return nil, os.ErrNotExist
	}

	node := env.nodeManager.GetNodesByFilePath(name)[0]
	file := NewRemoteFile(os.O_RDONLY, 0666, fileInfo)
	task := env.fileTransferor.CreateTask(node.Addr, fileInfo, file, consts.DataTaskType_Upload, 0)
	task.OnConnect = func() {
		file.conn = task.Conn
	}
	env.nodeManager.PrepareReadFile(task.Id, name)

	return file, nil
}

func (env *RemoteEnv) OpenFile(name string, flag int, perm os.FileMode) (types.FileHandle, error) {
	fileInfo, ok := env.getFileInfo(name)
	if !ok {
		return nil, os.ErrNotExist
	}

	node := env.nodeManager.GetNodesByFilePath(name)[0]
	file := NewRemoteFile(flag, perm, fileInfo)
	task := env.fileTransferor.CreateTask(node.Addr, fileInfo, file, consts.DataTaskType_Download, 0)
	task.OnConnect = func() {
		file.conn = task.Conn
	}
	env.nodeManager.PrepareReadFile(task.Id, name)

	return file, nil
}

func (env *RemoteEnv) Stat(name string) (*models.FileInfo, error) {
	fileInfo, ok := env.getFileInfo(name)
	if !ok {
		return fileInfo, os.ErrNotExist
	}

	return fileInfo, nil
}

func (env *RemoteEnv) ReadDir(dirname string) ([]*models.FileInfo, error) {
	entriesI, ok := env.metaMap.Load(dirname)
	if !ok {
		return nil, os.ErrNotExist
	}

	entries, ok := entriesI.([]string)
	if !ok {
		return nil, os.ErrInvalid
	}

	fileInfoList := []*models.FileInfo{}
	for _, entry := range entries {
		fileInfoI, ok := env.metaMap.Load(entry)
		if !ok {
			panic(fmt.Sprintf("forget to save the file info into metaMap: dirname=%s filepath=%s", dirname, entry))
		}

		fileInfo, ok := fileInfoI.(*models.FileInfo)
		if !ok { // it's a subfolder
			fileInfoList = append(fileInfoList, &models.FileInfo{
				FilePath: entry,
				IsDir:    true,
			})
		} else {
			fileInfoList = append(fileInfoList, fileInfo)
		}
	}

	return fileInfoList, nil
}

func (env *RemoteEnv) getFileInfo(filePath string) (*models.FileInfo, bool) {
	fileInfoI, ok := env.metaMap.Load(filePath)
	if !ok {
		return nil, false
	}

	fileInfo, ok := fileInfoI.(*models.FileInfo)
	return fileInfo, ok
}
