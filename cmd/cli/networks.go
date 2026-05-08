package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

func networkResourceArgsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) == 0 {
		networks, err := c.GetNetworks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, n := range networks {
			names = append(names, fmt.Sprintf("%s\t(id=%s)", n.Name, n.ID))
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
	netID, err := c.ResolveNetworkID(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	resources, err := c.GetNetworkResources(netID)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, r := range resources {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", r.Name, r.ID))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func resolveNetworkResourceID(networkID, nameOrID string) (string, error) {
	resources, err := c.GetNetworkResources(networkID)
	if err != nil {
		return "", err
	}
	for _, r := range resources {
		if r.Name == nameOrID || r.ID == nameOrID {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("resource %s not found in network", nameOrID)
}

var networksGetCmd = &cobra.Command{
	Use:               "networks [nom|id]",
	Aliases:           []string{"nw"},
	Short:             "List or display networks",
	ValidArgsFunction: validArgsFunc(networkNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			networks, err := c.GetNetworks()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			printOutput(networks)
			return
		}
		id, err := c.ResolveNetworkID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		network, err := c.GetNetworkByID(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(network)
	},
}

var networkCreateCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"nw"},
	Short:   "Create a network",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/networks", map[string]interface{}{
				"name":        "",
				"description": "",
			})
			if err != nil { exitErr("create network", err); return }
			if dryRunCheck(result) { return }
			data, err := c.PostRaw("/api/networks", result)
			if err != nil { exitErr("create network", err); return }
			var net client.Network
			json.Unmarshal(data, &net)
			fmt.Println("network created:")
			printOutput(net)
			return
		}
		req := &client.CreateNetworkRequest{
			Name:        nameFlag,
			Description: descFlag,
		}
		if dryRunCheck(req) {
			return
		}
		network, err := c.CreateNetwork(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("network created:")
		printOutput(network)
	},
}

var networkEditCmd = &cobra.Command{
	Use:               "network <nom|id>",
	Aliases:           []string{"nw"},
	Short:             "Edit a network",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(networkNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveNetworkID(args[0])
		if err != nil {
			exitErr("network", err)
			return
		}
		path := fmt.Sprintf("/api/networks/%s", id)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get network", err)
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
			exitErr("edit network", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var networkDeleteCmd = &cobra.Command{
	Use:               "network <nom|id>",
	Aliases:           []string{"nw"},
	Short:             "Delete a network",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(networkNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveNetworkID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete network %s", args[0])) {
			return
		}
		if err := c.DeleteNetwork(id); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("network deleted")
	},
}

var networkresourcesGetCmd = &cobra.Command{
	Use:               "networkresources <network-name|id>",
	Aliases:           []string{"nwr"},
	Short:             "List network resources",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(networkNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveNetworkID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		resources, err := c.GetNetworkResources(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(resources)
	},
}

var networkresourceCreateCmd = &cobra.Command{
	Use:     "networkresource",
	Aliases: []string{"nwr"},
	Short:   "Create a network resource",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/networks/<network-id>/resources", map[string]interface{}{
				"network_id":  networkFlag,
				"name":        "",
				"description": "",
				"address":     "",
				"enabled":     true,
				"groups":      []string{},
			})
			if err != nil {
				exitErr("create networkresource", err)
				return
			}
			if result == nil {
				return
			}
			netID, _ := result["network_id"].(string)
			delete(result, "network_id")
			if netID == "" {
				exitErr("network_id required in YAML", nil)
				return
			}
			id, err := c.ResolveNetworkID(netID)
			if err != nil {
				exitErr("network not found", err)
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw(fmt.Sprintf("/api/networks/%s/resources", id), result)
			if err != nil {
				exitErr("create networkresource", err)
				return
			}
			var res client.NetworkResource
			json.Unmarshal(data, &res)
			fmt.Println("resource created:")
			printOutput(res)
			return
		}
		id, err := c.ResolveNetworkID(networkFlag)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		req := &client.CreateNetworkResourceRequest{
			Name:        nameFlag,
			Description: descFlag,
			Address:     addrFlag,
			Enabled:     enabledFlag,
			Groups:      groupsFlag,
		}
		if dryRunCheck(req) {
			return
		}
		resource, err := c.CreateNetworkResource(id, req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("resource created:")
		printOutput(resource)
	},
}

var networkresourceEditCmd = &cobra.Command{
	Use:               "networkresource <network-name|id> <resource-id>",
	Aliases:           []string{"nwr"},
	Short:             "Edit a network resource",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: networkResourceArgsCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		netID, err := c.ResolveNetworkID(args[0])
		if err != nil {
			exitErr("network", err)
			return
		}
		resID, err := resolveNetworkResourceID(netID, args[1])
		if err != nil {
			exitErr("networkresource", err)
			return
		}
		path := fmt.Sprintf("/api/networks/%s/resources/%s", netID, resID)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get networkresource", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		if cmd.Flags().Changed("address") {
			overrides["address"] = addrFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("groups") {
			overrides["groups"] = groupsFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit networkresource", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var networkresourceDeleteCmd = &cobra.Command{
	Use:               "networkresource <network-name|id> <resource-id>",
	Aliases:           []string{"nwr"},
	Short:             "Delete a network resource",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: networkResourceArgsCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		netID, err := c.ResolveNetworkID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		resID, err := resolveNetworkResourceID(netID, args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg("network resource deletion in dry-run") {
			return
		}
		if err := c.DeleteNetworkResource(netID, resID); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("resource deleted")
	},
}

func init() {
	getCmd.AddCommand(networksGetCmd, networkresourcesGetCmd)
	createCmd.AddCommand(networkCreateCmd, networkresourceCreateCmd)
	editCmd.AddCommand(networkEditCmd, networkresourceEditCmd)
	deleteCmd.AddCommand(networkDeleteCmd, networkresourceDeleteCmd)

	networkCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Network name")
	networkCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")

	networkEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	networkEditCmd.Flags().StringVar(&descFlag, "desc", "", "New description")

	networkresourceEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	networkresourceEditCmd.Flags().StringVar(&descFlag, "desc", "", "New description")
	networkresourceEditCmd.Flags().StringVar(&addrFlag, "address", "", "Address (IP, CIDR, domain)")
	networkresourceEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	networkresourceEditCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs or names of groups")

	networkresourceCreateCmd.Flags().StringVar(&networkFlag, "network", "", "Network name or ID")
	networkresourceCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Resource name")
	networkresourceCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	networkresourceCreateCmd.Flags().StringVar(&addrFlag, "address", "", "Address (IP, CIDR, domain)")
	networkresourceCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	networkresourceCreateCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs or names of groups")
}
