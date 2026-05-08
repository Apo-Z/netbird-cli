package client

import (
	"fmt"
)

type Peer struct {
	ID                          string       `json:"id"`
	Name                        string       `json:"name"`
	CreatedAt                   string       `json:"created_at"`
	IP                          string       `json:"ip"`
	ConnectionIP                string       `json:"connection_ip"`
	Connected                   bool         `json:"connected"`
	LastSeen                    string       `json:"last_seen"`
	OS                          string       `json:"os"`
	KernelVersion               string       `json:"kernel_version"`
	GeonameID                   int          `json:"geoname_id"`
	Version                     string       `json:"version"`
	Groups                      []Group      `json:"groups"`
	SSHEnabled                  bool         `json:"ssh_enabled"`
	UserID                      string       `json:"user_id"`
	Hostname                    string       `json:"hostname"`
	UIVersion                   string       `json:"ui_version"`
	DNSLabel                    string       `json:"dns_label"`
	LoginExpirationEnabled      bool         `json:"login_expiration_enabled"`
	LoginExpired                bool         `json:"login_expired"`
	LastLogin                   string       `json:"last_login"`
	InactivityExpirationEnabled bool         `json:"inactivity_expiration_enabled"`
	ApprovalRequired            bool         `json:"approval_required"`
	DisapprovalReason           string       `json:"disapproval_reason"`
	CountryCode                 string       `json:"country_code"`
	CityName                    string       `json:"city_name"`
	SerialNumber                string       `json:"serial_number"`
	ExtraDNSLabels              []string     `json:"extra_dns_labels"`
	Ephemeral                   bool         `json:"ephemeral"`
	LocalFlags                  *LocalFlags  `json:"local_flags"`
	AccessiblePeersCount        int          `json:"accessible_peers_count,omitempty"`
}

type LocalFlags struct {
	RosenpassEnabled      bool `json:"rosenpass_enabled"`
	RosenpassPermissive   bool `json:"rosenpass_permissive"`
	ServerSSHAllowed      bool `json:"server_ssh_allowed"`
	DisableClientRoutes   bool `json:"disable_client_routes"`
	DisableServerRoutes   bool `json:"disable_server_routes"`
	DisableDNS            bool `json:"disable_dns"`
	DisableFirewall       bool `json:"disable_firewall"`
	BlockLANAccess        bool `json:"block_lan_access"`
	BlockInbound          bool `json:"block_inbound"`
	LazyConnectionEnabled bool `json:"lazy_connection_enabled"`
}

type UpdatePeerRequest struct {
	Name                        string `json:"name"`
	SSHEnabled                  bool   `json:"ssh_enabled"`
	LoginExpirationEnabled      bool   `json:"login_expiration_enabled"`
	InactivityExpirationEnabled bool   `json:"inactivity_expiration_enabled"`
	ApprovalRequired            *bool  `json:"approval_required,omitempty"`
	IP                          string `json:"ip,omitempty"`
}

type AccessiblePeer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	DNSLabel    string `json:"dns_label"`
	UserID      string `json:"user_id"`
	OS          string `json:"os"`
	CountryCode string `json:"country_code"`
	CityName    string `json:"city_name"`
	GeonameID   int    `json:"geoname_id"`
	Connected   bool   `json:"connected"`
	LastSeen    string `json:"last_seen"`
}

type CreateTemporaryAccessRequest struct {
	Name     string   `json:"name"`
	WGPubKey string   `json:"wg_pub_key"`
	Rules    []string `json:"rules"`
}

type TemporaryAccessPeer struct {
	Name  string   `json:"name"`
	ID    string   `json:"id"`
	Rules []string `json:"rules"`
}

func (c *Client) GetPeers(name, ip string) ([]Peer, error) {
	params := map[string]string{}
	if name != "" {
		params["name"] = name
	}
	if ip != "" {
		params["ip"] = ip
	}
	resp, err := c.doGet("/api/peers", params)
	if err != nil {
		return nil, fmt.Errorf("error GetPeers : %w", err)
	}
	return bodyToSlice[Peer](resp.Body)
}

func (c *Client) GetPeerByID(id string) (*Peer, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/peers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetPeerByID : %w", err)
	}
	return bodyToStructure[Peer](resp.Body)
}

func (c *Client) GetPeerByName(name string) (*Peer, error) {
	peers, err := c.GetPeers(name, "")
	if err != nil {
		return nil, fmt.Errorf("error lookup peer %s : %w", name, err)
	}
	if len(peers) == 0 {
		return nil, fmt.Errorf("peer %s not found", name)
	}
	c.cacheStore("peer:"+name, peers[0].ID)
	return &peers[0], nil
}

func (c *Client) ResolvePeerID(nameOrID string) (string, error) {
	if id, ok := c.cacheLookup("peer:" + nameOrID); ok {
		return id, nil
	}

	peer, err := c.GetPeerByName(nameOrID)
	if err == nil {
		c.cacheStore("peer:"+peer.Name, peer.ID)
		c.cacheStore("peer:"+peer.ID, peer.ID)
		return peer.ID, nil
	}

	peer, err = c.GetPeerByID(nameOrID)
	if err == nil {
		c.cacheStore("peer:"+peer.Name, peer.ID)
		c.cacheStore("peer:"+peer.ID, peer.ID)
		return peer.ID, nil
	}

	return "", fmt.Errorf("peer %s not found", nameOrID)
}

func (c *Client) UpdatePeer(id string, req *UpdatePeerRequest) (*Peer, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/peers/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdatePeer : %w", err)
	}
	return bodyToStructure[Peer](resp.Body)
}

func (c *Client) DeletePeer(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/peers/%s", id))
	if err != nil {
		return fmt.Errorf("error DeletePeer : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetAccessiblePeers(peerID string) ([]AccessiblePeer, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/peers/%s/accessible-peers", peerID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAccessiblePeers : %w", err)
	}
	return bodyToSlice[AccessiblePeer](resp.Body)
}

func (c *Client) CreateTemporaryAccessPeer(peerID string, req *CreateTemporaryAccessRequest) (*TemporaryAccessPeer, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/peers/%s/temporary-access", peerID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateTemporaryAccessPeer : %w", err)
	}
	return bodyToStructure[TemporaryAccessPeer](resp.Body)
}
