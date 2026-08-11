// Package backend defines the core interfaces for Apple Music backends.
package backend

import (
	"context"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
)

// Backend is the primary interface for controlling Apple Music.
// Implementations include: native Apple Events (musicapp),
// Accessibility automation, and Safari/MusicKit web automation.
type Backend interface {
	// Name returns a human-readable identifier for this backend.
	Name() string

	// Capabilities returns what this backend can do on the current system.
	Capabilities(ctx context.Context) (domain.CapabilitySet, error)

	// GetPlayerState returns the current playback state and track info.
	GetPlayerState(ctx context.Context) (*domain.PlaybackStatus, error)

	// Playback executes a playback command.
	Playback(ctx context.Context, cmd domain.PlaybackCommand) (*domain.PlaybackResult, error)

	// Preferences gets or sets player preferences.
	Preferences(ctx context.Context, cmd domain.PreferenceCommand) (*domain.PreferenceResult, error)

	// Search performs a search using this backend.
	Search(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error)

	// Queue manages the Up Next queue.
	// Returns domain.ErrUnsupportedCapability if the backend doesn't support queue operations.
	Queue(ctx context.Context, cmd domain.QueueCommand) (*domain.QueueResult, error)

	// Playlists manages playlists.
	Playlists(ctx context.Context, cmd domain.PlaylistCommand) (*domain.PlaylistResult, error)

	// Library manages the music library.
	Library(ctx context.Context, cmd domain.LibraryCommand) (*domain.LibraryResult, error)

	// Favorites manages favorites and ratings.
	Favorites(ctx context.Context, cmd domain.FavoriteCommand) (*domain.FavoriteResult, error)

	// IsAvailable returns true if this backend can be used right now.
	IsAvailable(ctx context.Context) bool
}

// QueueBackend is an optional interface for backends that support Up Next queue.
type QueueBackend interface {
	Backend
	// QueueBackendName returns the specific queue implementation name.
	QueueBackendName() string
}

// CatalogBackend is an optional interface for backends that support Apple Music catalog.
type CatalogBackend interface {
	Backend
	// CatalogSearch searches the Apple Music catalog.
	CatalogSearch(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error)
	// CatalogBackendName returns the specific catalog implementation name.
	CatalogBackendName() string
}
