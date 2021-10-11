package utils

import (
	"math"
	"sync"
	"sync/atomic"
)

type IdGenerator struct {
	minId int32
	maxId int32

	currentId int32

	idRefMap *sync.Map
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		minId:     1,
		maxId:     math.MaxInt32,
		currentId: 1,

		idRefMap: &sync.Map{},
	}
}

func (idGen *IdGenerator) SetMinId(minId int32) {
	atomic.StoreInt32(&idGen.minId, minId)

	for {
		oldId := atomic.LoadInt32(&idGen.currentId)
		if oldId < minId &&
			atomic.CompareAndSwapInt32(&idGen.currentId, oldId, minId) {
			break
		}
	}
}

func (idGen *IdGenerator) SetMaxId(maxId int32) {
	atomic.StoreInt32(&idGen.maxId, maxId)

	for {
		oldId := atomic.LoadInt32(&idGen.currentId)
		if oldId > maxId &&
			atomic.CompareAndSwapInt32(&idGen.currentId, oldId, idGen.minId) {
			break
		}
	}
}

func (idGen *IdGenerator) Next() int32 {
	var id int32
	for {
		id = atomic.AddInt32(&idGen.currentId, 1) - 1
		if _, ok := idGen.idRefMap.Load(id); !ok {
			break
		}
	}
	return id
}

func (idGen *IdGenerator) NextAndRef() int32 {
	id := idGen.Next()
	idGen.Ref(id)
	return id
}

func (idGen *IdGenerator) Ref(id int32) {
	idGen.idRefMap.Store(id, true)
}

func (idGen *IdGenerator) UnRef(id int32) {
	idGen.idRefMap.Delete(id)
}
