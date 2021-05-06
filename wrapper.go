package fzfwrapper

import (
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
}

func NewWrapper(input []string, pattern string) *Wrapper {
	patternSlice := strings.Split(pattern, " ")

	w := Wrapper{
		InputList: input,
		Pattern:   make([][]rune, 0, len(patternSlice)),
		Results:   make([]WrapperResult, 0),
	}

	for _, ps := range patternSlice {
		w.Pattern = append(w.Pattern, []rune(ps))
	}

	return &w
}

func (w *Wrapper) Fuzzy() (bool, error) {
	var itemIndex int32
	chunkList := MyNewChunkList(func(item *MyItem, data []byte) bool {

		item.text = util.ToChars(data)
		item.text.TrimTrailingWhitespaces()
		item.text.Index = itemIndex

		itemIndex++
		return true

	})

	for _, str := range w.InputList {
		chunkList.Push([]byte(str))
	}

	snapshot, _ := chunkList.Snapshot()

	merger := MScan(snapshot, w.Pattern)

	for _, list := range merger {
		for _, v := range list {
			result := WrapperResult{
				Text:  v.item.AsString(false),
				Score: v.points,
			}
			w.Results = append(w.Results, result)
		}
	}
	return true, nil
}
