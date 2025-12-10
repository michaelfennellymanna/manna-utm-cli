package riddp

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna.utm.cli/pkg/config"
	"manna.aero/manna.utm.cli/pkg/riddp_server"
)

var RidDP = &cobra.Command{
	Use:   "riddp",
	Short: "Start the Remote Id Display Provider server on <rid_dp_port>.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		writeRequests, err := cmd.Flags().GetBool("dump-requests")
		if err != nil {
			return err
		}

		println(writeRequests)
		c, err := config.LoadConfig("./config.yaml")
		if err != nil {
			return err
		}

		router := riddp_server.GetServer(*c, writeRequests)
		log.Debugf("attempting to start server on port: %d", c.RidDpPort)
		err = router.Run(fmt.Sprintf(":%d", c.RidDpPort))
		if err != nil {
			log.Fatalf("unable to start application server: %v", err)
		}

		log.Infof("Server listening on port: %d", c.RidDpPort)
		return nil
	},
}
