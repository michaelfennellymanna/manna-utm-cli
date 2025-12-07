package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/model/uspace"
	"manna.aero/manna-utm-geojson-api/model/utm"
	"manna.aero/manna-utm-geojson-api/sequence"
)

const PERSONAL_LIB_PATH = "./.libconfig/personal"

var Data = &cobra.Command{
	Use:   "data",
	Short: "Generate data for the configured simulations in <file>.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		config, err := sequence.LoadConfig(fromFile)
		if err != nil {
			log.Fatalf("error occurred loading config: %v", err)
		}

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			// create the GeoJson data
			err = ConvertToGeoJsonAndWriteToFile(config)
			if err != nil {
				log.Fatalf("error occurred writing GeoJson to file: %v", err)
			}
		}()

		err = os.MkdirAll(fmt.Sprintf("%s/utm", PERSONAL_LIB_PATH), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(fmt.Sprintf("%s/uspace", PERSONAL_LIB_PATH), os.ModePerm)
		if err != nil {
			panic(err)
		}

		for _, oiConfig := range config.OperationalIntentConfigs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// create the UTM data
				oi := utm.OperationalIntentFromConfig(&oiConfig)
				data, err := json.MarshalIndent(oi, "", "  ")
				if err != nil {
					log.Errorf("error occurred marshalling operational intent (name=%s) to JSON: %v", oiConfig.Name, err.Error())
					return
				}

				// create the file if it doesn't already exist
				err = os.WriteFile(fmt.Sprintf("%s/utm/%s.json", PERSONAL_LIB_PATH, oiConfig.Name), data, 0644)
				if err != nil {
					log.Errorf("error occurred writing operational intent (name=%s) to file: %v", oiConfig.Name, err.Error())
					return
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				// create the U-Space telemetry data
				oi := uspace.OperationalIntentFromConfig(&oiConfig)
				data, err := json.MarshalIndent(oi, "", "  ")
				if err != nil {
					log.Errorf("error occurred marshalling operational intent (name=%s) to JSON: %v", oiConfig.Name, err.Error())
					return
				}

				// create the file if it doesn't already exist
				err = os.WriteFile(fmt.Sprintf("%s/uspace/%s.json", PERSONAL_LIB_PATH, oiConfig.Name), data, 0644)
				if err != nil {
					log.Errorf("error occurred writing operational intent (name=%s) to file: %v", oiConfig.Name, err.Error())
					return
				}
			}()
		}

		wg.Wait()
		return nil
	},
}

func ConvertToGeoJsonAndWriteToFile(c *sequence.Config) error {
	geoJsonDir := fmt.Sprintf("%s/geojson", PERSONAL_LIB_PATH)
	err := os.MkdirAll(geoJsonDir, os.ModePerm)

	if err != nil {
		return fmt.Errorf("error occurred creating directory for configuration GeoJson %w", err)
	}

	configGeoJson := c.ToGeoJson()
	contents, err := json.MarshalIndent(configGeoJson, "", "  ")
	if err != nil {
		log.Fatalf("error occurred marshalling sequence to GeoJson %v", err)
	}

	err = os.MkdirAll("./.libconfig/personal/geojson", os.ModePerm)
	if err != nil {
		log.Fatalf("error occurred creating directory for configuration GeoJson %v", err)
	}

	outFileName := fmt.Sprintf("./.libconfig/personal/geojson/%s.geojson", c.Name)
	err = os.WriteFile(outFileName, contents, 0644)
	if err != nil {
		log.Fatalf("error occurred writing GeoJson to file: %v", err)
	}

	return nil
}
