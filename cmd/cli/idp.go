package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var idpsGetCmd = &cobra.Command{
	Use:               "identityproviders [id]",
	Aliases:           []string{"idp"},
	Short:             "List or display identity providers",
	ValidArgsFunction: validArgsFunc(idpNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			idps, err := c.GetIdentityProviders()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(idps)
			return
		}
		idp, err := c.GetIdentityProviderByID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(idp)
	},
}

var idpCreateCmd = &cobra.Command{
	Use:     "identityprovider",
	Aliases: []string{"idp"},
	Short:   "Create an identity provider",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/identity-providers", map[string]interface{}{
				"type":          "oidc",
				"name":          "",
				"issuer":        "",
				"client_id":     "",
				"client_secret": "",
			})
			if err != nil { exitErr("create idp", err); return }
			if result == nil { return }
			if dryRunCheck(result) { return }
			data, err := c.PostRaw("/api/identity-providers", result)
			if err != nil { exitErr("create idp", err); return }
			var idp client.IdentityProvider
			json.Unmarshal(data, &idp)
			fmt.Println("provider created:")
			printOutput(idp)
			return
		}
		req := &client.CreateIDPRequest{
			Type:         idpTypeFlag,
			Name:         nameFlag,
			Issuer:       issuerFlag,
			ClientID:     clientIDFlag,
			ClientSecret: clientSecretFlag,
		}
		if dryRunCheck(req) {
			return
		}
		idp, err := c.CreateIdentityProvider(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("identity provider created:")
		printOutput(idp)
	},
}

var idpEditCmd = &cobra.Command{
	Use:               "identityprovider <id>",
	Aliases:           []string{"idp"},
	Short:             "Edit an identity provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(idpNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/identity-providers/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get idp", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("issuer") {
			overrides["issuer"] = issuerFlag
		}
		if cmd.Flags().Changed("client-id") {
			overrides["client_id"] = clientIDFlag
		}
		if cmd.Flags().Changed("client-secret") {
			overrides["client_secret"] = clientSecretFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit idp", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var idpDeleteCmd = &cobra.Command{
	Use:               "identityprovider <id>",
	Aliases:           []string{"idp"},
	Short:             "Delete an identity provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(idpNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete identity provider %s", args[0])) {
			return
		}
		if err := c.DeleteIdentityProvider(args[0]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("identity provider deleted")
	},
}

func init() {
	getCmd.AddCommand(idpsGetCmd)
	createCmd.AddCommand(idpCreateCmd)
	editCmd.AddCommand(idpEditCmd)
	deleteCmd.AddCommand(idpDeleteCmd)

	idpCreateCmd.Flags().StringVar(&idpTypeFlag, "type", "oidc", "Type (oidc)")
	idpCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Nom")
	idpCreateCmd.Flags().StringVar(&issuerFlag, "issuer", "", "OIDC issuer URL")
	idpCreateCmd.Flags().StringVar(&clientIDFlag, "client-id", "", "OAuth2 client ID")
	idpCreateCmd.Flags().StringVar(&clientSecretFlag, "client-secret", "", "OAuth2 client secret")

	idpEditCmd.Flags().StringVar(&nameFlag, "name", "", "Nom")
	idpEditCmd.Flags().StringVar(&issuerFlag, "issuer", "", "Issuer URL")
	idpEditCmd.Flags().StringVar(&clientIDFlag, "client-id", "", "Client ID")
	idpEditCmd.Flags().StringVar(&clientSecretFlag, "client-secret", "", "Client secret")
}
