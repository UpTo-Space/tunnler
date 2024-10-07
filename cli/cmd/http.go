/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/UpTo-Space/tunnler/client/client"
	"github.com/spf13/cobra"
)

var (
	tunnlerServerAdress string
	tunnlerServerPort   string
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http [adress:port | port] [flags]",
	Short: "Create a http tunnel",
	Long:  `Create a http tunnel`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		portString := args[0]

		if strings.Contains(args[0], ":") {
			options := strings.Split(args[0], ":")
			portString = options[1]
		}

		if _, err := strconv.Atoi(portString); err != nil {
			return fmt.Errorf("%v is not a valid adress:port or port", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var hostString, portString string

		if strings.Contains(args[0], ":") {
			options := strings.Split(args[0], ":")
			hostString = options[0]
			portString = options[1]
		} else {
			hostString = "127.0.0.1"
			portString = args[0]
		}

		ci := client.TunnlerConnectionInfo{
			HostAdress:    hostString,
			HostPort:      portString,
			TunnlerAdress: tunnlerServerAdress,
			TunnlerPort:   tunnlerServerPort,
		}

		client := client.NewTunnlerClient(ci)
		client.Connect()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.PersistentFlags().StringVar(&tunnlerServerAdress, "server", "127.0.0.1", "IP / Domain of the tunnler server to connect to")
	httpCmd.PersistentFlags().StringVar(&tunnlerServerPort, "serverPort", "8888", "Port of the tunnler server to connect to")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
