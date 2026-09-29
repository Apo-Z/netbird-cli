package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func validArgsFunc(fetch func() ([]string, error)) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if c == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		items, err := fetch()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

func staticCompletion(options []string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return options, cobra.ShellCompDirectiveNoFileComp
	}
}

func userNames() ([]string, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, u := range users {
		names = append(names, u.Name)
		if u.Email != "" {
			names = append(names, u.Email)
		}
	}
	return names, nil
}

func groupNames() ([]string, error) {
	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, g := range groups {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", g.Name, g.ID))
	}
	return names, nil
}

func peerNames() ([]string, error) {
	peers, err := c.GetPeers("", "")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range peers {
		names = append(names, fmt.Sprintf("%s\t(ip=%s host=%s)", p.Name, p.IP, p.Hostname))
	}
	return names, nil
}

func policyNames() ([]string, error) {
	policies, err := c.GetPolicies()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range policies {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", p.Name, p.ID))
	}
	return names, nil
}

func networkNames() ([]string, error) {
	networks, err := c.GetNetworks()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, n := range networks {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", n.Name, n.ID))
	}
	return names, nil
}

func setupKeyNames() ([]string, error) {
	keys, err := c.GetSetupKeys()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, k := range keys {
		names = append(names, fmt.Sprintf("%s\t(id=%d)", k.Name, k.ID))
	}
	return names, nil
}

func routeNames() ([]string, error) {
	routes, err := c.GetRoutes()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, r := range routes {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", r.Description, r.ID))
	}
	return names, nil
}

func nameserverGroupNames() ([]string, error) {
	groups, err := c.GetNameserverGroups()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, g := range groups {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", g.Name, g.ID))
	}
	return names, nil
}

func dnsZoneNames() ([]string, error) {
	zones, err := c.GetDNSZones()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, z := range zones {
		names = append(names, fmt.Sprintf("%s\t(id=%s domain=%s)", z.Name, z.ID, z.Domain))
	}
	return names, nil
}

func serviceNames() ([]string, error) {
	services, err := c.GetServices()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, s := range services {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", s.Name, s.ID))
	}
	return names, nil
}

