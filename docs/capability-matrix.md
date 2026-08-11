# Apple Music MCP — Capability Matrix

**Generated:** 2026-08-11  
**macOS:** 26.5.2 (arm64)  
**Music.app:** 1.6.5  
**SDEF source:** `/System/Applications/Music.app` via `sdef`

## Analysis Method

The Music.app scripting dictionary was extracted with `sdef /System/Applications/Music.app` and analyzed line-by-line. Each capability was cross-referenced against:
- The SDEF class, property, and command definitions
- Runtime testing via `osascript`
- Apple's documentation for MusicKit and Apple Events

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Native — available via Apple Events (Music.sdef) |
| ⚠️ | Partial — limited or unstable |
| ❌ | Not available — not in SDEF, no public API |
| 🔑 | Requires Accessibility permission |
| 🎵 | Requires Apple Music subscription |
| 🔧 | Requires MusicKit developer token |
| 📖 | Read-only property |
| ✏️ | Read-write property |
| 🏷️ | Property on object (not a command) |

---

## 1. Playback Control

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Open/launch Music.app | ✅ | `run` command (Standard Suite) | Native |
| Play | ✅ | `play` command (hookPlay) | Can accept optional track specifier |
| Pause | ✅ | `pause` command (hookPaus) | Native |
| Toggle play/pause | ✅ | `playpause` command (hookPlPs) | Native |
| Stop | ✅ | `stop` command (hookStop) | Native |
| Next track | ✅ | `next track` command (hookNext) | Native |
| Previous track | ✅ | `previous track` command (hookPrev) | Native |
| Back track (restart) | ✅ | `back track` command (hookBack) | Repositions to beginning or previous |
| Fast forward | ✅ | `fast forward` command (hookFast) | Skips forward in playing track |
| Rewind | ✅ | `rewind` command (hookRwnd) | Skips backward in playing track |
| Resume | ✅ | `resume` command (hookResu) | Disables FF/rewind, resumes playback |
| Player state | ✅ | `player state` property (ePlS) | stopped/playing/paused/ff/rewinding |
| Player position | ✅✏️ | `player position` property | Read-write, in seconds |
| Sound volume | ✅✏️ | `sound volume` property (0-100) | Read-write |
| Mute | ✅✏️ | `mute` property | Read-write boolean |
| Shuffle enabled | ✅✏️ | `shuffle enabled` boolean | Read-write |
| Shuffle mode | ✅✏️ | `shuffle mode` (eShM) | songs/albums/groupings |
| Repeat mode | ✅✏️ | `song repeat` (eRpt) | off/one/all |

