package fzfwrapper

func sliceChunks(chunks []*MyChunk) [][]*MyChunk {
	perSlice := len(chunks) / partitions

	if perSlice == 0 {
		partitions = len(chunks)
		perSlice = 1
	}

	slices := make([][]*MyChunk, partitions)
	for i := 0; i < partitions; i++ {
		start := i * perSlice
		end := start + perSlice
		if i == partitions-1 {
			end = len(chunks)
		}
		slices[i] = chunks[start:end]
	}
	return slices
}

func MScan(chunks []*MyChunk, pattern [][]rune) [][]MyResult {

	slices := sliceChunks(chunks)

	var sliceMatches [][]MyResult

	for _, chunks := range slices {
		sliceMatches = make([][]MyResult, 0)

		for _, chunk := range chunks {
			matches := MMatchChunk(chunk, pattern)
			sliceMatches = append(sliceMatches, matches)
		}
	}

	allMatches := make([][]MyResult, 0)
	for _, slice := range sliceMatches {
		allMatches = append(allMatches, slice)
	}

	return allMatches //era fzf.NewMerger(nil, allMatches, false, false)
}