func countryNames() ([]string, error) {
	countries, err := c.GetCountries()
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func postureCheckNames() ([]string, error) {
	checks, err := c.GetPostureChecks()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, pc := range checks {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", pc.Name, pc.ID))
	}
	return names, nil
}

func idpNames() ([]string, error) {
	idps, err := c.GetIdentityProviders()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, i := range idps {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", i.Name, i.ID))
	}
	return names, nil
}

func inviteNames() ([]string, error) {
	invites, err := c.GetUserInvites()
	if err != nil {
		return nil, err
	}
	var items []string
	for _, inv := range invites {
		items = append(items, fmt.Sprintf("%s\t(name=%s)", inv.Email, inv.Name))
	}
	return items, nil
}

func agentProviderNames() ([]string, error) {
	providers, err := c.GetAgentProviders()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range providers {
		names = append(names, fmt.Sprintf("%s\t(%s)", p.Name, p.ProviderID))
	}
	return names, nil
}

func agentPolicyNames() ([]string, error) {
	policies, err := c.GetAgentPolicies()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range policies {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", p.Name, p.ID))
	}
	return names, nil
}

func guardrailNames() ([]string, error) {
	guardrails, err := c.GetAgentGuardrails()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, g := range guardrails {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", g.Name, g.ID))
	}
	return names, nil
}

func budgetRuleNames() ([]string, error) {
	rules, err := c.GetAgentBudgetRules()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, r := range rules {
		names = append(names, fmt.Sprintf("%s\t(id=%s)", r.Name, r.ID))
	}
	return names, nil
}

func agentCatalogIDs() ([]string, error) {
	catalog, err := c.GetAgentCatalogProviders()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, p := range catalog {
		ids = append(ids, fmt.Sprintf("%s\t%s", p.ID, p.Name))
	}
	return ids, nil
}

func eventStreamNames() ([]string, error) {
	streams, err := c.GetEventStreams()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, s := range streams {
		names = append(names, fmt.Sprintf("%d\t%s", s.ID, s.Platform))
	}
	return names, nil
}

func notificationTypeNames() ([]string, error) {
	types, err := c.GetNotificationTypes()
	if err != nil {
		return nil, err
	}
	var names []string
	for code, desc := range types {
		names = append(names, fmt.Sprintf("%s\t%s", code, desc))
	}
	return names, nil
}

func domainNames() ([]string, error) {
	domains, err := c.GetReverseProxyDomains()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, d := range domains {
		if d.Type == "custom" {
			names = append(names, d.Domain)
		}
	}
	return names, nil
}

func proxyTokenNames() ([]string, error) {
	tokens, err := c.GetProxyTokens()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, t := range tokens {
		if !t.Revoked {
			names = append(names, t.Name)
		}
	}
	return names, nil
}

func proxyClusterNames() ([]string, error) {
	clusters, err := c.GetProxyClusters()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, cl := range clusters {
		names = append(names, cl.Address)
	}
	return names, nil
}

func tenantNames() ([]string, error) {
	tenants, err := c.GetTenants()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, t := range tenants {
		names = append(names, fmt.Sprintf("%s\t%s (%s)", t.Name, t.Domain, t.Status))
	}
	return names, nil
}

// agentModelNames lists the models exposed by the account's agent providers (falling back to the catalog).
func agentModelNames() ([]string, error) {
	seen := map[string]bool{}
	var names []string
	add := func(id, desc string) {
		if id != "" && !seen[id] {
			seen[id] = true
			names = append(names, fmt.Sprintf("%s\t%s", id, desc))
		}
	}
	providers, err := c.GetAgentProviders()
	if err != nil {
		return nil, err
	}
	for _, p := range providers {
		for _, m := range p.Models {
			add(m.ID, p.Name)
		}
	}
	if len(names) == 0 {
		catalog, err := c.GetAgentCatalogProviders()
		if err != nil {
			return nil, err
		}
		for _, p := range catalog {
			for _, m := range p.Models {
				add(m.ID, p.Name)
			}
		}
	}
	return names, nil
}

// catalogModelCompletion completes --models on "create agentprovider" from the catalog of the chosen --type.
func catalogModelCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	catalog, err := c.GetAgentCatalogProviders()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, p := range catalog {
		if providerTypeFlag != "" && p.ID != providerTypeFlag {
			continue
		}
		for _, m := range p.Models {
			names = append(names, fmt.Sprintf("%s\t%s", m.ID, m.Label))
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func agentSessionIDs() ([]string, error) {
	resp, err := c.GetAgentAccessLogSessions(map[string]string{"page_size": "50"})
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, s := range resp.Data {
		if s.SessionID != "" {
			ids = append(ids, fmt.Sprintf("%s\t%s, %d requests, %s", s.SessionID, c.IDToName("user", s.UserID), s.RequestCount, s.StartedAt))
		}
	}
	return ids, nil
}

func notificationChannelIDs() ([]string, error) {
	channels, err := c.GetNotificationChannels()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, ch := range channels {
		ids = append(ids, fmt.Sprintf("%s\t%s %s", ch.ID, ch.Type, targetSummary(ch.Target)))
	}
	return ids, nil
}

func invoiceIDs() ([]string, error) {
	invoices, err := c.GetBillingInvoices()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, i := range invoices {
		ids = append(ids, fmt.Sprintf("%s\t%s → %s", i.ID, i.PeriodStart, i.PeriodEnd))
	}
	return ids, nil
}

func priceIDs() ([]string, error) {
	plans, err := c.GetBillingPlans()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, p := range plans {
		for _, pr := range p.Prices {
			ids = append(ids, fmt.Sprintf("%s\t%s, %s/%s", pr.PriceID, p.Name, formatMinorUnits(pr.Price, pr.Currency), pr.Unit))
		}
	}
	return ids, nil
}

func planTiers() ([]string, error) {
	plans, err := c.GetBillingPlans()
	if err != nil {
		return nil, err
	}
	var tiers []string
	for _, p := range plans {
		tiers = append(tiers, strings.ToLower(p.Name))
	}
	return tiers, nil
}

func userIDs() ([]string, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, u := range users {
		ids = append(ids, fmt.Sprintf("%s\t%s %s", u.ID, u.Name, u.Email))
	}
	return ids, nil
}

var tenantRoles = []string{"admin", "user", "auditor", "network_admin", "billing_admin"}

// tenantGroupCompletion completes <group>:<role> pairs: group names first, then roles after the colon.
func tenantGroupCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	// StringSlice flags: only complete the last comma-separated item.
	prefix := ""
	if i := strings.LastIndex(toComplete, ","); i >= 0 {
		prefix, toComplete = toComplete[:i+1], toComplete[i+1:]
	}
	if group, _, ok := strings.Cut(toComplete, ":"); ok {
		var out []string
		for _, r := range tenantRoles {
			out = append(out, prefix+group+":"+r)
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	}
	groups, err := c.GetGroups()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var out []string
	for _, g := range groups {
		out = append(out, prefix+g.Name+":")
	}
	return out, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

// idpSyncArgs completes <kind> then the IDs of that kind (kinds restricts the first word, e.g. google/azure for "sync idp").
func idpSyncArgs(kinds []string, offset int) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) - offset {
		case 0:
			return kinds, cobra.ShellCompDirectiveNoFileComp
		case 1:
			if c == nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			items, err := c.GetIdPSyncs(args[offset])
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var ids []string
			for _, s := range items {
				ids = append(ids, fmt.Sprintf("%d\t%s", s.ID, idpSyncDetail(s)))
			}
			return ids, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

// ingressPortNames completes allocation names of the peer given with --peer.
func ingressPortNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c == nil || peerFlag == "" || len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	peerID, err := c.ResolvePeerID(peerFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	ports, err := c.GetIngressPorts(peerID)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, p := range ports {
		names = append(names, fmt.Sprintf("%s\t%s", p.Name, p.IngressIP))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func noCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}
