package utils

import (
	"log"
	"testing"
)

func TestMax(t *testing.T) {
	m := Max(4, 5, 8, 0.0, 6.3, 8, 2, 7.1, 1, 9.4)
	log.Println(m)

	s := Max("4", "5", "18", "0", "6", "8", "2", "7", "1")
	log.Println(s)
}
