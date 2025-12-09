package manna_utm_client_cmds

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/manna_utm_uspace_client"
)

var Query4dVolume = &cobra.Command{
	Use:   "mutm-qv",
	Short: "Use the manna-utm client to query the a UTM airspace with provided 4d volume data.",
	Args:  cobra.MaximumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		if len(args) == 1 {
			port = 28080
			log.Infof("no port defined, configuring target USS port to default %d", port)
		}

		fromFilePath, err := cmd.Flags().GetString("file")
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}

		client, err := manna_utm_uspace_client.NewMannaUtmClient("localhost", port, writeRequests)
		if err != nil {
			log.Fatalf("unable to create USS client: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", port)
		allOperationsIn4dVolume, err := client.Query4dVolume(cmd.Context(), fromFilePath)
		if err != nil {
			log.Fatalf("unable to fetch most recent telemetry message from USS server: %v", err)
		}

		data, err := json.MarshalIndent(allOperationsIn4dVolume, "", "    ")
		if err != nil {
			log.Errorf("unable to marshal returned operational intent details from USS to JSON: %v", err)
		}

		fmt.Println(string(data))

		return nil
	},
}
