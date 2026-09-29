package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

// Agent Network: NetBird's gateway for AI agents (LLM providers, policies,
// guardrails, budgets, usage and access logs).

const agentPath = "/api/agent-network"

// --- Providers ---

var agentProvidersGetCmd = &cobra.Command{
	Use:               "agentproviders [name|id]",
	Aliases:           []string{"agentprovider", "ap"},
	Short:             "List or display Agent Network LLM providers",
	ValidArgsFunction: validArgsFunc(agentProviderNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			providers, err := c.GetAgentProviders()
			if err != nil {
				exitErr("agent providers", err)
				return
			}
			printOutput(providers)
			return
		}
		id, err := c.ResolveAgentProviderID(args[0])
		if err != nil {
			exitErr("agent provider", err)
			return
		}
		provider, err := c.GetAgentProviderByID(id)
		if err != nil {
			exitErr("agent provider", err)
			return
		}
		printOutput(provider)
	},
}

var agentProviderCreateCmd = &cobra.Command{
	Use:     "agentprovider",
	Aliases: []string{"ap"},
	Short:   "Create an Agent Network LLM provider",
	Long: `Create an Agent Network LLM provider.

--type is a catalog provider id (see "netbird get agentcatalog"), e.g. openai_api, anthropic_api.
--models restricts the exposed models (catalog prices are used); empty means all catalog models.`,
	Example: `  netbird create agentprovider --name openai --type openai_api --url https://api.openai.com --api-key sk-...
  netbird create agentprovider --edit`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var provider client.AgentProvider
			createAgentWithEditor("agentprovider", agentPath+"/providers", map[string]interface{}{
				"provider_id":           "openai_api",
				"name":                  "",
				"upstream_url":          "https://api.openai.com",
				"api_key":               "",
				"models":                []map[string]interface{}{},
				"enabled":               true,
				"skip_tls_verification": false,
				"metadata_disabled":     false,
			}, &provider)
			return
		}
		if nameFlag == "" || providerTypeFlag == "" || upstreamURLFlag == "" || apiKeyFlag == "" {
			exitErr("--name, --type, --url and --api-key are required (or use --edit)", nil)
			return
		}
		req := &client.AgentProviderRequest{
			ProviderID:          providerTypeFlag,
			Name:                nameFlag,
			UpstreamURL:         upstreamURLFlag,
			APIKey:              apiKeyFlag,
			Enabled:             &enabledFlag,
			SkipTLSVerification: skipTLSFlag,
			MetadataDisabled:    noMetadataFlag,
		}
		if len(modelsFlag) > 0 {
			models, err := catalogModelPrices(providerTypeFlag, modelsFlag)
			if err != nil {
				exitErr("models", err)
				return
			}
			req.Models = models
		}
		if dryRunCheck(redactedProviderRequest(req)) {
			return
		}
		provider, err := c.CreateAgentProvider(req)
		if err != nil {
			exitErr("create agentprovider", err)
			return
		}
		fmt.Println("agent provider created:")
		printOutput(provider)
	},
}

var agentProviderEditCmd = &cobra.Command{
	Use:               "agentprovider <name|id>",
	Aliases:           []string{"ap"},
	Short:             "Edit an Agent Network LLM provider",
	Long:              "Edit an Agent Network LLM provider. The API key is never returned; add an api_key field (or --api-key) to rotate it.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(agentProviderNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentProviderID(args[0])
		if err != nil {
			exitErr("agent provider", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("url") {
			overrides["upstream_url"] = upstreamURLFlag
		}
		if cmd.Flags().Changed("api-key") {
			overrides["api_key"] = apiKeyFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("skip-tls-verify") {
			overrides["skip_tls_verification"] = skipTLSFlag
		}
		if cmd.Flags().Changed("no-metadata") {
			overrides["metadata_disabled"] = noMetadataFlag
		}
		editAgentResource("agentprovider", fmt.Sprintf("%s/providers/%s", agentPath, id), overrides)
	},
}

var agentProviderDeleteCmd = &cobra.Command{
	Use:               "agentprovider <name|id>",
	Aliases:           []string{"ap"},
	Short:             "Delete an Agent Network LLM provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(agentProviderNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentProviderID(args[0])
		if err != nil {
			exitErr("agent provider", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete agent provider %s", args[0])) {
			return
		}
		if err := c.DeleteAgentProvider(id); err != nil {
			exitErr("delete agentprovider", err)
			return
		}
		fmt.Println("agent provider deleted")
	},
}

// --- Policies ---

var agentPoliciesGetCmd = &cobra.Command{
	Use:               "agentpolicies [name|id]",
	Aliases:           []string{"agentpolicy", "apol"},
	Short:             "List or display Agent Network policies",
	ValidArgsFunction: validArgsFunc(agentPolicyNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			policies, err := c.GetAgentPolicies()
			if err != nil {
				exitErr("agent policies", err)
				return
			}
			printOutput(policies)
			return
		}
		id, err := c.ResolveAgentPolicyID(args[0])
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		policy, err := c.GetAgentPolicyByID(id)
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		printOutput(policy)
	},
}

var agentPolicyCreateCmd = &cobra.Command{
	Use:     "agentpolicy",
	Aliases: []string{"apol"},
	Short:   "Create an Agent Network policy (which groups can reach which LLM providers)",
	Example: `  netbird create agentpolicy --name devs-openai --source-groups devs --providers openai
  netbird create agentpolicy --name capped --source-groups devs --providers openai \
      --guardrails no-pii --token-user-cap 200000 --budget-user-cap 5 --window 24h`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var policy client.AgentPolicy
			createAgentWithEditor("agentpolicy", agentPath+"/policies", map[string]interface{}{
				"name":                     "",
				"description":              "",
				"enabled":                  true,
				"source_groups":            []string{},
				"destination_provider_ids": []string{},
				"guardrail_ids":            []string{},
				"limits":                   defaultLimitsTemplate(),
			}, &policy)
			return
		}
		if nameFlag == "" || len(sourceGroupsFlag) == 0 || len(providersFlag) == 0 {
			exitErr("--name, --source-groups and --providers are required (or use --edit)", nil)
			return
		}
		req := &client.AgentPolicyRequest{
			Name:                   nameFlag,
			Description:            descFlag,
			Enabled:                &enabledFlag,
			SourceGroups:           resolveAll("group", sourceGroupsFlag),
			DestinationProviderIDs: resolveAll("agentprovider", providersFlag),
			GuardrailIDs:           resolveAll("guardrail", guardrailsFlag),
			Limits:                 limitsFromFlags(cmd),
		}
		if dryRunCheck(req) {
			return
		}
		policy, err := c.CreateAgentPolicy(req)
		if err != nil {
			exitErr("create agentpolicy", err)
			return
		}
		fmt.Println("agent policy created:")
		printOutput(policy)
	},
}

