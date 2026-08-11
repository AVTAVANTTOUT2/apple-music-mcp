-- AppleScript handler for searching the Music library
-- Usage: osascript search.applescript <query> [scope] [search_type]
-- Scope: library, all
-- Type: track, album, artist, playlist

on run argv
    if (count of argv) < 1 then
        return "{\"success\":false,\"error\":\"no search query specified\"}"
    end if

    set searchQuery to item 1 of argv
    set searchScope to "all"
    set searchType to "track"
    if (count of argv) >= 2 then set searchScope to item 2 of argv
    if (count of argv) >= 3 then set searchType to item 3 of argv

    try
        tell application "Music"
            set thePlaylist to first library playlist

            try
                if searchType is "track" then
                    set foundTracks to (search thePlaylist for searchQuery)
                    set trackList to "["
                    set firstItem to true
                    set count to 0
                    repeat with t in foundTracks
                        if count ≥ 50 then exit repeat
                        if not firstItem then set trackList to trackList & ","
                        set trackList to trackList & my trackToJSON(t)
                        set firstItem to false
                        set count to count + 1
                    end repeat
                    set trackList to trackList & "]"
                    return "{\"success\":true,\"tracks\":" & trackList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

                else if searchType is "album" then
                    set foundTracks to (search thePlaylist for searchQuery)
                    set albumNames to {}
                    set albumList to "["
                    set firstItem to true
                    set count to 0
                    repeat with t in foundTracks
                        if count ≥ 50 then exit repeat
                        try
                            set albName to album of t
                            if albName is not missing value and albName is not in albumNames then
                                set end of albumNames to albName
                                if not firstItem then set albumList to albumList & ","
                                set albArtist to ""
                                try
                                    set albArtist to album artist of t
                                end try
                                if albArtist is missing value then set albArtist to artist of t
                                try
                                    set albFav to album favorited of t
                                on error
                                    set albFav to false
                                end try
                                try
                                    set albYear to year of t
                                on error
                                    set albYear to 0
                                end try
                                set albumList to albumList & "{\"name\":" & my escapeJSON(albName) & ",\"artist\":" & my escapeJSON(albArtist) & ",\"favorited\":" & albFav & ",\"year\":" & albYear & "}"
                                set firstItem to false
                                set count to count + 1
                            end if
                        end try
                    end repeat
                    set albumList to albumList & "]"
                    return "{\"success\":true,\"albums\":" & albumList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

                else if searchType is "artist" then
                    set foundTracks to (search thePlaylist for searchQuery)
                    set artistNames to {}
                    set artistList to "["
                    set firstItem to true
                    set count to 0
                    repeat with t in foundTracks
                        if count ≥ 50 then exit repeat
                        try
                            set artName to artist of t
                            if artName is not missing value and artName is not in artistNames then
                                set end of artistNames to artName
                                if not firstItem then set artistList to artistList & ","
                                set artistList to artistList & "{\"name\":" & my escapeJSON(artName) & "}"
                                set firstItem to false
                                set count to count + 1
                            end if
                        end try
                    end repeat
                    set artistList to artistList & "]"
                    return "{\"success\":true,\"artists\":" & artistList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

                else if searchType is "playlist" then
                    set allPlaylists to every playlist
                    set playlistList to "["
                    set firstItem to true
                    set count to 0
                    repeat with pl in allPlaylists
                        if count ≥ 50 then exit repeat
                        try
                            set plName to name of pl
                            if plName contains searchQuery then
                                if not firstItem then set playlistList to playlistList & ","
                                set plPID to persistent ID of pl
                                try
                                    set plKind to special kind of pl as text
                                on error
                                    set plKind to "none"
                                end try
                                try
                                    set plCount to count of tracks of pl
                                on error
                                    set plCount to 0
                                end try
                                set playlistList to playlistList & "{\"name\":" & my escapeJSON(plName) & ",\"persistent_id\":\"" & plPID & "\",\"kind\":\"" & plKind & "\",\"track_count\":" & plCount & "}"
                                set firstItem to false
                                set count to count + 1
                            end if
                        end try
                    end repeat
                    set playlistList to playlistList & "]"
                    return "{\"success\":true,\"playlists\":" & playlistList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"
                end if

            on error errMsg
                return "{\"success\":false,\"error\":\"" & errMsg & "\",\"backend\":\"musicapp\"}"
            end try
        end tell
    on error errMsg number errNum
        return "{\"success\":false,\"error\":\"" & errMsg & "\",\"error_num\":" & errNum & ",\"backend\":\"musicapp\"}"
    end try
end run

on trackToJSON(t)
    try
        set tName to my escapeJSON(name of t)
    on error
        set tName to "\"\""
    end try
    try
        set tArtist to my escapeJSON(artist of t)
    on error
        set tArtist to "\"\""
    end try
    try
        set tAlbum to my escapeJSON(album of t)
    on error
        set tAlbum to "\"\""
    end try
    try
        set tPID to persistent ID of t
    on error
        set tPID to ""
    end try
    try
        set tDur to duration of t
    on error
        set tDur to 0
    end try
    try
        set tFav to favorited of t
    on error
        set tFav to false
    end try
    try
        set tYear to year of t
    on error
        set tYear to 0
    end try
    try
        set tGenre to my escapeJSON(genre of t)
    on error
        set tGenre to "\"\""
    end try
    return "{\"name\":" & tName & ",\"artist\":" & tArtist & ",\"album\":" & tAlbum & ",\"persistent_id\":\"" & tPID & "\",\"duration\":" & tDur & ",\"favorited\":" & tFav & ",\"year\":" & tYear & ",\"genre\":" & tGenre & "}"
end trackToJSON

on escapeJSON(str)
    if str is missing value then return "\"\""
    try
        set str to str as text
        set escaped to do shell script "python3 -c " & quoted form of ("import sys,json; print(json.dumps(sys.stdin.read().rstrip('\\n')))") with input str without altering line endings
        return escaped
    on error
        return "\"\""
    end try
end escapeJSON
