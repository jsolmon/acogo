package main

import (
	"fmt"
	"math"
)

type NodeType int

const (
	Home NodeType = iota
	Goal
	Path
)

// Graph struct holds nodes and edges that ants travel on
type Graph struct {
	// List of nodes containing edges
	Nodes []*Node
	// Index of the Home/start node
	HomeIdx int
	// Index of the Goal node
	GoalIdx int
	// How much the pheremone on each edge decreases after each round of ants
	// reaches the goal.
	DecayFactor float64
}

// NewGraph generates a new graph. The default graph at this time is a square of
// dim * dim nodes with each node having an edge to adjacent nodes above, below,
// left, right, and at all four diagonals.
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

// Run calls Run on each node which calls Run on each edge initializing go
// routines which pass ants from edge to edge in the graph.
func (g *Graph) Run() {
	for _, n := range g.Nodes {
		go func(n Node) { n.Run() }(*n)
	}
}

// generateEdges generates a slice of edges for a dim*dim graph such that each
// edge connects to adjacent nodes above, below, left, right, and on all four
// diagonals.
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

// generateNodes generates dim*dim nodes based and gives each the in/out edges
// mapped in the 2D slice of edges.
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

// Dissipate subtracts g.DecayFactor pheremone from each edge in the graph.
func (g *Graph) Dissipate() {
	for _, n := range g.Nodes {
		for _, e := range n.InEdges {
			e.AddPheremone(-1 * g.DecayFactor)
		}
	}
}

// MarkPath takes in a list of nodeIds representing the path an ant followed
// and adds depositAmt pheremone to each edge along the path.
func (g *Graph) MarkPath(steps []int, depositAmt float64) {
	lastStep := steps[0]
	for _, nodeId := range steps[0:] {
		node := g.Nodes[nodeId]
		node.MarkEdge(lastStep, depositAmt)
		lastStep = nodeId
	}
}

// Node struct represents a node in the graph. It contains slices of incoming
// and outgoing edges.
type Node struct {
	// Id is the numeric identity of the node.
	Id int
	// InEdges are edges coming into the node.
	InEdges []*Edge
	// OutEdges are edges going out of the node.
	OutEdges []*Edge
	// Type is the type of node, one of Goal, Path, or Home.
	Type NodeType
}

func NewNode(id int, inEdges []*Edge, outEdges []*Edge, t NodeType) *Node {
	return &Node{
		Id:       id,
		InEdges:  inEdges,
		OutEdges: outEdges,
		Type:     t,
	}
}

// Run starts up a go routine for each incoming edge in the node which will
// pull ants off of the edge's channel and push them onto their next chosen
// node.
func (n *Node) Run() {
	for _, e := range n.InEdges {
		edge := e
		go n.runAnts(edge)
	}
}

// runAnts pulls ants from the incoming channel and then pushes them off on
// their chosen path. If the ant is at the goal, it wil not be pushed to
// the next channel.
func (n *Node) runAnts(e *Edge) {
	for {
		ant := <-e.Path
		next, atGoal := ant.ChooseNext(n)
		if atGoal { //ant has reached goal - no more to do
			continue
		}
		next.Path <- ant
	}
}

// MarkEdge adds depositAmt pheremone to the correct incoming edge in the
// node.
func (n *Node) MarkEdge(from int, depositAmt float64) {
	for _, e := range n.InEdges {
		if e.StartNodeId == from {
			e.AddPheremone(depositAmt)
		}
	}
}

// Edge represents a directional edge in the graph.
type Edge struct {
	// Path is the channel on which the edge moves ants from the startNode
	// to the endNode
	Path chan Ant
	// StartNodeId is Id of the starting node in the edge
	StartNodeId int
	// EndNodeId is the Id of the ending node in the edge
	EndNodeId int

	// pheremone is the amount of pheremone currently on the edge
	pheremone float64
}

// NewEdge creates a new edge with the starting pheremone amount of 10.0.
func NewEdge(startId, endId int) *Edge {
	return &Edge{
		Path:        make(chan Ant, 5),
		StartNodeId: startId,
		EndNodeId:   endId,
		pheremone:   10.0,
	}
}

// Pheremone returns the amount of pheremone present on the edge.
func (e *Edge) Pheremone() float64 {
	return e.pheremone
}

// AddPheremone adds pheremone to the edge. AddPheremone will not allow the
// amount of pheremone on the edge to go below 0.1.
func (e *Edge) AddPheremone(f float64) {
	e.pheremone = math.Max(e.pheremone+f, 0.1)
}

// String prints edges as "StartNodeId -> EndNodeId: Pheremone".
func (e Edge) String() string {
	return fmt.Sprintf("%d -> %d: %.2f", e.StartNodeId, e.EndNodeId, e.pheremone)
}
