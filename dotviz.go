package main

import (
	"fmt"
	"strconv"

	"code.google.com/p/gographviz"
)

// ToDot takes in an acogo Graph and transforms it into a DOT language
// graph with the darkness of the edge color corresponding to the amount
// of pheromone on the edge.
func ToDot(g *Graph, max float64) *gographviz.Graph {
	gv := gographviz.NewGraph()
	gv.SetDir(true)

	for _, n := range g.Nodes {
		gv.AddNode(gv.Name, strconv.Itoa(n.Id), nodeAttrs(n))
		for _, e := range n.InEdges {
			gv.AddEdge(strconv.Itoa(e.StartNodeId), strconv.Itoa(e.EndNodeId), true, edgeAttrs(e, max))
		}
	}

	return gv
}

// nodeAttrs assigns DOT attributes to a node, assigning labels and
// colors based on whether they are home or goal nodes.
func nodeAttrs(n *Node) map[string]string {
	attrs := make(map[string]string, 2)

	switch n.Type {
	case Home:
		attrs["color"] = "\"#8B0000\"" // maroon
		attrs["label"] = fmt.Sprintf("\"%v: HOME\"", strconv.Itoa(n.Id))
	case Goal:
		attrs["color"] = "\"#008000\"" // green
		attrs["label"] = fmt.Sprintf("\"%v: GOAL\"", strconv.Itoa(n.Id))
	default: // path nodes
		attrs["color"] = "\"#D3D3D3\"" // light grey
	}
	return attrs
}

// edgeAttrs assigns DOT attributes to a node, assigning darker colors
// to nodes with more pheromone and lighter colors to nodes with less.
func edgeAttrs(e *Edge, max float64) map[string]string {
	attrs := make(map[string]string, 3)
	attrs["penwidth"] = "3.0"
	attrs["arrowType"] = "open"

	// sets a minimum alpha value of 10 so edges with no traffic will still appear in the graph
	alpha := int(e.pheromone/max*245) + 10
	color := fmt.Sprintf("\"#104E8B%X\"", alpha)
	attrs["color"] = color

	return attrs
}
