function CodeBlock(el)
	local is_file
	for i, class in pairs(el.classes) do
		if class:sub(1, 5) == "file:" then
			el.classes[i] = class:sub(6)
			is_file = true
		end
	end
	if not is_file then return end
	local f = io.open(el.text, "r")
	if not f then return end
	local text = f:read"*a"
	f:close()
	return pandoc.CodeBlock(text, el.attr)
end