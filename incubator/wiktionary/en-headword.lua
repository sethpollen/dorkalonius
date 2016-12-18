-- Invoking this script:
--   lua en-headword.lua <pos> <title> <args>
-- Where:
--   <pos> is "nouns", "verbs", etc.
--   <title> is the title of the Wiktionary page (usually the base form of the
--           word).
--   <args> are the arguments passed to the invoking template by the page.

local frame = {}
frame.parent = {}
frame.parent.args = {}
PAGENAME = ''
next_parent_args_position = 1

for k, v in pairs(arg) do
  if k <= 0 then
    -- Skip these args.
  elseif k == 1 then
    frame.args = {v}
  elseif k == 2 then
    PAGENAME = v
  else
    equals_sign = string.find(v, '=')
    if equals_sign == nil then
      -- This is a positional argument.
      frame.parent.args[next_parent_args_position] = v
      next_parent_args_position = next_parent_args_position + 1
    else
      -- This is a named argument.
      arg_name = string.sub(v, 1, equals_sign - 1)
      arg_value = string.sub(v, equals_sign + 1)      
      frame.parent.args[arg_name] = arg_value
    end
  end
end

local export = {}
local pos_functions = {}

-- The main entry point.
-- This is the only function that can be invoked from a template.
function export.show(frame)
	local poscat = frame.args[1] or error("Part of speech has not been specified. Please pass parameter 1 to the module invocation.")
	
	local params = {
		["head"] = {list = true, default = ""},
		["suff"] = {type = "boolean"},
	}
	
	if pos_functions[poscat] then
		for key, val in pairs(pos_functions[poscat].params) do
			params[key] = val
		end
	end
	
	local args = require("parameters").process(frame.parent.args, params)
	local data = {heads = args["head"], inflections = {}, categories = {}}
	
	if args["suff"] then
		table.insert(data.categories, "English suffixes")
		
		if poscat == "adjectives" then
			table.insert(data.categories, "English adjective-forming suffixes")
		elseif poscat == "adverbs" then
			table.insert(data.categories, "English adverb-forming suffixes")
		elseif poscat == "nouns" then
			table.insert(data.categories, "English noun-forming suffixes")
		elseif poscat == "verbs" then
			table.insert(data.categories, "English verb-forming suffixes")
		else
			error("No category exists for suffixes forming " .. poscat .. ".")
		end
	else
		table.insert(data.categories, "English " .. poscat)
	end
	
	if pos_functions[poscat] then
		pos_functions[poscat].func(args, data)
	end
	
	-- return require("Module:headword").full_headword(nil, data.heads, nil, nil, data.inflections, data.categories, nil)
  return data
end

-- This function does the common work between adjectives and adverbs
function make_comparatives(params, data)
	local comp_parts = {label = "[[Appendix:Glossary#comparable|comparative]]", accel = "comparative-form-of"}
	local sup_parts = {label = "[[Appendix:Glossary#comparable|superlative]]", accel = "superlative-form-of"}
	
	if #params == 0 then
		table.insert(params, {"more"})
	end
	
	-- To form the stem, replace -(e)y with -i and remove a final -e.
	local stem = PAGENAME:gsub("([^aeiou])e?y$", "%1i"):gsub("e$", "")
	
	-- Go over each parameter given and create a comparative and superlative form
	for i, val in ipairs(params) do
		local comp = val[1]
		local sup = val[2]
		
		if comp == "more" and PAGENAME ~= "many" and PAGENAME ~= "much" then
			table.insert(comp_parts, "[[more]] " .. PAGENAME)
			table.insert(sup_parts, "[[most]] " .. PAGENAME)
		elseif comp == "further" and PAGENAME ~= "far" then
			table.insert(comp_parts, "[[further]] " .. PAGENAME)
			table.insert(sup_parts, "[[furthest]] " .. PAGENAME)
		elseif comp == "er" then
			table.insert(comp_parts, stem .. "er")
			table.insert(sup_parts, stem .. "est")
		elseif comp == "-" or sup == "-" then
			-- Allowing '-' makes it more flexible to not have some forms
			if comp ~= "-" then
				table.insert(comp_parts, comp)
			end
			if sup ~= "-" then
				table.insert(sup_parts, sup)
			end
		else
			-- If the full comparative was given, but no superlative, then
			-- create it by replacing the ending -er with -est.
			if not sup then
				if comp:find("er$") then
					sup = comp:gsub("er$", "est")
				else
					error("The superlative of \"" .. comp .. "\" cannot be generated automatically. Please provide it with the \"sup" .. (i == 1 and "" or i) .. "=\" parameter.")
				end
			end
			
			table.insert(comp_parts, comp)
			table.insert(sup_parts, sup)
		end
	end
	
	table.insert(data.inflections, comp_parts)
	table.insert(data.inflections, sup_parts)
