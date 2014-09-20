package main

// TODO: Constructors/Generators for graph
// TODO: Stopping when finished
// TODO: Diminishing pheromones in edges

type Edge struct {
	Path        chan Ant
	Pheremone   int
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
		next := ant.ChoosePath(n.OutEdges)
		ant.AddPheremone(&next)
		next.Path <- ant
	}
}
