package uspace_client

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/uspace_client"
)

var Query4dVolume = &cobra.Command{
	Use:     "us-query-volume",
	Aliases: []string{"qv"},
	Short:   "Use the manna-utm client to query the a UTM airspace with a configured 4d volume.",
	Args:    cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		volName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}

		c, err := config.LoadConfig("./config.yaml")
		if err != nil {
			return err
		}

		client, err := uspace_client.NewMannaUtmClient("localhost", c.MannaUtmPort, writeRequests)
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", c.MannaUtmPort)
		allOperationsIn4dVolume, err := client.Query4dVolume(cmd.Context(), volName)
		if err != nil {
			log.Fatalf("an error occurred while querying 4d volume: %v", err)
		}

		data, err := json.MarshalIndent(allOperationsIn4dVolume, "", "    ")
		if err != nil {
			log.Errorf("unable to marshal returned operational intent details from USS to JSON: %v", err)
		}

		fmt.Println(string(data))

		return nil
	},
}
