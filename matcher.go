package fzfwrapper

func MScan(chunks []*MyChunk, pattern [][]rune) [][]MyResult {

	var sliceMatches [][]MyResult

	sliceMatches = make([][]MyResult, 0)

	for _, chunk := range chunks {
		matches := MMatchChunk(chunk, pattern)
		sliceMatches = append(sliceMatches, matches)
	}

	return sliceMatches
}
