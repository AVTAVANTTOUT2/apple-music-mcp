// Package domain defines the core types and interfaces for Apple Music control.
// These types are backend-agnostic and must not import platform-specific packages.
package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// --- Player State ---

// PlayerState represents the current state of the Music.app player.
type PlayerState string

const (
	PlayerStateStopped     PlayerState = "stopped"
	PlayerStatePlaying     PlayerState = "playing"
	PlayerStatePaused      PlayerState = "paused"
	PlayerStateFastForward PlayerState = "fast_forwarding"
	PlayerStateRewinding   PlayerState = "rewinding"
)

// RepeatMode represents the repeat mode of the player.
type RepeatMode string

const (
	RepeatOff RepeatMode = "off"
	RepeatOne RepeatMode = "one"
	RepeatAll RepeatMode = "all"
)

// ShuffleMode represents the shuffle mode of the player.
type ShuffleMode string

const (
	ShuffleOff    ShuffleMode = "off"
	ShuffleSongs  ShuffleMode = "songs"
	ShuffleAlbums ShuffleMode = "albums"
)

// PlaybackStatus is a snapshot of the current playback state.
type PlaybackStatus struct {
	AppRunning     bool        `json:"app_running"`
	PlayerState    PlayerState `json:"player_state"`
	Track          *TrackInfo  `json:"track,omitempty"`
	Volume         int         `json:"volume"` // 0-100
	Muted          bool        `json:"muted"`
	ShuffleMode    ShuffleMode `json:"shuffle_mode"`
	RepeatMode     RepeatMode  `json:"repeat_mode"`
	Position       float64     `json:"position"` // seconds
	BackendName    string      `json:"backend"`
	CapabilityMode string      `json:"capability_mode"` // native, accessibility, safari, emulated
}

// --- Track ---

// TrackInfo holds metadata about a single track.
type TrackInfo struct {
	Name         string  `json:"name"`
	Artist       string  `json:"artist"`
	Album        string  `json:"album"`
	AlbumArtist  string  `json:"album_artist,omitempty"`
	Composer     string  `json:"composer,omitempty"`
	Genre        string  `json:"genre,omitempty"`
	Year         int     `json:"year,omitempty"`
	TrackNumber  int     `json:"track_number,omitempty"`
	TrackCount   int     `json:"track_count,omitempty"`
	DiscNumber   int     `json:"disc_number,omitempty"`
	DiscCount    int     `json:"disc_count,omitempty"`
	Duration     float64 `json:"duration"`
	PersistentID string  `json:"persistent_id"`
	DatabaseID   int     `json:"database_id,omitempty"`
	PlayedCount  int     `json:"played_count,omitempty"`
	Rating       int     `json:"rating,omitempty"` // 0-100
	Favorited    bool    `json:"favorited"`
	Disliked     bool    `json:"disliked,omitempty"`
	BPM          int     `json:"bpm,omitempty"`
	BitRate      int     `json:"bit_rate,omitempty"`
	SampleRate   int     `json:"sample_rate,omitempty"`
	Size         int64   `json:"size,omitempty"`
	Kind         string  `json:"kind,omitempty"`
	CloudStatus  string  `json:"cloud_status,omitempty"`
	HasArtwork   bool    `json:"has_artwork"`
	ArtworkData  []byte  `json:"-"` // raw artwork bytes, not serialized to JSON
	ReleaseDate  string  `json:"release_date,omitempty"`
	DateAdded    string  `json:"date_added,omitempty"`
	MediaKind    string  `json:"media_kind,omitempty"` // song, music_video, etc.
}

// --- Playlist ---

// PlaylistKind represents the type of a playlist.
type PlaylistKind string

const (
	PlaylistKindUser         PlaylistKind = "user"
	PlaylistKindFolder       PlaylistKind = "folder"
	PlaylistKindLibrary      PlaylistKind = "library"
	PlaylistKindSubscription PlaylistKind = "subscription"
	PlaylistKindRadio        PlaylistKind = "radio"
	PlaylistKindAudioCD      PlaylistKind = "audio_cd"
	PlaylistKindSmart        PlaylistKind = "smart"
	PlaylistKindGenius       PlaylistKind = "genius"
)

