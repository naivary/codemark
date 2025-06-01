package converter

import (
	"errors"
)

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Pop() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}
	el := q.items[0]
	q.items = q.items[1:]
	return el, nil
}

type Node struct {
	value    string
	children []*Node
}

func bfs(root *Node, q Queue[string]) string {
	if len(root.children) == 0 {
		return _undefined
	}
	segment, err := q.Pop()
	if err != nil {
		return root.children[0].value
	}
	for _, child := range root.children {
		if child.value == segment {
			return bfs(child, q)
		}
	}
	return _undefined
}
