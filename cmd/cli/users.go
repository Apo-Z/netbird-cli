package main

import (
	"encoding/json"
	"fmt"

	"netbird-cli/internal/client"

	"github.com/spf13/cobra"
)

var usersGetCmd = &cobra.Command{
	Use:               "users [nom|email|id]",
	Aliases:           []string{"us"},
	Short:             "List or display users",
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			users, err := c.GetUsers()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(users)
			return
		}
		user, err := c.GetUser(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(user)
	},
}

var userCreateCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"usr"},
	Short:   "Create a user or send an invitation (use --invite for invitations)",
	Run: func(cmd *cobra.Command, args []string) {
		if editFlag {
			if inviteFlag {
				result, err := createWithEditor("/api/users/invites", map[string]interface{}{
					"email":       "",
					"name":        "",
					"role":        "user",
					"auto_groups": []string{},
					"expires_in":  259200,
				})
				if err != nil {
					exitErr("create invite", err)
					return
				}
				if result == nil {
					return
				}
				data, err := c.PostRaw("/api/users/invites", result)
				if err != nil {
					exitErr("create invite", err)
					return
				}
				var invite client.Invite
				json.Unmarshal(data, &invite)
				fmt.Println("invitation created:")
				printOutput(invite)
				return
			}
			result, err := createWithEditor("/api/users", map[string]interface{}{
				"email":          "",
				"name":           "",
				"role":           "user",
				"auto_groups":    []string{},
				"is_service_user": false,
			})
			if err != nil {
				exitErr("create user", err)
				return
			}
			if result == nil {
				return
			}
			data, err := c.PostRaw("/api/users", result)
			if err != nil {
				exitErr("create user", err)
				return
			}
			var user client.User
			json.Unmarshal(data, &user)
			fmt.Println("user created:")
			printOutput(user)
			return
		}
		if inviteFlag {
			req := &client.CreateInviteRequest{
				Email:      emailFlag,
				Name:       nameFlag,
				Role:       roleFlag,
				AutoGroups: autoGroupsFlag,
				ExpiresIn:  inviteExpireFlag,
			}
			if dryRunCheck(req) {
				return
			}
			invite, err := c.CreateUserInvite(req)
			if err != nil {
				printErr(err)
				return
			}
			printOutput(invite)
			return
		}
		req := &client.CreateUserRequest{
			Email:         emailFlag,
			Name:          nameFlag,
			Role:          roleFlag,
			AutoGroups:    autoGroupsFlag,
			IsServiceUser: isServiceUserFlag,
		}
		if dryRunCheck(req) {
			return
		}
		user, err := c.CreateUser(req)
		if err != nil {
			printErr(err)
			return
		}
		fmt.Println("user created:")
		printOutput(user)
	},
}

var userEditCmd = &cobra.Command{
	Use:               "user <nom|email|id>",
	Aliases:           []string{"usr"},
	Short:             "Edit a user",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := c.GetUser(args[0])
		if err != nil {
			exitErr("user", err)
			return
		}
		id := user.ID
		overrides := map[string]interface{}{}
		if cmd.Flags().Changed("role") {
			overrides["role"] = roleFlag
		}
		if cmd.Flags().Changed("auto-groups") {
			overrides["auto_groups"] = autoGroupsFlag
		}
		if cmd.Flags().Changed("blocked") {
			overrides["is_blocked"] = isBlockedFlag
		}
		if cmd.Flags().Changed("name") {
			overrides["name"] = nameFlag
		}
		data, _ := json.Marshal(user)
		result, err := editWithEditor(fmt.Sprintf("/api/users/%s", id), data, overrides)
		if err != nil {
			exitErr("edit user", err)
			return
		}
		if result != nil {
			printOutput(result)
		}
	},
}

var userDeleteCmd = &cobra.Command{
	Use:               "user <nom|email|id>",
	Aliases:           []string{"usr"},
	Short:             "Delete a user",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolveUserID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete user %s", args[0])) {
			return
		}
		if err := c.DeleteUser(id); err != nil {
			printErr(err)
			return
		}
		fmt.Println("user deleted")
	},
}

