package cmd

import (
	"github.com/spf13/cobra"
)

var (
	serviceName string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "service-cli",
	Short: "description",
	Long:  `description`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
	}
}

func init() {

	RootCmd.PersistentFlags().StringVarP(&serviceName, "service", "", "", "service")

}
