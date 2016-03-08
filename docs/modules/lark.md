#Module lark

##Description

The lark module provides the primary Lua interface to the lark system.  A
project defines its tasks by calling the lark.task() function.  Task
functions can call other tasks by calling lark.run() function.

The first task defined is assumed to be the default task and will be run
when a user does not explicitly specify a task name to the lark command.
If this behavior is not desired the default task can be set manually using
the module variable `default_task`.

##Variables

**verbose**

boolean -- Log more information then normal if this variable is true.

##Functions

**[environ](#function-larkenviron)**

Return a copy of the process environment as a table.

**[exec](#function-larkexec)**

Execute a command

**[get_name](#function-larkget_name)**

Return the name of the task corresponding to the given context.

**[get_param](#function-larkget_param)**

Return the value for the name parameter given to the task corresponding
to ctx.

**[get_pattern](#function-larkget_pattern)**

Return the regular expression that matched the executing task.

**[group](#function-larkgroup)**

Create a group with optional dependencies.

**[log](#function-larklog)**

Log a message to the standard error stream.

**[newpattern](#function-larknewpattern)**

This function has been deprecated.

**[newtask](#function-larknewtask)**

This function has been deprecated.

**[pattern](#function-larkpattern)**

Returns a decorator that associates the given patten with a function.

**[run](#function-larkrun)**

An alias for run() in module lark.

**[start](#function-larkstart)**

Start asynchronous execution of cmd.

**[task](#function-larktask)**

A decorator that creates an anonymous task from a function.

**[wait](#function-larkwait)**

Suspend execution until all processes in the specified groups have
terminated.

##Function lark.environ

###Signature

() => envmap

###Description

Return a copy of the process environment as a table.

##Function lark.exec

###Signature

cmd => output

###Description

Execute a command

###Parameters

**cmd**

array -- the command to run (e.g. {'gcc', '-c', 'foo.c'}

**cmd.dir**

string (optional) -- the directory cmd should execute in

**cmd.input**

string (optional) -- data written to the standard input stream

**cmd.stdin**

string (optional) -- A source filename to redirect into the standard input stream

**cmd.stdout**

string (optional) -- A destination filename to receive output redirected from the standard output stream

**cmd.stderr**

string (optional) -- A destination filename to receive output redirected from the standard error stream

**cmd.ignore**

boolean (optional) -- Do not terminate execution if cmd exits with an error

##Function lark.get_name

###Signature

ctx => string

###Description

Return the name of the task corresponding to the given context.

###Parameters

**ctx**

object -- the context argument of an executing task

##Function lark.get_param

###Signature

(ctx, name, [default]) => string

###Description

Return the value for the name parameter given to the task corresponding
to ctx.  If the task context has no name parameter then default is
returned.

###Parameters

**ctx**

object -- the context argument of an executing task

**name**

string -- the name of the task parameter

**default**

any -- returned when the task has no value for the parameter

##Function lark.get_pattern

###Signature

ctx => string

###Description

Return the regular expression that matched the executing task.  If the
task name was not matched against a pattern then nil is returned.

###Parameters

**ctx**

object -- the context argument of an executing task

##Function lark.group

###Signature

g => string

###Description

Create a group with optional dependencies.

###Parameters

**g.name**

string -- name of the group

**g.follows**

array (optional) -- wait for the specified groups before executing any group processes

**g.limit**

number (optional) -- limit parallel procceses among the group (in addition to global limits)

##Function lark.log

###Signature

{msg, [color = string]} => result

###Description

Log a message to the standard error stream.

###Parameters

**msg**

string -- The message to display

**color**

string -- The color to display the message as (red, blue, ...)

##Function lark.newpattern

###Description

This function has been deprecated. Use the pattern() function in
              module 'lark.task' instead.

    Returns a decorator that declares a pattern matching task.  The
    pattern and the matched string are accessible through the context
    argument of the decorated task function.

##Function lark.newtask

###Description

This function has been deprecated. Use lark.task instead.

    A decorator that declares a task.  Assign the result to a global
    variable to run the task by name.

##Function lark.pattern

###Signature

patt => fn => fn

###Description

Returns a decorator that associates the given patten with a function.

###Parameters

**patt**

string -- A regular expression to match against task names

**fn**

function -- A task function which may take a context argument

##Function lark.run

###Description

An alias for run() in module lark.task

##Function lark.start

###Signature

cmd => ()

###Description

Start asynchronous execution of cmd.  Except where noted the cmd argument is identical to the argument of lark.exec()

###Parameters

**cmd.group**

string (optional) -- the group that cmd should execute under

##Function lark.task

###Signature

fn => fn

###Description

A decorator that creates an anonymous task from a function.

To call an anonymous task by name assign it to global variable and call
run() with the name of the global function.

    > mytask = task.create(function() print("my task!") end)
    > task.run('mytask')
    my task!
    >

###Parameters

**fn**

function -- A task function

##Function lark.wait

###Signature

[group] => nil

###Description

Suspend execution until all processes in the specified groups have terminated.

###Parameters

**group** _the name of a group to wait for_

