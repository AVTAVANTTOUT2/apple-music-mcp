package applescript

import "testing"

func TestNormalizeJSONNumbers_FrenchLocaleDecimals(t *testing.T) {
	t.Parallel()

	raw := `{"duration":173,143005371094,"position":56,991}`
	want := `{"duration":173.143005371094,"position":56.991}`
	got := NormalizeJSONNumbers(raw)
	if got != want {
		t.Fatalf("NormalizeJSONNumbers() = %q, want %q", got, want)
	}
}

func TestGetKVFloat_LocaleDecimal(t *testing.T) {
	t.Parallel()

	kv := map[string]string{
		"POSITION":       "56,99100112915",
		"TRACK_DURATION": "173,143005371094",
	}
	if got := GetKVFloat(kv, "POSITION"); got < 56.99 || got > 56.992 {
		t.Fatalf("POSITION float = %v, want ~56.991", got)
	}
	if got := GetKVFloat(kv, "TRACK_DURATION"); got < 173.14 || got > 173.15 {
		t.Fatalf("TRACK_DURATION float = %v, want ~173.143", got)
	}
}

func TestParseJSONOutput_WithLocaleDecimals(t *testing.T) {
	t.Parallel()

	var resp struct {
		Success bool `json:"success"`
		Tracks  []struct {
			Duration float64 `json:"duration"`
		} `json:"tracks"`
	}

	raw := `{"success":true,"tracks":[{"duration":158,050003051758}]}`
	if err := ParseJSONOutput(raw, &resp); err != nil {
		t.Fatalf("ParseJSONOutput() error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if len(resp.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(resp.Tracks))
	}
	if resp.Tracks[0].Duration < 158.04 || resp.Tracks[0].Duration > 158.06 {
		t.Fatalf("duration = %v, want ~158.05", resp.Tracks[0].Duration)
	}
}
