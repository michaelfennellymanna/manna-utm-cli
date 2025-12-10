package uss_client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/uss_client"
)

var UssClientFetchTelemetry = &cobra.Command{
	Use:   "uss-tel",
	Short: "Use the USS client to fetch telemetry from the USS as the specified port.",
	Args:  cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		entityId, err := cmd.Flags().GetString("entityId")
		if err != nil {
			return err
		}

		appConfig, err := config.LoadConfig("./config.yaml")
		if err != nil {
			return err
		}

		client, err := uss_client.NewUssClient(fmt.Sprintf("http://localhost:%d", appConfig.MannaUtmPort))
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", appConfig.MannaUtmPort)
		err = client.GetLatestTelemetryForOperationalIntentByEntityId(cmd.Context(), entityId)
		if err != nil {
			log.Fatalf("unable to fetch most recent telemetry message from USS server: %v", err)
		}

		return nil
	},
}
