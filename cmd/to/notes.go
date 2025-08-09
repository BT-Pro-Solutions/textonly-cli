package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/textonlyio/textonly-cli/internal/notes"
)

func newNotesCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "notes", Short: "Manage notes"}

	// Pattern: to notes <id> make public|private
	cmd.Args = cobra.ArbitraryArgs
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) >= 3 && args[1] == "make" && (args[2] == "public" || args[2] == "private") {
			id := args[0]
			visibility := args[2]
			if err := notes.SetVisibility(id, visibility); err != nil {
				return err
			}
			fmt.Printf("note %s is now %s\n", id, visibility)
			return nil
		}
		return errors.New("usage: to notes <id> make public|private")
	}

	cmd.AddCommand(notes.NewListCommand())
	cmd.AddCommand(notes.NewViewCommand())
	cmd.AddCommand(notes.NewCreateCommand())
	cmd.AddCommand(notes.NewUpdateCommand())
	cmd.AddCommand(notes.NewDeleteCommand())
	cmd.AddCommand(notes.NewVisibilityCommand())
	cmd.AddCommand(notes.NewStatsCommand())
	cmd.AddCommand(notes.NewLinkCommand())
	return cmd
}
