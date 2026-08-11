package domain

import (
	"testing"
)

func TestValidatePersistentID_Valid(t *testing.T) {
	tests := []string{
		"ABCD1234EF567890",
		"abcdef0123456789",
		"1A2B3C",
	}
	for _, id := range tests {
		if err := ValidatePersistentID(id); err != nil {
			t.Errorf("expected valid persistent ID %q, got error: %v", id, err)
		}
	}
}

func TestValidatePersistentID_Invalid(t *testing.T) {
	tests := []string{
		"",
		"not-hex!!",
		"with spaces",
		"hello world",
	}
	for _, id := range tests {
		if err := ValidatePersistentID(id); err == nil {
			t.Errorf("expected error for invalid persistent ID %q", id)
		}
	}
}

func TestValidatePersistentID_TooLong(t *testing.T) {
	longID := make([]byte, 65)
	for i := range longID {
		longID[i] = 'A'
	}
	if err := ValidatePersistentID(string(longID)); err == nil {
		t.Error("expected error for too-long persistent ID")
	}
}

func TestValidateVolume(t *testing.T) {
	if err := ValidateVolume(0); err != nil {
		t.Errorf("volume 0 should be valid: %v", err)
	}
	if err := ValidateVolume(50); err != nil {
		t.Errorf("volume 50 should be valid: %v", err)
	}
	if err := ValidateVolume(100); err != nil {
		t.Errorf("volume 100 should be valid: %v", err)
	}
	if err := ValidateVolume(-1); err == nil {
		t.Error("volume -1 should be invalid")
	}
	if err := ValidateVolume(101); err == nil {
		t.Error("volume 101 should be invalid")
	}
}

func TestValidateRating(t *testing.T) {
	if err := ValidateRating(0); err != nil {
		t.Errorf("rating 0 should be valid: %v", err)
	}
	if err := ValidateRating(80); err != nil {
		t.Errorf("rating 80 should be valid: %v", err)
	}
	if err := ValidateRating(100); err != nil {
		t.Errorf("rating 100 should be valid: %v", err)
	}
	if err := ValidateRating(-1); err == nil {
		t.Error("rating -1 should be invalid")
	}
	if err := ValidateRating(101); err == nil {
		t.Error("rating 101 should be invalid")
	}
}

func TestValidateSearchQuery(t *testing.T) {
	// Valid
	q := &SearchQuery{Query: "test"}
	if err := ValidateSearchQuery(q); err != nil {
		t.Errorf("valid search should not error: %v", err)
	}
	if q.Limit != 25 {
		t.Errorf("default limit should be 25, got %d", q.Limit)
	}
	if len(q.Types) == 0 {
		t.Error("default types should be non-empty")
	}

	// Empty query
	q = &SearchQuery{Query: ""}
	if err := ValidateSearchQuery(q); err == nil {
		t.Error("empty query should error")
	}

	// Too long
	longQ := make([]byte, 1025)
	for i := range longQ {
		longQ[i] = 'a'
	}
	q = &SearchQuery{Query: string(longQ)}
	if err := ValidateSearchQuery(q); err == nil {
		t.Error("too-long query should error")
	}

	// Limit clamping
	q = &SearchQuery{Query: "test", Limit: 500}
	if err := ValidateSearchQuery(q); err != nil {
		t.Errorf("clamping limit should not error: %v", err)
	}
	if q.Limit != 200 {
		t.Errorf("limit should be clamped to 200, got %d", q.Limit)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds  float64
		expected string
	}{
		{0, "0:00"},
		{30, "0:30"},
		{60, "1:00"},
		{90, "1:30"},
		{3600, "1:00:00"},
		{3661, "1:01:01"},
		{-5, "0:00"},
	}
	for _, tc := range tests {
		result := FormatDuration(tc.seconds)
		if result != tc.expected {
			t.Errorf("FormatDuration(%.1f) = %q, want %q", tc.seconds, result, tc.expected)
		}
	}
}

func TestNewAPIResponse(t *testing.T) {
	data := map[string]string{"hello": "world"}
	resp := NewAPIResponse(data, "test", "native", "req-1")
	if !resp.OK {
		t.Error("expected ok=true")
	}
	if resp.Backend != "test" {
		t.Errorf("expected backend=test, got %s", resp.Backend)
	}
	if resp.CapabilityMode != "native" {
		t.Errorf("expected capability_mode=native, got %s", resp.CapabilityMode)
	}
	if len(resp.Data) == 0 {
		t.Error("expected data to be non-empty")
	}
}

func TestNewAPIError(t *testing.T) {
	errResp := NewAPIError(ErrNotFound, "item not found", true, "test", "req-2")
	if errResp.OK {
		t.Error("expected ok=false")
	}
	if errResp.Error.Code != ErrNotFound {
		t.Errorf("expected code=not_found, got %s", errResp.Error.Code)
	}
	if !errResp.Error.Recoverable {
		t.Error("expected recoverable=true")
	}
}

func TestNewAPIErrorWithAction(t *testing.T) {
	errResp := NewAPIErrorWithAction(ErrPermissionDenied, "permission needed", "grant automation", "test", "req-3")
	if errResp.OK {
		t.Error("expected ok=false")
	}
	if errResp.Error.RequiredAction != "grant automation" {
		t.Errorf("expected required_action, got %s", errResp.Error.RequiredAction)
	}
}

func TestErrorCodeImplementsError(t *testing.T) {
	var e error = ErrNotFound
	if e.Error() != "not_found" {
		t.Errorf("expected 'not_found', got '%s'", e.Error())
	}
}
