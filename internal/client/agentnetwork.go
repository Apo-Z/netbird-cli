package client

import (
	"fmt"
)

// Agent Network: NetBird's gateway for AI agents / LLM traffic.
// Endpoints live under /api/agent-network (and /api/integrations/agent-network for the managed gateway).

type AgentProviderModel struct {
	ID                 string   `json:"id" yaml:"id"`
	InputPer1k         float64  `json:"input_per_1k" yaml:"input_per_1k"`
	OutputPer1k        float64  `json:"output_per_1k" yaml:"output_per_1k"`
	CachedInputPer1k   *float64 `json:"cached_input_per_1k,omitempty" yaml:"cached_input_per_1k,omitempty"`
	CacheReadPer1k     *float64 `json:"cache_read_per_1k,omitempty" yaml:"cache_read_per_1k,omitempty"`
	CacheCreationPer1k *float64 `json:"cache_creation_per_1k,omitempty" yaml:"cache_creation_per_1k,omitempty"`
}

type AgentProvider struct {
	ID                   string               `json:"id"`
	ProviderID           string               `json:"provider_id"`
	Name                 string               `json:"name"`
	UpstreamURL          string               `json:"upstream_url"`
	Models               []AgentProviderModel `json:"models"`
	ExtraValues          map[string]string    `json:"extra_values,omitempty"`
	IdentityHeaderUserID string               `json:"identity_header_user_id"`
	IdentityHeaderGroups string               `json:"identity_header_groups"`
	Enabled              bool                 `json:"enabled"`
	SkipTLSVerification  bool                 `json:"skip_tls_verification"`
	MetadataDisabled     bool                 `json:"metadata_disabled"`
	CreatedAt            string               `json:"created_at"`
	UpdatedAt            string               `json:"updated_at"`
}

type AgentProviderRequest struct {
	ProviderID           string               `json:"provider_id" yaml:"provider_id"`
	Name                 string               `json:"name" yaml:"name"`
	UpstreamURL          string               `json:"upstream_url" yaml:"upstream_url"`
	APIKey               string               `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	Models               []AgentProviderModel `json:"models,omitempty" yaml:"models,omitempty"`
	ExtraValues          map[string]string    `json:"extra_values,omitempty" yaml:"extra_values,omitempty"`
	IdentityHeaderUserID string               `json:"identity_header_user_id,omitempty" yaml:"identity_header_user_id,omitempty"`
	IdentityHeaderGroups string               `json:"identity_header_groups,omitempty" yaml:"identity_header_groups,omitempty"`
	Enabled              *bool                `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	SkipTLSVerification  bool                 `json:"skip_tls_verification" yaml:"skip_tls_verification"`
	MetadataDisabled     bool                 `json:"metadata_disabled" yaml:"metadata_disabled"`
}

type AgentTokenLimit struct {
	Enabled       bool  `json:"enabled" yaml:"enabled"`
	GroupCap      int64 `json:"group_cap" yaml:"group_cap"`
	UserCap       int64 `json:"user_cap" yaml:"user_cap"`
	WindowSeconds int64 `json:"window_seconds" yaml:"window_seconds"`
}

type AgentBudgetLimit struct {
	Enabled       bool    `json:"enabled" yaml:"enabled"`
	GroupCapUSD   float64 `json:"group_cap_usd" yaml:"group_cap_usd"`
	UserCapUSD    float64 `json:"user_cap_usd" yaml:"user_cap_usd"`
	WindowSeconds int64   `json:"window_seconds" yaml:"window_seconds"`
}

type AgentLimits struct {
	TokenLimit  AgentTokenLimit  `json:"token_limit" yaml:"token_limit"`
	BudgetLimit AgentBudgetLimit `json:"budget_limit" yaml:"budget_limit"`
}

type AgentPolicy struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Description            string       `json:"description"`
	Enabled                bool         `json:"enabled"`
	SourceGroups           []string     `json:"source_groups"`
	DestinationProviderIDs []string     `json:"destination_provider_ids"`
	GuardrailIDs           []string     `json:"guardrail_ids"`
	Limits                 *AgentLimits `json:"limits,omitempty"`
	CreatedAt              string       `json:"created_at"`
	UpdatedAt              string       `json:"updated_at"`
}

