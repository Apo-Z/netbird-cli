package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

func dnsRecordArgsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) == 0 {
		zones, err := c.GetDNSZones()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []string
		for _, z := range zones {
			names = append(names, fmt.Sprintf("%s\t(id=%s domain=%s)", z.Name, z.ID, z.Domain))
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
	zoneID, err := c.ResolveDNSZoneID(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	recs, err := c.GetDNSRecords(zoneID)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, r := range recs {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", r.Name, r.ID))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func resolveDNSRecordID(zoneID, nameOrID string) (string, error) {
	recs, err := c.GetDNSRecords(zoneID)
	if err != nil {
		return "", err
	}
	for _, r := range recs {
		if r.Name == nameOrID || r.ID == nameOrID {
			return r.ID, nil
		}
	}
	return "", fmt.Errorf("record %s not found in zone", nameOrID)
}

var dnszonesGetCmd = &cobra.Command{
	Use:               "dnszones [id]",
	Aliases:           []string{"dz"},
	Short:             "List or display DNS zones",
	ValidArgsFunction: validArgsFunc(dnsZoneNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			zones, err := c.GetDNSZones()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			printOutput(zones)
			return
		}
		id, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		zone, err := c.GetDNSZoneByID(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(zone)
	},
}

var dnszoneCreateCmd = &cobra.Command{
	Use:     "dnszone",
	Aliases: []string{"dz"},
	Short:   "Create a DNS zone",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/dns/zones", map[string]interface{}{
				"name":                 "",
				"domain":               "",
				"enabled":              true,
				"enable_search_domain": false,
				"distribution_groups":  []string{},
			})
			if err != nil { exitErr("create dnszone", err); return }
			if result == nil { return }
			if dryRunCheck(result) { return }
			data, err := c.PostRaw("/api/dns/zones", result)
			if err != nil { exitErr("create dnszone", err); return }
			var zone client.DNSZone
			json.Unmarshal(data, &zone)
			fmt.Println("zone created:")
			printOutput(zone)
			return
		}
		req := &client.CreateDNSZoneRequest{
			Name:               nameFlag,
			Domain:             domainFlag,
			Enabled:            &enabledFlag,
			EnableSearchDomain: searchDomainFlag,
			DistributionGroups: distGroupsFlag,
		}
		if dryRunCheck(req) {
			return
		}
		zone, err := c.CreateDNSZone(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("DNS zone created:")
		printOutput(zone)
	},
}

var dnszoneEditCmd = &cobra.Command{
	Use:               "dnszone <id>",
	Aliases:           []string{"dz"},
	Short:             "Edit a DNS zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(dnsZoneNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			exitErr("dnszone", err)
			return
		}
		path := fmt.Sprintf("/api/dns/zones/%s", id)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get dnszone", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("domain") {
			overrides["domain"] = domainFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("search-domain") {
			overrides["enable_search_domain"] = searchDomainFlag
		}
		if cmd.Flags().Changed("groups") {
			overrides["distribution_groups"] = distGroupsFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit dnszone", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var dnszoneDeleteCmd = &cobra.Command{
	Use:               "dnszone <id>",
	Aliases:           []string{"dz"},
	Short:             "Delete a DNS zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(dnsZoneNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete DNS zone %s", args[0])) {
			return
		}
		if err := c.DeleteDNSZone(id); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("DNS zone deleted")
	},
}

var dnsrecordsGetCmd = &cobra.Command{
	Use:               "dnsrecords <zone-id>",
	Aliases:           []string{"dr"},
	Short:             "List DNS records for a zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(dnsZoneNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		recs, err := c.GetDNSRecords(id)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		printOutput(recs)
	},
}

var dnsrecordCreateCmd = &cobra.Command{
	Use:     "dnsrecord",
	Aliases: []string{"dr"},
	Short:   "Create a DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			result, err := createWithEditor("/api/dns/zones/<zone-id>/records", map[string]interface{}{
				"zone_id": zoneFlag,
				"name":    "",
				"type":    "A",
				"content": "",
				"ttl":     300,
			})
			if err != nil {
				exitErr("create dnsrecord", err)
				return
			}
			if result == nil {
				return
			}
			zoneID, _ := result["zone_id"].(string)
			delete(result, "zone_id")
			if zoneID == "" {
				exitErr("zone required in YAML", nil)
				return
			}
			if dryRunCheck(result) {
				return
			}
			data, err := c.PostRaw("/api/dns/zones/"+zoneID+"/records", result)
			if err != nil { exitErr("create dnsrecord", err); return }
			var rec client.DNSRecord
			json.Unmarshal(data, &rec)
			fmt.Println("record created:")
			printOutput(rec)
			return
		}
		zoneID, err := c.ResolveDNSZoneID(zoneFlag)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		req := &client.CreateDNSRecordRequest{
			Name:    nameFlag,
			Type:    dnsTypeFlag,
			Content: contentFlag,
			TTL:     ttlFlag,
		}
		if dryRunCheck(req) {
			return
		}
		rec, err := c.CreateDNSRecord(zoneID, req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("record created:")
		printOutput(rec)
	},
}

