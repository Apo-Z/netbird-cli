package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var policiesGetCmd = &cobra.Command{
	Use:               "policies [nom|id]",
	Aliases:           []string{"pol"},
	Short:             "List or display policies",
	ValidArgsFunction: validArgsFunc(policyNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			policies, err := c.GetPolicies()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			printOutput(policies)
			return
		}
		id, err := c.ResolvePolicyID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		policy, err := c.GetPolicyByID(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(policy)
	},
}

var policyCreateCmd = &cobra.Command{
	Use:     "policy",
	Aliases: []string{"pol"},
	Short:   "Create a policy",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/policies", map[string]interface{}{
				"name":                  "",
				"description":           "",
				"enabled":               true,
				"source_posture_checks": []string{},
				"rules":                 []interface{}{},
			})
			if err != nil {
				exitErr("create policy", err)
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/policies", result)
			if err != nil {
				exitErr("create policy", err)
				return
			}
			var pol client.Policy
			json.Unmarshal(data, &pol)
			fmt.Println("policy created:")
			printOutput(pol)
			return
		}
		req := &client.CreatePolicyRequest{
			Name:        nameFlag,
			Description: descFlag,
			Enabled:     enabledFlag,
		}
		if dryRunCheck(req) {
			return
		}
		policy, err := c.CreatePolicy(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("policy created:")
		printOutput(policy)
	},
}

var policyEditCmd = &cobra.Command{
	Use:               "policy <nom|id>",
	Aliases:           []string{"pol"},
	Short:             "Edit a policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(policyNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePolicyID(args[0])
		if err != nil {
			exitErr("policy", err)
			return
		}
		path := fmt.Sprintf("/api/policies/%s", id)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get policy", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit policy", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var policyDeleteCmd = &cobra.Command{
	Use:               "policy <nom|id>",
	Aliases:           []string{"pol"},
	Short:             "Delete a policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(policyNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePolicyID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete policy %s", args[0])) {
			return
		}
		if err := c.DeletePolicy(id); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("policy deleted")
	},
}

func init() {
	getCmd.AddCommand(policiesGetCmd)
	createCmd.AddCommand(policyCreateCmd)
	editCmd.AddCommand(policyEditCmd)
	deleteCmd.AddCommand(policyDeleteCmd)

	policyCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Policy name")
	policyCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	policyCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")

	policyEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	policyEditCmd.Flags().StringVar(&descFlag, "desc", "", "New description")
	policyEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
}
