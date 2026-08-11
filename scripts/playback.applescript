-- AppleScript for playback control
-- Returns simple formatted text
-- Usage: osascript playback.applescript <action> [args...]

on run argv
	if (count of argv) < 1 then
		set lf to character id 10
		return "SUCCESS:false" & lf & "ERROR:no action specified"
	end if
	
	set action to item 1 of argv
	set targetID to ""
	set targetType to ""
	set seekPos to "0"
	set onceFlag to "false"
	
	if (count of argv) >= 2 then set targetID to item 2 of argv
	if (count of argv) >= 3 then set targetType to item 3 of argv
	if (count of argv) >= 4 then set seekPos to item 4 of argv
	if (count of argv) >= 5 then set onceFlag to item 5 of argv
	
	set lf to character id 10
	
	try
		tell application "Music"
			if action is "open" then
				activate
				return "SUCCESS:true" & lf & "ACTION:open" & lf & "MESSAGE:Music activated"
			
			else if action is "play" then
				if targetID is not "" then
					if targetType is "track" then
						set t to (some track whose persistent ID is targetID)
						play t
					else if targetType is "album" then
						set matchingTracks to (every track whose album is targetID)
						if (count of matchingTracks) > 0 then
							play (item 1 of matchingTracks)
						else
							error "No tracks found for album: " & targetID
						end if
					else if targetType is "artist" then
						set matchingTracks to (every track whose artist is targetID)
						if (count of matchingTracks) > 0 then
							play (item 1 of matchingTracks)
						else
							error "No tracks found for artist: " & targetID
						end if
					else if targetType is "playlist" then
						try
							set pl to (some playlist whose persistent ID is targetID)
							play pl
						on error
							try
								set pl to playlist targetID
								play pl
							on error errMsg2
								error "Playlist not found: " & targetID
							end try
						end try
					else
						try
							set t to (some track whose persistent ID is targetID)
							play t
						on error
							error "Target not found: " & targetID
						end try
					end if
				else
					if onceFlag is "true" then
						play once true
					else
						play
					end if
				end if
				return "SUCCESS:true" & lf & "ACTION:play" & lf & "MESSAGE:Playing"
			
			else if action is "pause" then
				pause
				return "SUCCESS:true" & lf & "ACTION:pause" & lf & "MESSAGE:Paused"
			
			else if action is "toggle" then
				playpause
				return "SUCCESS:true" & lf & "ACTION:toggle" & lf & "MESSAGE:Toggled"
			
			else if action is "stop" then
				stop
				return "SUCCESS:true" & lf & "ACTION:stop" & lf & "MESSAGE:Stopped"
			
			else if action is "next" then
				next track
				return "SUCCESS:true" & lf & "ACTION:next" & lf & "MESSAGE:Next track"
			
			else if action is "previous" then
				previous track
				return "SUCCESS:true" & lf & "ACTION:previous" & lf & "MESSAGE:Previous track"
			
			else if action is "restart_current" then
				back track
				return "SUCCESS:true" & lf & "ACTION:restart_current" & lf & "MESSAGE:Restarted"
			
			else if action is "seek_absolute" then
				set seekSeconds to seekPos as real
				set player position to seekSeconds
				return "SUCCESS:true" & lf & "ACTION:seek_absolute"
			
			else if action is "seek_relative" then
				set currentPos to player position
				set delta to seekPos as real
				set newPos to currentPos + delta
				if newPos < 0 then set newPos to 0
				set player position to newPos
				return "SUCCESS:true" & lf & "ACTION:seek_relative"
			
			else if action is "play_url" then
				open location targetID
				return "SUCCESS:true" & lf & "ACTION:play_url"
			
			else if action is "reveal" then
				try
					set t to (some track whose persistent ID is targetID)
					reveal t
				on error
					try
						set pl to (some playlist whose persistent ID is targetID)
						reveal pl
					on error
						error "Item not found for reveal: " & targetID
					end try
				end try
				return "SUCCESS:true" & lf & "ACTION:reveal"
			
			else
				return "SUCCESS:false" & lf & "ERROR:unknown action: " & action
			end if
		end tell
	on error errMsg number errNum
		return "SUCCESS:false" & lf & "ERROR:" & errMsg & lf & "ERROR_NUM:" & errNum
	end try
end run
