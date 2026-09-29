package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

// Reverse proxy domains / proxy tokens / clusters, and ingress peers / port allocations.

// --- Domains ---

var domainsGetCmd = &cobra.Command{
	Use:     "domains",
	Aliases: []string{"domain", "dom"},
	Short:   "List reverse proxy service domains (free and custom)",
	Run: func(cmd *cobra.Command, args []string) {
		domains, err := c.GetReverseProxyDomains()
		if err != nil {
			exitErr("domains", err)
			return
		}
		printOutput(domains)
	},
}

var domainCreateCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"dom"},
	Short:   "Add a custom reverse proxy domain",
	Long: `Add a custom domain for reverse proxy services, validated against a proxy cluster
(see: netbird get proxyclusters). Point the domain's DNS to the cluster, then run
"netbird validate domain <domain>".`,
	Example: `  netbird create domain --domain apps.example.com --cluster eu.proxy.netbird.io`,
	Run: func(cmd *cobra.Command, args []string) {
		if domainFlag == "" || clusterFlag == "" {
			exitErr("--domain and --cluster are required", nil)
			return
		}
		req := &client.ReverseProxyDomainRequest{Domain: domainFlag, TargetCluster: clusterFlag}
		if dryRunCheck(req) {
			return
		}
		data, err := c.CreateReverseProxyDomain(req)
		if err != nil {
			exitErr("create domain", err)
			return
		}
		var out map[string]interface{}
		json.Unmarshal(data, &out)
		fmt.Println("domain added:")
		if outputFormat == "" {
			printJSON(out)
			return
		}
		printOutput(out)
	},
}

var domainDeleteCmd = &cobra.Command{
	Use:               "domain <domain|id>",
	Aliases:           []string{"dom"},
	Short:             "Delete a custom reverse proxy domain",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(domainNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveReverseProxyDomainID(args[0])
		if err != nil {
			exitErr("domain", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete domain %s", args[0])) {
			return
		}
		if err := c.DeleteReverseProxyDomain(id); err != nil {
			exitErr("delete domain", err)
			return
		}
		fmt.Println("domain deleted")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate domain <domain|id>",
	Short: "Trigger validation of a custom reverse proxy domain",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{"domain"}, cobra.ShellCompDirectiveNoFileComp
		}
		return validArgsFunc(domainNames)(cmd, args, toComplete)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "domain" && args[0] != "dom" {
			fmt.Println("usage: netbird validate domain <domain|id>")
			return
		}
		id, err := c.ResolveReverseProxyDomainID(args[1])
		if err != nil {
			exitErr("domain", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would validate domain %s", args[1])) {
			return
		}
		if err := c.ValidateReverseProxyDomain(id); err != nil {
			exitErr("validate domain", err)
			return
		}
		fmt.Printf("validation of %s started; check with \"netbird get domains\"\n", args[1])
	},
}

// --- Proxy tokens ---

var proxyTokensGetCmd = &cobra.Command{
	Use:     "proxytokens",
	Aliases: []string{"proxytoken", "pxt"},
	Short:   "List tokens used by self-hosted reverse proxies to join",
	Run: func(cmd *cobra.Command, args []string) {
		tokens, err := c.GetProxyTokens()
		if err != nil {
			exitErr("proxy tokens", err)
			return
		}
		printOutput(tokens)
	},
}

var proxyTokenCreateCmd = &cobra.Command{
	Use:     "proxytoken",
	Aliases: []string{"pxt"},
	Short:   "Create a token for a self-hosted reverse proxy (the token is shown once)",
	Example: `  netbird create proxytoken --name eu-proxy --expire 7776000`,
	Run: func(cmd *cobra.Command, args []string) {
		if nameFlag == "" {
			exitErr("--name is required", nil)
			return
		}
		req := &client.ProxyTokenRequest{Name: nameFlag, ExpiresIn: expireFlag}
		if dryRunCheck(req) {
			return
		}
		token, err := c.CreateProxyToken(req)
		if err != nil {
			exitErr("create proxytoken", err)
			return
		}
		if outputFormat != "" {
			printOutput(token)
			return
		}
		fmt.Printf("proxy token %s created (id %s)\n", token.Name, token.ID)
		fmt.Printf("token (shown once): %s\n", token.PlainToken)
	},
}

var proxyTokenDeleteCmd = &cobra.Command{
	Use:               "proxytoken <name|id>",
	Aliases:           []string{"pxt"},
	Short:             "Revoke a reverse proxy token",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(proxyTokenNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveProxyTokenID(args[0])
		if err != nil {
			exitErr("proxy token", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would revoke proxy token %s", args[0])) {
			return
		}
		if err := c.RevokeProxyToken(id); err != nil {
			exitErr("revoke proxytoken", err)
			return
		}
		fmt.Println("proxy token revoked")
	},
}

var proxyClusterDeleteCmd = &cobra.Command{
	Use:               "proxycluster <address>",
	Aliases:           []string{"pxc"},
	Short:             "Delete a self-hosted reverse proxy cluster",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(proxyClusterNames),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete proxy cluster %s", args[0])) {
			return
		}
		if err := c.DeleteProxyCluster(args[0]); err != nil {
			exitErr("delete proxycluster", err)
			return
		}
		fmt.Println("proxy cluster deleted")
	},
}

