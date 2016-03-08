#Lark compared to Python tools

Python tools like SCons, fabricate, and memoize.py are great because they give
your build system a powerful scripting environment to build increasingly
complex build and management logic as a project demands it.  Growth in
complexity is something that tools like `make` and `bash` are poor at without a
great deal of discipline from all developers.

As a full-featured programming environment Python has real module system and a
thriving ecosystem giving your scripts more power than they could reasonably
ask for.  But naively using a Python based systems can come with a number of
problems.

Producing consistent Python environments on different machines, or accounting
for those differences conversely, causes overhead and headaches.  Use of
Virtualenv can help with this, but incompatibilities between Python 2 and
Python 3 can still complicate things when using these systems.

One particular complication is the risk of tying your build system to your
application logic in undesirable ways.  If a Python application is depends on
the same modules as its (Python) build system, and they are stored in the same
module, you must update your build system in lock-step with your application
which introduces unnecessary risk of bugs in either application or build script
code.  Taking steps to avoids these potential problems is possible, but they
get complex.

The most painful case of this is probably in transition from Python 2.X to
Python 3.X. Forcing compliance and agreement between complex build scripts and
an application during a language version transition can, depending on
dependencies in use, significantly delay the transition and leave the project
using a more and more outdated interpreter version.

Obviously using Lua build scripts in a Python project will not introduce any of
these problems.  But Lark directly attempts to avoid analogous problems for
projects written primarily in Lua.  The interpreter included in lark is
self-contained, and module repositories for each project are isolated by
default.  It doesn't matter what versions of the Lua interpreter are installed
on developer machines (if any), or what global Lua modules are available.  The
interpreter used by Lark can be ensured to be consistent across developer
machines without interferring with normal project development.
