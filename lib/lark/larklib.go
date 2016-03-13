package lark

// LarkLib contains Lua source code for the lark module.
var LarkLib = `require 'string'
require 'table'
require 'os'

local core = require('lark.core')
local task = require('lark.task')
local fun = require('fun')
local doc = require('doc')

local function deprecation(old, new, mod)
    if mod then
        return string.format("deprecation warning: use %s in module '%s' instead of %s", new, mod, old)
    else
        return string.format("deprecation warning: use %s instead of %s", new, old)
    end
end

local lark =
    doc.desc[[
    The lark module provides the primary Lua interface to the lark system.  A
    project defines its tasks by calling the lark.task() function.  Task
    functions can call other tasks by calling lark.run() function.

    The first task defined is assumed to be the default task and will be run
    when a user does not explicitly specify a task name to the lark command.
    If this behavior is not desired the default task can be set manually using
    the module variable ` + "`" + `default_task` + "`" + `.
    ]] ..
    doc.var[[
    verbose boolean
    Log more information then normal if this variable is true.
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
    doc.param[[
             ctx  object
             The context argument of an executing task.
             ]] ..
    deprecated_alias(task.get_name, "lark.get_name()", "get_name()", "lark.task")

lark.get_pattern =
    doc.sig[[ctx => string]] ..
    doc.desc[[
        Return the regular expression that matched the executing task.  If the
        task name was not matched against a pattern then nil is returned.
        ]] ..
    doc.param[[
             ctx  object
             The context argument of an executing task.
             ]] ..
    deprecated_alias(task.get_pattern, "lark.get_pattern()", "get_pattern()", "lark.task")

lark.get_param =
    doc.sig[[(ctx, name, [default]) => string]] ..
    doc.desc[[
        Return the value for the name parameter given to the task corresponding
        to ctx.  If the task context has no name parameter then default is
        returned.
        ]] ..
    doc.param[[
             ctx      object
             The context argument of an executing task.
             ]] ..
    doc.param[[
             name     string
             The name of the task parameter.
             ]] ..
    doc.param[[
             default  any
             Returned when the task has no value for the parameter.
             ]] ..
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
    doc.param[[
             msg    string
             The message to display.
             ]] ..
    doc.param[[
             color  string -- The color to display the message as (red, blue, ...).
             ]] ..
    core.log

lark.exec =
    doc.sig[[(args, ..., opt) => output]] ..
    doc.desc[[
            Execute a command using the arguments given.  If opt named values
            are found in the last argument they are used with the following
            semantics.

                > lark.exec('echo', 'hello')
                echo hello
                hello
                > lark.exec('grep', 'xyz', path.glob('*.txt'), { ignore = true })
                grep xyz a.txt b.txt c.txt
                grep: exit status 1 (ignored)
                > lark.exec{'which', 'gcc', stdout = '/dev/null'}
                which gcc
                >
            ]] ..
    doc.param[[
             args        array or string
             The command to run (e.g. ('gcc', GCC_OPT, '-c', 'foo.c')).  Any
             nested arrays will be flattened to form a final array of string
             arguments.
             ]] ..
    doc.param[[
             opt         (optional) table
             Execution options interpreted by lark.  Options control logging,
             process initialization, redirection of standard I/O streams, and
             error handling.

             The opt table can contain command arguments as well for
             convenience.  So lark.exec() can be called using a single table
             argument, potentially using the special call syntax lark.exec{}.
             ]] ..
    doc.param[[
             opt.dir     string
             The directory cmd should execute in.
             ]] ..
    doc.param[[
             opt.input   string
             Data written to the standard input stream.
             ]] ..
    doc.param[[
             opt.stdin   string
             A source filename to redirect into the standard input stream.
             ]] ..
    doc.param[[
             opt.stdout  string
             A destination filename to receive output redirected from the standard output stream
             ]] ..
    doc.param[[
             opt.stderr  string
             A destination filename to receive output redirected from the standard error stream
             ]] ..
    doc.param[[
             opt.ignore  boolean
             Do not terminate execution if cmd exits with an error.
             ]] ..
    function (...)
        local args = {...}
        local opt = args[#args]
        local cmd = fun.flatten(args)
        if type(opt) == 'table' then
            for k, v in pairs(opt) do
                if type(k) == 'string' then
                    cmd[k] = v
                end
            end
        else
            opt = nil
        end

        cmd._str = shell_quote(cmd)
        local result = core.exec(cmd)
        local output = result.output
        local err = result.error
        if err then
            if opt and opt.ignore then
                if lark.verbose then
                    local msg = string.format('%s (ignored)', err)
                    lark.log{msg, color='yellow'}
                end
            elseif #cmd > 0 then
                error(string.format("%s: %s", cmd[1], err))
            else
                error(err)
            end
        end
        return output, err
    end

lark.start =
    doc.sig[[(args, ..., opt) => output]] ..
    doc.desc[[
            Start asynchronous execution of cmd.  Except where noted the cmd
            argument is identical to the argument of lark.exec()
            ]] ..
    doc.param[[
             opt.group  string (optional)
             The group that cmd should execute under.
             ]] ..
    function(...)
        local args = {...}
        local opt = args[#args]
        local cmd = fun.flatten(args)
        if type(opt) == 'table' then
            for k, v in pairs(opt) do
                if type(k) == 'string' then
                    cmd[k] = v
                end
            end
        else
            opt = nil
        end
        cmd._str = shell_quote(cmd) .. ' &'
        core.start(cmd)
    end

lark.group =
    doc.sig[[(name, opt) => name]] ..
    doc.desc[[
            Create a group with optional dependencies.  The name of the group
            is returned as if saving it to a variable is convenient for future
            calls to lark.start() and lark.wait().
            ]] ..
    doc.param[[
             name     string
             Name of the group.
             ]] ..
    doc.param[[
             opt     (optional) table
             Scheduling options for the group.  Options must be specified when
             the group is created.
             ]] ..
    doc.param[[
             opt.follows  array (optional)
             Wait for the specified groups before executing any group processes.
             ]] ..
    doc.param[[
             opt.limit    number (optional)
             Limit parallel procceses among the group (in addition to global
             limits).  The only exception to this is if the opt.limit is less
             then zero which tells lark to ignore the global limit and run an
             unlimited number of parallel processes in the group.
             ]] ..
    function (name, opt)
        if type(name) == 'table' then
            opt = name
            name = opt[1] or opt.name
        end
        if not name then
            error('no name given to group')
        end
        if not opt then
            opt = {}
        end
        opt.name = name
        core.make_group(opt)
        return name
    end

lark.wait =
    doc.sig[[(group, ...) => nil]] ..
    doc.desc[[
            Suspend execution until all processes in the specified groups have terminated.
            ]] ..
    doc.param[[
             group  string
             The name of a group to wait for.  An array of strings may be given
             instead of a single string the array and any nested arrays will be
             flattened instead a sequence of group names.
             ]] ..
    function (...)
        local args = fun.flatten({...})
        local result = core.wait(unpack(args))
        if result.error then
            error(result.error)
        end
    end

return lark
`
