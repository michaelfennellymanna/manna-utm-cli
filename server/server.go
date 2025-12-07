package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna-utm-geojson-api/sequence"
)

func GetServer(port int) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "asdb.fi mock server is healthy.",
		})
	})

	router.GET("/features", func(c *gin.Context) {
		_, config, err := sequence.LoadConfig("./.libconfig/sequence.yaml")
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/json", []byte(err.Error()))
		}
		geoJson := config.ToGeoJson()
		data, err := json.Marshal(geoJson)
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/json", []byte(err.Error()))
		}

		c.Data(http.StatusOK, "application/json", data)
	})

	router.GET("/features/events", func(c *gin.Context) {
		// CORS for SSE (adjust origin as needed)
		c.Writer.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://localhost:%d", port))
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

		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		log.Infof("operning sse connection")
		for {
			select {
			case <-ticker.C:
				payload := "payload"

				log.Infof("responding to client with payload")
				// Send event. Each SSE event is delimited by a blank line.
				// Note: "data:" supports multi-line; JSON is usually single-line.
				fmt.Fprintf(c.Writer, "event: features\ndata: %s\n\n", payload)
				flusher.Flush()
			}
		}
	})

	return router
}
