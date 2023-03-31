package utils

import (
	"fmt"
	"testing"
)

func TestJoinQuote(t *testing.T) {
	t.Log(JoinQuote([]string{"aa", "bb", "cc"}, "'", ","))
}
func TestReversed(t *testing.T) {
	t.Log()
	fmt.Println(Reversed([]any{"aa", "bb", "cc"}...))
}
