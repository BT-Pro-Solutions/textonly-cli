package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/textonlyio/textonly-cli/internal/config"
)

type Client struct {
	baseURL       string
	httpClient    *http.Client
	tokenProvider func() (string, error)
}

var userAgent = "to/dev (unknown/unknown)"

func SetUserAgent(ua string) { if ua != "" { userAgent = ua } }

func New(tokenProvider func() (string, error)) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment
	return &Client{
		baseURL:       config.APIBaseURL(),
		httpClient:    &http.Client{Timeout: 20 * time.Second, Transport: transport},
		tokenProvider: tokenProvider,
	}
}

func (c *Client) Do(method, path string, body any, requireAuth bool, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("User-Agent", userAgent)
	if requireAuth {
		if tokenEnv := os.Getenv("TO_TOKEN"); tokenEnv != "" {
			req.Header.Set("Authorization", "Bearer "+tokenEnv)
		} else if c.tokenProvider != nil {
			if tok, err := c.tokenProvider(); err == nil && tok != "" {
				req.Header.Set("Authorization", "Bearer "+tok)
			}
		}
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error: %d %s", resp.StatusCode, string(b))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func sleepWithJitter(base time.Duration, attempt int) {
	delay := base << attempt
	jitter := time.Duration(rand.Int63n(int64(delay/2))) - delay/4
	time.Sleep(delay + jitter)
}
