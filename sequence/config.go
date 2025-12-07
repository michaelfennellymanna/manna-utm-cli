package sequence

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Name                     string
	OperationalIntentConfigs []OperationalIntentConfig `yaml:"operational_intent_configs"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	var cfgName string
	cfgName = strings.TrimSuffix(path, ".yaml")
	cfgName = strings.TrimPrefix(cfgName, "./.libconfig/")
	cfg.Name = cfgName

	return &cfg, nil
}

func (c Config) ToGeoJson() *geojson.FeatureCollection {
	var featureCollection geojson.FeatureCollection
	for _, intent := range c.OperationalIntentConfigs {
		for _, feature := range *intent.geoJsonFeatureSlice() {
			featureCollection.Append(&feature)
		}
	}
	return &featureCollection
}

type OperationalIntentConfig struct {
	Name                string        `yaml:"name"`
	OwnerName           string        `yaml:"owner_name"`
	OwnerBaseURL        string        `yaml:"owner_baseurl"`
	Priority            uint16        `yaml:"priority"`
	ID                  uuid.UUID     `yaml:"id"`
	Duration            time.Duration `yaml:"duration"`
	WaypointCoordinates [][2]float64  `yaml:"waypoint_coordinates"`
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

func (oic OperationalIntentConfig) geoJsonFeatureSlice() *[]geojson.Feature {
	// Create all the 4d Volumes
	var fc []geojson.Feature
	startTime := time.Now()
	duration := oic.Duration / time.Duration(len(oic.WaypointCoordinates))
	for _, coordinate := range oic.WaypointCoordinates {
		startTime, endTime, polygon := CreateStd4dVolContents(startTime, duration, orb.Point{coordinate[1], coordinate[0]})
		// create a feature from the polygon
		f := geojson.NewFeature(polygon)
		// add metadata to the polygon, annotating start & end times
		f.Properties = map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
		}
		fc = append(fc, *f)
	}
	return &fc
}

// CreateStd4dVolContents computes and returns the start time, end time
//
//	and 2d polygon for a waypoint and time data.
func CreateStd4dVolContents(startTime time.Time, duration time.Duration, center orb.Point) (time.Time, time.Time, *orb.Polygon) {
	polygon := HexagonPlanar(center)
	endTime := startTime.Add(duration)
	return startTime, endTime, &polygon
}
