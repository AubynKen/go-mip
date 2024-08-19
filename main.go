package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"gomip/examples"
	"time"
)

func printTitle(title string) {
	fmt.Println(text.Bold.Sprint(text.FgHiGreen.Sprint("\n\n", title)))
}

func main() {
	printTitle("Knapsack Problem")
	examples.KnapsackProblem()

	printTitle("Transportation Problem")
	examples.TransportationProblem()

	printTitle("Production Planning Problem")
	examples.ProductionPlanningProblem()

	timeLimit := 10 * time.Second
	printTitle(fmt.Sprintf("SD-WAN link selection problem, time limit = %s", timeLimit))
	examples.RoutingCtrlLinkSelection(timeLimit)
}