type AgentPolicyRequest struct {
	Name                   string       `json:"name" yaml:"name"`
	Description            string       `json:"description,omitempty" yaml:"description,omitempty"`
	Enabled                *bool        `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	SourceGroups           []string     `json:"source_groups" yaml:"source_groups"`
	DestinationProviderIDs []string     `json:"destination_provider_ids" yaml:"destination_provider_ids"`
	GuardrailIDs           []string     `json:"guardrail_ids,omitempty" yaml:"guardrail_ids,omitempty"`
	Limits                 *AgentLimits `json:"limits,omitempty" yaml:"limits,omitempty"`
}

type AgentModelAllowlist struct {
	Enabled bool     `json:"enabled" yaml:"enabled"`
	Models  []string `json:"models" yaml:"models"`
}

type AgentPromptCapture struct {
	Enabled   bool `json:"enabled" yaml:"enabled"`
	RedactPII bool `json:"redact_pii" yaml:"redact_pii"`
}

type AgentGuardrailChecks struct {
	ModelAllowlist AgentModelAllowlist `json:"model_allowlist" yaml:"model_allowlist"`
	PromptCapture  AgentPromptCapture  `json:"prompt_capture" yaml:"prompt_capture"`
}

type AgentGuardrail struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Checks      AgentGuardrailChecks `json:"checks"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type AgentGuardrailRequest struct {
	Name        string               `json:"name" yaml:"name"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Checks      AgentGuardrailChecks `json:"checks" yaml:"checks"`
}

type AgentBudgetRule struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Enabled      bool        `json:"enabled"`
	TargetGroups []string    `json:"target_groups"`
	TargetUsers  []string    `json:"target_users"`
	Limits       AgentLimits `json:"limits"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
}

type AgentBudgetRuleRequest struct {
	Name         string      `json:"name" yaml:"name"`
	Enabled      *bool       `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	TargetGroups []string    `json:"target_groups" yaml:"target_groups"`
	TargetUsers  []string    `json:"target_users" yaml:"target_users"`
	Limits       AgentLimits `json:"limits" yaml:"limits"`
}

type AgentSettings struct {
	Endpoint               string `json:"endpoint"`
	ProxyAddress           string `json:"proxy_address"`
	Dedicated              bool   `json:"dedicated"`
	EnableLogCollection    bool   `json:"enable_log_collection"`
	EnablePromptCollection bool   `json:"enable_prompt_collection"`
	RedactPII              bool   `json:"redact_pii"`
	AccessLogRetentionDays int    `json:"access_log_retention_days"`
	CreatedAt              string `json:"created_at,omitempty"`
	UpdatedAt              string `json:"updated_at,omitempty"`
}

type AgentSettingsCreateRequest struct {
	ProxyAddress           string `json:"proxy_address,omitempty"`
	Endpoint               string `json:"endpoint,omitempty"`
	EnableLogCollection    *bool  `json:"enable_log_collection,omitempty"`
	EnablePromptCollection *bool  `json:"enable_prompt_collection,omitempty"`
	RedactPII              *bool  `json:"redact_pii,omitempty"`
	AccessLogRetentionDays *int   `json:"access_log_retention_days,omitempty"`
}

type AgentCatalogModel struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	InputPer1k         float64  `json:"input_per_1k"`
	OutputPer1k        float64  `json:"output_per_1k"`
	CachedInputPer1k   *float64 `json:"cached_input_per_1k,omitempty"`
	CacheReadPer1k     *float64 `json:"cache_read_per_1k,omitempty"`
	CacheCreationPer1k *float64 `json:"cache_creation_per_1k,omitempty"`
	ContextWindow      int      `json:"context_window"`
}

type AgentCatalogProvider struct {
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	DefaultHost        string                   `json:"default_host"`
	AuthHeaderTemplate string                   `json:"auth_header_template"`
	DefaultContentType string                   `json:"default_content_type"`
	BrandColor         string                   `json:"brand_color"`
	Kind               string                   `json:"kind"`
	ExtraHeaders       []map[string]interface{} `json:"extra_headers,omitempty"`
	IdentityInjection  map[string]interface{}   `json:"identity_injection,omitempty"`
	PricingSurfaces    []string                 `json:"pricing_surfaces,omitempty"`
	Models             []AgentCatalogModel      `json:"models"`
}

type AgentModelDiscoveryRequest struct {
	CatalogProviderID string `json:"catalog_provider_id"`
	UpstreamURL       string `json:"upstream_url,omitempty"`
	APIKey            string `json:"api_key,omitempty"`
	ProviderID        string `json:"provider_id,omitempty"`
}

type AgentDiscoveredModel struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label,omitempty"`
	PricingKnown       bool     `json:"pricing_known"`
	InputPer1k         float64  `json:"input_per_1k"`
	OutputPer1k        float64  `json:"output_per_1k"`
	CachedInputPer1k   *float64 `json:"cached_input_per_1k,omitempty"`
	CacheReadPer1k     *float64 `json:"cache_read_per_1k,omitempty"`
	CacheCreationPer1k *float64 `json:"cache_creation_per_1k,omitempty"`
}

type AgentModelDiscoveryResponse struct {
	Models []AgentDiscoveredModel `json:"models"`
}

type AgentConfigProvider struct {
	Name             string   `json:"name"`
	CatalogID        string   `json:"catalog_id"`
	APIFlavor        string   `json:"api_flavor"`
	AllModelsAllowed bool     `json:"all_models_allowed"`
	Models           []string `json:"models"`
}

type AgentConfig struct {
	Configured bool                  `json:"configured"`
	Endpoint   string                `json:"endpoint"`
	Providers  []AgentConfigProvider `json:"providers"`
}

type AgentConsumption struct {
	DimensionKind  string  `json:"dimension_kind"`
	DimensionID    string  `json:"dimension_id"`
	WindowSeconds  int64   `json:"window_seconds"`
	WindowStartUTC string  `json:"window_start_utc"`
	TokensInput    int64   `json:"tokens_input"`
	TokensOutput   int64   `json:"tokens_output"`
	CostUSD        float64 `json:"cost_usd"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}

// AgentUsageTotals holds the token/cost counters shared by usage buckets, access logs and sessions.
type AgentUsageTotals struct {
	InputTokens          int64   `json:"input_tokens"`
	OutputTokens         int64   `json:"output_tokens"`
	TotalTokens          int64   `json:"total_tokens"`
	CachedInputTokens    int64   `json:"cached_input_tokens"`
	CacheCreationTokens  int64   `json:"cache_creation_tokens"`
	CostUSD              float64 `json:"cost_usd"`
	InputCostUSD         float64 `json:"input_cost_usd"`
	CachedInputCostUSD   float64 `json:"cached_input_cost_usd"`
	CacheCreationCostUSD float64 `json:"cache_creation_cost_usd"`
	OutputCostUSD        float64 `json:"output_cost_usd"`
	CacheCostUSD         float64 `json:"cache_cost_usd"`
}

type AgentUsageBucket struct {
	PeriodStart string `json:"period_start"`
	AgentUsageTotals
}

type AgentAccessLog struct {
	ID                 string   `json:"id"`
	ServiceID          string   `json:"service_id"`
	Timestamp          string   `json:"timestamp"`
	UserID             string   `json:"user_id,omitempty"`
	Provider           string   `json:"provider,omitempty"`
	Model              string   `json:"model,omitempty"`
	Decision           string   `json:"decision,omitempty"`
	StatusCode         int      `json:"status_code"`
	DurationMS         int      `json:"duration_ms"`
	SourceIP           string   `json:"source_ip,omitempty"`
	Method             string   `json:"method,omitempty"`
	Host               string   `json:"host,omitempty"`
	Path               string   `json:"path,omitempty"`
	SessionID          string   `json:"session_id,omitempty"`
	ResolvedProviderID string   `json:"resolved_provider_id,omitempty"`
	SelectedPolicyID   string   `json:"selected_policy_id,omitempty"`
	DenyReason         string   `json:"deny_reason,omitempty"`
	Stream             bool     `json:"stream"`
	GroupIDs           []string `json:"group_ids,omitempty"`
	RequestPrompt      string   `json:"request_prompt,omitempty"`
	ResponseCompletion string   `json:"response_completion,omitempty"`
	AgentUsageTotals
}

type AgentAccessLogsResponse struct {
	Data         []AgentAccessLog `json:"data"`
	Page         int              `json:"page"`
	PageSize     int              `json:"page_size"`
	TotalRecords int              `json:"total_records"`
	TotalPages   int              `json:"total_pages"`
}

type AgentAccessLogSession struct {
	SessionID    string           `json:"session_id"`
	UserID       string           `json:"user_id,omitempty"`
	StartedAt    string           `json:"started_at"`
	EndedAt      string           `json:"ended_at"`
	RequestCount int              `json:"request_count"`
	Decision     string           `json:"decision"`
	Providers    []string         `json:"providers,omitempty"`
	Models       []string         `json:"models,omitempty"`
	GroupIDs     []string         `json:"group_ids,omitempty"`
	Entries      []AgentAccessLog `json:"entries,omitempty"`
	AgentUsageTotals
}

type AgentAccessLogSessionsResponse struct {
	Data         []AgentAccessLogSession `json:"data"`
	Page         int                     `json:"page"`
	PageSize     int                     `json:"page_size"`
	TotalRecords int                     `json:"total_records"`
	TotalPages   int                     `json:"total_pages"`
}

type AgentManagedProxy struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region,omitempty"`
	Message  string `json:"message,omitempty"`
}

// --- Providers ---

func (c *Client) GetAgentProviders() ([]AgentProvider, error) {
	resp, err := c.doGet("/api/agent-network/providers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentProviders : %w", err)
	}
	return bodyToSlice[AgentProvider](resp.Body)
}

