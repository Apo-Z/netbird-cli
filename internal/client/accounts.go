package client

import (
	"fmt"
)

type Account struct {
	ID             string          `json:"id"`
	Settings       *AccountSettings `json:"settings"`
	Domain         string          `json:"domain"`
	DomainCategory string          `json:"domain_category"`
	CreatedAt      string          `json:"created_at"`
	CreatedBy      string          `json:"created_by"`
	Onboarding     *Onboarding     `json:"onboarding,omitempty"`
}

type AccountSettings struct {
	PeerLoginExpirationEnabled       bool     `json:"peer_login_expiration_enabled"`
	PeerLoginExpiration              int      `json:"peer_login_expiration"`
	PeerInactivityExpirationEnabled  bool     `json:"peer_inactivity_expiration_enabled"`
	PeerInactivityExpiration         int      `json:"peer_inactivity_expiration"`
	RegularUsersViewBlocked          bool     `json:"regular_users_view_blocked"`
	GroupsPropagationEnabled         *bool    `json:"groups_propagation_enabled,omitempty"`
	JWTGroupsEnabled                 *bool    `json:"jwt_groups_enabled,omitempty"`
	JWTGroupsClaimName               string   `json:"jwt_groups_claim_name,omitempty"`
	JWTAllowGroups                   []string `json:"jwt_allow_groups,omitempty"`
	RoutingPeerDNSResolutionEnabled  *bool    `json:"routing_peer_dns_resolution_enabled,omitempty"`
	DNSDomain                        string   `json:"dns_domain,omitempty"`
	NetworkRange                     string   `json:"network_range,omitempty"`
	PeerExposeEnabled                bool     `json:"peer_expose_enabled"`
	PeerExposeGroups                 []string `json:"peer_expose_groups"`
	Extra                            *AccountExtra `json:"extra,omitempty"`
	LazyConnectionEnabled            *bool    `json:"lazy_connection_enabled,omitempty"`
	AutoUpdateVersion                string   `json:"auto_update_version,omitempty"`
	AutoUpdateAlways                 *bool    `json:"auto_update_always,omitempty"`
	EmbeddedIDPEnabled               *bool    `json:"embedded_idp_enabled,omitempty"`
	LocalAuthDisabled                *bool    `json:"local_auth_disabled,omitempty"`
}

type AccountExtra struct {
	PeerApprovalEnabled              bool     `json:"peer_approval_enabled"`
	UserApprovalRequired             bool     `json:"user_approval_required"`
	NetworkTrafficLogsEnabled        bool     `json:"network_traffic_logs_enabled"`
	NetworkTrafficLogsGroups         []string `json:"network_traffic_logs_groups"`
	NetworkTrafficPacketCounterEnabled bool   `json:"network_traffic_packet_counter_enabled"`
}

type Onboarding struct {
	SignupFormPending      bool `json:"signup_form_pending"`
	OnboardingFlowPending  bool `json:"onboarding_flow_pending"`
}

type UpdateAccountRequest struct {
	Settings   *AccountSettings `json:"settings"`
	Onboarding *Onboarding     `json:"onboarding,omitempty"`
}

func (c *Client) GetAccounts() ([]Account, error) {
	resp, err := c.doGet("/api/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAccounts : %w", err)
	}
	return bodyToSlice[Account](resp.Body)
}

func (c *Client) UpdateAccount(id string, req *UpdateAccountRequest) (*Account, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/accounts/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateAccount : %w", err)
	}
	return bodyToStructure[Account](resp.Body)
}

func (c *Client) DeleteAccount(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/accounts/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteAccount : %w", err)
	}
	defer resp.Body.Close()
	return nil
}
