Package acogo provides a framework for playing with ant colony optimization.

`acogo` will run a set number of ants from start to home in the graph and then
output a `DOT` language graph of the resulting graph pheromone levels.

Build
-----
    godep go build

Run
---
To write the output `DOT` graph to stdout, run

    ./acogo

To transform the output `DOT` graph to a png, run

    ./acogo | dot -Tpng -o antgraph.png


Command line options
---

	-antcount: The number of ants to create and run each iteration. Default 20.
	-depositamt: The amount of pheromone each ant deposits on their path. Default 1.0.
	-iterations: The number of times each ant runs from home to goal. Default 500.
	-decay: The amount that pheromone on each edge decreases after each iteration. Default 0.3.
	-dimension: The number of nodes on each side of the square graph. Default 6.
	-start: The index of the node from which the ants start. Default 0.
	-goal: The index of the node that ants are trying to reach. Default dimension * dimension.

Description
-----------

When run, acogo will create a square graph of size `dimension * dimension` with each
node having connections to adjacent nodes above, below, left, right, and on all
four diagonals. For example, if `dimension` is 3, the graph look like:

	0 1 2
	3 4 5
	6 7 8

`antcount` ants will be placed at the `start` node on the graph and travel in the graph
until they reach the `goal` node. At each node, the ant will probabilistically chose
which node to travel to next proportionate to the amount of pheromone on each
outgoing node. Ants will not return to the node they just left unless that is the
only option for exiting a particular node.

Once all ants reach the goal node in a given iteration, `depositamt` pheromone will
be added to each edge that each ant traveled on. After reaching the goal, ants
routes are unlooped, so an ant that traveled 1->4->5->2->4->8 would only lay down
pheromone from 1->4->8. After all ant pheromone has been laid down, `decay` pheromone
is subtracted from all edges.

When all ants have completed `iterations` iterations, a DOT language representation
of the final graph is written to `stdout`. Edge colors reflect how much pheromone
is on a given edge with darker edges representing more pheromone.
