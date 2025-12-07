package manna_utm_client_cmds

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/manna_utm_client"
	"manna.aero/manna-utm-geojson-api/model/uspace"
)

var CreateOperationalIntent = &cobra.Command{
	Use:   "mutm-coi",
	Short: "Create an operational intent for <mission_id> with the data from <file> in manna-utm listening on <port>.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		if len(args) == 1 {
			port = 28080
			log.Infof("no port defined, configuring target USS port to default %d", port)
		}
		entityId, err := cmd.Flags().GetString("mission_id")
		if err != nil {
			return err
		}
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		baseUrl := fmt.Sprintf("http://localhost:%d/", port)
		mannaUtmClient, err := manna_utm_client.NewMannaUtmClient(baseUrl)
		if err != nil {
			log.Fatalf("unable to create USS mannaUtmClient: %v", err)
		}

		log.Debugf("attempting to fetch most recent telemetry message from USS server on port: %d", port)

		// read the operational intent from the provided file
		intent, err := uspace.LoadOperationalIntentFromFile(filePath)
		if err != nil {
			log.Fatalf("unable to load operational intent from file: %v", err)
		}

		err = mannaUtmClient.CreateOperationalIntent(cmd.Context(), entityId, intent)
		if err != nil {
			log.Fatalf("unable to fetch most recent telemetry message from USS server: %v", err)
		}

		data, err := json.MarshalIndent(intent, "", "    ")
		if err != nil {
			log.Errorf("unable to marshal returned operational intent details from USS to JSON: %v", err)
		}

		fmt.Println(string(data))

		return nil
	},
}
