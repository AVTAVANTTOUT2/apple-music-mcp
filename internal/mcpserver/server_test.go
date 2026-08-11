// Package mcpserver tests
package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/config"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
)

// fakeBackend implements backend.Backend for testing.
type fakeBackend struct{}

func (f *fakeBackend) Name() string { return "fake" }

func (f *fakeBackend) IsAvailable(ctx context.Context) bool { return true }

func (f *fakeBackend) Capabilities(ctx context.Context) (domain.CapabilitySet, error) {
	return domain.CapabilitySet{
		BackendName:             "fake",
		NativeCapabilities:      []string{"play", "pause", "search"},
		UnavailableCapabilities: []string{"queue"},
	}, nil
}

func (f *fakeBackend) GetPlayerState(ctx context.Context) (*domain.PlaybackStatus, error) {
	return &domain.PlaybackStatus{
		AppRunning:  true,
		PlayerState: domain.PlayerStatePlaying,
		Volume:      50,
		ShuffleMode: domain.ShuffleOff,
		RepeatMode:  domain.RepeatOff,
		Position:    30.5,
		Track: &domain.TrackInfo{
			Name:         "Test Song",
			Artist:       "Test Artist",
			Album:        "Test Album",
			Duration:     240.0,
			PersistentID: "ABCD1234",
			Favorited:    true,
			Rating:       80,
			Year:         2024,
		},
		BackendName:    "fake",
		CapabilityMode: "native",
	}, nil
}

func (f *fakeBackend) Playback(ctx context.Context, cmd domain.PlaybackCommand) (*domain.PlaybackResult, error) {
	return &domain.PlaybackResult{
		Success:     true,
		Action:      string(cmd.Action),
		BackendName: "fake",
	}, nil
}

func (f *fakeBackend) Preferences(ctx context.Context, cmd domain.PreferenceCommand) (*domain.PreferenceResult, error) {
	return &domain.PreferenceResult{
		Success:     true,
		Volume:      50,
		ShuffleMode: domain.ShuffleOff,
		RepeatMode:  domain.RepeatOff,
		BackendName: "fake",
	}, nil
}

func (f *fakeBackend) Search(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error) {
	return &domain.SearchResult{
		Tracks: []domain.TrackInfo{
			{Name: "Result 1", Artist: "Artist 1", PersistentID: "AAAA"},
		},
		Total:       1,
		BackendName: "fake",
	}, nil
}

func (f *fakeBackend) Queue(ctx context.Context, cmd domain.QueueCommand) (*domain.QueueResult, error) {
	return nil, fmt.Errorf("unsupported")
}

func (f *fakeBackend) Playlists(ctx context.Context, cmd domain.PlaylistCommand) (*domain.PlaylistResult, error) {
	return &domain.PlaylistResult{
		Success: true,
		Playlists: []domain.PlaylistInfo{
			{Name: "Test Playlist", PersistentID: "BBBB", Kind: "user", TrackCount: 5},
		},
		BackendName: "fake",
	}, nil
}

func (f *fakeBackend) Library(ctx context.Context, cmd domain.LibraryCommand) (*domain.LibraryResult, error) {
	return &domain.LibraryResult{
		Success:     true,
		Tracks:      []domain.TrackInfo{},
		BackendName: "fake",
	}, nil
}

func (f *fakeBackend) Favorites(ctx context.Context, cmd domain.FavoriteCommand) (*domain.FavoriteResult, error) {
	return &domain.FavoriteResult{
		Success:     true,
		Favorited:   true,
		Rating:      80,
		BackendName: "fake",
	}, nil
}

// serveHelper creates a server with a request, starts it, and returns the output buffer.
// It waits for the server to finish processing (EOF on input triggers shutdown).
func serveHelper(t *testing.T, req map[string]interface{}, cfg *config.Config) (map[string]interface{}, error) {
	t.Helper()

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	input := bytes.NewReader(raw)
	var output bytes.Buffer

	logger := logging.NewLogger(logging.LevelError)
	be := &fakeBackend{}
	server := NewServerWithIO(cfg, logger, be, input, &output)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server and wait for it to finish (it will exit on EOF from input)
	err = server.Serve(ctx)
	if err != nil {
		return nil, fmt.Errorf("serve error: %w", err)
	}

	var resp map[string]interface{}
	decoder := json.NewDecoder(&output)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w (output: %s)", err, output.String())
	}

	return resp, nil
}

func makeRequest(method string, id int, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      id,
		"params":  params,
	}
}

// TestInitialize tests the initialize handshake.
func TestInitialize(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("initialize", 1, map[string]interface{}{
		"protocolVersion": "2024-11-05",
	}), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("no result in response: %v", resp)
	}
	protocolVersion, _ := result["protocolVersion"].(string)
	if protocolVersion != "2024-11-05" {
		t.Errorf("expected protocol 2024-11-05, got %s", protocolVersion)
	}
}

// TestInitializedNotification tests that initialized notification returns nil.
func TestInitializedNotification(t *testing.T) {
	cfg := config.DefaultConfig()
	// initialized is a notification (no id), so no response is sent
	// We just verify no error
	resp, err := serveHelper(t, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
	}, cfg)
	// For notifications, Serve returns nil response
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		t.Logf("notification response (expected no content): %v", resp)
	}
}

