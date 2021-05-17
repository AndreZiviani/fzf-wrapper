package fzfwrapper

import (
	"math"
	"sort"

	"github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

// Sort criteria
type criterion int

const (
	ByScore  criterion = iota // dont sort
	ByLength                  // sort by string length
	ByBegin                   // not implemented
	ByEnd                     // not implemented
)

type Result struct {
	item   *Item
	points [4]uint16
}

func buildResult(item *Item, offsets []fzf.Offset, score int) Result {
	sortCriteria := []criterion{ByScore, ByLength}
	if len(offsets) > 1 {
		sort.Sort(fzf.ByOrder(offsets))
	}

	result := Result{item: item}

	for idx, criterion := range sortCriteria {
		val := uint16(math.MaxUint16)

		switch criterion {
		case ByScore:
			val = math.MaxUint16 - util.AsUint16(score)
		case ByLength:
			val = item.TrimLength()

		}

		result.points[3-idx] = val
	}

	return result
}

// ByOrder is for sorting substring offsets
type ByOrder []fzf.Offset

func (a ByOrder) Len() int {
	return len(a)
}

func (a ByOrder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByOrder) Less(i, j int) bool {
	ioff := a[i]
	joff := a[j]
	return (ioff[0] < joff[0]) || (ioff[0] == joff[0]) && (ioff[1] <= joff[1])
}

// ByRelevance is for sorting Items
type ByRelevance []Result

func (a ByRelevance) Len() int {
	return len(a)
}

func (a ByRelevance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRelevance) Less(i, j int) bool {
	return compareRanks(a[i], a[j], false)
}

// ByRelevanceTac is for sorting Items
type ByRelevanceTac []Result

func (a ByRelevanceTac) Len() int {
	return len(a)
}

func (a ByRelevanceTac) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRelevanceTac) Less(i, j int) bool {
	return compareRanks(a[i], a[j], true)
}

func compareRanks(irank Result, jrank Result, tac bool) bool {
	for idx := 3; idx >= 0; idx-- {
		left := irank.points[idx]
		right := jrank.points[idx]
		if left < right {
			return true
		} else if left > right {
			return false
		}
	}

	return (irank.item.Index() <= jrank.item.Index()) != tac
}