## 2. Current Track Metadata

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Current track reference | ✅📖 | `current track` property | Returns track object |
| Track name | 🏷️✏️ | `name` property on item | Read-write on track |
| Artist | 🏷️✏️ | `artist` property | Read-write |
| Album | 🏷️✏️ | `album` property | Read-write |
| Album artist | 🏷️✏️ | `album artist` property | Read-write |
| Duration | 🏷️📖 | `duration` property | Seconds (real) |
| Track number | 🏷️✏️ | `track number` | Integer |
| Track count | 🏷️✏️ | `track count` | Integer |
| Disc number | 🏷️✏️ | `disc number` | Integer |
| Disc count | 🏷️✏️ | `disc count` | Integer |
| Genre | 🏷️✏️ | `genre` property | Read-write |
| Year | 🏷️✏️ | `year` property | Integer |
| Composer | 🏷️✏️ | `composer` property | Read-write |
| BPM | 🏷️✏️ | `bpm` property | Integer |
| Bit rate | 🏷️📖 | `bit rate` property | kbps |
| Sample rate | 🏷️📖 | `sample rate` property | Hz |
| Size | 🏷️📖 | `size` property | Bytes |
| Kind | 🏷️📖 | `kind` property | Text description |
| Persistent ID | 🏷️📖 | `persistent ID` property | Hex string, stable |
| Database ID | 🏷️📖 | `database ID` property | Integer, shared across playlists |
| Played count | 🏷️✏️ | `played count` | Read-write |
| Played date | 🏷️✏️ | `played date` | Read-write |
| Release date | 🏷️📖 | `release date` | Read-only |
| Modification date | 🏷️📖 | `modification date` | Read-only |
| Date added | 🏷️📖 | `date added` | Read-only |
| Comment | 🏷️✏️ | `comment` property | Freeform text |
| Lyrics | 🏷️✏️ | `lyrics` property | Read-write |
| Description | 🏷️✏️ | `description` property | Read-write |
| Compilation | 🏷️✏️ | `compilation` boolean | Read-write |
| Gapless | 🏷️✏️ | `gapless` boolean | Read-write |
| Bookmark | 🏷️✏️ | `bookmark` property | Seconds |
| Bookmarkable | 🏷️✏️ | `bookmarkable` boolean | Read-write |
| Start time | 🏷️✏️ | `start` property | Seconds |
| Finish time | 🏷️✏️ | `finish` property | Seconds |
| EQ preset | 🏷️✏️ | `EQ` property | Read-write |
| Volume adjustment | 🏷️✏️ | `volume adjustment` | -100% to +100% |
| Enabled | 🏷️✏️ | `enabled` boolean | Checkbox state |
| Shufflable | 🏷️✏️ | `shufflable` boolean | Read-write |
| Skipped count | 🏷️✏️ | `skipped count` | Read-write |
| Skipped date | 🏷️✏️ | `skipped date` | Read-write |
| Unplayed | 🏷️✏️ | `unplayed` boolean | Read-write |
| Cloud status | 🏷️📖 | `cloud status` (eClS) | iCloud status |
| Sort fields | 🏷️✏️ | sort album/artist/name/composer/show | Override strings |

## 3. Artwork

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Artwork object | ✅📖 | `artwork` class (cArt) | Element of track/playlist |
| Artwork data (picture) | 🏷️✏️ | `data` property | PICT format via AppleScript |
| Artwork raw data | 🏷️✏️ | `raw data` property | Original format |
| Artwork format | 🏷️📖 | `format` property | Data format type |
| Artwork downloaded | 🏷️📖 | `downloaded` property | Boolean |
| Artwork description | 🏷️✏️ | `description` property | Text |
| Artwork kind | 🏷️✏️ | `kind` property | Integer |

## 4. Search

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Search library playlist | ✅📖 | `search` command (hookSrch) | Searches a given playlist |
| Search by albums | ✅ | `only` parameter: albums (eSrA) | |
| Search by artists | ✅ | `only` parameter: artists (eSrA) | |
| Search by composers | ✅ | `only` parameter: composers (eSrA) | |
| Search by names | ✅ | `only` parameter: names (eSrA) | |
| Search all fields | ✅ | `only` parameter: all (eSrA) | Default |
| Search displayed fields | ✅ | `only` parameter: displayed (eSrA) | |
| Apple Music catalog search | ❌ | Not in SDEF | Requires MusicKit API |
| Search by persistent ID | ✅ | via `some track whose persistent ID is "..."` | Slow for large libraries |

