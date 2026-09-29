package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// NetBird comes as the managed cloud service or a self-hosted management server,
// and each has features the other lacks. Commands are tagged with the edition they
// need; the other edition hides them (help, completion) and refuses to run them.

const (
	editionCloud      = "cloud"
	editionSelfHosted = "selfhosted"

	// annEdition is the cobra annotation holding the edition a command requires.
	annEdition = "netbird-cli/edition"
)

var (
	edition       string // editionCloud or editionSelfHosted
	editionSource string // how edition was decided, for "netbird info"
)

// resolveEdition decides the edition from the configured value, falling back to the API URL:
// NetBird Cloud is served from *.netbird.io, anything else is a self-hosted server.
func resolveEdition(configured, apiURL string) (string, string, error) {
	switch strings.ToLower(strings.TrimSpace(configured)) {
	case "cloud":
		return editionCloud, "set in config / " + "NETBIRD_CLI_EDITION", nil
	case "selfhosted", "self-hosted", "self_hosted", "onprem":
		return editionSelfHosted, "set in config / " + "NETBIRD_CLI_EDITION", nil
	case "", "auto":
	default:
		return "", "", fmt.Errorf("invalid edition %q (expected cloud, selfhosted or auto)", configured)
	}
	host := apiURL
	if u, err := url.Parse(apiURL); err == nil && u.Host != "" {
		host = u.Hostname()
	}
	if host == "netbird.io" || strings.HasSuffix(host, ".netbird.io") {
		return editionCloud, fmt.Sprintf("detected: %s is NetBird Cloud", host), nil
	}
	return editionSelfHosted, fmt.Sprintf("detected: %s is not a netbird.io host", host), nil
}

// requireEdition tags commands as only available on the given edition.
func requireEdition(ed string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations[annEdition] = ed
	}
}

// requiredEdition returns the edition a command (or one of its parents) requires, "" if any.
func requiredEdition(cmd *cobra.Command) string {
	for ; cmd != nil; cmd = cmd.Parent() {
		if ed := cmd.Annotations[annEdition]; ed != "" {
			return ed
		}
	}
	return ""
}

// applyEdition hides the commands the current edition does not support, and parent
// verbs left without any visible subcommand (e.g. "regenerate" on self-hosted).
func applyEdition(root *cobra.Command) {
	var walk func(cmd *cobra.Command) bool
	walk = func(cmd *cobra.Command) bool {
		if ed := cmd.Annotations[annEdition]; ed != "" && ed != edition {
			cmd.Hidden = true
		}
		if !cmd.HasSubCommands() {
			return !cmd.Hidden
		}
		visible := false
		for _, sub := range cmd.Commands() {
			if walk(sub) {
				visible = true
			}
		}
		if !visible && !cmd.Runnable() && cmd != root {
			cmd.Hidden = true
		}
		return !cmd.Hidden
	}
	walk(root)

	if edition == editionCloud {
		// Invites only exist with the embedded IdP of a self-hosted server.
		userCreateCmd.Flags().MarkHidden("invite")
	}
}

// checkEdition stops a command the current edition does not support, with an explanation
// instead of the API's 404.
func checkEdition(cmd *cobra.Command) {
	need := requiredEdition(cmd)
	if edition == editionCloud && cmd == userCreateCmd && inviteFlag {
		need = editionSelfHosted
	}
	if need == "" || need == edition {
		return
	}
	name := strings.TrimPrefix(cmd.CommandPath(), rootCmd.Name()+" ")
	if cmd == userCreateCmd {
		name += " --invite"
	}
	if need == editionCloud {
		fmt.Fprintf(os.Stderr, "%q is only available on NetBird Cloud; %s is a self-hosted server (%s).\n", name, c.URL, editionSource)
	} else {
		fmt.Fprintf(os.Stderr, "%q is only available on self-hosted NetBird (embedded identity provider); %s is NetBird Cloud (%s).\n", name, c.URL, editionSource)
	}
	fmt.Fprintf(os.Stderr, "If the edition is wrong, set \"edition: cloud|selfhosted\" in the config file or NETBIRD_CLI_EDITION.\n")
	os.Exit(1)
}

