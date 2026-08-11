// Package mcpserver implements the Model Context Protocol server for Apple Music.
// It communicates exclusively via stdio (stdout for MCP, stderr for logs).
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/backend"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/config"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/version"
)

// Server is the MCP server instance.
type Server struct {
	cfg     *config.Config
	logger  *logging.Logger
	backend backend.Backend
	mu      sync.Mutex
	input   io.Reader
	output  io.Writer
	running bool
}

// NewServer creates a new MCP server.
func NewServer(cfg *config.Config, logger *logging.Logger, be backend.Backend) *Server {
	return &Server{
		cfg:     cfg,
		logger:  logger,
		backend: be,
		input:   os.Stdin,
		output:  os.Stdout,
	}
}

// NewServerWithIO creates a server with custom I/O (for testing).
func NewServerWithIO(cfg *config.Config, logger *logging.Logger, be backend.Backend, input io.Reader, output io.Writer) *Server {
	return &Server{
		cfg:     cfg,
		logger:  logger,
		backend: be,
		input:   input,
		output:  output,
	}
}

// Serve starts the MCP server loop on stdio.
func (s *Server) Serve(ctx context.Context) error {
	s.running = true
	defer func() { s.running = false }()

	s.logger.Info("apple-music-mcp v%s starting, transport=stdio", version.Short())

	decoder := json.NewDecoder(s.input)
	encoder := json.NewEncoder(s.output)
	encoder.SetEscapeHTML(false)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("server shutting down: %v", ctx.Err())
			return nil
		default:
		}

		var request map[string]interface{}
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				s.logger.Info("stdin closed, shutting down")
				return nil
			}
			s.logger.Error("failed to decode request: %v", err)
			continue
		}

		response := s.handleRequest(request)
		if response == nil {
			continue // notification, no response needed
		}

		if err := encoder.Encode(response); err != nil {
			s.logger.Error("failed to encode response: %v", err)
			return fmt.Errorf("write error: %w", err)
		}
	}
}

func (s *Server) handleRequest(req map[string]interface{}) map[string]interface{} {
	method, _ := req["method"].(string)
	id := req["id"]

	s.logger.Debug("received method=%s id=%v", method, id)

	switch method {
	case "initialize":
		return s.handleInitialize(id, req)
	case "initialized":
		// Notification — no response
		return nil
	case "tools/list":
		return s.handleToolsList(id)
	case "tools/call":
		return s.handleToolsCall(id, req)
	case "resources/list":
		return s.handleResourcesList(id)
	case "resources/read":
		return s.handleResourcesRead(id, req)
	case "ping":
		return s.makeResponse(id, map[string]interface{}{})
	default:
		s.logger.Warn("unknown method: %s", method)
		return s.makeError(id, -32601, fmt.Sprintf("method not found: %s", method))
	}
}

func (s *Server) handleInitialize(id interface{}, req map[string]interface{}) map[string]interface{} {
	return s.makeResponse(id, map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]interface{}{
			"name":    "apple-music-mcp",
			"version": version.Short(),
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
			"resources": map[string]interface{}{
				"subscribe":   false,
				"listChanged": false,
			},
		},
	})
}

