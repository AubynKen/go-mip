package examples

import (
	"fmt"
	"gomip/mip"
	"log"
	"math"
)

func ProductionPlanningProblem() {
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.ReleaseResources()

	products := []string{"chairs", "tables"}
	resources := []string{"wood", "labor"}

	// costs in resources for each product
	cost := map[string]map[string]float64{
		"chairs": {"wood": 5, "labor": 5},
		"tables": {"wood": 12, "labor": 6},
	}
	fmt.Println("We're making some chairs and tables.")
	fmt.Printf("Each chair requires %d wood and %d labor.\n", int(cost["chairs"]["wood"]), int(cost["chairs"]["labor"]))
	fmt.Printf("Each table requires %d wood and %d labor.\n", int(cost["tables"]["wood"]), int(cost["tables"]["labor"]))

	// our available resources
	available := map[string]float64{
		"wood":  1200,
		"labor": 800,
	}
	fmt.Printf("\nAvailable resource: \nWe have %d wood and %d labor available.\n", int(available["wood"]), int(available["labor"]))

	// profit in $ per unit for each product
	profit := map[string]float64{
		"chairs": 10,
		"tables": 20,
	}
	fmt.Printf("\nWe make $%d profit per chair and $%d profit per table.\n", int(profit["chairs"]), int(profit["tables"]))

	// vars[p] is the number of units that we want to produce for product p
	vars := make(map[string]*mip.Variable)
	for _, product := range products {
		vars[product] = solver.VarInt(fmt.Sprintf("x_%s", product), 0, math.MaxInt)
	}

	// used resource <= available resource for each resource
	for _, resource := range resources {
		exp := mip.NewLinearExpression()
		for _, product := range products {
			exp.AddTerm(vars[product], cost[product][resource])
		}
		solver.AddConstraintExpr(exp, mip.LessThanOrEqual, available[resource])
	}

	// objective: maximize total profit
	obj := mip.NewLinearExpression()
	for _, product := range products {
		obj.AddTerm(vars[product], profit[product])
	}
	solver.SetObjective(obj, mip.Maximize)

	_, err = solver.Solve(-1) // run until optimum found
	if err != nil {
		log.Fatalf("Error solving the problem: %v", err)
	}

	fmt.Println("\n\nProduction plan:")
	for product, v := range vars {
		// print out how much of each product we should produce
		fmt.Printf("%s: %.2f\n", product, v.Value())
	}

	fmt.Printf("\n\nTotal profit: %.2f\n", solver.ObjectiveValue())

	fmt.Println("\n\nResource usages:") // Print out the total resource usages in used/available format
	for resource := range available {
		var used float64
		// if we produce 5 chairs, and each chair requires 2 wood, then we used 5*2 = 10 wood for chairs
		// we sum over all products to get the total wood usage, and so on
		for _, product := range products {
			used += vars[product].Value() * cost[product][resource]
		}
		fmt.Printf("%s: %.2f/%.2f\n", resource, used, available[resource])
	}
}
