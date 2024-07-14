package utils

import (
	"cmp"
)

func Min[T cmp.Ordered](args ...T) T {
	min := args[0]
	for _, x := range args {
		if min > x {
			min = x
		}
	}
	return min

}
func Max[T cmp.Ordered](args ...T) T {
	max := args[0]
	for _, x := range args {
		if max < x {
			max = x
		}
	}
	return max
}
