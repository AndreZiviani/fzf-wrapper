package fzfwrapper

import (
	"sort"
	"strings"

	"github.com/junegunn/fzf/src/util"
)

type InputData interface {
	FzfInputList() []string
	FzfInputLen() int
}

type Wrapper struct {
	inputData InputData
	pattern   [][]rune
	sort      []criterion
}

type Option func(opt *opt)

type opt struct {
	sort []criterion
}

func WithSortBy(c ...criterion) Option {
	return func(o *opt) {
		sort := make([]criterion, 0)
		for _, s := range c {
			sort = append(sort, s)
		}
		o.sort = sort
	}
}

func NewWrapper(options ...Option) *Wrapper {
	opt := opt{[]criterion{}}
	for _, o := range options {
		o(&opt)
	}

	w := Wrapper{
		sort: opt.sort,
	}

	return &w
}

func (w *Wrapper) SetInput(input InputData) {
	w.inputData = input
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

	chunk := newChunk(w.inputData, trans)

	pattern := NewPattern(chunk, w.pattern, w.sort)
	results := pattern.MatchChunk()

	if len(w.sort) > 0 {
		sort.Sort(ByRelevance(results))
	}

	return results, nil
}
