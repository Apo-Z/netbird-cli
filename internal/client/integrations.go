package client

import (
	"fmt"
	"io"
)

// --- Event streaming (/api/event-streaming): export audit events to Datadog, S3, Firehose, HTTP ---

type EventStream struct {
	ID        int               `json:"id"`
	AccountID string            `json:"account_id,omitempty"`
	Platform  string            `json:"platform"`
	Enabled   bool              `json:"enabled"`
	Config    map[string]string `json:"config,omitempty"`
	CreatedAt string            `json:"created_at,omitempty"`
	UpdatedAt string            `json:"updated_at,omitempty"`
}

type EventStreamRequest struct {
	Platform string            `json:"platform"`
	Config   map[string]string `json:"config"`
	Enabled  bool              `json:"enabled"`
}

func (c *Client) GetEventStreams() ([]EventStream, error) {
	resp, err := c.doGet("/api/event-streaming", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetEventStreams : %w", err)
	}
	return bodyToSlice[EventStream](resp.Body)
}

func (c *Client) GetEventStream(id string) (*EventStream, error) {
	resp, err := c.doGet("/api/event-streaming/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("error GetEventStream : %w", err)
	}
	return bodyToStructure[EventStream](resp.Body)
}

// ResolveEventStreamID accepts a numeric ID or a platform name (when a single stream uses it).
func (c *Client) ResolveEventStreamID(idOrPlatform string) (string, error) {
	streams, err := c.GetEventStreams()
	if err != nil {
		return "", err
	}
	var matches []string
	for _, s := range streams {
		id := fmt.Sprintf("%d", s.ID)
		if id == idOrPlatform {
			return id, nil
		}
		if s.Platform == idOrPlatform {
			matches = append(matches, id)
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("event stream %s not found", idOrPlatform)
	}
	return "", fmt.Errorf("several %s event streams exist (ids %v), use the ID", idOrPlatform, matches)
}

func (c *Client) CreateEventStream(req *EventStreamRequest) (*EventStream, error) {
	resp, err := c.doPost("/api/event-streaming", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateEventStream : %w", err)
	}
	return bodyToStructure[EventStream](resp.Body)
}

func (c *Client) DeleteEventStream(id string) error {
	return c.deleteAndClose("/api/event-streaming/"+id, "DeleteEventStream")
}

// --- Notification channels (/api/integrations/notifications) ---

type NotificationChannel struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Target     map[string]interface{} `json:"target,omitempty"`
	EventTypes []string               `json:"event_types"`
	Enabled    bool                   `json:"enabled"`
}

type NotificationChannelRequest struct {
	Type       string                 `json:"type"`
	Target     map[string]interface{} `json:"target,omitempty"`
	EventTypes []string               `json:"event_types"`
	Enabled    bool                   `json:"enabled"`
}

func (c *Client) GetNotificationChannels() ([]NotificationChannel, error) {
	resp, err := c.doGet("/api/integrations/notifications/channels", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNotificationChannels : %w", err)
	}
	return bodyToSlice[NotificationChannel](resp.Body)
}

func (c *Client) GetNotificationChannel(id string) (*NotificationChannel, error) {
	resp, err := c.doGet("/api/integrations/notifications/channels/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNotificationChannel : %w", err)
	}
	return bodyToStructure[NotificationChannel](resp.Body)
}

func (c *Client) CreateNotificationChannel(req *NotificationChannelRequest) (*NotificationChannel, error) {
	resp, err := c.doPost("/api/integrations/notifications/channels", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateNotificationChannel : %w", err)
	}
	return bodyToStructure[NotificationChannel](resp.Body)
}

func (c *Client) DeleteNotificationChannel(id string) error {
	return c.deleteAndClose("/api/integrations/notifications/channels/"+id, "DeleteNotificationChannel")
}

// GetNotificationTypes returns event type code → description.
func (c *Client) GetNotificationTypes() (map[string]string, error) {
	resp, err := c.doGet("/api/integrations/notifications/types", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetNotificationTypes : %w", err)
	}
	m, err := bodyToStructure[map[string]string](resp.Body)
	if err != nil {
		return nil, err
	}
	return *m, nil
}

// --- Reverse proxy: domains, proxy tokens, clusters ---

type ReverseProxyDomain struct {
	ID                  string `json:"id"`
	Domain              string `json:"domain"`
	Type                string `json:"type"`
	Validated           bool   `json:"validated"`
	TargetCluster       string `json:"target_cluster,omitempty"`
	SupportsCustomPorts bool   `json:"supports_custom_ports"`
	RequireSubdomain    bool   `json:"require_subdomain"`
	SupportsCrowdsec    bool   `json:"supports_crowdsec"`
	SupportsPrivate     bool   `json:"supports_private"`
}

