package manna_utm_client

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"manna.aero/manna-utm-geojson-api/model/uspace"
	"manna.aero/manna-utm-geojson-api/model/utm"
)

type MannaUtmClient struct {
	baseUrl   *url.URL
	c         *http.Client
	UserAgent string
}

func (mutm *MannaUtmClient) Query4dVolume(context context.Context, fromFile string) ([]utm.OperationalIntentDetails, *MannaUtmError) {
	panic("implement me")
	return nil, nil
}

func (mutm *MannaUtmClient) CreateOperationalIntent(ctx context.Context, id string, intent *uspace.OperationalIntent) error {
	panic("implement me")
}

func NewMannaUtmClient(baseUrl string) (*MannaUtmClient, error) {
	u, err := url.Parse(strings.TrimRight(baseUrl, "/") + "/")
	if err != nil {
		return nil, err
	}
	return &MannaUtmClient{
		baseUrl: u,
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
