package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/textonlyio/textonly-cli/internal/config"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Manage configuration"}

	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Path())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			val, ok := config.Get(args[0])
			if !ok {
				return fmt.Errorf("not found: %s", args[0])
			}
			fmt.Println(val)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a config value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.Set(args[0], args[1])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "unset <key>",
		Short: "Unset a config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.Unset(args[0])
		},
	})

	return cmd
}
