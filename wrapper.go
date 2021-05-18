package fzfwrapper

import (
	"sort"
	"strings"

	"github.com/junegunn/fzf/src/util"
)

type WrapperResult struct {
	Text  string
	Score [4]uint16
}

type Wrapper struct {
	InputList []string
	Pattern   [][]rune
	Results   []WrapperResult
	sort      criterion
}

type Option func(opt *opt)

type opt struct {
	sort criterion
}

func WithSortBy(c criterion) Option {
	return func(o *opt) {
		o.sort = c
	}
}

func NewWrapper(input []string, pattern string, options ...Option) *Wrapper {
	opt := opt{-1}
	for _, o := range options {
		o(&opt)
	}

	patternSlice := strings.Split(pattern, " ")

	w := Wrapper{
		InputList: input,
		Pattern:   make([][]rune, 0, len(patternSlice)),
		Results:   make([]WrapperResult, 0),
		sort:      opt.sort,
	}

	for _, ps := range patternSlice {
		w.Pattern = append(w.Pattern, []rune(ps))
	}

	return &w
}

func (w *Wrapper) Fuzzy() ([]Result, error) {
	var itemIndex int32
	trans := (func(item *Item, data []byte) bool {

		item.text = util.ToChars(data)
		item.text.TrimTrailingWhitespaces()
		item.text.Index = itemIndex

		itemIndex++
		return true

	})

	chunk := newChunk(w.InputList, trans)

	merger := MatchChunk(chunk, w.Pattern)

	if w.sort != -1 {
		switch w.sort {
		case ByLength:
			sort.Sort(ByRelevance(merger))
		}
	}

	return merger, nil
}
