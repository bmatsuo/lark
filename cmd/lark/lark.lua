require 'string'
require 'os'

local core = require('lark.core')

local function flatten(...)
    local flat = {}
    for i, v in pairs(arg) do
        if i == 'n' then
            -- noop
        elseif type(v) == 'table' then
            for j, v_inner in pairs(flatten(unpack(v))) do
                table.insert(flat, v_inner)
            end
        else
            table.insert(flat, v)
        end
    end
    return flat
end

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

lark.run = function (...)
    local tasks = flatten(...)
    if table.getn(tasks) == 0 then
        tasks = {lark.default_task}
    end
    for i, name in pairs(tasks) do
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

lark.environ = core.environ

lark.log = core.log

lark.exec = function (args)
    local cmd_str = lark.shell_quote(args)

    args._str = lark.shell_quote(args)
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

lark.start = function(args)
    args._str = lark.shell_quote(args) .. ' &'

    core.start(args)
end

lark.group = function (args)
    if type(args) == 'string' then
        return args
    end
    if table.getn(args) == 1 then
        args.name = args[1]
    end
    if table.getn(args) > 1 then
        error('too many positional arguments given')
    end
    core.make_group(args)
    return args[1]
end

lark.wait = function (...)
    local args = arg
    if type(args) ~= 'table' then
        args = {arg}
    end
    local result = core.wait(unpack(flatten(args)))
    if result.error then
        error(result.error)
    end
end
