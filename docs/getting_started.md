#Getting Started with Lark

Here is a simple guide that will get new users up and running with Lark.  Lark
uses Lua for scripting but typical uses like those demonstrated in this guide
require minimal knowledge of the language.

This document is meant to accompany the command reference documentation
available [godoc.org](https://github.com/bmatsuo/lark/cmd/lark).  It will help
get the lark command installed and highlight the most basic and relevant parts
of its usage.  Finally, where the guide leaves off it will provide links to
more detailed documentation, outside of the command reference, if required.

##Installing Lark

Install Lark by downloading one of the precompiled
[executables](https://github.com/bmatsuo/lark/releases).  Unarchive the
executable and install it under a directory listed in your PATH environment
variable.

Instead installing of a stable release, the latest (unstable) version of lark
can be installed on the command line using `go get`.

    # this is not recommended for most users
    go get github.com/bmatsuo/lark/cmd/lark

Now the system has everything needed to run project tasks with Lark.  The Lua
interpreter is included in the Lark binary and does not need to be installed
separately.

##Creating tasks 

Most simple projects will just need to create a Lua file named `lark.lua` at
the project root.  When this file is executed with the `lark` command a module
named "lark" is accessible and allows users to define tasks.

**lark.lua**
```lua
build = lark.task .. function()
    lark.run('generate')
    lark.exec('go', 'build', './cmd/...')
end

generate = lark.task .. function()
    lark.exec('go', 'generate', './...')
end
```

The above `lark.lua` file defines a task called 'generate' that runs executes
the code generation tool `go generate ./...` and a task called 'build', that
builds executables with `go build ./cmd/...` after code generation has
completed successfully.  Tasks can be run using the `lark run` command.

```sh
lark run generate
lark run build
lark run           # runs the default task, "build".
```

The last line above executes the project's default task, the first task
defined.  The default task can also be set explicitly in the `lark.lua` file if
desired.

```lua
local task = require('lark.task')
task.default = 'build'
```

A user can query the list of available tasks at any time using the `lark list`
command.

```
$ lark list
=   build   (default)
=   generate
```

The '=' each lines first field indicates that the task matches the exact string
("build" or "generate").  The default task is indicated by the optional third
field.

##Executing commands

The `lark.lua` file above shows two examples of executing commands.  For a more
in depth look at how commands are executed see the [tutorial](exec.md) on the
lark.exec() function.

##Learning more about Lua

Depending on the use case, a lark script may not need to become any more
complex than this.  To learn more about Lua scripting a good first step is the
free book [Programming in Lua](http://www.lua.org/pil/contents.html).

##Builtin modules

As tasks become more complex users will want to make greater use of the Lua
module library provided by Lark.  Keep the scripting [reference](lua.md) on
hand for documention of all Lark programming APIs.

##Large projects

For projects with more tasks or making use of complex structures putting
everything in `lark.lua` can become hard to manage.  Additional task scripts
can be put in the `lark_tasks/` directory to modularize tasks and any custom
functions they need.  These task scripts are loaded in alphabetical order
following `lark.lua` (if it exists).

Files in `lark_tasks/` cannot see any of local variables set by `lark.lua`.
However global variables and package variables can be shared between `lark.lua`
and other task files.

**lark.lua**
```lua
local x = 1
y = 2
```

**lark_tasks/mytask.lua**
```lua
mytask = lark.task .. function()
    print(x)
    print(y)
end
```

In the above example the command `lark run mytask` will print "nil" followed by
"2" because `mytask.lua` cannot see the local variables from `lark.lua` and
only `y` is global.

##Custom modules

When scripting requirements grow beyond the Lua standard library and the
builtin Lark modules it may be desirable to look for third-party modules or to
write custom modules to reuse within an organization.  See the reference
[documentation](modules.md) for more information about modules.
