package main

import (
	"flag"
	"fmt"
	"sync"
)

type NodeType int

var (
	antCount       int
	stepCount      int
	graphDimension int
	startVertex    int
	goalVertex     int
)

func init() {
	flag.IntVar(&antCount, "ant count", 10, "the number of ants to create")
	flag.IntVar(&stepCount, "step count", 2000, "the number of steps each ant takes")
	flag.IntVar(&graphDimension, "graph dimension", 9, "the number of vertices on each side of the square graph")
	flag.IntVar(&startVertex, "start vertex", 0, "the index of the vertex where ants begin")
	flag.IntVar(&goalVertex, "goal vertex", 81, "the index of the vertex ants are trying to reach")
}

func main() {
	flag.Parse()

	// create and start graph
	graph := NewGraph(graphDimension, startVertex, goalVertex)
	doneChan := make(chan struct{})
	graph.Run(doneChan)

	// Add ants to graph via an in-edge on the home node
	homeNode := graph.Nodes[graph.HomeIdx]
	startEdge := homeNode.InEdges[0]

	var wg sync.WaitGroup
	wg.Add(antCount)
	for i := 0; i < antCount; i++ {
		startEdge.Path <- NewSimpleAnt(stepCount, &wg)
	}

	// wg completes once all ants have taken antCount steps
	wg.Wait()
	//Stop all graph nodes from decreasing pheremone
	close(doneChan)

	// TODO: Visualization
	for _, n := range graph.Nodes {
		for _, e := range n.InEdges {
			fmt.Printf("%v\n", e)
		}
	}
}
