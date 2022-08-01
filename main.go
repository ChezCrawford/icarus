package main

import (
	"context"
	"fmt"
	"icarus/internal/enphase"
	"icarus/internal/overwatch"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	fmt.Println("Let's become enlightened...")
	config := overwatch.LoadConfig()

	// client := enphase.NewClient(config.EnphaseAccessToken, config.EnphaseApiKey, config.EnphaseSystemId)
	client := enphase.NewFileClient()

	watch := overwatch.NewOverwatch(client)
	watch.Start()

	pdClient := pagerduty.NewClient(config.PdApiKey)

	var opts pagerduty.GetCurrentUserOptions
	ctx := context.Background()
	user, err := pdClient.GetCurrentUserWithContext(ctx, opts)

	fmt.Printf("User: %+v, Error: %+v", user, err)
}

func getSummary(client enphase.Client) {
	summary, _ := client.GetSystemSummary()

	fmt.Printf("System Summary: %+v\n", summary)
	startTime := overwatch.ParseTime(summary.LastReportAt)
	endTime := overwatch.ParseTime(summary.LastIntervalEndAt)
	fmt.Printf("Interval: [%+v , %+v] \n", startTime.UTC(), endTime.UTC())
}

func getConsumptionMeter(client enphase.Client) {
	consumption, _ := client.GetConsumptionMeter()

	fmt.Printf("Consumption Meter: %+v\n", consumption)
	startTime := overwatch.ParseTime(consumption.StartAt)
	endTime := overwatch.ParseTime(consumption.EndAt)
	fmt.Printf("Interval: [%+v , %+v] \n", startTime.UTC(), endTime.UTC())
}

func getProductionMeter(client enphase.Client) {
	production, _ := client.GetProductionMeter()

	fmt.Printf("Production Meter: %+v\n", production)
	startTime := overwatch.ParseTime(production.StartAt)
	endTime := overwatch.ParseTime(production.EndAt)
	fmt.Printf("Interval: [%+v , %+v] \n", startTime.UTC(), endTime.UTC())
}

func monitor(client enphase.Client) {
	consumption, _ := client.GetConsumptionMeter()
	production, _ := client.GetProductionMeter()

	generatingExcess := overwatch.Evaluate(*consumption, *production)
	fmt.Printf("Excess power being generated: %+v\n", generatingExcess)
}