func (s *Server) handleToolsList(id interface{}) map[string]interface{} {
	tools := []map[string]interface{}{
		{
			"name":        "music_get_state",
			"description": "Get the current playback state of Apple Music. Returns track info, player state, volume, shuffle, repeat, and position.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			"annotations": map[string]interface{}{
				"readOnlyHint": true,
			},
		},
		{
			"name":        "music_playback",
			"description": "Control Music.app playback. Actions: open, play, pause, toggle, stop, next, previous, restart_current, seek_absolute, seek_relative, play_track, play_album, play_artist, play_playlist, play_url, reveal.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Playback action to perform",
						"enum": []string{
							"open", "play", "pause", "toggle", "stop",
							"next", "previous", "restart_current",
							"seek_absolute", "seek_relative",
							"play_track", "play_album", "play_artist",
							"play_playlist", "play_url", "reveal",
						},
					},
					"target_id": map[string]interface{}{
						"type":        "string",
						"description": "Persistent ID of the track, album, artist, or playlist to target",
					},
					"target_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of the target: track, album, artist, playlist, url",
						"enum":        []string{"track", "album", "artist", "playlist", "url"},
					},
					"seek_position": map[string]interface{}{
						"type":        "number",
						"description": "Position in seconds for seek operations",
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"destructiveHint": true,
			},
		},
		{
			"name":        "music_preferences",
			"description": "Get or set Music.app preferences: volume, shuffle, repeat, and AirPlay output.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Preference action",
						"enum":        []string{"get", "set_volume", "set_shuffle", "set_repeat", "get_outputs", "set_output"},
					},
					"volume": map[string]interface{}{
						"type":        "integer",
						"description": "Volume level (0-100)",
						"minimum":     0,
						"maximum":     100,
					},
					"shuffle_mode": map[string]interface{}{
						"type":        "string",
						"description": "Shuffle mode: off, songs, albums",
						"enum":        []string{"off", "songs", "albums"},
					},
					"repeat_mode": map[string]interface{}{
						"type":        "string",
						"description": "Repeat mode: off, one, all",
						"enum":        []string{"off", "one", "all"},
					},
					"output_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the AirPlay device to set as output",
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"idempotentHint": true,
			},
		},
		{
			"name":        "music_search",
			"description": "Search your Music library for tracks, albums, artists, or playlists.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search text",
					},
					"scope": map[string]interface{}{
						"type":        "string",
						"description": "Search scope: library, all",
						"enum":        []string{"library", "all"},
						"default":     "all",
					},
					"types": map[string]interface{}{
						"type":        "array",
						"description": "Types to search for",
						"items": map[string]interface{}{
							"type": "string",
							"enum": []string{"track", "album", "artist", "playlist"},
						},
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum results (default 25, max 200)",
						"default":     25,
						"maximum":     200,
					},
				},
				"required": []string{"query"},
			},
			"annotations": map[string]interface{}{
				"readOnlyHint":  true,
				"openWorldHint": true,
			},
		},
		{
			"name":        "music_playlists",
			"description": "Manage Music playlists: list, get, create, rename, delete, add/remove tracks, copy, folders.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Playlist action",
						"enum": []string{
							"list", "get", "create", "rename", "delete",
							"add_tracks", "remove_tracks", "reorder_tracks",
							"copy", "list_folders", "create_folder", "move_to_folder",
						},
					},
					"playlist_id": map[string]interface{}{
						"type":        "string",
						"description": "Persistent ID of the playlist",
					},
					"new_name": map[string]interface{}{
						"type":        "string",
						"description": "New name for rename or create",
					},
					"track_ids": map[string]interface{}{
						"type":        "array",
						"description": "Array of track persistent IDs",
						"items":       map[string]interface{}{"type": "string"},
					},
					"target_folder_id": map[string]interface{}{
						"type":        "string",
						"description": "Persistent ID of target folder for move",
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"destructiveHint": true,
			},
		},
		{
			"name":        "music_library",
			"description": "Browse and search your music library: tracks, albums, artists, genres, recently added, recently played.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Library action",
						"enum": []string{
							"search", "recently_added", "recently_played",
							"list_tracks", "list_albums", "list_artists", "list_genres",
						},
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query (for search action)",
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"readOnlyHint": true,
			},
		},
		{
			"name":        "music_favorites",
			"description": "Manage favorites, ratings, and dislikes for tracks and albums.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Favorite action",
						"enum": []string{
							"get", "list", "favorite", "unfavorite",
							"suggest_less", "set_rating", "clear_rating",
						},
					},
					"target_id": map[string]interface{}{
						"type":        "string",
						"description": "Persistent ID of the track or album",
					},
					"target_type": map[string]interface{}{
						"type":        "string",
						"description": "Type: track or album",
						"enum":        []string{"track", "album"},
						"default":     "track",
					},
					"rating": map[string]interface{}{
						"type":        "integer",
						"description": "Rating (0-100, where 0=unrated, 20=1 star, 40=2 stars, 60=3 stars, 80=4 stars, 100=5 stars)",
						"minimum":     0,
						"maximum":     100,
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"destructiveHint": true,
			},
		},
		{
			"name":        "music_queue",
			"description": "Manage the Up Next queue. NOTE: The native Apple Events backend does NOT support the queue. Queue operations require the accessibility or MusicKit backend (not yet implemented).",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"description": "Queue action",
						"enum": []string{
							"list", "add_next", "add_later", "remove",
							"move", "jump", "clear", "get_autoplay", "set_autoplay",
						},
					},
					"track_id": map[string]interface{}{
						"type":        "string",
						"description": "Persistent ID of the track",
					},
					"from_index": map[string]interface{}{
						"type":        "integer",
						"description": "Source index for move",
					},
					"to_index": map[string]interface{}{
						"type":        "integer",
						"description": "Destination index for move",
					},
				},
				"required": []string{"action"},
			},
			"annotations": map[string]interface{}{
				"destructiveHint": true,
			},
		},
		{
			"name":        "music_capabilities",
			"description": "Report the current capabilities of the Apple Music MCP server based on your system configuration.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			"annotations": map[string]interface{}{
				"readOnlyHint": true,
			},
		},
		{
			"name":        "music_doctor",
			"description": "Run non-destructive diagnostics on the Apple Music integration. Checks Music.app availability, permissions, and backend status.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			"annotations": map[string]interface{}{
				"readOnlyHint": true,
			},
		},
	}

	return s.makeResponse(id, map[string]interface{}{
		"tools": tools,
	})
}

