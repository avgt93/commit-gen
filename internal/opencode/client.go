// Package opencode provides HTTP client for OpenCode API.
package opencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

type Session struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type HealthResponse struct {
	Healthy bool   `json:"healthy"`
	Version string `json:"version"`
}

type MessagePart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Model struct {
	ProviderID string `json:"providerID"`
	ModelID    string `json:"modelID"`
}

type PromptRequest struct {
	Model   *Model        `json:"model,omitempty"`
	Parts   []MessagePart `json:"parts"`
	NoReply bool          `json:"noReply,omitempty"`
}

type Message struct {
	Info struct {
		ID string `json:"id"`
	} `json:"info"`
	Parts []MessagePart `json:"parts"`
}

func NewClient(host string, port int, timeout int) *Client {
	baseURL := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (c *Client) CheckHealth() (bool, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/global/health", c.baseURL))
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return false, err
	}

	return health.Healthy, nil
}

func (c *Client) CreateSession(title string) (*Session, error) {
	reqBody := map[string]string{"title": title}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/session", c.baseURL),
		"application/json",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create session: %s (status %d)", string(body), resp.StatusCode)
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to parse session response: %w", err)
	}

	return &session, nil
}

func (c *Client) SendMessage(sessionID string, message string, model *Model) (string, error) {
	req := PromptRequest{
		Model: model,
		Parts: []MessagePart{
			{
				Type: "text",
				Text: message,
			},
		},
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/session/%s/message", c.baseURL, sessionID),
		"application/json",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to send message: %s (status %d)", string(body), resp.StatusCode)
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return "", fmt.Errorf("failed to parse message response: %w", err)
	}

	for _, part := range msg.Parts {
		if part.Type == "text" {
			return part.Text, nil
		}
	}

	return "", fmt.Errorf("no text response received")
}

func (c *Client) GetSession(sessionID string) (*Session, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/session/%s", c.baseURL, sessionID))
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("session not found (status %d)", resp.StatusCode)
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to parse session response: %w", err)
	}

	return &session, nil
}
