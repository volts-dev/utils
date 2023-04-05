package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestSliceEqual(t *testing.T) {
	log.Println(SliceEqual([]string{"aa", "bb", "cc"}, []string{"", ""}))
	log.Println(SliceEqual([]any{"aa", "bb", "cc"}, []any{"aa", "bb", "cc"}))

}

func TestJoinQuote(t *testing.T) {
	t.Log(JoinQuote([]string{"aa", "bb", "cc"}, "'", ","))
}

func TestReversed(t *testing.T) {
	t.Log()
	fmt.Println(Reversed([]any{"aa", "bb", "cc"}...))
}
