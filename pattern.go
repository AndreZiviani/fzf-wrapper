package fzfwrapper

import (
	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/algo"
)

type Pattern struct {
	chunk   *Chunk
	pattern [][]rune
	sortBy  []criterion
}

func NewPattern(chunk *Chunk, pattern [][]rune, sortBy []criterion) *Pattern {
	return &Pattern{chunk, pattern, sortBy}
}

func (p Pattern) MatchChunk() []Result {
	matches := []Result{}
	items := p.chunk.items
	for idx := 0; idx < p.chunk.count; idx++ {
		match, offsets, pos := p.MatchItem(&items[idx])
		if match != nil {
			match.Offsets = offsets
			match.Pos = pos
			matches = append(matches, *match)
		}
	}

	return matches
}

func (p Pattern) MatchItem(item *Item) (*Result, []fzf.Offset, *[]int) {
	offsets, bonus, pos := p.extendedMatch(item)
	if len(offsets) == len(p.pattern) {
		result := buildResult(item, offsets, bonus, p.sortBy)
		return &result, offsets, pos
	}

	return nil, nil, nil
}

func (p Pattern) extendedMatch(item *Item) ([]fzf.Offset, int, *[]int) {

	input := []Token{{&item.text, 0}}

	offsets := []fzf.Offset{}
	allPos := &[]int{}
	var totalScore int

	for _, term := range p.pattern {
		matched := false
		var offset fzf.Offset
		var currentScore int

		off, score, pos := p.iter(input, term)

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

func (p Pattern) iter(tokens []Token, pattern []rune) (fzf.Offset, int, *[]int) {
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
