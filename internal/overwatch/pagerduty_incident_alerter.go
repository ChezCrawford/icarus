package overwatch

import (
	"context"
	"log"

	"github.com/PagerDuty/go-pagerduty"
)

type PagerDutyGateway interface {
	ExcessPowerAlert() error
}

type PagerDutyIncidentAlerter struct {
	pdServiceId string
	pdUserEmail string
	pdClient    *pagerduty.Client
}

func NewPagerDutyIncidentAlerter(pdServiceId string, pdUserEmail string, pdApiKey string) *PagerDutyIncidentAlerter {
	pdClient := pagerduty.NewClient(pdApiKey)

	if err := verifyConnectivity(pdClient, pdServiceId); err != nil {
		log.Fatal(err)
	}

	return &PagerDutyIncidentAlerter{pdServiceId: pdServiceId, pdUserEmail: pdUserEmail, pdClient: pdClient}
}

func verifyConnectivity(pdClient *pagerduty.Client, serviceId string) error {
	ctx := context.Background()

	service, err := pdClient.GetServiceWithContext(ctx, serviceId, nil)

	log.Printf("Service: %+v, Error: %+v", service, err)
	if err != nil {
		return err
	}

	return nil
}

func (a PagerDutyIncidentAlerter) ExcessPowerAlert() error {
	ctx := context.Background()

	incidentRequest := &pagerduty.CreateIncidentOptions{
		Title: "We are losing power!!!",
		Service: &pagerduty.APIReference{
			ID:   a.pdServiceId,
			Type: "service_reference",
		},
	}

	incident, err := a.pdClient.CreateIncidentWithContext(ctx, a.pdUserEmail, incidentRequest)

	if err != nil {
		return err
	}

	log.Printf("Created incident with id: %+v, number: %+v, key: %+v", incident.ID, incident.IncidentNumber, incident.IncidentKey)

	return nil
}