// PlaylistInfo holds metadata about a playlist.
type PlaylistInfo struct {
	Name         string       `json:"name"`
	PersistentID string       `json:"persistent_id"`
	Kind         PlaylistKind `json:"kind"`
	TrackCount   int          `json:"track_count"`
	Duration     int          `json:"duration,omitempty"` // seconds
	Size         int64        `json:"size,omitempty"`     // bytes
	Favorited    bool         `json:"favorited,omitempty"`
	Disliked     bool         `json:"disliked,omitempty"`
	ParentID     string       `json:"parent_id,omitempty"`
	Smart        bool         `json:"smart"`
	Genius       bool         `json:"genius,omitempty"`
	Visible      bool         `json:"visible"`
	Description  string       `json:"description,omitempty"`
}

// --- Queue ---

// QueueItem represents an item in the Up Next queue.
type QueueItem struct {
	Index        int        `json:"index"`
	Track        *TrackInfo `json:"track"`
	WillPlayNext bool       `json:"will_play_next"` // true if this is the very next item
}

// QueueStatus describes the Up Next queue.
type QueueStatus struct {
	Items          []QueueItem `json:"items"`
	TotalCount     int         `json:"total_count"`
	Autoplay       bool        `json:"autoplay,omitempty"`
	BackendName    string      `json:"backend"`
	CapabilityMode string      `json:"capability_mode"` // native, accessibility, safari, emulated
}

// --- Search ---

// SearchScope defines where to search.
type SearchScope string

const (
	SearchScopeAll     SearchScope = "all"
	SearchScopeLibrary SearchScope = "library"
	SearchScopeCatalog SearchScope = "catalog"
)

// SearchType filters search results.
type SearchType string

const (
	SearchTypeTrack    SearchType = "track"
	SearchTypeAlbum    SearchType = "album"
	SearchTypeArtist   SearchType = "artist"
	SearchTypePlaylist SearchType = "playlist"
)

// SearchQuery represents a search request.
type SearchQuery struct {
	Query             string       `json:"query"`
	Scope             SearchScope  `json:"scope"`
	Types             []SearchType `json:"types"`
	Limit             int          `json:"limit"`
	Offset            int          `json:"offset"`
	Exact             bool         `json:"exact"`
	IncludeUnplayable bool         `json:"include_unplayable"`
}

// SearchResult holds search results.
type SearchResult struct {
	Tracks      []TrackInfo    `json:"tracks,omitempty"`
	Albums      []AlbumInfo    `json:"albums,omitempty"`
	Artists     []ArtistInfo   `json:"artists,omitempty"`
	Playlists   []PlaylistInfo `json:"playlists,omitempty"`
	Total       int            `json:"total"`
	BackendName string         `json:"backend"`
}

// --- Album ---

// AlbumInfo holds album metadata.
type AlbumInfo struct {
	Name         string `json:"name"`
	Artist       string `json:"artist"`
	PersistentID string `json:"persistent_id,omitempty"`
	TrackCount   int    `json:"track_count,omitempty"`
	Year         int    `json:"year,omitempty"`
	Genre        string `json:"genre,omitempty"`
	Favorited    bool   `json:"favorited,omitempty"`
	Rating       int    `json:"rating,omitempty"`
}

// --- Artist ---

// ArtistInfo holds artist metadata.
type ArtistInfo struct {
	Name         string `json:"name"`
	PersistentID string `json:"persistent_id,omitempty"`
	AlbumCount   int    `json:"album_count,omitempty"`
	TrackCount   int    `json:"track_count,omitempty"`
}

// --- AirPlay ---

// AirPlayDevice represents an AirPlay output device.
type AirPlayDevice struct {
	Name          string `json:"name"`
	Kind          string `json:"kind"` // computer, homepod, apple_tv, bluetooth, etc.
	Active        bool   `json:"active"`
	Available     bool   `json:"available"`
	Selected      bool   `json:"selected"`
	SupportsAudio bool   `json:"supports_audio"`
	Protected     bool   `json:"protected"`
	Volume        int    `json:"volume,omitempty"` // 0-100
}

