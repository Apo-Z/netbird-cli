package main

import (
	"fmt"
	"os"
	"strings"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

// MSP tenants and billing (NetBird Cloud).

// --- Tenants ---

var tenantsGetCmd = &cobra.Command{
	Use:               "tenants [name|domain|id]",
	Aliases:           []string{"tenant", "tn"},
	Short:             "List or display MSP tenants (cloud)",
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		tenants, err := c.GetTenants()
		if err != nil {
			exitErr("tenants", err)
			return
		}
		if len(args) == 0 {
			printOutput(tenants)
			return
		}
		for _, t := range tenants {
			if t.ID == args[0] || t.Name == args[0] || t.Domain == args[0] {
				printOutput(t)
				if t.Status == "pending" && t.DNSChallenge != "" && outputFormat == "" {
					fmt.Printf("\nTo verify the domain, add a TXT record on %s with:\n  %s\nthen run: netbird msp verify-dns %q\n", t.Domain, t.DNSChallenge, t.Name)
				}
				return
			}
		}
		exitErr(fmt.Sprintf("tenant %s not found", args[0]), nil)
	},
}

var tenantCreateCmd = &cobra.Command{
	Use:     "tenant",
	Aliases: []string{"tn"},
	Short:   "Create an MSP tenant (cloud)",
	Long: `Create an MSP tenant. --groups lists the MSP groups that can manage the tenant
and the role they assume there, as <group>:<role> (role: admin, user, ...).`,
	Example: `  netbird create tenant --name "Acme Corp" --domain acme.com --groups msp-admins:admin,msp-support:user`,
	Run: func(cmd *cobra.Command, args []string) {
		if nameFlag == "" || domainFlag == "" || len(tenantGroupsFlag) == 0 {
			exitErr("--name, --domain and --groups are required", nil)
			return
		}
		groups, err := parseTenantGroups(tenantGroupsFlag)
		if err != nil {
			exitErr("--groups", err)
			return
		}
		req := &client.TenantRequest{Name: nameFlag, Domain: domainFlag, Groups: groups}
		if dryRunCheck(req) {
			return
		}
		t, err := c.CreateTenant(req)
		if err != nil {
			exitErr("create tenant", err)
			return
		}
		fmt.Println("tenant created:")
		printOutput(t)
		if t.DNSChallenge != "" && outputFormat == "" {
			fmt.Printf("\nAdd a TXT record on %s with:\n  %s\nthen run: netbird msp verify-dns %q\n", t.Domain, t.DNSChallenge, t.Name)
		}
	},
}

var tenantEditCmd = &cobra.Command{
	Use:               "tenant <name|domain|id>",
	Aliases:           []string{"tn"},
	Short:             "Edit an MSP tenant (name, managing groups and roles)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveTenantID(args[0])
		if err != nil {
			exitErr("tenant", err)
			return
		}
		tenants, _ := c.GetTenants()
		var current *client.Tenant
		for i := range tenants {
			if tenants[i].ID == id {
				current = &tenants[i]
			}
		}
		if current == nil {
			exitErr(fmt.Sprintf("tenant %s not found", args[0]), nil)
			return
		}
		// No GET /tenants/{id}: build the request from the list and apply the flags.
		req := &client.TenantRequest{Name: current.Name, Groups: current.Groups}
		if cmd.Flags().Changed("name") {
			req.Name = nameFlag
		}
		if cmd.Flags().Changed("groups") {
			groups, err := parseTenantGroups(tenantGroupsFlag)
			if err != nil {
				exitErr("--groups", err)
				return
			}
			req.Groups = groups
		}
		if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("groups") {
			exitErr("nothing to change: pass --name and/or --groups", nil)
			return
		}
		if dryRunCheck(req) {
			return
		}
		t, err := c.UpdateTenant(id, req)
		if err != nil {
			exitErr("edit tenant", err)
			return
		}
		printOutput(t)
	},
}

var mspCmd = &cobra.Command{
	Use:   "msp",
	Short: "MSP tenant operations (cloud)",
}

var mspVerifyDNSCmd = &cobra.Command{
	Use:               "verify-dns <tenant>",
	Short:             "Verify the tenant domain DNS challenge (TXT record)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		tenantAction(args[0], "POST", "dns", nil, "DNS challenge verified")
	},
}

