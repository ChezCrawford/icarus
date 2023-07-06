package main

import (
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
	pdClient := pagerduty.NewClient(config.PdApiKey)

	if err := overwatch.CheckPagerDuty(pdClient); err != nil {
		log.Fatal(err)
	}

	watch := overwatch.NewOverwatch(client, pdClient, config.EnphasePollFrequency)
	watch.Start()

	log.Print("Done")
}
