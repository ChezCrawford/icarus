package enphase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

type FileClient struct {
	currentConsumptionInterval int
	currentProductionInterval  int
	consumption                ConsumptionMeter
	production                 ProductionMeter
}

func NewFileClient() *FileClient {
	// Interesting intervals:
	// 30 is where production begins
	// 97 intervals in total
	startInterval := 135

	consumptionFile := "test/data/2022-09-15-consumption-meter.json"
	productionFile := "test/data/2022-09-15-production-meter.json"

	var consumption ConsumptionMeter
	if err := loadFromFile(consumptionFile, &consumption); err != nil {
		panic(fmt.Sprintf("Unable to load data from file, error=%+v", err))
	}

	var production ProductionMeter
	if err := loadFromFile(productionFile, &production); err != nil {
		panic(fmt.Sprintf("Unable to load data from file, error=%+v", err))
	}

	return &FileClient{
		currentConsumptionInterval: startInterval,
		currentProductionInterval:  startInterval,
		consumption:                consumption,
		production:                 production,
	}
}

func (c *FileClient) GetConsumptionMeter() (*ConsumptionMeter, error) {
	c.currentConsumptionInterval += 1
	endInterval := min(c.currentConsumptionInterval, len(c.consumption.Intervals))

	consumptionCopy := c.consumption
	consumptionCopy.Intervals = c.consumption.Intervals[0:endInterval]

	return &consumptionCopy, nil
}

func (c *FileClient) GetProductionMeter() (*ProductionMeter, error) {
	c.currentProductionInterval += 1
	endInterval := min(c.currentProductionInterval, len(c.production.Intervals))
	productionCopy := c.production
	productionCopy.Intervals = c.production.Intervals[0:endInterval]

	return &productionCopy, nil
}

func (c *FileClient) GetSystemSummary() (*SystemSummary, error) {
	return nil, errors.New("not implemented")
}

func loadFromFile(fileName string, response interface{}) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &response); err != nil {
		return err
	}

	return nil
}

func min(a int, b int) int {
	c := math.Min(float64(a), float64(b))
	return int(c)
}
