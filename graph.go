package main

import (
	"sync"
)

// TODO: Constructors/Generators for graph
// TODO: Diminishing pheromones in edges

type Edge struct {
	Path        chan Ant
	pheremone   float64
	mu          sync.Mutex
	Weight      int
	StartNodeId int
	EndNodeId   int
}

type Node struct {
	Id       int
	InEdges  []Edge
	OutEdges []Edge
	Type     NodeType
}

type Graph struct {
	Nodes []Node
}

func (n *Node) Run() {
	for _, e := range n.InEdges {
		go n.runAnts(&e)
	}
}

func (n *Node) runAnts(e *Edge) {
	for {
		ant := <-e.Path
		ant.UpdateDestination(n)
		next, done := ant.ChoosePath(n.OutEdges)
		if done {
			break
		}
		e.AddPheremone(ant.PheremoneOut())
		next.Path <- ant
	}
}

func (e *Edge) Pheremone() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.pheremone
}

func (e *Edge) AddPheremone(f float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pheremone += f
}