// --- Capabilities ---

// CapabilitySet describes what the current backend can do.
type CapabilitySet struct {
	MacOSVersion            string   `json:"macos_version"`
	MacOSArch               string   `json:"macos_arch"`
	MusicAppVersion         string   `json:"music_app_version"`
	MusicAppInstalled       bool     `json:"music_app_installed"`
	AutomationGranted       bool     `json:"automation_granted"`
	AccessibilityGranted    bool     `json:"accessibility_granted"`
	AppleMusicSubscriber    bool     `json:"apple_music_subscriber,omitempty"`
	BackendName             string   `json:"backend_name"`
	QueueBackendName        string   `json:"queue_backend_name,omitempty"`
	CatalogBackendName      string   `json:"catalog_backend_name,omitempty"`
	NativeCapabilities      []string `json:"native_capabilities"`
	FallbackCapabilities    []string `json:"fallback_capabilities,omitempty"`
	UnavailableCapabilities []string `json:"unavailable_capabilities,omitempty"`
	Restrictions            []string `json:"restrictions,omitempty"`
}

// --- Playback Commands ---

// PlaybackAction is an action to perform on the player.
type PlaybackAction string

const (
	ActionOpen           PlaybackAction = "open"
	ActionPlay           PlaybackAction = "play"
	ActionPause          PlaybackAction = "pause"
	ActionToggle         PlaybackAction = "toggle"
	ActionStop           PlaybackAction = "stop"
	ActionNext           PlaybackAction = "next"
	ActionPrevious       PlaybackAction = "previous"
	ActionRestartCurrent PlaybackAction = "restart_current"
	ActionSeekAbsolute   PlaybackAction = "seek_absolute"
	ActionSeekRelative   PlaybackAction = "seek_relative"
	ActionPlayTrack      PlaybackAction = "play_track"
	ActionPlayAlbum      PlaybackAction = "play_album"
	ActionPlayArtist     PlaybackAction = "play_artist"
	ActionPlayPlaylist   PlaybackAction = "play_playlist"
	ActionPlayURL        PlaybackAction = "play_url"
	ActionReveal         PlaybackAction = "reveal"
)

// PlaybackCommand is a command to the playback controller.
type PlaybackCommand struct {
	Action       PlaybackAction `json:"action"`
	TargetID     string         `json:"target_id,omitempty"`     // persistent ID, database ID, or URL
	TargetType   string         `json:"target_type,omitempty"`   // track, album, artist, playlist, url
	SeekPosition float64        `json:"seek_position,omitempty"` // seconds
	Once         bool           `json:"once,omitempty"`          // play once and stop
}

// PlaybackResult is returned after a playback command.
type PlaybackResult struct {
	Success     bool   `json:"success"`
	Action      string `json:"action"`
	Message     string `json:"message,omitempty"`
	BackendName string `json:"backend"`
}

// --- Preferences Commands ---

// PreferenceAction is a preference operation.
type PreferenceAction string

const (
	PrefActionGet        PreferenceAction = "get"
	PrefActionSetVolume  PreferenceAction = "set_volume"
	PrefActionSetShuffle PreferenceAction = "set_shuffle"
	PrefActionSetRepeat  PreferenceAction = "set_repeat"
	PrefActionGetOutputs PreferenceAction = "get_outputs"
	PrefActionSetOutput  PreferenceAction = "set_output"
)

// PreferenceCommand is a command to modify player preferences.
type PreferenceCommand struct {
	Action      PreferenceAction `json:"action"`
	Volume      *int             `json:"volume,omitempty"` // 0-100
	ShuffleMode *ShuffleMode     `json:"shuffle_mode,omitempty"`
	RepeatMode  *RepeatMode      `json:"repeat_mode,omitempty"`
	OutputName  *string          `json:"output_name,omitempty"` // AirPlay device name
}

