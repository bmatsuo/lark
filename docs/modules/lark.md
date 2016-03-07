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

**get_name**

Return the name of the task corresponding to the given context.

**get_param**

Return the value for the name parameter given to the task corresponding
to ctx.

**group**

Create a group with optional dependencies.

**task**

Define a new task.

**run**

An alias for run() in module lark.

**exec**

Execute a command

**start**

Start asynchronous execution of cmd.

**wait**

Suspend execution until all processes in the specified groups have
terminated.

**environ**

Return a copy of the process environment as a table.

**get_pattern**

Return the regular expression that matched the executing task.

**log**

Log a message to the standard error stream.

