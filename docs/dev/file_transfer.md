# File Transfer Module

## 背景
文件传输是 CyDrive 的一个基础功能，在最初的版本中我们直接在 HTTP 的 request/response body 中存储文件数据，但是这在大文件传输中会非常不可靠，并且我们注意到一些 http client 实现会限制 body 的大小。因此我们需要实现一个文件传输功能，并能够支持断点续传。

## 设计
我们实现一个 FileTransferor 类来提供文件传输功能，并将上传/下载一个文件抽象为一个 DataTask 。以下载文件为例，新的实现的流程是：

1. client 访问 /file/*path 请求下载文件
2. master 创建一个 task，并将这个 task 添加到 FileTransferor 中。然后在 response 中附上 task_id 及 file_info 等信息
3. 建立一个 client 与 master/node 之间的 TCP 连接，client 将 task_id 写入 socket
4. master/node 根据 task_id 找到对应的任务，并将该连接与 task 对应。将相应的文件数据写入 socket

FileTransferor 会监听一个固定端口来接受 client 的连接，并负责管理所有的 task 。FileTransferor 和 DataTask 的定义如下：

```go
type FileTransferor struct {
	taskMap *sync.Map
	idGen   *utils.IdGenerator
}

type DataTask struct {
	// filled when the server deliver task id
	Id        int32
	ClientIp  string
	FileInfo  *models.FileInfo
	Account   *models.Account
	StartAt   time.Time
	Type      TaskType
	DoneBytes int64

	// filled when client connects to the server
	Conn            *net.TCPConn
    LastAccessTime  int64
}
``` 