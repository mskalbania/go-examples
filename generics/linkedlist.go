package generics

import "fmt"

type node[T comparable] struct {
	val  T
	next *node[T]
}

type LinkedList[T comparable] struct {
	start *node[T]
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) Add(e T) {
	newTail := &node[T]{e, nil}
	n := l.start
	if n == nil { //was empty
		l.start = newTail
		return
	}
	for {
		if n.next == nil { //travers to find tail and append
			n.next = newTail
			return
		}
		n = n.next
	}
}

// Delete deletes first occurence of element, returns true if found and deleted, false otherwise
func (l *LinkedList[T]) Delete(v T) bool {
	curr := l.start
	if curr == nil {
		return false
	}
	if curr.val == v {
		l.start = nil
		return true
	}
	for {
		//looking always at next, go doesn't support pointer arithmetic
		//can only set variable ref to nil or reattach to other ref, can't clean up memory at given address
		if curr.next == nil { //no more elements in list
			return false
		}
		if curr.next.val == v {
			if curr.next.next != nil {
				curr.next = curr.next.next //switch links
			} else {
				curr.next = nil //next is last in chain so just remove reference
			}
			return true
		}
		curr = curr.next
	}
}

// Search returns index of element, -1 if not found
func (l *LinkedList[T]) Search(v T) int {
	curr := l.start
	if curr == nil {
		return -1
	}
	var idx int
	for {
		if curr.val == v {
			return idx
		}
		if curr.next != nil {
			curr = curr.next
			idx++
		} else {
			return -1
		}
	}
}

func (l *LinkedList[T]) String() string {
	out := "["
	curr := l.start
	for {
		if curr == nil {
			out += "]"
			break
		}
		out += fmt.Sprintf("%v", curr.val)
		curr = curr.next
		if curr != nil {
			out += ", "
		}
	}
	return out
}
