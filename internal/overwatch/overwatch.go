package overwatch

import (
	"context"
	"icarus/internal/enphase"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type Overwatch interface {
	Start()
}

type overwatch struct {
	client        enphase.Client
	pdClient      *pagerduty.Client
	pollFrequency time.Duration
	ticker        *time.Ticker
}

func NewOverwatch(client enphase.Client, pdClient *pagerduty.Client, pollFrequency time.Duration) Overwatch {
	watch := overwatch{
		client:        client,
		pdClient:      pdClient,
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

func CheckPagerDuty(pdClient *pagerduty.Client) error {
	var opts pagerduty.ListAddonOptions
	ctx := context.Background()
	addons, err := pdClient.ListAddonsWithContext(ctx, opts)

	log.Printf("Addons: %+v, Error: %+v", addons, err)

	if err != nil {
		return err
	}

	return nil
}
