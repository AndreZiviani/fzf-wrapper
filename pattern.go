package fzfwrapper

import (
	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/algo"
)

func MMatchChunk(chunk *MyChunk, pattern [][]rune) []MyResult {
	matches := []MyResult{}
	count := len(s)
	for idx := 0; idx < count; idx++ {
		//items := GetUnexportedField(reflect.ValueOf(chunk).Elem().FieldByName("items")).([fzfChunkSize]fzf.Item) //TODO
		//MyItems := NewMyItemList(items)
		items := chunk.items
		match, _, _ := MMatchItem(&items[idx], pattern)
		if match != nil {
			matches = append(matches, *match)
		}
	}

	return matches
}

func MMatchItem(item *MyItem, pattern [][]rune) (*MyResult, []fzf.Offset, *[]int) {
	offsets, bonus, pos := MextendedMatch(item, pattern)
	if len(offsets) == len(pattern) {
		result := MbuildResult(item, offsets, bonus)
		return &result, offsets, pos
	}

	return nil, nil, nil
}

func MextendedMatch(item *MyItem, pattern [][]rune) ([]fzf.Offset, int, *[]int) {

	input := []MyToken{{&item.text, 0}}

	offsets := []fzf.Offset{}
	var totalScore int

	for _, term := range pattern {
		matched := false
		var offset fzf.Offset
		var currentScore int
		off, score, _ := Miter(input, term)
		if sidx := off[0]; sidx >= 0 {
			offset, currentScore = off, score
			matched = true
		}
		if matched {
			offsets = append(offsets, offset)
			totalScore += currentScore
		}
	}

	return offsets, totalScore, nil
}

func Miter(tokens []MyToken, mypattern []rune) (fzf.Offset, int, *[]int) {
	for _, part := range tokens {
		if res, pos := algo.FuzzyMatchV2(false, false, true, part.text, mypattern, false, nil); res.Start >= 0 {
			sidx := int32(res.Start) + part.prefixLength
			eidx := int32(res.End) + part.prefixLength
			if pos != nil {
				for idx := range *pos {
					(*pos)[idx] += int(part.prefixLength)
				}
			}
			return fzf.Offset{sidx, eidx}, res.Score, pos
		}
	}
	return fzf.Offset{-1, -1}, 0, nil
}
