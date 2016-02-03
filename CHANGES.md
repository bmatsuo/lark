#Change log

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