func (s *Server) handleResourcesList(id interface{}) map[string]interface{} {
	resources := []map[string]interface{}{
		{
			"uri":         "applemusic://now-playing",
			"name":        "Now Playing",
			"description": "Current playback status and track information",
			"mimeType":    "application/json",
		},
		{
			"uri":         "applemusic://playlists",
			"name":        "Playlists",
			"description": "List of all Music playlists",
			"mimeType":    "application/json",
		},
		{
			"uri":         "applemusic://capabilities",
			"name":        "Capabilities",
			"description": "System capability matrix for Apple Music control",
			"mimeType":    "application/json",
		},
	}

	return s.makeResponse(id, map[string]interface{}{
		"resources": resources,
	})
}

func (s *Server) handleResourcesRead(id interface{}, req map[string]interface{}) map[string]interface{} {
	params, _ := req["params"].(map[string]interface{})
	uri, _ := params["uri"].(string)

	ctx := context.Background()
	var contents []map[string]interface{}

	switch uri {
	case "applemusic://now-playing":
		status, err := s.backend.GetPlayerState(ctx)
		if err != nil {
			return s.makeError(id, -32000, fmt.Sprintf("failed to get state: %v", err))
		}
		data, _ := json.Marshal(status)
		contents = []map[string]interface{}{
			{"uri": uri, "mimeType": "application/json", "text": string(data)},
		}
	case "applemusic://playlists":
		result, err := s.backend.Playlists(ctx, domain.PlaylistCommand{Action: domain.PlaylistActionList})
		if err != nil {
			return s.makeError(id, -32000, fmt.Sprintf("failed to list playlists: %v", err))
		}
		data, _ := json.Marshal(result)
		contents = []map[string]interface{}{
			{"uri": uri, "mimeType": "application/json", "text": string(data)},
		}
	case "applemusic://capabilities":
		caps, err := s.backend.Capabilities(ctx)
		if err != nil {
			return s.makeError(id, -32000, fmt.Sprintf("failed to get capabilities: %v", err))
		}
		data, _ := json.Marshal(caps)
		contents = []map[string]interface{}{
			{"uri": uri, "mimeType": "application/json", "text": string(data)},
		}
	default:
		return s.makeError(id, -32001, fmt.Sprintf("resource not found: %s", uri))
	}

	return s.makeResponse(id, map[string]interface{}{
		"contents": contents,
	})
}

func (s *Server) handleToolsCall(id interface{}, req map[string]interface{}) map[string]interface{} {
	params, _ := req["params"].(map[string]interface{})
	toolName, _ := params["name"].(string)
	arguments, _ := params["arguments"].(map[string]interface{})

	s.logger.Debug("tool call: %s", toolName)

	ctx := context.Background()

	switch toolName {
	case "music_get_state":
		return s.callGetState(ctx, id)
	case "music_playback":
		return s.callPlayback(ctx, id, arguments)
	case "music_preferences":
		return s.callPreferences(ctx, id, arguments)
	case "music_search":
		return s.callSearch(ctx, id, arguments)
	case "music_playlists":
		return s.callPlaylists(ctx, id, arguments)
	case "music_library":
		return s.callLibrary(ctx, id, arguments)
	case "music_favorites":
		return s.callFavorites(ctx, id, arguments)
	case "music_queue":
		return s.callQueue(ctx, id, arguments)
	case "music_capabilities":
		return s.callCapabilities(ctx, id)
	case "music_doctor":
		return s.callDoctor(ctx, id)
	default:
		return s.makeError(id, -32601, fmt.Sprintf("unknown tool: %s", toolName))
	}
}

