package main

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// GraphTest ensures that NewGraph generates a 3x3 grid with each
// node connected to adjacent nodes above, below, left, right, and
// on all diagonals.
func TestGraph(t *testing.T) {
	g := NewGraph(3, 2, 6, 0.5)

	if len(g.Nodes) != 9 {
		t.Error(fmt.Sprintf("expected 9 nodes but found %v\n", len(g.Nodes)))
	}
	validateNode(g.Nodes[0], []int{1, 3, 4}, Path, t)
	validateNode(g.Nodes[1], []int{0, 2, 3, 4, 5}, Path, t)
	validateNode(g.Nodes[2], []int{1, 4, 5}, Home, t)
	validateNode(g.Nodes[3], []int{0, 1, 4, 6, 7}, Path, t)
	validateNode(g.Nodes[4], []int{0, 1, 2, 3, 5, 6, 7, 8}, Path, t)
	validateNode(g.Nodes[5], []int{1, 2, 4, 7, 8}, Path, t)
	validateNode(g.Nodes[6], []int{3, 4, 7}, Goal, t)
	validateNode(g.Nodes[7], []int{3, 4, 5, 6, 8}, Path, t)
	validateNode(g.Nodes[8], []int{4, 5, 7}, Path, t)
}

func validateNode(n *Node, edgesTo []int, nodeType NodeType, t *testing.T) {
	if n.Type != nodeType {
		t.Error(fmt.Sprintf("node %v should be type %v but was %v\n", n.Id, nodeType, n.Type))
	}
	if len(n.InEdges) != len(edgesTo) {
		t.Error(fmt.Sprintf("expected node %v to have %v edges but had %v", n.Id, len(edgesTo), len(n.InEdges)))
	}
	inEdgeList := make([]int, len(n.InEdges))
	for idx, e := range n.InEdges {
		inEdgeList[idx] = e.StartNodeId
	}
	outEdgeList := make([]int, len(n.InEdges))
	for idx, e := range n.OutEdges {
		outEdgeList[idx] = e.EndNodeId
	}

	sort.Sort(sort.IntSlice(inEdgeList))
	sort.Sort(sort.IntSlice(outEdgeList))
	sort.Sort(sort.IntSlice(edgesTo))

	if !reflect.DeepEqual(inEdgeList, edgesTo) {
		t.Error(fmt.Sprintf("expected node %v to have InEdges %v but got %v", n.Id, inEdgeList, edgesTo))
	}
	if !reflect.DeepEqual(outEdgeList, edgesTo) {
		t.Error(fmt.Sprintf("expected node %v to have OutEdges %v but got %v", n.Id, inEdgeList, edgesTo))
	}
}
