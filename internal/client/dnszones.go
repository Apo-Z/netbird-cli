package client

import (
	"fmt"
)

type DNSZone struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	Domain              string      `json:"domain"`
	Enabled             bool        `json:"enabled"`
	EnableSearchDomain  bool        `json:"enable_search_domain"`
	DistributionGroups  []string    `json:"distribution_groups"`
	Records             []DNSRecord `json:"records,omitempty"`
}

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type CreateDNSZoneRequest struct {
	Name               string   `json:"name"`
	Domain             string   `json:"domain"`
	Enabled            *bool    `json:"enabled,omitempty"`
	EnableSearchDomain bool     `json:"enable_search_domain"`
	DistributionGroups []string `json:"distribution_groups"`
}

type UpdateDNSZoneRequest struct {
	Name               string   `json:"name"`
	Domain             string   `json:"domain"`
	Enabled            *bool    `json:"enabled,omitempty"`
	EnableSearchDomain bool     `json:"enable_search_domain"`
	DistributionGroups []string `json:"distribution_groups"`
}

type CreateDNSRecordRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type UpdateDNSRecordRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

func (c *Client) GetDNSZones() ([]DNSZone, error) {
	resp, err := c.doGet("/api/dns/zones", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetDNSZones : %w", err)
	}
	return bodyToSlice[DNSZone](resp.Body)
}

func (c *Client) GetDNSZoneByID(id string) (*DNSZone, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/dns/zones/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetDNSZoneByID : %w", err)
	}
	return bodyToStructure[DNSZone](resp.Body)
}

func (c *Client) GetDNSZoneByName(name string) (*DNSZone, error) {
	zones, err := c.GetDNSZones()
	if err != nil {
		return nil, fmt.Errorf("error lookup dns-zone %s : %w", name, err)
	}
	for _, z := range zones {
		if z.Name == name {
			return &z, nil
		}
	}
	return nil, fmt.Errorf("DNS zone %s not found", name)
}

func (c *Client) ResolveDNSZoneID(nameOrID string) (string, error) {
	if id, ok := c.cacheLookup("dnszone:" + nameOrID); ok {
		return id, nil
	}

	zone, err := c.GetDNSZoneByName(nameOrID)
	if err == nil {
		c.cacheStore("dnszone:"+zone.Name, zone.ID)
		c.cacheStore("dnszone:"+zone.ID, zone.ID)
		return zone.ID, nil
	}

	zone, err = c.GetDNSZoneByID(nameOrID)
	if err == nil {
		c.cacheStore("dnszone:"+zone.Name, zone.ID)
		c.cacheStore("dnszone:"+zone.ID, zone.ID)
		return zone.ID, nil
	}

	return "", fmt.Errorf("DNS zone %s not found", nameOrID)
}

func (c *Client) CreateDNSZone(req *CreateDNSZoneRequest) (*DNSZone, error) {
	resp, err := c.doPost("/api/dns/zones", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateDNSZone : %w", err)
	}
	return bodyToStructure[DNSZone](resp.Body)
}

func (c *Client) UpdateDNSZone(id string, req *UpdateDNSZoneRequest) (*DNSZone, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/dns/zones/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateDNSZone : %w", err)
	}
	return bodyToStructure[DNSZone](resp.Body)
}

func (c *Client) DeleteDNSZone(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/dns/zones/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteDNSZone : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetDNSRecords(zoneID string) ([]DNSRecord, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/dns/zones/%s/records", zoneID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetDNSRecords : %w", err)
	}
	return bodyToSlice[DNSRecord](resp.Body)
}

func (c *Client) GetDNSRecord(zoneID, recordID string) (*DNSRecord, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/dns/zones/%s/records/%s", zoneID, recordID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetDNSRecord : %w", err)
	}
	return bodyToStructure[DNSRecord](resp.Body)
}

func (c *Client) CreateDNSRecord(zoneID string, req *CreateDNSRecordRequest) (*DNSRecord, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/dns/zones/%s/records", zoneID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateDNSRecord : %w", err)
	}
	return bodyToStructure[DNSRecord](resp.Body)
}

func (c *Client) UpdateDNSRecord(zoneID, recordID string, req *UpdateDNSRecordRequest) (*DNSRecord, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/dns/zones/%s/records/%s", zoneID, recordID), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateDNSRecord : %w", err)
	}
	return bodyToStructure[DNSRecord](resp.Body)
}

func (c *Client) DeleteDNSRecord(zoneID, recordID string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/dns/zones/%s/records/%s", zoneID, recordID))
	if err != nil {
		return fmt.Errorf("error DeleteDNSRecord : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
