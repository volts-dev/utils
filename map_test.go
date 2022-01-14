package utils

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap()
	cnt := 100
	var wg sync.WaitGroup
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			m.Set(fmt.Sprintf("map%d", n), n)
			wg.Done()
		}(i)
	}
	wg.Wait()
	itms := m.Items()
	t.Log(itms, len(itms))
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func(n int) {
			m.Delete(fmt.Sprintf("map%d", n))
			wg.Done()
		}(i)
	}
	wg.Wait()
	itms = m.Items()
	t.Log(itms, len(itms))

	t.Log(m.Count())
}
