package applescript

import (
	"reflect"
	"testing"
)

func TestBuildOsascriptArgs_AppleScriptFile(t *testing.T) {
	t.Parallel()

	got := buildOsascriptArgs("/tmp/search.applescript", "AppleScript", []string{"Werenoi", "library", "track"})
	want := []string{"/tmp/search.applescript", "Werenoi", "library", "track"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildOsascriptArgs() = %v, want %v", got, want)
	}
}

func TestBuildOsascriptArgs_AppleScriptFileNoArgs(t *testing.T) {
	t.Parallel()

	got := buildOsascriptArgs("/tmp/get_state.applescript", "AppleScript", nil)
	want := []string{"/tmp/get_state.applescript"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildOsascriptArgs() = %v, want %v", got, want)
	}
}

func TestBuildOsascriptArgs_JavaScriptFile(t *testing.T) {
	t.Parallel()

	got := buildOsascriptArgs("/tmp/handler.js", "JavaScript", []string{"play"})
	want := []string{"-l", "JavaScript", "/tmp/handler.js", "play"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildOsascriptArgs() = %v, want %v", got, want)
	}
}
