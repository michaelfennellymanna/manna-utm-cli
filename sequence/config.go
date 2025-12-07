package sequence

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb/geojson"
	"gopkg.in/yaml.v3"
)

type Simulation struct {
	startTime time.Time
	duration  time.Duration
}

type Config struct {
	OperationalIntents []OperationalIntentConfig `yaml:"operational_intents"`
}

func (c Config) ToGeoJson() []*geojson.FeatureCollection {
	var allFeatureCollections []*geojson.FeatureCollection
	for _, intent := range c.OperationalIntents {
		allFeatureCollections = append(allFeatureCollections, intent.geoJsonFc())
	}
	return allFeatureCollections
}

func LoadConfig(path string) (string, *Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return "", nil, fmt.Errorf("parse yaml: %w", err)
	}

	var cfgName string
	cfgName = strings.TrimSuffix(path, ".yaml")
	cfgName = strings.TrimPrefix(cfgName, "./.libconfig/")

	return cfgName, &cfg, nil
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

func (oic OperationalIntentConfig) geoJsonFc() *geojson.FeatureCollection {
	// Create all the 4d Volumes
	operationalIntent := UspaceOperationalIntentFromConfig(&oic)
	var fc geojson.FeatureCollection
	for _, v := range operationalIntent.Volumes {
		// create a feature from the polygon
		f := geojson.NewFeature(v.Polygon)
		// add metadata to the polygon, annotating start & end times
		f.Properties = map[string]interface{}{
			"start_time": v.TimeStart,
			"end_time":   v.TimeEnd,
		}
		fc.Append(f)
	}
	return &fc
}
