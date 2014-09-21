package main

import (
	"math/rand"
	"time"
)

type Ant interface {
	UpdateDestination(*Node)
	ChoosePath([]Edge) (*Edge, bool)
	PheremoneOut() float64
}

type SimpleAnt struct {
	LastNodeId  int
	Destination NodeType

	MaxSteps   int
	StepsCount int
}

func NewSimpleAnt(maxSteps int) *SimpleAnt {
	a := SimpleAnt{LastNodeId: 0,
		Destination: Goal,
		MaxSteps:    maxSteps,
		StepsCount:  0}
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

func (a *SimpleAnt) ChoosePath(edges []Edge) (*Edge, bool) {
	if a.StepsCount >= a.MaxSteps {
		return nil, true
	}
	a.StepsCount += 1

	total := a.sumPheremones(edges)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	choice := r.Float64()

	pos := 0.0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			pos += e.Pheremone()
			if choice <= pos/total {
				a.LastNodeId = e.StartNodeId
				return &e, false
			}
		}
	}
	return &edges[len(edges)-1], false
}

func (a *SimpleAnt) PheremoneOut() float64 {
	return 1.0
}

func (a *SimpleAnt) sumPheremones(edges []Edge) float64 {
	total := 0.0
	for _, e := range edges {
		if e.EndNodeId != a.LastNodeId {
			total += e.Pheremone()
		}
	}
	return total
}
