package overwatch

import (
	"icarus/internal/enphase"
	"log"
	"time"
)

type Overwatch interface {
	Start()
}

type overwatch struct {
	client        enphase.Client
	pollFrequency time.Duration
	ticker        *time.Ticker
}

func NewOverwatch(client enphase.Client, pollFrequency time.Duration) Overwatch {
	watch := overwatch{
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
}
