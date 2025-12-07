package sequence

import (
	"math"
	"time"

	"github.com/paulmach/orb"
	"manna.aero/manna-utm-geojson-api/model/uspace"
)

// UspaceOperationalIntentFromConfig constructs the JSON body for a
// request to create an operational intent in manna-utm, given a config file.
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L73-L91
func UspaceOperationalIntentFromConfig(config *OperationalIntentConfig) *uspace.OperationalIntent {
	// load the config
	// create standard hexagonal 4d volumes for each waypoint
	curTime := time.Now()
	timeIncrement := config.Duration / time.Duration(len(config.WaypointCoordinates))
	altitude := 0.0

	// loop through all coordinates, creating waypoints and volumes for each
	var volumes []uspace.Volume4d
	var waypoints []uspace.Waypoint
	for _, coordinate := range config.WaypointCoordinates {
		// create the 4d volume for the coordinate
		p := orb.Point{
			coordinate[0],
			coordinate[1],
		}
		volumes = append(volumes, createStandard4dVolume(curTime, timeIncrement, p))
		curTime.Add(timeIncrement)

		// create the waypoint for the coordinate
		var wp uspace.Waypoint
		wp.Delta = float64(timeIncrement)
		wp.Latitude = coordinate[0]
		wp.Longitude = coordinate[1]
		wp.Time = curTime
		wp.Altitude = altitude
		waypoints = append(waypoints, wp)
	}

	// construct the final request body type
	return &uspace.OperationalIntent{
		Priority:      config.Priority,
		DepartureTime: curTime,
		Volumes:       volumes,
		Waypoints:     waypoints,
	}
}

func createStandard4dVolume(startTime time.Time, duration time.Duration, center orb.Point) uspace.Volume4d {
	return uspace.Volume4d{
		TimeStart: startTime,
		TimeEnd:   startTime.Add(duration),
		Polygon:   hexagonPlanar(center, 1000),
	}
}

func hexagonPlanar(center orb.Point, radius float64) orb.Polygon {
	const sides = 6
	ring := make(orb.Ring, 0, sides+1)

	for i := 0; i < sides; i++ {
		angle := (math.Pi / 3.0) * float64(i) // 60° steps
		x := center[0] + radius*math.Cos(angle)
		y := center[1] + radius*math.Sin(angle)
		ring = append(ring, orb.Point{x, y})
	}

	// Close the ring (GeoJSON polygon requirement)
	ring = append(ring, ring[0])

	return orb.Polygon{ring}
}
