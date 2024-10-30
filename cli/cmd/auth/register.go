package auth

import (
	"fmt"
	"strings"

	"github.com/UpTo-Space/tunnler/client/client"
	"github.com/spf13/cobra"
)

// RegisterCmd represents the auth register command
var RegisterCmd = &cobra.Command{
	Use:   "register [username] [password] [email]",
	Short: "register user account",
	Long:  `Registration to the auth server.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(3)(cmd, args); err != nil {
			return err
		}

		email := args[2]
		if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			return fmt.Errorf("Your E-Mail seems to not be valid %v", email)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		email := args[2]

		ai := client.TunnlerAuthConnectionInfo{
			HostAdress: authServerAdress,
			HostPort:   authServerPort,
			HostScheme: authServerScheme,
		}

		authClient := client.NewAuthClient(ai)

		authClient.Register(username, password, email)
	},
}

func init() {
}
