package utils

import (
	"fmt"
	"testing"
)

func TestRandomCreateBytes(t *testing.T) {
	t.Log(string(RandomCreateBytes(16)))
}

func TestTitleCasedName(t *testing.T) {
	fmt.Println(TitleCasedName("hello_world"))

	fmt.Println(TitleCasedName("TmT"))
}
