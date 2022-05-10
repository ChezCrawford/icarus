package main

import (
	"fmt"
	"icarus/internal/enphase"
	"time"
)

const AccessToken string = "some-token"
const ApiKey string = "some-key"
const SystemId string = "123"

func main() {
	fmt.Println("Let's become enlightened...")
	getSummary()
}

func getSummary() {
	client := enphase.NewClient(AccessToken, ApiKey)
	summary, _ := client.GetSystemSummary(SystemId)

	// fmt.Printf("System Summary: %+v\n", summary)
	fmt.Printf("Last report time: %s\n", parseTime(summary.LastReportAt).String())
	fmt.Printf("Last interval time: %s\n", parseTime(summary.LastIntervalEndAt).String())
}

func parseTime(unixSeconds int64) time.Time {
	return time.Unix(unixSeconds, 0)
}
