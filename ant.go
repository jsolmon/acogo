package main

import (
	"math/rand"
	"time"
)

type Ant interface {
	UpdateDestination(*Node)
	ChoosePath([]Edge) Edge
	AddPheremone(*Edge)
}

type SimpleAnt struct {
	LastNodeId  int
	Destination NodeType
}

func NewSimpleAnt() *SimpleAnt {
	a := SimpleAnt{LastNodeId: 0, Destination: Goal}
	return &a
}

func (a *SimpleAnt) UpdateDestination(n *Node) {
	if a.Destination == n.Type {
		if a.Destination == Home {
			a.Destination = Goal
		} else if a.Destination == Goal {
			a.Destination = Home
		}
	}
}

func (a *SimpleAnt) ChoosePath(edges []Edge) *Edge {
	total := a.sumPheremones(edges)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	choice := r.Intn(total)

	pos := 0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			pos = pos + e.Pheremone
			if choice <= pos {
				a.LastNodeId = e.StartNodeId
				return &e
			}
		}
	}
	return &edges[len(edges)-1]
}

func (a *SimpleAnt) AddPheremone(e *Edge) {
	e.Weight = e.Weight + 1
}

func (a *SimpleAnt) sumPheremones(edges []Edge) int {
	total := 0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			total = total + e.Pheremone
		}
	}
	return total
}
