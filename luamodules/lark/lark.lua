require 'string'
require 'table'
require 'os'

local core = require('lark.core')
local doc = require('doc')

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

local lark =
    doc.desc[[
    The lark module provides the primary Lua interface to the lark system.  A
    project defines its tasks by calling the lark.task() function.  Task
    functions can call other tasks by calling lark.run() function.

    The first task defined is assumed to be the default task and will be run
    when a user does not explicitly specify a task name to the lark command.
    If this behavior is not desired the default task can be set manually using
    the module variable `default_task`.
    ]] ..
    {
        default_task = nil,
        tasks = {},
        patterns  = {},
    }

lark.task =
    doc.sig[[(name, fn) => ()]] ..
    doc.desc[[Define a new task.]] ..
    doc.param[[name  string -- the name of the task]] ..
    doc.param[[fn    string -- the function that performs the task]] ..
    function (name, fn)
        local pattern = nil
        local t = name
        if type(t) == 'table' then
            pattern = t.pattern
            if type(t[1]) == "string" then
                name = t[1]
            end
            fn = t[table.getn(t)]
        end

        if not lark.default_task then
            lark.default_task = name
        end

        if name then
            lark.tasks[name] = fn
        end
        if pattern then
            print('pattern task: ' .. pattern)
            for _, rec in pairs(lark.patterns) do
                if rec[1] == pattern then
                    error("pattern already defined: " .. pattern)
                end
            end
            local rec = { pattern, fn }
            table.insert(lark.patterns, rec)
        end
    end


local function run (task, ctx)
    local fn = lark.tasks[task]
    if not fn then
        for _, rec in pairs(lark.patterns) do
            if string.find(task, rec[1]) then
                ctx.pattern = rec[1]
                fn = rec[2]
                break
            end
        end
    end
    if not fn then
        error('no task matching ' .. task)
    end
    fn(ctx)
end

lark.run =
    doc.sig[[(task, params) => ()]] ..
    doc.desc[[Execute the given task.]] ..
    doc.param[[task    string | nil -- A task name.  If nil is given then default_task is used.]] ..
    doc.param[[params  (optional) table -- A map from parameter names to (string) values]] ..
    function (task, params)
        if not task then
            task = lark.default_task
			if not task then error('no task to run') end
        end
        if type(task) ~= 'string' then
            error('task is not a string')
        end

        local ctx = {name = task, params = params}
        run(task, ctx)
    end

lark.get_name =
    doc.sig[[ctx => string]] ..
    doc.desc[[Return the name of the task corresponding to the given context.]] ..
    doc.param[[ctx  object -- the context argument of an executing task]] ..
    function(ctx)
        if ctx then
            return ctx.name
        end
        return nil
    end

lark.get_pattern =
    doc.sig[[ctx => string]] ..
    doc.desc[[Return the regular expression that matched the executing task or nil if the task name was not matched against a pattern.]] ..
    doc.param[[ctx  object -- the context argument of an executing task]] ..
    function(ctx)
        if ctx then
            return ctx.pattern
        end
        return nil
    end

lark.get_param =
    doc.sig[[(ctx, name, [default]) => string]] ..
    doc.desc[[Return the value for the name parameter given to the task corresponding to ctx.]] ..
    doc.param[[ctx      object -- the context argument of an executing task]] ..
    doc.param[[name     string -- the name of the task parameter]] ..
    doc.param[[default  any -- returned when the task has no value for the parameter]] ..
    function(ctx, name, default)
        if ctx and ctx.params then
            return ctx.params[name] or default
        end
        return default
    end

local function shell_quote(args)
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
                str = str .. shell_quote(x)
            else
                error(string.format('cannot quote type: %s', type(x)))
                end
            end
        end
    end

    return str
end

lark.environ =
    doc.sig[[() => envmap]] ..
    doc.desc[[Return a copy of the process environment as a table.]] ..
    core.environ

lark.log =
    doc.sig[[{msg, [color = string]} => result]] ..
    doc.desc[[Log a message to the standard error stream.]] ..
    doc.param[[msg    string -- The message to display]] ..
    doc.param[[color  string -- The color to display the message as (red, blue, ...)]] ..
    core.log

lark.exec =
    doc.sig[[cmd => output]] ..
    doc.desc[[Execute a command]] ..
    doc.param[[cmd         array -- the command to run (e.g. {'gcc', '-c', 'foo.c'}]] ..
    doc.param[[cmd.dir     string (optional) -- the directory cmd should execute in]] ..
    doc.param[[cmd.input   string (optional) -- data written to the standard input stream]] ..
    doc.param[[cmd.stdin   string (optional) -- A source filename to redirect into the standard input stream]] ..
    doc.param[[cmd.stdout  string (optional) -- A destination filename to receive output redirected from the standard output stream]] ..
    doc.param[[cmd.stderr  string (optional) -- A destination filename to receive output redirected from the standard error stream]] ..
    doc.param[[cmd.ignore  boolean (optional) -- Do not terminate execution if cmd exits with an error]] ..
    function (args)
        local cmd_str = shell_quote(args)

        args._str = shell_quote(args)
        local result = core.exec(args)

		local output = result.output
		local err = result.err
		if args.ignore and err then
			if lark.verbose then
				local msg = string.format('%s (ignored)', err)
				lark.log{msg, color='yellow'}
			end
		end
		return output, err
	end

lark.start =
    doc.sig[[cmd => ()]] ..
    doc.desc[[Start asynchronous execution of cmd.  Except where noted the cmd argument is identical to the argument of lark.exec()]] ..
    doc.param[[cmd.group  string (optional) -- the group that cmd should execute under]] ..
    function(args)
        args._str = shell_quote(args) .. ' &'

        core.start(args)
    end

lark.group =
    doc.sig[[g => string]] ..
    doc.desc[[Create a group with optional dependencies.]] ..
    doc.param[[g.name     string -- name of the group]] ..
    doc.param[[g.follows  array (optional) -- wait for the specified groups before executing any group processes]] ..
    doc.param[[g.limit    number (optional) -- limit parallel procceses among the group (in addition to global limits)]] ..
    function (args)
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

lark.wait =
    doc.sig[[[group] => nil]] ..
    doc.desc[[Suspend execution until all processes in the specified groups have terminated.]] ..
    doc.param[[group  the name of a group to wait for]] ..
    function (...)
        local args = arg
        if type(args) ~= 'table' then
            args = {arg}
        end
        local result = core.wait(unpack(flatten(args)))
        if result.error then
            error(result.error)
        end
    end

return lark