var mspInviteCmd = &cobra.Command{
	Use:               "invite <tenant>",
	Short:             "Invite an existing NetBird account to become a tenant of this MSP",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		tenantAction(args[0], "POST", "invite", nil, "invitation sent")
	},
}

var mspAcceptCmd = &cobra.Command{
	Use:               "accept <tenant-id>",
	ValidArgsFunction: noCompletion,
	Short:             "Accept an MSP invitation (run as the invited account owner)",
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantActionByID(args[0], "PUT", "invite", map[string]string{"value": "accept"}, "invitation accepted")
	},
}

var mspDeclineCmd = &cobra.Command{
	Use:               "decline <tenant-id>",
	ValidArgsFunction: noCompletion,
	Short:             "Decline an MSP invitation (run as the invited account owner)",
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tenantActionByID(args[0], "PUT", "invite", map[string]string{"value": "decline"}, "invitation declined")
	},
}

var mspUnlinkCmd = &cobra.Command{
	Use:               "unlink <tenant> --owner <user-id>",
	Short:             "Unlink a tenant from the MSP, handing it over to a new owner",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		if ownerFlag == "" {
			exitErr("--owner is required (user ID of the tenant's new owner)", nil)
			return
		}
		tenantAction(args[0], "POST", "unlink", map[string]string{"owner": ownerFlag}, "tenant unlinked")
	},
}

var mspSubscribeCmd = &cobra.Command{
	Use:               "subscribe <tenant> --price <price-id>",
	Short:             "Create a subscription for a tenant (see: netbird billing plans)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(tenantNames),
	Run: func(cmd *cobra.Command, args []string) {
		if priceFlag == "" {
			exitErr("--price is required", nil)
			return
		}
		tenantAction(args[0], "POST", "subscription", map[string]string{"priceID": priceFlag}, "tenant subscription created")
	},
}

func tenantAction(nameOrID, method, action string, body interface{}, done string) {
	id, err := c.ResolveTenantID(nameOrID)
	if err != nil {
		exitErr("tenant", err)
		return
	}
	tenantActionByID(id, method, action, body, done)
}

func tenantActionByID(id, method, action string, body interface{}, done string) {
	if dryRun {
		fmt.Printf("[dry-run] would %s %s/%s/%s\n", method, "/api/integrations/msp/tenants", id, action)
		if body != nil {
			printJSON(body)
		}
		return
	}
	if _, err := c.TenantAction(method, id, action, body); err != nil {
		exitErr("msp "+action, err)
		return
	}
	fmt.Println(done)
}

// parseTenantGroups parses group:role pairs, resolving group names.
func parseTenantGroups(pairs []string) ([]client.TenantGroup, error) {
	var groups []client.TenantGroup
	for _, p := range pairs {
		g, role, ok := strings.Cut(p, ":")
		if !ok || role == "" {
			return nil, fmt.Errorf("%q: expected <group>:<role>", p)
		}
		id, err := c.ResolveGroupID(g)
		if err != nil {
			return nil, err
		}
		groups = append(groups, client.TenantGroup{ID: id, Role: role})
	}
	return groups, nil
}

// --- Billing ---

var billingCmd = &cobra.Command{
	Use:   "billing",
	Short: "Subscription, plans, usage and invoices (cloud)",
}

var billingUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Current billable usage (users, peers)",
	Run: func(cmd *cobra.Command, args []string) {
		u, err := c.GetBillingUsage()
		if err != nil {
			exitErr("billing usage", err)
			return
		}
		printOutput(u)
	},
}

var billingSubscriptionCmd = &cobra.Command{
	Use:     "subscription",
	Aliases: []string{"sub"},
	Short:   "Show the current subscription, or change it with --price / --plan",
	Run: func(cmd *cobra.Command, args []string) {
		if priceFlag != "" || planFlag != "" {
			if dryRunCheck(map[string]string{"priceID": priceFlag, "plan_tier": planFlag}) {
				return
			}
			if err := c.ChangeBillingSubscription(priceFlag, planFlag); err != nil {
				exitErr("change subscription", err)
				return
			}
			fmt.Println("subscription changed")
			return
		}
		s, err := c.GetBillingSubscription()
		if err != nil {
			exitErr("subscription", err)
			return
		}
		if outputFormat != "" {
			printOutput(s)
			return
		}
		printOutput(s)
		fmt.Printf("\n  price:                  %s / %s\n", formatMinorUnits(s.Price, s.Currency), s.PlanTier)
		if s.RemainingTrial > 0 {
			fmt.Printf("  trial remaining:        %s\n", formatWindow(int64(s.RemainingTrial)))
		}
	},
}

