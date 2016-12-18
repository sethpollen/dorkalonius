local export = {}

-- From https://github.com/wikimedia/mediawiki-extensions-Scribunto/blob/master/engines/LuaCommon/lualib/mw.text.lua.
local function trim(s)
  s = string.gsub(s, '^[\t\r\n\f ]*(.-)[\t\r\n\f ]*$', '%1')
  return s
end

-- A helper function to escape magic characters in a string
-- Magic characters: ^$()%.[]*+-?
local function plain(text)
	return string.gsub(text, "([%^%$%(%)%%%.%[%]%*%+%-%?])", "%%%1")
end

-- A helper function that removes empty numeric indexes in a table,
-- so that the values are tightly packed like in a normal Lua table.
function export.remove_holes(list)
	local new_list = {}
	
	for i = 1, list.maxindex do
		table.insert(new_list, list[i])
	end
	
	return new_list
end

function export.process(args, params, return_unknown)
	local args_new = {}
	
	-- Process parameters for specific properties
	local required = {}
	local patterns = {}
	local list_from_index = nil
	
	for name, param in pairs(params) do
		if param.required then
			required[name] = true
		end
		
		if param.list then
			if param.default ~= nil then
				args_new[type(name) == "string" and string.gsub(name, "=", "") or name] = {param.default, maxindex = 1}
			else
				args_new[type(name) == "string" and string.gsub(name, "=", "") or name] = {maxindex = 0}
			end
			
			if type(param.list) == "string" then
				-- If the list property is a string, then it represents the name
				-- to be used as the prefix for list items. This is for use with lists
				-- where the first item is a numbered parameter and the
				-- subsequent ones are named, such as 1, pl2, pl3.
				if string.match(param.list, "=") then
					patterns["^" .. string.gsub(plain(param.list), "=", "(%%d+)") .. "$"] = name
				else
					patterns["^" .. plain(param.list) .. "(%d+)$"] = name
				end
			elseif type(name) == "number" then
				-- If the name is a number, then all indexed parameters from
				-- this number onwards go in the list.
				list_from_index = name
			else
				if string.match(name, "=") then
					patterns["^" .. string.gsub(plain(name), "=", "(%%d+)") .. "$"] = string.gsub(name, "=", "")
				else
					patterns["^" .. plain(name) .. "(%d+)$"] = string.gsub(name, "=", "")
				end
			end
			
			if string.match(name, "=") then
				params[string.gsub(name, "=", "")] = params[name]
				params[name] = nil
			end
		elseif param.default ~= nil then
			args_new[name] = param.default
		end
	end
	
	-- Process the arguments
	local args_unknown = {}
	
	for name, val in pairs(args) do
		local index = nil
		
		if type(name) == "number" then
			if list_from_index ~= nil and name >= list_from_index then
				index = name - list_from_index + 1
				name = list_from_index
			end
		else
			-- Does this argument name match a pattern?
			for pattern, pname in pairs(patterns) do
				index = string.match(name, pattern)
				
				-- It matches, so store the parameter name and the
				-- numeric index extracted from the argument name.
				if index then
					index = tonumber(index)
					name = pname
					break
				end
			end
		end
		
		-- If no index was found, use 1 as the default index.
		-- This makes list parameters like g, g2, g3 put g at index 1.
		index = index or 1
		
		local param = params[name]
		
		-- If the argument is not in the list of parameters, trigger an error.
		-- return_unknown suppresses the error, and stores it in a separate list instead.
		if not param then
			if return_unknown then
				args_unknown[name] = val
			else
				error("The parameter \"" .. name .. "\" is not used by this template.")
			end
		else
			-- Remove leading and trailing whitespace
			val = trim(val)
			
			-- Empty string is equivalent to nil unless allow_empty is true.
			if val == "" and not param.allow_empty then
				val = nil
			end
			
			-- Convert to proper type if necessary.
			if param.type == "boolean" then
				val = not (not val or val == "" or val == "0" or val == "no" or val == "n" or val == "false")
			elseif param.type == "number" and val ~= nil then
				val = tonumber(val)
			end
			
			-- Can't use "if val" alone, because val may be a boolean false.
			if val ~= nil then
				-- Mark it as no longer required, as it is present.
				required[name] = nil
				
				-- Store the argument value.
				if param.list then
					-- If the parameter is an alias of another, store it as the original,
					-- but avoid overwriting it; the original takes precedence.
					if not param.alias_of then
						args_new[name][index] = val
						
						-- Store the highest index we find.
						args_new[name].maxindex = math.max(index, args_new[name].maxindex)
					elseif args[param.alias_of] == nil then
						if params[param.alias_of] and params[param.alias_of].list then
							args_new[param.alias_of][index] = val
							
							-- Store the highest index we find.
							args_new[param.alias_of].maxindex = math.max(1, args_new[param.alias_of].maxindex)
						else
							args_new[param.alias_of] = val
						end
					end
				else
					-- If the parameter is an alias of another, store it as the original,
					-- but avoid overwriting it; the original takes precedence.
					if not param.alias_of then
						args_new[name] = val
					elseif args[param.alias_of] == nil then
						if params[param.alias_of] and params[param.alias_of].list then
							args_new[param.alias_of][1] = val
							
							-- Store the highest index we find.
							args_new[param.alias_of].maxindex = math.max(1, args_new[param.alias_of].maxindex)
						else
							args_new[param.alias_of] = val
						end
					end
				end
			end
		end
	end
	
	-- The required table should now be empty.
	-- If any entry remains, trigger an error, unless we're in the template namespace.
	--if mw.title.getCurrentTitle().nsText ~= "Template" then
	--	for name, param in pairs(required) do
	--		error("The parameter \"" .. name .. "\" is required.")
	--	end
	--end
	
	-- Remove holes in any list parameters if needed.
	for name, val in pairs(args_new) do
		if type(val) == "table" and not params[name].allow_holes then
			args_new[name] = export.remove_holes(val)
		end
	end
	
	if return_unknown then
		return args_new, args_unknown
	else
		return args_new
	end
end

return export
