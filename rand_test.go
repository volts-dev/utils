package utils

import (
	//"fmt"
	"testing"
)

func TestRandomCreateBytes(t *testing.T) {
	t.Log(string(RandomCreateBytes(16)))
}