type billingPlanRow struct {
	Name   string `json:"name"`
	Free   bool   `json:"free"`
	Prices string `json:"prices"`
}

var billingPlansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Available plans and their price IDs",
	Run: func(cmd *cobra.Command, args []string) {
		plans, err := c.GetBillingPlans()
		if err != nil {
			exitErr("plans", err)
			return
		}
		if outputFormat != "" {
			printOutput(plans)
			return
		}
		var rows []billingPlanRow
		for _, p := range plans {
			var prices []string
			for _, pr := range p.Prices {
				prices = append(prices, fmt.Sprintf("%s/%s (%s)", formatMinorUnits(pr.Price, pr.Currency), pr.Unit, pr.PriceID))
			}
			rows = append(rows, billingPlanRow{p.Name, p.Free, strings.Join(prices, ", ")})
		}
		printOutput(rows)
	},
}

type billingInvoiceRow struct {
	InvoiceID   string `json:"invoice_id"`
	Type        string `json:"type"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
}

var billingInvoicesCmd = &cobra.Command{
	Use:   "invoices",
	Short: "List paid invoices",
	Run: func(cmd *cobra.Command, args []string) {
		invoices, err := c.GetBillingInvoices()
		if err != nil {
			exitErr("invoices", err)
			return
		}
		if outputFormat != "" {
			printOutput(invoices)
			return
		}
		var rows []billingInvoiceRow
		for _, i := range invoices {
			rows = append(rows, billingInvoiceRow{i.ID, i.Type, i.PeriodStart, i.PeriodEnd})
		}
		printOutput(rows)
	},
}

var billingInvoiceCmd = &cobra.Command{
	Use:               "invoice <id>",
	ValidArgsFunction: validArgsFunc(invoiceIDs),
	Short:             "Get an invoice: PDF link (default) or CSV (--csv, written to stdout or --file)",
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if csvFlag {
			data, err := c.GetBillingInvoiceCSV(args[0])
			if err != nil {
				exitErr("invoice csv", err)
				return
			}
			if fileFlag != "" {
				if err := os.WriteFile(fileFlag, data, 0o644); err != nil {
					exitErr("write", err)
					return
				}
				fmt.Printf("invoice written to %s\n", fileFlag)
				return
			}
			os.Stdout.Write(data)
			return
		}
		u, err := c.GetBillingInvoicePDF(args[0])
		if err != nil {
			exitErr("invoice pdf", err)
			return
		}
		fmt.Println(u.URL)
	},
}

var billingPortalCmd = &cobra.Command{
	Use:   "portal",
	Short: "Get a link to the customer billing portal (payment methods, invoices)",
	Run: func(cmd *cobra.Command, args []string) {
		u, err := c.GetBillingPortal(portalReturnFlag)
		if err != nil {
			exitErr("billing portal", err)
			return
		}
		fmt.Println(u.URL)
	},
}

var billingCheckoutCmd = &cobra.Command{
	Use:   "checkout --price <price-id>",
	Short: "Start a checkout session for a plan and print its URL",
	Run: func(cmd *cobra.Command, args []string) {
		if priceFlag == "" {
			exitErr("--price is required (see: netbird billing plans)", nil)
			return
		}
		if dryRunCheck(map[string]interface{}{"baseURL": baseURLFlag, "priceID": priceFlag, "enableTrial": trialFlag}) {
			return
		}
		u, err := c.CreateBillingCheckout(baseURLFlag, priceFlag, trialFlag)
		if err != nil {
			exitErr("checkout", err)
			return
		}
		fmt.Println(u.URL)
	},
}

var billingAWSCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS Marketplace subscription",
}

var billingAWSActivateCmd = &cobra.Command{
	Use:   "activate --plan <tier>",
	Short: "Activate an AWS Marketplace subscription",
	Run: func(cmd *cobra.Command, args []string) {
		if planFlag == "" {
			exitErr("--plan is required", nil)
			return
		}
		if dryRunMsg("would activate the AWS Marketplace subscription for plan " + planFlag) {
			return
		}
		if err := c.BillingAWSMarketplace("activate", map[string]string{"plan_tier": planFlag}); err != nil {
			exitErr("aws activate", err)
			return
		}
		fmt.Println("AWS Marketplace subscription activated")
	},
}

var billingAWSEnrichCmd = &cobra.Command{
	Use:   "enrich --aws-user-id <id>",
	Short: "Link the AWS Marketplace subscription to an AWS user ID",
	Run: func(cmd *cobra.Command, args []string) {
		if awsUserIDFlag == "" {
			exitErr("--aws-user-id is required", nil)
			return
		}
		if dryRunMsg("would link the AWS Marketplace subscription to " + awsUserIDFlag) {
			return
		}
		if err := c.BillingAWSMarketplace("enrich", map[string]string{"aws_user_id": awsUserIDFlag}); err != nil {
			exitErr("aws enrich", err)
			return
		}
		fmt.Println("AWS Marketplace subscription linked")
	},
}

// formatMinorUnits renders an amount in cents as e.g. 5.00 USD.
func formatMinorUnits(amount int, currency string) string {
	return fmt.Sprintf("%.2f %s", float64(amount)/100, strings.ToUpper(currency))
}

func init() {
	getCmd.AddCommand(tenantsGetCmd)
	createCmd.AddCommand(tenantCreateCmd)
	editCmd.AddCommand(tenantEditCmd)
	mspCmd.AddCommand(mspVerifyDNSCmd, mspInviteCmd, mspAcceptCmd, mspDeclineCmd, mspUnlinkCmd, mspSubscribeCmd)
	billingAWSCmd.AddCommand(billingAWSActivateCmd, billingAWSEnrichCmd)
	billingCmd.AddCommand(billingUsageCmd, billingSubscriptionCmd, billingPlansCmd, billingInvoicesCmd, billingInvoiceCmd,
		billingPortalCmd, billingCheckoutCmd, billingAWSCmd)
	rootCmd.AddCommand(mspCmd, billingCmd)

	for _, cmd := range []*cobra.Command{tenantCreateCmd, tenantEditCmd} {
		cmd.Flags().StringVar(&nameFlag, "name", "", "Tenant name")
		cmd.Flags().StringSliceVar(&tenantGroupsFlag, "groups", nil, "Managing MSP groups as <group>:<role>")
		cmd.RegisterFlagCompletionFunc("groups", tenantGroupCompletion)
	}
	tenantCreateCmd.Flags().StringVar(&domainFlag, "domain", "", "Tenant domain (verified by DNS TXT challenge)")
	mspUnlinkCmd.Flags().StringVar(&ownerFlag, "owner", "", "User ID of the tenant's new owner")
	mspSubscribeCmd.Flags().StringVar(&priceFlag, "price", "", "Price ID")

	billingSubscriptionCmd.Flags().StringVar(&priceFlag, "price", "", "Change to this price ID")
	billingSubscriptionCmd.Flags().StringVar(&planFlag, "plan", "", "Change to this plan tier")
	billingInvoiceCmd.Flags().BoolVar(&csvFlag, "csv", false, "Download the CSV instead of printing the PDF link")
	billingInvoiceCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Write the CSV to this file")
	billingPortalCmd.Flags().StringVar(&portalReturnFlag, "return-url", "https://app.netbird.io", "URL to return to after leaving the portal")
	billingCheckoutCmd.Flags().StringVar(&priceFlag, "price", "", "Price ID (see: netbird billing plans)")
	billingCheckoutCmd.Flags().StringVar(&baseURLFlag, "return-url", "https://app.netbird.io/plans/success", "URL to return to after checkout")
	billingCheckoutCmd.Flags().BoolVar(&trialFlag, "trial", false, "Start with a 14-day trial")
	billingAWSActivateCmd.Flags().StringVar(&planFlag, "plan", "", "Plan tier, e.g. business")
	billingAWSEnrichCmd.Flags().StringVar(&awsUserIDFlag, "aws-user-id", "", "AWS user ID")

	mspUnlinkCmd.RegisterFlagCompletionFunc("owner", validArgsFunc(userIDs))
	for _, cmd := range []*cobra.Command{mspSubscribeCmd, billingSubscriptionCmd, billingCheckoutCmd} {
		cmd.RegisterFlagCompletionFunc("price", validArgsFunc(priceIDs))
	}
	for _, cmd := range []*cobra.Command{billingSubscriptionCmd, billingAWSActivateCmd} {
		cmd.RegisterFlagCompletionFunc("plan", validArgsFunc(planTiers))
	}
}