type ReverseProxyDomainRequest struct {
	Domain        string `json:"domain"`
	TargetCluster string `json:"target_cluster"`
}

func (c *Client) GetReverseProxyDomains() ([]ReverseProxyDomain, error) {
	resp, err := c.doGet("/api/reverse-proxies/domains", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetReverseProxyDomains : %w", err)
	}
	return bodyToSlice[ReverseProxyDomain](resp.Body)
}

func (c *Client) ResolveReverseProxyDomainID(nameOrID string) (string, error) {
	return c.resolveByList("domain", nameOrID, func() (map[string]string, error) {
		items, err := c.GetReverseProxyDomains()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, d := range items {
			m[d.ID] = d.Domain
		}
		return m, nil
	})
}

func (c *Client) CreateReverseProxyDomain(req *ReverseProxyDomainRequest) ([]byte, error) {
	return c.PostRaw("/api/reverse-proxies/domains", req)
}

func (c *Client) DeleteReverseProxyDomain(id string) error {
	return c.deleteAndClose("/api/reverse-proxies/domains/"+id, "DeleteReverseProxyDomain")
}

// ValidateReverseProxyDomain triggers (asynchronous) validation of a custom domain.
func (c *Client) ValidateReverseProxyDomain(id string) error {
	resp, err := c.doGet(fmt.Sprintf("/api/reverse-proxies/domains/%s/validate", id), nil)
	if err != nil {
		return fmt.Errorf("error ValidateReverseProxyDomain : %w", err)
	}
	resp.Body.Close()
	return nil
}

type ProxyToken struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Revoked    bool   `json:"revoked"`
	ExpiresAt  string `json:"expires_at,omitempty"`
	LastUsed   string `json:"last_used,omitempty"`
	CreatedAt  string `json:"created_at"`
	PlainToken string `json:"plain_token,omitempty"`
}

