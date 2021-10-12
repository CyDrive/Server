package envs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
)

type FileHandle interface {
	Stat() (*models.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Chmod(mode os.FileMode) error
	Close() error
	io.Writer
	io.Reader
}

type LocalFile struct {
	path string
	file *os.File
}

func NewLocalFile(file *os.File, path string) *LocalFile {
	return &LocalFile{
		path: path,
		file: file,
	}
}

func (l *LocalFile) Stat() (*models.FileInfo, error) {
	inner, err := l.file.Stat()
	if err != nil {
		return &models.FileInfo{}, err
	}

	return utils.NewFileInfo(inner, l.path), nil
}

func (l *LocalFile) Seek(offset int64, whence int) (int64, error) {
	return l.file.Seek(offset, whence)
}

func (l *LocalFile) Chmod(mode os.FileMode) error {
	return l.file.Chmod(mode)
}

func (l *LocalFile) Close() error {
	return l.file.Close()
}

func (l *LocalFile) Write(p []byte) (n int, err error) {
	return l.file.Write(p)
}

func (l *LocalFile) Read(p []byte) (n int, err error) {
	return l.file.Read(p)
}

type Env interface {
	Open(name string) (FileHandle, error)
	OpenFile(name string, flag int, perm os.FileMode) (FileHandle, error)
	MkdirAll(path string, perm os.FileMode) error
	ReadDir(dirname string) ([]*models.FileInfo, error)
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Stat(name string) (*models.FileInfo, error)
}

type LocalEnv struct{}

func NewLocalEnv() *LocalEnv {
	return &LocalEnv{}
}

func (l *LocalEnv) Open(name string) (FileHandle, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return NewLocalFile(file, name), nil
}

func (l *LocalEnv) OpenFile(name string, flag int, perm os.FileMode) (FileHandle, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return NewLocalFile(file, name), nil
}

func (l *LocalEnv) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (l *LocalEnv) ReadDir(dirname string) ([]*models.FileInfo, error) {
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

func (l *LocalEnv) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (l *LocalEnv) Stat(name string) (*models.FileInfo, error) {
	inner, err := os.Stat(name)
	if err != nil {
		return &models.FileInfo{}, err
	}

	return utils.NewFileInfo(inner, name), nil
}
