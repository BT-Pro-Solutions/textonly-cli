package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/textonlyio/textonly-cli/internal/api"
	"github.com/textonlyio/textonly-cli/internal/config"
	"github.com/textonlyio/textonly-cli/internal/update"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to",
		Short: "TextOnly CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			api.SetUserAgent(fmt.Sprintf("to/%s (%s/%s)", version, runtime.GOOS, runtime.GOARCH))
			return config.Init()
		},
	}

	cmd.PersistentFlags().String("api", "", "Override API base URL (TO_API)")
	_ = viper.BindPFlag("api", cmd.PersistentFlags().Lookup("api"))

	// Top-level auth commands
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newLogoutCommand())
	cmd.AddCommand(newWhoAmICommand())

	// Grouped commands
	cmd.AddCommand(newAuthCommand())
	cmd.AddCommand(newNotesCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newCompletionCommand())
	cmd.AddCommand(newUpdateCommand())
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newDoctorCommand())

	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("to %s (%s %s %s)\n", version, runtime.GOOS, runtime.GOARCH, commit)
		},
	}
}

func newUpdateCommand() *cobra.Command {
	var checkOnly bool
	c := &cobra.Command{
		Use:   "update",
		Short: "Self-update the CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			newVer, updated, err := update.CheckAndApply(checkOnly, version)
			if err != nil {
				return err
			}
			if checkOnly {
				if updated {
					fmt.Printf("update available: %s\n", newVer)
				} else {
					fmt.Println("already up to date")
				}
				return nil
			}
			if updated {
				fmt.Printf("updated to %s\n", newVer)
			} else {
				fmt.Println("already up to date")
			}
			return nil
		},
	}
	c.Flags().BoolVar(&checkOnly, "check", false, "Check for updates only")
	return c
}

func newCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completions",
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
	return cmd
}
