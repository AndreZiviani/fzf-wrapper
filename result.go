package fzfwrapper

import (
	"math"
	"sort"

	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

type MyResult struct {
	item   *MyItem
	points [4]uint16
}

func MbuildResult(item *MyItem, offsets []fzf.Offset, score int) MyResult {
	if len(offsets) > 1 {
		sort.Sort(fzf.ByOrder(offsets))
	}

	result := MyResult{item: item}

	// only sort by score implemented
	val := uint16(math.MaxUint16)
	val = math.MaxUint16 - util.AsUint16(score)

	result.points[3] = val

	return result
}
