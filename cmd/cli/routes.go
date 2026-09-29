package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var routesGetCmd = &cobra.Command{
	Use:               "routes [id]",
	Aliases:           []string{"rt"},
	Short:             "List or display routes (deprecated)",
	ValidArgsFunction: validArgsFunc(routeNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			routes, err := c.GetRoutes()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(routes)
			return
		}
		route, err := c.GetRouteByID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(route)
	},
}

var routeCreateCmd = &cobra.Command{
	Use:     "route",
	Aliases: []string{"rt"},
	Short:   "Create a route",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/routes", map[string]interface{}{
				"description":    "",
				"network_id":     "",
				"enabled":        true,
				"network":        "",
				"domains":        []string{},
				"metric":         9999,
				"masquerade":     true,
				"groups":         []string{},
				"keep_route":     true,
				"skip_auto_apply": false,
			})
			if err != nil {
				exitErr("create route", err)
				return
			}
			if result == nil {
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/routes", result)
			if err != nil {
				exitErr("create route", err)
				return
			}
			var route client.Route
			json.Unmarshal(data, &route)
			fmt.Println("route created:")
			printOutput(route)
			return
		}
		req := &client.CreateRouteRequest{
			Description:   descFlag,
			NetworkID:     nameFlag,
			Enabled:       enabledFlag,
			Peer:          peerFlag,
			PeerGroups:    peersFlag,
			Network:       cidrFlag,
			Domains:       domainsFlag,
			Metric:        metricFlag,
			Masquerade:    masqueradeFlag,
			Groups:        groupsFlag,
			KeepRoute:     keepRouteFlag,
			SkipAutoApply: &skipAutoApplyFlag,
		}
		if dryRunCheck(req) {
			return
		}
		route, err := c.CreateRoute(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("route created:")
		printOutput(route)
	},
}

var routeEditCmd = &cobra.Command{
	Use:               "route <id>",
	Aliases:           []string{"rt"},
	Short:             "Edit a route",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(routeNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/routes/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get route", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		if cmd.Flags().Changed("network-id") {
			overrides["network_id"] = nameFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("peer") {
			overrides["peer"] = peerFlag
		}
		if cmd.Flags().Changed("peer-groups") {
			overrides["peer_groups"] = peersFlag
		}
		if cmd.Flags().Changed("network") {
			overrides["network"] = cidrFlag
		}
		if cmd.Flags().Changed("domains") {
			overrides["domains"] = domainsFlag
		}
		if cmd.Flags().Changed("metric") {
			overrides["metric"] = metricFlag
		}
		if cmd.Flags().Changed("masquerade") {
			overrides["masquerade"] = masqueradeFlag
		}
		if cmd.Flags().Changed("groups") {
			overrides["groups"] = groupsFlag
		}
		if cmd.Flags().Changed("keep-route") {
			overrides["keep_route"] = keepRouteFlag
		}
		if cmd.Flags().Changed("skip-auto-apply") {
			overrides["skip_auto_apply"] = skipAutoApplyFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit route", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:               "route <id>",
	Aliases:           []string{"rt"},
	Short:             "Delete a route",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(routeNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete route %s", args[0])) {
			return
		}
		if err := c.DeleteRoute(args[0]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("route deleted")
	},
}

func init() {
	getCmd.AddCommand(routesGetCmd)
	createCmd.AddCommand(routeCreateCmd)
	editCmd.AddCommand(routeEditCmd)
	deleteCmd.AddCommand(routeDeleteCmd)

	routeCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	routeCreateCmd.Flags().StringVar(&nameFlag, "network-id", "", "Network ID")
	routeCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	routeCreateCmd.Flags().StringVar(&peerFlag, "peer", "", "Peer ID or name")
	routeCreateCmd.Flags().StringSliceVar(&peersFlag, "peer-groups", nil, "IDs or names of peer groups")
	routeCreateCmd.Flags().StringVar(&cidrFlag, "network", "", "Network CIDR")
	routeCreateCmd.Flags().StringSliceVar(&domainsFlag, "domains", nil, "Domains")
	routeCreateCmd.Flags().IntVar(&metricFlag, "metric", 9999, "Metric (1-9999)")
	routeCreateCmd.Flags().BoolVar(&masqueradeFlag, "masquerade", true, "Masquerade")
	routeCreateCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs or names of routing groups")
	routeCreateCmd.Flags().BoolVar(&keepRouteFlag, "keep-route", true, "Keep route after DNS resolution")
	routeCreateCmd.Flags().BoolVar(&skipAutoApplyFlag, "skip-auto-apply", false, "Do not auto-apply")

	routeEditCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	routeEditCmd.Flags().StringVar(&nameFlag, "network-id", "", "Network ID")
	routeEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	routeEditCmd.Flags().StringVar(&peerFlag, "peer", "", "Peer ID or name")
	routeEditCmd.Flags().StringSliceVar(&peersFlag, "peer-groups", nil, "IDs or names of peer groups")
	routeEditCmd.Flags().StringVar(&cidrFlag, "network", "", "Network CIDR")
	routeEditCmd.Flags().StringSliceVar(&domainsFlag, "domains", nil, "Domains")
	routeEditCmd.Flags().IntVar(&metricFlag, "metric", 9999, "Metric")
	routeEditCmd.Flags().BoolVar(&masqueradeFlag, "masquerade", true, "Masquerade")
	routeEditCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs or names of groups")
	routeEditCmd.Flags().BoolVar(&keepRouteFlag, "keep-route", true, "Keep route")
	routeEditCmd.Flags().BoolVar(&skipAutoApplyFlag, "skip-auto-apply", false, "Skip auto-apply")
}
