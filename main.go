package main

import (
	"log"

	"github.com/spf13/cobra"
	"manna.aero/manna-utm-geojson-api/cmd"
)

var (
	port     int
	fromFile string
)

var rootCmd = &cobra.Command{
	Use:   "manna-utm-cli",
	Short: "A command line tool for working with manna-utm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	cmd.Start.Flags().IntVarP(&port, "port", "p", 38080, "Listen port to bind the server to.")
	cmd.Data.Flags().StringVar(&fromFile, "file", "./.libconfig/sequence.yaml", "The path to the directory that you want to write the JSON contents of the simulation data to.")
}

func main() {
	rootCmd.AddCommand(cmd.Start)
	rootCmd.AddCommand(cmd.Data)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
