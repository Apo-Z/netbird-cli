package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var servicesGetCmd = &cobra.Command{
	Use:               "services [id]",
	Aliases:           []string{"svc"},
	Short:             "List or display reverse proxy services",
	ValidArgsFunction: validArgsFunc(serviceNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			services, err := c.GetServices()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			printOutput(services)
			return
		}
		service, err := c.GetServiceByID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(service)
	},
}

var serviceCreateCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Create a reverse proxy service",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/reverse-proxies/services", map[string]interface{}{
				"name":                "",
				"domain":              "",
				"mode":                "http",
				"enabled":             true,
				"listen_port":         0,
				"pass_host_header":    false,
				"rewrite_redirects":  false,
				"targets":             []map[string]interface{}{},
			})
			if err != nil { exitErr("create service", err); return }
			if result == nil { return }
			if dryRunCheck(result) { return }
			data, err := c.PostRaw("/api/reverse-proxies/services", result)
			if err != nil { exitErr("create service", err); return }
			var svc client.Service
			json.Unmarshal(data, &svc)
			fmt.Println("service created:")
			printOutput(svc)
			return
		}
		req := &client.CreateServiceRequest{
			Name:     nameFlag,
			Domain:   domainFlag,
			Mode:     modeFlag,
			Enabled:  enabledFlag,
		}
		if cmd.Flags().Changed("listen-port") {
			req.ListenPort = listenPortFlag
		}
		if dryRunCheck(req) {
			return
		}
		service, err := c.CreateService(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("service created:")
		printOutput(service)
	},
}

var serviceEditCmd = &cobra.Command{
	Use:               "service <id>",
	Aliases:           []string{"svc"},
	Short:             "Edit a reverse proxy service",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(serviceNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/reverse-proxies/services/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get service", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("domain") {
			overrides["domain"] = domainFlag
		}
		if cmd.Flags().Changed("mode") {
			overrides["mode"] = modeFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit service", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var serviceDeleteCmd = &cobra.Command{
	Use:               "service <id>",
	Aliases:           []string{"svc"},
	Short:             "Delete a service",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(serviceNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete service %s", args[0])) {
			return
		}
		if err := c.DeleteService(args[0]); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("service deleted")
	},
}

var proxyclustersGetCmd = &cobra.Command{
	Use:     "proxyclusters",
	Aliases: []string{"pxc"},
	Short:   "List available proxy clusters",
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := c.GetProxyClusters()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(clusters)
	},
}

func init() {
	getCmd.AddCommand(servicesGetCmd, proxyclustersGetCmd)
	createCmd.AddCommand(serviceCreateCmd)
	editCmd.AddCommand(serviceEditCmd)
	deleteCmd.AddCommand(serviceDeleteCmd)

	serviceCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Service name")
	serviceCreateCmd.Flags().StringVar(&domainFlag, "domain", "", "Domaine")
	serviceCreateCmd.Flags().StringVar(&modeFlag, "mode", "http", "Mode (http/tcp/udp/tls)")
	serviceCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	serviceCreateCmd.Flags().IntVar(&listenPortFlag, "listen-port", 0, "Listening port (0=auto)")

	serviceEditCmd.Flags().StringVar(&nameFlag, "name", "", "Nom")
	serviceEditCmd.Flags().StringVar(&domainFlag, "domain", "", "Domaine")
	serviceEditCmd.Flags().StringVar(&modeFlag, "mode", "", "Mode")
	serviceEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
}