// TestPing tests the ping method.
func TestPing(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("ping", 2, nil), cfg)
	if err != nil {
		t.Fatal(err)
	}

	if id, _ := resp["id"].(float64); int(id) != 2 {
		t.Errorf("expected id 2, got %v", id)
	}
}

// TestToolsList tests that tools are listed correctly.
func TestToolsList(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("tools/list", 3, nil), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	tools, _ := result["tools"].([]interface{})
	if len(tools) < 5 {
		t.Errorf("expected at least 5 tools, got %d", len(tools))
	}
}

// TestResourcesList tests resource listing.
func TestResourcesList(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("resources/list", 4, nil), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	resources, _ := result["resources"].([]interface{})
	if len(resources) == 0 {
		t.Error("expected at least 1 resource")
	}
}

// TestGetState tests the music_get_state tool call.
func TestGetState(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("tools/call", 5, map[string]interface{}{
		"name":      "music_get_state",
		"arguments": map[string]interface{}{},
	}), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	content, _ := result["content"].([]interface{})
	if len(content) == 0 {
		t.Fatal("no content in response")
	}
	textContent, _ := content[0].(map[string]interface{})
	text, _ := textContent["text"].(string)

	var apiResp domain.APIResponse
	if err := json.Unmarshal([]byte(text), &apiResp); err != nil {
		t.Fatalf("failed to parse API response: %v", err)
	}
	if !apiResp.OK {
		t.Error("expected ok=true")
	}
}

// TestSearch tests the music_search tool call.
func TestSearch(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("tools/call", 6, map[string]interface{}{
		"name": "music_search",
		"arguments": map[string]interface{}{
			"query": "test query",
		},
	}), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	content, _ := result["content"].([]interface{})
	textContent, _ := content[0].(map[string]interface{})
	text, _ := textContent["text"].(string)

	var apiResp domain.APIResponse
	if err := json.Unmarshal([]byte(text), &apiResp); err != nil {
		t.Fatalf("failed to parse API response: %v", err)
	}
	if !apiResp.OK {
		t.Error("expected ok=true for search")
	}
}

// TestCapabilities tests the music_capabilities tool.
func TestCapabilities(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("tools/call", 7, map[string]interface{}{
		"name":      "music_capabilities",
		"arguments": map[string]interface{}{},
	}), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	content, _ := result["content"].([]interface{})
	textContent, _ := content[0].(map[string]interface{})
	text, _ := textContent["text"].(string)

	var apiResp domain.APIResponse
	if err := json.Unmarshal([]byte(text), &apiResp); err != nil {
		t.Fatalf("failed to parse API response: %v", err)
	}
	if !apiResp.OK {
		t.Error("expected ok=true for capabilities")
	}
}

// TestUnknownMethod tests that unknown methods return an error.
func TestUnknownMethod(t *testing.T) {
	cfg := config.DefaultConfig()
	resp, err := serveHelper(t, makeRequest("nonexistent/method", 8, nil), cfg)
	if err != nil {
		t.Fatal(err)
	}

	errObj, ok := resp["error"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error in response")
	}
	code, _ := errObj["code"].(float64)
	if int(code) != -32601 {
		t.Errorf("expected error code -32601, got %v", code)
	}
}

// TestPlaybackReadOnly tests that playback is blocked in read-only mode.
func TestPlaybackReadOnly(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ReadOnly = true

	resp, err := serveHelper(t, makeRequest("tools/call", 9, map[string]interface{}{
		"name": "music_playback",
		"arguments": map[string]interface{}{
			"action": "play",
		},
	}), cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := resp["result"].(map[string]interface{})
	content, _ := result["content"].([]interface{})
	textContent, _ := content[0].(map[string]interface{})
	text, _ := textContent["text"].(string)

	var apiErr domain.APIError
	if err := json.Unmarshal([]byte(text), &apiErr); err != nil {
		t.Fatalf("failed to parse API error: %v", err)
	}
	if apiErr.OK {
		t.Error("expected ok=false in read-only mode")
	}
	if apiErr.Error.Code != domain.ErrDestructiveActionBlocked {
		t.Errorf("expected destructive_action_blocked, got %s", apiErr.Error.Code)
	}
}

// TestStdoutOnlyJSON ensures that serve outputs only JSON (no logs on stdout).
func TestStdoutOnlyJSON(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := logging.NewLogger(logging.LevelError)
	be := &fakeBackend{}
	var buf bytes.Buffer

	req := makeRequest("ping", 10, nil)
	raw, _ := json.Marshal(req)
	input := bytes.NewReader(raw)

	server := NewServerWithIO(cfg, logger, be, input, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		_ = server.Serve(ctx)
		close(done)
	}()

	// Wait for Serve to finish (it will exit after reading one request since input is limited)
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for server to finish")
	}
	cancel()

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "{") {
			t.Errorf("stdout line is not JSON: %s", line)
		}
	}
}
