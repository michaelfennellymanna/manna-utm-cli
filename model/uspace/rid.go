package uspace

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// Telemetry is equivalent to https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceTelemetry.java
type Telemetry struct {
	Altitude      float64 `json:"altitude"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lng"`
	Heading       float64 `json:"heading"`
	Speed         float64 `json:"speed"`
	VerticalSpeed float64 `json:"vertical_speed"`
	TimeMeasured  int64   `json:"time_measured"`
	Mode          string  `json:"mode"`
	Armed         bool    `json:"armed"`
}

func (t Telemetry) GeoJsonFeature() *geojson.Feature {
	f := geojson.NewFeature(orb.Point{t.Latitude, t.Longitude})
	f.Properties["time_measured"] = t.TimeMeasured
	f.Properties["altitude"] = t.Altitude
	return f
}