## 5. Playlists

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| List all playlists | ✅📖 | `playlists` element of application | |
| List user playlists | ✅📖 | `user playlists` of source | |
| Get playlist by name | ✅📖 | `playlist "name"` specifier | |
| Get playlist by persistent ID | ✅ | `some playlist whose persistent ID is "..."` | |
| Playlist name | 🏷️✏️ | `name` property on item | Read-write |
| Playlist description | 🏷️✏️ | `description` property | Read-write |
| Playlist duration | 🏷️📖 | `duration` property | Total seconds |
| Playlist size (bytes) | 🏷️📖 | `size` property | Total bytes |
| Playlist time (formatted) | 🏷️📖 | `time` property | MM:SS format |
| Playlist tracks | ✅📖 | `tracks` element of playlist | |
| Playlist favorited | 🏷️✏️ | `favorited` boolean | Read-write |
| Playlist disliked | 🏷️✏️ | `disliked` boolean | Read-write |
| Playlist parent | 🏷️📖 | `parent` property | Folder containing playlist |
| Playlist special kind | 🏷️📖 | `special kind` (eSpK) | none/Library/Genius/etc |
| Playlist visible | 🏷️📖 | `visible` boolean | In source list |
| Create playlist | ✅ | `make` command (Standard Suite) | `make new user playlist with properties {name:"X"}` |
| Delete playlist | ✅ | `delete` command (Standard Suite) | |
| Rename playlist | ✅ | Set `name` property | |
| Duplicate playlist | ✅ | `duplicate` command (Standard Suite) | |
| Add track to playlist | ✅ | `add` command (hookAdd) | Adds files; `duplicate` for existing tracks |
| Remove track from playlist | ✅ | `delete` command on track reference within playlist | |
| Move playlist to folder | ✅ | `move` command (Standard Suite) | |
| Folder playlists | ✅ | `folder playlist` class (cFoP) | |
| Create folder | ✅ | `make new folder playlist` | |
| Smart playlists | ✅📖 | `smart` property (read-only) | Cannot create via AppleScript? |
| Genius playlists | ✅📖 | `genius` property (read-only) | |
| Shared playlists | 🏷️✏️ | `shared` boolean | Read-write |
| Library playlist | ✅📖 | `library playlist` class (cLiP) | The main library |
| Subscription playlists | ✅📖 | `subscription playlist` class (cSuP) | Apple Music playlists |
| Export playlist | ✅📖 | `export` command (hookExpt) | Various formats |

## 6. Queue (Up Next)

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| List Up Next queue | ❌ | **Not in SDEF** | No Apple Event API for queue |
| Add to Up Next (play next) | ❌ | **Not in SDEF** | No Apple Event API |
| Add to Up Next (play later) | ❌ | **Not in SDEF** | No Apple Event API |
| Remove from queue | ❌ | **Not in SDEF** | No Apple Event API |
| Move within queue | ❌ | **Not in SDEF** | No Apple Event API |
| Clear queue | ❌ | **Not in SDEF** | No Apple Event API |
| Jump to queue item | ❌ | **Not in SDEF** | No Apple Event API |
| Autoplay toggle | ❌ | **Not in SDEF** | No Apple Event API |

**Queue status: COMPLETELY UNAVAILABLE via Apple Events.**  
The Music.sdef contains zero queue-related classes, properties, or commands.  
Fallback options:
- Accessibility automation of Music.app UI (requires Accessibility permission) 🔑
- Safari/MusicKit web automation (requires Apple Music web session) 🔑🎵
- Emulated queue via temporary playlist (labeled `emulated_queue`)

## 7. Favorites & Ratings

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Track favorited | 🏷️✏️ | `favorited` boolean on track | Read-write |
| Track disliked | 🏷️✏️ | `disliked` boolean on track | Read-write |
| Album favorited | 🏷️✏️ | `album favorited` boolean | Read-write |
| Album disliked | 🏷️✏️ | `album disliked` boolean | Read-write |
| Track rating (0-100) | 🏷️✏️ | `rating` property on track | 0-100 scale |
| Album rating (0-100) | 🏷️✏️ | `album rating` property | 0-100 scale |
| Rating kind | 🏷️📖 | `rating kind` (eRtK) | user/computed |
| List favorites | ⚠️ | Filter tracks/albums by favorited property | Works but slow on large libraries |
| Suggest less | ⚠️ | Set `disliked` to true | Partial equivalent |
| Star rating (1-5) | ⚠️ | `rating` is 0-100, convert to stars | Needs mapping layer |

