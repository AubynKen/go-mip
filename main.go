package main

import (
	"fmt"
	"gomip/examples"
	"time"
)

func main() {
	fmt.Printf("Running example demos: \n\n\nKnapsack Problem\n")
	examples.KnapsackProblem()
	fmt.Println("=====================================")

	fmt.Println("\n\n\nTransportation Problem")
	examples.TransportationProblem()
	fmt.Println("=====================================")

	fmt.Println("\n\n\nProduction Planning Problem")
	examples.ProductionPlanningProblem()
	fmt.Println("=====================================")

	fmt.Println("\n\n\nSD-WAN link selection problem, time limit 10 seconds")
	examples.RoutingCtrlLinkSelection(10 * time.Second)
	fmt.Println("=====================================")
}