// PreferenceResult holds the result of a preference operation.
type PreferenceResult struct {
	Success     bool            `json:"success"`
	Volume      int             `json:"volume,omitempty"`
	ShuffleMode ShuffleMode     `json:"shuffle_mode,omitempty"`
	RepeatMode  RepeatMode      `json:"repeat_mode,omitempty"`
	Outputs     []AirPlayDevice `json:"outputs,omitempty"`
	BackendName string          `json:"backend"`
}

// --- Queue Commands ---

// QueueAction is an operation on the Up Next queue.
type QueueAction string

const (
	QueueActionList        QueueAction = "list"
	QueueActionAddNext     QueueAction = "add_next"
	QueueActionAddLater    QueueAction = "add_later"
	QueueActionRemove      QueueAction = "remove"
	QueueActionMove        QueueAction = "move"
	QueueActionJump        QueueAction = "jump"
	QueueActionClear       QueueAction = "clear"
	QueueActionGetAutoplay QueueAction = "get_autoplay"
	QueueActionSetAutoplay QueueAction = "set_autoplay"
)

// QueueCommand is a command to modify the Up Next queue.
type QueueCommand struct {
	Action    QueueAction `json:"action"`
	TrackID   string      `json:"track_id,omitempty"`   // persistent ID
	FromIndex int         `json:"from_index,omitempty"` // for move
	ToIndex   int         `json:"to_index,omitempty"`   // for move
	Autoplay  *bool       `json:"autoplay,omitempty"`   // for set_autoplay
}

// QueueResult holds the result of a queue operation.
type QueueResult struct {
	Success     bool         `json:"success"`
	Queue       *QueueStatus `json:"queue,omitempty"`
	BackendName string       `json:"backend"`
}

// --- Playlist Commands ---

// PlaylistAction enumerates playlist operations.
type PlaylistAction string

const (
	PlaylistActionList          PlaylistAction = "list"
	PlaylistActionGet           PlaylistAction = "get"
	PlaylistActionCreate        PlaylistAction = "create"
	PlaylistActionRename        PlaylistAction = "rename"
	PlaylistActionDelete        PlaylistAction = "delete"
	PlaylistActionAddTracks     PlaylistAction = "add_tracks"
	PlaylistActionRemoveTracks  PlaylistAction = "remove_tracks"
	PlaylistActionReorderTracks PlaylistAction = "reorder_tracks"
	PlaylistActionCopy          PlaylistAction = "copy"
	PlaylistActionListFolders   PlaylistAction = "list_folders"
	PlaylistActionCreateFolder  PlaylistAction = "create_folder"
	PlaylistActionMoveToFolder  PlaylistAction = "move_to_folder"
)

// PlaylistCommand is a command for playlist management.
type PlaylistCommand struct {
	Action         PlaylistAction `json:"action"`
	PlaylistID     string         `json:"playlist_id,omitempty"` // persistent ID or name
	NewName        string         `json:"new_name,omitempty"`
	TrackIDs       []string       `json:"track_ids,omitempty"` // persistent IDs
	FromIndex      int            `json:"from_index,omitempty"`
	ToIndex        int            `json:"to_index,omitempty"`
	TargetFolderID string         `json:"target_folder_id,omitempty"`
}

// PlaylistResult is returned from playlist operations.
type PlaylistResult struct {
	Success     bool           `json:"success"`
	Playlists   []PlaylistInfo `json:"playlists,omitempty"`
	Playlist    *PlaylistInfo  `json:"playlist,omitempty"`
	Tracks      []TrackInfo    `json:"tracks,omitempty"`
	Folders     []PlaylistInfo `json:"folders,omitempty"`
	BackendName string         `json:"backend"`
}

// --- Library Commands ---

// LibraryAction enumerates library operations.
type LibraryAction string

