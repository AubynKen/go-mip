package examples

import (
	"fmt"
	"gomip/mip"
	"log"
)

func ProductionPlanningProblem() {
	// Create a new solver
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.ReleaseResources()

	// Define the production planning problem parameters
	products := []string{"chairs", "tables"}
	resources := []string{"wood", "labor"}

	// Resource usage per unit of product
	usage := map[string]map[string]float64{
		"chairs": {"wood": 5, "labor": 5},
		"tables": {"wood": 12, "labor": 6},
	}

	// Available resources
	available := map[string]float64{
		"wood":  1200,
		"labor": 800,
	}

	// Profit per unit of product
	profit := map[string]float64{
		"chairs": 10,
		"tables": 20,
	}

	// Create variables for each product
	vars := make(map[string]*mip.Variable)
	for _, product := range products {
		vars[product] = solver.VarInt(fmt.Sprintf("x_%s", product), 0, 100000)
	}

	// Add resource constraints
	for _, resource := range resources {
		exp := mip.NewLinearExpression()
		for _, product := range products {
			exp.AddTerm(vars[product], usage[product][resource])
		}
		solver.AddConstraintExpr(exp, mip.LessThanOrEqual, available[resource])
	}

	// Set the objective function (maximize total profit)
	obj := mip.NewLinearExpression()
	for _, product := range products {
		obj.AddTerm(vars[product], profit[product])
	}
	solver.SetObjective(obj, mip.Maximize)

	// Solve the problem
	_, err = solver.Solve(-1)
	if err != nil {
		log.Fatalf("Error solving the problem: %v", err)
	}

	// Print out the production plan
	fmt.Println("Production plan:")
	for product, v := range vars {
		fmt.Printf("%s: %.2f\n", product, v.Value())
	}

	fmt.Printf("Total profit: %.2f\n", solver.ObjectiveValue())

	fmt.Println("Resource usage:")
	for resource := range usage["chairs"] {
		totalUsage := 0.
		for product, v := range vars {
			totalUsage += usage[product][resource] * v.Value()
		}
		fmt.Printf("%s: %.2f\n", resource, totalUsage)
	}

	fmt.Println("Available resources:")
	for resource, amount := range available {
		fmt.Printf("%s: %.2f\n", resource, amount)
	}

	fmt.Println("Profit per unit of product:")
	for product, p := range profit {
		fmt.Printf("%s: %.2f\n", product, p)
	}
}