var agentPolicyEditCmd = &cobra.Command{
	Use:               "agentpolicy <name|id>",
	Aliases:           []string{"apol"},
	Short:             "Edit an Agent Network policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(agentPolicyNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentPolicyID(args[0])
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("source-groups") {
			overrides["source_groups"] = resolveAll("group", sourceGroupsFlag)
		}
		if cmd.Flags().Changed("providers") {
			overrides["destination_provider_ids"] = resolveAll("agentprovider", providersFlag)
		}
		if cmd.Flags().Changed("guardrails") {
			overrides["guardrail_ids"] = resolveAll("guardrail", guardrailsFlag)
		}
		editAgentResource("agentpolicy", fmt.Sprintf("%s/policies/%s", agentPath, id), overrides)
	},
}

var agentPolicyDeleteCmd = &cobra.Command{
	Use:               "agentpolicy <name|id>",
	Aliases:           []string{"apol"},
	Short:             "Delete an Agent Network policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(agentPolicyNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentPolicyID(args[0])
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete agent policy %s", args[0])) {
			return
		}
		if err := c.DeleteAgentPolicy(id); err != nil {
			exitErr("delete agentpolicy", err)
			return
		}
		fmt.Println("agent policy deleted")
	},
}

// --- Guardrails ---

var guardrailsGetCmd = &cobra.Command{
	Use:               "guardrails [name|id]",
	Aliases:           []string{"guardrail", "gr"},
	Short:             "List or display Agent Network guardrails",
	ValidArgsFunction: validArgsFunc(guardrailNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			guardrails, err := c.GetAgentGuardrails()
			if err != nil {
				exitErr("guardrails", err)
				return
			}
			printOutput(guardrails)
			return
		}
		id, err := c.ResolveAgentGuardrailID(args[0])
		if err != nil {
			exitErr("guardrail", err)
			return
		}
		guardrail, err := c.GetAgentGuardrailByID(id)
		if err != nil {
			exitErr("guardrail", err)
			return
		}
		printOutput(guardrail)
	},
}

var guardrailCreateCmd = &cobra.Command{
	Use:     "guardrail",
	Aliases: []string{"gr"},
	Short:   "Create an Agent Network guardrail (model allowlist, prompt capture)",
	Example: `  netbird create guardrail --name cheap-models --models gpt-4o-mini,claude-haiku-4-5
  netbird create guardrail --name audit --prompt-capture --redact-pii`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var guardrail client.AgentGuardrail
			createAgentWithEditor("guardrail", agentPath+"/guardrails", map[string]interface{}{
				"name":        "",
				"description": "",
				"checks": map[string]interface{}{
					"model_allowlist": map[string]interface{}{"enabled": false, "models": []string{}},
					"prompt_capture":  map[string]interface{}{"enabled": false, "redact_pii": false},
				},
			}, &guardrail)
			return
		}
		if nameFlag == "" {
			exitErr("--name is required (or use --edit)", nil)
			return
		}
		req := &client.AgentGuardrailRequest{
			Name:        nameFlag,
			Description: descFlag,
			Checks: client.AgentGuardrailChecks{
				ModelAllowlist: client.AgentModelAllowlist{Enabled: len(modelsFlag) > 0, Models: nonNil(modelsFlag)},
				PromptCapture:  client.AgentPromptCapture{Enabled: promptCaptureFlag, RedactPII: redactPIIFlag},
			},
		}
		if dryRunCheck(req) {
			return
		}
		guardrail, err := c.CreateAgentGuardrail(req)
		if err != nil {
			exitErr("create guardrail", err)
			return
		}
		fmt.Println("guardrail created:")
		printOutput(guardrail)
	},
}

var guardrailEditCmd = &cobra.Command{
	Use:               "guardrail <name|id>",
	Aliases:           []string{"gr"},
	Short:             "Edit an Agent Network guardrail",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(guardrailNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentGuardrailID(args[0])
		if err != nil {
			exitErr("guardrail", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("desc") {
			overrides["description"] = descFlag
		}
		editAgentResource("guardrail", fmt.Sprintf("%s/guardrails/%s", agentPath, id), overrides)
	},
}

var guardrailDeleteCmd = &cobra.Command{
	Use:               "guardrail <name|id>",
	Aliases:           []string{"gr"},
	Short:             "Delete an Agent Network guardrail",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(guardrailNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentGuardrailID(args[0])
		if err != nil {
			exitErr("guardrail", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete guardrail %s", args[0])) {
			return
		}
		if err := c.DeleteAgentGuardrail(id); err != nil {
			exitErr("delete guardrail", err)
			return
		}
		fmt.Println("guardrail deleted")
	},
}

// --- Budget rules ---

var budgetRulesGetCmd = &cobra.Command{
	Use:               "budgetrules [name|id]",
	Aliases:           []string{"budgetrule", "br"},
	Short:             "List or display Agent Network budget rules",
	ValidArgsFunction: validArgsFunc(budgetRuleNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			rules, err := c.GetAgentBudgetRules()
			if err != nil {
				exitErr("budget rules", err)
				return
			}
			printOutput(rules)
			return
		}
		id, err := c.ResolveAgentBudgetRuleID(args[0])
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		rule, err := c.GetAgentBudgetRuleByID(id)
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		printOutput(rule)
	},
}

var budgetRuleCreateCmd = &cobra.Command{
	Use:     "budgetrule",
	Aliases: []string{"br"},
	Short:   "Create an Agent Network budget rule (token / USD caps for groups or users)",
	Long: `Create an Agent Network budget rule.
Without --groups and --users the rule is account-wide. At least one cap is required.`,
	Example: `  netbird create budgetrule --name devs-daily --groups devs --budget-user-cap 10 --window 24h
  netbird create budgetrule --name account-monthly --budget-group-cap 500 --window 720h`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var rule client.AgentBudgetRule
			createAgentWithEditor("budgetrule", agentPath+"/budget-rules", map[string]interface{}{
				"name":          "",
				"enabled":       true,
				"target_groups": []string{},
				"target_users":  []string{},
				"limits":        defaultLimitsTemplate(),
			}, &rule)
			return
		}
		limits := limitsFromFlags(cmd)
		if nameFlag == "" || limits == nil {
			exitErr("--name and at least one --token-*-cap / --budget-*-cap are required (or use --edit)", nil)
			return
		}
		req := &client.AgentBudgetRuleRequest{
			Name:         nameFlag,
			Enabled:      &enabledFlag,
			TargetGroups: nonNil(resolveAll("group", groupsFlag)),
			TargetUsers:  nonNil(resolveAll("user", usersFlag)),
			Limits:       *limits,
		}
		if dryRunCheck(req) {
			return
		}
		rule, err := c.CreateAgentBudgetRule(req)
		if err != nil {
			exitErr("create budgetrule", err)
			return
		}
		fmt.Println("budget rule created:")
		printOutput(rule)
	},
}

