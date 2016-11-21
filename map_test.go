package utils

import (
	"fmt"
	"testing"
)

type (
	Test struct {
		A string
		b int
		//	c []string
	}
)

func TestConvert(t *testing.T) {
	fmt.Println("ffff")
	ts := &Test{A: "fff", b: 44 /*c: []string{"gdf", "dfg"}*/}
	fmt.Println("ffff", ts, Map(ts))
}
