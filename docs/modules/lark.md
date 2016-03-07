#Module lark

##Description

The lark module provides the primary Lua interface to the lark system.  A
project defines its tasks by calling the lark.task() function.  Task
functions can call other tasks by calling lark.run() function.

The first task defined is assumed to be the default task and will be run
when a user does not explicitly specify a task name to the lark command.
If this behavior is not desired the default task can be set manually using
the module variable `default_task`.

##Variables

**verbose**

boolean -- Log more information then normal if this variable is true.

##Functions

**get_pattern**

Return the regular expression that matched the executing task.

**environ**

Return a copy of the process environment as a table.

**log**

Log a message to the standard error stream.

**group**

Create a group with optional dependencies.

**newpattern**

Returns a decorator that associates the given patten with a function.

**task**

Define a new task.

**run**

An alias for run() in module lark.

**get_param**

Return the value for the name parameter given to the task corresponding
to ctx.

**wait**

Suspend execution until all processes in the specified groups have
terminated.

**pattern**

Returns a decorator that associates the given patten with a function.

**get_name**

Return the name of the task corresponding to the given context.

**start**

Start asynchronous execution of cmd.

**exec**

Execute a command

