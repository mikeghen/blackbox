package cmd

import (
	"github.com/crypdex/blackbox/cmd/blackbox"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logsCmd)
}

// cleanupCmd represents the cleanup command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs of all running containers",

	Run: func(cmd *cobra.Command, args []string) {
		status := blackbox.Logs()

		log("info", status.Stdout...)
		// log("error", status.Stderr...)
	},
}
