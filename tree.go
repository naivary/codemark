package codemark

import (
	"errors"
	"fmt"
	"strings"
)

const (
	_rootNodeValue = "ROOT_NODE"
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

type tree struct {
	root *node
}

func (t *tree) add(current *node, convNode *node, q Queue[string]) {
	seg, err := q.Pop()
	if err != nil {
		current.conv = convNode.conv
		return
	}
	for _, childNode := range current.children {
		if childNode.value == seg {
			t.add(childNode, convNode, q)
		}
	}
	nn := &node{value: seg, children: []*node{}}
	current.children = append(current.children, nn)
	t.add(nn, convNode, q)
}

func (t *tree) Add(convNode *node) error {
	typeID := convNode.value
	if convNode.conv == nil {
		return errors.New("conveter cannot be nil")
	}
	if typeID == "" {
		return errors.New("empty type id is not valid")
	}
	q := Queue[string]{strings.Split(typeID, ".")}
	kind, err := q.Pop()
	if err != nil {
		return err
	}
	for _, childNode := range t.root.children {
		if childNode.value == kind {
			t.add(childNode, convNode, q)
			return nil
		}
	}
	nn := &node{value: kind, children: []*node{}}
	t.root.children = append(t.root.children, nn)
	t.add(nn, convNode, q)
	return nil
}

func (t *tree) getConverter(n *node, q Queue[string]) Converter {
	kind, err := q.Pop()
	if err != nil {
		return n.conv
	}
	for _, child := range n.children {
		if child.value == kind {
			return t.getConverter(child, q)
		}
	}
	return nil
}

func (t *tree) GetConverter(typeID string) (Converter, error) {
	if typeID == "" {
		return nil, errors.New("typeid cannot be empty")
	}
	q := Queue[string]{strings.Split(typeID, ".")}
	conv := t.getConverter(t.root, q)
	if conv == nil {
		return nil, fmt.Errorf("no converter found for type id `%s`", typeID)
	}
	return conv, nil
}