var budgetRuleEditCmd = &cobra.Command{
	Use:               "budgetrule <name|id>",
	Aliases:           []string{"br"},
	Short:             "Edit an Agent Network budget rule",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(budgetRuleNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentBudgetRuleID(args[0])
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("groups") {
			overrides["target_groups"] = resolveAll("group", groupsFlag)
		}
		if cmd.Flags().Changed("users") {
			overrides["target_users"] = resolveAll("user", usersFlag)
		}
		editAgentResource("budgetrule", fmt.Sprintf("%s/budget-rules/%s", agentPath, id), overrides)
	},
}

var budgetRuleDeleteCmd = &cobra.Command{
	Use:               "budgetrule <name|id>",
	Aliases:           []string{"br"},
	Short:             "Delete an Agent Network budget rule",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(budgetRuleNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentBudgetRuleID(args[0])
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete budget rule %s", args[0])) {
			return
		}
		if err := c.DeleteAgentBudgetRule(id); err != nil {
			exitErr("delete budgetrule", err)
			return
		}
		fmt.Println("budget rule deleted")
	},
}

// --- Settings ---

var agentSettingsGetCmd = &cobra.Command{
	Use:     "agentsettings",
	Aliases: []string{"as"},
	Short:   "Display Agent Network settings (endpoint, log/prompt collection, retention)",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := c.GetAgentSettings()
		if err != nil {
			exitErr("agent settings", err)
			return
		}
		printOutput(settings)
	},
}

var agentSettingsCreateCmd = &cobra.Command{
	Use:     "agentsettings",
	Aliases: []string{"as"},
	Short:   "Bootstrap Agent Network for the account (allocates the gateway endpoint)",
	Long: `Bootstrap Agent Network for the account.
Exactly one of --proxy-address (endpoint allocated under that proxy cluster) or
--endpoint (self-addressed dedicated endpoint) is required. Both are immutable afterwards.
On NetBird Cloud, "netbird create agentgateway" provisions a managed gateway instead.`,
	Example: `  netbird create agentsettings --proxy-address proxy.example.com
  netbird create agentsettings --endpoint ai.example.com --prompt-collection --redact-pii`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var settings client.AgentSettings
			createAgentWithEditor("agentsettings", agentPath+"/settings", map[string]interface{}{
				"proxy_address":             "",
				"enable_log_collection":     true,
				"enable_prompt_collection":  false,
				"redact_pii":                false,
				"access_log_retention_days": 30,
			}, &settings)
			return
		}
		if (proxyAddrFlag == "") == (endpointFlag == "") {
			exitErr("exactly one of --proxy-address or --endpoint is required (or use --edit)", nil)
			return
		}
		req := &client.AgentSettingsCreateRequest{
			ProxyAddress: proxyAddrFlag,
			Endpoint:     endpointFlag,
		}
		if cmd.Flags().Changed("log-collection") {
			req.EnableLogCollection = &logCollectionFlag
		}
		if cmd.Flags().Changed("prompt-collection") {
			req.EnablePromptCollection = &promptCollectFlag
		}
		if cmd.Flags().Changed("redact-pii") {
			req.RedactPII = &redactPIIFlag
		}
		if cmd.Flags().Changed("retention-days") {
			req.AccessLogRetentionDays = &retentionDaysFlag
		}
		if dryRunCheck(req) {
			return
		}
		settings, err := c.CreateAgentSettings(req)
		if err != nil {
			exitErr("create agentsettings", err)
			return
		}
		fmt.Println("agent network bootstrapped:")
		printOutput(settings)
	},
}

var agentSettingsEditCmd = &cobra.Command{
	Use:     "agentsettings",
	Aliases: []string{"as"},
	Short:   "Edit Agent Network settings (endpoint and proxy_address are immutable)",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("log-collection") {
			overrides["enable_log_collection"] = logCollectionFlag
		}
		if cmd.Flags().Changed("prompt-collection") {
			overrides["enable_prompt_collection"] = promptCollectFlag
		}
		if cmd.Flags().Changed("redact-pii") {
			overrides["redact_pii"] = redactPIIFlag
		}
		if cmd.Flags().Changed("retention-days") {
			overrides["access_log_retention_days"] = retentionDaysFlag
		}
		editAgentResource("agentsettings", agentPath+"/settings", overrides)
	},
}

var agentSettingsDeleteCmd = &cobra.Command{
	Use:     "agentsettings",
	Aliases: []string{"as"},
	Short:   "Delete Agent Network settings (tears down the account's Agent Network setup)",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg("would delete agent network settings") {
			return
		}
		if err := c.DeleteAgentSettings(); err != nil {
			exitErr("delete agentsettings", err)
			return
		}
		fmt.Println("agent network settings deleted")
	},
}

// --- Managed gateway (cloud) ---

var agentGatewayGetCmd = &cobra.Command{
	Use:     "agentgateway",
	Aliases: []string{"agw"},
	Short:   "Display the NetBird-managed Agent Network gateway status (cloud)",
	Run: func(cmd *cobra.Command, args []string) {
		proxy, err := c.GetAgentManagedProxy()
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				fmt.Println("No managed Agent Network gateway (cloud-only; provision it with \"netbird create agentgateway\").")
				return
			}
			exitErr("agent gateway", err)
			return
		}
		printOutput(proxy)
	},
}

