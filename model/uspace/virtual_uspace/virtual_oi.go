package virtual_uspace

import (
	"sync"
	"time"

	"github.com/paulmach/orb/geojson"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/pkg/config"
)

// OperationalIntentManager is a structure that interpolates an
// operational intent (OI), given the waypoints and duration of the
// OI.
//
// # Defining an OI.
//
// The `uspace.OperationalIntent` structure contains certain temporal
// and atemporal properties:
//
// 1. Priority of the OI.
// 2. The departure time.
// 3. The 4d volumes.
// 4. The waypoints.
// 5. The telemetry messages.
//
// # API, and purposefully limited data before init.
//
// The user defines their operational intent, using a protocol
// (U-Space/UTM) agnostic configuration. This configuration
// is serialized on program initialization, into the `config.OperationalIntentConfig`
// structure.
//
// # Sub-parts and concurrent initialization.
//
// During initialization via the structures APIs, it breaks up the
// operational intent into parts, and constructs each subpart in parallel.
// Namely, each OperationalIntentManager is constructed of many
// `virtualSubOperationalIntent` instances. Each of these instances is
// the aforementioned 'subpart' of this OperationalIntentManager.
type OperationalIntentManager struct {
	// telemetry series
	df int

	telemetryLock sync.Mutex
	telemetry     []uspace.Telemetry
	waypointLock  sync.Mutex
	waypoints     []uspace.Waypoint
	volumeLock    sync.Mutex
	volumes       []uspace.Volume4d

	departureTime time.Time
}

// NewOperationalIntentManager constructs the JSON body for a
// request to create an operational intent in manna-utm, given a config file.
//
// See [UTMController]
//
// [UTMController]: https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L73-L91
func NewOperationalIntentManager(oiCnf *config.OperationalIntentConfig, df int) *OperationalIntentManager {
	nP := len(oiCnf.WaypointCoordinates) - 1

	voi := OperationalIntentManager{
		df:        df,
		telemetry: make([]uspace.Telemetry, nP*df),
		waypoints: make([]uspace.Waypoint, nP),
		volumes:   make([]uspace.Volume4d, nP),
	}

	var wg sync.WaitGroup
	curTime := time.Now()
	timeIncrement := oiCnf.Duration / time.Duration(len(oiCnf.WaypointCoordinates))
	for i := 0; i+1 <= nP; i++ {
		nextTime := curTime.Add(timeIncrement)
		curCoordinate := oiCnf.WaypointCoordinates[i]
		nextCoordinate := oiCnf.WaypointCoordinates[i+1]
		vsoi := voi.newVirtualSubOi(curTime, nextTime, curCoordinate, nextCoordinate, i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := vsoi.interpolateFeatures()
			if err != nil {
				log.Errorf("an error occurred interpolating features for subarray at index %d: %v", vsoi.index, err)
			}
		}()
	}

	wg.Wait()

	voi.departureTime = time.Now()

	return &voi
}

func (oim *OperationalIntentManager) getOi() uspace.OperationalIntent {
	return uspace.OperationalIntent{
		Priority:      0,
		DepartureTime: oim.departureTime,
		Volumes:       oim.volumes,
		Waypoints:     oim.waypoints,
	}
}

func (oim *OperationalIntentManager) GeoJson(includeVols bool, includeWaypoints bool, includeTelemetry bool) *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()

	if includeVols {
		for _, vol := range oim.volumes {
			fc.Append(vol.GeoJsonFeature())
		}
	}
	if includeWaypoints {
		for _, wp := range oim.waypoints {
			fc.Append(wp.GeoJsonFeature())
		}
	}
	if includeTelemetry {
		for _, t := range oim.telemetry {
			fc.Append(t.GeoJsonFeature())
		}
	}

	return fc
}
