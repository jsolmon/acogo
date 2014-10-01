package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var antCount = *flag.Int("antcount", 20, "the number of ants to create")
	var depositAmt = *flag.Float64("depositamt", 1.0, "amount of pheremone deposited by ant")
	var iterations = *flag.Int("iterations", 500, "the number times each ant will find the goal node")
	var decayFactor = *flag.Float64("decay", -0.3, "the amount of pheremone dissipated after each round")
	var dimension = *flag.Int("dimension", 6, "the number of vertices on each side of the square graph")
	var startNode = *flag.Int("start", 0, "the index of the node where ants begin")
	var goalNode = *flag.Int("goal", dimension*dimension-1, "the index of the vertex ants are trying to reach")

	// initialize randomness source

	// create and start graph
	graph := NewGraph(dimension, startNode, goalNode, decayFactor)
	graph.Run()

	// initialize source of randomness
	randSource := RandomSource{make(chan chan float64), rand.New(rand.NewSource(time.Now().Unix()))}
	go randSource.Run()

	// add ants to graph via an in-edge on the home node
	homeNode := graph.Nodes[graph.HomeIdx]
	startEdge := homeNode.InEdges[0]

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		wg.Add(antCount)
		ants := make([]Ant, 0, antCount)

		for i := 0; i < antCount; i++ {
			ants = append(ants, NewSimpleAnt(graph.HomeIdx, depositAmt, randSource.requestChan, &wg))
			startEdge.Path <- ants[i]
		}

		// wg completes once all ants have reached the goal node
		wg.Wait()

		// cycle through ants and update pheremone based on their paths
		for _, ant := range ants {
			ant.MarkPath(graph)
		}

		graph.Dissipate()
	}

	// TODO: Visualization
	for _, n := range graph.Nodes {
		for _, e := range n.InEdges {
			fmt.Printf("%v\n", e)
		}
	}
}

type RandomSource struct {
	requestChan chan chan float64
	rand        *rand.Rand
}

func (r *RandomSource) Run() {
	for {
		respChan := <-r.requestChan
		respChan <- r.rand.Float64()
	}
}
