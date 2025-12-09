package manna_utm_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
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

func (mutm *MannaUtmClient) CreateOperationalIntent(ctx context.Context, uavId int, entityId string, intent *uspace.OperationalIntent) error {

	requestUrl, err := url.JoinPath(mutm.baseUrl.String(), path.Join("/operationalintent", strconv.Itoa(uavId), entityId))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, nil)
	if err != nil {
		log.Errorf("An error occurred attempting to create operational intent in manna-utm: %v", err)
	}

	logRequestContents(req)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", mutm.UserAgent)

	resp, err := mutm.c.Do(req)
	if err != nil {
		if resp == nil {
			return &MannaUtmError{StatusCode: 0, Body: err.Error()}
		}

		errJson, err := json.MarshalIndent(err.Error(), "", "   ")
		if err != nil {
			log.Errorf("error trying to parse error response from manna-utm: %v", err)
			return &MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
		}
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: string(errJson)}
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
	}

	var operationalIntentDetails utm.OperationalIntentDetails
	if err := json.Unmarshal(b, &operationalIntentDetails); err != nil {
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
	}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields() // optional; helps catch API changes
	if err := dec.Decode(resp.Body); err != nil {
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
	}
	return nil
}

func NewMannaUtmClient(host string, port int) (*MannaUtmClient, error) {
	baseUrlStr := fmt.Sprintf("http://%s:%d", host, port)
	baseUrl, err := url.Parse(baseUrlStr)
	if err != nil {
		return nil, err
	}

	return &MannaUtmClient{
		baseUrl: baseUrl,
		c: &http.Client{
			Timeout: 15 * time.Second,
		},
		UserAgent: "manna-utm-cli",
	}, nil
}

type MannaUtmError struct {
	StatusCode int
	Body       string
}

func (m MannaUtmError) Error() string {
	return fmt.Sprintf("manna-utm error: status=%d body=%q", m.StatusCode, m.Body)
}

func logRequestContents(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true) // true includes body (consumes it)
	if err != nil {
		log.Debugf("dump error: %v", err)
		return
	}
	log.Debugf("REQUEST:\n%s", dump)
}