func (s *Server) callGetState(ctx context.Context, id interface{}) map[string]interface{} {
	status, err := s.backend.GetPlayerState(ctx)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, status, s.backend.Name(), "native")
}

func (s *Server) callPlayback(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	action, _ := args["action"].(string)
	if action == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "action is required", s.backend.Name())
	}

	cmd := domain.PlaybackCommand{
		Action: domain.PlaybackAction(action),
	}
	if v, ok := args["target_id"].(string); ok {
		cmd.TargetID = v
	}
	if v, ok := args["target_type"].(string); ok {
		cmd.TargetType = v
	}
	if v, ok := args["seek_position"].(float64); ok {
		cmd.SeekPosition = v
	}

	if s.cfg.ReadOnly && cmd.Action != domain.ActionReveal {
		return s.formatToolError(id, domain.ErrDestructiveActionBlocked,
			"read-only mode is enabled. Set APPLE_MUSIC_MCP_READ_ONLY=0 to allow mutations.",
			s.backend.Name())
	}

	result, err := s.backend.Playback(ctx, cmd)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callPreferences(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	action, _ := args["action"].(string)
	if action == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "action is required", s.backend.Name())
	}

	cmd := domain.PreferenceCommand{
		Action: domain.PreferenceAction(action),
	}
	if v, ok := args["volume"].(float64); ok {
		vol := int(v)
		cmd.Volume = &vol
	}
	if v, ok := args["shuffle_mode"].(string); ok {
		mode := domain.ShuffleMode(v)
		cmd.ShuffleMode = &mode
	}
	if v, ok := args["repeat_mode"].(string); ok {
		mode := domain.RepeatMode(v)
		cmd.RepeatMode = &mode
	}
	if v, ok := args["output_name"].(string); ok {
		cmd.OutputName = &v
	}

	if s.cfg.ReadOnly && cmd.Action != domain.PrefActionGet && cmd.Action != domain.PrefActionGetOutputs {
		return s.formatToolError(id, domain.ErrDestructiveActionBlocked,
			"read-only mode is enabled", s.backend.Name())
	}

	result, err := s.backend.Preferences(ctx, cmd)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callSearch(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	query, _ := args["query"].(string)
	if query == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "query is required", s.backend.Name())
	}

	sq := domain.SearchQuery{
		Query: query,
		Limit: 25,
	}
	if v, ok := args["scope"].(string); ok {
		sq.Scope = domain.SearchScope(v)
	}
	if v, ok := args["limit"].(float64); ok {
		sq.Limit = int(v)
	}
	if types, ok := args["types"].([]interface{}); ok {
		for _, t := range types {
			if ts, ok := t.(string); ok {
				sq.Types = append(sq.Types, domain.SearchType(ts))
			}
		}
	}

	result, err := s.backend.Search(ctx, sq)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callPlaylists(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	action, _ := args["action"].(string)
	if action == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "action is required", s.backend.Name())
	}

	cmd := domain.PlaylistCommand{
		Action: domain.PlaylistAction(action),
	}
	if v, ok := args["playlist_id"].(string); ok {
		cmd.PlaylistID = v
	}
	if v, ok := args["new_name"].(string); ok {
		cmd.NewName = v
	}
	if v, ok := args["target_folder_id"].(string); ok {
		cmd.TargetFolderID = v
	}
	if ids, ok := args["track_ids"].([]interface{}); ok {
		for _, tid := range ids {
			if ts, ok := tid.(string); ok {
				cmd.TrackIDs = append(cmd.TrackIDs, ts)
			}
		}
	}

	destructiveActions := map[domain.PlaylistAction]bool{
		domain.PlaylistActionDelete:       true,
		domain.PlaylistActionRemoveTracks: true,
	}
	if s.cfg.ReadOnly && destructiveActions[cmd.Action] {
		return s.formatToolError(id, domain.ErrDestructiveActionBlocked,
			"read-only mode is enabled", s.backend.Name())
	}

	result, err := s.backend.Playlists(ctx, cmd)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callLibrary(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	action, _ := args["action"].(string)
	if action == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "action is required", s.backend.Name())
	}

	cmd := domain.LibraryCommand{
		Action: domain.LibraryAction(action),
	}
	if v, ok := args["query"].(string); ok {
		cmd.Query = v
	}

	result, err := s.backend.Library(ctx, cmd)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callFavorites(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	action, _ := args["action"].(string)
	if action == "" {
		return s.formatToolError(id, domain.ErrInvalidArgument, "action is required", s.backend.Name())
	}

	cmd := domain.FavoriteCommand{
		Action: domain.FavoriteAction(action),
	}
	if v, ok := args["target_id"].(string); ok {
		cmd.TargetID = v
	}
	if v, ok := args["target_type"].(string); ok {
		cmd.TargetType = v
	}
	if v, ok := args["rating"].(float64); ok {
		r := int(v)
		cmd.Rating = &r
	}

	if s.cfg.ReadOnly && cmd.Action != domain.FavoriteActionGet && cmd.Action != domain.FavoriteActionList {
		return s.formatToolError(id, domain.ErrDestructiveActionBlocked,
			"read-only mode is enabled", s.backend.Name())
	}

	result, err := s.backend.Favorites(ctx, cmd)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, result, s.backend.Name(), "native")
}

