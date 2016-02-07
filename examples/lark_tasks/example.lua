local path = require("path")

local cmd_reusable = {'python', '-c', 'exit(1)'}

-- When `lark run` is not given any arguments lark.default_task will be
-- executed with the default parameter values.  When lark.default_task is not
-- set the first task defined will be used as the default.
lark.default_task = 'demo'

lark.task{'fail', function ()
    lark.exec{cmd_reusable}
end}

lark.task{'demo', function ()
	local file = path.join('abc', 'def')
	lark.log{file, color='blue'}

    -- Start a processes in the 'setup' execution group.
    lark.start{'sh', '-c', 'sleep 5', group='setup'}

    -- Create a group for parallel execution that cannot execute programs until
    -- after  all processes in the 'setup' group have terminated.
    build = lark.group{'build', follows='setup'}
    lark.start{'cc', '--version', group=build}

    -- Start an independent process that can begin executing before the 'build'
    -- (or 'setup') groups have completed.
    lark.start{'echo', 'an independent task'}

    -- Wait for all processes in the build group to terminate.
    lark.wait{build}

    -- If a command may terminate with a non-zero exit code the 'ignore' named
    -- argument will ensure that it does not cause the lark task to terminate.
    lark.exec{cmd_reusable, ignore=true}

    -- Wait for all outstanding execution groups.
    lark.wait()

    -- Simple logging is provided with terminal colorization.  TTY devices are
    -- detected and color is ignored when output is written to a file.
    lark.log{'everything works!', color='green'}
end}
