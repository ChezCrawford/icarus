package overwatch

import (
	"icarus/internal/enphase"
	"log"
	"math"
	"time"
)

type EnergyDelta struct {
	IntervalEnd time.Time
	DeltaWh     int64
}

func Evaluate(consumptionMeter enphase.ConsumptionMeter, productionMeter enphase.ProductionMeter) bool {
	deltas := CalculateNets(consumptionMeter, productionMeter)
	logDeltas(deltas)

	return generatingExcess(deltas)
}

func CalculateNets(consumptionMeter enphase.ConsumptionMeter, productionMeter enphase.ProductionMeter) []EnergyDelta {
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
	deltas := make([]EnergyDelta, intervals)

	for i := intervals - 1; i >= 0; i-- {
		consumption = consumptions[i]
		production = productions[i]

		netPower := production.WhDel - consumption.Enwh
		time := ParseTime(production.EndAt)
		delta := EnergyDelta{IntervalEnd: time, DeltaWh: netPower}
		deltas[i] = delta
	}

	return deltas
}

func logDeltas(deltas []EnergyDelta) {
	for _, delta := range deltas {
		log.Printf("Delta: %+v", delta)
	}

	for _, delta := range positiveDeltas(deltas) {
		log.Printf("Positive Delta: %+v", delta)
	}
}

func positiveDeltas(deltas []EnergyDelta) (positives []EnergyDelta) {
	for _, delta := range deltas {
		if delta.DeltaWh > 0 {
			positives = append(positives, delta)
		}
	}

	return positives
}

func generatingExcess(deltas []EnergyDelta) bool {
	intervals := len(deltas)
	lastFew := deltas[max(intervals-3, 0):intervals]

	excess := true
	for _, d := range lastFew {
		if d.DeltaWh < 0 {
			excess = false
		}
	}

	return excess
}

func ParseTime(unixSeconds int64) time.Time {
	return time.Unix(unixSeconds, 0)
}

func max(a int, b int) int {
	c := math.Max(float64(a), float64(b))
	return int(c)
}
