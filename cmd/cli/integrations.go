package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

// Integrations: event streaming, notification channels, EDR, IdP sync.

// --- Event streaming ---

type eventStreamRow struct {
	StreamID int    `json:"stream_id"`
	Platform string `json:"platform"`
	Enabled  bool   `json:"enabled"`
	Config   string `json:"config_keys"`
}

var eventStreamsGetCmd = &cobra.Command{
	Use:               "eventstreams [id|platform]",
	Aliases:           []string{"eventstream", "es"},
	Short:             "List or display event streaming integrations (Datadog, S3, Firehose, HTTP)",
	ValidArgsFunction: validArgsFunc(eventStreamNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			id, err := c.ResolveEventStreamID(args[0])
			if err != nil {
				exitErr("event stream", err)
				return
			}
			s, err := c.GetEventStream(id)
			if err != nil {
				exitErr("event stream", err)
				return
			}
			printOutput(s)
			return
		}
		streams, err := c.GetEventStreams()
		if err != nil {
			exitErr("event streams", err)
			return
		}
		if outputFormat != "" {
			printOutput(streams)
			return
		}
		var rows []eventStreamRow
		for _, s := range streams {
			rows = append(rows, eventStreamRow{s.ID, s.Platform, s.Enabled, strings.Join(sortedKeys(s.Config), ", ")})
		}
		printOutput(rows)
	},
}

var eventStreamCreateCmd = &cobra.Command{
	Use:     "eventstream",
	Aliases: []string{"es"},
	Short:   "Create an event streaming integration",
	Long: `Stream NetBird audit events to an external platform.
--config takes the platform-specific key=value settings (credentials, region, bucket, URL...).`,
	Example: `  netbird create eventstream --platform datadog --config api_key=$DD_API_KEY,api_url=https://http-intake.logs.datadoghq.eu
  netbird create eventstream --platform generic_http --config url=https://siem.example.com/ingest`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var s client.EventStream
			createAgentWithEditor("eventstream", "/api/event-streaming", map[string]interface{}{
				"platform": "datadog",
				"enabled":  true,
				"config":   map[string]string{"api_key": "", "api_url": ""},
			}, &s)
			return
		}
		if platformFlag == "" || len(configFlag) == 0 {
			exitErr("--platform and --config are required (or use --edit)", nil)
			return
		}
		req := &client.EventStreamRequest{Platform: platformFlag, Config: configFlag, Enabled: enabledFlag}
		if dryRunCheck(maskedEventStream(req)) {
			return
		}
		s, err := c.CreateEventStream(req)
		if err != nil {
			exitErr("create eventstream", err)
			return
		}
		fmt.Println("event stream created:")
		printOutput(s)
	},
}

var eventStreamEditCmd = &cobra.Command{
	Use:               "eventstream <id|platform>",
	Aliases:           []string{"es"},
	Short:             "Edit an event streaming integration",
	Long:              "Edit an event stream. Secret config values come back masked (****); they are left out of the editor and kept server-side unless you set them again.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(eventStreamNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveEventStreamID(args[0])
		if err != nil {
			exitErr("event stream", err)
			return
		}
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		editTransformed("eventstream", "/api/event-streaming/"+id, func(m map[string]interface{}) {
			delete(m, "account_id")
			delete(m, "created_at")
			delete(m, "updated_at")
			if cfg, ok := m["config"].(map[string]interface{}); ok {
				stripMasked(cfg)
				for k, v := range configFlag {
					cfg[k] = v
				}
			}
		}, overrides)
	},
}

var eventStreamDeleteCmd = &cobra.Command{
	Use:               "eventstream <id|platform>",
	Aliases:           []string{"es"},
	Short:             "Delete an event streaming integration",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(eventStreamNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveEventStreamID(args[0])
		if err != nil {
			exitErr("event stream", err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete event stream %s", id)) {
			return
		}
		if err := c.DeleteEventStream(id); err != nil {
			exitErr("delete eventstream", err)
			return
		}
		fmt.Println("event stream deleted")
	},
}

// --- Notification channels ---

type notificationChannelRow struct {
	ChannelID  string   `json:"channel_id"`
	Type       string   `json:"type"`
	Target     string   `json:"target_summary"`
	EventTypes []string `json:"event_types"`
	Enabled    bool     `json:"enabled"`
}

var notificationChannelsGetCmd = &cobra.Command{
	Use:     "notificationchannels [id]",
	Aliases: []string{"notificationchannel", "nc"},
	Short:   "List or display notification channels (email, webhook)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			ch, err := c.GetNotificationChannel(args[0])
			if err != nil {
				exitErr("notification channel", err)
				return
			}
			printOutput(ch)
			return
		}
		channels, err := c.GetNotificationChannels()
		if err != nil {
			exitErr("notification channels", err)
			return
		}
		if outputFormat != "" {
			printOutput(channels)
			return
		}
		var rows []notificationChannelRow
		for _, ch := range channels {
			rows = append(rows, notificationChannelRow{ch.ID, ch.Type, targetSummary(ch.Target), ch.EventTypes, ch.Enabled})
		}
		printOutput(rows)
	},
}

