package manna_utm_client_cmds

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/config"
	"manna.aero/manna-utm-geojson-api/manna_utm_uspace_client"
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
		entityId, err := cmd.Flags().GetString("mission_id")
		if err != nil {
			return err
		}
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}
		oiName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		uavId, err := cmd.Flags().GetInt("uav_id")
		if err != nil {
			return err
		}

		mannaUtmClient, err := manna_utm_uspace_client.NewMannaUtmClient("localhost", port, writeRequests)
		if err != nil {
			log.Fatalf("unable to create USS mannaUtmClient: %v", err)
		}

		log.Debugf("attempting to create operational intent via manna-utm U-Space interface on port: %d", port)

		// load config
		appCnf, err := config.LoadConfig("./config.yaml")
		if err != nil {
			log.Fatalf("error occurred loading config: %v", err)
		}

		oiCnf, err := appCnf.GetOperationalIntentConfigByName(oiName)
		if err != nil {
			log.Fatalf("error occurred loading operational intent config: %v", err)
		}

		oi := uspace.OperationalIntentFromConfig(oiCnf)

		err = mannaUtmClient.CreateOperationalIntent(cmd.Context(), uavId, entityId, oi)
		if err != nil {
			log.Fatalf("failed to create operational intent: %v", err.Error())
		}

		return nil
	},
}