type ProxyTokenRequest struct {
	Name      string `json:"name"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

func (c *Client) GetProxyTokens() ([]ProxyToken, error) {
	resp, err := c.doGet("/api/reverse-proxies/proxy-tokens", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetProxyTokens : %w", err)
	}
	return bodyToSlice[ProxyToken](resp.Body)
}

func (c *Client) ResolveProxyTokenID(nameOrID string) (string, error) {
	return c.resolveByList("proxytoken", nameOrID, func() (map[string]string, error) {
		items, err := c.GetProxyTokens()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, t := range items {
			if !t.Revoked {
				m[t.ID] = t.Name
			}
		}
		return m, nil
	})
}

func (c *Client) CreateProxyToken(req *ProxyTokenRequest) (*ProxyToken, error) {
	resp, err := c.doPost("/api/reverse-proxies/proxy-tokens", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateProxyToken : %w", err)
	}
	return bodyToStructure[ProxyToken](resp.Body)
}

func (c *Client) RevokeProxyToken(id string) error {
	return c.deleteAndClose("/api/reverse-proxies/proxy-tokens/"+id, "RevokeProxyToken")
}

func (c *Client) DeleteProxyCluster(address string) error {
	return c.deleteAndClose("/api/reverse-proxies/clusters/"+address, "DeleteProxyCluster")
}

// --- Ingress peers & port allocations (cloud-only) ---

type IngressPeer struct {
	ID             string `json:"id"`
	PeerID         string `json:"peer_id"`
	IngressIP      string `json:"ingress_ip"`
	Region         string `json:"region"`
	Enabled        bool   `json:"enabled"`
	Connected      bool   `json:"connected"`
	Fallback       bool   `json:"fallback"`
	AvailablePorts struct {
		TCP int `json:"tcp"`
		UDP int `json:"udp"`
	} `json:"available_ports"`
}

type IngressPeerRequest struct {
	PeerID   string `json:"peer_id,omitempty"`
	Enabled  bool   `json:"enabled"`
	Fallback bool   `json:"fallback"`
}

func (c *Client) GetIngressPeers() ([]IngressPeer, error) {
	resp, err := c.doGet("/api/ingress/peers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIngressPeers : %w", err)
	}
	return bodyToSlice[IngressPeer](resp.Body)
}

func (c *Client) GetIngressPeer(id string) (*IngressPeer, error) {
	resp, err := c.doGet("/api/ingress/peers/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIngressPeer : %w", err)
	}
	return bodyToStructure[IngressPeer](resp.Body)
}

// ResolveIngressPeerID accepts an ingress peer ID, or the name/ID of the underlying peer.
func (c *Client) ResolveIngressPeerID(nameOrID string) (string, error) {
	peers, err := c.GetIngressPeers()
	if err != nil {
		return "", err
	}
	peerID, _ := c.ResolvePeerID(nameOrID)
	for _, p := range peers {
		if p.ID == nameOrID || p.PeerID == nameOrID || (peerID != "" && p.PeerID == peerID) {
			return p.ID, nil
		}
	}
	return "", fmt.Errorf("ingress peer %s not found", nameOrID)
}

func (c *Client) CreateIngressPeer(req *IngressPeerRequest) (*IngressPeer, error) {
	resp, err := c.doPost("/api/ingress/peers", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateIngressPeer : %w", err)
	}
	return bodyToStructure[IngressPeer](resp.Body)
}

func (c *Client) DeleteIngressPeer(id string) error {
	return c.deleteAndClose("/api/ingress/peers/"+id, "DeleteIngressPeer")
}

type IngressPortMapping struct {
	TranslatedStart int    `json:"translated_start"`
	TranslatedEnd   int    `json:"translated_end"`
	IngressStart    int    `json:"ingress_start"`
	IngressEnd      int    `json:"ingress_end"`
	Protocol        string `json:"protocol"`
}

type IngressPortAllocation struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	IngressPeerID     string               `json:"ingress_peer_id"`
	IngressIP         string               `json:"ingress_ip"`
	Region            string               `json:"region"`
	Enabled           bool                 `json:"enabled"`
	PortRangeMappings []IngressPortMapping `json:"port_range_mappings"`
}

type IngressPortRange struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Protocol string `json:"protocol"`
}

type IngressDirectPort struct {
	Count    int    `json:"count"`
	Protocol string `json:"protocol"`
}

type IngressPortAllocationRequest struct {
	Name       string             `json:"name"`
	Enabled    bool               `json:"enabled"`
	PortRanges []IngressPortRange `json:"port_ranges,omitempty"`
	DirectPort *IngressDirectPort `json:"direct_port,omitempty"`
}

func (c *Client) GetIngressPorts(peerID string) ([]IngressPortAllocation, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/peers/%s/ingress/ports", peerID), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIngressPorts : %w", err)
	}
	return bodyToSlice[IngressPortAllocation](resp.Body)
}

func (c *Client) ResolveIngressPortID(peerID, nameOrID string) (string, error) {
	ports, err := c.GetIngressPorts(peerID)
	if err != nil {
		return "", err
	}
	for _, p := range ports {
		if p.ID == nameOrID || p.Name == nameOrID {
			return p.ID, nil
		}
	}
	return "", fmt.Errorf("port allocation %s not found on peer %s", nameOrID, peerID)
}

func (c *Client) CreateIngressPort(peerID string, req *IngressPortAllocationRequest) (*IngressPortAllocation, error) {
	resp, err := c.doPost(fmt.Sprintf("/api/peers/%s/ingress/ports", peerID), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateIngressPort : %w", err)
	}
	return bodyToStructure[IngressPortAllocation](resp.Body)
}

func (c *Client) DeleteIngressPort(peerID, id string) error {
	return c.deleteAndClose(fmt.Sprintf("/api/peers/%s/ingress/ports/%s", peerID, id), "DeleteIngressPort")
}

// --- EDR integrations (/api/integrations/edr/{vendor}) and peer compliance bypass ---

// EDRVendors lists the supported EDR integration path segments.
var EDRVendors = []string{"intune", "sentinelone", "falcon", "huntress", "fleetdm"}

// EDRIntegration is the union of the per-vendor EDR responses.
type EDRIntegration struct {
	Vendor             string                 `json:"vendor,omitempty"`
	ID                 int                    `json:"id"`
	Enabled            bool                   `json:"enabled"`
	Groups             []PeerRef              `json:"groups"`
	LastSyncedAt       string                 `json:"last_synced_at,omitempty"`
	LastSyncedInterval int                    `json:"last_synced_interval,omitempty"`
	ClientID           string                 `json:"client_id,omitempty"`
	TenantID           string                 `json:"tenant_id,omitempty"`
	CloudID            string                 `json:"cloud_id,omitempty"`
	ZTAScoreThreshold  *int                   `json:"zta_score_threshold,omitempty"`
	APIURL             string                 `json:"api_url,omitempty"`
	MatchAttributes    map[string]interface{} `json:"match_attributes,omitempty"`
	AccountID          string                 `json:"account_id,omitempty"`
	CreatedBy          string                 `json:"created_by,omitempty"`
	CreatedAt          string                 `json:"created_at,omitempty"`
	UpdatedAt          string                 `json:"updated_at,omitempty"`
}

func edrPath(vendor string) string { return "/api/integrations/edr/" + vendor }

func (c *Client) GetEDRIntegration(vendor string) (*EDRIntegration, error) {
	resp, err := c.doGet(edrPath(vendor), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetEDRIntegration %s : %w", vendor, err)
	}
	e, err := bodyToStructure[EDRIntegration](resp.Body)
	if err == nil {
		e.Vendor = vendor
	}
	return e, err
}

func (c *Client) CreateEDRIntegration(vendor string, req map[string]interface{}) (*EDRIntegration, error) {
	resp, err := c.doPost(edrPath(vendor), req)
	if err != nil {
		return nil, fmt.Errorf("error CreateEDRIntegration %s : %w", vendor, err)
	}
	e, err := bodyToStructure[EDRIntegration](resp.Body)
	if err == nil {
		e.Vendor = vendor
	}
	return e, err
}

func (c *Client) DeleteEDRIntegration(vendor string) error {
	return c.deleteAndClose(edrPath(vendor), "DeleteEDRIntegration "+vendor)
}

type EDRBypass struct {
	PeerID string `json:"peer_id"`
}

func (c *Client) GetEDRBypassedPeers() ([]EDRBypass, error) {
	resp, err := c.doGet("/api/peers/edr/bypassed", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetEDRBypassedPeers : %w", err)
	}
	return bodyToSlice[EDRBypass](resp.Body)
}

func (c *Client) BypassPeerEDR(peerID string) error {
	resp, err := c.doRequest("POST", fmt.Sprintf("/api/peers/%s/edr/bypass", peerID), nil)
	if err != nil {
		return fmt.Errorf("error BypassPeerEDR : %w", err)
	}
	resp.Body.Close()
	return nil
}

func (c *Client) RevokePeerEDRBypass(peerID string) error {
	return c.deleteAndClose(fmt.Sprintf("/api/peers/%s/edr/bypass", peerID), "RevokePeerEDRBypass")
}

// --- IdP sync integrations (Google Workspace, Azure AD / Entra, Okta SCIM, generic SCIM) ---

// IdPSyncKinds maps CLI kind → API path segment.
var IdPSyncKinds = map[string]string{
	"google": "google-idp",
	"azure":  "azure-idp",
	"okta":   "okta-scim-idp",
	"scim":   "scim-idp",
}

type IdPSync struct {
	Kind              string   `json:"kind,omitempty"`
	ID                int      `json:"id"`
	Enabled           bool     `json:"enabled"`
	Provider          string   `json:"provider,omitempty"`
	Prefix            string   `json:"prefix,omitempty"`
	CustomerID        string   `json:"customer_id,omitempty"`
	ClientID          string   `json:"client_id,omitempty"`
	TenantID          string   `json:"tenant_id,omitempty"`
	Host              string   `json:"host,omitempty"`
	SyncInterval      int      `json:"sync_interval,omitempty"`
	LastSyncedAt      string   `json:"last_synced_at,omitempty"`
	GroupPrefixes     []string `json:"group_prefixes,omitempty"`
	UserGroupPrefixes []string `json:"user_group_prefixes,omitempty"`
	ConnectorID       string   `json:"connector_id,omitempty"`
	AuthToken         string   `json:"auth_token,omitempty"`
}

type IdPSyncLog struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func idpSyncPath(kind string) (string, error) {
	seg, ok := IdPSyncKinds[kind]
	if !ok {
		return "", fmt.Errorf("unknown IdP sync type %q (google, azure, okta, scim)", kind)
	}
	return "/api/integrations/" + seg, nil
}

func (c *Client) GetIdPSyncs(kind string) ([]IdPSync, error) {
	path, err := idpSyncPath(kind)
	if err != nil {
		return nil, err
	}
	resp, err := c.doGet(path, nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIdPSyncs %s : %w", kind, err)
	}
	items, err := bodyToSlice[IdPSync](resp.Body)
	for i := range items {
		items[i].Kind = kind
	}
	return items, err
}

func (c *Client) CreateIdPSync(kind string, req map[string]interface{}) (*IdPSync, error) {
	path, err := idpSyncPath(kind)
	if err != nil {
		return nil, err
	}
	resp, err := c.doPost(path, req)
	if err != nil {
		return nil, fmt.Errorf("error CreateIdPSync %s : %w", kind, err)
	}
	s, err := bodyToStructure[IdPSync](resp.Body)
	if err == nil {
		s.Kind = kind
	}
	return s, err
}

func (c *Client) DeleteIdPSync(kind, id string) error {
	path, err := idpSyncPath(kind)
	if err != nil {
		return err
	}
	return c.deleteAndClose(path+"/"+id, "DeleteIdPSync "+kind)
}

// SyncIdP triggers an immediate sync (google, azure).
func (c *Client) SyncIdP(kind, id string) (string, error) {
	path, err := idpSyncPath(kind)
	if err != nil {
		return "", err
	}
	resp, err := c.doRequest("POST", path+"/"+id+"/sync", nil)
	if err != nil {
		return "", fmt.Errorf("error SyncIdP %s : %w", kind, err)
	}
	r, err := bodyToStructure[struct {
		Result string `json:"result"`
	}](resp.Body)
	if err != nil {
		return "", err
	}
	return r.Result, nil
}

// RegenerateIdPSyncToken issues a new SCIM token (okta, scim).
func (c *Client) RegenerateIdPSyncToken(kind, id string) (string, error) {
	path, err := idpSyncPath(kind)
	if err != nil {
		return "", err
	}
	resp, err := c.doRequest("POST", path+"/"+id+"/token", nil)
	if err != nil {
		return "", fmt.Errorf("error RegenerateIdPSyncToken %s : %w", kind, err)
	}
	r, err := bodyToStructure[struct {
		AuthToken string `json:"auth_token"`
	}](resp.Body)
	if err != nil {
		return "", err
	}
	return r.AuthToken, nil
}

func (c *Client) GetIdPSyncLogs(kind, id string) ([]IdPSyncLog, error) {
	path, err := idpSyncPath(kind)
	if err != nil {
		return nil, err
	}
	resp, err := c.doGet(path+"/"+id+"/logs", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetIdPSyncLogs %s : %w", kind, err)
	}
	return bodyToSlice[IdPSyncLog](resp.Body)
}

// --- MSP tenants (/api/integrations/msp/tenants) ---

type TenantGroup struct {
	ID   string `json:"id" yaml:"id"`
	Role string `json:"role" yaml:"role"`
}

type Tenant struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Domain       string        `json:"domain"`
	Status       string        `json:"status"`
	Groups       []TenantGroup `json:"groups"`
	DNSChallenge string        `json:"dns_challenge,omitempty"`
	ActivatedAt  string        `json:"activated_at,omitempty"`
	InvitedAt    string        `json:"invited_at,omitempty"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
}