type notificationTypeRow struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

var notificationTypesGetCmd = &cobra.Command{
	Use:     "notificationtypes",
	Aliases: []string{"nt"},
	Short:   "List the event types a notification channel can subscribe to",
	Run: func(cmd *cobra.Command, args []string) {
		types, err := c.GetNotificationTypes()
		if err != nil {
			exitErr("notification types", err)
			return
		}
		if outputFormat != "" {
			printOutput(types)
			return
		}
		var rows []notificationTypeRow
		for _, k := range sortedKeys(types) {
			rows = append(rows, notificationTypeRow{k, types[k]})
		}
		printOutput(rows)
	},
}

var notificationChannelCreateCmd = &cobra.Command{
	Use:     "notificationchannel",
	Aliases: []string{"nc"},
	Short:   "Create a notification channel",
	Example: `  netbird create notificationchannel --type email --emails ops@example.com --events user.join,peer.add
  netbird create notificationchannel --type webhook --url https://hooks.example.com/nb --header Authorization="Bearer x" --events user.join`,
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			var ch client.NotificationChannel
			createAgentWithEditor("notificationchannel", "/api/integrations/notifications/channels", map[string]interface{}{
				"type":        "email",
				"target":      map[string]interface{}{"emails": []string{}},
				"event_types": []string{},
				"enabled":     true,
			}, &ch)
			return
		}
		target, err := notificationTarget(notifTypeFlag)
		if err != nil {
			exitErr("create notificationchannel", err)
			return
		}
		if len(eventTypesFlag) == 0 {
			exitErr("--events is required (see: netbird get notificationtypes)", nil)
			return
		}
		req := &client.NotificationChannelRequest{Type: notifTypeFlag, Target: target, EventTypes: eventTypesFlag, Enabled: enabledFlag}
		if dryRunCheck(req) {
			return
		}
		ch, err := c.CreateNotificationChannel(req)
		if err != nil {
			exitErr("create notificationchannel", err)
			return
		}
		fmt.Println("notification channel created:")
		printOutput(ch)
	},
}

var notificationChannelEditCmd = &cobra.Command{
	Use:     "notificationchannel <id>",
	Aliases: []string{"nc"},
	Short:   "Edit a notification channel",
	Long:    "Edit a notification channel. Webhook header values are write-only (masked on read): they are left out of the editor; pass --header to set them again.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		if cmd.Flags().Changed("events") {
			overrides["event_types"] = eventTypesFlag
		}
		editTransformed("notificationchannel", "/api/integrations/notifications/channels/"+args[0], func(m map[string]interface{}) {
			target, _ := m["target"].(map[string]interface{})
			if target == nil {
				return
			}
			if cmd.Flags().Changed("emails") {
				target["emails"] = emailsFlag
			}
			if cmd.Flags().Changed("url") {
				target["url"] = urlFlag
			}
			if headers, ok := target["headers"].(map[string]interface{}); ok {
				stripMasked(headers)
				if len(headers) == 0 {
					delete(target, "headers")
				}
			}
			if len(headersFlag) > 0 {
				target["headers"] = headersFlag
			}
		}, overrides)
	},
}

var notificationChannelDeleteCmd = &cobra.Command{
	Use:     "notificationchannel <id>",
	Aliases: []string{"nc"},
	Short:   "Delete a notification channel",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete notification channel %s", args[0])) {
			return
		}
		if err := c.DeleteNotificationChannel(args[0]); err != nil {
			exitErr("delete notificationchannel", err)
			return
		}
		fmt.Println("notification channel deleted")
	},
}

// --- EDR integrations ---

var edrGetCmd = &cobra.Command{
	Use:       "edr [vendor]",
	Short:     "Show EDR integrations (intune, sentinelone, falcon, huntress, fleetdm)",
	ValidArgs: client.EDRVendors,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			e, err := c.GetEDRIntegration(args[0])
			if err != nil {
				if is404(err) {
					fmt.Printf("No %s integration configured.\n", args[0])
					return
				}
				exitErr("edr", err)
				return
			}
			printOutput(e)
			return
		}
		var items []client.EDRIntegration
		for _, v := range client.EDRVendors {
			e, err := c.GetEDRIntegration(v)
			if err != nil {
				if !is404(err) {
					fmt.Fprintf(os.Stderr, "warning: %s: %s\n", v, err)
				}
				continue
			}
			items = append(items, *e)
		}
		if len(items) == 0 && outputFormat == "" {
			fmt.Println("No EDR integration configured.")
			return
		}
		printOutput(items)
	},
}

