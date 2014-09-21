package main

type NodeType int

const ( // named nodes
	Home NodeType = iota
	Goal
	Path
)

func main() {
	// command line input for:
	//   number of ants
	//   number of steps allowed each ant
	//   ant type ??

	// create graph
	// create ants
	// make each ant take step from home in graph
	// run ants and wait for waitgroup to be done
	// when it's done output image of final state???

}
