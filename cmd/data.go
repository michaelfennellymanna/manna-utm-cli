package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/sequence"
)

var Data = &cobra.Command{
	Use:   "data",
	Short: "Generate data for the configured simulations.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fromFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		configName, c, err := sequence.LoadConfig(fromFile)
		if err != nil {
			return err
		}
		operationalIntent := c.ToGeoJson()
		contents, err := json.MarshalIndent(operationalIntent, "", "  ")
		if err != nil {
			log.Fatalf("error occurred marshalling sequence to GeoJson %v", err)
		}

		err = os.MkdirAll("./.libconfig/personal/geojson", os.ModePerm)
		if err != nil {
			log.Fatalf("error occurred creating directory for configuration GeoJson %v", err)
		}

		outFileName := fmt.Sprintf("./.libconfig/personal/geojson/%s.geojson", configName)
		err = os.WriteFile(outFileName, contents, 0644)
		if err != nil {
			log.Fatalf("error occurred writing GeoJson to file: %v", err)
		}

		log.Printf("successfully wrote contents of geojson to file")

		return nil
	},
}