var edrCreateCmd = &cobra.Command{
	Use:       "edr <vendor>",
	Short:     "Create an EDR integration (peers must be compliant to connect)",
	ValidArgs: client.EDRVendors,
	Args:      cobra.ExactArgs(1),
	Example: `  netbird create edr intune --client-id <app-id> --tenant-id <tenant> --secret $SECRET --groups laptops
  netbird create edr sentinelone --api-url https://x.sentinelone.net --api-token $TOKEN --groups laptops
  netbird create edr falcon --client-id <id> --secret $SECRET --cloud-id eu-1 --zta-threshold 60 --groups laptops
  netbird create edr fleetdm --edit`,
	Run: func(cmd *cobra.Command, args []string) {
		vendor := args[0]
		template, ok := edrTemplates[vendor]
		if !ok {
			exitErr(fmt.Sprintf("unknown EDR vendor %s (%s)", vendor, strings.Join(client.EDRVendors, ", ")), nil)
			return
		}
		if editFlag {
			var e client.EDRIntegration
			createAgentWithEditor("edr "+vendor, "/api/integrations/edr/"+vendor, template(), &e)
			return
		}
		req := template()
		for key, flag := range edrFlagFields {
			if _, used := req[key]; used && cmd.Flags().Changed(flag) {
				req[key] = edrFlagValue(flag)
			}
		}
		req["groups"] = nonNil(resolveAll("group", groupsFlag))
		req["enabled"] = enabledFlag
		if dryRunCheck(maskSecrets(req)) {
			return
		}
		e, err := c.CreateEDRIntegration(vendor, req)
		if err != nil {
			exitErr("create edr", err)
			return
		}
		fmt.Printf("%s integration created:\n", vendor)
		printOutput(e)
	},
}

var edrEditCmd = &cobra.Command{
	Use:       "edr <vendor>",
	Short:     "Edit an EDR integration",
	Long:      "Edit an EDR integration. The API credentials are never returned: set them again with the flags (--secret, --api-token, ...) or add them in the editor.",
	ValidArgs: client.EDRVendors,
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vendor := args[0]
		overrides := map[string]interface{}{}
		for key, flag := range edrFlagFields {
			if cmd.Flags().Changed(flag) {
				overrides[key] = edrFlagValue(flag)
			}
		}
		if cmd.Flags().Changed("groups") {
			overrides["groups"] = resolveAll("group", groupsFlag)
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		editTransformed("edr "+vendor, "/api/integrations/edr/"+vendor, func(m map[string]interface{}) {
			for _, k := range []string{"id", "account_id", "created_by", "created_at", "updated_at", "last_synced_at"} {
				delete(m, k)
			}
			m["groups"] = groupObjectsToIDs(m["groups"])
		}, overrides)
	},
}

var edrDeleteCmd = &cobra.Command{
	Use:       "edr <vendor>",
	Short:     "Delete an EDR integration",
	ValidArgs: client.EDRVendors,
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if dryRunMsg(fmt.Sprintf("would delete the %s EDR integration", args[0])) {
			return
		}
		if err := c.DeleteEDRIntegration(args[0]); err != nil {
			exitErr("delete edr", err)
			return
		}
		fmt.Printf("%s integration deleted\n", args[0])
	},
}

var edrBypassedGetCmd = &cobra.Command{
	Use:     "edrbypassed",
	Aliases: []string{"bypassed"},
	Short:   "List peers whose EDR compliance check is bypassed",
	Run: func(cmd *cobra.Command, args []string) {
		peers, err := c.GetEDRBypassedPeers()
		if err != nil {
			exitErr("edr bypassed peers", err)
			return
		}
		printOutput(peers)
	},
}

var bypassCmd = &cobra.Command{
	Use:               "bypass peer <name|id>",
	Short:             "Bypass the EDR compliance check for a non-compliant peer",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: verbPeerArgs,
	Run: func(cmd *cobra.Command, args []string) {
		peerEDRBypass(args, true)
	},
}

var unbypassCmd = &cobra.Command{
	Use:               "unbypass peer <name|id>",
	Short:             "Revoke the EDR compliance bypass of a peer",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: verbPeerArgs,
	Run: func(cmd *cobra.Command, args []string) {
		peerEDRBypass(args, false)
	},
}

