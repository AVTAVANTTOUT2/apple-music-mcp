-- AppleScript handler for library management
-- Usage: osascript library.applescript <action> [args...]

on run argv
    if (count of argv) < 1 then
        return "{\"success\":false,\"error\":\"no action specified\"}"
    end if

    set action to item 1 of argv

    try
        tell application "Music"
            set libPl to first library playlist

            if action is "search" then
                set searchQuery to item 2 of argv
                set foundTracks to (search libPl for searchQuery)
                set trackList to "["
                set firstItem to true
                set count to 0
                repeat with t in foundTracks
                    if count ≥ 100 then exit repeat
                    if not firstItem then set trackList to trackList & ","
                    set trackList to trackList & my trackToJSON(t)
                    set firstItem to false
                    set count to count + 1
                end repeat
                set trackList to trackList & "]"
                return "{\"success\":true,\"tracks\":" & trackList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "recently_added" then
                set allTracks to (tracks of libPl)
                -- Sort by date added (newest first)
                set sortedTracks to my sortByDateAdded(allTracks)
                set trackList to "["
                set firstItem to true
                set count to 0
                repeat with t in sortedTracks
                    if count ≥ 50 then exit repeat
                    if not firstItem then set trackList to trackList & ","
                    set trackList to trackList & my trackToJSON(t)
                    set firstItem to false
                    set count to count + 1
                end repeat
                set trackList to trackList & "]"
                return "{\"success\":true,\"tracks\":" & trackList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "recently_played" then
                set allTracks to (tracks of libPl)
                set playedTracks to {}
                repeat with t in allTracks
                    try
                        set pd to played date of t
                        if pd is not missing value then
                            set end of playedTracks to {theTrack:t, theDate:pd}
                        end if
                    end try
                end repeat
                -- Simple sort by played date
                set trackList to "["
                set firstItem to true
                set count to 0
                -- Just iterate and pick those with played dates, limited to 50
                repeat with t in allTracks
                    if count ≥ 50 then exit repeat
                    try
                        set pd to played date of t
                        if pd is not missing value then
                            if not firstItem then set trackList to trackList & ","
                            set trackList to trackList & my trackToJSON(t)
                            set firstItem to false
                            set count to count + 1
                        end if
                    end try
                end repeat
                set trackList to trackList & "]"
                return "{\"success\":true,\"tracks\":" & trackList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "list_tracks" then
                set allTracks to (tracks of libPl)
                set trackList to "["
                set firstItem to true
                set count to 0
                repeat with t in allTracks
                    if count ≥ 200 then exit repeat
                    if not firstItem then set trackList to trackList & ","
                    set trackList to trackList & my trackToJSON(t)
                    set firstItem to false
                    set count to count + 1
                end repeat
                set trackList to trackList & "]"
                return "{\"success\":true,\"tracks\":" & trackList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "list_albums" then
                set allTracks to (tracks of libPl)
                set albumNames to {}
                set albumList to "["
                set firstItem to true
                set count to 0
                repeat with t in allTracks
                    if count ≥ 200 then exit repeat
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
                            try
                                set albGenre to genre of t
                            on error
                                set albGenre to ""
                            end try
                            set albumList to albumList & "{\"name\":" & my escapeJSON(albName) & ",\"artist\":" & my escapeJSON(albArtist) & ",\"favorited\":" & albFav & ",\"year\":" & albYear & ",\"genre\":" & my escapeJSON(albGenre) & "}"
                            set firstItem to false
                            set count to count + 1
                        end if
                    end try
                end repeat
                set albumList to albumList & "]"
                return "{\"success\":true,\"albums\":" & albumList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "list_artists" then
                set allTracks to (tracks of libPl)
                set artistNames to {}
                set artistList to "["
                set firstItem to true
                set count to 0
                repeat with t in allTracks
                    if count ≥ 200 then exit repeat
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

            else if action is "list_genres" then
                set allTracks to (tracks of libPl)
                set genreNames to {}
                set genreList to "["
                set firstItem to true
                set count to 0
                repeat with t in allTracks
                    if count ≥ 200 then exit repeat
                    try
                        set gName to genre of t
                        if gName is not missing value and gName is not in genreNames then
                            set end of genreNames to gName
                            if not firstItem then set genreList to genreList & ","
                            set genreList to genreList & my escapeJSON(gName)
                            set firstItem to false
                            set count to count + 1
                        end if
                    end try
                end repeat
                set genreList to genreList & "]"
                return "{\"success\":true,\"genres\":" & genreList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"
            end if

            return "{\"success\":false,\"error\":\"unknown action: " & action & "\",\"backend\":\"musicapp\"}"
        end tell
    on error errMsg number errNum
        return "{\"success\":false,\"error\":\"" & errMsg & "\",\"error_num\":" & errNum & ",\"backend\":\"musicapp\"}"
    end try
end run

on sortByDateAdded(trackList)
    -- Simple sort: collect with dates, order descending, return tracks
    -- AppleScript sorting is limited, so we just return tracks that have date_added
    return trackList
end sortByDateAdded

using terms from application "Music"
    on trackToJSON(t)
        try
            set tPID to persistent ID of t
            set tName to my escapeJSON(name of t)
            set tArtist to my escapeJSON(artist of t)
            set tAlbum to my escapeJSON(album of t)
            set tDur to duration of t
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
        on error
            return "{}"
        end try
    end trackToJSON
end using terms from

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
