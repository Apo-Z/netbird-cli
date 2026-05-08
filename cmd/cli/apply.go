package main

import (
	"fmt"
	"os"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply declarative configuration (create/update)",
	Run: func(cmd *cobra.Command, args []string) {
		if fileFlag == "" {
			fmt.Fprintln(os.Stderr, "usage: netbird apply -f <file.yaml>")
			return
		}
		data, err := os.ReadFile(fileFlag)
		if err != nil {
			exitErr("file read", err)
			return
		}

		var spec struct {
			Groups    []applyGroup    `yaml:"groups"`
			Users     []applyUser     `yaml:"users"`
			Peers     []applyPeer     `yaml:"peers"`
			Policies  []applyPolicy   `yaml:"policies"`
		SetupKeys []applySetupKey `yaml:"setupkeys"`
		Accounts  []applyAccount  `yaml:"accounts"`
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
			exitErr("YAML parsing", err)
			return
		}

		// Groups first (dependencies)
		for _, g := range spec.Groups {
			resolvedPeers := resolveAll("peer", g.Peers)
			existing, _ := c.GetGroupByName(g.Name)
			if existing != nil {
				if dryRun {
					fmt.Printf("[dry-run] would update group %s (id=%s)\n", g.Name, existing.ID)
					fmt.Printf("  peers: %v\n", resolvedPeers)
					continue
				}
				_, err := c.UpdateGroup(existing.ID, &client.UpdateGroupRequest{Name: g.Name, Peers: resolvedPeers})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error updating group %s: %s\n", g.Name, err)
				} else {
					fmt.Printf("group %s updated\n", g.Name)
				}
			} else {
				if dryRun {
					fmt.Printf("[dry-run] would create group %s\n", g.Name)
					fmt.Printf("  peers: %v\n", resolvedPeers)
					continue
				}
				_, err := c.CreateGroup(&client.CreateGroupRequest{Name: g.Name, Peers: resolvedPeers})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error creating group %s: %s\n", g.Name, err)
				} else {
					fmt.Printf("group %s created\n", g.Name)
				}
			}
		}

		// Users
		for _, u := range spec.Users {
			resolvedGroups := resolveAll("group", u.AutoGroups)
			existing, _ := c.GetUserByEmail(u.Email)
			if existing != nil {
				if dryRun {
					fmt.Printf("[dry-run] would update user %s\n", u.Email)
					continue
				}
				_, err := c.UpdateUser(existing.ID, &client.UpdateUserRequest{Role: u.Role, AutoGroups: resolvedGroups})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error updating user %s: %s\n", u.Email, err)
				} else {
					fmt.Printf("user %s updated\n", u.Email)
				}
			} else {
				if dryRun {
					fmt.Printf("[dry-run] would create user %s (%s)\n", u.Name, u.Email)
					continue
				}
				_, err := c.CreateUser(&client.CreateUserRequest{
					Email: u.Email, Name: u.Name, Role: u.Role,
					AutoGroups: resolvedGroups, IsServiceUser: u.IsServiceUser,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error creating user %s: %s\n", u.Email, err)
				} else {
					fmt.Printf("user %s created\n", u.Email)
				}
			}
		}

		// Peers
		for _, p := range spec.Peers {
			id, err := c.ResolvePeerID(p.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "peer %s not found\n", p.Name)
				continue
			}
			if dryRun {
				fmt.Printf("[dry-run] would update peer %s\n", p.Name)
				continue
			}
			_, err = c.UpdatePeer(id, &client.UpdatePeerRequest{
				Name: p.Name, SSHEnabled: p.SSH,
				LoginExpirationEnabled: p.LoginExp, InactivityExpirationEnabled: p.InactivityExp,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "error updating peer %s: %s\n", p.Name, err)
			} else {
				fmt.Printf("peer %s updated\n", p.Name)
			}
		}

		// Policies
		for _, pol := range spec.Policies {
			id, err := c.ResolvePolicyID(pol.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "policy %s not found\n", pol.Name)
				continue
			}
			if dryRun {
				fmt.Printf("[dry-run] would update policy %s (enabled=%v)\n", pol.Name, pol.Enabled)
				continue
			}
			policy, _ := c.GetPolicyByID(id)
			if policy != nil {
				policy.Enabled = pol.Enabled
				// PUT via raw to avoid typed struct complexity
				c.PutRaw(fmt.Sprintf("/api/policies/%s", id), policy)
				fmt.Printf("policy %s updated\n", pol.Name)
			}
		}

		// SetupKeys
		for _, sk := range spec.SetupKeys {
			resolvedGroups := resolveAll("group", sk.AutoGroups)
			existing, _ := c.GetSetupKeyByName(sk.Name)
			if existing != nil {
				if dryRun {
					fmt.Printf("[dry-run] would update key %s\n", sk.Name)
					continue
				}
				c.UpdateSetupKey(fmt.Sprintf("%d", existing.ID), &client.UpdateSetupKeyRequest{AutoGroups: resolvedGroups})
				fmt.Printf("key %s updated\n", sk.Name)
			} else {
				if dryRun {
					fmt.Printf("[dry-run] would create key %s (type=%s, expire=%d)\n", sk.Name, sk.Type, sk.ExpiresIn)
					continue
				}
				c.CreateSetupKey(&client.CreateSetupKeyRequest{
					Name: sk.Name, Type: sk.Type, ExpiresIn: sk.ExpiresIn,
					AutoGroups: resolvedGroups, UsageLimit: sk.UsageLimit,
				})
				fmt.Printf("key %s created\n", sk.Name)
			}
		}
		// Accounts
		for _, acct := range spec.Accounts {
			accounts, err := c.GetAccounts()
			if err != nil || len(accounts) == 0 {
				fmt.Fprintf(os.Stderr, "no account found\n")
				continue
			}
			if dryRun {
				fmt.Printf("[dry-run] would update account %s\n", accounts[0].ID)
				continue
			}
			c.PutRaw(fmt.Sprintf("/api/accounts/%s", accounts[0].ID), acct)
			fmt.Println("account updated")
		}
	},
}

type applyGroup struct {
	Name  string   `yaml:"name"`
	Peers []string `yaml:"peers"`
}

type applyUser struct {
	Name         string   `yaml:"name"`
	Email        string   `yaml:"email"`
	Role         string   `yaml:"role"`
	AutoGroups   []string `yaml:"auto_groups"`
	IsServiceUser bool    `yaml:"is_service_user"`
}

type applyPeer struct {
	Name         string `yaml:"name"`
	SSH          bool   `yaml:"ssh"`
	LoginExp     bool   `yaml:"login_expiration"`
	InactivityExp bool  `yaml:"inactivity_expiration"`
}

type applyPolicy struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

type applySetupKey struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	ExpiresIn  int      `yaml:"expires_in"`
	AutoGroups []string `yaml:"auto_groups"`
	UsageLimit int      `yaml:"usage_limit"`
}

type applyAccount struct {
	Settings *client.AccountSettings `yaml:"settings" json:"settings"`
}

func resolveAll(resourceType string, names []string) []string {
	var result []string
	for _, n := range names {
		var id string
		var err error
		switch resourceType {
		case "group":
			id, err = c.ResolveGroupID(n)
		case "user":
			id, err = c.ResolveUserID(n)
		case "peer":
			id, err = c.ResolvePeerID(n)
		default:
			id = n
		}
		if err != nil {
			id = n
		}
		result = append(result, id)
	}
	return result
}

func init() {
	applyCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "YAML configuration file")
	rootCmd.AddCommand(applyCmd)
}
