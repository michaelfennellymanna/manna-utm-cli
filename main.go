package main

import (
	"log"

	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/cmd"
	"manna.aero/manna-utm-geojson-api/cmd/manna_utm_client_cmds"
	"manna.aero/manna-utm-geojson-api/cmd/uss_client_cmds"
)

var (
	port     int
	fromFile string
	entityId string
)

var rootCmd = &cobra.Command{
	Use:   "manna-utm-cli",
	Short: "A command line tool for working with manna-utm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	cmd.RidDP.Flags().IntVarP(&port, "port", "p", 38080, "Listen port to bind the server to.")
	cmd.Data.Flags().StringVar(&fromFile, "file", "./.libconfig/sequence.yaml", "The path to the directory that you want to write the JSON contents of the simulation data to.")

	uss_client_cmds.UssClientFetchTelemetry.Flags().StringVar(&fromFile, "file", "", "The file that contains the JSON for the telemetry message required to send.")
	uss_client_cmds.GetOperationalIntentDetails.Flags().StringVar(&entityId, "entityId", "", "The entityId of the operational intent to fetch latest telemetry for.")

	manna_utm_client_cmds.Query4dVolume.Flags().StringVar(&fromFile, "file", "", "The file that contains the JSON for the 4d volume you want to send.")
	manna_utm_client_cmds.CreateOperationalIntent.Flags().StringVar(&fromFile, "file", "", "The file that contains the JSON for the operational intent you want to create.")
}

func main() {
	rootCmd.AddCommand(cmd.RidDP)
	rootCmd.AddCommand(cmd.Data)

	rootCmd.AddCommand(uss_client_cmds.UssClientFetchTelemetry)
	rootCmd.AddCommand(uss_client_cmds.GetOperationalIntentDetails)

	rootCmd.AddCommand(manna_utm_client_cmds.Query4dVolume)
	rootCmd.AddCommand(manna_utm_client_cmds.CreateOperationalIntent)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
