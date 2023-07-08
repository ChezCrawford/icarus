package overwatch

import (
	"icarus/internal/enphase"
	"icarus/internal/utils"
	"log"
	"math"
	"time"
)

type energyDelta struct {
	IntervalEnd time.Time
	DeltaWh     int64
}

func IsExcessProduction(consumptionMeter enphase.ConsumptionMeter, productionMeter enphase.ProductionMeter) bool {
	deltas := CalculateNets(consumptionMeter, productionMeter)

	return generatingExcess(deltas)
}

func CalculateNets(consumptionMeter enphase.ConsumptionMeter, productionMeter enphase.ProductionMeter) []energyDelta {
	consumptions := consumptionMeter.Intervals
	productions := productionMeter.Intervals

	consumption := consumptions[len(consumptions)-1]
	production := productions[len(productions)-1]

	if consumption.EndAt != production.EndAt {
		// Right now we are not very smrt.
		log.Panic("Meter data not taken at matching intervals.")
		return nil
	}

	intervals := len(productions)
	deltas := make([]energyDelta, intervals)

	for i := intervals - 1; i >= 0; i-- {
		consumption = consumptions[i]
		production = productions[i]

		netPower := production.WhDel - consumption.Enwh
		time := utils.ParseUnixSeconds(production.EndAt)
		delta := energyDelta{IntervalEnd: time, DeltaWh: netPower}
		deltas[i] = delta
	}

	return deltas
}

func generatingExcess(deltas []energyDelta) bool {
	intervals := len(deltas)
	lastFew := deltas[max(intervals-3, 0):intervals]

	logDeltas(lastFew)

	excess := true
	for _, d := range lastFew {
		if d.DeltaWh < 0 {
			excess = false
		}
	}

	return excess
}

func logDeltas(deltas []energyDelta) {
	for _, delta := range deltas {
		log.Printf("Delta: %+v", delta)
	}
}

func max(a int, b int) int {
	c := math.Max(float64(a), float64(b))
	return int(c)
}
