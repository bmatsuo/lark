#Executing Commands

Executing commands is the primary purpose of Lark. And this document provides a
complete tutorial on the general principals and guides for how to accomplish
tasks which may not be obvious.

Lark is opinionated about how commands are executed and it may not be what
users coming from a scripting or dynamic language background expect.  The goals
of lark are that commands are executed as safely as possible to ensure that
build scripts utilizing it are reliable.  And the default execution semantics
of typical scripting languages (bash included) are prone to subtle user error
which reduces their reliability and effectiveness in providing a critical piece
of the development process, the build pipeline.

The examples in this tutorial are given assuming the user is following along
running commands in the provided interactive interpreter, `lark lua`, but they
will work equivalently in a project's `lark.lua` script file.

##Command Basics

Lark provides one primary method of executing a command, `lark.exec()`. Lark
accepts an array of strings representing the command, one string for each
argument.

    > lark.exec('cp', 'a.txt', 'b.txt')

The above example executes the "cp" command with two arguments, the string
"a.txt" and the string "b.txt".  The lark.exec() function will not accept the
command as a single string like the os.system() function requires.

    > os.system('cp a.txt b.txt')

The os.system() function is convenient but is also prone to errors especially
as commands get more complex.  But before that, consider the above command a
little longer.  When there is no file named "a.txt" cp will exit with a
non-zero status code (failure).  If the user forgets to check that the returned
value is 0 then a critical piece of the build may be missing.  Bash has similar
default semantics but does provide a facility to achieve larks behavior, `set
-e`.  The lark.exec() function will raise an exception for any command that
terminates unsuccessfully.  If command failure is known to be benign then the
user must explicitly declare this by using the pcall() function to call
lark.exec() or by passing an option table with the lark.exec() function call.

    > _, err = lark.exec('cp', 'a.txt', 'b.txt', {ignore = true})

When told to **ignore** errors the lark.exec() returns any error that occurred
as the second return value.  The first return value is used for output captured
from the program.  But in the above case the first return value will be nil
because lark.exec() was not asked to capture any output streams.  To read the
output of a command into a string another option is passed to the lark.exec()
function.


    > out = lark.exec('cat', 'b.txt', {stdout='$'})

The **stdout** redirection option can place the command's output in a specified
file, but here it uses the special sigil '$' to tell lark.exec to return the
bytes from the 'cat' program's stdout stream as a string for the processing by
the script.


##Command Construction

