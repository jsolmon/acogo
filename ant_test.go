package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestChooseNext(t *testing.T) {
	randSource := RandomSource{make(chan chan float64), rand.New(rand.NewSource(time.Now().Unix()))}
	go randSource.Run()

	var wg sync.WaitGroup
	ant := NewSimpleAnt(1, 1.0, randSource.RequestChan, &wg)

	edges := []*Edge{NewEdge(0, 1), NewEdge(0, 6), NewEdge(0, 3), NewEdge(0, 10)}
	edges[0].pheremone = 5.0
	edges[1].pheremone = 3.0
	edges[2].pheremone = 6.0
	edges[3].pheremone = 1.0

	node := NewNode(0, []*Edge{}, edges, Path)

	edge0Count := 0
	edge1Count := 0
	edge2Count := 0
	edge3Count := 0

	for i := 0; i < 1000; i++ {
		ant.LastNodeId = 1
		choice, _ := ant.ChooseNext(node)
		switch choice {
		case edges[0]:
			edge0Count++
		case edges[1]:
			edge1Count++
		case edges[2]:
			edge2Count++
		case edges[3]:
			edge3Count++
		}
	}
	// Should never choose edge0 because that is the LastNodeId
	if edge0Count != 0 {
		t.Error(fmt.Sprintf("edge0 count should be 0 but was %v\n", edge0Count))
	}
	if edge1Count > 350 || edge1Count < 250 {
		t.Error(fmt.Sprintf("edge1 count should be between 250 and 350 but was %v\n", edge1Count))
	}
	if edge2Count > 650 || edge2Count < 550 {
		t.Error(fmt.Sprintf("edge2 count should be between 550 and 650 but was %v\n", edge2Count))
	}
	if edge3Count > 130 || edge1Count < 70 {
		t.Error(fmt.Sprintf("edge3 count should be between 70 and 130 but was %v\n", edge3Count))
	}
}

func TestSumPheremones(t *testing.T) {
	var wg sync.WaitGroup
	ant := NewSimpleAnt(1, 1.0, make(chan chan float64), &wg)
	edges := []*Edge{NewEdge(0, 1), NewEdge(0, 6), NewEdge(0, 3), NewEdge(0, 10)}
	edges[0].pheremone = 1.0
	edges[1].pheremone = 2.0
	edges[2].pheremone = 4.0
	edges[3].pheremone = 8.0

	// should sum all but 0->1 node which matches ant's LastNodeId
	pheremone := ant.sumPheremones(edges)
	expected := 14.0

	if expected != pheremone {
		t.Error(fmt.Sprintf("Pheremone sum should be %v, got %v", expected, pheremone))
	}
}

func TestUnloop(t *testing.T) {
	testUnloop([]int{1, 2, 3, 2, 5, 4, 7, 5, 6, 10}, []int{1, 2, 5, 6, 10}, t)
	testUnloop([]int{3, 4, 5, 2, 3, 5, 4, 3, 5, 6, 7}, []int{3, 5, 6, 7}, t)
	testUnloop([]int{3, 4, 5, 4, 5, 4, 3, 6, 2}, []int{3, 6, 2}, t)
}

func testUnloop(input, expected []int, t *testing.T) {
	unlooped := unloop(input)
	if !reflect.DeepEqual(expected, unlooped) {
		t.Error(fmt.Sprintf("Unloop failed: expected %v, got %v", expected, unlooped))
	}
}
