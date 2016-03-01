#Lark Scripting Reference

This document provides a reference for Lua modules specific to Lark scripts.

##Module: lark

The lark module is available as `lark`.

###lark.default_task=string

The task run when no name is given on the command line.

###lark.task{[name], fn, [pattern=string]}

Define a task that can be executed using lark.run() function (and the `lark`
command).

- **name** -- An exact name that is given to lark.run().  The name may be
  omitted if a **pattern** is defined.

- **pattern** -- A regular expression matching strings given to lark.run().

- **fn** -- The function to execute with lark.run() when the task is matched.
  The function may take a _context_ argument.

###lark.get_name(ctx)

Passed the first argument of a task function, return the name of the task being
executed.

###lark.get_pattern(ctx)

Passed the first argument of a task function, return the pattern that matched
the task name if any.

###lark.get_param(ctx, param, [default])

Passed the first argument of a task function and a parameter name, return the
parameter value given when the task was invoked.  If no value was specified
explicity with the invokation then default is returned.

###lark.environ()

Return a table containing a copy of all environment variables for the process.
The table may be altered and passed to future calls to `lark.exec{}`.

###lark.exec{command, ..., [ignore=bool]}

Execute the named command with any arguments given.

- **dir** -- The directory from which to execute the command.

- **env** -- A table containing all environment variables for the command.

- **input** -- A raw string to pass into the process standand input
  stream.

- **stdin** -- A source filename to redirect into the process standard
  input stream.

- **stdout** -- A destination filename to redirect output from the
  process standard output stream.

- **stderr** -- A destination filename to redirect output from the
  process standard error stream.

- **ignore** -- Do not terminate the task if the command exits with an error.

###lark.start{command, ..., [ignore=bool], [group=obj]}

Execute a command asynchronously.

- **group** -- Start the process in a specific group that can be recalled
  later.

###lark.run(...)

Run given tasks.  A task may be a string or a table with params to be provided
through the task's context argument.  If no task with the exact name is defined
then patterns are matched in the order they are defined.  The first task with
matching pattern will be executed.

###lark.group{name, [follows={...}]}

Get a handle on a named execution group, creating the group if necessary.

###lark.wait{[group, ...]}

Wait for outstanding processing in the named group.

##Module: path

```lua
local path = require('path')
```

The path module defines a common set of path functions that are not available
in standard lua libraries.

##path.glob(pattern)

Return a table containing paths matching the given glob pattern.

##path.join(...)

Join the given path segments.

##path.base(filepath)

Return the basename for the filepath.

##path.dir(filepath)

Return the parent directory for the filepath.

##path.exists(filepath)

Returns true if a file at filepath exists.

##path.is_dir(filepath)

Returns true if the file at filepath is a directory.
