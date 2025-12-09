package uss_client_cmds

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/uss_client"
)

var UssClientFetchTelemetry = &cobra.Command{
	Use:   "uss-tel",
	Short: "Use the USS client to fetch telemetry from the USS as the specified port.",
	Args:  cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}

		entityId, err := cmd.Flags().GetString("entity_id")
		if err != nil {
			return err
		}
		log.SetLevel(log.InfoLevel)

		baseUrl := fmt.Sprintf("http://localhost:%d/", port)
		client, err := uss_client.NewUssClient(baseUrl)
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", port)
		err = client.GetLatestTelemetryForOperationalIntentByEntityId(cmd.Context(), entityId)
		if err != nil {
			log.Fatalf("unable to fetch most recent telemetry message from USS server: %v", err)
		}

		return nil
	},
}
