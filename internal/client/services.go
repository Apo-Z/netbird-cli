package client

import (
	"fmt"
)

type Service struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Domain             string              `json:"domain"`
	Mode               string              `json:"mode"`
	ListenPort         int                 `json:"listen_port"`
	PortAutoAssigned   bool                `json:"port_auto_assigned"`
	ProxyCluster       string              `json:"proxy_cluster"`
	Targets            []ServiceTarget     `json:"targets"`
	Enabled            bool                `json:"enabled"`
	Terminated         bool                `json:"terminated"`
	PassHostHeader     bool                `json:"pass_host_header"`
	RewriteRedirects   bool                `json:"rewrite_redirects"`
	Auth               *ServiceAuth        `json:"auth,omitempty"`
	AccessRestrictions *AccessRestrictions `json:"access_restrictions,omitempty"`
	Meta               *ServiceMeta        `json:"meta"`
}

type ServiceTarget struct {
	TargetID   string          `json:"target_id"`
	TargetType string          `json:"target_type"`
	Path       string          `json:"path,omitempty"`
	Protocol   string          `json:"protocol"`
	Host       string          `json:"host,omitempty"`
	Port       int             `json:"port"`
	Enabled    bool            `json:"enabled"`
	Options    *TargetOptions  `json:"options,omitempty"`
}

type TargetOptions struct {
	SkipTLSVerify       *bool              `json:"skip_tls_verify,omitempty"`
	RequestTimeout      string             `json:"request_timeout,omitempty"`
	PathRewrite         string             `json:"path_rewrite,omitempty"`
	CustomHeaders       map[string]string  `json:"custom_headers,omitempty"`
	ProxyProtocol       *bool              `json:"proxy_protocol,omitempty"`
	SessionIdleTimeout  string             `json:"session_idle_timeout,omitempty"`
}

type ServiceAuth struct {
	PasswordAuth *PasswordAuth `json:"password_auth,omitempty"`
	PINAuth      *PINAuth      `json:"pin_auth,omitempty"`
	BearerAuth   *BearerAuth   `json:"bearer_auth,omitempty"`
	LinkAuth     *LinkAuth     `json:"link_auth,omitempty"`
	HeaderAuths  []HeaderAuth  `json:"header_auths,omitempty"`
}

type PasswordAuth struct {
	Enabled  bool   `json:"enabled"`
	Password string `json:"password"`
}

type PINAuth struct {
	Enabled bool   `json:"enabled"`
	PIN     string `json:"pin"`
}

type BearerAuth struct {
	Enabled            bool     `json:"enabled"`
	DistributionGroups []string `json:"distribution_groups,omitempty"`
}

type LinkAuth struct {
	Enabled bool `json:"enabled"`
}

type HeaderAuth struct {
	Enabled bool   `json:"enabled"`
	Header  string `json:"header"`
	Value   string `json:"value"`
}

type AccessRestrictions struct {
	AllowedCIDRs    []string `json:"allowed_cidrs,omitempty"`
	BlockedCIDRs    []string `json:"blocked_cidrs,omitempty"`
	AllowedCountries []string `json:"allowed_countries,omitempty"`
	BlockedCountries []string `json:"blocked_countries,omitempty"`
	CrowdsecMode    string   `json:"crowdsec_mode,omitempty"`
}

type ServiceMeta struct {
	CreatedAt           string `json:"created_at"`
	CertificateIssuedAt string `json:"certificate_issued_at"`
	Status              string `json:"status"`
}

type ProxyCluster struct {
	Address           string `json:"address"`
	ConnectedProxies  int    `json:"connected_proxies"`
}

type CreateServiceRequest struct {
	Name               string              `json:"name"`
	Domain             string              `json:"domain"`
	Mode               string              `json:"mode,omitempty"`
	ListenPort         int                 `json:"listen_port,omitempty"`
	Targets            []ServiceTarget     `json:"targets,omitempty"`
	Enabled            bool                `json:"enabled"`
	PassHostHeader     *bool               `json:"pass_host_header,omitempty"`
	RewriteRedirects   *bool               `json:"rewrite_redirects,omitempty"`
	Auth               *ServiceAuth        `json:"auth,omitempty"`
	AccessRestrictions *AccessRestrictions `json:"access_restrictions,omitempty"`
}

type UpdateServiceRequest struct {
	Name               string              `json:"name"`
	Domain             string              `json:"domain"`
	Mode               string              `json:"mode,omitempty"`
	ListenPort         int                 `json:"listen_port,omitempty"`
	Targets            []ServiceTarget     `json:"targets,omitempty"`
	Enabled            bool                `json:"enabled"`
	PassHostHeader     *bool               `json:"pass_host_header,omitempty"`
	RewriteRedirects   *bool               `json:"rewrite_redirects,omitempty"`
	Auth               *ServiceAuth        `json:"auth,omitempty"`
	AccessRestrictions *AccessRestrictions `json:"access_restrictions,omitempty"`
}

func (c *Client) GetProxyClusters() ([]ProxyCluster, error) {
	resp, err := c.doGet("/api/reverse-proxies/clusters", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetProxyClusters : %w", err)
	}
	return bodyToSlice[ProxyCluster](resp.Body)
}

func (c *Client) GetServices() ([]Service, error) {
	resp, err := c.doGet("/api/reverse-proxies/services", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetServices : %w", err)
	}
	return bodyToSlice[Service](resp.Body)
}

func (c *Client) GetServiceByID(id string) (*Service, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/reverse-proxies/services/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetServiceByID : %w", err)
	}
	return bodyToStructure[Service](resp.Body)
}

func (c *Client) GetServiceByName(name string) (*Service, error) {
	services, err := c.GetServices()
	if err != nil {
		return nil, fmt.Errorf("error lookup service %s : %w", name, err)
	}
	for _, s := range services {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("service %s not found", name)
}

func (c *Client) CreateService(req *CreateServiceRequest) (*Service, error) {
	resp, err := c.doPost("/api/reverse-proxies/services", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateService : %w", err)
	}
	return bodyToStructure[Service](resp.Body)
}

func (c *Client) UpdateService(id string, req *UpdateServiceRequest) (*Service, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/reverse-proxies/services/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateService : %w", err)
	}
	return bodyToStructure[Service](resp.Body)
}

func (c *Client) DeleteService(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/reverse-proxies/services/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteService : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
