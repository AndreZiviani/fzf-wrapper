package fzfwrapper

import (
	"sync"
)

type MyItemBuilder func(*MyItem, []byte) bool

type MyChunkList struct {
	chunks []*MyChunk
	mutex  sync.Mutex
	trans  MyItemBuilder
}

func MyNewChunkList(trans MyItemBuilder) *MyChunkList {
	return &MyChunkList{
		chunks: []*MyChunk{},
		mutex:  sync.Mutex{},
		trans:  trans,
	}
}

func (cl *MyChunkList) lastChunk() *MyChunk {
	return cl.chunks[len(cl.chunks)-1]
}

func (cl *MyChunkList) Push(data []byte) bool {
	cl.mutex.Lock()

	if len(cl.chunks) == 0 || cl.lastChunk().IsFull() {
		cl.chunks = append(cl.chunks, &MyChunk{})
	}

	ret := cl.lastChunk().push(cl.trans, data)
	cl.mutex.Unlock()
	return ret
}

func (cl *MyChunkList) Clear() {
	cl.mutex.Lock()
	cl.chunks = nil
	cl.mutex.Unlock()
}

func CountItems(cs []*MyChunk) int {
	if len(cs) == 0 {
		return 0
	}
	return fzfChunkSize*(len(cs)-1) + cs[len(cs)-1].count
}

func (cl *MyChunkList) Snapshot() ([]*MyChunk, int) {
	cl.mutex.Lock()

	ret := make([]*MyChunk, len(cl.chunks))
	copy(ret, cl.chunks)

	// Duplicate the last chunk
	if cnt := len(ret); cnt > 0 {
		newChunk := *ret[cnt-1]
		ret[cnt-1] = &newChunk
	}

	cl.mutex.Unlock()
	return ret, CountItems(ret)
}

type MyChunk struct {
	items [fzfChunkSize]MyItem
	count int
}

func (c *MyChunk) push(trans MyItemBuilder, data []byte) bool {
	if trans(&c.items[c.count], data) {
		c.count++
		return true
	}
	return false
}

func (c *MyChunk) IsFull() bool {
	return c.count == fzfChunkSize
}
