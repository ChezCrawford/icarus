package main

import (
	"icarus/internal/enphase"
	"icarus/internal/overwatch"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Print("Let's become enlightened...")

	config := overwatch.LoadConfig()

	client := enphase.NewClient(config.EnphaseAccessToken, config.EnphaseApiKey, config.EnphaseSystemId)
	alerter := overwatch.NewPagerDutyIncidentAlerter(config.PdServiceId, config.PdUserEmail, config.PdApiKey)

	// Use these to fake stuff.
	// client := enphase.NewFileClient()
	// alerter := overwatch.NewLogAlerter()

	watch := overwatch.NewOverwatch(client, config.EnphasePollFrequency, alerter)
	watch.Start()

	log.Print("Done")
}
