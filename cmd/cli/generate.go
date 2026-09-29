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

var generateAgentNetworkCmd = &cobra.Command{
	Use:     "agentnetwork",
	Aliases: []string{"agent", "ai"},
	Short:   "Agent Network manifest (LLM providers, guardrails, agent policies, budget rules)",
	Run: func(cmd *cobra.Command, args []string) {
		limits := func(userUSD float64) map[string]interface{} {
			return map[string]interface{}{
				"token_limit":  map[string]interface{}{"enabled": false, "group_cap": 0, "user_cap": 0, "window_seconds": 86400},
				"budget_limit": map[string]interface{}{"enabled": true, "group_cap_usd": 0, "user_cap_usd": userUSD, "window_seconds": 86400},
			}
		}
		printManifest(map[string]interface{}{
			"agentproviders": []map[string]interface{}{
				{
					"name":         "openai",
					"provider_id":  "openai_api",
					"upstream_url": "https://api.openai.com",
					"api_key":      "${OPENAI_API_KEY}",
					"models": []map[string]interface{}{
						{"id": "gpt-4o-mini", "input_per_1k": 0.00015, "output_per_1k": 0.0006},
					},
				},
			},
			"guardrails": []map[string]interface{}{
				{
					"name": "audited",
					"checks": map[string]interface{}{
						"model_allowlist": map[string]interface{}{"enabled": true, "models": []string{"gpt-4o-mini"}},
						"prompt_capture":  map[string]interface{}{"enabled": true, "redact_pii": true},
					},
				},
			},
			"agentpolicies": []map[string]interface{}{
				{
					"name":                     "devs-llm",
					"source_groups":            []string{"devs"},
					"destination_provider_ids": []string{"openai"},
					"guardrail_ids":            []string{"audited"},
					"limits":                   limits(10),
				},
			},
			"budgetrules": []map[string]interface{}{
				{
					"name":          "devs-daily",
					"target_groups": []string{"devs"},
					"target_users":  []string{},
					"limits":        limits(5),
				},
			},
		})
	},
}

func init() {
	generateCmd.AddCommand(generateUserCmd, generateGroupCmd, generatePeerCmd,
		generatePolicyCmd, generateSetupKeyCmd, generateAccountCmd, generateAgentNetworkCmd, generateAllCmd)
	rootCmd.AddCommand(generateCmd)
}
