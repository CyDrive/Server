package envs

import (
	"io"
	"os"

	"github.com/CyDrive/models"
	"github.com/CyDrive/types"
)

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

func (file *RemoteFile) Stat() (models.FileInfo, error) {
	return *file.FileInfo, nil
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
	return file.writer.Close()
}

// write the data from node to the buffer
// the err is always nil
func (file *RemoteFile) Write(p []byte) (n int, err error) {
	return file.writer.Write(p)
}

func (file *RemoteFile) Read(p []byte) (n int, err error) {
	return file.reader.Read(p)
}