func (c *Client) GetAgentProviderByID(id string) (*AgentProvider, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/agent-network/providers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentProviderByID : %w", err)
	}
	return bodyToStructure[AgentProvider](resp.Body)
}

func (c *Client) ResolveAgentProviderID(nameOrID string) (string, error) {
	return c.resolveByList("agentprovider", nameOrID, func() (map[string]string, error) {
		items, err := c.GetAgentProviders()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, p := range items {
			m[p.ID] = p.Name
		}
		return m, nil
	})
}

func (c *Client) CreateAgentProvider(req *AgentProviderRequest) (*AgentProvider, error) {
	resp, err := c.doPost("/api/agent-network/providers", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentProvider : %w", err)
	}
	return bodyToStructure[AgentProvider](resp.Body)
}

func (c *Client) UpdateAgentProvider(id string, req *AgentProviderRequest) (*AgentProvider, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/agent-network/providers/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateAgentProvider : %w", err)
	}
	return bodyToStructure[AgentProvider](resp.Body)
}

func (c *Client) DeleteAgentProvider(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/agent-network/providers/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteAgentProvider : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Policies ---

func (c *Client) GetAgentPolicies() ([]AgentPolicy, error) {
	resp, err := c.doGet("/api/agent-network/policies", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentPolicies : %w", err)
	}
	return bodyToSlice[AgentPolicy](resp.Body)
}

func (c *Client) GetAgentPolicyByID(id string) (*AgentPolicy, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/agent-network/policies/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentPolicyByID : %w", err)
	}
	return bodyToStructure[AgentPolicy](resp.Body)
}

