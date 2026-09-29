package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var posturechecksGetCmd = &cobra.Command{
	Use:               "posturechecks [id]",
	Aliases:           []string{"pc"},
	Short:             "List or display posture checks",
	ValidArgsFunction: validArgsFunc(postureCheckNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			checks, err := c.GetPostureChecks()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(checks)
			return
		}
		check, err := c.GetPostureCheckByID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(check)
	},
}

var posturecheckCreateCmd = &cobra.Command{
	Use:     "posturecheck",
	Aliases: []string{"pc"},
	Short:   "Create a posture check",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/posture-checks", map[string]interface{}{
				"name":        "",
				"description": "",
			})
			if err != nil {
				exitErr("create posturecheck", err)
				return
			}
			if result == nil {
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/posture-checks", result)
			if err != nil {
				exitErr("create posturecheck", err)
				return
			}
			var pc client.PostureCheck
			json.Unmarshal(data, &pc)
			fmt.Println("posture check created:")
			printOutput(pc)
			return
		}
		req := &client.CreatePostureCheckRequest{
			Name:        nameFlag,
			Description: descFlag,
		}
		if dryRunCheck(req) {
			return
		}
		check, err := c.CreatePostureCheck(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("posture check created:")
		printOutput(check)
	},
}

var posturecheckEditCmd = &cobra.Command{
	Use:               "posturecheck <id>",
	Aliases:           []string{"pc"},
	Short:             "Edit a posture check",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(postureCheckNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/posture-checks/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get posturecheck", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit posturecheck", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var posturecheckDeleteCmd = &cobra.Command{
	Use:               "posturecheck <id>",
	Aliases:           []string{"pc"},
	Short:             "Delete a posture check",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(postureCheckNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete posture check %s", args[0])) {
			return
		}
		if err := c.DeletePostureCheck(args[0]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("posture check deleted")
	},
}

func init() {
	getCmd.AddCommand(posturechecksGetCmd)
	createCmd.AddCommand(posturecheckCreateCmd)
	editCmd.AddCommand(posturecheckEditCmd)
	deleteCmd.AddCommand(posturecheckDeleteCmd)

	posturecheckCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Posture check name")
	posturecheckCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")

	posturecheckEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	posturecheckEditCmd.Flags().StringVar(&descFlag, "desc", "", "New description")
}
