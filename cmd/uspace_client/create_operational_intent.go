package uspace_client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/model/uspace"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/uspace_client"
)

var CreateOperationalIntent = &cobra.Command{
	Use:     "us-create-operational-intent",
	Aliases: []string{"coi"},
	Short:   "Create an operational intent for <operation_name> in manna-utm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}
		oiName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		appConfig, err := config.LoadConfig("./config.yaml")
		if err != nil {
			return err
		}

		mannaUtmClient, err := uspace_client.NewMannaUtmClient("localhost", appConfig.MannaUtmPort, writeRequests)
		if err != nil {
			log.Fatalf("unable to create USS mannaUtmClient: %v", err)
		}

		log.Debugf("attempting to create operational intent via manna-utm U-Space interface on port: %d", appConfig.MannaUtmPort)

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

		err = mannaUtmClient.CreateOperationalIntent(cmd.Context(), oiCnf.UavId, oiCnf.MissionId.String(), oi)
		if err != nil {
			log.Fatalf("failed to create operational intent: %v", err.Error())
		}

		return nil
	},
}
