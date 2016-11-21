package utils

import (
	"fmt"
	//"fmt"
	"testing"
)

func Test_add_sub_item(t *testing.T) {
	lLst := NewStringList("sdfa")
	if lLst.Len() == 1 {
		t.Log("Create ok", lLst.String())
	}

	lLst.PushString("asdf")
	if lLst.Len() == 2 {
		t.Log("PushString ok", lLst.String())
	}

	lPop := lLst.Pop()
	if lLst.Len() == 1 {
		t.Log("Pop ok", lPop, lLst.String())
	}

	lShift := lLst.Shift()
	if lLst.Len() == 0 {
		t.Log("Shift ok", lShift, lLst.String())
	}

	lLst.PushString(fmt.Sprintf("%s and %s", lPop.String(), lShift.String()))
	t.Log("Shift ok", lLst.String())

	lLst2 := NewStringList()
	lLst2.Push(lLst)
	if lLst2.Len() == 1 {
		t.Log("Push ok", lLst2.String())
	}

}
