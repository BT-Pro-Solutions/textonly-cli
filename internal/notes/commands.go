package notes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/textonlyio/textonly-cli/internal/api"
	"github.com/textonlyio/textonly-cli/internal/auth"
)

func SetVisibility(id string, visibility string) error {
	v := map[string]any{"visibility": visibility}
	client := api.New(auth.LoadToken)
	return client.Do("POST", "/notes/"+id+"/visibility", v, true, nil)
}

func NewListCommand() *cobra.Command {
	var (
		pub bool
		priv bool
		asJSON bool
	)
	c := &cobra.Command{
		Use:   "list",
		Short: "List notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pub && priv { return errors.New("cannot set both --public and --private") }
			client := api.New(auth.LoadToken)
			qs := ""
			if pub { qs += "?public=1" }
			if priv { if qs == "" { qs += "?private=1" } else { qs += "&private=1" } }
			path := "/notes" + qs
			if asJSON {
				var v []map[string]any
				if err := client.Do("GET", path, nil, true, &v); err != nil { return err }
				b, _ := json.MarshalIndent(v, "", "  ")
				fmt.Println(string(b))
				return nil
			}
			var v []struct{ ID int `json:"id"`; Title string `json:"title"` }
			if err := client.Do("GET", path, nil, true, &v); err != nil { return err }
			for _, n := range v { fmt.Printf("%d\t%s\n", n.ID, n.Title) }
			return nil
		},
	}
	c.Flags().BoolVar(&pub, "public", false, "Show only public notes")
	c.Flags().BoolVar(&priv, "private", false, "Show only private notes")
	c.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return c
}

func NewViewCommand() *cobra.Command {
	var raw, openFlag, asJSON bool
	c := &cobra.Command{
		Use:   "view <id>",
		Args:  cobra.ExactArgs(1),
		Short: "View a note",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.New(auth.LoadToken)
			var v map[string]any
			if err := client.Do("GET", "/notes/"+args[0], nil, true, &v); err != nil { return err }
			if asJSON { b,_ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)); return nil }
			content, _ := v["content"].(string)
			title, _ := v["title"].(string)
			if raw {
				fmt.Println(content)
				return nil
			}
			// Default: title + content
			fmt.Println(title)
			if content != "" {
				fmt.Println()
				fmt.Println(content)
			}
			if openFlag {
				if guid, ok := v["guid"].(string); ok && guid != "" {
					url := "https://textonly.io/n/" + guid
					_ = openURL(url)
				}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&raw, "raw", false, "Print raw content only")
	c.Flags().BoolVar(&openFlag, "open", false, "Open in browser if public")
	c.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return c
}

func NewCreateCommand() *cobra.Command {
	var title, file string
	var stdin, asJSON bool
	var pub, priv bool
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a note",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pub && priv { return errors.New("cannot set both --public and --private") }
			content, err := readContent(file, stdin)
			if err != nil { return err }
			payload := map[string]any{"title": title, "content": content}
			if pub { payload["public"] = true } else if priv { payload["public"] = false }
			client := api.New(auth.LoadToken)
			var out map[string]any
			if err := client.Do("POST", "/notes", payload, true, &out); err != nil { return err }
			if asJSON { b,_ := json.MarshalIndent(out, "", "  "); fmt.Println(string(b)); return nil }
			if id, ok := out["id"].(float64); ok { fmt.Println(strconv.Itoa(int(id))) } else { fmt.Println("created") }
			return nil
		},
	}
	c.Flags().StringVar(&title, "title", "", "Title")
	c.Flags().StringVar(&file, "file", "", "File with content")
	c.Flags().BoolVar(&stdin, "stdin", false, "Read content from stdin")
	c.Flags().BoolVar(&pub, "public", false, "Set visibility to public")
	c.Flags().BoolVar(&priv, "private", false, "Set visibility to private")
	c.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return c
}

func NewUpdateCommand() *cobra.Command {
	var title, file string
	var stdin bool
	var pub, priv bool
	c := &cobra.Command{
		Use:   "update <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Update a note",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pub && priv { return errors.New("cannot set both --public and --private") }
			payload := map[string]any{}
			if title != "" { payload["title"] = title }
			if file != "" || stdin { c, err := readContent(file, stdin); if err != nil { return err }; payload["content"] = c }
			if pub { payload["public"] = true } else if priv { payload["public"] = false }
			client := api.New(auth.LoadToken)
			return client.Do("PATCH", "/notes/"+args[0], payload, true, nil)
		},
	}
	c.Flags().StringVar(&title, "title", "", "New title")
	c.Flags().StringVar(&file, "file", "", "File with content")
	c.Flags().BoolVar(&stdin, "stdin", false, "Read content from stdin")
	c.Flags().BoolVar(&pub, "public", false, "Set visibility to public")
	c.Flags().BoolVar(&priv, "private", false, "Set visibility to private")
	return c
}

func NewDeleteCommand() *cobra.Command {
	var yes bool
	c := &cobra.Command{
		Use:   "delete <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a note",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes { return errors.New("use --yes to confirm") }
			client := api.New(auth.LoadToken)
			return client.Do("DELETE", "/notes/"+args[0], nil, true, nil)
		},
	}
	c.Flags().BoolVar(&yes, "yes", false, "Confirm deletion")
	return c
}

func NewVisibilityCommand() *cobra.Command {
	var pub, priv bool
	c := &cobra.Command{
		Use:   "visibility <id> --public|--private",
		Args:  cobra.ExactArgs(1),
		Short: "Change note visibility",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pub == priv { return errors.New("must set exactly one of --public or --private") }
			v := map[string]any{}
			if pub { v["visibility"] = "public" } else { v["visibility"] = "private" }
			client := api.New(auth.LoadToken)
			return client.Do("POST", "/notes/"+args[0]+"/visibility", v, true, nil)
		},
	}
	c.Flags().BoolVar(&pub, "public", false, "Set visibility to public")
	c.Flags().BoolVar(&priv, "private", false, "Set visibility to private")
	return c
}

func NewStatsCommand() *cobra.Command {
	var asJSON bool
	c := &cobra.Command{
		Use:   "stats <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show note stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.New(auth.LoadToken)
			var v map[string]any
			if err := client.Do("GET", "/notes/"+args[0]+"/stats", nil, true, &v); err != nil { return err }
			if asJSON { b,_ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)); return nil }
			for k, vv := range v { fmt.Printf("%s: %v\n", k, vv) }
			return nil
		},
	}
	c.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return c
}

func NewLinkCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "link <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Print share/public URL if enabled",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.New(auth.LoadToken)
			var v map[string]any
			if err := client.Do("GET", "/notes/"+args[0], nil, true, &v); err != nil { return err }
			if guid, ok := v["guid"].(string); ok && guid != "" { fmt.Println("https://textonly.io/n/"+guid); return nil }
			return errors.New("note is not public or has no link")
		},
	}
	return c
}

func readContent(file string, stdin bool) (string, error) {
	if stdin {
		b, err := io.ReadAll(os.Stdin)
		if err != nil { return "", err }
		return string(b), nil
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil { return "", err }
		return string(b), nil
	}
	return "", errors.New("provide --file or --stdin")
}

var openURL = func(u string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", u).Start()
	case "linux":
		return exec.Command("xdg-open", u).Start()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
