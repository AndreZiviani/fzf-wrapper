package fzfwrapper

import (
	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/algo"
)

func MatchChunk(chunk *Chunk, pattern [][]rune) []Result {
	matches := []Result{}
	items := chunk.items
	for idx := 0; idx < chunk.count; idx++ {
		match, offsets, pos := MatchItem(&items[idx], pattern)
		if match != nil {
			match.Offsets = offsets
			match.Pos = pos
			matches = append(matches, *match)
		}
	}

	return matches
}

func MatchItem(item *Item, pattern [][]rune) (*Result, []fzf.Offset, *[]int) {
	offsets, bonus, pos := extendedMatch(item, pattern)
	if len(offsets) == len(pattern) {
		result := buildResult(item, offsets, bonus)
		return &result, offsets, pos
	}

	return nil, nil, nil
}

func extendedMatch(item *Item, pattern [][]rune) ([]fzf.Offset, int, *[]int) {

	input := []Token{{&item.text, 0}}

	offsets := []fzf.Offset{}
	allPos := &[]int{}
	var totalScore int

	for _, term := range pattern {
		matched := false
		var offset fzf.Offset
		var currentScore int

		off, score, pos := iter(input, term)

		if sidx := off[0]; sidx >= 0 {
			offset, currentScore = off, score
			matched = true

			if pos != nil {
				*allPos = append(*allPos, *pos...)
			} else {
				for idx := off[0]; idx < off[1]; idx++ {
					*allPos = append(*allPos, int(idx))
				}
			}
		}

		if matched {
			offsets = append(offsets, offset)
			totalScore += currentScore
		}
	}

	return offsets, totalScore, allPos
}

func iter(tokens []Token, pattern []rune) (fzf.Offset, int, *[]int) {
	for _, part := range tokens {
		if res, pos := algo.FuzzyMatchV2(false, true, true, part.text, pattern, true, nil); res.Start >= 0 {
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
