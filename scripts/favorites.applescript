-- AppleScript handler for favorites and ratings
-- Usage: osascript favorites.applescript <action> [target_id] [target_type] [rating]

on run argv
    if (count of argv) < 1 then
        return "{\"success\":false,\"error\":\"no action specified\"}"
    end if

    set action to item 1 of argv
    set targetID to ""
    set targetType to "track"
    set ratingVal to "0"

    if (count of argv) >= 2 then set targetID to item 2 of argv
    if (count of argv) >= 3 then set targetType to item 3 of argv
    if (count of argv) >= 4 then set ratingVal to item 4 of argv

    try
        tell application "Music"
            if action is "get" then
                try
                    set t to (some track whose persistent ID is targetID)
                    set isFav to favorited of t
                    try
                        set rat to rating of t
                    on error
                        set rat to 0
                    end try
                    try
                        set isDisliked to disliked of t
                    on error
                        set isDisliked to false
                    end try
                    return "{\"success\":true,\"favorited\":" & isFav & ",\"rating\":" & rat & ",\"disliked\":" & isDisliked & ",\"backend\":\"musicapp\"}"
                on error
                    return "{\"success\":false,\"error\":\"Track not found: " & targetID & "\",\"backend\":\"musicapp\"}"
                end try

            else if action is "list" then
                set allTracks to (tracks of first library playlist)
                set favList to "["
                set firstItem to true
                set count to 0
                repeat with t in allTracks
                    if count ≥ 200 then exit repeat
                    try
                        if favorited of t then
                            if not firstItem then set favList to favList & ","
                            set favList to favList & my trackToJSON(t)
                            set firstItem to false
                            set count to count + 1
                        end if
                    end try
                end repeat
                set favList to favList & "]"
                return "{\"success\":true,\"favorites\":" & favList & ",\"total\":" & count & ",\"backend\":\"musicapp\"}"

            else if action is "favorite" then
                if targetType is "track" then
                    set t to (some track whose persistent ID is targetID)
                    set favorited of t to true
                else if targetType is "album" then
                    set t to (some track whose persistent ID is targetID)
                    set album favorited of t to true
                end if
                return "{\"success\":true,\"action\":\"favorite\",\"favorited\":true,\"backend\":\"musicapp\"}"

            else if action is "unfavorite" then
                if targetType is "track" then
                    set t to (some track whose persistent ID is targetID)
                    set favorited of t to false
                else if targetType is "album" then
                    set t to (some track whose persistent ID is targetID)
                    set album favorited of t to false
                end if
                return "{\"success\":true,\"action\":\"unfavorite\",\"favorited\":false,\"backend\":\"musicapp\"}"

            else if action is "suggest_less" then
                set t to (some track whose persistent ID is targetID)
                set disliked of t to true
                return "{\"success\":true,\"action\":\"suggest_less\",\"disliked\":true,\"backend\":\"musicapp\"}"

            else if action is "set_rating" then
                set t to (some track whose persistent ID is targetID)
                set newRating to ratingVal as integer
                if newRating < 0 then set newRating to 0
                if newRating > 100 then set newRating to 100
                set rating of t to newRating
                return "{\"success\":true,\"action\":\"set_rating\",\"rating\":" & newRating & ",\"backend\":\"musicapp\"}"

            else if action is "clear_rating" then
                set t to (some track whose persistent ID is targetID)
                set rating of t to 0
                return "{\"success\":true,\"action\":\"clear_rating\",\"rating\":0,\"backend\":\"musicapp\"}"
            end if

            return "{\"success\":false,\"error\":\"unknown action: " & action & "\",\"backend\":\"musicapp\"}"
        end tell
    on error errMsg number errNum
        return "{\"success\":false,\"error\":\"" & errMsg & "\",\"error_num\":" & errNum & ",\"backend\":\"musicapp\"}"
    end try
end run

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
        try
            set tYear to year of t
        on error
            set tYear to 0
        end try
        return "{\"name\":" & tName & ",\"artist\":" & tArtist & ",\"album\":" & tAlbum & ",\"persistent_id\":\"" & tPID & "\",\"duration\":" & tDur & ",\"favorited\":" & tFav & ",\"rating\":" & tRating & ",\"year\":" & tYear & "}"
    on error
        return "{}"
    end try
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
