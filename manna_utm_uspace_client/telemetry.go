package manna_utm_uspace_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"manna.aero/manna-utm-geojson-api/model/uspace"
)

// SendTelemetry interfaces with the manna-utm telemetry interface
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L165-L197
func (mutm *MannaUtmClient) SendTelemetry(ctx context.Context, message *uspace.Telemetry, missionId string, uavId int) error {

	requestUrl := fmt.Sprintf("%s/ussClient/v1/operational_intents/%s", mutm.baseUrl, missionId)

	messageContents, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewBuffer(messageContents))
	if err != nil {
		log.Errorf("An error occurred fetching operational intent operationalIntentDetails from USS: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", mutm.UserAgent)

	resp, err := mutm.c.Do(req)
	if err != nil {
		return MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return MannaUtmError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	resp, err = mutm.c.Do(req)
	if err != nil {
		return MannaUtmError{StatusCode: resp.StatusCode, Body: err.Error()}
	}

	return nil
}
