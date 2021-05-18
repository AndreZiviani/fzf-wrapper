package fzfwrapper

import (
	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

type Item struct {
	text        util.Chars
	transformed *[]fzf.Token
	origText    *[]byte
}

func (item *Item) TrimLength() uint16 {
	return item.text.TrimLength()
}

// AsString returns the original string
func (item *Item) AsString(stripAnsi bool) string {
	if item.origText != nil {
		return string(*item.origText)
	}
	return item.text.ToString()
}

func (item *Item) Index() int32 {
	return item.text.Index
}
