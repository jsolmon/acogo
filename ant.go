package main

import (
	"sync"
)

type Ant interface {
	ChoosePath(*Node) (*Edge, bool)
	MarkPath(*Graph)
}

type SimpleAnt struct {
	LastNodeId int

	StepsTaken []int
	DepositAmt float64
	RandomSrc  chan chan float64
	waitGroup  *sync.WaitGroup
}

func NewSimpleAnt(lastNodeId int, depositAmt float64, randSrc chan chan float64, wg *sync.WaitGroup) *SimpleAnt {
	return &SimpleAnt{
		LastNodeId: lastNodeId,
		DepositAmt: depositAmt,
		StepsTaken: make([]int, 0, 100),
		RandomSrc:  randSrc,
		waitGroup:  wg,
	}
}

func (a *SimpleAnt) ChoosePath(node *Node) (*Edge, bool) {
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

func (a *SimpleAnt) MarkPath(g *Graph) {
	unlooped := unloop(a.StepsTaken)
	g.MarkPath(unlooped, a.DepositAmt)
}

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

func (a *SimpleAnt) PheremoneAmt() float64 {
	return 1.0
}

func (a *SimpleAnt) sumPheremones(edges []*Edge) float64 {
	total := 0.0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			total += e.Pheremone()
		}
	}
	return total
}
