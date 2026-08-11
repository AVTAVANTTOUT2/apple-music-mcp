// Package musicapp implements the native Apple Events backend for Music.app.
// It uses static AppleScript files with parameterized arguments — never
// interpolates user input into script source.
package musicapp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/applescript"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
)

// Backend implements backend.Backend using Apple Events via osascript.
type Backend struct {
	executor   *applescript.Executor
	logger     *logging.Logger
	scriptsDir string
}

// NewBackend creates a new native Music.app backend.
// scriptsDir is the path to the scripts/ directory containing the .applescript files.
func NewBackend(logger *logging.Logger) *Backend {
	return &Backend{
		executor: applescript.NewExecutor(logger),
		logger:   logger,
	}
}

// SetScriptsDir sets the directory containing the AppleScript files.
func (b *Backend) SetScriptsDir(dir string) {
	b.scriptsDir = dir
}

func (b *Backend) scriptPath(name string) string {
	if b.scriptsDir == "" {
		// Default: look relative to the binary or current directory
		if exe, err := os.Executable(); err == nil {
			dir := filepath.Dir(exe)
			candidate := filepath.Join(dir, "scripts", name)
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		return filepath.Join("scripts", name)
	}
	return filepath.Join(b.scriptsDir, name)
}

// Name returns "musicapp".
func (b *Backend) Name() string {
	return "musicapp"
}

// IsAvailable returns true if Music.app is installed and accessible.
func (b *Backend) IsAvailable(ctx context.Context) bool {
	_, err := b.executor.RunAppleScriptString(ctx,
		`tell application "Music" to get name`)
	return err == nil
}

// Capabilities detects what this backend can do on the current system.
func (b *Backend) Capabilities(ctx context.Context) (domain.CapabilitySet, error) {
	cap := domain.CapabilitySet{
		BackendName: "musicapp",
		NativeCapabilities: []string{
			"play", "pause", "toggle", "stop",
			"next_track", "previous_track", "restart_current",
			"seek_absolute", "seek_relative",
			"get_state", "get_metadata",
			"volume", "shuffle", "repeat", "mute",
			"search_library", "search_playlist",
			"list_playlists", "create_playlist", "delete_playlist",
			"rename_playlist", "add_tracks_playlist", "remove_tracks_playlist",
			"duplicate_playlist", "folder_playlists",
			"favorite", "unfavorite", "suggest_less",
			"set_rating", "list_favorites",
			"list_airplay_devices", "set_airplay_output",
			"reveal_item", "open_location",
			"track_metadata", "artwork",
		},
		UnavailableCapabilities: []string{
			"queue_list", "queue_add_next", "queue_add_later",
			"queue_remove", "queue_move", "queue_jump", "queue_clear",
			"autoplay_get", "autoplay_set",
			"catalog_search",
		},
	}

	if cap.MusicAppInstalled {
		// Get version
		out, err := b.executor.RunAppleScriptString(ctx,
			`tell application "Music" to get version`)
		if err == nil {
			cap.MusicAppVersion = strings.TrimSpace(out)
		}
	}

	cap.MusicAppInstalled = b.IsAvailable(ctx)
	cap.AutomationGranted = b.IsAvailable(ctx)

	return cap, nil
}

// GetPlayerState returns the current playback state.
func (b *Backend) GetPlayerState(ctx context.Context) (*domain.PlaybackStatus, error) {
	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("get_state.applescript"))
	if err != nil {
		return nil, fmt.Errorf("failed to get player state: %w", err)
	}

	kv := applescript.ParseKV(out)
	status := applescript.ParsePlayerState(kv)
	status.BackendName = b.Name()

	return status, nil
}

// Playback executes a playback command.
func (b *Backend) Playback(ctx context.Context, cmd domain.PlaybackCommand) (*domain.PlaybackResult, error) {
	args := []string{string(cmd.Action)}
	if cmd.TargetID != "" {
		args = append(args, cmd.TargetID)
		if cmd.TargetType != "" {
			args = append(args, cmd.TargetType)
		}
	}
	// For seek operations, pass position
	if cmd.Action == domain.ActionSeekAbsolute || cmd.Action == domain.ActionSeekRelative {
		args = append(args, fmt.Sprintf("%.3f", cmd.SeekPosition))
	}
	if cmd.Once {
		args = append(args, "true")
	}

	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("playback.applescript"), args...)
	if err != nil {
		return &domain.PlaybackResult{
			Success:     false,
			Action:      string(cmd.Action),
			BackendName: b.Name(),
		}, fmt.Errorf("playback command failed: %w", err)
	}

	kv := applescript.ParseKV(out)
	result := applescript.ParsePlaybackResult(kv)
	result.BackendName = b.Name()

	return result, nil
}

