#Lark Scripting Reference

This document provides a reference for Lua modules specific to Lark scripts.

##Module: lark

The lark module is available as `lark`.

###lark.default_task=string

The task run when no name is given on the command line.

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

**NOTE**: `lark.start{}` is currently aliased to `lark.exec{}` and will not run
processes concurrently.

Execute a command asynchronously and return its execution group.

- **group** -- Start the process in a specific group that can be recalled
  later.

###lark.run{...}

Run given tasks.  All arguments given are flattened and each task named run.

###lark.group{name}

Get a handle on a named execution group, creating the group if necessary.

###lark.wait{group}

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