// cloudHint explains a 404 on self-hosted servers, where optional features (reverse proxy,
// Agent Network, embedded IdP) may simply not be deployed.
func cloudHint(err error) string {
	if !is404(err) || edition != editionSelfHosted {
		return ""
	}
	return "\nhint: this self-hosted server does not expose this endpoint; the feature may be disabled, not deployed, or need a newer management version (netbird info)."
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show the API URL, NetBird edition (cloud / self-hosted) and server version",
	Run: func(cmd *cobra.Command, args []string) {
		info := map[string]interface{}{
			"url":            c.URL,
			"edition":        edition,
			"edition_source": editionSource,
		}
		if v, err := c.GetInstanceVersion(); err == nil {
			info["management_version"] = v.ManagementCurrentVersion
			if v.ManagementUpdateAvailable {
				info["management_update"] = v.ManagementAvailableVersion
			}
		}
		var hidden []string
		for _, verb := range rootCmd.Commands() {
			collectHidden(verb, &hidden)
		}
		if outputFormat != "" {
			info["unavailable_commands"] = hidden
			printOutput(info)
			return
		}
		label := map[string]string{editionCloud: "NetBird Cloud", editionSelfHosted: "self-hosted"}[edition]
		fmt.Printf("URL        : %s\n", c.URL)
		fmt.Printf("Edition    : %s  (%s)\n", bold(label), editionSource)
		if v, ok := info["management_version"]; ok {
			line := fmt.Sprintf("Management : %v", v)
			if u, ok := info["management_update"]; ok {
				line += "  " + yellow(fmt.Sprintf("(update available: %v)", u))
			}
			fmt.Println(line)
		}
		if len(hidden) > 0 {
			other := map[string]string{editionCloud: "self-hosted", editionSelfHosted: "NetBird Cloud"}[edition]
			fmt.Printf("\nOnly on %s, hidden here (%d):\n  %s\n", other, len(hidden), strings.Join(hidden, "\n  "))
		}
	},
}

// collectHidden lists the edition-restricted commands hidden in this edition.
func collectHidden(cmd *cobra.Command, out *[]string) {
	if ed := cmd.Annotations[annEdition]; ed != "" && ed != edition {
		*out = append(*out, strings.TrimPrefix(cmd.CommandPath(), rootCmd.Name()+" "))
		return
	}
	for _, sub := range cmd.Commands() {
		collectHidden(sub, out)
	}
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) { checkEdition(cmd) }

	// From the OpenAPI spec: operations and tags flagged x-cloud-only.
	requireEdition(editionCloud,
		trafficCmd,
		eventStreamsGetCmd, eventStreamCreateCmd, eventStreamEditCmd, eventStreamDeleteCmd,
		notificationChannelsGetCmd, notificationChannelCreateCmd, notificationChannelEditCmd, notificationChannelDeleteCmd, notificationTypesGetCmd,
		edrGetCmd, edrCreateCmd, edrEditCmd, edrDeleteCmd, edrBypassedGetCmd, bypassCmd, unbypassCmd,
		idpSyncsGetCmd, idpSyncCreateCmd, idpSyncEditCmd, idpSyncDeleteCmd, idpSyncLogsGetCmd, syncCmd, idpSyncTokenCmd,
		ingressPeersGetCmd, ingressPeerCreateCmd, ingressPeerEditCmd, ingressPeerDeleteCmd,
		ingressPortsGetCmd, ingressPortCreateCmd, ingressPortEditCmd, ingressPortDeleteCmd,
		tenantsGetCmd, tenantCreateCmd, tenantEditCmd, mspCmd,
		billingCmd,
		// The managed Agent Network gateway is run by NetBird itself.
		agentGatewayGetCmd, agentGatewayCreateCmd,
	)

	// Self-hosted only: instance setup and the embedded IdP (invites).
	requireEdition(editionSelfHosted,
		setupCmd, instanceStatusCmd,
		inviteCmd, invitesGetCmd, invitesDeleteCmd,
	)
}
