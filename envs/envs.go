package envs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/CyDrive/models"
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
	SetFileInfo(name string, fileInfo *models.FileInfo) error
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

func (env *LocalEnv) SetFileInfo(name string, fileInfo *models.FileInfo) error {
	return nil
}
