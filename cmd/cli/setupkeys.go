package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var setupkeysGetCmd = &cobra.Command{
	Use:               "setupkeys [id]",
	Aliases:           []string{"sk"},
	Short:             "List or display setup keys",
	ValidArgsFunction: validArgsFunc(setupKeyNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			keys, err := c.GetSetupKeys()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(keys)
			return
		}
		key, err := c.GetSetupKeyByID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(key)
	},
}

var setupkeyCreateCmd = &cobra.Command{
	Use:     "setupkey",
	Aliases: []string{"sk"},
	Short:   "Create a setup key",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/setup-keys", map[string]interface{}{
				"name":                   "",
				"type":                   "reusable",
				"expires_in":             2592000,
				"auto_groups":            []string{},
				"usage_limit":            0,
				"ephemeral":              false,
				"allow_extra_dns_labels": false,
			})
			if err != nil { exitErr("create setupkey", err); return }
			if result == nil { return }
			if dryRunCheck(result) { return }
			data, err := c.PostRaw("/api/setup-keys", result)
			if err != nil { exitErr("create setupkey", err); return }
			var key client.SetupKey
			json.Unmarshal(data, &key)
			fmt.Println("key created:")
			printOutput(key)
			return
		}
		req := &client.CreateSetupKeyRequest{
			Name:       nameFlag,
			Type:       keyTypeFlag,
			ExpiresIn:  expireFlag,
			AutoGroups: autoGroupsFlag,
			UsageLimit: usageLimitFlag,
		}
		if cmd.Flags().Changed("ephemeral") {
			req.Ephemeral = &ephemeralFlag
		}
		if cmd.Flags().Changed("allow-dns-labels") {
			req.AllowExtraDNSLabels = &allowDNSLabelsFlag
		}
		if dryRunCheck(req) {
			return
		}
		key, err := c.CreateSetupKey(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("key created:")
		printOutput(key)
	},
}

var setupkeyEditCmd = &cobra.Command{
	Use:               "setupkey <id>",
	Aliases:           []string{"sk"},
	Short:             "Edit a setup key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(setupKeyNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/setup-keys/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get setupkey", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("revoked") {
			overrides["revoked"] = revokedFlag
		}
		if cmd.Flags().Changed("auto-groups") {
			overrides["auto_groups"] = autoGroupsFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit setupkey", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var setupkeyDeleteCmd = &cobra.Command{
	Use:               "setupkey <id>",
	Aliases:           []string{"sk"},
	Short:             "Delete a setup key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(setupKeyNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete key %s", args[0])) {
			return
		}
		if err := c.DeleteSetupKey(args[0]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("key deleted")
	},
}

func init() {
	getCmd.AddCommand(setupkeysGetCmd)
	createCmd.AddCommand(setupkeyCreateCmd)
	editCmd.AddCommand(setupkeyEditCmd)
	deleteCmd.AddCommand(setupkeyDeleteCmd)

	setupkeyCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Key name")
	setupkeyCreateCmd.Flags().StringVar(&keyTypeFlag, "type", "reusable", "Type (one-off/reusable)")
	setupkeyCreateCmd.Flags().IntVar(&expireFlag, "expire", 86400*30, "Expiration in seconds")
	setupkeyCreateCmd.Flags().StringSliceVar(&autoGroupsFlag, "auto-groups", nil, "IDs or names of auto-assigned groups")
	setupkeyCreateCmd.Flags().IntVar(&usageLimitFlag, "usage-limit", 0, "Usage limit (0=unlimited)")
	setupkeyCreateCmd.Flags().BoolVar(&ephemeralFlag, "ephemeral", false, "Ephemeral peer")
	setupkeyCreateCmd.Flags().BoolVar(&allowDNSLabelsFlag, "allow-dns-labels", false, "Allow extra DNS labels")

	setupkeyEditCmd.Flags().BoolVar(&revokedFlag, "revoked", false, "Revoke key")
	setupkeyEditCmd.Flags().StringSliceVar(&autoGroupsFlag, "auto-groups", nil, "IDs or names of auto-assigned groups")

	setupkeyCreateCmd.RegisterFlagCompletionFunc("type", staticCompletion([]string{"one-off", "reusable"}))
	setupkeyCreateCmd.RegisterFlagCompletionFunc("auto-groups", validArgsFunc(groupNames))
	setupkeyEditCmd.RegisterFlagCompletionFunc("auto-groups", validArgsFunc(groupNames))
}
