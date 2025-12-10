package geo

import (
	"math"
	"sync"
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

func PolygonFromCoords(coords [][2]float64) orb.Polygon {
	ring := make(orb.Ring, 0, len(coords)+1)
	for _, c := range coords {
		ring = append(ring, orb.Point{c[0], c[1]})
	}

	// Ensure closed ring (first point == last point)
	if len(ring) > 0 && ring[0] != ring[len(ring)-1] {
		ring = append(ring, ring[0])
	}

	// Polygon = []Ring (first ring is outer, others are holes)
	return orb.Polygon{ring}
}

type pointSeries struct {
	points                      [][2]float64
	seriesLock                  sync.Mutex
	interpolatedTelemetrySeries []orb.Point
}

func InterpolateCoordinateSeries(initialSeries [][2]float64) []orb.Point {
	fp := pointSeries{
		points: initialSeries,
	}

	finalArrayCapacity := (len(fp.points) * 3) + 1
	fp.interpolatedTelemetrySeries = make([]orb.Point, 0, finalArrayCapacity)

	var wg sync.WaitGroup
	for i := 0; i < len(fp.points)-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fp.interpolateWindow(i)
		}()
	}
	wg.Wait()

	// set the last waypoint message to the existing final point
	fp.interpolatedTelemetrySeries[len(fp.interpolatedTelemetrySeries)] = fp.points[len(fp.points)]

	return fp.interpolatedTelemetrySeries
}

func (fp *pointSeries) interpolateWindow(index int) {
	start := fp.points[index]
	end := fp.points[index+1]
	oneThirdLat := (end[0] - start[0]) / 3
	oneThirdLong := (end[1] - start[1]) / 3

	p1 := orb.Point{
		start[0] + oneThirdLat,
		start[1] + oneThirdLong,
	}

	p2 := orb.Point{
		start[0] + (2 * oneThirdLat),
		start[1] + (2 + oneThirdLong),
	}

	fp.seriesLock.Lock()
	fp.interpolatedTelemetrySeries[index*3] = start
	fp.interpolatedTelemetrySeries[(index*3)+1] = p1
	fp.interpolatedTelemetrySeries[(index*3)+2] = p2
	fp.seriesLock.Unlock()
}
