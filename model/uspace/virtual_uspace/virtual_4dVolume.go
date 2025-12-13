package virtual_uspace

import (
	"time"

	"github.com/paulmach/orb"
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/pkg/geo"
)

func create4dVol(startTime time.Time, duration time.Duration, point orb.Point) uspace.Volume4d {
	start, end, polygon := geo.CreateStd4dVolContents(startTime, duration, point)

	return uspace.Volume4d{
		TimeStart:     start,
		TimeEnd:       end,
		AltitudeLower: AltLower,
		AltitudeUpper: AltUpper,
		Polygon:       *polygon,
		Wsg84:         0,
	}
}
