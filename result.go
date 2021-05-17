package fzfwrapper

import (
	"math"
	"sort"

	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

type Result struct {
	item   *Item
	points [4]uint16
}

func buildResult(item *Item, offsets []fzf.Offset, score int) Result {
	if len(offsets) > 1 {
		sort.Sort(fzf.ByOrder(offsets))
	}

	result := Result{item: item}

	// only sort by score implemented
	val := uint16(math.MaxUint16)
	val = math.MaxUint16 - util.AsUint16(score)

	result.points[3] = val

	return result
}