var approveCmd = &cobra.Command{
	Use:               "approve user <nom|email|id>",
	Short:             "Approve a pending user",
	Args:              cobra.MinimumNArgs(2),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "user" && args[0] != "usr" {
			fmt.Println("usage: netbird approve user <name|email|id>")
			return
		}
		id, err := c.ResolveUserID(args[1])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would approve user %s", args[1])) {
			return
		}
		user, err := c.ApproveUser(id)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(user)
	},
}

var blockCmd = &cobra.Command{
	Use:               "block user <nom|email|id>",
	Short:             "Block a user",
	Args:              cobra.MinimumNArgs(2),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "user" && args[0] != "usr" {
			fmt.Println("usage: netbird block user <name|email|id>")
			return
		}
		id, err := c.ResolveUserID(args[1])
		if err != nil {
			printErr(err)
			return
		}
		blocked := true
		if dryRunCheck(&client.UpdateUserRequest{IsBlocked: &blocked}) {
			return
		}
		req := &client.UpdateUserRequest{IsBlocked: &blocked}
		user, err := c.UpdateUser(id, req)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(user)
	},
}

var unblockCmd = &cobra.Command{
	Use:               "unblock user <nom|email|id>",
	Short:             "Unblock a user",
	Args:              cobra.MinimumNArgs(2),
	ValidArgsFunction: validArgsFunc(userNames),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "user" && args[0] != "usr" {
			fmt.Println("usage: netbird unblock user <name|email|id>")
			return
		}
		id, err := c.ResolveUserID(args[1])
		if err != nil {
			printErr(err)
			return
		}
		blocked := false
		if dryRunCheck(&client.UpdateUserRequest{IsBlocked: &blocked}) {
			return
		}
		req := &client.UpdateUserRequest{IsBlocked: &blocked}
		user, err := c.UpdateUser(id, req)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(user)
	},
}

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Manage invitations",
}

var inviteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create invitation",
	Run: func(cmd *cobra.Command, args []string) {
		req := &client.CreateInviteRequest{
			Email:      emailFlag,
			Name:       nameFlag,
			Role:       roleFlag,
			AutoGroups: autoGroupsFlag,
			ExpiresIn:  inviteExpireFlag,
		}
		if dryRunCheck(req) {
			return
		}
		invite, err := c.CreateUserInvite(req)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(invite)
	},
}

var inviteDeleteCmd = &cobra.Command{
	Use:               "delete <name|email|id>",
	Short:             "Delete an invitation",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(inviteNames),
	Run: func(cmd *cobra.Command, args []string) {
		invite, err := c.GetInvite(args[0])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete invitation %s", args[0])) {
			return
		}
		if err := c.DeleteUserInvite(invite.ID); err != nil {
			printErr(err)
			return
		}
		fmt.Println("invitation deleted")
	},
}

var inviteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List invitations",
	Run: func(cmd *cobra.Command, args []string) {
		invites, err := c.GetUserInvites()
		if err != nil {
			printErr(err)
			return
		}
		printOutput(invites)
	},
}

var inviteRegenerateCmd = &cobra.Command{
	Use:               "regenerate <name|email|id>",
	Aliases:           []string{"regen"},
	Short:             "Regenerate an invitation token",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(inviteNames),
	Run: func(cmd *cobra.Command, args []string) {
		invite, err := c.GetInvite(args[0])
		if err != nil {
			printErr(err)
			return
		}
		req := &client.RegenerateInviteRequest{ExpiresIn: inviteExpireFlag}
		if dryRunCheck(req) {
			return
		}
		regenerated, err := c.RegenerateUserInvite(invite.ID, req)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(regenerated)
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := c.GetCurrentUser()
		if err != nil {
			printErr(err)
			return
		}
		printOutput(user)
	},
}

