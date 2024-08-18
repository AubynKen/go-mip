package examples

import (
	"fmt"
	"gomip/mip"
	"log"
	"math"
)

func TransportationProblem() {
	// Create a new solver
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.ReleaseResources()

	// Define the transportation problem parameters
	sources := []string{"Factory1", "Factory2"}
	destinations := []string{"Store1", "Store2", "Store3"}

	// Supply at each source
	supply := map[string]float64{
		"Factory1": 100,
		"Factory2": 150,
	}

	// Demand at each destination
	demand := map[string]float64{
		"Store1": 80,
		"Store2": 70,
		"Store3": 90,
	}

	// Transportation cost from each source to each destination
	cost := map[string]map[string]float64{
		"Factory1": {"Store1": 2, "Store2": 3, "Store3": 1},
		"Factory2": {"Store1": 5, "Store2": 4, "Store3": 6},
	}

	// Create variables for transportation quantities
	vars := make(map[string]map[string]*mip.Variable)
	for _, source := range sources {
		vars[source] = make(map[string]*mip.Variable)
		for _, dest := range destinations {
			vars[source][dest] = solver.VarFloat(fmt.Sprintf("x_%s_%s", source, dest), 0, math.MaxFloat64)
		}
	}

	// Add supply constraints
	for _, source := range sources {
		exp := mip.NewLinearExpression()
		for _, dest := range destinations {
			exp.AddVar(vars[source][dest])
		}
		solver.AddConstraintExpr(exp, mip.LessThanOrEqual, supply[source])
	}

	// Add demand constraints
	for _, dest := range destinations {
		exp := mip.NewLinearExpression()
		for _, source := range sources {
			exp.AddVar(vars[source][dest])
		}
		solver.AddConstraintExpr(exp, mip.GreaterThanOrEqual, demand[dest])
	}

	// Set the objective function (minimize total transportation cost)
	obj := mip.NewLinearExpression()
	for _, source := range sources {
		for _, dest := range destinations {
			obj.AddTerm(vars[source][dest], cost[source][dest])
		}
	}
	solver.SetObjective(obj, mip.Minimize)

	// Solve the problem
	_, err = solver.Solve(-1)
	if err != nil {
		log.Fatalf("Error solving the problem: %v", err)
	}

	// Print out the solution that we've found
	fmt.Println("Transportation plan:")
	for _, source := range sources {
		for _, dest := range destinations {
			fmt.Printf("From %s to %s: %.2f\n", source, dest, vars[source][dest].Value())
		}
	}

	fmt.Printf("Total cost: %.2f\n", solver.ObjectiveValue())
}
