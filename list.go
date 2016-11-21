package utils

/*

type(
	TList struct {
	Items []interface{} // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
)

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
func (l *TList) MoveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.insert(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
func (l *TList) MoveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.insert(l.remove(e), l.root.prev)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *TList) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
func (l *TList) Remove(e *Element) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *TList) Len() int { return l.len }
*/
