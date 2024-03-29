package envs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/CyDrive/consts"
	. "github.com/CyDrive/envs"
	"github.com/CyDrive/master/managers"
	"github.com/CyDrive/models"
	"github.com/CyDrive/network"
	"github.com/CyDrive/types"
)

var (
	_ Env = (*LocalEnv)(nil)
	_ Env = (*RemoteEnv)(nil)
)

type RemoteEnv struct {
	nodeManager    *managers.NodeManager
	fileTransferor *network.FileTransferor
	metaMap        *sync.Map // map: filePath -> *FileInfo or *[]string
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

	file := NewRemoteFile(os.O_RDONLY, 0666, fileInfo)
	task := env.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Upload, 0)
	task.OnEnd = func() {
		file.Close()
	}
	err := env.nodeManager.PrepareReadFile(task.Id, name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (env *RemoteEnv) OpenFile(name string, flag int, perm os.FileMode) (types.FileHandle, error) {
	fileInfo, ok := env.getFileInfo(name)
	if !ok {
		panic("forget to update metaMap?")
	}

	file := NewRemoteFile(flag, perm, fileInfo)
	task := env.fileTransferor.CreateTask(fileInfo, file, consts.DataTaskType_Download, 0)
	task.OnEnd = func() {
		file.Close()

		dir := filepath.Dir(name)
		entriesI, _ := env.metaMap.Load(dir)
		subFolders := entriesI.(*[]string)
		*subFolders = append(*subFolders, name)
		fileInfo, _ := file.Stat()
		env.SetFileInfo(name, fileInfo)
	}
	err := env.nodeManager.PrepareWriteFile(task.Id, name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (env *RemoteEnv) RemoveAll(path string) error {
	env.nodeManager.NotifyDeleteFile(path)
	env.RemoveAllMeta(path)

	return env.removeEntryFromFolder(path)
}

func (env *RemoteEnv) RemoveAllMeta(path string) error {
	entriesI, ok := env.metaMap.Load(path)
	if !ok {
		return nil
	}

	entries, ok := entriesI.([]string)
	env.metaMap.Delete(path)
	if !ok { // this is a file
		return nil
	}
	for _, entry := range entries {
		env.RemoveAllMeta(entry)
	}

	return nil
}

func (env *RemoteEnv) removeEntryFromFolder(path string) error {
	dir := filepath.Dir(path)
	entriesI, ok := env.metaMap.Load(dir)
	if !ok {
		return os.ErrNotExist
	}

	entries, ok := entriesI.(*[]string)
	if !ok {
		return os.ErrInvalid
	}

	for i, entry := range *entries {
		if entry == path {
			*entries = append((*entries)[:i], (*entries)[i+1:]...)
			break
		}
	}

	return nil
}

func (env *RemoteEnv) MkdirAll(path string, perm os.FileMode) error {
	for path != "." { //
		_, exist := env.metaMap.LoadOrStore(path, &[]string{})
		if exist {
			return nil
		}

		path = filepath.Dir(path)
	}

	return nil
}

func (env *RemoteEnv) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func (env *RemoteEnv) Stat(name string) (models.FileInfo, error) {
	fileInfo, ok := env.getFileInfo(name)
	if !ok {
		return models.FileInfo{}, os.ErrNotExist
	}

	return *fileInfo, nil
}

func (env *RemoteEnv) ReadDir(dirname string) ([]models.FileInfo, error) {
	entriesI, ok := env.metaMap.Load(dirname)
	if !ok {
		return nil, os.ErrNotExist
	}

	entries, ok := entriesI.(*[]string)
	if !ok {
		return nil, os.ErrInvalid
	}

	fileInfoList := []models.FileInfo{}
	for _, entry := range *entries {
		fileInfoI, ok := env.metaMap.Load(entry)
		if !ok {
			panic(fmt.Sprintf("forget to save the file info into metaMap: dirname=%s filepath=%s", dirname, entry))
		}

		fileInfo, ok := fileInfoI.(*models.FileInfo)
		if !ok { // it's a subfolder
			fileInfoList = append(fileInfoList, models.FileInfo{
				FilePath: entry,
				IsDir:    true,
			})
		} else {
			fileInfoList = append(fileInfoList, *fileInfo)
		}
	}

	return fileInfoList, nil
}

func (env *RemoteEnv) SetFileInfo(name string, fileInfo models.FileInfo) error {
	err := env.MkdirAll(filepath.Dir(name), 0666)
	if err != nil {
		return err
	}

	_, ok := env.metaMap.Load(name)
	isNewEntry := !ok

	if fileInfo.IsDir {
		env.metaMap.LoadOrStore(name, &[]string{})
	} else {
		env.metaMap.Store(name, fileInfo)
	}

	if isNewEntry {
		dir := filepath.Dir(name)
		if dir != "." {
			entriesI, ok := env.metaMap.Load(dir)
			if !ok {
				panic("forget to mkdir for this folder: " + dir + ", the filepath is " + name)
			}

			entries, ok := entriesI.(*[]string)
			if !ok {
				panic("not a folder: " + dir)
			}

			*entries = append(*entries, name)
		}
	}

	return nil
}

func (env *RemoteEnv) getFileInfo(filePath string) (*models.FileInfo, bool) {
	fileInfoI, ok := env.metaMap.Load(filePath)
	if !ok {
		return nil, false
	}

	fileInfo, ok := fileInfoI.(*models.FileInfo)
	return fileInfo, ok
}

func (env *RemoteEnv) addToDir(filePath string) {
	_, ok := env.metaMap.Load(filePath)

	// The file is already in the dirs
	if ok {
		return
	}

	dir := filepath.Dir(filePath)

	entriesI, ok := env.metaMap.Load(dir)
	if !ok {
		entriesI = &[]string{}
	}
	entries := entriesI.(*[]string)

	*entries = append(*entries, filePath)

	env.metaMap.Store(dir, entries)
}
