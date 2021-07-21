package utils

import (
	"testing"
)

func TestJoinQuote(t *testing.T) {
	t.Log(JoinQuote([]string{"aa", "bb", "cc"}, "'", ","))
}
