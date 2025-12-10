package uss_client

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/uss_client"
)

var GetOperationalIntentDetails = &cobra.Command{
	Use:   "uss-oid",
	Short: "Get operational intent details for <entity_id> from the USS listening on <port>.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		entityId, err := cmd.Flags().GetString("entity_id")
		if err != nil {
			return err
		}

		c, err := config.LoadConfig("./config.yaml")
		if err != nil {
			return err
		}

		baseUrl := fmt.Sprintf("http://localhost:%d/", c.MannaUtmPort)
		client, err := uss_client.NewUssClient(baseUrl)
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch operational intent details from USS server on port: %d", c.MannaUtmPort)
		details, err := client.GetOperationalIntentDetailsByEntityId(cmd.Context(), c.MannaUtmPort, entityId)
		if err != nil {
			log.Fatalf("unable to fetch operational intent details from USS server: %v", err)
		}

		data, err := json.MarshalIndent(details, "", "    ")
		if err != nil {
			log.Errorf("unable to marshal returned operational intent details from USS to JSON: %v", err)
		}

		fmt.Println(string(data))

		return nil
	},
}
