package uspace

import (
	"time"

	"github.com/paulmach/orb"
)

// Volume4d is equivalent to the type MannaUspaceVolume4d in manna-utm.
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceVolume4d.java
type Volume4d struct {
	TimeStart     time.Time
	TimeEnd       time.Time
	Polygon       orb.Polygon
	altitudeLower float64
	altitudeUpper float64
}

// Waypoint is equivalent to https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceWaypoint.java
type Waypoint struct {
	Altitude  float64   `json:"altitude"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Delta     float64   `json:"delta"`
	Time      time.Time `json:"time"`
}

// OperationalIntent is equivalent to the type MannUspaceOperationalIntent
// in manna-utm.
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceOperationalIntent.java
type OperationalIntent struct {
	Priority      uint16     `json:"priority"`
	DepartureTime time.Time  `json:"departure_time"`
	Waypoints     []Waypoint `json:"waypoints"`
	Volumes       []Volume4d `json:"volumes"`
}
