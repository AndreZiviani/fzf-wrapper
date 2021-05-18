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
	inputList []string
	pattern   [][]rune
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

func NewWrapper(options ...Option) *Wrapper {
	opt := opt{-1}
	for _, o := range options {
		o(&opt)
	}

	w := Wrapper{
		Results: make([]WrapperResult, 0),
		sort:    opt.sort,
	}

	return &w
}

func (w *Wrapper) SetInput(input []string) {
	w.inputList = input
}

func (w *Wrapper) SetPattern(pattern string) {
	patternSlice := strings.Split(pattern, " ")

	w.pattern = make([][]rune, 0, len(patternSlice))

	for _, ps := range patternSlice {
		w.pattern = append(w.pattern, []rune(ps))
	}

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

	chunk := newChunk(w.inputList, trans)

	merger := MatchChunk(chunk, w.pattern)

	if w.sort != -1 {
		switch w.sort {
		case ByLength:
			sort.Sort(ByRelevance(merger))
		}
	}

	return merger, nil
}
