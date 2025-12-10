package uspace_client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/uspace_client"
)

var EndOperationalIntent = &cobra.Command{
	Use:     "us-end-operational-intent",
	Aliases: []string{"eoi"},
	Short:   "End an operational intent with the name <name> in config.yaml.",
	RunE: func(cmd *cobra.Command, args []string) error {
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}
		oiName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		// load config
		appCnf, err := config.LoadConfig("./config.yaml")
		if err != nil {
			log.Fatalf("error occurred loading config: %v", err)
		}
		mannaUtmClient, err := uspace_client.NewMannaUtmClient("localhost", appCnf.MannaUtmPort, writeRequests)
		if err != nil {
			log.Fatalf("unable to create USS mannaUtmClient: %v", err)
		}

		log.Debugf("attempting to create operational intent via manna-utm U-Space interface on port: %d", appCnf.MannaUtmPort)

		oiCnf, err := appCnf.GetOperationalIntentConfigByName(oiName)
		if err != nil {
			log.Fatalf("error occurred getting operational intent by the name %s from config: %v", oiName, err)
		}

		err = mannaUtmClient.EndOperationalIntent(cmd.Context(), oiCnf.MissionId.String())
		if err != nil {
			log.Fatalf("failed to end operational intent: %v", err.Error())
		}

		return nil
	},
}
