package utm

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	log "github.com/sirupsen/logrus"
	"manna.aero/manna-utm-geojson-api/sequence"
)

type OperationalIntent struct {
	Reference OperationalIntentReference `json:"reference"`
	Details   OperationalIntentDetails   `json:"details"`
}

type OperationalIntentReference struct {
	ID uuid.UUID `json:"id"`
}

type OperationalIntentDetails struct {
	Volumes           []Volume4d `json:"volumes"`
	OffNominalVolumes []Volume4d `json:"off_nominal_volumes"`
	Priority          uint16     `json:"priority"`
}

// Volume4d is the equivalent of the manna-utm type UtmVolume4d
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/astm/dss/model/operationalintent/UtmVolume4D.java
type Volume4d struct {
	Volume    Volume3d  `json:"volume"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
}

// Volume3d is the equivalent of the manna-utm type UtmVolume3d
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/astm/dss/model/operationalintent/UtmVolume3D.java
type Volume3d struct {
	OutlinePolygon orb.Polygon `json:"outline_polygon"`
	AltitudeLower  float64     `json:"altitude_lower"`
	AltitudeUpper  float64     `json:"altitude_upper"`
}

func OperationalIntentFromConfig(oicnf *sequence.OperationalIntentConfig) *OperationalIntent {
	log.Tracef("constructing UTM operational intent for operational intent config: %s", oicnf.Name)
	// construct the volumes
	var vols []Volume4d
	startTime := time.Now()
	volDuration := oicnf.Duration / time.Duration(len(oicnf.WaypointCoordinates))
	for _, coord := range oicnf.WaypointCoordinates {
		vols = append(vols, *getVolume4dFromCoordinate(coord[0], coord[1], startTime, volDuration))
		startTime.Add(volDuration)
	}

	log.Tracef("returning UTM operational intent: %s", oicnf.Name)
	return &OperationalIntent{
		Reference: OperationalIntentReference{
			ID: oicnf.ID,
		},
		Details: OperationalIntentDetails{
			Volumes:  vols,
			Priority: oicnf.Priority,
		},
	}
}

func getVolume4dFromCoordinate(lat float64, lng float64, startTime time.Time, duration time.Duration) *Volume4d {
	polygon := sequence.HexagonPlanar(orb.Point{lng, lat})
	var vol3d Volume3d
	vol3d.OutlinePolygon = polygon
	vol3d.AltitudeLower = 0.0
	vol3d.AltitudeUpper = 0.0

	return &Volume4d{
		Volume:    vol3d,
		TimeStart: startTime,
		TimeEnd:   startTime.Add(duration),
	}
}
