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
		E string `field:"ee"`
		F int
		//	c []string
	}
	Test struct {
		SubTest `field:"-"` // 忽略
		SubTest2
		MyName string
		Age    int `field:"myage"`
		//	c []string
	}
)

func TestConvert(t *testing.T) {
	ts := &Test{MyName: "fff", Age: 44 /*c: []string{"gdf", "dfg"}*/}
	ts.C = "asd"
	ts.D = 12
	ts.E = "E"
	ts.F = 13
	n := NewStructMapper(ts)
	//fmt.Println("ffff", ts, n.Map())

	n.IsFlat = true
	fmt.Println(ts)
	fmt.Println(n.Map())

}