func (c *Client) ResolveAgentPolicyID(nameOrID string) (string, error) {
	return c.resolveByList("agentpolicy", nameOrID, func() (map[string]string, error) {
		items, err := c.GetAgentPolicies()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, p := range items {
			m[p.ID] = p.Name
		}
		return m, nil
	})
}

func (c *Client) CreateAgentPolicy(req *AgentPolicyRequest) (*AgentPolicy, error) {
	resp, err := c.doPost("/api/agent-network/policies", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentPolicy : %w", err)
	}
	return bodyToStructure[AgentPolicy](resp.Body)
}

func (c *Client) UpdateAgentPolicy(id string, req *AgentPolicyRequest) (*AgentPolicy, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/agent-network/policies/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateAgentPolicy : %w", err)
	}
	return bodyToStructure[AgentPolicy](resp.Body)
}

func (c *Client) DeleteAgentPolicy(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/agent-network/policies/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteAgentPolicy : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Guardrails ---

func (c *Client) GetAgentGuardrails() ([]AgentGuardrail, error) {
	resp, err := c.doGet("/api/agent-network/guardrails", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentGuardrails : %w", err)
	}
	return bodyToSlice[AgentGuardrail](resp.Body)
}

func (c *Client) GetAgentGuardrailByID(id string) (*AgentGuardrail, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/agent-network/guardrails/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentGuardrailByID : %w", err)
	}
	return bodyToStructure[AgentGuardrail](resp.Body)
}