var agentGatewayCreateCmd = &cobra.Command{
	Use:     "agentgateway",
	Aliases: []string{"agw"},
	Short:   "Provision a NetBird-managed Agent Network gateway (cloud, idempotent)",
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg("would provision a managed agent network gateway") {
			return
		}
		proxy, err := c.CreateAgentManagedProxy()
		if err != nil {
			exitErr("create agentgateway", err)
			return
		}
		fmt.Println("agent gateway:")
		printOutput(proxy)
	},
}

// --- Catalog, models, caller config ---

type agentCatalogRow struct {
	CatalogID   string `json:"catalog_id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	DefaultHost string `json:"default_host"`
}

type agentModelRow struct {
	Model         string  `json:"model"`
	Label         string  `json:"label"`
	InputPer1k    float64 `json:"input_per_1k"`
	OutputPer1k   float64 `json:"output_per_1k"`
	ContextWindow int     `json:"context_window,omitempty"`
	PricingKnown  *bool   `json:"pricing_known,omitempty"`
}

var agentCatalogGetCmd = &cobra.Command{
	Use:               "agentcatalog [catalog-id]",
	Aliases:           []string{"acat"},
	Short:             "List supported LLM provider types, or the models (and prices) of one",
	ValidArgsFunction: validArgsFunc(agentCatalogIDs),
	Run: func(cmd *cobra.Command, args []string) {
		catalog, err := c.GetAgentCatalogProviders()
		if err != nil {
			exitErr("agent catalog", err)
			return
		}
		if len(args) == 0 {
			if outputFormat != "" {
				printOutput(catalog)
				return
			}
			var rows []agentCatalogRow
			for _, p := range catalog {
				rows = append(rows, agentCatalogRow{p.ID, p.Name, p.Kind, p.DefaultHost})
			}
			printOutput(rows)
			return
		}
		for _, p := range catalog {
			if p.ID != args[0] && p.Name != args[0] {
				continue
			}
			if outputFormat != "" {
				printOutput(p)
				return
			}
			fmt.Printf("%s (%s) — %s\n  default host: %s\n\n", p.Name, p.ID, p.Description, p.DefaultHost)
			var rows []agentModelRow
			for _, m := range p.Models {
				rows = append(rows, agentModelRow{Model: m.ID, Label: m.Label, InputPer1k: m.InputPer1k, OutputPer1k: m.OutputPer1k, ContextWindow: m.ContextWindow})
			}
			printOutput(rows)
			return
		}
		exitErr(fmt.Sprintf("catalog provider %s not found", args[0]), nil)
	},
}

var agentModelsGetCmd = &cobra.Command{
	Use:     "agentmodels",
	Aliases: []string{"am"},
	Short:   "Discover the models an LLM provider credential can reach",
	Long: `Ask the upstream vendor which models a credential can reach, with NetBird's default pricing.
Use --provider to query with an existing provider's stored credential, or --type + --api-key
(and optionally --url) for a provider that has not been created yet.`,
	Example: `  netbird get agentmodels --provider openai
  netbird get agentmodels --type anthropic_api --api-key sk-ant-...`,
	Run: func(cmd *cobra.Command, args []string) {
		req := &client.AgentModelDiscoveryRequest{
			CatalogProviderID: providerTypeFlag,
			UpstreamURL:       upstreamURLFlag,
			APIKey:            apiKeyFlag,
		}
		if providerFlag != "" {
			id, err := c.ResolveAgentProviderID(providerFlag)
			if err != nil {
				exitErr("agent provider", err)
				return
			}
			req.ProviderID = id
			if req.CatalogProviderID == "" {
				if p, err := c.GetAgentProviderByID(id); err == nil {
					req.CatalogProviderID = p.ProviderID
				}
			}
		}
		if req.CatalogProviderID == "" || (req.ProviderID == "" && req.APIKey == "") {
			exitErr("--provider, or --type with --api-key, is required", nil)
			return
		}
		resp, err := c.DiscoverAgentModels(req)
		if err != nil {
			exitErr("discover models", err)
			return
		}
		if outputFormat != "" {
			printOutput(resp.Models)
			return
		}
		var rows []agentModelRow
		for _, m := range resp.Models {
			known := m.PricingKnown
			rows = append(rows, agentModelRow{Model: m.ID, Label: m.Label, InputPer1k: m.InputPer1k, OutputPer1k: m.OutputPer1k, PricingKnown: &known})
		}
		printOutput(rows)
	},
}

var agentConfigGetCmd = &cobra.Command{
	Use:     "agentconfig",
	Aliases: []string{"acfg"},
	Short:   "Show the Agent Network endpoint and providers available to you (the calling user)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := c.GetAgentConfig()
		if err != nil {
			exitErr("agent config", err)
			return
		}
		if outputFormat != "" {
			printOutput(cfg)
			return
		}
		if !cfg.Configured {
			fmt.Println("Agent Network is not set up on this account.")
			return
		}
		fmt.Printf("Endpoint : %s  (reachable over the NetBird tunnel only)\n", cfg.Endpoint)
		if len(cfg.Providers) == 0 {
			fmt.Println("\nNo agent policy grants you access to a provider yet.")
			return
		}
		fmt.Println("\nProviders:")
		for _, p := range cfg.Providers {
			flavor := p.APIFlavor
			if flavor == "" {
				flavor = "by URL path"
			}
			models := strings.Join(p.Models, ", ")
			if p.AllModelsAllowed {
				models = "all models (" + models + ")"
			}
			fmt.Printf("  %s  [%s, api: %s]\n    models: %s\n", bold(p.Name), p.CatalogID, flavor, models)
		}
	},
}

// --- Observability ---

var agentUsageGetCmd = &cobra.Command{
	Use:     "agentusage",
	Aliases: []string{"au"},
	Short:   "Agent Network token and cost usage per day/week/month",
	Example: `  netbird get agentusage --since 30d
  netbird get agentusage --granularity month --user alice@example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		params, ok := agentFilterParams()
		if !ok {
			return
		}
		if granularityFlag != "" {
			params["granularity"] = granularityFlag
		}
		buckets, err := c.GetAgentUsage(params)
		if err != nil {
			exitErr("agent usage", err)
			return
		}
		printOutput(buckets)
		if outputFormat == "" && len(buckets) > 0 {
			var total client.AgentUsageTotals
			for _, b := range buckets {
				total.InputTokens += b.InputTokens
				total.OutputTokens += b.OutputTokens
				total.TotalTokens += b.TotalTokens
				total.CostUSD += b.CostUSD
			}
			fmt.Printf("\ntotal: %d tokens (in %d / out %d), %s\n", total.TotalTokens, total.InputTokens, total.OutputTokens, formatFloat("cost_usd", total.CostUSD))
		}
	},
}

