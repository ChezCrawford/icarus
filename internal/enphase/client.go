package enphase

import (
	"encoding/json"
	"errors"
	"icarus/internal/utils"
	"log"
	"net/http"
	"time"
)

const ApiUrl string = "https://api.enphaseenergy.com/api/v4"

type SystemSummary struct {
	SystemId          int64  `json:"system_id,omitempty"`
	CurrentPower      int64  `json:"current_power,omitempty"`
	EnergyLifetime    int64  `json:"energy_lifetime,omitempty"`
	EnergyToday       int64  `json:"energy_today,omitempty"`
	LastIntervalEndAt int64  `json:"last_interval_end_at,omitempty"`
	LastReportAt      int64  `json:"last_report_at,omitempty"`
	Modules           int64  `json:"modules,omitempty"`
	OperationalAt     int64  `json:"operational_at,omitempty"`
	SizeW             int64  `json:"size_w,omitempty"`
	Source            string `json:"source,omitempty"`
	Status            string `json:"status,omitempty"`
	SummaryDate       string `json:"summary_date,omitempty"`
}

type ConsumptionMeter struct {
	SystemID     int64  `json:"system_id"`
	Granularity  string `json:"granularity"`
	TotalDevices int64  `json:"total_devices"`
	StartAt      int64  `json:"start_at"`
	EndAt        int64  `json:"end_at"`
	Items        string `json:"items"`
	Intervals    []struct {
		EndAt            int64 `json:"end_at"`
		DevicesReporting int64 `json:"devices_reporting"`
		Enwh             int64 `json:"enwh"`
	} `json:"intervals"`
	Meta struct {
		Status        string `json:"status"`
		LastReportAt  int64  `json:"last_report_at"`
		LastEnergyAt  int64  `json:"last_energy_at"`
		OperationalAt int64  `json:"operational_at"`
	} `json:"meta"`
}

type ProductionMeter struct {
	SystemID     int64  `json:"system_id"`
	Granularity  string `json:"granularity"`
	TotalDevices int64  `json:"total_devices"`
	StartAt      int64  `json:"start_at"`
	EndAt        int64  `json:"end_at"`
	Items        string `json:"items"`
	Intervals    []struct {
		EndAt            int64 `json:"end_at"`
		DevicesReporting int64 `json:"devices_reporting"`
		WhDel            int64 `json:"wh_del"`
	} `json:"intervals"`
}

type Client interface {
	GetConsumptionMeter() (*ConsumptionMeter, error)
	GetProductionMeter() (*ProductionMeter, error)
	GetSystemSummary() (*SystemSummary, error)
}

type EnphaseClient struct {
	accessToken string
	apiKey      string
	client      *http.Client
	systemId    string
}

func NewClient(accessToken string, apiKey string, systemId string) *EnphaseClient {
	enphaseClient := &EnphaseClient{
		accessToken: accessToken,
		apiKey:      apiKey,
		client:      &http.Client{},
		systemId:    systemId,
	}

	system, err := enphaseClient.GetSystemSummary()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found Enphase system with id: %+v, status: %+v, and last report: %+v",
		system.SystemId, system.Status, time.Unix(system.LastReportAt, 0))
	log.Printf("System last reported: %+v, last interval ended: %+v",
		utils.ParseUnixSeconds(system.LastReportAt), utils.ParseUnixSeconds(system.LastIntervalEndAt))

	return enphaseClient
}

func (c EnphaseClient) GetConsumptionMeter() (*ConsumptionMeter, error) {
	path := "/systems/" + c.systemId + "/telemetry/consumption_meter"

	path = ApiUrl + path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return nil, err
	}

	var response ConsumptionMeter
	err = c.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c EnphaseClient) GetProductionMeter() (*ProductionMeter, error) {
	path := "/systems/" + c.systemId + "/telemetry/production_meter"

	path = ApiUrl + path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return nil, err
	}

	var response ProductionMeter
	err = c.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}

func (c EnphaseClient) GetSystemSummary() (*SystemSummary, error) {
	path := "/systems/" + c.systemId + "/summary"

	path = ApiUrl + path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return nil, err
	}

	var response SystemSummary
	err = c.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c EnphaseClient) doRequest(req *http.Request, response interface{}) error {
	authorization := "Bearer " + c.accessToken

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error performing HTTP request: %v", err)
		return err
	} else if res.StatusCode >= 400 {
		log.Printf("Error response %v", res.StatusCode)
		return errors.New("HTTP Error")
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
