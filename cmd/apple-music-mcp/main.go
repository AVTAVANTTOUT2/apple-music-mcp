// Package main is the entry point for apple-music-mcp.
// Commands: serve, doctor, capabilities, version, install, uninstall, configure, test-live
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/backend/musicapp"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/config"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/domain"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/mcpserver"
	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/version"
)

func main() {
	cfg := config.LoadFromEnv()

	// Setup logging (stderr only, stdout reserved for MCP)
	logLevel := logging.LevelInfo
	if cfg.Verbose {
		logLevel = logging.LevelDebug
	}
	logger := logging.NewLogger(logLevel)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		cmdServe(cfg, logger)
	case "doctor":
		cmdDoctor(cfg, logger)
	case "capabilities":
		cmdCapabilities(cfg, logger)
	case "version":
		fmt.Println(version.Info())
	case "install":
		cmdInstall(cfg, logger)
	case "uninstall":
		cmdUninstall(cfg, logger)
	case "configure":
		cmdConfigure(cfg, logger)
	case "test-live":
		cmdTestLive(cfg, logger)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `apple-music-mcp — Native Apple Music MCP server for macOS

Usage:
  apple-music-mcp serve         Start the MCP server (stdio transport)
  apple-music-mcp doctor        Run non-destructive diagnostics
  apple-music-mcp capabilities  Show capability matrix
  apple-music-mcp version       Print version information
  apple-music-mcp install       Install MCP configuration for clients
  apple-music-mcp uninstall     Remove MCP configuration
  apple-music-mcp configure     Configure MCP client integration
  apple-music-mcp test-live     Run live tests against Music.app

Environment:
  APPLE_MUSIC_MCP_READ_ONLY=1     Disable all mutations
  APPLE_MUSIC_MCP_VERBOSE=1       Enable debug logging to stderr
  APPLE_MUSIC_MCP_TRANSPORT=stdio Transport mode (stdio only for v0.1.0)

`)
}

func cmdServe(cfg *config.Config, logger *logging.Logger) {
	logger.Info("starting apple-music-mcp server v%s", version.Short())

	backend := musicapp.NewBackend(logger)
	if cfg.ScriptsDir != "" {
		backend.SetScriptsDir(cfg.ScriptsDir)
	} else {
		// Look for scripts directory relative to the binary
		if exe, err := os.Executable(); err == nil {
			scriptsDir := exe + "/../../scripts"
			if _, err := os.Stat(scriptsDir); err == nil {
				backend.SetScriptsDir(scriptsDir)
			}
		}
		// Fallback: look in current directory and parent
		for _, dir := range []string{"scripts", "../scripts", "../../scripts"} {
			if _, err := os.Stat(dir); err == nil {
				backend.SetScriptsDir(dir)
				break
			}
		}
	}

	server := mcpserver.NewServer(cfg, logger, backend)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGTERM/SIGINT gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		logger.Info("received shutdown signal")
		cancel()
	}()

	if err := server.Serve(ctx); err != nil {
		logger.Error("server error: %v", err)
		os.Exit(1)
	}
}

