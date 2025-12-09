package manna_utm_uspace_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"manna.aero/manna-utm-geojson-api/model/uspace"
	"manna.aero/manna-utm-geojson-api/model/utm"
)

type MannaUtmClient struct {
	baseUrl       *url.URL
	c             *http.Client
	UserAgent     string
	writeRequests bool
}

func NewMannaUtmClient(host string, port int, writeRequests bool) (*MannaUtmClient, error) {
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
		UserAgent:     "manna-utm-cli",
		writeRequests: writeRequests,
	}, nil
}

func (mutm *MannaUtmClient) Query4dVolume(context context.Context, volName string) ([]utm.OperationalIntentDetails, *MannaUtmError) {
	panic("implement me")
	return nil, nil
}

type CreateOperationalIntentRequest struct {
	time      time.Time
	oi        *uspace.OperationalIntent
	uav       int
	missionId string
}

// CreateOperationalIntent interfaces with the manna-utm U-Space interface to create an operational intent.
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L56-L91
func (mutm *MannaUtmClient) CreateOperationalIntent(ctx context.Context, uavId int, entityId string, intent *uspace.OperationalIntent) error {

	requestUrl, err := url.JoinPath(mutm.baseUrl.String(), path.Join("/operationalintent", strconv.Itoa(uavId), entityId))
	if err != nil {
		return err
	}

	reader, err := intent.ToReader()
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read intent body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Errorf("An error occurred attempting to create operational intent in manna-utm: %v", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", mutm.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logReq, _ := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewReader(bodyBytes))
		logReq.Header = req.Header.Clone()
		mutm.logRequestContents(*logReq, CreateOperationalIntentRequest{
			time:      time.Now(),
			oi:        intent,
			uav:       uavId,
			missionId: entityId,
		})
	}()

	errChannel := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := mutm.c.Do(req)
		if err != nil {
			if resp == nil {
				errChannel <- &MannaUtmError{StatusCode: 0, Body: err.Error()}
				return
			}

			errJson, err := json.MarshalIndent(err.Error(), "", "   ")
			if err != nil {
				errChannel <- fmt.Errorf("error trying to parse error response from manna-utm: %v", err)
				return
			}
			errChannel <- &MannaUtmError{StatusCode: resp.StatusCode, Body: string(errJson)}
			return
		}
		defer resp.Body.Close()

		// Handle non-2xx responses with useful errors
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// read a limited amount so you don’t blow memory on huge error bodies
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
			errChannel <- &MannaUtmError{StatusCode: resp.StatusCode, Body: string(b)}
			return
		} else {
			log.Infof("successfully created operational intent in manna-utm")
		}

		return
	}()

	wg.Wait()
	select {
	case err := <-errChannel:
		return err
	default:
		return nil
	}
}

type MannaUtmError struct {
	StatusCode int
	Body       string
}

func (m MannaUtmError) Error() string {
	return fmt.Sprintf("manna-utm error: status=%d body=%q", m.StatusCode, m.Body)
}

func (mutm *MannaUtmClient) logRequestContents(r http.Request, request CreateOperationalIntentRequest) {
	reqDirName := "./.requests"
	// ensure outDir exists
	err := os.MkdirAll(reqDirName, os.ModePerm)
	if err != nil {
		log.Errorf("an error occurred creating requests output directory: %v", err)
	}

	dump, err := httputil.DumpRequest(&r, true)
	if err != nil {
		log.Debugf("dump error: %v", err)
		return
	}
	log.Debugf("REQUEST:\n%s", dump)

	if mutm.writeRequests {

		// write the operational intent by mission id and time
		outFileName := fmt.Sprintf("%v-COI-%v.http", request.time.Format(time.RFC3339), request.missionId)
		log.Debugf("writing request to http file: %s", outFileName)
		err = os.WriteFile(path.Join(reqDirName, outFileName), dump, 0644)
		if err != nil {
			log.Errorf("error occurred writing operational intent to file: %v", err)
			return
		}
	}
}
