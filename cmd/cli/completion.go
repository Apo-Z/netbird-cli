package main

import (
	"fmt"

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
