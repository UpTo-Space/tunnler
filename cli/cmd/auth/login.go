package auth

import (
	"github.com/UpTo-Space/tunnler/client/client"
	"github.com/spf13/cobra"
)

// LoginCmd represents the auth register command
var LoginCmd = &cobra.Command{
	Use:   "login [username] [password]",
	Short: "login with user account",
	Long:  `Login a user to the auth server.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]

		ai := client.TunnlerAuthConnectionInfo{
			HostAdress: authServerAdress,
			HostPort:   authServerPort,
			HostScheme: authServerScheme,
		}

		authClient := client.NewAuthClient(ai)

		authClient.Login(username, password)
	},
}

func init() {
}
