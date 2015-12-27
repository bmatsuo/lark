#Getting Started with Lark

Here is a simple guide that will get you up and running using Lark.  Lark uses
Lua for scripting but typical uses require minimal knowledge of the language.

If you do want to learn more about Lua scripting the best first step is the
free book [Programming in Lua](http://www.lua.org/pil/contents.html)

##Installing Lark

Install Lark by downloading one of the precompiled executable binaries.
Unarchive the binary and install it under a directory listed in your PATH
environment variable.

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

##More complex projects

For projects with more tasks or more complex task structure putting everything
in `lark.lua` can become hard to manage.  Additional task scripts can be put in
the `lark_tasks/` directory to modularize tasks and any custom functions they
need.  These task scripts are loaded in alphabetical order following `lark.lua`
(if it exists).

As tasks become more complex users will want to make greater use of the modules
provided by Lark.  Keep the scripting [reference](lua.md) on hand for
documention of all Lark programming APIs.
