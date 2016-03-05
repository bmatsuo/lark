#Lark [![Build Status](https://travis-ci.org/bmatsuo/lark.svg?branch=master)](https://travis-ci.org/bmatsuo/lark)

**NOTE:  Until version 1 is released the Lark Lua API should be considered
unstable.  As features are developed and practical experience is gained the Lua
API will change to better suit the needs of developers.  Any incompatible
change will be preceded by a deprecation warning and migration plan (if
necessary).  The list of open [issues](issues) of contains the of planned
upcoming incompatabilities (breaking changes).**

**BREAKING CHANGE: In lark v0.5.0 the task API is going to be altered.
Affected portions are now guarded by deprecation warnings.  To avoid these
warnings use the `lark.newtask` decorator as described in the [Getting
Started](docs/getting_started.md) guide. For information about a more permanent
migration strategy see the relevant section of the issue
[#24](https://github.com/bmatsuo/lark/issues/24)**

Lark is a Lua scripting environment that targets build systems and project
management tasks.  Lark is inspired by `make` and several build systems written
in Python.  The goal of Lark is to provide the ease and flexibility of a full
scripting environment in a portable, self-contained, and easy to integrate
package.

Python tools are great but producing consistent Python environments on
different machines, or accounting for those differences conversely, causes
overhead and headaches.  Use of Virtualenv can help with this, but
incompatibilities between Python 2 and Python 3 can still complicate things
when using these systems.

Lark attempts to avoid the problems of using Python by shipping a
self-contained interpreter and isolating module repositories for each project.
It doesn't matter what versions of the Lua interpreter are installed on
developer machines (if any).  The interpreter used by Lark can be ensured to be
consistent without interferring with normal project development.

```lua
-- global variables that can be used during project tasks
name = 'foobaz'
objects = {
    'main.o',
    'util.o',
}

-- define the "build" project task.
-- build can be executed on the command line with `lark run build`.
build = lark.newtask .. function()
    -- compile each object.
    for _, o in pairs(objects) do lark.run(o) end

    -- compile the application.
    lark.exec('gcc', '-o', name, objects)
end

-- regular expressions can match sets of task names.
build_object = lark.newpattern[[%.o$]] .. function(ctx)
    -- get the object name, construct the source path, and compile the object.
    local o = lark.newtask.get_name(ctx)
    local c = string.gsub(o, '%.o$', '.c')
    lark.exec('gcc', '-c', c)
end
```

##Core Features

- A simple to install, self-contained system.
- Pattern matching tasks (a la make).
- Parameterized tasks.
- Agnostic to target programing languages.  It doesn't matter if you are
  writing Rust, Python, or Lua.
- Builtin modules to simplify shell-style scripting.
- Extensible module system allows sharing sophisticated scripting logic
  targeted at a specific use case.
- The `LUA_PATH` environment variable is ignored. The module search directory
  is rooted under the project root for repeatable builds.
- Optional command dependency checking and memoization through [external
  tools](docs/memoize.md).
- Explicit parallel processing with execution groups for synchronization.

##Roadmap features

- More idiomatic Lua API.
- System for vendored third-party modules.  Users opt out of repeatable builds
  by explicitly ignoring the module directory in their VCS. 
- Integrated dependency checking in the same spirit of the fabricate and
  memoize.py projects.

##Documentation

New users should read the guide to [Getting Started](docs/getting_started.md).
After getting comfortable with the basics, users should consult the command
reference documentation on
[godoc.org](https://godoc.org/github.com/bmatsuo/lark/cmd/lark) to familiarize
themselves more deeply with lark command and the included Lua library.

To fully leverage Lua in Lark tasks it is recommended users unfamiliar with the
language read relevant sections of the free book [Programming in
Lua](http://www.lua.org/pil/contents.html).


##Development

For developers interested in contributing to the project see the
[docs](docs/development.md) for information about setting up a development
environment.  Also make sure to review the project contribution
[guidelines](CONTRIBUTING.md) before opening any pull requests.
