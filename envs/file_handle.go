package envs

import (
	"io"
	"os"

	"github.com/CyDrive/models"
	"github.com/CyDrive/utils"
)

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

func (l *LocalFile) Truncate(size int64) error {
	return l.file.Truncate(size)
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

type PipeFile struct {
	fileInfo *models.FileInfo
	reader   *io.PipeReader
	writer   *io.PipeWriter
}

func NewPipeFile(fileInfo *models.FileInfo) *PipeFile {
	reader, writer := io.Pipe()
	return &PipeFile{
		fileInfo: fileInfo,
		reader:   reader,
		writer:   writer,
	}
}

func (l *PipeFile) Stat() (*models.FileInfo, error) {
	return l.fileInfo, nil
}

func (l *PipeFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (l *PipeFile) Truncate(size int64) error {
	return nil
}

func (l *PipeFile) Chmod(mode os.FileMode) error {
	return nil
}

func (l *PipeFile) Close() error {
	return l.writer.Close()
}

func (l *PipeFile) Write(p []byte) (n int, err error) {
	return l.writer.Write(p)
}

func (l *PipeFile) Read(p []byte) (n int, err error) {
	return l.reader.Read(p)
}
