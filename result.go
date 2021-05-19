package fzfwrapper

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/junegunn/fzf/src"
)

// Sort criteria
type criterion int

const (
	ByScore    criterion = iota // sort by score
	ByLength                    // sort by string length
	ByPosition                  // sort by match location

	maxInt = int(^uint(0) >> 1)
)

type Result struct {
	Item    *Item
	Points  [4]int
	Offsets []fzf.Offset
	Pos     *[]int
}

func buildResult(item *Item, offsets []fzf.Offset, score int, sortBy []criterion) Result {
	if len(offsets) > 1 {
		sort.Sort(fzf.ByOrder(offsets))
	}

	result := Result{Item: item}

	for idx, criterion := range sortBy {
		val := maxInt

		switch criterion {
		case ByScore:
			val = maxInt - score
		case ByPosition:
			val = int(offsets[0][0]) // position of first match
		case ByLength:
			val = int(item.TrimLength())
		}

		result.Points[3-idx] = val
	}

	return result
}

func (r *Result) HighlightResult() string {
	text := r.Item.AsString(false)

	if r.Pos == nil {
		return text
	}

	// sort offsets
	var pos []int
	pos = *r.Pos

	if len(pos) == 0 {
		return text
	}

	charOffsets := make([]fzf.Offset, len(pos))
	for idx, p := range pos {
		offset := fzf.Offset{int32(p), int32(p + 1)}
		charOffsets[idx] = offset
	}
	sort.Sort(ByOrder(charOffsets))

	offsets := MergeOffsets(charOffsets)

	m := color.New(color.FgRed)
	m.EnableColor()
	matchColor := m.SprintFunc()
	//matchColor := color.New(108).SprintFunc()

	begin := offsets[0][0]
	end := offsets[0][1]
	idx := end

	dest := bytes.NewBufferString("")

	if begin == 0 {
		// match is already at first char
		// start color and copy the first matched substring
		fmt.Fprintf(dest, "%s", matchColor(text[begin:end]))
	} else {
		// match is not first char
		// copy the substring until the first match,
		// begin color and copy the first matched substring
		fmt.Fprintf(dest, "%s%s", text[:begin], matchColor(text[begin:end]))
	}

	l := len(offsets)

	for off := 1; off < l; off++ {
		begin = offsets[off][0]

		if idx > begin {
			panic(fmt.Sprintf("idx: %d, off: %d, offsets: %v\npos: %v\n", idx, off, offsets, *r.Pos))
		}
		if idx != begin { // print not matched substring
			fmt.Fprintf(dest, "%s", text[idx:begin])
		}

		end = offsets[off][1] // somehow idx was referencing end instead of copying its value

		// end color
		fmt.Fprintf(dest, "%s", matchColor(text[begin:end]))

		idx = end
	}

	fmt.Fprintf(dest, "%s", text[idx:])
	return dest.String()
}

func MergeOffsets(matchOffsets []fzf.Offset) []fzf.Offset {
	offsets := []fzf.Offset{}

	begin := matchOffsets[0][0]
	end := matchOffsets[0][0]
	for _, off := range matchOffsets {
		if begin <= off[0] && end == off[1] {
			// duplicated offset
			continue
		}

		if end == off[0] {
			end = off[1]
			continue
		}
		offsets = append(offsets, fzf.Offset{begin, end})
		begin = off[0]
		end = off[1]
	}

	offsets = append(offsets, fzf.Offset{begin, end})

	return offsets
}

func beginColor() string {
	return fmt.Sprintf("\033[1;108m")
}
func endColor() string {
	return fmt.Sprintf("\033[0m")
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
		left := irank.Points[idx]
		right := jrank.Points[idx]
		if left < right {
			return true
		} else if left > right {
			return false
		}
	}

	return (irank.Item.Index() <= jrank.Item.Index()) != tac
}
