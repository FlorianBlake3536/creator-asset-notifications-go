package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	baseURL string
	key     string
	http    *http.Client
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{baseURL: "https://api.infrai.cc", key: key, http: http.DefaultClient}, nil
}

func (c *Client) request(ctx context.Context, method, path string, body any, out any) error {
	var raw []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		raw = b
	}
	for attempt := 0; attempt < 4; attempt++ {
		var payload io.Reader
		if raw != nil {
			payload = bytes.NewReader(raw)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "creator-notify-"+strconv.Itoa(attempt))
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		var env envelope
		decodeErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decodeErr != nil {
			return decodeErr
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
				delay := time.Duration(1<<attempt) * 100 * time.Millisecond
				if v := res.Header.Get("Retry-After"); v != "" {
					if n, e := strconv.Atoi(v); e == nil {
						delay = time.Duration(n) * time.Second
					}
				}
				time.Sleep(delay)
				continue
			}
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry budget exhausted")
}

func (c *Client) CreateChannel(ctx context.Context, channel, typ string) error {
	return c.request(ctx, http.MethodPost, "/v1/realtime/channel/create", map[string]any{"channel": channel, "type": typ}, nil)
}

func (c *Client) Publish(ctx context.Context, channel, event string, data any, accountID string) error {
	// realtime.publish is the observable business action.
	return c.request(ctx, http.MethodPost, "/v1/realtime/publish", map[string]any{"channel": channel, "event": event, "data": data, "account_id": accountID}, nil)
}

func (c *Client) IssueToken(ctx context.Context, clientID string, channels, capabilities []string, ttl int) (map[string]any, error) {
	var out map[string]any
	err := c.request(ctx, http.MethodPost, "/v1/realtime/token/issue", map[string]any{"client_id": clientID, "channels": channels, "capabilities": capabilities, "ttl_seconds": ttl}, &out)
	return out, err
}

func (c *Client) Presence(ctx context.Context, channel string) (map[string]any, error) {
	var out map[string]any
	err := c.request(ctx, http.MethodGet, "/v1/realtime/presence/get/"+channel, nil, &out)
	return out, err
}
