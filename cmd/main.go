package main

import (
	"context"
	"icarus/internal/enphase"
	"icarus/internal/overwatch"
	"log"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Print("Let's become enlightened...")

	config := overwatch.LoadConfig()

	// client := enphase.NewClient(config.EnphaseAccessToken, config.EnphaseApiKey, config.EnphaseSystemId)
	client := enphase.NewFileClient()

	watch := overwatch.NewOverwatch(client)
	watch.Start()

	pdClient := pagerduty.NewClient(config.PdApiKey)

	var opts pagerduty.GetCurrentUserOptions
	ctx := context.Background()
	user, err := pdClient.GetCurrentUserWithContext(ctx, opts)

	log.Printf("User: %+v, Error: %+v", user, err)
}