func peerEDRBypass(args []string, bypass bool) {
	verb := map[bool]string{true: "bypass", false: "unbypass"}[bypass]
	if args[0] != "peer" && args[0] != "pr" {
		fmt.Printf("usage: netbird %s peer <name|id>\n", verb)
		return
	}
	id, err := c.ResolvePeerID(args[1])
	if err != nil {
		exitErr("peer", err)
		return
	}
	if dryRunMsg(fmt.Sprintf("would %s EDR compliance for peer %s", verb, args[1])) {
		return
	}
	if bypass {
		err = c.BypassPeerEDR(id)
	} else {
		err = c.RevokePeerEDRBypass(id)
	}
	if err != nil {
		exitErr(verb+" peer", err)
		return
	}
	if bypass {
		fmt.Printf("EDR compliance bypassed for peer %s\n", args[1])
	} else {
		fmt.Printf("EDR compliance bypass revoked for peer %s\n", args[1])
	}
}

// edrTemplates gives the request fields of each vendor, with defaults.
var edrTemplates = map[string]func() map[string]interface{}{
	"intune": func() map[string]interface{} {
		return map[string]interface{}{"client_id": "", "tenant_id": "", "secret": "", "groups": []string{}, "last_synced_interval": 24, "enabled": true}
	},
	"sentinelone": func() map[string]interface{} {
		return map[string]interface{}{"api_url": "", "api_token": "", "groups": []string{}, "last_synced_interval": 24, "enabled": true,
			"match_attributes": map[string]interface{}{"active_threats": 0, "infected": false, "is_active": true}}
	},
	"falcon": func() map[string]interface{} {
		return map[string]interface{}{"client_id": "", "secret": "", "cloud_id": "us-1", "groups": []string{}, "zta_score_threshold": 50, "enabled": true}
	},
	"huntress": func() map[string]interface{} {
		return map[string]interface{}{"api_key": "", "api_secret": "", "groups": []string{}, "last_synced_interval": 24, "enabled": true,
			"match_attributes": map[string]interface{}{"defender_status": "Healthy"}}
	},
	"fleetdm": func() map[string]interface{} {
		return map[string]interface{}{"api_url": "", "api_token": "", "groups": []string{}, "last_synced_interval": 24, "enabled": true,
			"match_attributes": map[string]interface{}{"disk_encryption_enabled": true, "failing_policies_count_max": 0, "status_online": true}}
	},
}

// edrFlagFields maps EDR request fields to the CLI flag that sets them.
var edrFlagFields = map[string]string{
	"client_id":            "client-id",
	"tenant_id":            "tenant-id",
	"secret":               "secret",
	"api_url":              "api-url",
	"api_token":            "api-token",
	"api_key":              "api-key",
	"api_secret":           "api-secret",
	"cloud_id":             "cloud-id",
	"zta_score_threshold":  "zta-threshold",
	"last_synced_interval": "sync-interval",
}

func edrFlagValue(flag string) interface{} {
	switch flag {
	case "client-id":
		return clientIDFlag
	case "tenant-id":
		return tenantIDFlag
	case "secret":
		return secretFlag
	case "api-url":
		return upstreamURLFlag
	case "api-token":
		return apiTokenFlag
	case "api-key":
		return apiKeyFlag
	case "api-secret":
		return apiSecretFlag
	case "cloud-id":
		return cloudIDFlag
	case "zta-threshold":
		return ztaThresholdFlag
	case "sync-interval":
		return syncIntervalFlag
	}
	return nil
}

// --- IdP sync (Google Workspace, Azure/Entra, Okta SCIM, SCIM) ---

var idpSyncKinds = []string{"google", "azure", "okta", "scim"}

type idpSyncRow struct {
	Kind         string `json:"kind"`
	SyncID       int    `json:"sync_id"`
	Enabled      bool   `json:"enabled"`
	Detail       string `json:"detail"`
	SyncInterval int    `json:"sync_interval,omitempty"`
	LastSyncedAt string `json:"last_synced_at"`
}

var idpSyncsGetCmd = &cobra.Command{
	Use:     "idpsyncs [kind] [id]",
	Aliases: []string{"idpsync", "ids"},
	Short:   "List IdP user/group sync integrations (google, azure, okta, scim)",
	Args:    cobra.MaximumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return idpSyncKinds, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		kinds := idpSyncKinds
		if len(args) > 0 {
			kinds = []string{args[0]}
		}
		var items []client.IdPSync
		for _, k := range kinds {
			list, err := c.GetIdPSyncs(k)
			if err != nil {
				if len(args) == 0 && is404(err) {
					continue
				}
				exitErr("idp syncs", err)
				return
			}
			items = append(items, list...)
		}
		if len(args) == 2 {
			for _, s := range items {
				if strconv.Itoa(s.ID) == args[1] {
					printOutput(s)
					return
				}
			}
			exitErr(fmt.Sprintf("%s sync %s not found", args[0], args[1]), nil)
			return
		}
		if outputFormat != "" {
			printOutput(items)
			return
		}
		var rows []idpSyncRow
		for _, s := range items {
			rows = append(rows, idpSyncRow{s.Kind, s.ID, s.Enabled, idpSyncDetail(s), s.SyncInterval, s.LastSyncedAt})
		}
		printOutput(rows)
	},
}

