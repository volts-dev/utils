package utils

import (
	"fmt"
	"testing"
)

type (
	SubTest struct {
		C string
		D int
		//	c []string
	}
	SubTest2 struct {
		E string `field:"aa"`
		F int
		//	c []string
	}
	Test struct {
		SubTest `field:"-"`
		SubTest2
		A string `field:"aa"`
		b int    `field:"aa"`
		//	c []string
	}
)

func TestConvert(t *testing.T) {
	fmt.Println("ffff")
	ts := &Test{A: "fff", b: 44 /*c: []string{"gdf", "dfg"}*/}
	ts.C = "asd"
	ts.D = 12
	ts.E = "E"
	ts.F = 13
	n := NewStructMapper(ts)
	//fmt.Println("ffff", ts, n.Map())

	n.IsFlat = true
	fmt.Println("ffff", ts, n.Map())
}
