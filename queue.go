package utils

import "sync"

type (
	Queue struct {
		sync.RWMutex
		start  *node
		end    *node
		pos    *node
		length int
		cap    int
	}

	node struct {
		value  interface{}
		next   *node
		isLast bool
	}
)

// Create a new queue
func NewQueue() *Queue {
	return &Queue{}
}

// Take the next item off the front of the queue
func (self *Queue) Dequeue() interface{} {
	self.Lock()
	defer self.Unlock()

	if self.length == 0 {
		return nil
	}
	n := self.start
	if self.length == 1 {
		self.start = nil
		self.end = nil
		self.pos = nil
	} else {
		if self.pos == self.start {
			self.pos = self.start.next
		}
		self.start = self.start.next
		self.end.next = self.start
	}

	self.length--
	return n.value
}

// Put an item on the end of a queue
func (self *Queue) Enqueue(value interface{}) {
	self.Lock()

	if self.length == 0 {
		n := &node{value, nil, true}
		self.start = n
		self.end = n
		self.pos = n
	} else {
		n := &node{value, self.start, true}
		self.end.isLast = false
		self.end.next = n
		self.end = n
	}

	self.length++
	self.Unlock()
}

// Return the number of items in the queue
func (self *Queue) Len() int {
	self.Lock()
	defer self.Unlock()
	return self.length
}

func (self *Queue) IsEmpty() bool {
	return self.length == 0
}

func (self *Queue) IsFull() bool {
	self.Lock()
	defer self.Unlock()
	return self.length >= self.cap
}

func (self *Queue) IsLast() bool {
	return self.pos.isLast
}

// move peek position to first node
func (self *Queue) First() {
	self.pos = self.start
}

// move peek position to last node
func (self *Queue) Last() {
	self.pos = self.end
}

// Return the first item in the queue without removing it
func (self *Queue) Peek() (value interface{}, islast bool) {
	self.Lock()
	n := self.pos
	self.pos = self.pos.next
	self.Unlock()
	return n.value, n.isLast
}
