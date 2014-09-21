package main

import (
	"sync"
	"time"
)

// TODO: Constructors/Generators for graph

type Graph struct {
	Nodes []Node
}

type Node struct {
	Id          int
	InEdges     []Edge
	OutEdges    []Edge
	Type        NodeType
	DecayFactor float64
}

func (n *Node) Run(done chan int) {
	for _, e := range n.InEdges {
		go n.runAnts(&e)
	}

	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-done: //all ants have completed, exit function
			return
		case <-ticker.C:
			for _, e := range n.InEdges {
				e.AddPheremone(n.DecayFactor)
			}
		}
	}
}

func (n *Node) runAnts(e *Edge) {
	for {
		ant := <-e.Path
		next, err := ant.ChoosePath(n)
		if err {
			continue
		}
		e.AddPheremone(ant.PheremoneAmt())
		next.Path <- ant
	}
}

type Edge struct {
	Path        chan Ant
	StartNodeId int
	EndNodeId   int

	pheremone float64
	mu        sync.Mutex
}

//TODO: Do this using channels instead of locks to make more "go-like"
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
