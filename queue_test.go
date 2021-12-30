package utils

import (
	"log"
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup
	read := func(q *Queue, id int) {
		for {
			if q.Len() == 0 {
				break
			}

			item := q.Dequeue()
			if item == nil {
				// log.Println(err)
				// time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Println("id", id, "------", item.(int))
		}
		wg.Done()
	}

	q := NewQueue()
	for i := 0; i < 1000; i++ {
		//time.Sleep(50 * time.Millisecond)
		q.Enqueue(i)
	}

	wg.Add(3)
	go read(q, 1)
	go read(q, 2)
	go read(q, 3)

	wg.Wait()
}

func TestPeak(t *testing.T) {
	var wg sync.WaitGroup
	fn := func(q *Queue, id int) {
		for i := 0; i < 1; i++ {
			for {
				item, islast := q.Peek()

				if item == nil {
					// log.Println(err)
					// time.Sleep(100 * time.Millisecond)
					continue
				}
				t.Log("id", id, "------", item.(int))
				if islast {
					t.Log("last id", id, "------", item.(int))
					break
				}
			}
		}
		wg.Done()
	}

	q := NewQueue()
	for i := 0; i < 1000; i++ {
		q.Enqueue(i)
	}

	wg.Add(1)
	go fn(q, 1)
	wg.Add(1)
	go fn(q, 2)

	t.Log("count", q.Len())

	wg.Wait()
}
