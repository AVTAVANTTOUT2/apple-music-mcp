-- Common JSON escaping utility for AppleScript
-- Replaces special characters for JSON string values

on escapeJSON(str)
	if str is missing value then return "\"\""
	try
		set str to str as text
	on error
		return "\"\""
	end try
	
	-- Build escaped string manually
	set resultStr to ""
	repeat with i from 1 to count of str
		set c to character i of str
		if c is "\"" then
			set resultStr to resultStr & "\\\""
		else if c is "\\" then
			set resultStr to resultStr & "\\\\"
		else if c is "/" then
			set resultStr to resultStr & "\\/"
		else if c is return then
			set resultStr to resultStr & "\\n"
		else if c is tab then
			set resultStr to resultStr & "\\t"
		else if (ASCII number of c) < 32 then
			set resultStr to resultStr & " "
		else
			set resultStr to resultStr & c
		end if
	end repeat
	return "\"" & resultStr & "\""
end escapeJSON
