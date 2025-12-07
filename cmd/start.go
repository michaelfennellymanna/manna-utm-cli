package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/server"
)

var Start = &cobra.Command{
	Use:   "start",
	Short: "Start the application listening on the provided server port.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		if len(args) == 1 {
			port = 38080
			log.Infof("no port defined, starting application server on port %d", port)
		}

		log.SetLevel(log.InfoLevel)

		router := server.GetServer(port)
		log.Debugf("attempting to start server on port: %d", port)
		err = router.Run(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("unable to start application server: %v", err)
		}

		log.Infof("Server listening on port: %d", port)
		return nil
	},
}
