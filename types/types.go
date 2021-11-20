package types

import (
	"io"
	"os"

	"github.com/CyDrive/models"
)

type EnvType = string
type TaskId = int32

type ReadIndex struct {
	FilePath string
	Offset   int64
	Count    int64
}
type FileHandle interface {
	Stat() (*models.FileInfo, error)
	Truncate(size int64) error
	Chmod(mode os.FileMode) error
	io.Writer
	io.Reader
	io.Seeker
	io.Closer
}
