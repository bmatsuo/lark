#Change log

##v0.5.0-dev

- Code generation makes sure that a valid version number is always defined.
  This means that installing the lark command with `go get` is a little safer
  (it is still a bad idea for production systems though).

- New command `lark lua` that behaves more like a standard interpreter than
  `lark repl`, which is now deprecated.

```sh
lark lua myscript.lua arg1 arg2 arg3
lark lua -c 'print("hello world")'
lark lua  # identical to command `lark repl`
```

- Run accepts a new -C flag that sets the working directory before loading
  files.  This option is useful for working with sub-projects or working inside
  a nested project directory.

- Documentation has been added to the project's custom modules (in
  `lark_modules/`).  This serves as an example of how to document objects in
  lua code.

##v0.4.0

- Fix documentation of available modules available through the help() command.

##v0.4.0-beta1

- Documentation is now accessible through the REPL (`lark repl`) using the
  global function `help()`.

- Modules can be documented using decoractors in the "doc" module.

- Remove the function `lark.shell_quote()`.  It's direct use was never
  recommended.

- New module "lark.task" that will deprecate task-related functions and
  variables in the "lark" module and replace the API/syntax used for
  `lark.task()`.  See #24, #25, and #30.

- Project reorganization.  The contents of the `luamodules/` directory can now
  be found in the `lib/` directory.

##v0.3.1

- Fix bug retreiving captured output from lark.exec().  The lark.exec()
  function returns two values, captured output and any (ignored) error
  encountered.

##v0.3.0

- Added optional filename parameters to `lark.exec{}`: **stdin**, **stdout**,
  and **stderr**.  See the [docs](docs/lua.md) for more information.
- Added parallel processing, limited by the -j flag to `lark run`/`lark make`.
- Added the ability to capture exec output.  Additionally, exec output can tee
  to standard streams and files/memory.
- Tasks take an optional context parameter that allows access to task metadata
  and parameters.
- Tasks can now define a regular expression pattern instead of a single name to
  match multiple values and behave dynamically according to the names they
  match.
- Simple REPL for developers to experiment.  There are some quirks but it is
  essentially the same as the lua5.1 interpreter.
  

##v0.2.0

- The `lark.run` function can accept variable arguments without wrapping them
  in a table.  This makes `lark.run(...)` equivalent to `lark.run{...}`.
- Added subcommand "make" as an alias of "run" to make invocations more
  idiomatic (e.g. `lark make release`).
- Added optional parameters to `lark.exec{}`: **dir**, and **env**
- Builtin lua modules now have nominal tests in place with a test framework for
  future modules.
- Dependencies have been updated.  See the [glide.lock](glide.lock) file for
  exact version information.

##v0.1.0

This is the first release (MVP) of Lark. Basic scripting functionality is
available in lark scripts, and the scripts are evaluated basically correctly.

The build scripts for Lark and several other test projects have been written as
Lark scripts.  This has served as validation of the system as working and
something that feels natural and can be built on.

- Portable command evaluation though a custom Go function.
- Builtin 'path' module for basic path manipulation facilities.
- Logging with colors for TTY devices.
- Environment `LUA_PATH` is ignored, all modules loaded from the
  `lark_modules/` directory.
