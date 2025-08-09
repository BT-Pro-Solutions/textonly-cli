package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/textonlyio/textonly-cli/internal/auth"
)

func newAuthCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "Authentication"}
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	cmd.AddCommand(newWhoAmICommand())
	return cmd
}

func newLoginCommand() *cobra.Command {
	var noOpen bool
	c := &cobra.Command{
		Use:   "login",
		Short: "Log in via magic link",
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.Login(noOpen)
		},
	}
	c.Flags().BoolVar(&noOpen, "no-open", false, "Do not open the browser automatically")
	return c
}

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out and revoke token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.Logout(); err != nil {
				return err
			}
			fmt.Println("logged out")
			return nil
		},
	}
}

func newWhoAmICommand() *cobra.Command {
	var asJSON bool
	c := &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.WhoAmI(asJSON)
		},
	}
	c.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return c
}