var idpSyncCreateCmd = &cobra.Command{
	Use:       "idpsync <google|azure|okta|scim>",
	Aliases:   []string{"ids"},
	Short:     "Create an IdP user/group sync integration",
	ValidArgs: idpSyncKinds,
	Args:      cobra.ExactArgs(1),
	Example: `  netbird create idpsync google --customer-id C01abc --service-account-key-file sa.json --group-prefixes eng-
  netbird create idpsync azure --client-id <app> --client-secret $SECRET --tenant-id <tenant>
  netbird create idpsync okta --connection-name okta-prod      # prints the SCIM token
  netbird create idpsync scim --provider jumpcloud --prefix jc-`,
	Run: func(cmd *cobra.Command, args []string) {
		kind := args[0]
		template, ok := idpSyncTemplates[kind]
		if !ok {
			exitErr(fmt.Sprintf("unknown IdP sync type %s (google, azure, okta, scim)", kind), nil)
			return
		}
		path := "/api/integrations/" + client.IdPSyncKinds[kind]
		if editFlag {
			var s client.IdPSync
			createAgentWithEditor("idpsync "+kind, path, template(), &s)
			return
		}
		req := template()
		for key, flag := range idpSyncFlagFields {
			if _, used := req[key]; used && cmd.Flags().Changed(flag) {
				v, err := idpSyncFlagValue(flag)
				if err != nil {
					exitErr(flag, err)
					return
				}
				req[key] = v
			}
		}
		if dryRunCheck(maskSecrets(req)) {
			return
		}
		s, err := c.CreateIdPSync(kind, req)
		if err != nil {
			exitErr("create idpsync", err)
			return
		}
		fmt.Printf("%s sync created:\n", kind)
		printOutput(s)
		if s.AuthToken != "" && outputFormat == "" {
			fmt.Printf("\nSCIM token (shown once, configure it in your IdP): %s\n", s.AuthToken)
		}
	},
}

var idpSyncEditCmd = &cobra.Command{
	Use:       "idpsync <google|azure|okta|scim> [id]",
	Aliases:   []string{"ids"},
	Short:     "Edit an IdP sync integration (id optional when there is only one of that kind)",
	ValidArgs: idpSyncKinds,
	Args:      cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		kind, id, ok := resolveIdPSync(args)
		if !ok {
			return
		}
		overrides := map[string]interface{}{}
		for key, flag := range idpSyncFlagFields {
			if cmd.Flags().Changed(flag) && key != "host" && key != "connection_name" && key != "provider" {
				v, err := idpSyncFlagValue(flag)
				if err != nil {
					exitErr(flag, err)
					return
				}
				overrides[key] = v
			}
		}
		if cmd.Flags().Changed("enabled") {
			overrides["enabled"] = enabledFlag
		}
		editTransformed("idpsync "+kind, fmt.Sprintf("/api/integrations/%s/%s", client.IdPSyncKinds[kind], id), func(m map[string]interface{}) {
			for _, k := range []string{"id", "last_synced_at", "auth_token", "host", "provider"} {
				delete(m, k)
			}
		}, overrides)
	},
}

var idpSyncDeleteCmd = &cobra.Command{
	Use:       "idpsync <google|azure|okta|scim> [id]",
	Aliases:   []string{"ids"},
	Short:     "Delete an IdP sync integration",
	ValidArgs: idpSyncKinds,
	Args:      cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		kind, id, ok := resolveIdPSync(args)
		if !ok {
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete %s sync %s", kind, id)) {
			return
		}
		if err := c.DeleteIdPSync(kind, id); err != nil {
			exitErr("delete idpsync", err)
			return
		}
		fmt.Printf("%s sync deleted\n", kind)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync idp <google|azure> [id]",
	Short: "Trigger an immediate IdP user/group sync",
	Args:  cobra.RangeArgs(2, 3),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return []string{"idp"}, cobra.ShellCompDirectiveNoFileComp
		case 1:
			return []string{"google", "azure"}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "idp" {
			fmt.Println("usage: netbird sync idp <google|azure> [id]")
			return
		}
		kind, id, ok := resolveIdPSync(args[1:])
		if !ok {
			return
		}
		if dryRunMsg(fmt.Sprintf("would sync %s integration %s", kind, id)) {
			return
		}
		result, err := c.SyncIdP(kind, id)
		if err != nil {
			exitErr("sync", err)
			return
		}
		fmt.Printf("%s sync triggered: %s\n", kind, result)
	},
}

