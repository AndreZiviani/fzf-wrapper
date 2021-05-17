package fzfwrapper

import (
	"github.com/junegunn/fzf/src/util"
)

const (
	// hardcoded from: https://github.com/junegunn/fzf/blob/764316a53d0eb60b315f0bbcd513de58ed57a876/src/constants.go#L38
	fzfChunkSize int = 100
)

type Token struct {
	text         *util.Chars
	prefixLength int32
}
