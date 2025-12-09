package geo

import (
	"math"
	"time"

	"github.com/paulmach/orb"
)

// CreateStd4dVolContents computes and returns the start time, end time
//
//	and 2d polygon for a waypoint and time data.
func CreateStd4dVolContents(startTime time.Time, duration time.Duration, center orb.Point) (time.Time, time.Time, *orb.Polygon) {
	polygon := HexagonPlanar(center)
	endTime := startTime.Add(duration)
	return startTime, endTime, &polygon
}

func HexagonPlanar(center orb.Point) orb.Polygon {
	const radius = 0.001
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