var dnsrecordEditCmd = &cobra.Command{
	Use:               "dnsrecord <zone-id> <record-id>",
	Aliases:           []string{"dr"},
	Short:             "Edit a DNS record",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: dnsRecordArgsCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		zoneID, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			exitErr("dnsrecord", err)
			return
		}
		recID, err := resolveDNSRecordID(zoneID, args[1])
		if err != nil {
			exitErr("dnsrecord", err)
			return
		}
		path := fmt.Sprintf("/api/dns/zones/%s/records/%s", zoneID, recID)
		data, err := c.GetRaw(path)
		if err != nil {
			exitErr("get dnsrecord", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("type") {
			overrides["type"] = dnsTypeFlag
		}
		if cmd.Flags().Changed("content") {
			overrides["content"] = contentFlag
		}
		if cmd.Flags().Changed("ttl") {
			overrides["ttl"] = ttlFlag
		}
		result, err := editWithEditor(path, data, overrides)
		if err != nil {
			exitErr("edit dnsrecord", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var dnsrecordDeleteCmd = &cobra.Command{
	Use:               "dnsrecord <zone-id> <record-id>",
	Aliases:           []string{"dr"},
	Short:             "Delete a DNS record",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: dnsRecordArgsCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		zoneID, err := c.ResolveDNSZoneID(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		recID, err := resolveDNSRecordID(zoneID, args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		if dryRunMsg("DNS record deletion in dry-run") {
			return
		}
		if err := c.DeleteDNSRecord(zoneID, recID); err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Println("record deleted")
	},
}

func init() {
	getCmd.AddCommand(dnszonesGetCmd, dnsrecordsGetCmd)
	createCmd.AddCommand(dnszoneCreateCmd, dnsrecordCreateCmd)
	editCmd.AddCommand(dnszoneEditCmd, dnsrecordEditCmd)
	deleteCmd.AddCommand(dnszoneDeleteCmd, dnsrecordDeleteCmd)

	dnszoneCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Zone name")
	dnszoneCreateCmd.Flags().StringVar(&domainFlag, "domain", "", "Domain (FQDN)")
	dnszoneCreateCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	dnszoneCreateCmd.Flags().BoolVar(&searchDomainFlag, "search-domain", false, "Search zone")
	dnszoneCreateCmd.Flags().StringSliceVar(&distGroupsFlag, "groups", nil, "IDs or names of distribution groups")

dnszoneEditCmd.Flags().StringVar(&nameFlag, "name", "", "Name")
dnszoneEditCmd.Flags().StringVar(&domainFlag, "domain", "", "Domain")
dnszoneEditCmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	dnszoneEditCmd.Flags().BoolVar(&searchDomainFlag, "search-domain", false, "Zone de recherche")
	dnszoneEditCmd.Flags().StringSliceVar(&distGroupsFlag, "groups", nil, "IDs or names of groups")

	dnsrecordCreateCmd.Flags().StringVar(&zoneFlag, "zone", "", "DNS zone name or ID")
	dnsrecordCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Record FQDN")
	dnsrecordCreateCmd.Flags().StringVar(&dnsTypeFlag, "type", "A", "Type (A, AAAA, CNAME)")
	dnsrecordCreateCmd.Flags().StringVar(&contentFlag, "content", "", "Content (IP, domain)")
	dnsrecordCreateCmd.Flags().IntVar(&ttlFlag, "ttl", 300, "TTL")

	dnsrecordEditCmd.Flags().StringVar(&nameFlag, "name", "", "FQDN")
	dnsrecordEditCmd.Flags().StringVar(&dnsTypeFlag, "type", "A", "Type")
	dnsrecordEditCmd.Flags().StringVar(&contentFlag, "content", "", "Contenu")
	dnsrecordEditCmd.Flags().IntVar(&ttlFlag, "ttl", 300, "TTL")
}
