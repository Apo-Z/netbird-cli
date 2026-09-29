package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var accountsGetCmd = &cobra.Command{
	Use:     "accounts [id]",
	Aliases: []string{"ac"},
	Short:   "List accounts",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			accounts, err := c.GetAccounts()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(accounts)
			return
		}
		// There's no GET /api/accounts/{id} — list and filter
		accounts, err := c.GetAccounts()
		if err != nil {
			printErr(err)
			return
		}
		for _, a := range accounts {
			if a.ID == args[0] {
				printOutput(a)
				return
			}
		}
		fmt.Printf("compte %s not found\n", args[0])
	},
}

var accountEditCmd = &cobra.Command{
	Use:     "account <id>",
	Aliases: []string{"ac"},
	Short:   "Edit account settings (opens editor)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := c.GetAccounts()
		if err != nil {
			exitErr("get accounts", err)
			return
		}
		var accountID string
		for _, a := range accounts {
			if a.ID == args[0] {
				accountID = a.ID
				break
			}
		}
		if accountID == "" {
			exitErr("account not found", nil)
			return
		}

		path := fmt.Sprintf("/api/accounts/%s", accountID)
		// GetRaw doesn't work (no GET /api/accounts/{id}) — marshal from list
		var accountData []byte
		for _, a := range accounts {
			if a.ID == accountID {
				accountData, _ = json.Marshal(a)
				break
			}
		}

		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("peer-login-expiration") {
			overrides["peer_login_expiration"] = peerLoginExpFlag
		}
		if cmd.Flags().Changed("peer-inactivity-expiration") {
			overrides["peer_inactivity_expiration"] = peerInactivityExpFlag
		}
		if cmd.Flags().Changed("peer-approval") {
			overrides["peer_approval_enabled"] = peerApprovalFlag
		}
		if cmd.Flags().Changed("dns-domain") {
			overrides["dns_domain"] = domainFlag
		}
		if cmd.Flags().Changed("network-range") {
			overrides["network_range"] = cidrFlag
		}

		result, err := editWithEditor(path, accountData, overrides)
		if err != nil {
			exitErr("edit account", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var tokensGetCmd = &cobra.Command{
	Use:               "tokens <user-id>",
	Aliases:           []string{"tk"},
	Short:             "List user tokens",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := c.ResolveUserID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		tokens, err := c.GetTokens(userID)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(tokens)
	},
}

var tokenCreateCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"tk"},
	Short:   "Create a token",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/users/<user-id>/tokens", map[string]interface{}{
				"user_id":    userFlag,
				"name":       "",
				"expires_in": 30,
			})
			if err != nil {
				exitErr("create token", err)
				return
			}
			if result == nil {
				return
			}
			userID, _ := result["user_id"].(string)
			delete(result, "user_id")
			if userID == "" {
				exitErr("user_id required in YAML", nil)
				return
			}
			resolvedID, err := c.ResolveUserID(userID)
			if err == nil {
				userID = resolvedID
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/users/"+userID+"/tokens", result)
			if err != nil { exitErr("create token", err); return }
			var token client.CreateTokenResponse
			json.Unmarshal(data, &token)
			fmt.Println("token created:")
			printOutput(token)
			return
		}
		req := &client.CreateTokenRequest{
			Name:      nameFlag,
			ExpiresIn: expireFlag,
		}
		if dryRunCheck(req) {
			return
		}
		userID, err := c.ResolveUserID(userFlag)
		if err != nil {
			printErr(err)
			return
		}
		result, err := c.CreateToken(userID, req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("token created (save the plain_token value):")
		printOutput(result)
	},
}

var tokenDeleteCmd = &cobra.Command{
	Use:               "token <user-id> <token-id>",
	Aliases:           []string{"tk"},
	Short:             "Delete a token",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := c.ResolveUserID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete token %s", args[1])) {
			return
		}
		if err := c.DeleteToken(userID, args[1]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("token deleted")
	},
}

func init() {
	getCmd.AddCommand(accountsGetCmd, tokensGetCmd)
	createCmd.AddCommand(tokenCreateCmd)
	editCmd.AddCommand(accountEditCmd)
	deleteCmd.AddCommand(tokenDeleteCmd)

	accountEditCmd.Flags().IntVar(&peerLoginExpFlag, "peer-login-expiration", 0, "Peer login expiration (seconds)")
	accountEditCmd.Flags().IntVar(&peerInactivityExpFlag, "peer-inactivity-expiration", 0, "Peer inactivity expiration (seconds)")
	accountEditCmd.Flags().BoolVar(&peerApprovalFlag, "peer-approval", false, "Enable peer approval")
	accountEditCmd.Flags().StringVar(&domainFlag, "dns-domain", "", "Custom DNS domain")
	accountEditCmd.Flags().StringVar(&cidrFlag, "network-range", "", "Network CIDR range")

	tokenCreateCmd.Flags().StringVar(&userFlag, "user", "", "User ID or email")
	tokenCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Token name")
	tokenCreateCmd.Flags().IntVar(&expireFlag, "expire", 30, "Expiration in days")
}