// --- Ingress peers (cloud) ---

const ingressCloudOnly = "Ingress peers are only available on NetBird Cloud."

var ingressPeersGetCmd = &cobra.Command{
	Use:               "ingresspeers [peer|id]",
	Aliases:           []string{"ingresspeer", "ip"},
	Short:             "List or display ingress peers (public ingress IPs forwarding ports to your network, cloud)",
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			id, err := c.ResolveIngressPeerID(args[0])
			if err != nil {
				exitErr("ingress peer", err)
				return
			}
			p, err := c.GetIngressPeer(id)
			if err != nil {
				exitErr("ingress peer", err)
				return
			}
			printOutput(p)
			return
		}
		peers, err := c.GetIngressPeers()
		if err != nil {
			if is404(err) {
				fmt.Println(ingressCloudOnly)
				return
			}
			exitErr("ingress peers", err)
			return
		}
		printOutput(peers)
	},
}

var ingressPeerCreateCmd = &cobra.Command{
	Use:               "ingresspeer <peer>",
	Aliases:           []string{"ip"},
	Short:             "Turn a peer into an ingress peer (cloud)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		peerID, err := c.ResolvePeerID(args[0])
		if err != nil {
			exitErr("peer", err)
			return
		}
		req := &client.IngressPeerRequest{PeerID: peerID, Enabled: enabledFlag, Fallback: fallbackFlag}
		if dryRunCheck(req) {
			return
		}
		p, err := c.CreateIngressPeer(req)
		if err != nil {
			exitErr("create ingresspeer", err)
			return
		}
		fmt.Println("ingress peer created:")
		printOutput(p)
	},
}

var ingressPeerEditCmd = &cobra.Command{
	Use:               "ingresspeer <peer|id>",
	Aliases:           []string{"ip"},
	Short:             "Edit an ingress peer (enabled, fallback)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveIngressPeerID(args[0])
		if err != nil {
			exitErr("ingress peer", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("fallback") {
			overrides["fallback"] = fallbackFlag
		}
		editTransformed("ingresspeer", "/api/ingress/peers/"+id, func(m map[string]interface{}) {
			for k := range m {
				if k != "enabled" && k != "fallback" {
					delete(m, k)
				}
			}
		}, overrides)
	},
}

var ingressPeerDeleteCmd = &cobra.Command{
	Use:               "ingresspeer <peer|id>",
	Aliases:           []string{"ip"},
	Short:             "Remove an ingress peer",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveIngressPeerID(args[0])
		if err != nil {
			exitErr("ingress peer", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete ingress peer %s", args[0])) {
			return
		}
		if err := c.DeleteIngressPeer(id); err != nil {
			exitErr("delete ingresspeer", err)
			return
		}
		fmt.Println("ingress peer deleted")
	},
}

// --- Ingress port allocations (cloud) ---

type ingressPortRow struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	IngressIP string `json:"ingress_ip"`
	Region    string `json:"region"`
	Ports     string `json:"port_mappings"`
}

