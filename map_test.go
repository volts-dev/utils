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
	Test struct {
		SubTest
		A string
		b int
		//	c []string
	}
)

func TestConvert(t *testing.T) {
	fmt.Println("ffff")
	ts := &Test{A: "fff", b: 44 /*c: []string{"gdf", "dfg"}*/}
	ts.C = "asd"
	ts.D = 12
	fmt.Println("ffff", ts, Map(ts))
}
