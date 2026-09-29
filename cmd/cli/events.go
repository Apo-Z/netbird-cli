package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "View events",
}

var auditCmd = &cobra.Command{
	Use:     "audit",
	Aliases: []string{"ae"},
	Short:   "Audit events (default: last 50)",
	Run: func(cmd *cobra.Command, args []string) {
		events, err := c.GetAuditEvents()
		if err != nil {
			printErr(err)
			return
		}
		if pageSizeFlag > 0 && len(events) > pageSizeFlag {
			events = events[:pageSizeFlag]
		}
		printOutput(events)
	},
}

var trafficCmd = &cobra.Command{
	Use:     "traffic",
	Aliases: []string{"te"},
	Short:   "Network traffic events (cloud-only)",
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]string{}
		if pageFlag > 0 {
			params["page"] = fmt.Sprintf("%d", pageFlag)
		}
		if pageSizeFlag > 0 {
			params["page_size"] = fmt.Sprintf("%d", pageSizeFlag)
		}
		if searchFlag != "" {
			params["search"] = searchFlag
		}
		events, err := c.GetTrafficEvents(params)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				fmt.Println("Network traffic events are only available on Netbird Cloud, not on self-hosted instances.")
				return
			}
			printErr(err)
			return
		}
		printOutput(events)
	},
}

var proxylogsGetCmd = &cobra.Command{
	Use:   "proxylogs",
	Short: "Reverse proxy access logs",
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]string{}
		if pageFlag > 0 {
			params["page"] = fmt.Sprintf("%d", pageFlag)
		}
		if pageSizeFlag > 0 {
			params["page_size"] = fmt.Sprintf("%d", pageSizeFlag)
		}
		logs, err := c.GetProxyAccessLogs(params)
		if err != nil {
			exitErr("proxy logs", err)
			return
		}
		printOutput(logs)
	},
}

func init() {
	eventsCmd.AddCommand(auditCmd, trafficCmd)
	rootCmd.AddCommand(eventsCmd)

	getCmd.AddCommand(proxylogsGetCmd)

	auditCmd.Flags().IntVar(&pageSizeFlag, "page-size", 50, "Number of events (0 = all)")

	trafficCmd.Flags().IntVar(&pageFlag, "page", 1, "Page number")
	trafficCmd.Flags().IntVar(&pageSizeFlag, "page-size", 50, "Page size")
	trafficCmd.Flags().StringVar(&searchFlag, "search", "", "Text search")

	proxylogsGetCmd.Flags().IntVar(&pageFlag, "page", 1, "Page number")
	proxylogsGetCmd.Flags().IntVar(&pageSizeFlag, "page-size", 50, "Page size")
}