var agentConsumptionGetCmd = &cobra.Command{
	Use:     "agentconsumption",
	Aliases: []string{"acons"},
	Short:   "Current Agent Network consumption counters (per user/group and window, vs. caps)",
	Run: func(cmd *cobra.Command, args []string) {
		rows, err := c.GetAgentConsumption()
		if err != nil {
			exitErr("agent consumption", err)
			return
		}
		printOutput(rows)
	},
}

type agentLogRow struct {
	Timestamp   string  `json:"timestamp"`
	UserID      string  `json:"user_id"`
	Provider    string  `json:"provider"`
	Model       string  `json:"model"`
	Decision    string  `json:"decision"`
	StatusCode  int     `json:"status_code"`
	TotalTokens int64   `json:"total_tokens"`
	CostUSD     float64 `json:"cost_usd"`
	DenyReason  string  `json:"reason"`
}

type agentSessionRow struct {
	StartedAt    string   `json:"started_at"`
	UserID       string   `json:"user_id"`
	RequestCount int      `json:"request_count"`
	Models       []string `json:"models"`
	Decision     string   `json:"decision"`
	TotalTokens  int64    `json:"total_tokens"`
	CostUSD      float64  `json:"cost_usd"`
	SessionID    string   `json:"session_id"`
}

var agentLogsCmd = &cobra.Command{
	Use:     "agent",
	Aliases: []string{"ai", "agentlogs"},
	Short:   "Agent Network (LLM) access logs, optionally grouped by session",
	Example: `  netbird events agent --since 24h
  netbird events agent --decision deny
  netbird events agent --sessions --user alice@example.com
  netbird events agent --session <session-id> -o yaml    # full entries, prompts if captured`,
	Run: func(cmd *cobra.Command, args []string) {
		params, ok := agentFilterParams()
		if !ok {
			return
		}
		params["page"] = strconv.Itoa(pageFlag)
		params["page_size"] = strconv.Itoa(pageSizeFlag)
		if searchFlag != "" {
			params["search"] = searchFlag
		}
		if decisionFlag != "" {
			params["decision"] = decisionFlag
		}
		if sortByFlag != "" {
			params["sort_by"] = sortByFlag
		}
		if sortOrderFlag != "" {
			params["sort_order"] = sortOrderFlag
		}

		if sessionsFlag {
			resp, err := c.GetAgentAccessLogSessions(params)
			if err != nil {
				exitErr("agent sessions", err)
				return
			}
			if outputFormat != "" {
				printOutput(resp)
				return
			}
			var rows []agentSessionRow
			for _, s := range resp.Data {
				rows = append(rows, agentSessionRow{s.StartedAt, s.UserID, s.RequestCount, s.Models, s.Decision, s.TotalTokens, s.CostUSD, s.SessionID})
			}
			printOutput(rows)
			printPageFooter(resp.Page, resp.TotalPages, resp.TotalRecords, "sessions")
			return
		}

		resp, err := c.GetAgentAccessLogs(params)
		if err != nil {
			exitErr("agent access logs", err)
			return
		}
		if outputFormat != "" {
			printOutput(resp)
			return
		}
		var rows []agentLogRow
		for _, l := range resp.Data {
			rows = append(rows, agentLogRow{l.Timestamp, l.UserID, l.Provider, l.Model, l.Decision, l.StatusCode, l.TotalTokens, l.CostUSD, l.DenyReason})
		}
		printOutput(rows)
		printPageFooter(resp.Page, resp.TotalPages, resp.TotalRecords, "requests")
	},
}

// --- Describe ---

var describeAgentPolicyCmd = &cobra.Command{
	Use:               "agentpolicy <name|id>",
	Aliases:           []string{"apol"},
	Short:             "Detail an Agent Network policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(agentPolicyNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentPolicyID(args[0])
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		policy, err := c.GetAgentPolicyByID(id)
		if err != nil {
			exitErr("agent policy", err)
			return
		}
		if outputFormat != "" {
			printOutput(policy)
			return
		}

		fmt.Println("=== Agent policy ===")
		fmt.Printf("Name    : %s\n", policy.Name)
		if policy.Description != "" {
			fmt.Printf("Desc.   : %s\n", policy.Description)
		}
		fmt.Printf("Status  : %s\n", colorBool(policy.Enabled, "enabled", "disabled"))
		fmt.Printf("ID      : %s\n", policy.ID)

		fmt.Println("\nWho (source groups):")
		for _, gid := range policy.SourceGroups {
			fmt.Printf("  %s (%s)\n", nameOrID("group", gid), gid)
		}

		fmt.Println("\nCan reach (providers):")
		providers, _ := c.GetAgentProviders()
		for _, pid := range policy.DestinationProviderIDs {
			line := fmt.Sprintf("  %s (%s)", nameOrID("agentprovider", pid), pid)
			for _, p := range providers {
				if p.ID == pid {
					line += fmt.Sprintf(" — %s, %s, %s", p.ProviderID, p.UpstreamURL, colorBool(p.Enabled, "enabled", "disabled"))
				}
			}
			fmt.Println(line)
		}

		if len(policy.GuardrailIDs) > 0 {
			fmt.Println("\nGuardrails:")
			for _, gid := range policy.GuardrailIDs {
				line := fmt.Sprintf("  %s (%s)", nameOrID("guardrail", gid), gid)
				if g, err := c.GetAgentGuardrailByID(gid); err == nil {
					line += " — " + describeGuardrailChecks(g.Checks)
				}
				fmt.Println(line)
			}
		}

		if policy.Limits != nil {
			fmt.Println("\nLimits:")
			fmt.Printf("  tokens : %s\n", describeTokenLimit(policy.Limits.TokenLimit))
			fmt.Printf("  budget : %s\n", describeBudgetLimit(policy.Limits.BudgetLimit))
		}
	},
}

