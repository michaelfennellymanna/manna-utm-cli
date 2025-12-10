package uspace_client

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
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/model/utm"
	"manna.aero/manna.utm.cli/pkg/config"
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

// Query4dVolume uses the manna-utm U-Space interface to query a given 4d volume.
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L141-L163
func (mutm *MannaUtmClient) Query4dVolume(ctx context.Context, volName string) ([]utm.OperationalIntentDetails, error) {
	c, err := config.LoadConfig("./config.yaml")
	if err != nil {
		return nil, err
	}

	volCnf, err := c.Get4dVolumeConfigByName(volName)
	if err != nil {
		return nil, err
	}

	vol := uspace.GetVolume4dFromConfig(*volCnf)

	reader, err := vol.ToReader()
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read intent body: %w", err)
	}

	requestUrl := fmt.Sprintf("http://localhost:%d/operationalintent/query", c.MannaUtmPort)
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Errorf("An error occurred attempting to create operational intent in manna-utm: %v", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", mutm.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	mutm.logRequestContents(*req, IdAndTimeRecord{startTime: time.Now(), missionId: volName})

	resp, err := mutm.c.Do(req)
	if err != nil {
		if resp == nil {
			return nil, err
		}

		errJson, err := json.Marshal(err.Error())
		if err != nil {
			return nil, fmt.Errorf("error trying to parse error response from manna-utm: %v", err)
		}
		return nil, &MannaUtmError{StatusCode: resp.StatusCode, Body: string(errJson)}
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return nil, &MannaUtmError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	return nil, nil
}

type CreateOperationalIntentRequest struct {
	time      time.Time
	oi        *uspace.OperationalIntent
	uav       int
	missionId string
}

func (coir CreateOperationalIntentRequest) getTime() time.Time {
	return coir.time
}

func (coir CreateOperationalIntentRequest) getMissionId() string {
	return coir.missionId
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

			errJson, err := json.Marshal(err.Error())
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

type IdAndTimeRecord struct {
	missionId string
	startTime time.Time
}

func (eoir IdAndTimeRecord) getTime() time.Time {
	return eoir.startTime
}

func (eoir IdAndTimeRecord) getMissionId() string {
	return eoir.missionId
}

// EndOperationalIntent interfaces with the manna-utm U-Space interface to end the operational intent associated with <missionId>
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L117-L139
//
// NOTE that the UTMController class specifies the argument to this method as the 'operationId'.
// In InterUSS, there are only entityIds, and there can be 1 entity for 1 operation.
// When this is routed to InterUSS, manna-utm will remove all versions of the entity (specified by OVN) from the DSS.
func (mutm *MannaUtmClient) EndOperationalIntent(ctx context.Context, missionId string) error {
	requestUrl, err := url.JoinPath(mutm.baseUrl.String(), path.Join("/operationalintent", missionId, "end"))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", requestUrl, nil)
	if err != nil {
		log.Errorf("An error occurred attempting to create operational intent in manna-utm: %v", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", mutm.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	mutm.logRequestContents(*req, IdAndTimeRecord{startTime: time.Now(), missionId: missionId})

	resp, err := mutm.c.Do(req)
	if err != nil {
		if resp == nil {
			return &MannaUtmError{StatusCode: 0, Body: err.Error()}
		}

		errJson, err := json.Marshal(err.Error())
		if err != nil {
			return fmt.Errorf("error trying to parse error response from manna-utm: %v", err)
		}
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: string(errJson)}
	}
	defer resp.Body.Close()

	// Handle non-2xx responses with useful errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read a limited amount so you don’t blow memory on huge error bodies
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return &MannaUtmError{StatusCode: resp.StatusCode, Body: string(b)}
	} else {
		log.Infof("successfully ended operational intent in manna-utm")
	}

	return nil
}

type MannaUtmError struct {
	StatusCode int
	Body       string
}

func (m MannaUtmError) Error() string {
	return fmt.Sprintf("manna-utm error: status=%d body=%q", m.StatusCode, m.Body)
}

type HasTimeAndMissionId interface {
	getTime() time.Time
	getMissionId() string
}

func (mutm *MannaUtmClient) logRequestContents(r http.Request, request HasTimeAndMissionId) {
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

	if mutm.writeRequests == true {
		// write the operational intent by mission id and time
		outFileName := fmt.Sprintf("%v-%v.http", request.getTime().UnixMilli(), request.getMissionId())
		log.Debugf("writing request to http file: %s", outFileName)
		err = os.WriteFile(path.Join(reqDirName, outFileName), dump, 0644)
		if err != nil {
			log.Errorf("error occurred writing operational intent to file: %v", err)
			return
		}
	}
}