## 8. Library Management

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Add file to library | ✅ | `add` command to library playlist | Adds file(s) to library |
| Remove track from library | ✅ | `delete` track | Removes from library |
| Recently added | ⚠️ | Filter by `date added`, sort descending | No native endpoint |
| Recently played | ⚠️ | Filter by `played date`, sort descending | No native endpoint |
| List all tracks | ✅📖 | `tracks` of library playlist | |
| List all albums | ⚠️ | Group tracks by album | No native album list API |
| List all artists | ⚠️ | Group tracks by artist | No native artist list API |
| List genres | ⚠️ | Group tracks by genre | No native genre list API |
| Add Apple Music catalog item to library | ❌ | **Not in SDEF** | Requires MusicKit API |
| Download cloud track | ✅ | `download` command (hookDwnl) | Downloads iCloud tracks |
| Convert track | ✅ | `convert` command (hookConv) | Convert format |
| Refresh track metadata | ✅ | `refresh` command (hookRfrs) | Update from file |

## 9. AirPlay & Audio Output

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| List AirPlay devices | ✅📖 | `AirPlay devices` element | All known devices |
| Current AirPlay devices | ✅✏️ | `current AirPlay devices` property | Selected devices |
| AirPlay enabled | ✅📖 | `AirPlay enabled` boolean | Read-only |
| Device active | 🏷️📖 | `active` boolean | Currently playing to |
| Device available | 🏷️📖 | `available` boolean | |
| Device selected | 🏷️✏️ | `selected` boolean | Read-write |
| Device sound volume | 🏷️✏️ | `sound volume` on AirPlay device | 0-100 |
| Device name | 🏷️📖 | `name` on item | |
| Device kind | 🏷️📖 | `kind` (eAPD) | computer/HomePod/Apple TV/etc |
| Device supports audio | 🏷️📖 | `supports audio` boolean | |
| Device protected | 🏷️📖 | `protected` boolean | Password-protected |

## 10. User Interface

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| Reveal item | ✅ | `reveal` command (hookRevl) | Shows and selects in Music |
| Select item | ✅ | `select` command | Selects without revealing |
| Browser window selection | ✅📖 | `selection` property | User's visible selection |
| Frontmost | 🏷️✏️ | `frontmost` boolean | Is app active |
| Full screen | 🏷️✏️ | `full screen` boolean | |
| Visuals enabled | 🏷️✏️ | `visuals enabled` boolean | |
| Open location (URL) | ✅ | `open location` command (Internet Suite) | iTunes Store / stream URLs |

## 11. Misc Application

| Capability | Status | SDEF Reference | Notes |
|------------|--------|---------------|-------|
| App version | 🏷️📖 | `version` property | |
| App name | 🏷️📖 | `name` property | "Music" |
| Current stream title | 🏷️📖 | `current stream title` | Streaming server provided |
| Current stream URL | 🏷️📖 | `current stream URL` | Streaming server provided |
| EQ enabled | 🏷️✏️ | `EQ enabled` boolean | |
| Current EQ preset | 🏷️✏️ | `current EQ preset` | |
| Current encoder | 🏷️✏️ | `current encoder` | |
| Converting status | 🏷️📖 | `converting` boolean | |
| Fixed indexing | 🏷️✏️ | `fixed indexing` boolean | |
| Print playlist | ✅ | `print` command | Various formats |

---

## Summary

| Category | Native | Partial | Unavailable |
|----------|--------|---------|-------------|
| Playback Control | 16 | 0 | 0 |
| Track Metadata | 35+ | 0 | 0 |
| Artwork | 7 | 0 | 0 |
| Search | 6 | 0 | 1 (catalog) |
| Playlists | 20+ | 0 | 0 |
| Queue (Up Next) | 0 | 0 | **8** |
| Favorites & Ratings | 8 | 2 | 0 |
| Library Management | 8 | 4 | 1 (catalog add) |
| AirPlay & Audio | 10 | 0 | 0 |
| User Interface | 6 | 0 | 0 |
| Misc | 9 | 0 | 0 |
| **TOTAL** | **125+** | **6** | **10** |

## Key Gaps Requiring Fallback Engines

1. **Up Next Queue** (8 capabilities) — requires Accessibility or Safari/MusicKit fallback
2. **Apple Music Catalog Search** — requires MusicKit API (developer token)
3. **Add catalog item to library** — requires MusicKit API
4. **Autoplay** — not scriptable, potentially via Accessibility
