package uspace

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/paulmach/orb"
)

// Volume4d is equivalent to the type MannaUspaceVolume4d in manna-utm.
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceVolume4d.java
type Volume4d struct {
	TimeStart     time.Time   `json:"time_start"`
	TimeEnd       time.Time   `json:"time_end"`
	AltitudeLower float64     `json:"altitude_lower"`
	AltitudeUpper float64     `json:"altitude_upper"`
	Polygon       orb.Polygon `json:"polygon"`
	Wsg84         float64     `json:"wsg_84"`
}
type volume4dJSON struct {
	TimeStart     int64    `json:"time_start"`
	TimeEnd       int64    `json:"time_end"`
	AltitudeLower float32  `json:"altitude_lower"`
	AltitudeUpper float32  `json:"altitude_upper"`
	Polygon       []Vertex `json:"polygon"`
	Wsg84         float64  `json:"wsg_84"`
}

type Vertex struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (vol Volume4d) MarshalJSON() ([]byte, error) {
	var polygon []Vertex

	ring := vol.Polygon[0]
	firstAndLastVertexOverlap := len(ring) > 1 && ring[0] == ring[len(ring)-1]

	if firstAndLastVertexOverlap {
		for i := 0; i < len(ring)-1; i++ {
			polygon = append(polygon, Vertex{
				Latitude:  float32(ring[i][1]),
				Longitude: float32(ring[i][0]),
			})
		}
	}

	return json.Marshal(volume4dJSON{
		TimeStart:     vol.TimeStart.UnixMilli(),
		TimeEnd:       vol.TimeEnd.UnixMilli(),
		Polygon:       polygon,
		AltitudeLower: float32(vol.AltitudeLower),
		AltitudeUpper: float32(vol.AltitudeUpper),
		Wsg84:         vol.Wsg84,
	})
}

// Waypoint is equivalent to https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceWaypoint.java
type Waypoint struct {
	Altitude  float64   `json:"altitude"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Delta     float64   `json:"delta"`
	Time      time.Time `json:"time"`
}

type waypointJSON struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Time      int64   `json:"time"` // Unix millis
	Delta     float64 `json:"delta"`
}

func (wp Waypoint) MarshalJSON() ([]byte, error) {
	out := waypointJSON{
		Altitude:  wp.Altitude,
		Latitude:  wp.Latitude,
		Longitude: wp.Longitude,
		Delta:     wp.Delta,
		Time:      wp.Time.UnixMilli(),
	}

	// If you don't want pretty JSON, use json.Marshal(out) instead.
	return json.Marshal(out)
}

// OperationalIntent is equivalent to the type MannUspaceOperationalIntent
// in manna-utm.
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceOperationalIntent.java
type OperationalIntent struct {
	Priority      uint16     `json:"priority"`
	DepartureTime time.Time  `json:"departure_time"`
	Volumes       []Volume4d `json:"volumes"`
	Waypoints     []Waypoint `json:"waypoints"`
}

type operationalIntentJSON struct {
	Priority      uint16     `json:"priority"`
	DepartureTime int64      `json:"departure_time"`
	Volumes       []Volume4d `json:"volumes"`
	Waypoints     []Waypoint `json:"waypoints"`
}

func (oi OperationalIntent) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&operationalIntentJSON{
		Priority:      oi.Priority,
		DepartureTime: oi.DepartureTime.UnixMilli(),
		Volumes:       oi.Volumes,
		Waypoints:     oi.Waypoints,
	}, "", "  ")
}

func (oi OperationalIntent) ToReader() (*bytes.Reader, error) {
	jsonBytes, err := oi.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}
