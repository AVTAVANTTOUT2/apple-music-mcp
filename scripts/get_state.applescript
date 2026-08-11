-- get_state.applescript - Returns player state as key:value lines
-- Each line: KEY:VALUE (value is everything after first colon)
-- Usage: osascript get_state.applescript

on run argv
	tell application "Music"
		try
			set pState to player state as text
		on error errMsg
			return "APP_RUNNING:true" & character id 10 & "ERROR:" & errMsg
		end try
		
		try
			set vol to sound volume
		on error
			set vol to 0
		end try
		
		try
			set mutedState to mute
		on error
			set mutedState to false
		end try
		
		try
			set shuf to shuffle enabled
		on error
			set shuf to false
		end try
		
		try
			set shufM to shuffle mode as text
		on error
			set shufM to "off"
		end try
		
		try
			set rptM to song repeat as text
		on error
			set rptM to "off"
		end try
		
		try
			set pos to player position
		on error
			set pos to 0
		end try
		
		set lf to character id 10
		
		set out to "APP_RUNNING:true"
		set out to out & lf & "PLAYER_STATE:" & pState
		set out to out & lf & "VOLUME:" & vol
		set out to out & lf & "MUTED:" & mutedState
		set out to out & lf & "SHUFFLE_MODE:" & shufM
		set out to out & lf & "REPEAT_MODE:" & rptM
		set out to out & lf & "POSITION:" & pos
		
		try
			set ct to current track
			if ct is not missing value then
				set out to out & lf & "HAS_TRACK:true"
				try
					set out to out & lf & "TRACK_NAME:" & (name of ct)
				end try
				try
					set out to out & lf & "TRACK_ARTIST:" & (artist of ct)
				end try
				try
					set out to out & lf & "TRACK_ALBUM:" & (album of ct)
				end try
				try
					set out to out & lf & "TRACK_DURATION:" & (duration of ct)
				end try
				try
					set out to out & lf & "TRACK_PID:" & (persistent ID of ct)
				end try
				try
					set out to out & lf & "TRACK_DBID:" & (database ID of ct)
				end try
				try
					set out to out & lf & "TRACK_FAVORITED:" & (favorited of ct)
				end try
				try
					set out to out & lf & "TRACK_RATING:" & (rating of ct)
				end try
				try
					set out to out & lf & "TRACK_YEAR:" & (year of ct)
				end try
				try
					set out to out & lf & "TRACK_GENRE:" & (genre of ct)
				end try
				try
					set out to out & lf & "TRACK_TRACK_NUMBER:" & (track number of ct)
				end try
				try
					set out to out & lf & "TRACK_TRACK_COUNT:" & (track count of ct)
				end try
				try
					set out to out & lf & "TRACK_DISC_NUMBER:" & (disc number of ct)
				end try
				try
					set out to out & lf & "TRACK_DISC_COUNT:" & (disc count of ct)
				end try
				try
					set out to out & lf & "TRACK_COMPOSER:" & (composer of ct)
				end try
				try
					set out to out & lf & "TRACK_PLAYED_COUNT:" & (played count of ct)
				end try
				try
					set out to out & lf & "TRACK_KIND:" & (kind of ct)
				end try
				try
					set artCnt to count of artworks of ct
					if artCnt > 0 then
						set out to out & lf & "TRACK_HAS_ARTWORK:true"
					else
						set out to out & lf & "TRACK_HAS_ARTWORK:false"
					end if
				on error
					set out to out & lf & "TRACK_HAS_ARTWORK:false"
				end try
			else
				set out to out & lf & "HAS_TRACK:false"
			end if
		on error
			set out to out & lf & "HAS_TRACK:false"
		end try
		
		return out
	end tell
end run
