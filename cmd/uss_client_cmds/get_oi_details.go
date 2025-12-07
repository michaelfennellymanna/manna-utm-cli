package uss_client_cmds

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/uss_client"
)

var GetOperationalIntentDetails = &cobra.Command{
	Use:   "uss-oid",
	Short: "Get operational intent details for <entity_id> from the USS listening on <port>.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}

		entityId, err := cmd.Flags().GetString("entity_id")
		if err != nil {
			return err
		}
		if len(args) == 1 {
			port = 28080
			log.Infof("no port defined, configuring target USS port to default %d", port)
		}

		baseUrl := fmt.Sprintf("http://localhost:%d/", port)
		client, err := uss_client.NewUssClient(baseUrl)
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", port)
		details, err := client.GetOperationalIntentDetailsByEntityId(cmd.Context(), entityId)
		if err != nil {
			log.Fatalf("unable to fetch most recent telemetry message from USS server: %v", err)
		}

		data, err := json.MarshalIndent(details, "", "    ")
		if err != nil {
			log.Errorf("unable to marshal returned operational intent details from USS to JSON: %v", err)
		}

		fmt.Println(string(data))

		return nil
	},
}
