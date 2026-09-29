package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var nameserversGetCmd = &cobra.Command{
	Use:               "nameservers [id]",
	Aliases:           []string{"ns"},
	Short:             "List or display nameserver groups",
	ValidArgsFunction: validArgsFunc(nameserverGroupNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			groups, err := c.GetNameserverGroups()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(groups)
			return
		}
		group, err := c.GetNameserverGroupByID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(group)
	},
}

var nameserverCreateCmd = &cobra.Command{
	Use:     "nameserver",
	Aliases: []string{"ns"},
	Short:   "Create a nameserver group",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/dns/nameservers", map[string]interface{}{
				"name":                   "",
				"description":            "",
				"nameservers":            []map[string]interface{}{{"ip": "8.8.8.8", "ns_type": "udp", "port": 53}},
				"enabled":                true,
				"groups":                 []string{},
				"primary":                false,
				"domains":                []string{},
				"search_domains_enabled": false,
			})
			if err != nil {
				exitErr("create nameserver", err)
				return
			}
			if result == nil {
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/dns/nameservers", result)
			if err != nil {
				exitErr("create nameserver", err)
				return
			}
			var nsg client.NameserverGroup
			json.Unmarshal(data, &nsg)
			fmt.Println("group created:")
			printOutput(nsg)
			return
		}
		var nss []client.Nameserver
		for _, ip := range nameserversFlag {
			nss = append(nss, client.Nameserver{IP: ip, NSType: dnsTypeFlag, Port: portFlag})
		}
		req := &client.CreateNameserverGroupRequest{
			Name:                 nameFlag,
			Description:          descFlag,
			Nameservers:          nss,
			Enabled:              enabledFlag,
			Groups:               groupsFlag,
			Primary:              primaryFlag,
			Domains:              domainsFlag,
			SearchDomainsEnabled: searchDomainFlag,
		}
		if dryRunCheck(req) {
			return
		}
		group, err := c.CreateNameserverGroup(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("nameserver group created:")
		printOutput(group)
	},
}

var nameserverEditCmd = &cobra.Command{
	Use:               "nameserver <id>",
	Aliases:           []string{"ns"},
	Short:             "Edit a nameserver group",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(nameserverGroupNames),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("/api/dns/nameservers/%s", args[0])
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get nameserver", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		if cmd.Flags().Changed("ns") {
			var nss []map[string]interface{}
			for _, ip := range nameserversFlag {
				nss = append(nss, map[string]interface{}{"ip": ip, "ns_type": dnsTypeFlag, "port": portFlag})
			}
			overrides["nameservers"] = nss
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("groups") {
			overrides["groups"] = groupsFlag
		}
		if cmd.Flags().Changed("primary") {
			overrides["primary"] = primaryFlag
		}
		if cmd.Flags().Changed("domains") {
			overrides["domains"] = domainsFlag
		}
		if cmd.Flags().Changed("search-domains") {
			overrides["search_domains_enabled"] = searchDomainFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit nameserver", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var nameserverDeleteCmd = &cobra.Command{
	Use:               "nameserver <id>",
	Aliases:           []string{"ns"},
	Short:             "Delete a nameserver group",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(nameserverGroupNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete nameserver %s", args[0])) {
			return
		}
		if err := c.DeleteNameserverGroup(args[0]); err != nil {
			printErr(err)
			return
		}
		fmt.Println("nameserver group deleted")
	},
}

var dnssettingsGetCmd = &cobra.Command{
	Use:     "dnssettings",
	Aliases: []string{"ds"},
	Short:   "Show DNS settings",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := c.GetDNSSettings()
		if err != nil {
			printErr(err)
			return
		}
		printOutput(settings)
	},
}

var dnssettingsEditCmd = &cobra.Command{
	Use:     "dnssettings",
	Aliases: []string{"ds"},
	Short:   "Edit DNS settings",
	Run: func(cmd *cobra.Command, args []string) {
		path := "/api/dns/settings"
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get dnssettings", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("disabled-groups") {
			overrides["disabled_management_groups"] = groupsFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit dnssettings", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

func init() {
	getCmd.AddCommand(nameserversGetCmd, dnssettingsGetCmd)
	createCmd.AddCommand(nameserverCreateCmd)
	editCmd.AddCommand(nameserverEditCmd, dnssettingsEditCmd)
	deleteCmd.AddCommand(nameserverDeleteCmd)

	nameserverCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Group name")
	nameserverCreateCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	nameserverCreateCmd.Flags().StringSliceVar(&nameserversFlag, "ns", nil, "IPs des nameservers")
	nameserverCreateCmd.Flags().StringVar(&dnsTypeFlag, "ns-type", "udp", "Type (udp/tcp)")
	nameserverCreateCmd.Flags().IntVar(&portFlag, "port", 53, "Port")
	nameserverCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	nameserverCreateCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs or names of groups")
	nameserverCreateCmd.Flags().BoolVar(&primaryFlag, "primary", false, "Primary group")
	nameserverCreateCmd.Flags().StringSliceVar(&domainsFlag, "domains", nil, "Domains")
	nameserverCreateCmd.Flags().BoolVar(&searchDomainFlag, "search-domains", false, "Domain search")

	nameserverEditCmd.Flags().StringVar(&nameFlag, "name", "", "Nom")
	nameserverEditCmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	nameserverEditCmd.Flags().StringSliceVar(&nameserversFlag, "ns", nil, "IPs")
	nameserverEditCmd.Flags().StringVar(&dnsTypeFlag, "ns-type", "udp", "Type")
	nameserverEditCmd.Flags().IntVar(&portFlag, "port", 53, "Port")
	nameserverEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	nameserverEditCmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "IDs ou noms des groupes")
	nameserverEditCmd.Flags().BoolVar(&primaryFlag, "primary", false, "Primary")
	nameserverEditCmd.Flags().StringSliceVar(&domainsFlag, "domains", nil, "Domaines")
	nameserverEditCmd.Flags().BoolVar(&searchDomainFlag, "search-domains", false, "Recherche de domaines")

	dnssettingsEditCmd.Flags().StringSliceVar(&groupsFlag, "disabled-groups", nil, "Groups with DNS management disabled")
}
