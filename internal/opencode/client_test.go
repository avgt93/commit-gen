package opencode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestClientCreation tests creating a new OpenCode client
func TestClientCreation(t *testing.T) {
	client := NewClient("localhost", 4096, 30)

	if client == nil {
		t.Error("NewClient returned nil")
	}

	if client.baseURL != "http://localhost:4096" {
		t.Errorf("Base URL incorrect: got %q, expected %q", client.baseURL, "http://localhost:4096")
	}

	t.Log("✓ OpenCode client created successfully")
}

// TestClientBaseURL tests different host/port combinations
func TestClientBaseURL(t *testing.T) {
	tests := []struct {
		host     string
		port     int
		expected string
	}{
		{"localhost", 4096, "http://localhost:4096"},
		{"127.0.0.1", 4096, "http://127.0.0.1:4096"},
		{"example.com", 8080, "http://example.com:8080"},
	}

	for _, tt := range tests {
		client := NewClient(tt.host, tt.port, 30)
		if client.baseURL != tt.expected {
			t.Errorf("Base URL mismatch: got %q, expected %q", client.baseURL, tt.expected)
		} else {
			t.Logf("✓ Base URL correct for %s:%d", tt.host, tt.port)
		}
	}
}

// TestCheckHealthSuccess tests successful health check
func TestCheckHealthSuccess(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/global/health" {
			t.Errorf("Wrong path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Healthy: true,
			Version: "1.0.0",
		})
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	healthy, err := client.CheckHealth()
	if err != nil {
		t.Fatalf("CheckHealth failed: %v", err)
	}

	if !healthy {
		t.Error("Expected healthy=true, got false")
	}

	t.Log("✓ Health check passed")
}

// TestCheckHealthFailure tests failed health check
func TestCheckHealthFailure(t *testing.T) {
	// Create a mock server that returns unhealthy
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Healthy: false,
			Version: "1.0.0",
		})
	}))
	defer server.Close()

	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	healthy, err := client.CheckHealth()
	if err != nil {
		t.Fatalf("CheckHealth failed: %v", err)
	}

	if healthy {
		t.Error("Expected healthy=false, got true")
	}

	t.Log("✓ Unhealthy response detected correctly")
}

// TestCreateSessionSuccess tests successful session creation
func TestCreateSessionSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session" {
			t.Errorf("Wrong path: %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Wrong method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Session{
			ID:    "session-123",
			Title: "Test Session",
		})
	}))
	defer server.Close()

	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	session, err := client.CreateSession("Test Session")
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if session.ID != "session-123" {
		t.Errorf("Session ID mismatch: got %q, expected %q", session.ID, "session-123")
	}

	t.Logf("✓ Session created: %s", session.ID)
}

// TestSendMessageSuccess tests successful message sending
func TestSendMessageSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Message{
			Info: struct {
				ID string `json:"id"`
			}{ID: "msg-123"},
			Parts: []MessagePart{
				{
					Type: "text",
					Text: "Generated commit message",
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	model := &Model{
		ProviderID: "anthropic",
		ModelID:    "claude-3-5-sonnet-20241022",
	}

	response, err := client.SendMessage("session-123", "Test message", model)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if response != "Generated commit message" {
		t.Errorf("Response mismatch: got %q, expected %q", response, "Generated commit message")
	}

	t.Logf("✓ Message sent and response received: %s", response)
}

// TestSendMessageExtractsFirstTextPart tests that SendMessage extracts text correctly
func TestSendMessageExtractsFirstTextPart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Message{
			Info: struct {
				ID string `json:"id"`
			}{ID: "msg-456"},
			Parts: []MessagePart{
				{Type: "code", Text: "some code"},
				{Type: "text", Text: "feat: add feature"},
				{Type: "text", Text: "should not be used"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	response, err := client.SendMessage("session-123", "test", nil)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if response != "feat: add feature" {
		t.Errorf("Should extract first text part: got %q", response)
	}

	t.Log("✓ Correctly extracts first text part from response")
}

// TestGetSessionSuccess tests successful session retrieval
func TestGetSessionSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Session{
			ID:    "session-123",
			Title: "Retrieved Session",
		})
	}))
	defer server.Close()

	client := NewClient("localhost", 9999, 5)
	client.baseURL = server.URL

	session, err := client.GetSession("session-123")
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if session.ID != "session-123" {
		t.Errorf("Session ID mismatch: got %q, expected %q", session.ID, "session-123")
	}

	t.Logf("✓ Session retrieved: %s", session.ID)
}

// TestClientTimeout tests that client timeout is set
func TestClientTimeout(t *testing.T) {
	client := NewClient("localhost", 4096, 15)

	if client.timeout != 15*time.Second {
		t.Errorf("Timeout mismatch: got %v, expected 15s", client.timeout)
	}

	t.Logf("✓ Client timeout set correctly: %v", client.timeout)
}

// TestMessagePartTypes tests different message part types
func TestMessagePartTypes(t *testing.T) {
	types := []string{"text", "code", "image", "json"}

	for _, partType := range types {
		part := MessagePart{
			Type: partType,
			Text: "test content",
		}

		if part.Type == "" {
			t.Errorf("Message part type is empty for %s", partType)
		} else {
			t.Logf("✓ Message part type: %s", part.Type)
		}
	}
}

// TestModelConfiguration tests model struct configuration
func TestModelConfiguration(t *testing.T) {
	model := &Model{
		ProviderID: "anthropic",
		ModelID:    "claude-3-5-sonnet-20241022",
	}

	if model.ProviderID == "" {
		t.Error("Provider ID is empty")
	}

	if model.ModelID == "" {
		t.Error("Model ID is empty")
	}

	t.Logf("✓ Model configured: %s/%s", model.ProviderID, model.ModelID)
}