func (s *Server) callQueue(ctx context.Context, id interface{}, args map[string]interface{}) map[string]interface{} {
	// Queue is not supported by the native backend
	return s.formatToolError(id, domain.ErrUnsupportedCapability,
		"Up Next queue is not available via Apple Events. The queue requires an accessibility or MusicKit backend which is not yet implemented in v0.1.0. This is a known limitation.",
		s.backend.Name())
}

func (s *Server) callCapabilities(ctx context.Context, id interface{}) map[string]interface{} {
	caps, err := s.backend.Capabilities(ctx)
	if err != nil {
		return s.formatToolError(id, domain.ErrAutomationError, err.Error(), s.backend.Name())
	}
	return s.formatToolResult(id, caps, s.backend.Name(), "native")
}

func (s *Server) callDoctor(ctx context.Context, id interface{}) map[string]interface{} {
	diag := map[string]interface{}{
		"music_app_installed":   false,
		"music_app_path":        "",
		"music_app_version":     "",
		"automation_working":    false,
		"accessibility_granted": false,
		"backend_name":          s.backend.Name(),
		"mcp_protocol":          "2024-11-05",
		"transport":             s.cfg.Transport,
		"read_only":             s.cfg.ReadOnly,
		"version":               version.Short(),
	}

	// Check Music.app exists
	paths := []string{"/System/Applications/Music.app", "/Applications/Music.app"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			diag["music_app_installed"] = true
			diag["music_app_path"] = p
			break
		}
	}

	// Check automation
	if s.backend.IsAvailable(ctx) {
		diag["automation_working"] = true

		// Get version
		out, err := s.backend.GetPlayerState(ctx)
		if err == nil {
			diag["music_app_version"] = "responding (state available)"
		}
		_ = out
	}

	return s.formatToolResult(id, diag, s.backend.Name(), "native")
}

// --- Response helpers ---

func (s *Server) formatToolResult(id interface{}, data interface{}, backend string, mode string) map[string]interface{} {
	resp := domain.NewAPIResponse(data, backend, mode, fmt.Sprintf("%v", id))
	return s.makeResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": mustMarshal(resp),
			},
		},
	})
}

func (s *Server) formatToolError(id interface{}, code domain.ErrorCode, message string, backend string) map[string]interface{} {
	errResp := domain.NewAPIError(code, message, code != domain.ErrAutomationError, backend, fmt.Sprintf("%v", id))
	return s.makeResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": mustMarshal(errResp),
			},
		},
		"isError": true,
	})
}

func (s *Server) makeResponse(id interface{}, result map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}

func (s *Server) makeError(id interface{}, code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

func mustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"ok":false,"error":{"code":"internal_error","message":"%s"}}`, err.Error())
	}
	return string(data)
}
