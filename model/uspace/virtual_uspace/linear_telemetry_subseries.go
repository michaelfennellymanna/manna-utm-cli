package virtual_uspace

import (
	"time"

	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/model/uspace"
)

func (vsoi *virtualSubOi) initTelemetry() {
	numMessages := float64(vsoi.parentOi.df)

	t1 := createTelemetryMessage(vsoi.p1, vsoi.startTime)
	t2 := createTelemetryMessage(vsoi.p2, vsoi.endTime)
	tTransformation := uspace.Telemetry{
		Altitude:      (t2.Altitude - t1.Altitude) / numMessages,
		Latitude:      (t2.Latitude - t1.Latitude) / numMessages,
		Longitude:     (t2.Longitude - t1.Longitude) / numMessages,
		Heading:       (t2.Heading - t1.Heading) / numMessages,
		Speed:         (t2.Speed - t1.Speed) / numMessages,
		VerticalSpeed: (t2.VerticalSpeed - t1.VerticalSpeed) / numMessages,
		TimeMeasured:  (t2.TimeMeasured - t1.TimeMeasured) / int64(numMessages),
		Mode:          "null-mode",
		Armed:         true,
	}

	for i := 0; i < vsoi.parentOi.df; i++ {
		thisMessage := transformTelemetry(t1, tTransformation, float64(i))
		indexOfMessageInParent := (vsoi.index * vsoi.parentOi.df) + i
		vsoi.parentOi.telemetryLock.Lock()
		// where to add in the parent index is defined by the subarray's index
		// in the parent and the index of the message within the subarray
		vsoi.parentOi.telemetry[indexOfMessageInParent] = thisMessage
		log.Tracef("adding telemetry message for index %d, and offset %d", vsoi.index, i)
		vsoi.parentOi.telemetryLock.Unlock()
	}
}

func transformTelemetry(input uspace.Telemetry, transformation uspace.Telemetry, factor float64) uspace.Telemetry {
	out := input

	out.Altitude = input.Altitude + (transformation.Altitude * factor)
	out.Latitude = input.Latitude + (transformation.Latitude * factor)
	out.Longitude = input.Longitude + (transformation.Longitude * factor)
	out.Heading = input.Heading + (transformation.Heading * factor)
	out.Speed = input.Speed + (transformation.Speed * factor)
	out.VerticalSpeed = input.VerticalSpeed + (transformation.VerticalSpeed * factor)
	out.TimeMeasured = input.TimeMeasured + (transformation.TimeMeasured * int64(factor))

	return out
}

func createTelemetryMessage(twoDimensionalPoint orb.Point, timeMeasured time.Time) uspace.Telemetry {
	altitudeDelta := AltUpper - AltLower
	altitude := float64(AltLower + (altitudeDelta / 2))
	return uspace.Telemetry{
		Altitude:      altitude,
		Latitude:      twoDimensionalPoint.Lat(),
		Longitude:     twoDimensionalPoint.Lon(),
		Heading:       0,
		Speed:         100,
		VerticalSpeed: 0,
		TimeMeasured:  timeMeasured.UnixMilli(),
		Mode:          "null-mode",
		Armed:         true,
	}
}
