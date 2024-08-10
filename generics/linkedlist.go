package generics

import "fmt"

type node[T any] struct {
	val  T
	next *node[T]
}

type LinkedList[T any] struct {
	start *node[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
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

func (l *LinkedList[T]) String() string {
	out := "["
	curr := l.start
	for {
		if curr == nil {
			break
		}
		out += fmt.Sprintf("%v", curr.val)
		curr = curr.next
		if curr != nil {
			out += " "
		} else {
			out += "]"
		}
	}
	return out
}
