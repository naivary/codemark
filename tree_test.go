package codemark

import (
	"reflect"
	"strings"
	"testing"
)

func TestQueue(t *testing.T) {
	example := "slice.ptr.string"
	slice := strings.Split(example, TypeIDSep)
	q := Queue[string]{slice}
	if len(q.items) != 3 {
		n := len(strings.Split(example, TypeIDSep))
		t.Fatalf("length of queue is wrong. expected: %d\n", n)
	}
	item, err := q.Pop()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	if item != "slice" {
		t.Fatalf("wrong first item. expected: %s\n", slice[0])
	}
	if len(q.items) != 2 {
		t.Fatalf("wrong queue size. expected: %d\n", len(slice)-1)
	}
	if slice[1] != "ptr" || slice[2] != "string" {
		t.Fatalf("queue does not contain the correct elements after resize. expected: %s and %s\n", slice[1], slice[2])
	}
}

func newConvNode(typ any, conv Converter) *node {
	id := TypeID(reflect.TypeOf(typ))
	return &node{
		value: id,
		conv:  conv,
	}
}

func exampleTree() *tree {
	i := &node{value: "int"}
	str := &node{value: "string", conv: &stringConverter{}}
	ptr := &node{value: "ptr", children: []*node{str}}
	t := newTree(i, ptr)
	return t
}

func TestTree_Add(t *testing.T) {
	tests := []struct {
		name     string
		convNode *node
		t        *tree
		query    string
	}{
		{
			name:     "add to empty tree",
			convNode: newConvNode(new(string), &stringConverter{}),
			t:        newTree(),
			query:    "ptr.string",
		},
		{
			name:     "add to existing node",
			convNode: newConvNode(new(string), &stringConverter{}),
			t:        newTree(&node{value: "ptr"}),
			query:    "ptr.string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.t.Add(tc.convNode); err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			conv, err := tc.t.GetConverter(tc.query)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			if conv == nil {
				t.Fatalf("couldn't get converter after adding it. query: %s\n", tc.query)
			}
		})
	}
}

func TestTree_GetConverter(t *testing.T) {
	tests := []struct {
		name       string
		typeID     string
		isExisting bool
	}{
		{
			name:       "string node",
			typeID:     "ptr.string",
			isExisting: true,
		},
		{
			name:   "non existing node",
			typeID: "slice.ptr.int",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tr := exampleTree()
			conv, err := tr.GetConverter(tc.typeID)
			if err != nil && tc.isExisting {
				t.Errorf("err occured: %s\n", err)
			}
			if err == nil && !tc.isExisting {
				t.Fatalf("expected no result; got: %v\n", conv)
			}
		})
	}
}
