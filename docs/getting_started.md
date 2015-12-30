#Getting Started with Lark

Here is a simple guide that will get you up and running using Lark.  Lark uses
Lua for scripting but typical uses require minimal knowledge of the language.

##Installing Lark

Install Lark by downloading one of the precompiled
[executables](https://github.com/bmatsuo/lark/releases).  Unarchive the
executable and install it under a directory listed in your PATH environment
variable.

This will give you everything you need to run tasks with Lark.  The Lua
interpreter is included in the Lark binary and does not need to be installed
separately.

##Creating tasks 

Most simple projects will just need to create a Lua file named `lark.lua` at
their project's root.  When this file is executed a module named "lark" is
accessible and allows users to define tasks.

```lua
lark.task{'build', function()
    lark.run('generate')
    lark.exec{'go', 'build', './cmd/...'}
end}

lark.task{'generate', function()
    lark.exec{'go', 'generate', './...'}
end}
```

The above `lark.lua` file defines a task called 'generate' that runs code
generation and a task called 'build', that depends on code generation, that
builds executables.  Tasks can be run using the `lark run` command.

```sh
    lark run generate
    lark run build
    lark run           # runs the default task for the build.
```

The last line above executes the project's default task, the first task
defined.  The default task can also be set explicitly in the `lark.lua` file if
desired.

```lua
lark.default_task = 'build'
```

A user can query the list of available tasks at any time using the `lark list`
command.

```
$ lark list
  build (default)
  generate
```

##Learning more about Lua

Depending on the use case, a lark script may not need to become any more
complex than this.  To learn more about Lua scripting a good first step is the
free book [Programming in Lua](http://www.lua.org/pil/contents.html).

##Bultin modules

As tasks become more complex users will want to make greater use of the modules
provided by Lark.  Keep the scripting [reference](lua.md) on hand for
documention of all Lark programming APIs.

##Large projects

For projects with more tasks or making use of complex structures putting
everything in `lark.lua` can become hard to manage.  Additional task scripts
can be put in the `lark_tasks/` directory to modularize tasks and any custom
functions they need.  These task scripts are loaded in alphabetical order
following `lark.lua` (if it exists).

##Custom modules

When scripting requirements grow beyond the Lua standard library and the
builtin Lark modules it may be desirable to look for third-party modules or to
write custom modules to reuse within an organization.  See the reference
[documentation](modules.md) for more information about modules.
