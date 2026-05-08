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
