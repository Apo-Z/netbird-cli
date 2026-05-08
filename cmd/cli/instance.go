package main

import (
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var instanceStatusCmd = &cobra.Command{
	Use:   "instancestatus",
	Short: "Show instance status (setup required?)",
	Run: func(cmd *cobra.Command, args []string) {
		status, err := c.GetInstanceStatus()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(status)
	},
}

var instanceVersionCmd = &cobra.Command{
	Use:   "instanceversion",
	Short: "Show instance versions",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := c.GetInstanceVersion()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(version)
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize instance (first admin + token)",
	Run: func(cmd *cobra.Command, args []string) {
		req := &client.SetupRequest{
			Email:    emailFlag,
			Password: passwordFlag,
			Name:     nameFlag,
		}
		if cmd.Flags().Changed("create-pat") {
			pat := true
			req.CreatePAT = &pat
		}
		if cmd.Flags().Changed("pat-expire") {
			req.PatExpireIn = &patExpireFlag
		}
		if dryRunCheck(req) {
			return
		}
		result, err := c.SetupInstance(req)
		if err != nil {
			exitErr("setup", err)
			return
		}
		fmt.Println("instance initialized:")
		printOutput(result)
	},
}

func init() {
	getCmd.AddCommand(instanceStatusCmd, instanceVersionCmd)
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().StringVar(&emailFlag, "email", "", "Admin email")
	setupCmd.Flags().StringVar(&passwordFlag, "password", "", "Admin password (≥8 chars)")
	setupCmd.Flags().StringVar(&nameFlag, "name", "", "Admin name")
	setupCmd.Flags().Bool("create-pat", false, "Create admin PAT")
	setupCmd.Flags().IntVar(&patExpireFlag, "pat-expire", 1, "PAT expiration in days")
}