func (c *Client) ResolveAgentGuardrailID(nameOrID string) (string, error) {
	return c.resolveByList("guardrail", nameOrID, func() (map[string]string, error) {
		items, err := c.GetAgentGuardrails()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, g := range items {
			m[g.ID] = g.Name
		}
		return m, nil
	})
}

func (c *Client) CreateAgentGuardrail(req *AgentGuardrailRequest) (*AgentGuardrail, error) {
	resp, err := c.doPost("/api/agent-network/guardrails", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentGuardrail : %w", err)
	}
	return bodyToStructure[AgentGuardrail](resp.Body)
}

func (c *Client) UpdateAgentGuardrail(id string, req *AgentGuardrailRequest) (*AgentGuardrail, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/agent-network/guardrails/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateAgentGuardrail : %w", err)
	}
	return bodyToStructure[AgentGuardrail](resp.Body)
}

func (c *Client) DeleteAgentGuardrail(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/agent-network/guardrails/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteAgentGuardrail : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Budget rules ---

func (c *Client) GetAgentBudgetRules() ([]AgentBudgetRule, error) {
	resp, err := c.doGet("/api/agent-network/budget-rules", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentBudgetRules : %w", err)
	}
	return bodyToSlice[AgentBudgetRule](resp.Body)
}

func (c *Client) GetAgentBudgetRuleByID(id string) (*AgentBudgetRule, error) {
	resp, err := c.doGet(fmt.Sprintf("/api/agent-network/budget-rules/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentBudgetRuleByID : %w", err)
	}
	return bodyToStructure[AgentBudgetRule](resp.Body)
}

func (c *Client) ResolveAgentBudgetRuleID(nameOrID string) (string, error) {
	return c.resolveByList("budgetrule", nameOrID, func() (map[string]string, error) {
		items, err := c.GetAgentBudgetRules()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(items))
		for _, r := range items {
			m[r.ID] = r.Name
		}
		return m, nil
	})
}

func (c *Client) CreateAgentBudgetRule(req *AgentBudgetRuleRequest) (*AgentBudgetRule, error) {
	resp, err := c.doPost("/api/agent-network/budget-rules", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentBudgetRule : %w", err)
	}
	return bodyToStructure[AgentBudgetRule](resp.Body)
}

func (c *Client) UpdateAgentBudgetRule(id string, req *AgentBudgetRuleRequest) (*AgentBudgetRule, error) {
	resp, err := c.doPut(fmt.Sprintf("/api/agent-network/budget-rules/%s", id), req)
	if err != nil {
		return nil, fmt.Errorf("error UpdateAgentBudgetRule : %w", err)
	}
	return bodyToStructure[AgentBudgetRule](resp.Body)
}