func cmdDoctor(cfg *config.Config, logger *logging.Logger) {
	fmt.Println("🔍 Apple Music MCP Doctor")
	fmt.Println("========================")
	fmt.Println()

	// System info
	fmt.Printf("macOS:     checking...\n")
	fmt.Printf("Arch:      checking...\n")
	fmt.Printf("Go:        %s\n", version.Info())
	fmt.Println()

	// Music.app check
	fmt.Println("🎵 Music.app")
	musicPath := ""
	for _, p := range []string{"/System/Applications/Music.app", "/Applications/Music.app"} {
		if _, err := os.Stat(p); err == nil {
			musicPath = p
			fmt.Printf("  Path:      %s ✅\n", p)
			break
		}
	}
	if musicPath == "" {
		fmt.Println("  Path:      NOT FOUND ❌")
		fmt.Println("  → Install Apple Music from the App Store or reinstall macOS.")
		return
	}

	// Automation check
	fmt.Println()
	fmt.Println("🔐 Permissions")
	backend := musicapp.NewBackend(logger)
	ctx := context.Background()
	if backend.IsAvailable(ctx) {
		fmt.Println("  Automation: GRANTED ✅")
	} else {
		fmt.Println("  Automation: NOT GRANTED ❌")
		fmt.Println("  → Open System Settings > Privacy & Security > Automation")
		fmt.Println("  → Enable the terminal/client that runs this tool for Music.app")
	}

	// Get player state
	fmt.Println()
	fmt.Println("▶️  Player State")
	status, err := backend.GetPlayerState(ctx)
	if err != nil {
		fmt.Printf("  Error: %v ❌\n", err)
	} else {
		fmt.Printf("  State:    %s\n", status.PlayerState)
		fmt.Printf("  Volume:   %d%%\n", status.Volume)
		fmt.Printf("  Shuffle:  %s\n", status.ShuffleMode)
		fmt.Printf("  Repeat:   %s\n", status.RepeatMode)
		if status.Track != nil {
			fmt.Printf("  Track:    %s — %s\n", status.Track.Name, status.Track.Artist)
			fmt.Printf("  Album:    %s\n", status.Track.Album)
			fmt.Printf("  Duration: %s\n", domain.FormatDuration(status.Track.Duration))
			fmt.Printf("  PID:      %s\n", status.Track.PersistentID)
		} else {
			fmt.Println("  Track:    (no track loaded)")
		}
	}

	// Capabilities summary
	fmt.Println()
	fmt.Println("📋 Capabilities")
	caps, err := backend.Capabilities(ctx)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Native capabilities:     %d\n", len(caps.NativeCapabilities))
		fmt.Printf("  Unavailable capabilities: %d\n", len(caps.UnavailableCapabilities))
		fmt.Println("  Key unavailable:")
		for _, c := range caps.UnavailableCapabilities {
			fmt.Printf("    - %s\n", c)
		}
	}

	fmt.Println()
	fmt.Println("✅ Doctor check complete.")
}

func cmdCapabilities(cfg *config.Config, logger *logging.Logger) {
	backend := musicapp.NewBackend(logger)
	ctx := context.Background()
	caps, err := backend.Capabilities(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Backend: %s\n", caps.BackendName)
	fmt.Println()
	fmt.Println("Native capabilities:")
	for _, c := range caps.NativeCapabilities {
		fmt.Printf("  ✅ %s\n", c)
	}
	fmt.Println()
	fmt.Println("Unavailable:")
	for _, c := range caps.UnavailableCapabilities {
		fmt.Printf("  ❌ %s\n", c)
	}
}

func cmdInstall(cfg *config.Config, logger *logging.Logger) {
	fmt.Println("Installation is not yet automated for v0.1.0.")
	fmt.Println()
	fmt.Println("Manual configuration for Cursor:")
	fmt.Println("  Add to ~/.cursor/mcp.json:")
	fmt.Println(`  {
    "mcpServers": {
      "apple-music": {
        "command": "/path/to/apple-music-mcp",
        "args": ["serve"]
      }
    }
  }`)
	fmt.Println()
	fmt.Println("Manual configuration for Claude Desktop:")
	fmt.Println(`  Add to ~/Library/Application Support/Claude/claude_desktop_config.json:
  {
    "mcpServers": {
      "apple-music": {
        "command": "/path/to/apple-music-mcp",
        "args": ["serve"]
      }
    }
  }`)
}

func cmdUninstall(cfg *config.Config, logger *logging.Logger) {
	fmt.Println("To uninstall, remove the apple-music-mcp binary and delete the MCP configuration entry from your client's config file.")
	fmt.Println("No files are installed outside the binary location.")
}

func cmdConfigure(cfg *config.Config, logger *logging.Logger) {
	client := "cursor"
	if len(os.Args) > 2 {
		client = os.Args[2]
	}
	fmt.Printf("Configuration for %s:\n", client)
	fmt.Println()
	exe, _ := os.Executable()
	fmt.Printf(`{
  "mcpServers": {
    "apple-music": {
      "command": "%s",
      "args": ["serve"]
    }
  }
}`, exe)
	fmt.Println()
}

func cmdTestLive(cfg *config.Config, logger *logging.Logger) {
	if !cfg.AllowLiveTests {
		fmt.Println("Live tests are disabled. Set APPLE_MUSIC_MCP_LIVE_TESTS=1 to enable.")
		fmt.Println("WARNING: Live tests will modify your Music.app state temporarily and attempt to restore it.")
		return
	}
	fmt.Println("Live tests not yet implemented for v0.1.0.")
	fmt.Println("This command will run real tests against your Music.app, capturing and restoring state.")
}
