# MIP in Go demo project

This project implements a simple mixed-integer programming (MIP) model in Go using the OR-Tools library.
We use the C++ api of OR-Tools through cgo.

## Structure
`bridge`: the C++ bridging/translation layer code with extern C that interfaces with the OR-Tools library.

`mip`: the Go wrapper code that interfaces with the bridging layer C code with CGO.

`examples`: example MIP models that can be solved using the Go wrapper code.

## How to run the demo
You may manually install the OR-Tools library and its dependencies, or use the provided Dockerfile to build a Docker image that contains the OR-Tools library and the Go code.

To run the demo using Docker, execute the following commands:
```bash
docker build -t mip-go-demo .
docker run mip-go-demo
```

The following should normally work on arm chips:
```bash
 docker run --rm --platform linux/amd64 -it $(docker build --platform linux/amd64 -q .) 
```

## Architecture
Note since Google does not provide the binaries for arm64 architecture, we've run the tests on amd64 architecture only.

We've tested the docker image on MacOS with Apple Silicon under Rosetta 2. No testing was done on native amd64 architecture machines.

## Super brief explanations for some of the examples in `examples` folder

### Production Planning
The Production Planning Problem involves deciding how much of each product to manufacture to maximize profit, given limited resources. Key components:

Variables: Quantities of each product to produce (chairs and tables)
Constraints: Limited resources (wood and labor)
Objective: Maximize total profit

This problem is common in manufacturing, where companies need to optimize their product mix based on resource availability and profitability.

### Transportation Problem
The Transportation Problem focuses on finding the most cost-effective way to transport goods from multiple sources to multiple destinations. Key components:

Variables: Quantity to ship from each source to each destination
Constraints: Supply limits at sources, demand requirements at destinations
Objective: Minimize total transportation cost

This problem is crucial in logistics, helping companies optimize their distribution networks and reduce shipping costs.

### KnapSack Problem
The Bin Packing Problem involves packing items of different sizes into a minimum number of fixed-size bins. Key components:

Variables: Binary decisions for placing each item in each bin
Constraints: Ensure each item is placed in exactly one bin and bin capacities are not exceeded
Objective: Minimize the number of bins used

This problem has applications in various fields, including logistics (packing goods into containers), computer science (allocating processes to servers), and manufacturing (cutting stock problem).

## SD-WAN Link Selection Optimization, detailed explanation

### Problem Overview

The SD-WAN link selection problem involves optimizing the assignment of network links to different prefix groups while balancing performance and device usage. This is a complex optimization problem with many variables and constraints.

### Key Components

Prefix Groups: Different network destinations or services.

Links: Network connections, each with its own capacity, latency, and packet loss rate.

Devices: Network equipment that links are associated with.

### Optimization Goals

Assign links to prefix groups to meet capacity requirements.

Minimize latency and packet loss.

Balance the load across devices.

note: for more detailed explanation of the mathematical formulation, please look at the python demo jupyter notebook.

### Solving Process and Time Limits
In ideal circumstances, we would solve the problem to find the mathematically optimal solution. However, for large instances like this SD-WAN problem, finding the absolute best solution can be extremely time-consuming or even impractical.

To address this, we opt for an [anytime optimization algorithm](https://en.wikipedia.org/wiki/Anytime_algorithm):

Instead of making the underlying solver run until the optimal is found, we set a time limit (in this case, 10 seconds) for the solver to find a good solution. - it allows us to get a good, feasible solution within a reasonable timeframe, even if it's not guaranteed to be the absolute best.

### Understanding and Measuring Solution Quality: Bounds and Gap
To evaluate how good our solution is, we use the concept of bounds:

Upper Bound (of the optimal objective value): Intuitively, an upper bound (there's an infinity of them) is a value that is guaranteed to be greater (worse, in the case of an minimization problem) than the optimal. In our case, any feasible solution provides an objective value that is an upper bound of the optimal objective vlaue.

Lower Bound: (of the optimum) This is a theoretical lower limit on the optimal objective value. The solver calculates this based on relaxations of the problem and other mathematical techniques. Initially, the lower bound is minus infinity. As the solver progresses, it improves the lower bound to get closer to the optimal objective value. For example, if at some point, we stop the optimization process, and the lower bound is 100, it means the optimal objective value is proven mathematically by the solver to be at least 100.

Optimality Gap: This is the difference between the upper and lower bounds, usually expressed as a percentage. It tells us how close our current best solution is to the theoretical optimum.
Gap = (Upper Bound - Lower Bound) / Lower Bound * 100%

For example, if we have an upper bound (objective value of the best feasible solution so far, as fore-mentioned) of 115 and a lower bound of 100 (i.e. the solver proves that the optimal objective value is at least 100), the gap is 15%. This means the optimal objective lies somewhere between 100 and 115, and we are AT MOST 15% away from the optimal (it could be less, but we don't know).

Interpreting the Results
In the code, we see these concepts applied in the SD-WAN example demo.

```go
fmt.Printf("Best Objective Value Found: %f\n", solver.ObjectiveValue())
// This is our upper bound - the best solution found within the time limit

if isOptimal {
    fmt.Println("The solution found is proven to be optimal!")
} else {
    fmt.Println("The solution found is not guaranteed to be optimal.")
    fmt.Printf("The solver proved that the optimal objective is no less than %f\n", solver.BestBound())
    // This is our lower bound
    fmt.Printf("Which means that our solution is within %.2f%% of the optimal.\n",
        solver.Gap()*100)
    // This prints the optimality gap as a percentage
}
```

If isOptimal is true, we've found and proved the optimal solution.
If not, we report the best solution found (upper bound), the best bound (lower bound), and the gap.

Note: Only tested on CBC solver.
No guarantee it would work with other solvers.

### Why This Matters

Practical Solutions: We get a good, usable solution within a reasonable time.

Quality Assurance: The gap tells us how much room for improvement there might be.
Decision-Making: Understanding the gap helps in deciding whether to accept the current solution or allocate more time to potentially improve it.

By using this approach, we balance the need for a good solution with the practical constraints of time and computational resources in complex network optimization problems.

### What I would have added if I had more time

It would be interesting to evaluate how much the current heuristic deviates from the optimum to assess whether it's worth it to use a sophisticated approach
like mixed-integer programming. 

I would also have tried to use a heuristic to generate a warm-start solution for the MIP model, to initiate the MIP solver with a good starting point. It might be interesting performance wise.
