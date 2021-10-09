package consts

const (
	HttpListenPort    = 6454
	HttpListenPortStr = ":6454"
	RpcListenPort     = 6455
	RpcListenPortStr  = ":6455"
	FtmListenPort     = 6456
	FtmListenPortStr  = ":6456"
	UserDataDir       = "user_data"
)

const (
	// The size of file must be not greater than 1GB
	FileSizeLimit int64 = 1 << 30

	// A file with not small than 100MB size should be compressed
	CompressBaseline int64 = 100 << 20
)
