// Package livetest runs non-destructive integration checks against the real Music.app.
// Enable with APPLE_MUSIC_MCP_LIVE_TESTS=1.
package livetest

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/backend"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/backend/musicapp"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
	sembed "github.com/AVTAVANTTOUT2/apple-music-mcp/scripts"
)

const (
	defaultSearchArtist = "Werenoi"
	defaultPlaylistName = "Jarvis MCP Live Test"
)

// Options tune the live test run.
type Options struct {
	SearchQuery  string
	PlaylistName string
	KeepPlaylist bool
}

// Result summarizes a live test run.
type Result struct {
	Passed int
	Failed int
	Lines  []string
}

// Run executes the full live workflow used by JARVIS integrations.
func Run(ctx context.Context, opts Options) (*Result, error) {
	if opts.SearchQuery == "" {
		opts.SearchQuery = envOr("APPLE_MUSIC_MCP_LIVE_SEARCH", defaultSearchArtist)
	}
	if opts.PlaylistName == "" {
		opts.PlaylistName = envOr("APPLE_MUSIC_MCP_LIVE_PLAYLIST", defaultPlaylistName)
	}
	if !opts.KeepPlaylist {
		opts.KeepPlaylist = os.Getenv("APPLE_MUSIC_MCP_LIVE_KEEP_PLAYLIST") == "1"
	}

	logger := logging.NewLogger(logging.LevelInfo)
	be := musicapp.NewBackend(logger)

	scriptsDir, err := musicapp.ExtractEmbeddedScripts(sembed.FS)
	if err != nil {
		return nil, fmt.Errorf("extract embedded scripts: %w", err)
	}
	be.SetScriptsDir(scriptsDir)

	result := &Result{}
	step := func(name string, fn func() error) {
		if err := fn(); err != nil {
			result.Failed++
			result.Lines = append(result.Lines, fmt.Sprintf("❌ %s — %v", name, err))
			return
		}
		result.Passed++
		result.Lines = append(result.Lines, fmt.Sprintf("✅ %s", name))
	}

	step("Music.app automation available", func() error {
		if !be.IsAvailable(ctx) {
			return fmt.Errorf("automation permission missing for Music.app")
		}
		return nil
	})

	var track domain.TrackInfo
	step(fmt.Sprintf("search library for %q", opts.SearchQuery), func() error {
		search, err := be.Search(ctx, domain.SearchQuery{
			Query: opts.SearchQuery,
			Types: []domain.SearchType{domain.SearchTypeTrack},
			Limit: 5,
		})
		if err != nil {
			return err
		}
		if search.Total == 0 || len(search.Tracks) == 0 {
			return fmt.Errorf("no tracks found for query %q", opts.SearchQuery)
		}
		track = search.Tracks[0]
		if track.PersistentID == "" {
			return fmt.Errorf("first track has empty persistent_id")
		}
		result.Lines = append(result.Lines, fmt.Sprintf("   → track: %s — %s (%s)", track.Name, track.Artist, track.PersistentID))
		return nil
	})

	if track.PersistentID == "" {
		return result, nil
	}

	step("play track via play_track", func() error {
		playback, err := be.Playback(ctx, domain.PlaybackCommand{
			Action:     domain.ActionPlayTrack,
			TargetID:   track.PersistentID,
			TargetType: "track",
		})
		if err != nil {
			return err
		}
		if !playback.Success {
			return fmt.Errorf("playback returned success=false")
		}
		time.Sleep(2 * time.Second)
		return nil
	})

	step("read playback state", func() error {
		state, err := be.GetPlayerState(ctx)
		if err != nil {
			return err
		}
		if state.PlayerState != domain.PlayerStatePlaying && state.PlayerState != domain.PlayerStatePaused {
			return fmt.Errorf("unexpected player state: %s", state.PlayerState)
		}
		if state.Track == nil || state.Track.PersistentID == "" {
			return fmt.Errorf("current track metadata unavailable")
		}
		return nil
	})

	step("favorite current track", func() error {
		fav, err := be.Favorites(ctx, domain.FavoriteCommand{
			Action:     domain.FavoriteActionFavorite,
			TargetID:   track.PersistentID,
			TargetType: "track",
		})
		if err != nil {
			return err
		}
		if !fav.Success || !fav.Favorited {
			return fmt.Errorf("favorite not applied")
		}
		return nil
	})

	var playlistID string
	step(fmt.Sprintf("create playlist %q", opts.PlaylistName), func() error {
		_ = deletePlaylistIfExists(ctx, be, opts.PlaylistName)

		created, err := be.Playlists(ctx, domain.PlaylistCommand{
			Action:  domain.PlaylistActionCreate,
			NewName: opts.PlaylistName,
		})
		if err != nil {
			return err
		}
		if !created.Success || created.Playlist == nil || created.Playlist.PersistentID == "" {
			return fmt.Errorf("playlist creation failed")
		}
		playlistID = created.Playlist.PersistentID
		return nil
	})

	if playlistID != "" {
		step("add track to playlist", func() error {
			added, err := be.Playlists(ctx, domain.PlaylistCommand{
				Action:     domain.PlaylistActionAddTracks,
				PlaylistID: playlistID,
				TrackIDs:   []string{track.PersistentID},
			})
			if err != nil {
				return err
			}
			if !added.Success {
				return fmt.Errorf("add_tracks returned success=false")
			}

			got, err := be.Playlists(ctx, domain.PlaylistCommand{
				Action:     domain.PlaylistActionGet,
				PlaylistID: playlistID,
			})
			if err != nil {
				return err
			}
			if !got.Success || len(got.Tracks) == 0 {
				return fmt.Errorf("playlist has no tracks after add")
			}
			return nil
		})
	}

	if playlistID != "" && !opts.KeepPlaylist {
		step("cleanup test playlist", func() error {
			deleted, err := be.Playlists(ctx, domain.PlaylistCommand{
				Action:     domain.PlaylistActionDelete,
				PlaylistID: playlistID,
			})
			if err != nil {
				return err
			}
			if !deleted.Success {
				return fmt.Errorf("delete playlist failed")
			}
			return nil
		})
	}

	return result, nil
}

func deletePlaylistIfExists(ctx context.Context, be backend.Backend, name string) error {
	list, err := be.Playlists(ctx, domain.PlaylistCommand{Action: domain.PlaylistActionList})
	if err != nil || !list.Success {
		return err
	}
	for _, pl := range list.Playlists {
		if strings.EqualFold(pl.Name, name) {
			_, _ = be.Playlists(ctx, domain.PlaylistCommand{
				Action:     domain.PlaylistActionDelete,
				PlaylistID: pl.PersistentID,
			})
		}
	}
	return nil
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
