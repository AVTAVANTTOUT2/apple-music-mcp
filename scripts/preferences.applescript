-- AppleScript handler for preferences: volume, shuffle, repeat, AirPlay
-- Usage: osascript preferences.applescript <action> [value]

on escapeJSON(str)
	if str is missing value then return "\"\""
	try
		set str to str as text
	on error
		return "\"\""
	end try
	set pythonCmd to "import sys,json; print(json.dumps(sys.stdin.read().rstrip(chr(10))))"
	try
		return do shell script "echo " & quoted form of str & " | python3 -c " & quoted form of pythonCmd
	on error
		return "\"" & str & "\""
	end try
end escapeJSON

on run argv
	if (count of argv) < 1 then
		return "{\"success\":false,\"error\":\"no action specified\"}"
	end if
	
	set action to item 1 of argv
	set val to ""
	if (count of argv) >= 2 then set val to item 2 of argv
	
	try
		tell application "Music"
			if action is "get" then
				set vol to sound volume
				set mut to mute
				set shuf to shuffle enabled
				set shufM to shuffle mode as text
				set rpt to song repeat as text
				set resultStr to "{\"success\":true,\"volume\":" & vol & ",\"muted\":" & mut & ",\"shuffle_enabled\":" & shuf & ",\"shuffle_mode\":\"" & shufM & "\",\"repeat_mode\":\"" & rpt & "\",\"backend\":\"musicapp\"}"
				return resultStr
			end if
			
			if action is "set_volume" then
				set newVol to val as integer
				if newVol < 0 then set newVol to 0
				if newVol > 100 then set newVol to 100
				set sound volume to newVol
				return "{\"success\":true,\"action\":\"set_volume\",\"volume\":" & newVol & ",\"backend\":\"musicapp\"}"
			end if
			
			if action is "set_shuffle" then
				if val is "off" then
					set shuffle enabled to false
				else if val is "songs" then
					set shuffle enabled to true
					set shuffle mode to songs
				else if val is "albums" then
					set shuffle enabled to true
					set shuffle mode to albums
				else
					return "{\"success\":false,\"error\":\"Invalid shuffle mode: " & val & "\",\"backend\":\"musicapp\"}"
				end if
				return "{\"success\":true,\"action\":\"set_shuffle\",\"shuffle_mode\":\"" & val & "\",\"backend\":\"musicapp\"}"
			end if
			
			if action is "set_repeat" then
				if val is "off" then
					set song repeat to off
				else if val is "one" then
					set song repeat to one
				else if val is "all" then
					set song repeat to all
				else
					return "{\"success\":false,\"error\":\"Invalid repeat mode: " & val & "\",\"backend\":\"musicapp\"}"
				end if
				return "{\"success\":true,\"action\":\"set_repeat\",\"repeat_mode\":\"" & val & "\",\"backend\":\"musicapp\"}"
			end if
			
			if action is "get_outputs" then
				set outputList to "["
				set firstItem to true
				try
					set allDevices to AirPlay devices
					repeat with dev in allDevices
						if not firstItem then set outputList to outputList & ","
						set devName to name of dev
						try
							set devKind to kind of dev as text
						on error
							set devKind to "unknown"
						end try
						try
							set devActive to active of dev
						on error
							set devActive to false
						end try
						try
							set devAvail to available of dev
						on error
							set devAvail to false
						end try
						try
							set devSel to selected of dev
						on error
							set devSel to false
						end try
						try
							set devVol to sound volume of dev
						on error
							set devVol to 0
						end try
						try
							set devProt to protected of dev
						on error
							set devProt to false
						end try
						
						set outputList to outputList & "{\"name\":" & my escapeJSON(devName) & ",\"kind\":\"" & devKind & "\",\"active\":" & devActive & ",\"available\":" & devAvail & ",\"selected\":" & devSel & ",\"protected\":" & devProt & ",\"volume\":" & devVol & "}"
						set firstItem to false
					end repeat
				end try
				set outputList to outputList & "]"
				return "{\"success\":true,\"action\":\"get_outputs\",\"outputs\":" & outputList & ",\"backend\":\"musicapp\"}"
			end if
			
			return "{\"success\":false,\"error\":\"unknown action: " & action & "\",\"backend\":\"musicapp\"}"
		end tell
	on error errMsg number errNum
		return "{\"success\":false,\"error\":\"" & errMsg & "\",\"error_num\":" & errNum & ",\"backend\":\"musicapp\"}"
	end try
end run
