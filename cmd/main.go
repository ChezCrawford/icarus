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

	// client := enphase.NewClient(config.EnphaseAccessToken, config.EnphaseApiKey, config.EnphaseSystemId)
	client := enphase.NewFileClient()
	alerter := overwatch.NewPagerDutyIncidentAlerter(config.PdServiceId, config.PdUserEmail, config.PdApiKey)

	watch := overwatch.NewOverwatch(client, config.EnphasePollFrequency, alerter)
	watch.Start()

	log.Print("Done")
}