// Preferences gets or sets player preferences.
func (b *Backend) Preferences(ctx context.Context, cmd domain.PreferenceCommand) (*domain.PreferenceResult, error) {
	args := []string{string(cmd.Action)}

	switch cmd.Action {
	case domain.PrefActionSetVolume:
		if cmd.Volume == nil {
			return nil, fmt.Errorf("volume is required for set_volume")
		}
		if err := domain.ValidateVolume(*cmd.Volume); err != nil {
			return nil, err
		}
		args = append(args, fmt.Sprintf("%d", *cmd.Volume))
	case domain.PrefActionSetShuffle:
		if cmd.ShuffleMode == nil {
			return nil, fmt.Errorf("shuffle_mode is required for set_shuffle")
		}
		args = append(args, string(*cmd.ShuffleMode))
	case domain.PrefActionSetRepeat:
		if cmd.RepeatMode == nil {
			return nil, fmt.Errorf("repeat_mode is required for set_repeat")
		}
		args = append(args, string(*cmd.RepeatMode))
	case domain.PrefActionSetOutput:
		if cmd.OutputName == nil {
			return nil, fmt.Errorf("output_name is required for set_output")
		}
		args = append(args, applescript.SanitizeArg(*cmd.OutputName))
	}

	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("preferences.applescript"), args...)
	if err != nil {
		return &domain.PreferenceResult{Success: false, BackendName: b.Name()}, fmt.Errorf("preference command failed: %w", err)
	}

	var result domain.PreferenceResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return &domain.PreferenceResult{Success: false, BackendName: b.Name()}, fmt.Errorf("failed to parse preference result: %w", err)
	}
	result.BackendName = b.Name()

	return &result, nil
}

// Search performs a search using the native backend (library only).
func (b *Backend) Search(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error) {
	if err := domain.ValidateSearchQuery(&query); err != nil {
		return nil, err
	}

	sanitized := applescript.SanitizeArg(query.Query)

	// Search each requested type
	result := &domain.SearchResult{BackendName: b.Name()}

	for _, st := range query.Types {
		args := []string{sanitized, "library", string(st)}
		out, err := b.executor.RunAppleScript(ctx, b.scriptPath("search.applescript"), args...)
		if err != nil {
			b.logger.Warn("search for type %s failed: %v", st, err)
			continue
		}

		switch st {
		case domain.SearchTypeTrack:
			var resp struct {
				Tracks []domain.TrackInfo `json:"tracks"`
				Total  int                `json:"total"`
			}
			if err := json.Unmarshal([]byte(out), &resp); err == nil {
				result.Tracks = resp.Tracks
				result.Total += resp.Total
			}
		case domain.SearchTypeAlbum:
			var resp struct {
				Albums []domain.AlbumInfo `json:"albums"`
				Total  int                `json:"total"`
			}
			if err := json.Unmarshal([]byte(out), &resp); err == nil {
				result.Albums = resp.Albums
				result.Total += resp.Total
			}
		case domain.SearchTypeArtist:
			var resp struct {
				Artists []domain.ArtistInfo `json:"artists"`
				Total   int                 `json:"total"`
			}
			if err := json.Unmarshal([]byte(out), &resp); err == nil {
				result.Artists = resp.Artists
				result.Total += resp.Total
			}
		case domain.SearchTypePlaylist:
			var resp struct {
				Playlists []domain.PlaylistInfo `json:"playlists"`
				Total     int                   `json:"total"`
			}
			if err := json.Unmarshal([]byte(out), &resp); err == nil {
				result.Playlists = resp.Playlists
				result.Total += resp.Total
			}
		}
	}

	return result, nil
}

// Queue is not supported by the native backend.
func (b *Backend) Queue(ctx context.Context, cmd domain.QueueCommand) (*domain.QueueResult, error) {
	return nil, fmt.Errorf("%s: Up Next queue is not available via Apple Events. Consider enabling the accessibility or MusicKit backend", domain.ErrUnsupportedCapability)
}

