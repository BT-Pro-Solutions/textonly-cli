package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/textonlyio/textonly-cli/internal/auth"
	"github.com/textonlyio/textonly-cli/internal/config"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run connectivity and auth checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			base := config.APIBaseURL()
			fmt.Println("API:", base)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(base + "/me")
			if err != nil {
				fmt.Println("network:", err)
			} else {
				_ = resp.Body.Close()
				if resp.StatusCode == 200 || resp.StatusCode == 401 {
					fmt.Println("network: ok")
				} else {
					fmt.Println("network: unexpected status", resp.StatusCode)
				}
			}
			if _, err := auth.LoadToken(); err != nil {
				fmt.Println("auth:", err)
			} else {
				fmt.Println("auth: token present")
			}
			return nil
		},
	}
}
