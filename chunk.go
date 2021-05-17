package fzfwrapper

type ItemBuilder func(*Item, []byte) bool

type Chunk struct {
	items [fzfChunkSize]Item
	count int
}

func (c *Chunk) push(trans ItemBuilder, data []byte) bool {
	if trans(&c.items[c.count], data) {
		c.count++
		return true
	}
	return false
}
