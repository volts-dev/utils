package utils

import (
	"fmt"
	//"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	one := NewStringList()
	one.PushString("a")
	one.PushString("in")
	one.AddSubList([]string{"1", "2", "3"}...)

	fmt.Println("result1", one.String())

	one = NewStringList()
	one.PushString("a")
	one.PushString("=")
	one.AddSubList("true")

	fmt.Println("result1", one.String())
}

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

	t.Log("add_sub_item successfull", lLst.String())
}

func Test_has(t *testing.T) {
	Lst := NewStringList("sdfa")
	if Lst.Len() == 1 {
		t.Log("Create ok", Lst.String())
	}
	Item0 := NewStringList(`"res_partner" as "res_user__partner_id"`)
	Item1 := NewStringList("Item1")
	Lst.Push(Item0)
	Lst.Push(Item1)

	if !Lst.Has(`"res_partner" as "res_user__partner_id"`) {
		t.Log("Has failure string", Lst.String())
	}

	if !Lst.Has(Item1) {
		t.Log("Has failure object", Lst.String())
	}

	t.Log("Has successfull", Lst.String())
}
