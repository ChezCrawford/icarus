package overwatch

import (
	"icarus/internal/enphase"
	"log"
	"time"
)

type Overwatch interface {
	Start()
}

type Alerter interface {
	ExcessPowerAlert() error
}

type overwatch struct {
	alerter       Alerter
	client        enphase.Client
	pollFrequency time.Duration
	ticker        *time.Ticker
}

func NewOverwatch(client enphase.Client, pollFrequency time.Duration, alerter Alerter) Overwatch {
	watch := overwatch{
		alerter:       alerter,
		client:        client,
		pollFrequency: pollFrequency,
		ticker:        nil,
	}

	return &watch
}

func (ow *overwatch) Start() {
	ow.ticker = time.NewTicker(ow.pollFrequency)

	for range ow.ticker.C {
		ow.monitorEnergy()
		ow.ticker.Reset(ow.pollFrequency)
	}
}

func (ow *overwatch) monitorEnergy() {
	consumption, _ := ow.client.GetConsumptionMeter()
	production, _ := ow.client.GetProductionMeter()

	generatingExcess := Evaluate(*consumption, *production)
	log.Printf("Excess power being generated: %+v", generatingExcess)

	if generatingExcess {
		if err := ow.alerter.ExcessPowerAlert(); err != nil {
			log.Fatal(err)
		}
	}
}
