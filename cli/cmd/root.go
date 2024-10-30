package cmd

import (
	"os"

	"github.com/UpTo-Space/tunnler/client/cmd/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "Create a connection to an UpToSpace Tunnel Server",
	Long: `Establish a connection to a UpToSpace Tunnel Server.
	Expose local ports of your applications and make them avaiable in the WWW
	without worring about firewalls etc.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
}
