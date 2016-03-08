require 'string'
require 'table'
require 'os'

local core = require('lark.core')
local task = require('lark.task')
local doc = require('doc')

local function deprecation(old, new, mod)
    if mod then
        return string.format("deprecation warning: use %s in module '%s' instead of %s", new, mod, old)
    else
        return string.format("deprecation warning: use %s instead of %s", new, old)
    end
end

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
    doc.var[[
    verbose
    boolean -- Log more information then normal if this variable is true.
    ]] ..
    {
        default_task = nil,
        tasks = {},
        patterns  = {},
    }

lark.pattern = task.pattern
lark.task = task.create

lark.run =
    doc.desc[[An alias for run() in module lark.task]] ..
    function(...)
        if lark.default_task then
            task.default = lark.default_task
            lark.default_task = nil

            local msg = deprecation("lark.default_task", "variable default", 'lark.task')
            lark.log{msg, color='yellow'}
        end

        return task.run(unpack(arg))
    end

local function deprecated_alias(fn, old, new, mod)
    return function(...)
        local msg = deprecation(old, new, mod)
        lark.log{msg, color='yellow'}

        return fn(unpack(arg))
    end
end

lark.newpattern =
    doc.desc[[This function has been deprecated. Use the pattern() function in
              module 'lark.task' instead.

              Returns a decorator that declares a pattern matching task.  The
              pattern and the matched string are accessible through the context
              argument of the decorated task function.
              ]] ..
    deprecated_alias(task.pattern, 'lark.newpattern()', 'pattern()', 'lark.task')
lark.newtask =
    doc.desc[[This function has been deprecated. Use lark.task instead.

              A decorator that declares a task.  Assign the result to a global
              variable to run the task by name.
              ]] ..
    deprecated_alias(task.create, 'lark.newtask', 'create', 'lark.task')

lark.get_name =
    doc.sig[[ctx => string]] ..
    doc.desc[[Return the name of the task corresponding to the given context.]] ..
    doc.param[[ctx  object -- the context argument of an executing task]] ..
    deprecated_alias(task.get_name, "lark.get_name()", "get_name()", "lark.task")

lark.get_pattern =
    doc.sig[[ctx => string]] ..
    doc.desc[[
        Return the regular expression that matched the executing task.  If the
        task name was not matched against a pattern then nil is returned.
        ]] ..
    doc.param[[ctx  object -- the context argument of an executing task]] ..
    deprecated_alias(task.get_pattern, "lark.get_pattern()", "get_pattern()", "lark.task")

lark.get_param =
    doc.sig[[(ctx, name, [default]) => string]] ..
    doc.desc[[
        Return the value for the name parameter given to the task corresponding
        to ctx.  If the task context has no name parameter then default is
        returned.
        ]] ..
    doc.param[[ctx      object -- the context argument of an executing task]] ..
    doc.param[[name     string -- the name of the task parameter]] ..
    doc.param[[default  any -- returned when the task has no value for the parameter]] ..
    deprecated_alias(task.get_param, "lark.get_param()", "get_param()", "lark.task")

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
        local err = result.error
        if err then
            if args.ignore then
                if lark.verbose then
                    local msg = string.format('%s (ignored)', err)
                    lark.log{msg, color='yellow'}
                end
            else
                error(err)
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