Sometimes commands need to be constructed piecemeal, or parameters may need to
be inserted into the rest of the command.  The lark.exec() function makes this
painless.  Commands can contain nested arrays which are flattened to construct
the final argument sequence.

    > CC_OPTS = {'-O2', '-W', -Wall'}
    > lark.exec('gcc', CC_OPTS, 'foo.c') 

If some options are conditional then Lua's builtin table manipulation
facilities can assist building an argument sequence.

    > cmd = {'go', 'test'}
    > if race then
    >>   table.insert(cmd, '-race')
    >> end
    > table.insert(cmd, test_path)
    > lark.exec(cmd)

These tasks can be performed using bash or the builtin os.execute() function.
But without care such construction can suffer several pitfalls in obscure
situations.  Take for instance a naive attempt to issue the above `cmd` which joins the strings.

    > os.execute(table.concat(cmd, ' '))

In common situations the above invokation works as desired.  But if the
`test_path` variable contains white space then the command being executed will
not be command expected.  To account for this using os.execute() arguments must
be quoted, but quoting everything correctly can be tricky especially if pipes
or other special shell syntax are involved in the command.

##Shell Commands

A distinguishing feature of os.execute() is that it relies on the system shell
to parse and execute command strings.  This can be extremely convenient because
the shell handles pipelines and redirection natively with syntax that is about
as concise as possible.

The lark.exec() function does not have direct support for pipelines, and its
redirection syntax is far less concise.  But if shell evaluation is ever
desired it is possible to execute the shell directly.

    > function shell(cmdstr) lark.exec('sh', '-c', cmdstr) end
    > shell('cat b.txt | grep "hello pipes"')

This is of course not recommended for reasons mentioned previously.  However
when it is necessary it could be slightly safer and easy to use.  Quoting
strings can be a pain but the shell is actually better at doing it than other
languages.  Consider the following command.

    > lark.exec('sh', '-c', 'echo "$0"', msg)

The command above will execute correctly regardless of the msg variable's
contents (even if it is `";/usr/bin/forkbomb "`).  This kind of construction is
not possible using the os.execute() function, so using a custom quoting
function is the only option.  Taking this into account we can rewrite the
shell() function from above to allow for safely quoting substitutions.

    > function shell(...) lark.exec('sh', '-c', ...) end
    > shell('cat b.txt | grep "$0"', 'hello pipes')

See the -c option in the `man sh` for more information.

## Asynchronous Commands

Some build tasks will require many independent commands to complete before they
may procede.  In these tasks lark scripts can use the lark.start() start
function to begin each asynchronously and finally issue a lark.wait() to
synchronize with their termination before continuing.

    > for _, name in pairs(SOURCES) do lark.start('gcc', '-c', name .. '.c') end
    > objects = {}
    > for _, name in pairs(SOURCES) do table.insert(objects, name .. '.o') end
    > lark.wait()
    > lark.exec('gcc', '-o', BIN, objects)

The above example starts compiling all of a project's C source files into
object files.  While `gcc` is compiling object files the script continues
processing and constructs the list of object filenames needed in the final
compilation step.  But, before the final binary can be compiled the lark.wait()
function which will return after all objects have been compiled successfully or
otherwise raise an error if a process terminated due to an error.  Finally the
binary is compiled.

The lark command will always call the lark.wait() function after all tasks have
terminated to clean up any asynchronous processes still running.  So it doesn't
matter if final binary in the example above is compiled using the lark.start()
function.

By default the lark command limits the number of parallel asynchronous commands
at any time to avoid overwhelming the OS/CPU with processes and dying in
context switch hell.  The default limit is equal to the number of CPU cores the
lark command can detect.  The limit can be adjusted by passing the -j option to
the `lark run` command.

Another way to provide limits with finer granularity is through the use of
_execution groups_.  An execution group is a label for commands passed to
lark.start() which makes lark schedule their execution differently.

    > lark.group{'bin', limit=2}
    > lark.start('go', 'build', './cmd/foo', {group = 'bin'})
    > lark.start('go', 'build', './cmd/bar', {group = 'bin'})
    > lark.start('go', 'build', './cmd/baz', {group = 'bin'})
    > lark.start('./compile_docs.sh')
    > lark.wait('bin')
    > lark.exec('mv', 'foo', 'bar', 'baz', 'release/bin')
    > lark.wait()
    > lark.exec('tar', '-C', 'release', '-cvzf, 'foo.tar.gz', '.')

In the above example no more than two `go build` commands will execute
simultaneously while the script `./compile_docs.sh` is free to execute in
parallel with any of them, because it has no execution group other than the
global group.  Furthermore, the first call to lark.wait() will only wait for
the `go build` commands to finish.  And the resulting binaries may be relocated
before `./compile_docs.sh` has completed.  The second call to wait will not
return until all outstanding asynchronous commands have finished, regardless of
their group.

Execution groups can also declare other groups as their dependencies and delay
their commands from executing until all processes in those groups have
completed.

    > lark.group{'bin', limit=2}
    > lark.start('go', 'build', './cmd/foo', {group = 'bin'})
    > lark.start('go', 'build', './cmd/bar', {group = 'bin'})
    > lark.start('go', 'build', './cmd/baz', {group = 'bin'})
    > lark.start('./compile_docs.sh', {group = 'docs'})
    > lark.group{'bin_reloc', follows = 'bin'}
    > lark.start('mv', 'foo', 'bar', 'baz', 'release/bin', {group = 'bin_reloc'})
    > lark.group{'archive', follows = {'docs', 'bin_reloc'}}
    > lark.start('tar', '-C', 'release', '-cvzf, 'foo.tar.gz', '.', {group = 'archive'})
    > lark.wait('archive')

The above somewhat pathological example assigns everything to a group and
starts all commands asynchronously.  This ensures that all commands execute as
soon as safely possible.  This is often not necessary or completely recommended
as developement with this style of programming can become confusing.  But
productivity gains from reduced build times may justify the introduction of
parallel command execution for some projects.
