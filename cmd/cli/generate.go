package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var generateCmd = &cobra.Command{
	Use:   "generate <resource>",
	Short: "Generate a template YAML manifest",
	Long:  `Generate a YAML manifest file (or JSON with -o json) ready to edit and apply with "netbird apply -f".`,
}

func printManifest(v interface{}) {
	switch outputFormat {
	case "json":
		printJSON(v)
	default:
		data, _ := yaml.Marshal(v)
		fmt.Print(string(data))
	}
}

var generateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "User manifest",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name":          "John Doe",
					"email":         "john@example.com",
					"role":          "user",
					"auto_groups":   []string{"devs"},
					"is_service_user": false,
				},
			},
		})
	},
}

var generateGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group manifest",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"groups": []map[string]interface{}{
				{
					"name":  "devs",
					"peers": []string{"alice-laptop", "stage-host-1"},
				},
			},
		})
	},
}

var generatePeerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Peer manifest",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"peers": []map[string]interface{}{
				{
					"name":                "stage-host-1",
					"ssh":                 true,
					"login_expiration":     false,
					"inactivity_expiration": false,
				},
			},
		})
	},
}

var generatePolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Policy manifest",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"policies": []map[string]interface{}{
				{
					"name":    "Default Policy",
					"enabled": true,
				},
			},
		})
	},
}

var generateSetupKeyCmd = &cobra.Command{
	Use:   "setupkey",
	Short: "Setup key manifest",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"setupkeys": []map[string]interface{}{
				{
					"name":        "default",
					"type":        "reusable",
					"expires_in":  2592000,
					"auto_groups": []string{"devs"},
					"usage_limit": 0,
				},
			},
			"accounts": []map[string]interface{}{
				{
					"settings": map[string]interface{}{
						"peer_login_expiration_enabled":       true,
						"peer_login_expiration":              43200,
						"peer_inactivity_expiration_enabled": true,
						"peer_inactivity_expiration":         43200,
						"peer_approval_enabled":              false,
						"dns_domain":                         "",
						"network_range":                      "100.64.0.0/16",
						"auto_update_version":                "latest",
					},
				},
			},
		})
	},
}

var generateAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account manifest (global settings)",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"accounts": []map[string]interface{}{
				{
					"settings": map[string]interface{}{
						"peer_login_expiration_enabled":       true,
						"peer_login_expiration":              43200,
						"peer_inactivity_expiration_enabled": true,
						"peer_inactivity_expiration":         43200,
						"peer_approval_enabled":              false,
						"user_approval_required":             false,
						"dns_domain":                         "",
						"network_range":                      "100.64.0.0/16",
						"peer_expose_enabled":                false,
						"peer_expose_groups":                 []string{},
						"auto_update_version":                "latest",
						"auto_update_always":                 false,
						"network_traffic_logs_enabled":       false,
						"network_traffic_logs_groups":        []string{},
					},
				},
			},
		})
	},
}

var generateAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Global manifest (all types)",
	Run: func(cmd *cobra.Command, args []string) {
		printManifest(map[string]interface{}{
			"groups": []map[string]interface{}{
				{
					"name":  "devs",
					"peers": []string{"alice-laptop", "stage-host-1"},
				},
				{
					"name":  "ops",
					"peers": []string{"stage-host-1"},
				},
			},
			"users": []map[string]interface{}{
				{
					"name":          "Alice",
					"email":         "alice@example.com",
					"role":          "admin",
					"auto_groups":   []string{"devs"},
					"is_service_user": false,
				},
				{
					"name":          "Bob",
					"email":         "bob@example.com",
					"role":          "user",
					"auto_groups":   []string{"ops"},
					"is_service_user": false,
				},
			},
			"peers": []map[string]interface{}{
				{
					"name":                "stage-host-1",
					"ssh":                 true,
					"login_expiration":     false,
					"inactivity_expiration": false,
				},
			},
			"policies": []map[string]interface{}{
				{
					"name":    "Default Policy",
					"enabled": true,
				},
			},
			"setupkeys": []map[string]interface{}{
				{
					"name":        "default",
					"type":        "reusable",
					"expires_in":  2592000,
					"auto_groups": []string{"devs"},
					"usage_limit": 0,
				},
			},
		})
	},
}

func init() {
	generateCmd.AddCommand(generateUserCmd, generateGroupCmd, generatePeerCmd,
		generatePolicyCmd, generateSetupKeyCmd, generateAccountCmd, generateAllCmd)
	rootCmd.AddCommand(generateCmd)
}
