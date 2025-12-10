package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/cmd"
	"manna.aero/manna.utm.cli/cmd/riddp"
	"manna.aero/manna.utm.cli/cmd/uspace_client"
	"manna.aero/manna.utm.cli/cmd/uss_client"
)

var (
	port                    int
	fromFile                string
	entityId                string
	oiName                  string
	logLevel                string = "info"
	writeRequestsToHttpFile bool   = false
	volName                 string
)

const ConfigPath = "./config.yaml"

var rootCmd = &cobra.Command{
	Use:   "manna-utm-cli",
	Short: "A command line tool for working with manna-utm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	riddp.RidDP.Flags().IntVarP(&port, "port", "p", 38080, "Listen port to bind the server to.")
	cmd.Data.Flags().StringVar(&fromFile, "file", ConfigPath, "The path to the directory that you want to write the JSON contents of the simulation data to.")

	uss_client.UssClientFetchTelemetry.Flags().StringVar(&fromFile, "file", "", "The file that contains the JSON for the telemetry message required to send.")
	uss_client.GetOperationalIntentDetails.Flags().StringVar(&entityId, "entityId", "", "The entityId of the operational intent to fetch latest telemetry for.")

	uspace_client.Query4dVolume.Flags().StringVarP(&volName, "name", "n", "", "The name of the 4d volume in config.yaml to query.")
	uspace_client.Query4dVolume.Flags().BoolVarP(&writeRequestsToHttpFile, "dump-requests", "d", false, "Specify true/false to enable/disable writing requests to http files.")

	uspace_client.CreateOperationalIntent.Flags().StringVarP(&oiName, "name", "n", "", "The name of the operational intent that you want to create.")
	uspace_client.CreateOperationalIntent.Flags().BoolVarP(&writeRequestsToHttpFile, "dump-requests", "d", false, "Specify true/false to enable/disable writing requests to http files.")

	uspace_client.EndOperationalIntent.Flags().BoolVarP(&writeRequestsToHttpFile, "dump-requests", "d", false, "Specify true/false to enable/disable writing requests to http files.")
	uspace_client.EndOperationalIntent.Flags().StringVarP(&oiName, "name", "n", "", "the name of the operational intent that you want to delete.")

	riddp.RidDP.Flags().BoolVarP(&writeRequestsToHttpFile, "dump-requests", "d", false, "Specify true/false to enable/disable writing requests to http files.")
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
	rootCmd.AddCommand(riddp.RidDP)
	rootCmd.AddCommand(cmd.Data)

	rootCmd.AddCommand(uss_client.UssClientFetchTelemetry)
	rootCmd.AddCommand(uss_client.GetOperationalIntentDetails)

	rootCmd.AddCommand(uspace_client.Query4dVolume)
	rootCmd.AddCommand(uspace_client.CreateOperationalIntent)
	rootCmd.AddCommand(uspace_client.EndOperationalIntent)
	rootCmd.AddCommand(uspace_client.CancelOperationalIntent)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
