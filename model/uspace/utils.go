package uspace

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/paulmach/orb"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/geo"
)

// OperationalIntentFromConfig constructs the JSON body for a
// request to create an operational intent in manna-utm, given a config file.
//
// see https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/controller/UTMController.java#L73-L91
func OperationalIntentFromConfig(oiCnf *config.OperationalIntentConfig) *OperationalIntent {
	// load the oiCnf
	// create standard hexagonal 4d volumes for each waypoint
	curTime := time.Now()
	timeIncrement := oiCnf.Duration / time.Duration(len(oiCnf.WaypointCoordinates))
	altitude := 0.0

	// loop through all coordinates, creating waypoints and volumes for each
	var volumes []Volume4d
	var waypoints []Waypoint
	for _, coordinate := range oiCnf.WaypointCoordinates {
		// create the 4d volume for the coordinate
		startTime, endTime, polygon := geo.CreateStd4dVolContents(curTime, timeIncrement, orb.Point{coordinate[1], coordinate[0]})
		volumes = append(volumes, *create4dVolFromStdContents(startTime, endTime, polygon))
		curTime = endTime

		// create the waypoint for the coordinate
		wp := Waypoint{
			Delta:     float64(timeIncrement),
			Latitude:  coordinate[0],
			Longitude: coordinate[1],
			Time:      curTime,
			Altitude:  altitude,
		}
		waypoints = append(waypoints, wp)
	}

	// construct the final request body type
	return &OperationalIntent{
		Priority:      oiCnf.Priority,
		DepartureTime: curTime,
		Volumes:       volumes,
		Waypoints:     waypoints,
	}
}

func create4dVolFromStdContents(startTime time.Time, endTime time.Time, polyGon *orb.Polygon) *Volume4d {
	return &Volume4d{
		TimeStart:     startTime,
		TimeEnd:       endTime,
		Polygon:       *polyGon,
		AltitudeLower: 0.0,
		AltitudeUpper: 0.0,
	}
}

func LoadOperationalIntentFromFile(path string) (*OperationalIntent, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var intent OperationalIntent
	if err := json.Unmarshal(b, &intent); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	var cfgName string
	cfgName = strings.TrimSuffix(path, ".yaml")
	cfgName = strings.TrimPrefix(cfgName, "./.libconfig/")

	return &intent, nil
}