type TenantRequest struct {
	Name   string        `json:"name"`
	Domain string        `json:"domain,omitempty"`
	Groups []TenantGroup `json:"groups"`
}

const tenantsPath = "/api/integrations/msp/tenants"

func (c *Client) GetTenants() ([]Tenant, error) {
	resp, err := c.doGet(tenantsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("error GetTenants : %w", err)
	}
	return bodyToSlice[Tenant](resp.Body)
}

// ResolveTenantID accepts a tenant ID, name or domain.
func (c *Client) ResolveTenantID(nameOrID string) (string, error) {
	tenants, err := c.GetTenants()
	if err != nil {
		return "", err
	}
	for _, t := range tenants {
		if t.ID == nameOrID || t.Name == nameOrID || t.Domain == nameOrID {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("tenant %s not found", nameOrID)
}

func (c *Client) CreateTenant(req *TenantRequest) (*Tenant, error) {
	resp, err := c.doPost(tenantsPath, req)
	if err != nil {
		return nil, fmt.Errorf("error CreateTenant : %w", err)
	}
	return bodyToStructure[Tenant](resp.Body)
}

func (c *Client) UpdateTenant(id string, req *TenantRequest) (*Tenant, error) {
	resp, err := c.doPut(tenantsPath+"/"+id, req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateTenant : %w", err)
	}
	return bodyToStructure[Tenant](resp.Body)
}

// TenantAction POSTs/PUTs to /tenants/{id}/{action} (unlink, dns, subscription, invite).
func (c *Client) TenantAction(method, id, action string, body interface{}) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", tenantsPath, id, action)
	if body == nil {
		resp, err := c.doRequest(method, path, nil)
		if err != nil {
			return nil, err
		}
		return bodyToBytes(resp.Body)
	}
	if method == "PUT" {
		return c.PutRaw(path, body)
	}
	return c.PostRaw(path, body)
}

// --- Billing (/api/integrations/billing, cloud) ---

type BillingUsage struct {
	ActiveUsers int `json:"active_users"`
	TotalUsers  int `json:"total_users"`
	ActivePeers int `json:"active_peers"`
	TotalPeers  int `json:"total_peers"`
}

type BillingSubscription struct {
	Active         bool     `json:"active"`
	PlanTier       string   `json:"plan_tier"`
	PriceID        string   `json:"price_id"`
	Price          int      `json:"price"`
	Currency       string   `json:"currency"`
	Provider       string   `json:"provider"`
	RemainingTrial int      `json:"remaining_trial,omitempty"`
	Features       []string `json:"features,omitempty"`
	UpdatedAt      string   `json:"updated_at"`
}

type BillingPrice struct {
	PriceID  string `json:"price_id"`
	Currency string `json:"currency"`
	Price    int    `json:"price"`
	Unit     string `json:"unit"`
}

type BillingPlan struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Free        bool           `json:"free"`
	Features    []string       `json:"features"`
	Prices      []BillingPrice `json:"prices"`
}

