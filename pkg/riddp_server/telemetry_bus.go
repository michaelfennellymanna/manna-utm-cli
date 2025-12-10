package riddp_server

import (
	"time"

	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/pkg/config"
)

type TelemetryBus struct {
	TelemetryEvents chan uspace.Telemetry
}

// GetTelemetryBusFromConfig constructs a TelemetryBus object, and a producer
// that periodically starts producing telemetry messages to the bus over the specified
// produceInterval.
func GetTelemetryBusFromConfig(appConfig config.Config, produceInterval time.Duration) *TelemetryBus {
	// get the interpolated U-Space series from the appConfig.

	bufferSize := 100

	// construct the bus
	bus := TelemetryBus{
		TelemetryEvents: make(chan uspace.Telemetry, bufferSize),
	}

	log.Debugf("starting telemetry producers for all operational intents in the config")

	for _, oiCnf := range appConfig.OperationalIntentConfigs {
		// create a telemetry producer
		go startTelemetryProducerForOi(produceInterval, &bus, oiCnf)
	}

	return &bus
}

func startTelemetryProducerForOi(producerInterval time.Duration, bus *TelemetryBus, intentConfig config.OperationalIntentConfig) {
	// get telemetry from the config
	allTelemetryMessages := uspace.GetUspaceTelemetryFromConfig(intentConfig)

	for _, msg := range allTelemetryMessages {
		// write the message to the channel
		log.Tracef("writing telemetry message (pos=%f,%f) to bus", msg.Latitude, msg.Longitude)
		bus.TelemetryEvents <- msg

		time.Sleep(producerInterval)
	}
}
