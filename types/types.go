package types

type AccountStoreType = string
type MessageStoreType = string
type EnvType = string
type TaskId = int32

type ReadIndex struct {
	FilePath string
	Offset   int64
	Count    int64
}