var idpSyncLogsGetCmd = &cobra.Command{
	Use:       "idpsynclogs <google|azure|okta|scim> [id]",
	Aliases:   []string{"idpsynclog"},
	Short:     "Show the sync logs of an IdP sync integration",
	ValidArgs: idpSyncKinds,
	Args:      cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		kind, id, ok := resolveIdPSync(args)
		if !ok {
			return
		}
		logs, err := c.GetIdPSyncLogs(kind, id)
		if err != nil {
			exitErr("idp sync logs", err)
			return
		}
		printOutput(logs)
	},
}

var idpSyncTokenCmd = &cobra.Command{
	Use:   "scimtoken <okta|scim> [id]",
	Short: "Regenerate the SCIM token of an Okta SCIM / SCIM sync integration",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		kind, id, ok := resolveIdPSync(args)
		if !ok {
			return
		}
		if dryRunMsg(fmt.Sprintf("would regenerate the SCIM token of %s sync %s", kind, id)) {
			return
		}
		token, err := c.RegenerateIdPSyncToken(kind, id)
		if err != nil {
			exitErr("regenerate scim token", err)
			return
		}
		fmt.Printf("new SCIM token (shown once, update it in your IdP): %s\n", token)
	},
}

var idpSyncTemplates = map[string]func() map[string]interface{}{
	"google": func() map[string]interface{} {
		return map[string]interface{}{"customer_id": "", "service_account_key": "", "sync_interval": 300, "group_prefixes": []string{}, "user_group_prefixes": []string{}}
	},
	"azure": func() map[string]interface{} {
		return map[string]interface{}{"client_id": "", "client_secret": "", "tenant_id": "", "host": "microsoft.com", "sync_interval": 300, "group_prefixes": []string{}, "user_group_prefixes": []string{}}
	},
	"okta": func() map[string]interface{} {
		return map[string]interface{}{"connection_name": "", "group_prefixes": []string{}, "user_group_prefixes": []string{}}
	},
	"scim": func() map[string]interface{} {
		return map[string]interface{}{"provider": "", "prefix": "", "group_prefixes": []string{}, "user_group_prefixes": []string{}}
	},
}

var idpSyncFlagFields = map[string]string{
	"customer_id":         "customer-id",
	"service_account_key": "service-account-key-file",
	"client_id":           "client-id",
	"client_secret":       "client-secret",
	"tenant_id":           "tenant-id",
	"host":                "host",
	"sync_interval":       "sync-interval",
	"group_prefixes":      "group-prefixes",
	"user_group_prefixes": "user-group-prefixes",
	"connector_id":        "connector-id",
	"connection_name":     "connection-name",
	"provider":            "provider",
	"prefix":              "prefix",
}

func idpSyncFlagValue(flag string) (interface{}, error) {
	switch flag {
	case "customer-id":
		return customerIDFlag, nil
	case "service-account-key-file":
		data, err := os.ReadFile(fileFlag)
		if err != nil {
			return nil, err
		}
		return string(data), nil
	case "client-id":
		return clientIDFlag, nil
	case "client-secret":
		return clientSecretFlag, nil
	case "tenant-id":
		return tenantIDFlag, nil
	case "host":
		return hostFlag, nil
	case "sync-interval":
		return syncIntervalFlag, nil
	case "group-prefixes":
		return nonNil(groupPrefixesFlag), nil
	case "user-group-prefixes":
		return nonNil(userGroupPrefixesFlag), nil
	case "connector-id":
		return connectorIDFlag, nil
	case "connection-name":
		return connectionNameFlag, nil
	case "provider":
		return providerFlag, nil
	case "prefix":
		return prefixFlag, nil
	}
	return nil, nil
}

// resolveIdPSync turns <kind> [id] into a kind and ID; the ID may be omitted when only one integration of that kind exists.
func resolveIdPSync(args []string) (string, string, bool) {
	kind := args[0]
	if _, ok := client.IdPSyncKinds[kind]; !ok {
		exitErr(fmt.Sprintf("unknown IdP sync type %s (google, azure, okta, scim)", kind), nil)
		return "", "", false
	}
	if len(args) > 1 {
		return kind, args[1], true
	}
	items, err := c.GetIdPSyncs(kind)
	if err != nil {
		exitErr("idp syncs", err)
		return "", "", false
	}
	switch len(items) {
	case 0:
		exitErr(fmt.Sprintf("no %s sync integration", kind), nil)
		return "", "", false
	case 1:
		return kind, strconv.Itoa(items[0].ID), true
	}
	exitErr(fmt.Sprintf("several %s sync integrations exist, pass the ID (netbird get idpsyncs %s)", kind, kind), nil)
	return "", "", false
}

func idpSyncDetail(s client.IdPSync) string {
	switch s.Kind {
	case "google":
		return "customer " + s.CustomerID
	case "azure":
		return fmt.Sprintf("tenant %s (%s)", s.TenantID, s.Host)
	case "scim":
		return strings.TrimSpace(s.Provider + " " + s.Prefix)
	}
	return ""
}