var invitesGetCmd = &cobra.Command{
	Use:               "invites [name|email|id]",
	Aliases:           []string{"inv"},
	Short:             "List or display invitations",
	ValidArgsFunction: validArgsFunc(inviteNames),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			invites, err := c.GetUserInvites()
			if err != nil {
				printErr(err)
				return
			}
			printOutput(invites)
			return
		}
		invite, err := c.GetInvite(args[0])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(invite)
	},
}

var invitesDeleteCmd = &cobra.Command{
	Use:               "invites <name|email|id>",
	Aliases:           []string{"inv"},
	Short:             "Delete an invitation",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(inviteNames),
	Run: func(cmd *cobra.Command, args []string) {
		invite, err := c.GetInvite(args[0])
		if err != nil {
			printErr(err)
			return
		}
		if dryRunMsg(fmt.Sprintf("would delete invitation %s", args[0])) {
			return
		}
		if err := c.DeleteUserInvite(invite.ID); err != nil {
			printErr(err)
			return
		}
		fmt.Println("invitation deleted")
	},
}

func init() {
	getCmd.AddCommand(usersGetCmd)
	getCmd.AddCommand(invitesGetCmd)
	createCmd.AddCommand(userCreateCmd)
	editCmd.AddCommand(userEditCmd)
	deleteCmd.AddCommand(userDeleteCmd)
	deleteCmd.AddCommand(invitesDeleteCmd)
	rootCmd.AddCommand(approveCmd, blockCmd, unblockCmd, inviteCmd, whoamiCmd)
	inviteCmd.AddCommand(inviteListCmd, inviteCreateCmd, inviteDeleteCmd, inviteRegenerateCmd)

	userCreateCmd.Flags().StringVar(&emailFlag, "email", "", "Email")
	userCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Full name")
	userCreateCmd.Flags().StringVar(&roleFlag, "role", "user", "Role (admin/user)")
	userCreateCmd.Flags().StringSliceVar(&autoGroupsFlag, "auto-groups", nil, "IDs or names of auto-assigned groups")
	userCreateCmd.Flags().BoolVar(&isServiceUserFlag, "service", false, "Service user")
	userCreateCmd.Flags().BoolVar(&inviteFlag, "invite", false, "Send an invitation instead of creating directly")
	userCreateCmd.Flags().IntVar(&inviteExpireFlag, "expire", 259200, "Invitation expiration in seconds (only with --invite)")

	userEditCmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	userEditCmd.Flags().StringVar(&roleFlag, "role", "", "New role")
	userEditCmd.Flags().StringSliceVar(&autoGroupsFlag, "auto-groups", nil, "IDs or names of auto-assigned groups")
	userEditCmd.Flags().BoolVar(&isBlockedFlag, "blocked", false, "Block user")

	inviteCreateCmd.Flags().StringVar(&emailFlag, "email", "", "Email")
	inviteCreateCmd.Flags().StringVar(&nameFlag, "name", "", "Full name")
	inviteCreateCmd.Flags().StringVar(&roleFlag, "role", "user", "Role")
	inviteCreateCmd.Flags().StringSliceVar(&autoGroupsFlag, "auto-groups", nil, "IDs or names of groups")
	inviteCreateCmd.Flags().IntVar(&inviteExpireFlag, "expire", 259200, "Expiration in seconds")

	inviteRegenerateCmd.Flags().IntVar(&inviteExpireFlag, "expire", 259200, "New expiration in seconds")

	userCreateCmd.RegisterFlagCompletionFunc("role", staticCompletion([]string{"admin", "user"}))
	userCreateCmd.RegisterFlagCompletionFunc("auto-groups", validArgsFunc(groupNames))
	userEditCmd.RegisterFlagCompletionFunc("role", staticCompletion([]string{"admin", "user"}))
	userEditCmd.RegisterFlagCompletionFunc("auto-groups", validArgsFunc(groupNames))
	inviteCreateCmd.RegisterFlagCompletionFunc("role", staticCompletion([]string{"admin", "user"}))
	inviteCreateCmd.RegisterFlagCompletionFunc("auto-groups", validArgsFunc(groupNames))
}
