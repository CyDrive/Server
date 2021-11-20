package envs

import (
	"io"
	"os"

	"github.com/CyDrive/models"
	"github.com/CyDrive/types"
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

type RemoteFile struct {
	Flag        int
	Perm        os.FileMode
	FileInfo    *models.FileInfo
	CallOnStart func(taskId types.TaskId)

	// pipe
	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewRemoteFile(flag int, perm os.FileMode, fileInfo *models.FileInfo) *RemoteFile {
	reader, writer := io.Pipe()
	return &RemoteFile{
		Flag:     flag,
		Perm:     perm,
		FileInfo: fileInfo,

		reader: reader,
		writer: writer,
	}
}

func (file *RemoteFile) Stat() (*models.FileInfo, error) {
	return file.FileInfo, nil
}

func (file *RemoteFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (file *RemoteFile) Truncate(size int64) error {
	return nil
}

// unimplemented
func (file *RemoteFile) Chmod(mode os.FileMode) error {
	return nil
}

func (file *RemoteFile) Close() error {
	// var err error = nil
	// if err = file.writer.Close(); err != nil {
	// 	return err
	// }
	// err = file.reader.Close()
	// return err
	return file.writer.Close()
}

// write the data from node to the buffer
// the err is always nil
func (file *RemoteFile) Write(p []byte) (n int, err error) {
	return file.writer.Write(p)
}

func (file *RemoteFile) Read(p []byte) (n int, err error) {
	return file.reader.Read(p)
	// if file.cacheFile == nil {
	// }
	// n, err = file.cacheFile.Read(p)

	// // we think of the err = io.EOF as err = nil
	// // and always return the file.Err if there're both errors
	// // +--------+----------+--------+
	// // |  err   | file.Err | return |
	// // +--------+----------+--------+
	// // | nil    | nil      | nil    |
	// // | io.EOF | nil      | nil    |
	// // | io.EOF | error    | error  |
	// // | error  | nil      | error  |
	// // | error1 | error2   | error2 |
	// // +--------+----------+--------+
	// if err == io.EOF {
	// 	err = nil
	// }
	// if file.Err != nil {
	// 	err = file.Err
	// }

	// return n, err
}