end

pos_functions["adjectives"] = {
	params = {
		[1] = {list = true, allow_holes = true},
		["sup"] = {list = true, allow_holes = true},
		},
	func = function(args, data)
		local shift = 0
		local is_not_comparable = 0
		local is_comparative_only = 0
		
		-- If the first parameter is ?, then don't show anything, just return.
		if args[1][1] == "?" then
			return
		-- If the first parameter is -, then move all parameters up one position.
		elseif args[1][1] == "-" then
			shift = 1
			is_not_comparable = 1
		-- If the only argument is +, then remember this and clear parameters
		elseif args[1][1] == "+" and args[1].maxindex == 1 then
			shift = 1
			is_comparative_only = 1
		end
		
		-- Gather all the comparative and superlative parameters.
		local params = {}
		
		for i = 1, args[1].maxindex - shift do
			local comp = args[1][i + shift]
			local sup = args["sup"][i]
			
			if comp or sup then
				table.insert(params, {comp, sup})
			end
		end
		
		if shift == 1 then
			-- If the first parameter is "-" but there are no parameters,
			-- then show "not comparable" only and return.
			-- If there are parameters, then show "not generally comparable"
			-- before the forms.
			if #params == 0 then
				if is_not_comparable == 1 then
					table.insert(data.inflections, {label = "not [[Appendix:Glossary#comparable|comparable]]"})
					table.insert(data.categories, "English uncomparable adjectives")
					return
				end
				if is_comparative_only == 1 then
					table.insert(data.inflections, {label = "[[Appendix:Glossary#comparable|comparative]] form only"})
					table.insert(data.categories, "English comparative-only adjectives")
					return
				end
			else
				table.insert(data.inflections, {label = "not generally [[Appendix:Glossary#comparable|comparable]]"})
			end
		end
		
		-- Process the parameters
		make_comparatives(params, data)
	end
}

pos_functions["adverbs"] = {
	params = {
		[1] = {list = true, allow_holes = true},
		["sup"] = {list = true, allow_holes = true},
		},
	func = function(args, data)
		local shift = 0
		
		-- If the first parameter is ?, then don't show anything, just return.
		if args[1][1] == "?" then
			return
		-- If the first parameter is -, then move all parameters up one position.
		elseif args[1][1] == "-" then
			shift = 1
		end
		
		-- Gather all the comparative and superlative parameters.
		local params = {}
		
		for i = 1, args[1].maxindex - shift do
			local comp = args[1][i + shift]
			local sup = args["sup"][i]
			
			if comp or sup then
				table.insert(params, {comp, sup})
			end
		end
		
		if shift == 1 then
			-- If the first parameter is "-" but there are no parameters,
			-- then show "not comparable" only and return. If there are parameters,
			-- then show "not generally comparable" before the forms.
			if #params == 0 then
				table.insert(data.inflections, {label = "not [[Appendix:Glossary#comparable|comparable]]"})
				table.insert(data.categories, "English uncomparable adverbs")
				return
			else
				table.insert(data.inflections, {label = "not generally [[Appendix:Glossary#comparable|comparable]]"})
			end
		end
		
		-- Process the parameters
		make_comparatives(params, data)
	end
}

