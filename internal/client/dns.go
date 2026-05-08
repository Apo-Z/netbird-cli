package client

import (
	"fmt"
)

type NameserverGroup struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Nameservers          []Nameserver `json:"nameservers"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Primary              bool         `json:"primary"`
	Domains              []string     `json:"domains"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

type Nameserver struct {
	IP     string `json:"ip"`
	NSType string `json:"ns_type"`
	Port   int    `json:"port"`
}

type CreateNameserverGroupRequest struct {
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Nameservers          []Nameserver `json:"nameservers"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Primary              bool         `json:"primary"`
	Domains              []string     `json:"domains"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

type UpdateNameserverGroupRequest struct {
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Nameservers          []Nameserver `json:"nameservers"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Primary              bool         `json:"primary"`
	Domains              []string     `json:"domains"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

type DNSSettings struct {
	DisabledManagementGroups []string `json:"disabled_management_groups"`
}

type UpdateDNSSettingsRequest struct {
	DisabledManagementGroups []string `json:"disabled_management_groups"`
}

func (c *Client) GetNameserverGroups() ([]NameserverGroup, error) {
	resp, err := c.doGet("/api/dns/nameservers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNameserverGroups : %w", err)
	}
	return bodyToSlice[NameserverGroup](resp.Body)
}

func (c *Client) GetNameserverGroupByID(id string) (*NameserverGroup, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/dns/nameservers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNameserverGroupByID : %w", err)
	}
	return bodyToStructure[NameserverGroup](resp.Body)
}

func (c *Client) GetNameserverGroupByName(name string) (*NameserverGroup, error) {
	groups, err := c.GetNameserverGroups()
	if err != nil {
		return nil, fmt.Errorf("error lookup nameserver-group %s : %w", name, err)
	}
	for _, nsg := range groups {
		if nsg.Name == name {
			return &nsg, nil
		}
	}
	return nil, fmt.Errorf("nameserver group %s not found", name)
}

func (c *Client) CreateNameserverGroup(req *CreateNameserverGroupRequest) (*NameserverGroup, error) {
	resp, err := c.doPost("/api/dns/nameservers", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateNameserverGroup : %w", err)
	}
	return bodyToStructure[NameserverGroup](resp.Body)
}

func (c *Client) UpdateNameserverGroup(id string, req *UpdateNameserverGroupRequest) (*NameserverGroup, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/dns/nameservers/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateNameserverGroup : %w", err)
	}
	return bodyToStructure[NameserverGroup](resp.Body)
}

func (c *Client) DeleteNameserverGroup(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/dns/nameservers/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteNameserverGroup : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetDNSSettings() (*DNSSettings, error) {
	resp, err := c.doGet("/api/dns/settings", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetDNSSettings : %w", err)
	}
	return bodyToStructure[DNSSettings](resp.Body)
}

func (c *Client) UpdateDNSSettings(req *UpdateDNSSettingsRequest) (*DNSSettings, error) {
	resp, err := c.doPut("/api/dns/settings", req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateDNSSettings : %w", err)
	}
	return bodyToStructure[DNSSettings](resp.Body)
}
