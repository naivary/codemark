package codemark

import (
	"errors"
	"strings"
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

type node struct {
	value    string
	children []*node
	conv     Converter
}

type Tree struct {
	root *node
}

func (t *Tree) add(n *node, q Queue[string]) {
	seg, err := q.Pop()
	if err != nil {
		n.children = append(n.children, n)
		return
	}
	for _, child := range n.children {
		if child.value == seg {
			t.add(child, q)
		}
	}

	nn := &node{value: seg, children: []*node{}}
	n.children = append(n.children, nn)
	t.add(nn, q)
}

func (t *Tree) Add(n *node) {
	typeID := n.value
	q := Queue[string]{strings.Split(typeID, ".")}
	t.add(n, q)
}

func bfs(root *node, q Queue[string]) *node {
	if len(root.children) == 0 {
		return nil
	}
	segment, err := q.Pop()
	if err != nil {
		return root
	}
	for _, child := range root.children {
		if child.value == segment {
			return bfs(child, q)
		}
	}
	return nil
}