pos_functions["nouns"] = {
	params = {
		[1] = {list = true, allow_holes = true},
		
		["plqual"] = {},
		["pl2qual"] = {},
		["pl3qual"] = {},
		["pl4qual"] = {},
		["pl5qual"] = {},
		},
	func = function(args, data)
		-- Gather all the plural parameters from the numbered parameters.
		local plurals = {}
		
		for i = 1, args[1].maxindex do
			local pl = args[1][i]
			
			if pl then
				local qual = args["pl" .. (i == 1 and "" or i) .. "qual"]
				
				if qual then
					table.insert(plurals, {term = pl, qualifiers = {qual}})
				else
					table.insert(plurals, pl)
				end
			end
		end
		
		-- Decide what to do next...
		local mode = nil
		
		if plurals[1] == "?" or plurals[1] == "!" or plurals[1] == "-" or plurals[1] == "~" then
			mode = plurals[1]
			table.remove(plurals, 1)  -- Remove the mode parameter
		end
		
		-- Plural is unknown
		if mode == "?" then
			table.insert(data.categories, "English nouns with unknown or uncertain plurals")
			return
		-- Plural is not attested
		elseif mode == "!" then
			table.insert(data.inflections, {label = "plural not attested"})
			table.insert(data.categories, "English nouns with unattested plurals")
			return
		-- Uncountable noun; may occasionally have a plural
		elseif mode == "-" then
			table.insert(data.categories, "English uncountable nouns")
			
			-- If plural forms were given explicitly, then show "usually"
			if #plurals > 0 then
				table.insert(data.inflections, {label = "usually [[Appendix:Glossary#uncountable|uncountable]]"})
				table.insert(data.categories, "English countable nouns")
			else
				table.insert(data.inflections, {label = "[[Appendix:Glossary#uncountable|uncountable]]"})
			end
		-- Mixed countable/uncountable noun, always has a plural
		elseif mode == "~" then
			table.insert(data.inflections, {label = "[[Appendix:Glossary#countable|countable]] and [[Appendix:Glossary#uncountable|uncountable]]"})
			table.insert(data.categories, "English uncountable nouns")
			table.insert(data.categories, "English countable nouns")
			
			-- If no plural was given, add a default one now
			if #plurals == 0 then
				plurals = {"s"}
			end
		-- The default, always has a plural
		else
			table.insert(data.categories, "English countable nouns")
			
			-- If no plural was given, add a default one now
			if #plurals == 0 then
				plurals = {"s"}
			end
			if plural and not mw.title.new(plural).exists then
				table.insert(categories, "English nouns with missing plurals")
			end
		end
		
		-- If there are no plurals to show, return now
		if #plurals == 0 then
			return
		end
		
		-- There are plural forms to show, so show them
		local pl_parts = {label = "plural", accel = "plural-form-of"}
		
		local stem = PAGENAME
		
		for i, pl in ipairs(plurals) do
			if pl == "s" then
				table.insert(pl_parts, stem .. "s")
			elseif pl == "es" then
				table.insert(pl_parts, stem .. "es")
			else
				table.insert(pl_parts, pl)
			end
		end
		
		table.insert(data.inflections, pl_parts)
	end
}

pos_functions["proper nouns"] = {
	params = {
		[1] = {list = true},
		},
	func = function(args, data)
		local plurals = args[1]
		
		-- Decide what to do next...
		local mode = nil
		
		if plurals[1] == "?" or plurals[1] == "!" or plurals[1] == "-" or plurals[1] == "~" then
			mode = plurals[1]
			table.remove(plurals, 1)  -- Remove the mode parameter
		end
		
		-- Plural is unknown
		if mode == "?" then
			table.insert(data.categories, "English proper nouns with unknown or uncertain plurals")
			return
		-- Plural is not attested
		elseif mode == "!" then
			table.insert(data.inflections, {label = "plural not attested"})
			table.insert(data.categories, "English proper nouns with unattested plurals")
			return
		-- Uncountable noun; may occasionally have a plural
		elseif mode == "-" then
			-- If plural forms were given explicitly, then show "usually"
			if #plurals > 0 then
				table.insert(data.inflections, {label = "usually [[Appendix:Glossary#uncountable|uncountable]]"})
				table.insert(data.categories, "English countable proper nouns")
			else
				table.insert(data.inflections, {label = "[[Appendix:Glossary#uncountable|uncountable]]"})
			end
		-- Mixed countable/uncountable noun, always has a plural
		elseif mode == "~" then
			table.insert(data.inflections, {label = "[[Appendix:Glossary#countable|countable]] and [[Appendix:Glossary#uncountable|uncountable]]"})
			table.insert(data.categories, "English countable proper nouns")
			
			-- If no plural was given, add a default one now
			if #plurals == 0 then
				plurals = {"s"}
			end
		elseif #plurals > 0 then
			table.insert(data.categories, "English countable proper nouns")
		end
		
		-- If there are no plurals to show, return now
		if #plurals == 0 then
			return
		end
		
		-- There are plural forms to show, so show them
		local pl_parts = {label = "plural", accel = "proper-noun-plural-form-of"}
		
		local stem = PAGENAME
		
		for i, pl in ipairs(plurals) do
			if pl == "s" then
				table.insert(pl_parts, stem .. "s")
			elseif pl == "es" then
				table.insert(pl_parts, stem .. "es")
			else
				table.insert(pl_parts, pl)
			end
	
		end
		
		table.insert(data.inflections, pl_parts)
	end
}