type BillingInvoice struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
}

// BillingURL is returned by checkout, portal and invoice PDF endpoints.
type BillingURL struct {
	SessionID string `json:"session_id,omitempty"`
	URL       string `json:"url"`
}

const billingPath = "/api/integrations/billing"

func (c *Client) GetBillingUsage() (*BillingUsage, error) {
	resp, err := c.doGet(billingPath+"/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingUsage : %w", err)
	}
	return bodyToStructure[BillingUsage](resp.Body)
}

func (c *Client) GetBillingSubscription() (*BillingSubscription, error) {
	resp, err := c.doGet(billingPath+"/subscription", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingSubscription : %w", err)
	}
	return bodyToStructure[BillingSubscription](resp.Body)
}

func (c *Client) ChangeBillingSubscription(priceID, planTier string) error {
	body := map[string]string{}
	if priceID != "" {
		body["priceID"] = priceID
	}
	if planTier != "" {
		body["plan_tier"] = planTier
	}
	_, err := c.PutRaw(billingPath+"/subscription", body)
	return err
}

func (c *Client) GetBillingPlans() ([]BillingPlan, error) {
	resp, err := c.doGet(billingPath+"/plans", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingPlans : %w", err)
	}
	return bodyToSlice[BillingPlan](resp.Body)
}

