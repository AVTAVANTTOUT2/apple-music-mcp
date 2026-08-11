// Package applescript provides safe execution and parsing of AppleScript output.
package applescript

import (
	"strconv"
	"strings"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
)

// ParseKV parses AppleScript output in KEY:VALUE format (one per line).
// The value is everything after the first colon on each line.
func ParseKV(output string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		idx := strings.Index(line, ":")
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		result[key] = value
	}
	return result
}

// GetKV returns a value from a KV map with a default.
func GetKV(kv map[string]string, key, defaultVal string) string {
	if v, ok := kv[key]; ok {
		return v
	}
	return defaultVal
}

// GetKVBool returns a boolean value.
func GetKVBool(kv map[string]string, key string) bool {
	v := GetKV(kv, key, "false")
	return v == "true"
}

// GetKVInt returns an integer value.
func GetKVInt(kv map[string]string, key string) int {
	v := GetKV(kv, key, "0")
	n, _ := strconv.Atoi(v)
	return n
}

// GetKVFloat returns a float value.
func GetKVFloat(kv map[string]string, key string) float64 {
	v := GetKV(kv, key, "0")
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

// ParsePlayerState parses a KV map into a PlaybackStatus.
func ParsePlayerState(kv map[string]string) *domain.PlaybackStatus {
	status := &domain.PlaybackStatus{
		AppRunning:     GetKVBool(kv, "APP_RUNNING"),
		PlayerState:    domain.PlayerState(GetKV(kv, "PLAYER_STATE", "stopped")),
		Volume:         GetKVInt(kv, "VOLUME"),
		Muted:          GetKVBool(kv, "MUTED"),
		ShuffleMode:    domain.ShuffleMode(GetKV(kv, "SHUFFLE_MODE", "off")),
		RepeatMode:     domain.RepeatMode(GetKV(kv, "REPEAT_MODE", "off")),
		Position:       GetKVFloat(kv, "POSITION"),
		BackendName:    "musicapp",
		CapabilityMode: "native",
	}

	if GetKVBool(kv, "HAS_TRACK") {
		status.Track = &domain.TrackInfo{
			Name:         GetKV(kv, "TRACK_NAME", "Unknown"),
			Artist:       GetKV(kv, "TRACK_ARTIST", ""),
			Album:        GetKV(kv, "TRACK_ALBUM", ""),
			Duration:     GetKVFloat(kv, "TRACK_DURATION"),
			PersistentID: GetKV(kv, "TRACK_PID", ""),
			DatabaseID:   GetKVInt(kv, "TRACK_DBID"),
			Favorited:    GetKVBool(kv, "TRACK_FAVORITED"),
			Rating:       GetKVInt(kv, "TRACK_RATING"),
			Year:         GetKVInt(kv, "TRACK_YEAR"),
			Genre:        GetKV(kv, "TRACK_GENRE", ""),
			TrackNumber:  GetKVInt(kv, "TRACK_TRACK_NUMBER"),
			TrackCount:   GetKVInt(kv, "TRACK_TRACK_COUNT"),
			DiscNumber:   GetKVInt(kv, "TRACK_DISC_NUMBER"),
			DiscCount:    GetKVInt(kv, "TRACK_DISC_COUNT"),
			Composer:     GetKV(kv, "TRACK_COMPOSER", ""),
			PlayedCount:  GetKVInt(kv, "TRACK_PLAYED_COUNT"),
			Kind:         GetKV(kv, "TRACK_KIND", ""),
			HasArtwork:   GetKVBool(kv, "TRACK_HAS_ARTWORK"),
		}
	}

	return status
}

// ParsePlaybackResult parses a KV map into a PlaybackResult.
func ParsePlaybackResult(kv map[string]string) *domain.PlaybackResult {
	return &domain.PlaybackResult{
		Success:     GetKVBool(kv, "SUCCESS"),
		Action:      GetKV(kv, "ACTION", ""),
		Message:     GetKV(kv, "MESSAGE", ""),
		BackendName: "musicapp",
	}
}

// ParsePreferencesResult parses KV output from preferences script.
func ParsePreferencesResult(kv map[string]string) *domain.PreferenceResult {
	return &domain.PreferenceResult{
		Success:     GetKVBool(kv, "SUCCESS"),
		Volume:      GetKVInt(kv, "VOLUME"),
		ShuffleMode: domain.ShuffleMode(GetKV(kv, "SHUFFLE_MODE", "off")),
		RepeatMode:  domain.RepeatMode(GetKV(kv, "REPEAT_MODE", "off")),
		BackendName: "musicapp",
	}
}

// ParseTrackLines parses track info from KV format.
func ParseTrackLines(kv map[string]string, prefix string) *domain.TrackInfo {
	return &domain.TrackInfo{
		Name:         GetKV(kv, prefix+"NAME", ""),
		Artist:       GetKV(kv, prefix+"ARTIST", ""),
		Album:        GetKV(kv, prefix+"ALBUM", ""),
		Duration:     GetKVFloat(kv, prefix+"DURATION"),
		PersistentID: GetKV(kv, prefix+"PID", ""),
		DatabaseID:   GetKVInt(kv, prefix+"DBID"),
		Favorited:    GetKVBool(kv, prefix+"FAVORITED"),
		Rating:       GetKVInt(kv, prefix+"RATING"),
		Year:         GetKVInt(kv, prefix+"YEAR"),
		Genre:        GetKV(kv, prefix+"GENRE", ""),
		TrackNumber:  GetKVInt(kv, prefix+"TRACK_NUMBER"),
		TrackCount:   GetKVInt(kv, prefix+"TRACK_COUNT"),
	}
}