pos_functions["verbs"] = {
	params = {
		[1] = {list = "pres_3sg", allow_holes = true},
		
		["pres_3sg_qual"] = {},
		["pres_3sg2_qual"] = {},
		["pres_3sg3_qual"] = {},
		["pres_3sg4_qual"] = {},
		["pres_3sg5_qual"] = {},
		
		[2] = {list = "pres_ptc", allow_holes = true},
		
		["pres_ptc_qual"] = {},
		["pres_ptc2_qual"] = {},
		["pres_ptc3_qual"] = {},
		["pres_ptc4_qual"] = {},
		["pres_ptc5_qual"] = {},
		
		[3] = {list = "past", allow_holes = true},
		
		["past_qual"] = {},
		["past2_qual"] = {},
		["past3_qual"] = {},
		["past4_qual"] = {},
		["past5_qual"] = {},
		
		[4] = {list = "past_ptc", allow_holes = true},
		
		["past_ptc_qual"] = {},
		["past_ptc2_qual"] = {},
		["past_ptc3_qual"] = {},
		["past_ptc4_qual"] = {},
		["past_ptc5_qual"] = {},
		},
	func = function(args, data)
		-- Get parameters
		local par1 = args[1][1]
		local par2 = args[2][1]
		local par3 = args[3][1]
		local par4 = args[4][1]
		
		local pres_3sg_forms = {label = "third-person singular simple present", accel = "third-person-singular-form-of"}
		local pres_ptc_forms = {label = "present participle", accel = "present-participle-form-of"}
		local past_forms = {label = "simple past", accel = "simple-past-form-of"}
		local pres_3sg_form = par1 or PAGENAME .. "s"
		local pres_ptc_form = par2 or PAGENAME .. "ing"
		local past_form = par3 or PAGENAME .. "ed"
		
		if par1 and not par2 and not par3 then
			-- This is the "new" format, which uses only the first parameter.
			if par1 == "es" then
				pres_3sg_form = PAGENAME .. "es"
				pres_ptc_form = PAGENAME .. "ing"
				past_form = PAGENAME .. "ed"
			elseif par1 == "ies" then
				if not string.find(PAGENAME, "y$") then
					error("The first parameter is \"ies\" but the verb does not end in -y.")
				end
				
				local stem = string.gsub(PAGENAME, "y$", "")
				pres_3sg_form = stem .. "ies"
				pres_ptc_form = stem .. "ying"
				past_form = stem .. "ied"
			elseif par1 == "d" then
				pres_3sg_form = PAGENAME .. "s"
				pres_ptc_form = PAGENAME .. "ing"
				past_form = PAGENAME .. "d"
			else
				pres_3sg_form = PAGENAME .. "s"
				pres_ptc_form = par1 .. "ing"
				past_form = par1 .. "ed"
			end
		else
			-- This is the "legacy" format, using the second and third parameters as well.
			-- It is included here for backwards compatibility and to ease the transition.
			if par3 then
				if par3 == "es" then
					pres_3sg_form = par1 .. par2 .. "es"
					pres_ptc_form = par1 .. par2 .. "ing"
					past_form = par1 .. par2 .. "ed"
				elseif par3 == "ing" then
					pres_3sg_form = PAGENAME .. "s"
					pres_ptc_form = par1 .. par2 .. "ing"
					
					if par2 == "y" then
						past_form = PAGENAME .. "d"
					else
						past_form = par1 .. par2 .. "ed"
					end
				elseif par3 == "ed" then
					
					if par2 == "i" then
						pres_3sg_form = par1 .. par2 .. "es"
						pres_ptc_form = PAGENAME .. "ing"
					else
						pres_3sg_form = PAGENAME .. "s"
						pres_ptc_form = par1 .. par2 .. "ing"
					end
					
					past_form = par1 .. par2 .. "ed"
				elseif par3 == "d" then
					pres_3sg_form = PAGENAME .. "s"
					pres_ptc_form = par1 .. par2 .. "ing"
					past_form = par1 .. par2 .. "d"
				else
				end
			else
				if par2 == "es" then
					pres_3sg_form = par1 .. "es"
					pres_ptc_form = par1 .. "ing"
					past_form = par1 .. "ed"
				elseif par2 == "ies" then
					
					if par1 .. "y" ~= PAGENAME then
					end
					
					pres_3sg_form = par1 .. "ies"
					pres_ptc_form = par1 .. "ying"
					past_form = par1 .. "ied"
				elseif par2 == "ing" then
					pres_3sg_form = PAGENAME .. "s"
					pres_ptc_form = par1 .. "ing"
					past_form = par1 .. "ed"
				elseif par2 == "ed" then
					pres_3sg_form = PAGENAME .. "s"
					pres_ptc_form = par1 .. "ing"
					past_form = par1 .. "ed"
				elseif par2 == "d" then
					
					if par1 ~= PAGENAME then
					end
					
					pres_3sg_form = PAGENAME .. "s"
					pres_ptc_form = par1 .. "ing"
					past_form = par1 .. "d"
				else
				end
			end
		end
		
		local pres_3sg_qual = args["pres_3sg_qual"]
		local pres_ptc_qual = args["pres_ptc_qual"]
		local past_qual = args["past_qual"]
		
		table.insert(pres_ptc_forms, {term = pres_ptc_form, qualifiers = {pres_ptc_qual}})
		table.insert(pres_3sg_forms, {term = pres_3sg_form, qualifiers = {pres_3sg_qual}})
		table.insert(past_forms, {term = past_form, qualifiers = {past_qual}})
		
		-- Present 3rd singular
		for i = 2, args[1].maxindex do
			local form = args[1][i]
			local qual = args["pres_3sg" .. i .. "_qual"]
			
			if form then
				table.insert(pres_3sg_forms, {term = form, qualifiers = {qual}})
			end
		end
		
		-- Present participle
		for i = 2, args[2].maxindex do
			local form = args[2][i]
			local qual = args["pres_ptc" .. i .. "_qual"]
			
			if form then
				table.insert(pres_ptc_forms, {term = form, qualifiers = {qual}})
			end
		end
		
		-- Past
		for i = 2, args[3].maxindex do
			local form = args[3][i]
			local qual = args["past" .. i .. "_qual"]
			
			if form then
				table.insert(past_forms, {term = form, qualifiers = {qual}})
			end
		end
		
		-- Past participle
		local past_ptc_forms = {label = "past participle", accel = "past-participle-form-of"}
		
		if par4 then
			local qual = args["past_ptc_qual"]; if qual == "" then qual = nil end
			table.insert(past_ptc_forms, {term = par4, qualifiers = {qual}})
			
			for i = 2, args[4].maxindex do
				local form = args[4][i]
				local qual = args["past_ptc" .. i .. "_qual"]
				
				if form then
					table.insert(past_ptc_forms, {term = form, qualifiers = {qual}})
				end
			end
		end
		
		-- Are the past forms identical to the past participle forms?
		local identical = true
		
		if #past_forms ~= #past_ptc_forms then
			identical = false
		else
			for key, val in ipairs(past_forms) do
				if past_ptc_forms[key].term ~= val.term or past_ptc_forms[key].qual ~= val.qual then
					identical = false
					break
				end
			end
		end
		
		-- Insert the forms
		table.insert(data.inflections, pres_3sg_forms)
		table.insert(data.inflections, pres_ptc_forms)
		
		if #past_ptc_forms == 0 or identical then
			past_forms.label = "simple past and past participle"
			past_forms.accel = "simple-past-and-participle-form-of"
			table.insert(data.inflections, past_forms)
		else
			table.insert(data.inflections, past_forms)
			table.insert(data.inflections, past_ptc_forms)
		end
	end
}

-- Print results.
result = export.show(frame, pagename)
for f, entry in pairs(result.inflections) do
  label = entry['accel']
  for i=1,1000 do
    if entry[i] == nil then
      break
    end
    inflection = entry[i]
    if type(inflection) == "string" then
      print(label .. ': ' .. inflection)
    else
      print(label .. ': ' .. inflection.term)
    end
  end
end
