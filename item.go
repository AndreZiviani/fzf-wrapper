package fzfwrapper

import (
	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

type MyItem struct {
	text        util.Chars
	transformed *[]fzf.Token
	origText    *[]byte
}

func (item *MyItem) TrimLength() uint16 {
	return item.text.TrimLength()
}

// AsString returns the original string
func (item *MyItem) AsString(stripAnsi bool) string {
	if item.origText != nil {
		return string(*item.origText)
	}
	return item.text.ToString()
}
