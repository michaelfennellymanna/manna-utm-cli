package virtual_uspace

import (
	"fmt"
	"sync"
	"time"

	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/pkg/geo"
)

type virtualSubOi struct {
	parentOi  *OperationalIntentManager
	startTime time.Time
	endTime   time.Time
	p1        orb.Point
	p2        orb.Point

	index int
}

func (oim *OperationalIntentManager) newVirtualSubOi(startTime time.Time, endTime time.Time, p1 orb.Point, p2 orb.Point, index int) virtualSubOi {
	return virtualSubOi{
		parentOi:  oim,
		startTime: startTime,
		endTime:   endTime,
		p1:        p1,
		p2:        p2,
		index:     index,
	}
}

func (vsoi *virtualSubOi) interpolateFeatures() error {
	if vsoi.parentOi.df < 0 {
		return fmt.Errorf("invalid argument. The detail factor is less than 0")
	}
	if vsoi.index < 0 || vsoi.index > len(vsoi.parentOi.waypoints) {
		return fmt.Errorf("invalid argument. Index of virutal subarray is out of bounds of parent size")
	}

	log.Tracef("interpolating features for virtualSubOi of index: %d", vsoi.index)

	// create the 4d volumes
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		vsoi.initVolume4ds()
	}()

	// create the telemetry
	wg.Add(1)
	go func() {
		defer wg.Done()
		vsoi.initTelemetry()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		vsoi.initWaypoints()
	}()

	wg.Wait()
	return nil
}

func (vsoi *virtualSubOi) initVolume4ds() {
	// create a 4d volume for the start point
	vol := uspace.Volume4d{
		TimeStart:     vsoi.startTime,
		TimeEnd:       vsoi.endTime,
		AltitudeLower: AltLower,
		AltitudeUpper: AltUpper,
		Polygon:       geo.HexagonPlanar(geo.Midpoint(vsoi.p1, vsoi.p2)),
	}
	vsoi.parentOi.volumeLock.Lock()
	vsoi.parentOi.volumes[vsoi.index] = vol
	vsoi.parentOi.volumeLock.Unlock()
}

func (vsoi *virtualSubOi) initWaypoints() {
	vsoi.parentOi.waypointLock.Lock()
	vsoi.parentOi.waypoints[vsoi.index] = uspace.Waypoint{
		Altitude: AltLower + ((AltUpper - AltLower) / 2),
	}
	vsoi.parentOi.waypointLock.Unlock()
}