const (
	LibraryActionSearch         LibraryAction = "search"
	LibraryActionAdd            LibraryAction = "add"
	LibraryActionRemove         LibraryAction = "remove"
	LibraryActionRecentlyAdded  LibraryAction = "recently_added"
	LibraryActionRecentlyPlayed LibraryAction = "recently_played"
	LibraryActionListTracks     LibraryAction = "list_tracks"
	LibraryActionListAlbums     LibraryAction = "list_albums"
	LibraryActionListArtists    LibraryAction = "list_artists"
	LibraryActionListGenres     LibraryAction = "list_genres"
)

// LibraryCommand is a command for library management.
type LibraryCommand struct {
	Action   LibraryAction `json:"action"`
	Query    string        `json:"query,omitempty"`
	TrackIDs []string      `json:"track_ids,omitempty"`
	Limit    int           `json:"limit,omitempty"`
	Offset   int           `json:"offset,omitempty"`
}

// LibraryResult holds the result of a library operation.
type LibraryResult struct {
	Success     bool         `json:"success"`
	Tracks      []TrackInfo  `json:"tracks,omitempty"`
	Albums      []AlbumInfo  `json:"albums,omitempty"`
	Artists     []ArtistInfo `json:"artists,omitempty"`
	Genres      []string     `json:"genres,omitempty"`
	Total       int          `json:"total,omitempty"`
	BackendName string       `json:"backend"`
}

// --- Favorite Commands ---

// FavoriteAction enumerates favorite/rating operations.
type FavoriteAction string

const (
	FavoriteActionGet         FavoriteAction = "get"
	FavoriteActionList        FavoriteAction = "list"
	FavoriteActionFavorite    FavoriteAction = "favorite"
	FavoriteActionUnfavorite  FavoriteAction = "unfavorite"
	FavoriteActionSuggestLess FavoriteAction = "suggest_less" // set disliked
	FavoriteActionSetRating   FavoriteAction = "set_rating"
	FavoriteActionClearRating FavoriteAction = "clear_rating"
)

// FavoriteCommand is a command for favorite/rating operations.
type FavoriteCommand struct {
	Action     FavoriteAction `json:"action"`
	TargetID   string         `json:"target_id,omitempty"`   // persistent ID
	TargetType string         `json:"target_type,omitempty"` // track, album, artist
	Rating     *int           `json:"rating,omitempty"`      // 0-100
}

// FavoriteResult holds the result of a favorite operation.
type FavoriteResult struct {
	Success     bool        `json:"success"`
	Favorited   bool        `json:"favorited,omitempty"`
	Rating      int         `json:"rating,omitempty"`
	Favorites   []TrackInfo `json:"favorites,omitempty"`
	BackendName string      `json:"backend"`
}

// --- Response Envelope ---

// APIResponse is the standard response envelope for all MCP tool calls.
type APIResponse struct {
	OK             bool            `json:"ok"`
	Backend        string          `json:"backend"`
	CapabilityMode string          `json:"capability_mode"`
	Data           json.RawMessage `json:"data,omitempty"`
	Warnings       []string        `json:"warnings,omitempty"`
	RequestID      string          `json:"request_id"`
}

// APIError is the standard error structure.
type APIError struct {
	OK        bool        `json:"ok"`
	Error     ErrorDetail `json:"error"`
	Backend   string      `json:"backend"`
	RequestID string      `json:"request_id"`
}

// ErrorDetail describes an error with metadata.
type ErrorDetail struct {
	Code           ErrorCode `json:"code"`
	Message        string    `json:"message"`
	Recoverable    bool      `json:"recoverable"`
	RequiredAction string    `json:"required_action,omitempty"`
}

// Error implements the error interface for ErrorCode.
func (e ErrorCode) Error() string {
	return string(e)
}

// ErrorCode is a typed error identifier.
type ErrorCode string

