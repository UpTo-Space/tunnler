package auth

import (
	"github.com/spf13/cobra"
)

var (
	authServerAdress string
	authServerPort   string
	authServerScheme string
)

// authCmd represents the auth command
var AuthCmd = &cobra.Command{
	Use:   "auth [register | login]",
	Short: "authentication commands",
	Long:  `Login or Registration to the auth server.`,
}

func init() {
	AuthCmd.AddCommand(RegisterCmd)
	AuthCmd.AddCommand(LoginCmd)
	AuthCmd.PersistentFlags().StringVar(&authServerAdress, "server", "127.0.0.1", "IP / Domain of the tunnler auth server to connect to")
	AuthCmd.PersistentFlags().StringVar(&authServerPort, "serverPort", "8887", "Port of the tunnler auth server to connect to")
	AuthCmd.PersistentFlags().StringVar(&authServerScheme, "serverScheme", "http", "Scheme of the tunnler auth server")
}
