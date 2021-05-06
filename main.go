package fzfwrapper

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/junegunn/fzf/src/util"
)

var (
	s          []string
	partitions int = util.Min(8*runtime.NumCPU(), 32)
)

const (
	fzfChunkSize int = 100 // https://github.com/junegunn/fzf/blob/764316a53d0eb60b315f0bbcd513de58ed57a876/src/constants.go#L38
)

type MyToken struct {
	text         *util.Chars
	prefixLength int32
}

func main() {
	// tmp
	s = make([]string, 0)
	f, _ := os.Open("/tmp/c")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	f.Close()

	/////////////////////////////////////////////////////////////////////////////

	// p.termSets eh basicamente uma lista dos termos digitados: ["apt", "install", "-t"]
	// term.text eh o termo em si: "apt"

	mypattern := "apt install -t"

	var itemIndex int32
	chunkList := MyNewChunkList(func(item *MyItem, data []byte) bool {

		item.text = util.ToChars(data)
		item.text.TrimTrailingWhitespaces()
		item.text.Index = itemIndex

		itemIndex++
		return true

	})

	for _, str := range s {
		chunkList.Push([]byte(str))
	}

	snapshot, _ := chunkList.Snapshot()

	patternSlice := strings.Split(mypattern, " ")
	patternRunes := make([][]rune, 0, len(patternSlice))
	for _, ps := range patternSlice {
		patternRunes = append(patternRunes, []rune(ps))
	}
	merger := MScan(snapshot, patternRunes)

	for _, list := range merger {
		for _, v := range list {
			fmt.Println(v.item.AsString(false))
		}
	}
	return

	/////////////////////////////////////////////////////////////////////////////
}
