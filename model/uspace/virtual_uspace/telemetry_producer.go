package virtual_uspace

import (
	"time"

	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/model/uspace"
)

func (oim *OperationalIntentManager) ProduceTelemetryMessagesToBus(bus *TelemetryBus, interval time.Duration) {
	for _, t := range oim.telemetry {
		log.Tracef("sending telemetry message to bus for operational intent.")
		bus.TelemetryEvents <- t
		time.Sleep(interval)
	}
}

type TelemetryBus struct {
	TelemetryEvents chan uspace.Telemetry
}

func (tb TelemetryBus) NewBus(bufferSize int) *TelemetryBus {
	return &TelemetryBus{
		TelemetryEvents: make(chan uspace.Telemetry, bufferSize),
	}
}
