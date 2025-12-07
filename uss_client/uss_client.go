package uss_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"manna.aero/manna-utm-geojson-api/model/utm"
)

type UssClient struct {
	ussBaseUrl *url.URL
	c          *http.Client
	UserAgent  string
}

func NewUssClient(base string) (*UssClient, error) {
	u, err := url.Parse(strings.TrimRight(base, "/") + "/")
	if err != nil {
		return nil, err
	}
	return &UssClient{
		ussBaseUrl: u,
		c: &http.Client{
			Timeout: 15 * time.Second,
		},
		UserAgent: "manna-utm-cli",
	}, nil
}

type UssClientError struct {
	StatusCode int
	Body       string
}

func (e *UssClientError) Error() string {
	return fmt.Sprintf("uss client error: status=%d body=%q", e.StatusCode, e.Body)
}

// GetOperationalIntentDetailsByEntityId from the USS.
func (ussClient *UssClient) GetOperationalIntentDetailsByEntityId(ctx context.Context, entityId string) (*utm.OperationalIntentDetails, *UssClientError) {

	requestUrl := fmt.Sprintf("%s/uss/v1/operational_intents/%s", ussClient.ussBaseUrl.String(), entityId)

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		log.Errorf("An error occurred fetching operational intent operationalIntentDetails from USS: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", ussClient.UserAgent)

	resp, err := ussClient.c.Do(req)
	if err != nil {
		return nil, &UssClientError{StatusCode: resp.StatusCode, Body: err.Error()}
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return nil, &UssClientError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &UssClientError{StatusCode: resp.StatusCode, Body: err.Error()}
	}

	var operationalIntentDetails utm.OperationalIntentDetails
	if err := json.Unmarshal(b, &operationalIntentDetails); err != nil {
		return nil, &UssClientError{StatusCode: resp.StatusCode, Body: err.Error()}
	}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields() // optional; helps catch API changes
	if err := dec.Decode(resp.Body); err != nil {
		return nil, &UssClientError{StatusCode: resp.StatusCode, Body: err.Error()}
	}
	return &operationalIntentDetails, nil
}

// GetLatestTelemetryForOperationalIntentByEntityId gets the latest telemetry
// message for the specified operational intent, from the USS.
func (ussClient *UssClient) GetLatestTelemetryForOperationalIntentByEntityId(ctx context.Context, entityId string) error {
	requestUrl := fmt.Sprintf("%s/ussClient/v1/operational_intents/%s", ussClient.ussBaseUrl.String(), entityId)

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, nil)
	if err != nil {
		log.Errorf("An error occurred fetching operational intent details from USS: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", ussClient.UserAgent)

	resp, err := ussClient.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return &UssClientError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields() // optional; helps catch API changes
	if err := dec.Decode(resp.Body); err != nil {
		return err
	}
	return nil
}