var ingressPortsGetCmd = &cobra.Command{
	Use:     "ingressports --peer <peer>",
	Aliases: []string{"ingressport", "iport"},
	Short:   "List the port allocations forwarded to a peer through ingress peers (cloud)",
	Run: func(cmd *cobra.Command, args []string) {
		peerID, ok := requirePeerFlag()
		if !ok {
			return
		}
		ports, err := c.GetIngressPorts(peerID)
		if err != nil {
			if is404(err) {
				fmt.Println(ingressCloudOnly)
				return
			}
			exitErr("ingress ports", err)
			return
		}
		if outputFormat != "" {
			printOutput(ports)
			return
		}
		var rows []ingressPortRow
		for _, p := range ports {
			var maps []string
			for _, m := range p.PortRangeMappings {
				maps = append(maps, fmt.Sprintf("%s → %s/%s", portRange(m.IngressStart, m.IngressEnd), portRange(m.TranslatedStart, m.TranslatedEnd), m.Protocol))
			}
			rows = append(rows, ingressPortRow{p.Name, p.Enabled, p.IngressIP, p.Region, strings.Join(maps, ", ")})
		}
		printOutput(rows)
	},
}

var ingressPortCreateCmd = &cobra.Command{
	Use:     "ingressport --peer <peer>",
	Aliases: []string{"iport"},
	Short:   "Forward ports from an ingress peer to a peer (cloud)",
	Example: `  netbird create ingressport --peer web-1 --name https --range 443/tcp
  netbird create ingressport --peer game-1 --name game --range 27015-27030/udp
  netbird create ingressport --peer web-1 --name random --direct 2/tcp     # 2 ports picked by NetBird`,
	Run: func(cmd *cobra.Command, args []string) {
		peerID, ok := requirePeerFlag()
		if !ok {
			return
		}
		if nameFlag == "" || (len(rangesFlag) == 0 && directFlag == "") {
			exitErr("--name and --range or --direct are required", nil)
			return
		}
		req := &client.IngressPortAllocationRequest{Name: nameFlag, Enabled: enabledFlag}
		for _, r := range rangesFlag {
			pr, err := parsePortRange(r)
			if err != nil {
				exitErr("--range", err)
				return
			}
			req.PortRanges = append(req.PortRanges, pr)
		}
		if directFlag != "" {
			count, proto, _ := strings.Cut(directFlag, "/")
			n, err := strconv.Atoi(count)
			if err != nil || proto == "" {
				exitErr(fmt.Sprintf("--direct %q: expected <count>/<tcp|udp|tcp/udp>", directFlag), nil)
				return
			}
			req.DirectPort = &client.IngressDirectPort{Count: n, Protocol: proto}
		}
		if dryRunCheck(req) {
			return
		}
		p, err := c.CreateIngressPort(peerID, req)
		if err != nil {
			exitErr("create ingressport", err)
			return
		}
		fmt.Println("port allocation created:")
		printOutput(p)
	},
}

