package fzfwrapper

type ItemBuilder func(*Item, []byte) bool

type Chunk struct {
	items []Item
	count int
}

func newChunk(input []string, trans ItemBuilder) *Chunk {
	size := len(input)
	chunk := Chunk{}
	chunk.items = make([]Item, size)

	for _, str := range input {
		chunk.push(trans, []byte(str))
	}

	return &chunk
}

func (c *Chunk) push(trans ItemBuilder, data []byte) bool {
	if trans(&c.items[c.count], data) {
		c.count++
		return true
	}
	return false
}
