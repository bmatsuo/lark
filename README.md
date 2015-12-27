#Lark

Lark is a modern extensible build system scripted using Lua.

##Core Features

- A simple to install, self-contained system.
- Builtin modules to simplify shell-style scripting.
- Optional dependency checking through external tools.

##Roadmap features
- More idiomatic Lua API.
- Parameterized tasks.
- Pattern matching tasks (a la make).
- System for vendored third-party modules.  Users opt out of repeatable builds
  by explicitly ignoring the module directory in their VCS. 
- Parallel processing (aspirations for builtin race detection).
- Integrated dependency checking in the same spirit of the fabricate and
  memoize.py projects.

##Documentation

New users should read the guide to [Getting Started](docs/getting_started.md).
After getting comfortable with the basics, users should consult the [Lua
Scripting Reference](docs/lua.md) to familiarize themselves with the facilities
provided in a lark task.

To fully leverage Lua in Lark tasks it is recommended users unfamiliar with the
language read relevant sections of the free book [Programming in
Lua](http://www.lua.org/pil/contents.html).

##Dependency checking using fabricate or memoize.py

The python projects [fabricate](https://github.com/SimonAlfie/fabricate) or
[memoize.py](https://github.com/kgaughan/memoize.py) can be used as programs to
check the dependencies of commands.

```lua
lark.exec{'fabricate.py', 'cc', CC_OPTS, '-o', BIN, OBJECTS}
lark.exec{'memoize.py', 'cc', CC_OPTS, '-o', BIN, OBJECTS}
```