const (
	ErrMusicNotInstalled        ErrorCode = "music_not_installed"
	ErrMusicNotRunning          ErrorCode = "music_not_running"
	ErrNoTrackLoaded            ErrorCode = "no_track_loaded"
	ErrPermissionDenied         ErrorCode = "permission_denied"
	ErrAccessibilityRequired    ErrorCode = "accessibility_permission_required"
	ErrAuthenticationRequired   ErrorCode = "authentication_required"
	ErrSubscriptionRequired     ErrorCode = "subscription_required"
	ErrUnsupportedCapability    ErrorCode = "unsupported_capability"
	ErrUnsupportedMacOSVersion  ErrorCode = "unsupported_macos_version"
	ErrNotFound                 ErrorCode = "not_found"
	ErrAmbiguousResult          ErrorCode = "ambiguous_result"
	ErrInvalidReference         ErrorCode = "invalid_reference"
	ErrInvalidArgument          ErrorCode = "invalid_argument"
	ErrTimeout                  ErrorCode = "timeout"
	ErrBackendUnavailable       ErrorCode = "backend_unavailable"
	ErrAutomationError          ErrorCode = "automation_error"
	ErrConflict                 ErrorCode = "conflict"
	ErrDestructiveActionBlocked ErrorCode = "destructive_action_blocked"
)

// --- Helpers ---

// NewAPIResponse creates a standard success response.
func NewAPIResponse(data interface{}, backend string, capabilityMode string, requestID string) APIResponse {
	raw, _ := json.Marshal(data)
	return APIResponse{
		OK:             true,
		Backend:        backend,
		CapabilityMode: capabilityMode,
		Data:           raw,
		Warnings:       []string{},
		RequestID:      requestID,
	}
}

// NewAPIError creates a standard error response.
func NewAPIError(code ErrorCode, message string, recoverable bool, backend string, requestID string) APIError {
	return APIError{
		OK: false,
		Error: ErrorDetail{
			Code:        code,
			Message:     message,
			Recoverable: recoverable,
		},
		Backend:   backend,
		RequestID: requestID,
	}
}

// NewAPIErrorWithAction creates an error with a required user action.
func NewAPIErrorWithAction(code ErrorCode, message string, action string, backend string, requestID string) APIError {
	return APIError{
		OK: false,
		Error: ErrorDetail{
			Code:           code,
			Message:        message,
			Recoverable:    true,
			RequiredAction: action,
		},
		Backend:   backend,
		RequestID: requestID,
	}
}

// FormatDuration formats seconds as MM:SS or H:MM:SS.
func FormatDuration(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}
	s := int(seconds)
	h := s / 3600
	m := (s % 3600) / 60
	sec := s % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%d:%02d", m, sec)
}

// ValidateSearchQuery validates and normalizes a search query.
func ValidateSearchQuery(q *SearchQuery) error {
	if q.Query == "" {
		return fmt.Errorf("search query must not be empty")
	}
	if len(q.Query) > 1024 {
		return fmt.Errorf("search query too long (max 1024 characters)")
	}
	if q.Limit <= 0 {
		q.Limit = 25
	}
	if q.Limit > 200 {
		q.Limit = 200
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Scope == "" {
		q.Scope = SearchScopeAll
	}
	if len(q.Types) == 0 {
		q.Types = []SearchType{SearchTypeTrack, SearchTypeAlbum, SearchTypeArtist}
	}
	return nil
}

// ValidatePersistentID checks if a string looks like a valid persistent ID (hex).
func ValidatePersistentID(id string) error {
	if id == "" {
		return fmt.Errorf("persistent ID must not be empty")
	}
	if len(id) > 64 {
		return fmt.Errorf("persistent ID too long")
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return fmt.Errorf("persistent ID contains invalid character: %c", c)
		}
	}
	return nil
}

// ValidateVolume validates a volume value between 0 and 100.
func ValidateVolume(v int) error {
	if v < 0 || v > 100 {
		return fmt.Errorf("volume must be between 0 and 100, got %d", v)
	}
	return nil
}

// ValidateRating validates a rating value between 0 and 100.
func ValidateRating(r int) error {
	if r < 0 || r > 100 {
		return fmt.Errorf("rating must be between 0 and 100, got %d", r)
	}
	return nil
}

// Now returns the current time, used for request IDs and logging.
func Now() time.Time {
	return time.Now()
}
