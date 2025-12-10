package riddp_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/pkg/config"
)

func GetServer(appConfig config.Config, writeRequests bool) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "manna-utm-cli geojson server is healthy.",
		})
	})

	router.GET("/features", func(c *gin.Context) {
		geoJson := appConfig.ToGeoJson()
		data, err := json.Marshal(geoJson)
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/json", []byte(err.Error()))
		}

		c.Data(http.StatusOK, "application/json", data)
	})

	router.GET("/features/events", func(c *gin.Context) {
		// CORS for SSE (adjust origin as needed)
		c.Writer.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://localhost:3001"))
		// If you use cookies/auth:
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"error": "streaming unsupported"})
			return
		}

		// initial write helps some setups
		_, _ = c.Writer.WriteString(": connected\n\n")
		flusher.Flush()

		// Optional: tell the client how long to wait before auto-reconnect (ms)
		_, _ = c.Writer.WriteString("retry: 3000\n\n")
		flusher.Flush()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		log.Infof("operning sse connection")

		operationalIntentGeoJsonFc := appConfig.ToGeoJson()
		operationalIntentGeoJsonFcBytes, err := json.Marshal(operationalIntentGeoJsonFc)
		if err != nil {
			log.Fatalf("error occurred marshalling config to GeoJson %v", err)
		}

		if writeRequests == true {
			outFileName := ".requests/riddp_init.geojson"
			err := os.WriteFile(outFileName, operationalIntentGeoJsonFcBytes, os.ModePerm)
			if err != nil {
				log.Errorf("unable to GeoJSON data to " + outFileName)
			}
		}

		// create a telemetry bus, that we can listen to for the configured flights
		telemetryBus := GetTelemetryBusFromConfig(appConfig, 1000)

		for {
			select {
			case <-ticker.C: // if we receive a ticker message in the channel.
				log.Debugf("responding to client with operationalIntentGeoJsonFcBytes")
				// Send event. Each SSE event is delimited by a blank line.
				// Note: "data:" supports multi-line; JSON is usually single-line.
				_, err := fmt.Fprintf(c.Writer, "event: features\ndata: %s\n\n", operationalIntentGeoJsonFcBytes)
				if err != nil {
					log.Errorf("error writing event features to client: %v", err)
					return
				}
				flusher.Flush()
				break
			case telemetryMessage, ok := <-telemetryBus.TelemetryEvents:
				if !ok {
					log.Warnf("telemetry channel closed")
					return
				}
				log.Tracef("telemetry message received from mission: %s", telemetryMessage.GetMissionId())

				// convert it to geojson.
				operationalIntentGeoJsonFc.Append(telemetryMessage.ToGeoJsonFeature())
				operationalIntentGeoJsonFcBytes, err = json.Marshal(operationalIntentGeoJsonFc)
				if err != nil {
					log.Fatalf("error occurred marshalling config to GeoJson %v", err)
				}

				_, err := fmt.Fprintf(c.Writer, "event: features\ndata: %s\n\n", operationalIntentGeoJsonFcBytes)
				if err != nil {
					log.Errorf("error writing event features to client: %v", err)
					return
				}
				flusher.Flush()
				// write it to the client channel.
				break
			}
		}
	})

	return router
}
