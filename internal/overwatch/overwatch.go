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
	client    enphase.Client
	frequency int
	ticker    *time.Ticker
}

func NewOverwatch(client enphase.Client) Overwatch {
	frequency := 5
	watch := overwatch{
		client:    client,
		frequency: frequency,
		ticker:    nil,
	}

	return &watch
}

func (ow *overwatch) Start() {
	ow.ticker = time.NewTicker(time.Duration(ow.frequency) * time.Second)

	for range ow.ticker.C {
		ow.monitorEnergy()
		ow.ticker.Reset(time.Duration(ow.frequency) * time.Second)
	}
}

func (ow *overwatch) monitorEnergy() {
	consumption, _ := ow.client.GetConsumptionMeter()
	production, _ := ow.client.GetProductionMeter()

	generatingExcess := Evaluate(*consumption, *production)

	log.Printf("Excess power being generated: %+v", generatingExcess)
}