func (c *Client) DeleteAgentBudgetRule(id string) error {
	resp, err := c.doDelete(fmt.Sprintf("/api/agent-network/budget-rules/%s", id))
	if err != nil {
		return fmt.Errorf("error DeleteAgentBudgetRule : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Settings ---

func (c *Client) GetAgentSettings() (*AgentSettings, error) {
	resp, err := c.doGet("/api/agent-network/settings", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentSettings : %w", err)
	}
	return bodyToStructure[AgentSettings](resp.Body)
}

func (c *Client) CreateAgentSettings(req *AgentSettingsCreateRequest) (*AgentSettings, error) {
	resp, err := c.doPost("/api/agent-network/settings", req)
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentSettings : %w", err)
	}
	return bodyToStructure[AgentSettings](resp.Body)
}

func (c *Client) DeleteAgentSettings() error {
	resp, err := c.doDelete("/api/agent-network/settings")
	if err != nil {
		return fmt.Errorf("error DeleteAgentSettings : %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// --- Catalog, discovery, caller config ---

func (c *Client) GetAgentCatalogProviders() ([]AgentCatalogProvider, error) {
	resp, err := c.doGet("/api/agent-network/catalog/providers", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentCatalogProviders : %w", err)
	}
	return bodyToSlice[AgentCatalogProvider](resp.Body)
}

func (c *Client) DiscoverAgentModels(req *AgentModelDiscoveryRequest) (*AgentModelDiscoveryResponse, error) {
	resp, err := c.doPost("/api/agent-network/catalog/providers/models", req)
	if err != nil {
		return nil, fmt.Errorf("error DiscoverAgentModels : %w", err)
	}
	return bodyToStructure[AgentModelDiscoveryResponse](resp.Body)
}

func (c *Client) GetAgentConfig() (*AgentConfig, error) {
	resp, err := c.doGet("/api/agent-network/agent-config", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentConfig : %w", err)
	}
	return bodyToStructure[AgentConfig](resp.Body)
}

// --- Observability ---

func (c *Client) GetAgentConsumption() ([]AgentConsumption, error) {
	resp, err := c.doGet("/api/agent-network/consumption", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentConsumption : %w", err)
	}
	return bodyToSlice[AgentConsumption](resp.Body)
}

func (c *Client) GetAgentUsage(params map[string]string) ([]AgentUsageBucket, error) {
	resp, err := c.doGet("/api/agent-network/usage/overview", params)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentUsage : %w", err)
	}
	return bodyToSlice[AgentUsageBucket](resp.Body)
}

func (c *Client) GetAgentAccessLogs(params map[string]string) (*AgentAccessLogsResponse, error) {
	resp, err := c.doGet("/api/agent-network/access-logs", params)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentAccessLogs : %w", err)
	}
	return bodyToStructure[AgentAccessLogsResponse](resp.Body)
}

func (c *Client) GetAgentAccessLogSessions(params map[string]string) (*AgentAccessLogSessionsResponse, error) {
	resp, err := c.doGet("/api/agent-network/access-log-sessions", params)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentAccessLogSessions : %w", err)
	}
	return bodyToStructure[AgentAccessLogSessionsResponse](resp.Body)
}

// --- Managed gateway (cloud) ---

func (c *Client) GetAgentManagedProxy() (*AgentManagedProxy, error) {
	resp, err := c.doGet("/api/integrations/agent-network/managed-proxy", nil)
	if err != nil {
		return nil, fmt.Errorf("error GetAgentManagedProxy : %w", err)
	}
	return bodyToStructure[AgentManagedProxy](resp.Body)
}

func (c *Client) CreateAgentManagedProxy() (*AgentManagedProxy, error) {
	resp, err := c.doPost("/api/integrations/agent-network/managed-proxy", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("error CreateAgentManagedProxy : %w", err)
	}
	return bodyToStructure[AgentManagedProxy](resp.Body)
}

// resolveByList resolves a name or ID using a single list call (id → name map),
// caching both directions. kind prefixes the cache keys (e.g. "guardrail:").
func (c *Client) resolveByList(kind, nameOrID string, list func() (map[string]string, error)) (string, error) {
	if id, ok := c.cacheLookup(kind + ":" + nameOrID); ok {
		return id, nil
	}
	items, err := list()
	if err != nil {
		return "", err
	}
	found := ""
	c.mu.Lock()
	for id, name := range items {
		c.idToName[id] = name
		c.nameCache[kind+":"+id] = id
		if name != "" {
			c.nameCache[kind+":"+name] = id
		}
		if id == nameOrID || (found == "" && name == nameOrID) {
			found = id
		}
	}
	c.mu.Unlock()
	if found == "" {
		return "", fmt.Errorf("%s %s not found", kind, nameOrID)
	}
	return found, nil
}