var describeBudgetRuleCmd = &cobra.Command{
	Use:               "budgetrule <name|id>",
	Aliases:           []string{"br"},
	Short:             "Detail an Agent Network budget rule",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(budgetRuleNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveAgentBudgetRuleID(args[0])
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		rule, err := c.GetAgentBudgetRuleByID(id)
		if err != nil {
			exitErr("budget rule", err)
			return
		}
		if outputFormat != "" {
			printOutput(rule)
			return
		}

		fmt.Println("=== Budget rule ===")
		fmt.Printf("Name    : %s\n", rule.Name)
		fmt.Printf("Status  : %s\n", colorBool(rule.Enabled, "enabled", "disabled"))
		fmt.Printf("ID      : %s\n", rule.ID)

		fmt.Println("\nApplies to:")
		if len(rule.TargetGroups) == 0 && len(rule.TargetUsers) == 0 {
			fmt.Println("  the whole account")
		}
		for _, gid := range rule.TargetGroups {
			fmt.Printf("  group %s (%s)\n", nameOrID("group", gid), gid)
		}
		for _, uid := range rule.TargetUsers {
			fmt.Printf("  user  %s (%s)\n", nameOrID("user", uid), uid)
		}

		fmt.Println("\nLimits:")
		fmt.Printf("  tokens : %s\n", describeTokenLimit(rule.Limits.TokenLimit))
		fmt.Printf("  budget : %s\n", describeBudgetLimit(rule.Limits.BudgetLimit))
	},
}

// --- Helpers ---

// createAgentWithEditor runs the create --edit flow for an Agent Network resource and prints the result into out's type.
func createAgentWithEditor(kind, path string, template map[string]interface{}, out interface{}) {
	result, err := createWithEditor(path, template)
	if err != nil {
		exitErr("create "+kind, err)
		return
	}
	if result == nil {
		return
	}
	if dryRunCheck(result) {
		return
	}
	data, err := c.PostRaw(path, result)
	if err != nil {
		exitErr("create "+kind, err)
		return
	}
	if err := json.Unmarshal(data, out); err != nil {
		exitErr("create "+kind, err)
		return
	}
	fmt.Printf("%s created:\n", kind)
	printOutput(out)
}

func editAgentResource(kind, path string, overrides map[string]interface{}) {
	data, err := c.GetRaw(path)
	if err != nil {
		exitErr("get "+kind, err)
		return
	}
	result, err := editWithEditor(path, data, overrides)
	if err != nil {
		exitErr("edit "+kind, err)
		return
	}
	if result != nil {
		printOutput(result)
	}
}

func defaultLimitsTemplate() map[string]interface{} {
	return map[string]interface{}{
		"token_limit":  map[string]interface{}{"enabled": false, "group_cap": 0, "user_cap": 0, "window_seconds": 86400},
		"budget_limit": map[string]interface{}{"enabled": false, "group_cap_usd": 0, "user_cap_usd": 0, "window_seconds": 86400},
	}
}

// limitsFromFlags builds limits from the --token-*/--budget-* flags; nil when none was given.
func limitsFromFlags(cmd *cobra.Command) *client.AgentLimits {
	tokens := cmd.Flags().Changed("token-user-cap") || cmd.Flags().Changed("token-group-cap")
	budget := cmd.Flags().Changed("budget-user-cap") || cmd.Flags().Changed("budget-group-cap")
	if !tokens && !budget {
		return nil
	}
	window := int64(windowFlag.Seconds())
	if window < 60 {
		window = 60
	}
	return &client.AgentLimits{
		TokenLimit:  client.AgentTokenLimit{Enabled: tokens, UserCap: tokenUserCapFlag, GroupCap: tokenGroupCapFlag, WindowSeconds: window},
		BudgetLimit: client.AgentBudgetLimit{Enabled: budget, UserCapUSD: budgetUserCapFlag, GroupCapUSD: budgetGroupCapFlag, WindowSeconds: window},
	}
}

// catalogModelPrices turns model IDs into provider model entries priced from the catalog.
func catalogModelPrices(catalogID string, ids []string) ([]client.AgentProviderModel, error) {
	catalog, err := c.GetAgentCatalogProviders()
	if err != nil {
		return nil, err
	}
	prices := map[string]client.AgentCatalogModel{}
	for _, p := range catalog {
		if p.ID == catalogID {
			for _, m := range p.Models {
				prices[m.ID] = m
			}
		}
	}
	var models []client.AgentProviderModel
	for _, id := range ids {
		m, ok := prices[id]
		if !ok {
			return nil, fmt.Errorf("model %s is not in the %s catalog (price unknown); use --edit to set prices", id, catalogID)
		}
		models = append(models, client.AgentProviderModel{
			ID: id, InputPer1k: m.InputPer1k, OutputPer1k: m.OutputPer1k,
			CachedInputPer1k: m.CachedInputPer1k, CacheReadPer1k: m.CacheReadPer1k, CacheCreationPer1k: m.CacheCreationPer1k,
		})
	}
	return models, nil
}

func redactedProviderRequest(req *client.AgentProviderRequest) *client.AgentProviderRequest {
	r := *req
	if r.APIKey != "" {
		r.APIKey = "<redacted>"
	}
	return &r
}

// agentFilterParams builds the shared user/group/provider/model/session/date query filters.
func agentFilterParams() (map[string]string, bool) {
	params := map[string]string{}
	if userFlag != "" {
		id, err := c.ResolveUserID(userFlag)
		if err != nil {
			exitErr("user", err)
			return nil, false
		}
		params["user_id"] = id
	}
	if groupFlag != "" {
		id, err := c.ResolveGroupID(groupFlag)
		if err != nil {
			exitErr("group", err)
			return nil, false
		}
		params["group_id"] = id
	}
	if providerFlag != "" {
		id, err := c.ResolveAgentProviderID(providerFlag)
		if err != nil {
			exitErr("agent provider", err)
			return nil, false
		}
		params["provider_id"] = id
	}
	if modelFlag != "" {
		params["model"] = modelFlag
	}
	if sessionFlag != "" {
		params["session_id"] = sessionFlag
	}
	start := startFlag
	if sinceFlag != "" {
		d, err := parseSince(sinceFlag)
		if err != nil {
			exitErr("--since", err)
			return nil, false
		}
		start = time.Now().UTC().Add(-d).Format(time.RFC3339)
	}
	for key, val := range map[string]string{"start_date": start, "end_date": endFlag} {
		if val == "" {
			continue
		}
		t, err := parseDate(val)
		if err != nil {
			exitErr(key, err)
			return nil, false
		}
		params[key] = t
	}
	return params, true
}

