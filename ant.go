package acogo

import (
	"sync"
)

// Ant is an interface for all Ants used in acogo
type Ant interface {
	ChooseNext(*Node) (*Edge, bool)
	MarkPath(*Graph)
}

// SimpleAnt is the most basic Ant. It probabilistically chooses a path based on
// the relative amount of pheremone in all paths.
type SimpleAnt struct {
	// Last node visited
	LastNodeId int

	// A slice of NodeIds visited in the current move toward goal
	StepsTaken []int
	// Amount of pheremone left at each step along the path
	DepositAmt float64
	// Source of randomness for making probabilistic path decisions
	RandomSrc chan chan float64
	// WaitGroup for reporting back to the system when the goal has been reached
	waitGroup *sync.WaitGroup
}

// NewSimpleAnt creates a SimpleAnt with the input parameters.
func NewSimpleAnt(lastNodeId int, depositAmt float64, randSrc chan chan float64, wg *sync.WaitGroup) *SimpleAnt {
	return &SimpleAnt{
		LastNodeId: lastNodeId,
		DepositAmt: depositAmt,
		StepsTaken: make([]int, 0, 100),
		RandomSrc:  randSrc,
		waitGroup:  wg,
	}
}

// ChooseNext probabilistically chooses the next edge in the graph to move down
// based on the amount of pheremone on each edge. ChooseNext will not choose
// an edge leading to the node the ant just visited unless it is the only edge
// available on the node.
func (a *SimpleAnt) ChooseNext(node *Node) (*Edge, bool) {
	a.StepsTaken = append(a.StepsTaken, node.Id)

	if node.Type == Goal {
		a.waitGroup.Done()
		return nil, true
	}

	total := a.sumPheremones(node.OutEdges)

	// use RandomSrc to get random float64
	randChan := make(chan float64)
	a.RandomSrc <- randChan
	choice := <-randChan

	pos := 0.0
	for _, e := range node.OutEdges {
		if e.EndNodeId != a.LastNodeId {
			pos += e.Pheremone()
			if choice <= pos/total {
				a.LastNodeId = node.Id
				return e, false
			}
		}
	}
	a.LastNodeId = node.Id
	return node.OutEdges[len(node.OutEdges)-1], false
}

// MarkPath lays down pheremone based on the path the ant took to from home
// to goal. MarkPath will "unloop" the path meaning that any loops in the
// original path will be eliminated.
func (a *SimpleAnt) MarkPath(g *Graph) {
	unlooped := unloop(a.StepsTaken)
	g.MarkPath(unlooped, a.DepositAmt)
}

// unloop takes the steps taken and eliminates any loops. This is done by always
// choosing the last instance that a given node was passed through, e.g. if the
// sequence of nodes is [1, 2, 4, 5, 2, 6, 7], it will become [1, 2, 6, 7].
func unloop(steps []int) []int {
	// create a map of the last index at which each node occurs
	lastIndices := make(map[int]int, 100)
	for idx, node := range steps {
		lastIndices[node] = idx
	}

	idx := 0
	unlooped := make([]int, 0, len(steps))
	for idx < len(steps) {
		node := steps[idx]
		unlooped = append(unlooped, node)
		if idx == lastIndices[node] {
			idx++
		} else {
			idx = lastIndices[node] + 1
		}
	}
	return unlooped
}

// sumPheremones totals the pheremones present on the slice of edges passed in.
// sumPheremones does not count pheremones from the edge the ant most recently
// visited.
func (a *SimpleAnt) sumPheremones(edges []*Edge) float64 {
	total := 0.0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			total += e.Pheremone()
		}
	}
	return total
}