func (c *Client) GetBillingInvoices() ([]BillingInvoice, error) {
	resp, err := c.doGet(billingPath+"/invoices", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingInvoices : %w", err)
	}
	return bodyToSlice[BillingInvoice](resp.Body)
}

func (c *Client) GetBillingInvoicePDF(id string) (*BillingURL, error) {
	resp, err := c.doGet(fmt.Sprintf("%s/invoices/%s/pdf", billingPath, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingInvoicePDF : %w", err)
	}
	return bodyToStructure[BillingURL](resp.Body)
}

func (c *Client) GetBillingInvoiceCSV(id string) ([]byte, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("%s/invoices/%s/csv", billingPath, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetBillingInvoiceCSV : %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) GetBillingPortal(baseURL string) (*BillingURL, error) {
	resp, err := c.doGet(billingPath+"/portal", map[string]string{"baseURL": baseURL})
	if err != nil {
		return nil, fmt.Errorf("error GetBillingPortal : %w", err)
	}
	return bodyToStructure[BillingURL](resp.Body)
}

func (c *Client) CreateBillingCheckout(baseURL, priceID string, trial bool) (*BillingURL, error) {
	resp, err := c.doPost(billingPath+"/checkout", map[string]interface{}{"baseURL": baseURL, "priceID": priceID, "enableTrial": trial})
	if err != nil {
		return nil, fmt.Errorf("error CreateBillingCheckout : %w", err)
	}
	return bodyToStructure[BillingURL](resp.Body)
}

func (c *Client) BillingAWSMarketplace(action string, body map[string]string) error {
	_, err := c.PostRaw(billingPath+"/aws/marketplace/"+action, body)
	return err
}

func (c *Client) deleteAndClose(path, op string) error {
	resp, err := c.doDelete(path)
	if err != nil {
		return fmt.Errorf("error %s : %w", op, err)
	}
	resp.Body.Close()
	return nil
}