// parseSince accepts Go durations plus a "d" (days) suffix: 90m, 24h, 7d.
func parseSince(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q", s)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

// parseDate accepts RFC3339 or YYYY-MM-DD and returns RFC3339.
func parseDate(s string) (string, error) {
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		return s, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return "", fmt.Errorf("invalid date %q (expected RFC3339 or YYYY-MM-DD)", s)
	}
	return t.UTC().Format(time.RFC3339), nil
}

func printPageFooter(page, totalPages, totalRecords int, what string) {
	if totalPages > 1 {
		fmt.Printf("\npage %d/%d — %d %s (use --page to see more)\n", page, totalPages, totalRecords, what)
	}
}

func nameOrID(resource, id string) string {
	if name := c.IDToName(resource, id); name != "" {
		return name
	}
	return id
}

func describeGuardrailChecks(ch client.AgentGuardrailChecks) string {
	var parts []string
	if ch.ModelAllowlist.Enabled {
		parts = append(parts, "models: "+strings.Join(ch.ModelAllowlist.Models, ", "))
	}
	if ch.PromptCapture.Enabled {
		p := "prompt capture"
		if ch.PromptCapture.RedactPII {
			p += " (PII redacted)"
		}
		parts = append(parts, p)
	}
	if len(parts) == 0 {
		return "no checks enabled"
	}
	return strings.Join(parts, "; ")
}

func describeTokenLimit(l client.AgentTokenLimit) string {
	if !l.Enabled {
		return "none"
	}
	return fmt.Sprintf("%s per user, %s per group, every %s",
		capStr(fmt.Sprintf("%d", l.UserCap), l.UserCap == 0), capStr(fmt.Sprintf("%d", l.GroupCap), l.GroupCap == 0), formatWindow(l.WindowSeconds))
}

func describeBudgetLimit(l client.AgentBudgetLimit) string {
	if !l.Enabled {
		return "none"
	}
	return fmt.Sprintf("%s per user, %s per group, every %s",
		capStr(fmt.Sprintf("$%.2f", l.UserCapUSD), l.UserCapUSD == 0), capStr(fmt.Sprintf("$%.2f", l.GroupCapUSD), l.GroupCapUSD == 0), formatWindow(l.WindowSeconds))
}

// formatWindow renders a window in seconds as 7d, 24h, 30m or 90s.
func formatWindow(seconds int64) string {
	switch {
	case seconds > 0 && seconds%86400 == 0 && seconds >= 2*86400:
		return fmt.Sprintf("%dd", seconds/86400)
	case seconds > 0 && seconds%3600 == 0:
		return fmt.Sprintf("%dh", seconds/3600)
	case seconds > 0 && seconds%60 == 0:
		return fmt.Sprintf("%dm", seconds/60)
	}
	return fmt.Sprintf("%ds", seconds)
}

func capStr(s string, uncapped bool) string {
	if uncapped {
		return "uncapped"
	}
	return s
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func addLimitFlags(cmd *cobra.Command) {
	cmd.Flags().Int64Var(&tokenUserCapFlag, "token-user-cap", 0, "Max tokens per user per window (0 = uncapped)")
	cmd.Flags().Int64Var(&tokenGroupCapFlag, "token-group-cap", 0, "Max tokens per source group per window (0 = uncapped)")
	cmd.Flags().Float64Var(&budgetUserCapFlag, "budget-user-cap", 0, "Max USD per user per window (0 = uncapped)")
	cmd.Flags().Float64Var(&budgetGroupCapFlag, "budget-group-cap", 0, "Max USD per group per window (0 = uncapped)")
	cmd.Flags().DurationVar(&windowFlag, "window", 24*time.Hour, "Cap reset window (e.g. 1h, 24h, 720h; min 1m)")
}

func addAgentFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&userFlag, "user", "", "Filter by user (name, email or ID)")
	cmd.Flags().StringVar(&groupFlag, "group", "", "Filter by authorising group (name or ID)")
	cmd.Flags().StringVar(&providerFlag, "provider", "", "Filter by agent provider (name or ID)")
	cmd.Flags().StringVar(&modelFlag, "model", "", "Filter by model")
	cmd.Flags().StringVar(&sessionFlag, "session", "", "Filter to one conversation / coding session ID")
	cmd.Flags().StringVar(&sinceFlag, "since", "", "Only since this long ago (e.g. 90m, 24h, 7d)")
	cmd.Flags().StringVar(&startFlag, "start", "", "Start date (RFC3339 or YYYY-MM-DD)")
	cmd.Flags().StringVar(&endFlag, "end", "", "End date (RFC3339 or YYYY-MM-DD)")
	cmd.RegisterFlagCompletionFunc("user", validArgsFunc(userNames))
	cmd.RegisterFlagCompletionFunc("group", validArgsFunc(groupNames))
	cmd.RegisterFlagCompletionFunc("provider", validArgsFunc(agentProviderNames))
}

