package testing

import (
	"container/list"
	"github.com/TrueFMartin/engotut/systems"
	"log"
	"testing"
)

func TestCopyListUntil(t *testing.T) {
	startList := list.New()
	startList.PushFront(0)
	startList.PushFront(1)
	stopElem := startList.PushFront(2)
	startList.PushFront(3)
	copyList := systems.CopyListUntil(startList, stopElem)

	copyList.PushFront(9)

	got := make([]int, 0)
	for e := copyList.Back(); e != nil; e = e.Prev() {
		val, ok := e.Value.(int)
		if ok != true {
			log.Fatal("List assertion not int")
		}
		got = append(got, val)
	}
	want := []int{0, 1, 2, 9}
	for i, _ := range got {
		if got[i] != want[i] {
			t.Error("Got ", got, " Want ", want)
		}
	}
}
