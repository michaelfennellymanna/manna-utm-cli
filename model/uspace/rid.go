package uspace

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/geo"
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
	missionId     uuid.UUID
}

func GetUspaceTelemetryFromConfig(cnf config.OperationalIntentConfig) []Telemetry {
	coords := geo.InterpolateCoordinateSeries(cnf.WaypointCoordinates)
	telemetrySeries := make([]Telemetry, 0, len(coords))
	for i, c := range coords {
		telemetryMessage := Telemetry{
			Altitude:      200,
			Latitude:      c.Lat(),
			Longitude:     c.Lon(),
			Heading:       0,
			Speed:         100,
			VerticalSpeed: 0,
			TimeMeasured:  time.Now().UnixMilli(),
			Mode:          "null-mode",
			Armed:         true,
			missionId:     cnf.MissionId,
		}
		telemetrySeries[i] = telemetryMessage
	}

	return telemetrySeries
}

func (t Telemetry) ToGeoJsonFeature() *geojson.Feature {
	f := geojson.NewFeature(orb.Point{t.Longitude, t.Latitude})

	f.Properties = map[string]interface{}{
		"flight_id":     t.missionId,
		"altitude":      t.Altitude,
		"heading":       t.Heading,
		"speed":         t.Speed,
		"time_measured": t.TimeMeasured,
	}

	return f
}

func (t Telemetry) GetMissionId() uuid.UUID {
	return t.missionId
}