func init() {
	getCmd.AddCommand(agentProvidersGetCmd, agentPoliciesGetCmd, guardrailsGetCmd, budgetRulesGetCmd,
		agentSettingsGetCmd, agentGatewayGetCmd, agentCatalogGetCmd, agentModelsGetCmd, agentConfigGetCmd,
		agentUsageGetCmd, agentConsumptionGetCmd)
	createCmd.AddCommand(agentProviderCreateCmd, agentPolicyCreateCmd, guardrailCreateCmd, budgetRuleCreateCmd,
		agentSettingsCreateCmd, agentGatewayCreateCmd)
	editCmd.AddCommand(agentProviderEditCmd, agentPolicyEditCmd, guardrailEditCmd, budgetRuleEditCmd, agentSettingsEditCmd)
	deleteCmd.AddCommand(agentProviderDeleteCmd, agentPolicyDeleteCmd, guardrailDeleteCmd, budgetRuleDeleteCmd, agentSettingsDeleteCmd)
	describeCmd.AddCommand(describeAgentPolicyCmd, describeBudgetRuleCmd)
	eventsCmd.AddCommand(agentLogsCmd)

	// providers
	for _, cmd := range []*cobra.Command{agentProviderCreateCmd, agentProviderEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Provider display name")
		cmd.Flags().StringVar(&upstreamURLFlag, "url", "", "Upstream URL (with scheme), e.g. https://api.openai.com")
		cmd.Flags().StringVar(&apiKeyFlag, "api-key", "", "Upstream API key (sealed server-side, never returned)")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.Flags().BoolVar(&skipTLSFlag, "skip-tls-verify", false, "Skip upstream TLS verification (self-hosted gateways)")
		cmd.Flags().BoolVar(&noMetadataFlag, "no-metadata", false, "Do not inject caller identity metadata upstream")
	}
	agentProviderCreateCmd.Flags().StringVar(&providerTypeFlag, "type", "", "Catalog provider id (see: netbird get agentcatalog)")
	agentProviderCreateCmd.Flags().StringSliceVar(&modelsFlag, "models", nil, "Restrict to these catalog models (default: all)")
	agentProviderCreateCmd.RegisterFlagCompletionFunc("type", validArgsFunc(agentCatalogIDs))

	// policies
	for _, cmd := range []*cobra.Command{agentPolicyCreateCmd, agentPolicyEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Policy name")
		cmd.Flags().StringVar(&descFlag, "desc", "", "Description")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.Flags().StringSliceVar(&sourceGroupsFlag, "source-groups", nil, "Groups allowed to use the providers (names or IDs)")
		cmd.Flags().StringSliceVar(&providersFlag, "providers", nil, "Agent providers reachable (names or IDs)")
		cmd.Flags().StringSliceVar(&guardrailsFlag, "guardrails", nil, "Guardrails to attach (names or IDs)")
		cmd.RegisterFlagCompletionFunc("source-groups", validArgsFunc(groupNames))
		cmd.RegisterFlagCompletionFunc("providers", validArgsFunc(agentProviderNames))
		cmd.RegisterFlagCompletionFunc("guardrails", validArgsFunc(guardrailNames))
	}
	addLimitFlags(agentPolicyCreateCmd)

	// guardrails
	for _, cmd := range []*cobra.Command{guardrailCreateCmd, guardrailEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Guardrail name")
		cmd.Flags().StringVar(&descFlag, "desc", "", "Description")
	}
	guardrailCreateCmd.Flags().StringSliceVar(&modelsFlag, "models", nil, "Model allowlist (enables the allowlist check)")
	guardrailCreateCmd.Flags().BoolVar(&promptCaptureFlag, "prompt-capture", false, "Capture prompts/completions (also needs agentsettings --prompt-collection)")
	guardrailCreateCmd.Flags().BoolVar(&redactPIIFlag, "redact-pii", false, "Redact PII from captured prompts")

	// budget rules
	for _, cmd := range []*cobra.Command{budgetRuleCreateCmd, budgetRuleEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Budget rule name")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enforced")
		cmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "Target groups (names or IDs)")
		cmd.Flags().StringSliceVar(&usersFlag, "users", nil, "Target users (names, emails or IDs)")
		cmd.RegisterFlagCompletionFunc("groups", validArgsFunc(groupNames))
		cmd.RegisterFlagCompletionFunc("users", validArgsFunc(userNames))
	}
	addLimitFlags(budgetRuleCreateCmd)

	// settings
	agentSettingsCreateCmd.Flags().StringVar(&proxyAddrFlag, "proxy-address", "", "Proxy cluster address to allocate the endpoint under")
	agentSettingsCreateCmd.Flags().StringVar(&endpointFlag, "endpoint", "", "Hostname to claim as a dedicated endpoint")
	for _, cmd := range []*cobra.Command{agentSettingsCreateCmd, agentSettingsEditCmd} {
		cmd.Flags().BoolVar(&logCollectionFlag, "log-collection", true, "Collect per-request access logs")
		cmd.Flags().BoolVar(&promptCollectFlag, "prompt-collection", false, "Allow prompt capture (master switch, also needs a guardrail)")
		cmd.Flags().BoolVar(&redactPIIFlag, "redact-pii", false, "Redact PII in captured prompts")
		cmd.Flags().IntVar(&retentionDaysFlag, "retention-days", 30, "Access log retention in days (0 = forever)")
	}

	// discovery
	agentModelsGetCmd.Flags().StringVar(&providerFlag, "provider", "", "Existing agent provider (name or ID) whose credential to use")
	agentModelsGetCmd.Flags().StringVar(&providerTypeFlag, "type", "", "Catalog provider id (for a provider not created yet)")
	agentModelsGetCmd.Flags().StringVar(&apiKeyFlag, "api-key", "", "API key to query with (for a provider not created yet)")
	agentModelsGetCmd.Flags().StringVar(&upstreamURLFlag, "url", "", "Upstream URL override")
	agentModelsGetCmd.RegisterFlagCompletionFunc("provider", validArgsFunc(agentProviderNames))
	agentModelsGetCmd.RegisterFlagCompletionFunc("type", validArgsFunc(agentCatalogIDs))

	// usage & logs
	addAgentFilterFlags(agentUsageGetCmd)
	agentUsageGetCmd.Flags().StringVar(&granularityFlag, "granularity", "day", "Bucket width: day, week, month")
	agentUsageGetCmd.RegisterFlagCompletionFunc("granularity", staticCompletion([]string{"day", "week", "month"}))

	addAgentFilterFlags(agentLogsCmd)
	agentLogsCmd.Flags().BoolVar(&sessionsFlag, "sessions", false, "Group requests by conversation / coding session")
	agentLogsCmd.Flags().StringVar(&decisionFlag, "decision", "", "Filter by policy decision: allow, deny")
	agentLogsCmd.Flags().StringVar(&searchFlag, "search", "", "Search in ID, host, path, model, user")
	agentLogsCmd.Flags().StringVar(&sortByFlag, "sort-by", "", "Sort field: timestamp, model, provider, status_code, duration, cost_usd, total_tokens, user_id, decision")
	agentLogsCmd.Flags().StringVar(&sortOrderFlag, "sort-order", "", "Sort order: asc, desc")
	agentLogsCmd.Flags().IntVar(&pageFlag, "page", 1, "Page number")
	agentLogsCmd.Flags().IntVar(&pageSizeFlag, "page-size", 50, "Page size (max 100)")
	agentLogsCmd.RegisterFlagCompletionFunc("decision", staticCompletion([]string{"allow", "deny"}))
	agentLogsCmd.RegisterFlagCompletionFunc("sort-by", staticCompletion([]string{"timestamp", "model", "provider", "status_code", "duration", "cost_usd", "total_tokens", "user_id", "decision"}))
	agentLogsCmd.RegisterFlagCompletionFunc("sort-order", staticCompletion([]string{"asc", "desc"}))
}
