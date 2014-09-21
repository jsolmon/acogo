package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Ant interface {
	ChoosePath(*Node) (*Edge, bool)
	PheremoneAmt() float64
}

type SimpleAnt struct {
	LastNodeId  int
	Destination NodeType

	MaxSteps   int
	StepsCount int
	waitGroup  *sync.WaitGroup
}

func NewSimpleAnt(maxSteps int, wg *sync.WaitGroup) *SimpleAnt {
	a := SimpleAnt{
		LastNodeId:  0,
		Destination: Goal,
		MaxSteps:    maxSteps,
		StepsCount:  0,
		waitGroup:   wg}
	return &a
}

func (a *SimpleAnt) updateDestination(n *Node) {
	if a.Destination == n.Type {
		if a.Destination == Home {
			a.Destination = Goal
		} else if a.Destination == Goal {
			a.Destination = Home
		}
	}
}

func (a *SimpleAnt) ChoosePath(node *Node) (*Edge, error) {
	if a.StepsCount >= a.MaxSteps {
		a.waitGroup.Done()
		return nil, errors.New("maximum steps reached")
	}
	a.StepsCount += 1
	a.updateDestination(node)

	total := a.sumPheremones(node.OutEdges)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	choice := r.Float64()

	pos := 0.0
	for _, e := range node.OutEdges {
		if e.EndNodeId != a.LastNodeId {
			pos += e.Pheremone()
			if choice <= pos/total {
				a.LastNodeId = e.StartNodeId
				return &e, nil
			}
		}
	}
	return &node.OutEdges[len(node.OutEdges)-1], nil
}

func (a *SimpleAnt) PheremoneAmt() float64 {
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
