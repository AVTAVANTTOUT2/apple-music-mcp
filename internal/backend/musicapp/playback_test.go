package musicapp

import (
	"testing"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
)

func TestNormalizePlaybackCommand_PlayTrack(t *testing.T) {
	t.Parallel()

	cmd := normalizePlaybackCommand(domain.PlaybackCommand{
		Action:   domain.ActionPlayTrack,
		TargetID: "ABCD1234",
	})

	if cmd.Action != domain.ActionPlay {
		t.Fatalf("action = %q, want %q", cmd.Action, domain.ActionPlay)
	}
	if cmd.TargetType != "track" {
		t.Fatalf("target_type = %q, want track", cmd.TargetType)
	}
}

func TestNormalizePlaybackCommand_PlayURLUnchanged(t *testing.T) {
	t.Parallel()

	cmd := normalizePlaybackCommand(domain.PlaybackCommand{
		Action:   domain.ActionPlayURL,
		TargetID: "https://music.apple.com/fr/artist/werenoi/1551258720",
	})

	if cmd.Action != domain.ActionPlayURL {
		t.Fatalf("action = %q, want play_url", cmd.Action)
	}
}
