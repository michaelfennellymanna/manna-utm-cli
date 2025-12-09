package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/cmd"
	"manna.aero/manna-utm-geojson-api/cmd/manna_utm_client_cmds"
	"manna.aero/manna-utm-geojson-api/cmd/uss_client_cmds"
)

var (
	port                    int
	fromFile                string
	entityId                string
	missionId               string
	oiName                  string
	uavId                   int
	logLevel                string = "info"
	writeRequestsToHttpFile bool   = false
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
	manna_utm_client_cmds.Query4dVolume.Flags().BoolVar(&writeRequestsToHttpFile, "w", false, "Specify true/false to enable/disable writing requests to http files.")

	manna_utm_client_cmds.CreateOperationalIntent.Flags().StringVarP(&oiName, "name", "n", "", "The name of the operational intent that you want to create.")
	manna_utm_client_cmds.CreateOperationalIntent.Flags().IntVar(&port, "port", 28082, "The port that manna-utm is listening on.")
	manna_utm_client_cmds.CreateOperationalIntent.Flags().StringVarP(&missionId, "mission_id", "m", "", "The ID of the mission that you want to create.")
	manna_utm_client_cmds.CreateOperationalIntent.Flags().IntVarP(&uavId, "uav_id", "u", 1, "The ID of the UAV flying the payload in the operational intent.")
	manna_utm_client_cmds.CreateOperationalIntent.Flags().BoolVarP(&writeRequestsToHttpFile, "dump-requests", "d", false, "Specify true/false to enable/disable writing requests to http files.")
}

func configureLogging(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Warnf("Unknown log level %s, falling back to info", level)
	}
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "The log level that you want to run your command with.")
	cobra.OnInitialize(func() { configureLogging(logLevel) })
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
