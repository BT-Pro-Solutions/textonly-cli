package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/keyring"

	"github.com/textonlyio/textonly-cli/internal/api"
	"github.com/textonlyio/textonly-cli/internal/config"
)

const keyringService = "textonly-cli"

func Login(noOpen bool) error {
	c := api.New(LoadToken)
	var startResp struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		Interval        int    `json:"interval"`
		ExpiresIn       int    `json:"expires_in"`
	}
	if err := c.Do("POST", "/cli/login/start", nil, false, &startResp); err != nil { return err }
	if startResp.Interval <= 0 { startResp.Interval = 3 }
	fmt.Printf("Go to %s and enter code: %s\n", startResp.VerificationURI, startResp.UserCode)
	if !noOpen { _ = openBrowser(startResp.VerificationURI) }

	deadline := time.Now().Add(time.Duration(startResp.ExpiresIn) * time.Second)
	for time.Now().Before(deadline) {
		var pollResp struct {
			Status      string `json:"status"`
			AccessToken string `json:"access_token"`
		}
		err := c.Do("POST", "/cli/login/poll", map[string]string{"device_code": startResp.DeviceCode}, false, &pollResp)
		if err == nil {
			if pollResp.AccessToken != "" {
				if err := SaveToken(pollResp.AccessToken); err != nil { return err }
				fmt.Println("login successful")
				return nil
			}
		}
		sleep := time.Duration(startResp.Interval) * time.Second
		j := sleep / 4
		sleep = sleep + (time.Duration(int64(time.Now().UnixNano())%int64(2*j)) - j)
		time.Sleep(sleep)
	}
	return errors.New("login timed out")
}

func Logout() error {
	tok, _ := LoadToken()
	if tok != "" {
		c := api.New(LoadToken)
		_ = c.Do("POST", "/auth/logout", nil, true, nil)
	}
	return ClearToken()
}

func WhoAmI(asJSON bool) error {
	c := api.New(LoadToken)
	var me map[string]any
	if err := c.Do("GET", "/me", nil, true, &me); err != nil { return err }
	if asJSON { b,_ := json.MarshalIndent(me, "", "  "); fmt.Println(string(b)); return nil }
	if v, ok := me["email"].(string); ok { fmt.Println(v); return nil }
	fmt.Println("authenticated")
	return nil
}

func openBrowser(url string) error {
	cmds := [][]string{{"open", url}, {"xdg-open", url}}
	for _, c := range cmds {
		cmd := exec.Command(c[0], c[1])
		if err := cmd.Start(); err == nil { return nil }
	}
	return fmt.Errorf("please open %s manually", url)
}

func tokenFilePath() string { return filepath.Join(config.Dir(), "token") }

func SaveToken(token string) error {
	r, err := keyring.Open(keyring.Config{ServiceName: keyringService})
	if err == nil {
		if e := r.Set(keyring.Item{Key: "token", Data: []byte(token)}); e == nil {
			return nil
		}
	}
	if err := os.MkdirAll(config.Dir(), 0o755); err != nil { return err }
	return os.WriteFile(tokenFilePath(), []byte(token), 0o600)
}

func LoadToken() (string, error) {
	r, err := keyring.Open(keyring.Config{ServiceName: keyringService})
	if err == nil {
		it, e := r.Get("token")
		if e == nil { return string(it.Data), nil }
	}
	b, err := os.ReadFile(tokenFilePath())
	if err != nil { return "", errors.New("not logged in") }
	return strings.TrimSpace(string(b)), nil
}

func ClearToken() error {
	r, err := keyring.Open(keyring.Config{ServiceName: keyringService})
	if err == nil { _ = r.Remove("token") }
	_ = os.Remove(tokenFilePath())
	return nil
}