// --- Helpers ---

// editTransformed runs the edit flow after letting transform clean up the fetched object
// (drop read-only or masked fields, reshape fields the PUT expects differently).
func editTransformed(kind, path string, transform func(map[string]interface{}), overrides map[string]interface{}) {
	data, err := c.GetRaw(path)
	if err != nil {
		exitErr("get "+kind, err)
		return
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		exitErr("get "+kind, err)
		return
	}
	transform(m)
	data, _ = json.Marshal(m)
	result, err := editWithEditor(path, data, overrides)
	if err != nil {
		exitErr("edit "+kind, err)
		return
	}
	if result != nil {
		printOutput(result)
	}
}

// stripMasked removes values the API returns masked ("****"), so they are not written back.
func stripMasked(m map[string]interface{}) {
	for k, v := range m {
		if s, ok := v.(string); ok && s != "" && strings.Trim(s, "*") == "" {
			delete(m, k)
		}
	}
}

var secretKeys = []string{"secret", "api_token", "api_key", "api_secret", "client_secret", "service_account_key", "password", "token"}

// maskSecrets returns a copy of req with secret-looking values replaced, for --dry-run output.
func maskSecrets(req map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(req))
	for k, v := range req {
		out[k] = v
		for _, s := range secretKeys {
			if strings.Contains(k, s) {
				if str, ok := v.(string); ok && str != "" {
					out[k] = "<redacted>"
				}
			}
		}
	}
	return out
}

func maskedEventStream(req *client.EventStreamRequest) *client.EventStreamRequest {
	r := *req
	cfg := map[string]interface{}{}
	for k, v := range req.Config {
		cfg[k] = v
	}
	r.Config = map[string]string{}
	for k, v := range maskSecrets(cfg) {
		r.Config[k] = v.(string)
	}
	return &r
}

func notificationTarget(kind string) (map[string]interface{}, error) {
	switch kind {
	case "email":
		if len(emailsFlag) == 0 {
			return nil, fmt.Errorf("--emails is required for an email channel")
		}
		return map[string]interface{}{"emails": emailsFlag}, nil
	case "webhook":
		if urlFlag == "" {
			return nil, fmt.Errorf("--url is required for a webhook channel")
		}
		t := map[string]interface{}{"url": urlFlag}
		if len(headersFlag) > 0 {
			t["headers"] = headersFlag
		}
		return t, nil
	}
	return nil, fmt.Errorf("--type must be email or webhook")
}

func targetSummary(t map[string]interface{}) string {
	if u, ok := t["url"].(string); ok {
		return u
	}
	if emails, ok := t["emails"].([]interface{}); ok {
		var parts []string
		for _, e := range emails {
			parts = append(parts, fmt.Sprint(e))
		}
		return strings.Join(parts, ", ")
	}
	return ""
}

// groupObjectsToIDs turns [{id, name, ...}] (as returned by GET) into [id] (as expected by PUT).
func groupObjectsToIDs(v interface{}) []string {
	ids := []string{}
	list, _ := v.([]interface{})
	for _, g := range list {
		switch g := g.(type) {
		case map[string]interface{}:
			if id, ok := g["id"].(string); ok {
				ids = append(ids, id)
			}
		case string:
			ids = append(ids, g)
		}
	}
	return ids
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func is404(err error) bool {
	return err != nil && strings.Contains(err.Error(), "HTTP 404")
}

func verbPeerArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return []string{"peer"}, cobra.ShellCompDirectiveNoFileComp
	}
	return validArgsFunc(peerNames)(cmd, args, toComplete)
}

