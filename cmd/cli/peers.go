package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var peersGetCmd = &cobra.Command{
	Use:               "peers [nom|id]",
	Aliases:           []string{"pr"},
	Short:             "List or display peers",
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			peers, err := c.GetPeers(nameFlag, ipFlag)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			printOutput(peers)
			return
		}
		id, err := c.ResolvePeerID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		peer, err := c.GetPeerByID(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(peer)
	},
}

var peerEditCmd = &cobra.Command{
	Use:               "peer <nom|id>",
	Aliases:           []string{"pr"},
	Short:             "Edit a peer",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePeerID(args[0])
		if err != nil {
			exitErr("peer", err)
			return
		}
		path := fmt.Sprintf("/api/peers/%s", id)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get peer", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("ssh") {
			overrides["ssh_enabled"] = sshFlag
		}
		if cmd.Flags().Changed("login-expiration") {
			overrides["login_expiration_enabled"] = loginExpFlag
		}
		if cmd.Flags().Changed("inactivity-expiration") {
			overrides["inactivity_expiration_enabled"] = inactivityExpFlag
		}
		if cmd.Flags().Changed("ip") {
			overrides["ip"] = ipFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit peer", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var peerDeleteCmd = &cobra.Command{
	Use:               "peer <nom|id>",
	Aliases:           []string{"pr"},
	Short:             "Delete a peer",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePeerID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete peer %s", args[0])) {
			return
		}
		if err := c.DeletePeer(id); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("peer deleted")
	},
}

func init() {
	getCmd.AddCommand(peersGetCmd)
	editCmd.AddCommand(peerEditCmd)
	deleteCmd.AddCommand(peerDeleteCmd)

	peersGetCmd.Flags().StringVar(&nameFlag, "name", "", "Filter by name")
	peersGetCmd.Flags().StringVar(&ipFlag, "ip", "", "Filter by IP")

	peerEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	peerEditCmd.Flags().BoolVar(&sshFlag, "ssh", false, "Enable SSH")
	peerEditCmd.Flags().BoolVar(&loginExpFlag, "login-expiration", false, "Login expiration")
	peerEditCmd.Flags().BoolVar(&inactivityExpFlag, "inactivity-expiration", false, "Inactivity expiration")
	peerEditCmd.Flags().StringVar(&ipFlag, "ip", "", "IP address")
}
