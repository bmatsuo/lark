#Lark [![Build Status](https://travis-ci.org/bmatsuo/lark.svg?branch=master)](https://travis-ci.org/bmatsuo/lark)

Lark is a modern extensible build system scripted using Lua.  Lark is inspired
by `make` and several build systems written in Python.  The goal of Lark is to
provide the ease and flexibility of a full scripting environment in a portable,
self-contained, and easy to integrate package.

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

##Core Features

- A simple to install, self-contained system.
- Agnostic to target programing languages.  It doesn't matter if you are
  writing Rust, Python, or Lua.
- Builtin modules to simplify shell-style scripting.
- Extensible module system allows sharing sophisticated scripting logic
  targeted at a specific use case.
- The `LUA_PATH` environment variable is ignored. The module search directory
  is rooted under the project root for repeatable builds.
- Optional command dependency checking and memoization through [external
  tools](docs/memoize.md).

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


##Development

For developers interested in contributing to the project see the
[docs](docs/development.md) for information about setting up a development
environment.  Also make sure to review the project contribution
[guidelines](CONTRIBUTING.md) before opening any pull requests.
