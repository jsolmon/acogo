package main

import (
	"fmt"
	"math"
	"sync"
)

type NodeType int

const (
	Home NodeType = iota
	Goal
	Path
)

type Graph struct {
	Nodes       []*Node
	HomeIdx     int
	GoalIdx     int
	DecayFactor float64
}

func NewGraph(dimension, homeNode, goalNode int, decayFactor float64) *Graph {
	// generate edges and put them in a [][]*Edge 2d slice
	edges := generateEdges(dimension)

	// iterate through the [][]*Edge 2D slice to generate nodes
	nodes := generateNodes(edges, dimension, homeNode, goalNode)

	// return graph from list of nodes
	return &Graph{
		Nodes:       nodes,
		HomeIdx:     homeNode,
		GoalIdx:     goalNode,
		DecayFactor: decayFactor,
	}
}

func (g *Graph) Run() {
	for _, n := range g.Nodes {
		go func(n Node) { n.Run() }(*n)
	}
}

func generateEdges(dim int) [][]*Edge {
	edges := make([][]*Edge, dim*dim)
	for i := 0; i < dim*dim; i++ {
		edges[i] = make([]*Edge, dim*dim)
	}

	for n := 0; n < dim*dim; n++ {
		// edge to node above
		if n >= dim {
			edges[n][n-dim] = NewEdge(n, n-dim)
		}
		// edge to node above/right
		if n >= dim && n%dim != dim-1 {
			edges[n][n-dim+1] = NewEdge(n, n-dim+1)
		}
		// edge to node right
		if n%dim != dim-1 {
			edges[n][n+1] = NewEdge(n, n+1)
		}
		// edge to node below/right
		if dim*dim-n > dim && n%dim != dim-1 {
			edges[n][n+dim+1] = NewEdge(n, n+dim+1)
		}
		// edge to node below
		if dim*dim-n > dim {
			edges[n][n+dim] = NewEdge(n, n+dim)
		}
		// edge to node below/left
		if dim*dim-n > dim && n%dim != 0 {
			edges[n][n+dim-1] = NewEdge(n, n+dim-1)
		}
		// edge to node left
		if n%dim != 0 {
			edges[n][n-1] = NewEdge(n, n-1)
		}
		// edge to node above/left
		if n%dim != 0 && n >= dim {
			edges[n][n-dim-1] = NewEdge(n, n-dim-1)
		}
	}
	return edges
}

func generateNodes(edges [][]*Edge, dim, homeNode, goalNode int) []*Node {
	// in edges, each row is outgoing edges, each column is incoming edges for a given node
	nodes := make([]*Node, dim*dim)
	for n, row := range edges {
		outEdges := make([]*Edge, 0, 8)
		inEdges := make([]*Edge, 0, 8)

		// get all active edges in row for OutEdges in node
		for _, out := range row {
			if out != nil {
				outEdges = append(outEdges, out)
			}
		}
		// pull active edges in column for InEdges in node
		for r := 0; r < dim*dim; r++ {
			if edges[r][n] != nil {
				inEdges = append(inEdges, edges[r][n])
			}
		}
		// if n == homeNode or goalNode, set NodeType as home or goal, otherwise path
		nodeType := Path
		if n == homeNode {
			nodeType = Home
		}
		if n == goalNode {
			nodeType = Goal
		}
		// create node and add to nodes
		nodes[n] = NewNode(n, inEdges, outEdges, nodeType)
	}

	return nodes
}

func (g *Graph) Dissipate() {
	for _, n := range g.Nodes {
		for _, e := range n.InEdges {
			e.AddPheremone(g.DecayFactor)
		}
	}
}

func (g *Graph) MarkPath(steps []int, depositAmt float64) {
	lastStep := steps[0]
	for _, nodeId := range steps[0:] {
		node := g.Nodes[nodeId]
		node.MarkEdge(lastStep, depositAmt)
		lastStep = nodeId
	}
}

type Node struct {
	Id       int
	InEdges  []*Edge
	OutEdges []*Edge
	Type     NodeType
}

func NewNode(id int, inEdges []*Edge, outEdges []*Edge, t NodeType) *Node {
	return &Node{
		Id:       id,
		InEdges:  inEdges,
		OutEdges: outEdges,
		Type:     t,
	}
}

func (n *Node) Run() {
	for _, e := range n.InEdges {
		edge := e
		go n.runAnts(edge)
	}
}

func (n *Node) runAnts(e *Edge) {
	for {
		ant := <-e.Path
		next, atGoal := ant.ChoosePath(n)
		if atGoal { //ant has reached goal - no more to do
			continue
		}
		next.Path <- ant
	}
}

func (n *Node) MarkEdge(from int, depositAmt float64) {
	for _, e := range n.InEdges {
		if e.StartNodeId == from {
			e.AddPheremone(depositAmt)
		}
	}
}

type Edge struct {
	Path        chan Ant
	StartNodeId int
	EndNodeId   int

	pheremone float64
	mu        sync.RWMutex
}

func NewEdge(startId, endId int) *Edge {
	return &Edge{
		Path:        make(chan Ant, 5),
		StartNodeId: startId,
		EndNodeId:   endId,
		pheremone:   10.0,
	}
}

func (e *Edge) Pheremone() float64 {
	return e.pheremone
}

func (e *Edge) AddPheremone(f float64) {
	e.pheremone = math.Max(e.pheremone+f, 0.1)
}

func (e Edge) String() string {
	return fmt.Sprintf("%d -> %d: %.2f", e.StartNodeId, e.EndNodeId, e.pheremone)
}
