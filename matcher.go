package fzfwrapper

func Scan(chunks *Chunk, pattern [][]rune) [][]Result {

	var sliceMatches [][]Result

	sliceMatches = make([][]Result, 0)

	matches := MatchChunk(chunks, pattern)
	sliceMatches = append(sliceMatches, matches)

	return sliceMatches
}
