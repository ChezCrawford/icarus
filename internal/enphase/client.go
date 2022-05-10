package enphase

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

type Client struct {
	accessToken string
	apiKey      string
	client      *http.Client
}

func NewClient(accessToken string, apiKey string) *Client {
	return &Client{
		accessToken: accessToken,
		apiKey:      apiKey,
		client:      &http.Client{},
	}
}

func (c Client) GetSystemSummary(systemId string) (*SystemSummary, error) {
	path := "/systems/" + systemId + "/summary"

	path = ApiUrl + path
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return nil, err
	}

	authorization := "Bearer " + c.accessToken

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error performing HTTP request: %v", err)
		return nil, err
	} else if res.StatusCode >= 400 {
		log.Printf("Error response %v", res.StatusCode)
		return nil, errors.New("HTTP Error")
	}

	defer res.Body.Close()

	var summary SystemSummary
	json.NewDecoder(res.Body).Decode(&summary)

	return &summary, nil
}

func GetData() int {
	return 1
}
