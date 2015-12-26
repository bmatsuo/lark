#Lark

Lark is a modern extensible build system scripted using Lua.

##Core Features

- A simple to install, self-contained system.
- Dependency checking.
- Tool for managing vendored third-party modules.
- Users opt out of repeatable builds by explicitly ignoring the module
  directory in their VCS. 
- Parallel processing (aspirations for builtin race detection).

##MVP

- [x] Run lua a named task in `lark.lua` or `lark_tasks/*.lua`.
- [x] Execute a default task that is specified by the user, otherwise the first
  task encountered.
- [x] Define multiple tasks in a lua script.
- [ ]  A task can easily and safely spawn processes and glob files.  The default
  behavior should terminate the task and exit non-zero if spawned processes
  exit non-zero.  There should be a way to ignore the exit code of a
  process.

##Tasks

The user can define an arbitrary number of named tasks.  Tasks have parameters
which are provided as unordered name-value pairs.  Parameters may have a
default value, otherwise the value for the task parameter will be nil.

Tasks are invoked from the `lark` command line tool using their name followed
by explicit values for any desired parameters using command line flag syntax.
Multiple tasks can be specified, in which case they will be excuted in serial.
To disambiguate parameter values from tasks '--' may be used to separate task
invocations.

```
lark build -O2
lark build install
lark build test --full -- release
```

##Dependency checking using strace

**NOT IMPLEMENTED**

Systems that have the `strace` command will use it by default to check the
dependencies of commands before spawning their processes.  This follows the
general idea that a program's output is dependent only on its input.  There are
exceptions to this rule.  Programs that do not take any input are expected to
be run every time.
