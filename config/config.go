package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"gopkg.in/yaml.v3"
	"manna.aero/manna-utm-geojson-api/geo"
)

type Config struct {
	Name                     string                    `yaml:"name"`
	MannaUtmPort             int                       `yaml:"manna_utm_port"`
	OperationalIntentConfigs []OperationalIntentConfig `yaml:"operational_intent_configs"`
	FourDVolumes             []Volume4dConfig          `yaml:"4d_volumes"`
}

func (appCnf *Config) GetOperationalIntentConfigByName(name string) (*OperationalIntentConfig, error) {
	for _, oiCnf := range appCnf.OperationalIntentConfigs {
		if oiCnf.Name == name {
			return &oiCnf, nil
		}
	}

	return nil, fmt.Errorf("no operational intent is configured by the name: %s", name)
}

func (appCnf *Config) Get4dVolumeConfigByName(name string) (*Volume4dConfig, error) {
	for _, vCnf := range appCnf.FourDVolumes {
		if vCnf.Name == name {
			return &vCnf, nil
		}
	}

	return nil, fmt.Errorf("no operational intent is configured by the name: %s", name)
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

func (appCnf *Config) ToGeoJson() *geojson.FeatureCollection {
	var featureCollection geojson.FeatureCollection
	for _, intent := range appCnf.OperationalIntentConfigs {
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
	MissionId           uuid.UUID     `yaml:"mission_id"`
	UavId               int           `yaml:"uav_id"`
	Duration            time.Duration `yaml:"duration"`
	WaypointCoordinates [][2]float64  `yaml:"waypoint_coordinates"`
}

type Volume4dConfig struct {
	Name          string        `yaml:"name"`
	Duration      time.Duration `yaml:"duration"`
	AltLower      float64       `yaml:"alt_lower"`
	AltUpper      float64       `yaml:"alt_upper"`
	PolygonCoords [][2]float64  `yaml:"polygon_coords"`
}

func (oic OperationalIntentConfig) geoJsonFeatureSlice() *[]geojson.Feature {
	// Create all the 4d Volumes
	var fc []geojson.Feature
	startTime := time.Now()
	duration := oic.Duration / time.Duration(len(oic.WaypointCoordinates))
	for _, coordinate := range oic.WaypointCoordinates {
		startTime, endTime, polygon := geo.CreateStd4dVolContents(startTime, duration, orb.Point{coordinate[1], coordinate[0]})
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