func init() {
	getCmd.AddCommand(eventStreamsGetCmd, notificationChannelsGetCmd, notificationTypesGetCmd, edrGetCmd, edrBypassedGetCmd, idpSyncsGetCmd, idpSyncLogsGetCmd)
	createCmd.AddCommand(eventStreamCreateCmd, notificationChannelCreateCmd, edrCreateCmd, idpSyncCreateCmd)
	editCmd.AddCommand(eventStreamEditCmd, notificationChannelEditCmd, edrEditCmd, idpSyncEditCmd)
	deleteCmd.AddCommand(eventStreamDeleteCmd, notificationChannelDeleteCmd, edrDeleteCmd, idpSyncDeleteCmd)
	rootCmd.AddCommand(bypassCmd, unbypassCmd, syncCmd)
	regenerateCmd.AddCommand(idpSyncTokenCmd)

	// event streams
	eventStreamCreateCmd.Flags().StringVar(&platformFlag, "platform", "", "Platform: datadog, s3, firehose, generic_http")
	eventStreamCreateCmd.RegisterFlagCompletionFunc("platform", staticCompletion([]string{"datadog", "s3", "firehose", "generic_http"}))
	for _, cmd := range []*cobra.Command{eventStreamCreateCmd, eventStreamEditCmd} {
		cmd.Flags().StringToStringVar(&configFlag, "config", nil, "Platform config as key=value pairs")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
	}

	// notification channels
	notificationChannelCreateCmd.Flags().StringVar(&notifTypeFlag, "type", "email", "Channel type: email, webhook")
	notificationChannelCreateCmd.RegisterFlagCompletionFunc("type", staticCompletion([]string{"email", "webhook"}))
	for _, cmd := range []*cobra.Command{notificationChannelCreateCmd, notificationChannelEditCmd} {
		cmd.Flags().StringSliceVar(&emailsFlag, "emails", nil, "Email recipients (email channel)")
		cmd.Flags().StringVar(&urlFlag, "url", "", "Webhook URL (webhook channel)")
		cmd.Flags().StringToStringVar(&headersFlag, "header", nil, "Webhook HTTP headers as key=value")
		cmd.Flags().StringSliceVar(&eventTypesFlag, "events", nil, "Event types to notify on (see: netbird get notificationtypes)")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.RegisterFlagCompletionFunc("events", validArgsFunc(notificationTypeNames))
	}

	// EDR
	for _, cmd := range []*cobra.Command{edrCreateCmd, edrEditCmd} {
		cmd.Flags().StringSliceVar(&groupsFlag, "groups", nil, "Groups the compliance check applies to (names or IDs)")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.Flags().StringVar(&clientIDFlag, "client-id", "", "API client ID (intune, falcon)")
		cmd.Flags().StringVar(&tenantIDFlag, "tenant-id", "", "Azure tenant ID (intune)")
		cmd.Flags().StringVar(&secretFlag, "secret", "", "API client secret (intune, falcon)")
		cmd.Flags().StringVar(&upstreamURLFlag, "api-url", "", "API base URL (sentinelone, fleetdm)")
		cmd.Flags().StringVar(&apiTokenFlag, "api-token", "", "API token (sentinelone, fleetdm)")
		cmd.Flags().StringVar(&apiKeyFlag, "api-key", "", "API key (huntress)")
		cmd.Flags().StringVar(&apiSecretFlag, "api-secret", "", "API secret (huntress)")
		cmd.Flags().StringVar(&cloudIDFlag, "cloud-id", "", "CrowdStrike cloud: us-1, us-2, eu-1, ... (falcon)")
		cmd.Flags().IntVar(&ztaThresholdFlag, "zta-threshold", 50, "Minimum Zero Trust Assessment score 0-100 (falcon)")
		cmd.Flags().IntVar(&syncIntervalFlag, "sync-interval", 24, "Device last-sync requirement in hours, min 24")
		cmd.RegisterFlagCompletionFunc("groups", validArgsFunc(groupNames))
	}

	// IdP sync
	for _, cmd := range []*cobra.Command{idpSyncCreateCmd, idpSyncEditCmd} {
		cmd.Flags().StringVar(&customerIDFlag, "customer-id", "", "Google Workspace customer ID (google)")
		cmd.Flags().StringVar(&fileFlag, "service-account-key-file", "", "Google service account JSON key file (google)")
		cmd.Flags().StringVar(&clientIDFlag, "client-id", "", "Application client ID (azure)")
		cmd.Flags().StringVar(&clientSecretFlag, "client-secret", "", "Application client secret (azure)")
		cmd.Flags().StringVar(&tenantIDFlag, "tenant-id", "", "Tenant ID (azure)")
		cmd.Flags().StringVar(&hostFlag, "host", "microsoft.com", "microsoft.com or microsoft.us (azure)")
		cmd.Flags().IntVar(&syncIntervalFlag, "sync-interval", 300, "Sync interval in seconds (google, azure)")
		cmd.Flags().StringSliceVar(&groupPrefixesFlag, "group-prefixes", nil, "Only sync groups starting with these prefixes")
		cmd.Flags().StringSliceVar(&userGroupPrefixesFlag, "user-group-prefixes", nil, "Only sync users of groups starting with these prefixes")
		cmd.Flags().StringVar(&connectorIDFlag, "connector-id", "", "DEX connector ID (embedded IdP setups)")
		cmd.Flags().StringVar(&connectionNameFlag, "connection-name", "", "Connection name (okta)")
		cmd.Flags().StringVar(&providerFlag, "provider", "", "SCIM identity provider name (scim)")
		cmd.Flags().StringVar(&prefixFlag, "prefix", "", "Group name prefix (scim)")
		cmd.Flags().BoolVar(&enabledFlag, "enabled", true, "Enabled")
		cmd.RegisterFlagCompletionFunc("host", staticCompletion([]string{"microsoft.com", "microsoft.us"}))
	}
}
