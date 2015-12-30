package main

// LarkLib contains Lua source code for the lark module.
var LarkLib = `require 'string'
require 'os'

local core = require('lark.core')

lark = {}

lark.default_task = nil
lark.tasks = {}

lark.task = function (name, fn)
	local t = name
	if type(t) == 'table' then
		name = t[1]
		fn = t[2]
	end

    -- print('created task: ' .. name)
	if not lark.default_task then
		lark.default_task = name
	end

	lark.tasks[name] = fn
end


local function run (name)
	local fn = lark.tasks[name]
	if not fn then
		error('no task named ' .. name)
	end
	fn()
end

lark.run = function (name)
	local t = name
    if not name then
        t = {lark.default_task}
    elseif type(t) ~= 'table' then
        t = {t}
	end
    if table.getn(t) == 0 then
        error('no tasks to run')
    end

    for i, name in pairs(t) do
        run(name)
    end
end

lark.shell_quote = function (args)
    local q = function (s)
        s = string.gsub(s, '\\', '\\\\')
        s = string.gsub(s, '"', '\\"')
        if string.find(s, '%s') then
            s = '"' .. s .. '"'
        end
        return s
    end

    local str = ''
    for i, x in pairs(args) do
        if type(i) == 'number' then
            if i > 1  then
                str = str .. ' '
            end
            if type(x) == 'string' then
                str = str .. q(x)
            else if type(x) == 'table' then
                str = str .. lark.shell_quote(x)
            else
                error(string.format('cannot quote type: %s', type(x)))
                end
            end
        end
    end

    return str
end

lark.log = core.log

lark.exec = function (args)
    local cmd_str = lark.shell_quote(args)
    lark.log{cmd_str, color='green'}

    -- This is weird... The docs online do not indicate that os.execute should
    -- return three arguments.
    local result = core.exec(args)

    if args.ignore and result.error then
		if lark.verbose then 
            local msg = string.format('%s (ignored)', result.error)
			lark.log{msg, color='yellow'}
		end
	elseif result.error then
		error(result.error)
    end
end

lark.start = lark.exec

lark.group = function (args)
    print('created group ' .. args[1])
end

lark.wait = function (args)
    local group_name = args
    if type(args) == 'table' then
        group_name = args[1]
    end
    if group_name then
        print('joined group' .. group_name)
    else
        print('joined all outstanding groups')
    end
end
`