var ingressPortEditCmd = &cobra.Command{
	Use:     "ingressport <name|id> --peer <peer>",
	Aliases: []string{"iport"},
	Short:   "Edit a port allocation (cloud)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		peerID, ok := requirePeerFlag()
		if !ok {
			return
		}
		id, err := c.ResolveIngressPortID(peerID, args[0])
		if err != nil {
			exitErr("port allocation", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		editTransformed("ingressport", fmt.Sprintf("/api/peers/%s/ingress/ports/%s", peerID, id), func(m map[string]interface{}) {
			// GET returns resolved mappings; PUT takes the requested port ranges.
			var ranges []map[string]interface{}
			if list, ok := m["port_range_mappings"].([]interface{}); ok {
				for _, item := range list {
					if pm, ok := item.(map[string]interface{}); ok {
						ranges = append(ranges, map[string]interface{}{"start": pm["translated_start"], "end": pm["translated_end"], "protocol": pm["protocol"]})
					}
				}
			}
			for k := range m {
				if k != "name" && k != "enabled" {
					delete(m, k)
				}
			}
			m["port_ranges"] = ranges
		}, overrides)
	},
}

var ingressPortDeleteCmd = &cobra.Command{
	Use:     "ingressport <name|id> --peer <peer>",
	Aliases: []string{"iport"},
	Short:   "Delete a port allocation (cloud)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		peerID, ok := requirePeerFlag()
		if !ok {
			return
		}
		id, err := c.ResolveIngressPortID(peerID, args[0])
		if err != nil {
			exitErr("port allocation", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete port allocation %s", args[0])) {
			return
		}
		if err := c.DeleteIngressPort(peerID, id); err != nil {
			exitErr("delete ingressport", err)
			return
		}
		fmt.Println("port allocation deleted")
	},
}

func requirePeerFlag() (string, bool) {
	if peerFlag == "" {
		exitErr("--peer is required", nil)
		return "", false
	}
	id, err := c.ResolvePeerID(peerFlag)
	if err != nil {
		exitErr("peer", err)
		return "", false
	}
	return id, true
}

// parsePortRange parses 443/tcp or 8000-8010/udp.
func parsePortRange(s string) (client.IngressPortRange, error) {
	ports, proto, ok := strings.Cut(s, "/")
	if !ok || proto == "" {
		return client.IngressPortRange{}, fmt.Errorf("%q: expected <port>[-<port>]/<tcp|udp|tcp/udp>", s)
	}
	startStr, endStr, isRange := strings.Cut(ports, "-")
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return client.IngressPortRange{}, fmt.Errorf("%q: invalid port", s)
	}
	end := start
	if isRange {
		if end, err = strconv.Atoi(endStr); err != nil {
			return client.IngressPortRange{}, fmt.Errorf("%q: invalid port", s)
		}
	}
	return client.IngressPortRange{Start: start, End: end, Protocol: proto}, nil
}

func portRange(start, end int) string {
	if start == end {
		return strconv.Itoa(start)
	}
	return fmt.Sprintf("%d-%d", start, end)
}

func init() {
	getCmd.AddCommand(domainsGetCmd, proxyTokensGetCmd, ingressPeersGetCmd, ingressPortsGetCmd)
	createCmd.AddCommand(domainCreateCmd, proxyTokenCreateCmd, ingressPeerCreateCmd, ingressPortCreateCmd)
	editCmd.AddCommand(ingressPeerEditCmd, ingressPortEditCmd)
	deleteCmd.AddCommand(domainDeleteCmd, proxyTokenDeleteCmd, proxyClusterDeleteCmd, ingressPeerDeleteCmd, ingressPortDeleteCmd)
	rootCmd.AddCommand(validateCmd)

	domainCreateCmd.Flags().StringVar(&domainFlag, "domain", "", "Custom domain name")
	domainCreateCmd.Flags().StringVar(&clusterFlag, "cluster", "", "Proxy cluster address to validate against")
	domainCreateCmd.RegisterFlagCompletionFunc("cluster", validArgsFunc(proxyClusterNames))

	proxyTokenCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Token name")
	proxyTokenCreateCmd.Flags().IntVar(&expireFlag, "expire", 0, "Expiration in seconds (0 = never)")

	for _, cmd := range []*cobra.Command{ingressPeerCreateCmd, ingressPeerEditCmd} {
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.Flags().BoolVar(&fallbackFlag, "fallback", false, "Use as fallback when no ingress peer exists in the region")
	}

	for _, cmd := range []*cobra.Command{ingressPortsGetCmd, ingressPortCreateCmd, ingressPortEditCmd, ingressPortDeleteCmd} {
		cmd.Flags().StringVar(&peerFlag, "peer", "", "Target peer (name or ID)")
		cmd.RegisterFlagCompletionFunc("peer", validArgsFunc(peerNames))
	}
	for _, cmd := range []*cobra.Command{ingressPortCreateCmd, ingressPortEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Allocation name")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	}
	ingressPortCreateCmd.Flags().StringSliceVar(&rangesFlag, "range", nil, "Port or range to forward, e.g. 443/tcp, 8000-8010/udp (repeatable)")
	ingressPortCreateCmd.Flags().StringVar(&directFlag, "direct", "", "Let NetBird pick <count>/<protocol> ports, e.g. 2/tcp")
}
