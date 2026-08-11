-- AppleScript handler for playlist management
-- Usage: osascript playlists.applescript <action> [args...]
-- Actions: list, get, create, rename, delete, add_tracks, remove_tracks, copy, list_folders, create_folder, move_to_folder

on run argv
    if (count of argv) < 1 then
        return "{\"success\":false,\"error\":\"no action specified\"}"
    end if

    set action to item 1 of argv

    try
        tell application "Music"
            if action is "list" then
                set allPlaylists to every playlist
                set playlistList to "["
                set firstItem to true
                set count to 0
                repeat with pl in allPlaylists
                    if not firstItem then set playlistList to playlistList & ","
                    set plName to name of pl
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
                    try
                        set plDur to duration of pl
                    on error
                        set plDur to 0
                    end try
                    try
                        set plFav to favorited of pl
                    on error
                        set plFav to false
                    end try
                    try
                        set plSmart to smart of pl
                    on error
                        set plSmart to false
                    end try
                    try
                        set plVis to visible of pl
                    on error
                        set plVis to true
                    end try
                    set playlistList to playlistList & "{\"name\":" & my escapeJSON(plName) & ",\"persistent_id\":\"" & plPID & "\",\"kind\":\"" & plKind & "\",\"track_count\":" & plCount & ",\"duration\":" & plDur & ",\"favorited\":" & plFav & ",\"smart\":" & plSmart & ",\"visible\":" & plVis & "}"
                    set firstItem to false
                end repeat
                set playlistList to playlistList & "]"
                return "{\"success\":true,\"playlists\":" & playlistList & ",\"backend\":\"musicapp\"}"

            else if action is "get" then
                set plID to item 2 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                set plName to name of pl
                set plPID to persistent ID of pl
                try
                    set plKind to special kind of pl as text
                on error
                    set plKind to "none"
                end try
                set plCount to count of tracks of pl
                set trackList to "["
                set firstItem to true
                set idx to 0
                repeat with t in (tracks of pl)
                    if idx ≥ 200 then exit repeat
                    if not firstItem then set trackList to trackList & ","
                    set trackList to trackList & my trackToJSON(t)
                    set firstItem to false
                    set idx to idx + 1
                end repeat
                set trackList to trackList & "]"
                return "{\"success\":true,\"playlist\":{\"name\":" & my escapeJSON(plName) & ",\"persistent_id\":\"" & plPID & "\",\"kind\":\"" & plKind & "\",\"track_count\":" & plCount & "},\"tracks\":" & trackList & ",\"backend\":\"musicapp\"}"

            else if action is "create" then
                set plName to item 2 of argv
                set newPl to make new user playlist with properties {name:plName}
                return "{\"success\":true,\"playlist\":{\"name\":" & my escapeJSON(plName) & ",\"persistent_id\":\"" & persistent ID of newPl & "\",\"kind\":\"none\",\"track_count\":0},\"backend\":\"musicapp\"}"

            else if action is "rename" then
                set plID to item 2 of argv
                set newName to item 3 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                set name of pl to newName
                return "{\"success\":true,\"action\":\"rename\",\"message\":\"Renamed to " & newName & "\",\"backend\":\"musicapp\"}"

            else if action is "delete" then
                set plID to item 2 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                delete pl
                return "{\"success\":true,\"action\":\"delete\",\"message\":\"Playlist deleted\",\"backend\":\"musicapp\"}"

            else if action is "add_tracks" then
                set plID to item 2 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                set addedCount to 0
                repeat with i from 3 to (count of argv)
                    set trackID to item i of argv
                    try
                        set t to (some track whose persistent ID is trackID)
                        duplicate t to pl
                        set addedCount to addedCount + 1
                    end try
                end repeat
                return "{\"success\":true,\"action\":\"add_tracks\",\"added\":" & addedCount & ",\"backend\":\"musicapp\"}"

            else if action is "remove_tracks" then
                set plID to item 2 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                set removedCount to 0
                repeat with i from 3 to (count of argv)
                    set trackID to item i of argv
                    repeat with t in (tracks of pl)
                        if persistent ID of t is trackID then
                            delete t
                            set removedCount to removedCount + 1
                            exit repeat
                        end if
                    end repeat
                end repeat
                return "{\"success\":true,\"action\":\"remove_tracks\",\"removed\":" & removedCount & ",\"backend\":\"musicapp\"}"

            else if action is "copy" then
                set plID to item 2 of argv
                set newName to item 3 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                set newPl to duplicate pl
                set name of newPl to newName
                return "{\"success\":true,\"playlist\":{\"name\":" & my escapeJSON(newName) & ",\"persistent_id\":\"" & persistent ID of newPl & "\"},\"backend\":\"musicapp\"}"

            else if action is "list_folders" then
                set folderPls to every folder playlist
                set folderList to "["
                set firstItem to true
                repeat with fp in folderPls
                    if not firstItem then set folderList to folderList & ","
                    set fpName to name of fp
                    set fpPID to persistent ID of fp
                    set folderList to folderList & "{\"name\":" & my escapeJSON(fpName) & ",\"persistent_id\":\"" & fpPID & "\",\"kind\":\"folder\"}"
                    set firstItem to false
                end repeat
                set folderList to folderList & "]"
                return "{\"success\":true,\"folders\":" & folderList & ",\"backend\":\"musicapp\"}"

            else if action is "create_folder" then
                set folderName to item 2 of argv
                set newFolder to make new folder playlist with properties {name:folderName}
                return "{\"success\":true,\"folder\":{\"name\":" & my escapeJSON(folderName) & ",\"persistent_id\":\"" & persistent ID of newFolder & "\"},\"backend\":\"musicapp\"}"

            else if action is "move_to_folder" then
                set plID to item 2 of argv
                set folderID to item 3 of argv
                try
                    set pl to (some playlist whose persistent ID is plID)
                on error
                    try
                        set pl to playlist plID
                    on error
                        return "{\"success\":false,\"error\":\"Playlist not found: " & plID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                try
                    set targetFolder to (some folder playlist whose persistent ID is folderID)
                on error
                    try
                        set targetFolder to folder playlist folderID
                    on error
                        return "{\"success\":false,\"error\":\"Folder not found: " & folderID & "\",\"backend\":\"musicapp\"}"
                    end try
                end try
                move pl to targetFolder
                return "{\"success\":true,\"action\":\"move_to_folder\",\"backend\":\"musicapp\"}"
            end if

            return "{\"success\":false,\"error\":\"unknown action: " & action & "\",\"backend\":\"musicapp\"}"
        end tell
    on error errMsg number errNum
        return "{\"success\":false,\"error\":\"" & errMsg & "\",\"error_num\":" & errNum & ",\"backend\":\"musicapp\"}"
    end try
end run

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
                set tRating to rating of t
            on error
                set tRating to 0
            end try
            return "{\"name\":" & tName & ",\"artist\":" & tArtist & ",\"album\":" & tAlbum & ",\"persistent_id\":\"" & tPID & "\",\"duration\":" & tDur & ",\"favorited\":" & tFav & ",\"rating\":" & tRating & "}"
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
