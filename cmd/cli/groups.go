package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var groupsGetCmd = &cobra.Command{
	Use:               "groups [nom|id]",
	Aliases:           []string{"grp"},
	Short:             "List or display groups",
	ValidArgsFunction: validArgsFunc(groupNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			groups, err := c.GetGroups()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(groups)
			return
		}
		id, err := c.ResolveGroupID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		group, err := c.GetGroupByID(id)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(group)
	},
}

var groupCreateCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"grp"},
	Short:   "Create a group",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/groups", map[string]interface{}{
				"name":  "",
				"peers": []string{},
			})
			if err != nil {
				exitErr("create group", err)
				return
			}
			if result == nil {
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/groups", result)
			if err != nil {
				exitErr("create group", err)
				return
			}
			var group client.Group
			json.Unmarshal(data, &group)
			fmt.Println("group created:")
			printOutput(group)
			return
		}
		req := &client.CreateGroupRequest{
			Name:  nameFlag,
			Peers: peersFlag,
		}
		if dryRunCheck(req) {
			return
		}
		group, err := c.CreateGroup(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("group created:")
		printOutput(group)
	},
}

var groupEditCmd = &cobra.Command{
	Use:               "group <nom|id>",
	Aliases:           []string{"grp"},
	Short:             "Edit a group",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(groupNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveGroupID(args[0])
		if err != nil {
			exitErr("group", err)
			return
		}
		path := fmt.Sprintf("/api/groups/%s", id)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get group", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("peers") {
			overrides["peers"] = peersFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit group", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:               "group <nom|id>",
	Aliases:           []string{"grp"},
	Short:             "Delete a group",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(groupNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveGroupID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete group %s", args[0])) {
			return
		}
		if err := c.DeleteGroup(id); err != nil {
			printErr(err)
			return
		}
		fmt.Println("group deleted")
	},
}

func init() {
	getCmd.AddCommand(groupsGetCmd)
	createCmd.AddCommand(groupCreateCmd)
	editCmd.AddCommand(groupEditCmd)
	deleteCmd.AddCommand(groupDeleteCmd)

	groupCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Group name")
	groupCreateCmd.Flags().StringSliceVar(&peersFlag, "peers", nil, "IDs or names of peers")

	groupEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	groupEditCmd.Flags().StringSliceVar(&peersFlag, "peers", nil, "IDs or names of peers")
}
