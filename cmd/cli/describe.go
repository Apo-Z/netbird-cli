package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe <resource> <name>",
	Short: "Display full details of a resource (IDs resolved to names)",
}

var describeUserCmd = &cobra.Command{
	Use:               "user <nom|email|id>",
	Aliases:           []string{"usr"},
	Short:             "Detail a user",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := c.GetUser(args[0])
		if err != nil {
			exitErr("user", err)
			return
		}

		if outputFormat != "" {
			printOutput(user)
			return
		}

		fmt.Println("=== User ===")
		fmt.Printf("Name   : %s\n", user.Name)
		fmt.Printf("Email  : %s\n", user.Email)
		fmt.Printf("Role   : %s\n", user.Role)
		fmt.Printf("Status : %s\n", user.Status)
		fmt.Printf("ID     : %s\n", user.ID)

		if len(user.AutoGroups) > 0 {
			fmt.Println("\nAuto-assigned groups:")
			c.IDToName("group", "") // warm cache
			for _, gid := range user.AutoGroups {
				name := c.IDToName("group", gid)
				if name == "" {
					name = gid
				}
				group, _ := c.GetGroupByID(gid)
				if group != nil {
					fmt.Printf("  %s (%s) — %d peers, %d resources\n", name, gid, group.PeersCount, group.ResourcesCount)
				} else {
					fmt.Printf("  %s (%s)\n", name, gid)
				}
			}
		}

		peers, _ := c.GetPeers("", "")
		var userPeers []string
		for _, p := range peers {
			if p.UserID == user.ID {
				status := "disconnected"
				if p.Connected {
					status = "connected"
				}
				userPeers = append(userPeers, fmt.Sprintf("  %s (%s) — %s, %s", p.Name, p.IP, status, p.OS))
			}
		}
		if len(userPeers) > 0 {
			fmt.Println("\nRegistered peers:")
			for _, p := range userPeers {
				fmt.Println(p)
			}
		}
	},
}

var describeGroupCmd = &cobra.Command{
	Use:               "group <nom|id>",
	Aliases:           []string{"grp"},
	Short:             "Detail a group",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(groupNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := c.ResolveGroupID(args[0])
		group, err := c.GetGroupByID(id)
		if err != nil {
			exitErr("group", err)
			return
		}

		if outputFormat != "" {
			printOutput(group)
			return
		}

		fmt.Println("=== Group ===")
		fmt.Printf("Name    : %s\n", group.Name)
		fmt.Printf("ID      : %s\n", group.ID)
		fmt.Printf("Peers   : %d\n", group.PeersCount)
		fmt.Printf("Ress.   : %d\n", group.ResourcesCount)

		if len(group.Peers) > 0 {
			fmt.Println("\nPeers :")
			for _, p := range group.Peers {
				name := c.IDToName("peer", p.ID)
				if name == "" {
					name = p.ID
				}
				fmt.Printf("  %s (%s)\n", name, p.ID)
			}
		}

		if len(group.Resources) > 0 {
			fmt.Println("\nResources:")
			for _, r := range group.Resources {
				fmt.Printf("  %s (%s)\n", r.ID, r.Type)
			}
		}

		users, _ := c.GetUsers()
		var groupUsers []string
		for _, u := range users {
			for _, g := range u.AutoGroups {
				if g == group.ID || g == group.Name {
					groupUsers = append(groupUsers, u.Name)
				}
			}
		}
		if len(groupUsers) > 0 {
			fmt.Printf("\nAuto-assigned users: %s\n", strings.Join(groupUsers, ", "))
		}
	},
}

var describePeerCmd = &cobra.Command{
	Use:               "peer <nom|id>",
	Aliases:           []string{"pr"},
	Short:             "Detail a peer",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := c.ResolvePeerID(args[0])
		peer, err := c.GetPeerByID(id)
		if err != nil {
			exitErr("peer", err)
			return
		}

		if outputFormat != "" {
			printOutput(peer)
			return
		}

		fmt.Println("=== Peer ===")
		fmt.Printf("Name     : %s\n", peer.Name)
		fmt.Printf("Hostname : %s\n", peer.Hostname)
		fmt.Printf("IP       : %s\n", peer.IP)
		fmt.Printf("OS       : %s\n", peer.OS)
		fmt.Printf("Version  : %s\n", peer.Version)
		fmt.Printf("ID       : %s\n", peer.ID)
		status := "disconnected"
		if peer.Connected {
			status = "connected"
		}
		fmt.Printf("Status   : %s\n", status)

		userName := c.IDToName("user", peer.UserID)
		if userName == "" {
			userName = peer.UserID
		}
		fmt.Printf("User: %s (%s)\n", userName, peer.UserID)

		if len(peer.Groups) > 0 {
			fmt.Println("\nGroups:")
			for _, g := range peer.Groups {
				fmt.Printf("  %s (%s) — %d peers\n", g.Name, g.ID, g.PeersCount)
			}
		}

		fmt.Printf("\nSSH        : %v\n", peer.SSHEnabled)
		fmt.Printf("DNS label  : %s\n", peer.DNSLabel)
		fmt.Printf("Country/City: %s / %s\n", peer.CountryCode, peer.CityName)
	},
}

func init() {
	describeCmd.AddCommand(describeUserCmd, describeGroupCmd, describePeerCmd)
	rootCmd.AddCommand(describeCmd)
}
