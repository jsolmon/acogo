// Copyright (c) 2014 Joanna Solmon. All rights reserved.
// Use of this source code is governed by the MIT License found in the LICENSE file.

/*
Package acogo provides a framework for playing with ant colony optimization.

Command line options:

	antcount: The number of ants to create and run each iteration. Default 20.
	depositamt: The amount of pheremone each ant deposits on their path. Default 1.0.
	iterations: The number of times each ant runs from home to goal. Default 500.
	decay: The amount that pheromone on each edge decreases after each iteration. Default 0.3.
	dimension: The number of nodes on each side of the square graph. Default 6.
	start: The index of the node from which the ants start. Default 0.
	goal: The index of the node that ants are trying to reach. Default dimension * dimension.

When run, acogo will create a square graph of size dimension * dimension with each
node having connections to adjacent nodes above, below, left, right, and on all
four diagonals. For example, if dimension is 3, the graph look like:

	0 1 2
	3 4 5
	6 7 8

antcount ants will be placed at the start node on the graph and travel in the graph
until they reach the goal node. At each node, the ant will probabilitstically chose
which node to travel to next proportionate to the amount of pheremone on each
outgoing node. Ants will not return to the node they just left unless that is the
only option for exiting a particular node.

Once all ants reach the goal node in a given iteration, depositamt pheremone will
be added to each edge that each ant traveled on. After reaching the goal, ants
routes are unlooped, so an ant that traveled 1->4->5->2->4->8 would only lay down
pheremone from 1->4->8. After all ant pheremone has been laid down, decay pheremone
is subtracted from all edges.

When all ants have completed iterations iterations, the pheremone values for all
edges are printed to the screen. (Visualization will be added soon!)
*/
package acogo

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
	var decayFactor = *flag.Float64("decay", 0.3, "the amount of pheremone dissipated after each round")
	var dimension = *flag.Int("dimension", 6, "the number of nodes on each side of the square graph")
	var startNode = *flag.Int("start", 0, "the index of the node where ants begin")
	var goalNode = *flag.Int("goal", dimension*dimension-1, "the index of the vertex ants are trying to reach")

	// initialize randomness source

	// create and start graph
	graph := NewGraph(dimension, startNode, goalNode, decayFactor)
	graph.Run()

	// initialize source of randomness
	randSource := RandomSource{make(chan chan float64), rand.New(rand.NewSource(time.Now().Unix()))}
	randSource.Run()

	// add ants to graph via an in-edge on the home node
	homeNode := graph.Nodes[graph.HomeIdx]
	startEdge := homeNode.InEdges[0]

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		wg.Add(antCount)
		ants := make([]Ant, 0, antCount)

		for i := 0; i < antCount; i++ {
			ants = append(ants, NewSimpleAnt(graph.HomeIdx, depositAmt, randSource.RequestChan, &wg))
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

// RandomSource provides a shared source of randomness for ants via requests
// (chan float64s) sent to its RequestChan.
type RandomSource struct {
	RequestChan chan chan float64
	rand        *rand.Rand
}

// Run initializes a RandomSource to run and wait for requests on its RequestChan.
func (r *RandomSource) Run() {
	go func() {
		for {
			respChan := <-r.RequestChan
			respChan <- r.rand.Float64()
		}
	}()
}