// Playlists manages playlists via Apple Events.
func (b *Backend) Playlists(ctx context.Context, cmd domain.PlaylistCommand) (*domain.PlaylistResult, error) {
	args := []string{string(cmd.Action)}

	switch cmd.Action {
	case domain.PlaylistActionGet, domain.PlaylistActionDelete, domain.PlaylistActionRename, domain.PlaylistActionCopy:
		if cmd.PlaylistID == "" {
			return nil, fmt.Errorf("playlist_id is required for %s", cmd.Action)
		}
		if err := domain.ValidatePersistentID(cmd.PlaylistID); err != nil {
			return nil, fmt.Errorf("invalid playlist_id: %w", err)
		}
		args = append(args, cmd.PlaylistID)
		if cmd.NewName != "" {
			args = append(args, applescript.SanitizeArg(cmd.NewName))
		}
	case domain.PlaylistActionCreate, domain.PlaylistActionCreateFolder:
		if cmd.NewName == "" {
			return nil, fmt.Errorf("new_name is required for %s", cmd.Action)
		}
		args = append(args, applescript.SanitizeArg(cmd.NewName))
	case domain.PlaylistActionAddTracks, domain.PlaylistActionRemoveTracks:
		if cmd.PlaylistID == "" {
			return nil, fmt.Errorf("playlist_id is required for %s", cmd.Action)
		}
		if err := domain.ValidatePersistentID(cmd.PlaylistID); err != nil {
			return nil, fmt.Errorf("invalid playlist_id: %w", err)
		}
		args = append(args, cmd.PlaylistID)
		for _, tid := range cmd.TrackIDs {
			if err := domain.ValidatePersistentID(tid); err != nil {
				return nil, fmt.Errorf("invalid track_id %q: %w", tid, err)
			}
			args = append(args, tid)
		}
	case domain.PlaylistActionMoveToFolder:
		if cmd.PlaylistID == "" || cmd.TargetFolderID == "" {
			return nil, fmt.Errorf("playlist_id and target_folder_id are required for move_to_folder")
		}
		args = append(args, cmd.PlaylistID, cmd.TargetFolderID)
	}

	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("playlists.applescript"), args...)
	if err != nil {
		return &domain.PlaylistResult{Success: false, BackendName: b.Name()}, fmt.Errorf("playlist command failed: %w", err)
	}

	var result domain.PlaylistResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return &domain.PlaylistResult{Success: false, BackendName: b.Name()}, fmt.Errorf("failed to parse playlist result: %w", err)
	}
	result.BackendName = b.Name()

	return &result, nil
}

// Library manages the music library via Apple Events.
func (b *Backend) Library(ctx context.Context, cmd domain.LibraryCommand) (*domain.LibraryResult, error) {
	args := []string{string(cmd.Action)}

	if cmd.Query != "" {
		args = append(args, applescript.SanitizeArg(cmd.Query))
	}

	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("library.applescript"), args...)
	if err != nil {
		return &domain.LibraryResult{Success: false, BackendName: b.Name()}, fmt.Errorf("library command failed: %w", err)
	}

	var result domain.LibraryResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return &domain.LibraryResult{Success: false, BackendName: b.Name()}, fmt.Errorf("failed to parse library result: %w", err)
	}
	result.BackendName = b.Name()

	return &result, nil
}

// Favorites manages favorites and ratings.
func (b *Backend) Favorites(ctx context.Context, cmd domain.FavoriteCommand) (*domain.FavoriteResult, error) {
	args := []string{string(cmd.Action)}

	if cmd.TargetID != "" {
		if err := domain.ValidatePersistentID(cmd.TargetID); err != nil {
			return nil, fmt.Errorf("invalid target_id: %w", err)
		}
		args = append(args, cmd.TargetID)
		if cmd.TargetType != "" {
			args = append(args, cmd.TargetType)
		}
	}

	if cmd.Action == domain.FavoriteActionSetRating && cmd.Rating != nil {
		if err := domain.ValidateRating(*cmd.Rating); err != nil {
			return nil, err
		}
		args = append(args, fmt.Sprintf("%d", *cmd.Rating))
	}

	out, err := b.executor.RunAppleScript(ctx, b.scriptPath("favorites.applescript"), args...)
	if err != nil {
		return &domain.FavoriteResult{Success: false, BackendName: b.Name()}, fmt.Errorf("favorite command failed: %w", err)
	}

	var result domain.FavoriteResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return &domain.FavoriteResult{Success: false, BackendName: b.Name()}, fmt.Errorf("failed to parse favorite result: %w", err)
	}
	result.BackendName = b.Name()

	return &result, nil
}
