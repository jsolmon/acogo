Package acogo provides a framework for playing with ant colony optimization.

Build
-----
    go build

Run
---
    ./acogo

Command line options:

	antcount: The number of ants to create and run each iteration. Default 20.
	depositamt: The amount of pheremone each ant deposits on their path. Default 1.0.
	iterations: The number of times each ant runs from home to goal. Default 500.
	decay: The amount that pheromone on each edge decreases after each iteration. Default 0.3.
	dimension: The number of nodes on each side of the square graph. Default 6.
	start: The index of the node from which the ants start. Default 0.
	goal: The index of the node that ants are trying to reach. Default dimension * dimension.

Description
-----------

When run, acogo will create a square graph of size `dimension * dimension` with each
node having connections to adjacent nodes above, below, left, right, and on all
four diagonals. For example, if `dimension` is 3, the graph look like:

	0 1 2
	3 4 5
	6 7 8

`antcount` ants will be placed at the `start` node on the graph and travel in the graph
until they reach the `goal` node. At each node, the ant will probabilitstically chose
which node to travel to next proportionate to the amount of pheremone on each
outgoing node. Ants will not return to the node they just left unless that is the
only option for exiting a particular node.

Once all ants reach the goal node in a given iteration, `depositamt` pheremone will
be added to each edge that each ant traveled on. After reaching the goal, ants
routes are unlooped, so an ant that traveled 1->4->5->2->4->8 would only lay down
pheremone from 1->4->8. After all ant pheremone has been laid down, decay pheremone
is subtracted from all edges.

When all ants have completed `iterations` iterations, the pheremone values for all
edges are printed to the screen. (Visualization will be added soon!)
