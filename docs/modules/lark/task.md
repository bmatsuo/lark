#Module lark.task

##Description

The task module manages lark tasks.  It provides decorators that can be
used to define new tasks.  To locate and execute tasks by name the
find() and run() functions are provided respectively.

For convenience the module table acts as a decorator aliasing
task.create.  These tasks are located through the global variable index


    local task = require('task')
    mytask = task .. function() print('my task!') end
    task.run('mytask')

##Variables

**default**

string -- The task to perform when lark.run() is given no arguments.

##Functions

**create**

A decorator that creates an anonymous task from a function.

**dump**

Write all task names and patterns to standard output.

**find**

Return the task matching the given name.

**get_name**

Retrieve the name of a (running) task from the task's context.

**get_param**

Retrieve the value of a task parameter (typically passed in through the
command line).

**get_pattern**

Retrieve the regular expression that matched a (running) task from the
task's context.

**name**

Return a decorator that gives a task function an explicit name.

**pattern**

Returns a decorator that associates the given patten with a function.

**run**

Find and run the task with the given name.

##Function lark.task.create

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

##Function lark.task.dump

###Signature

() => ()

###Description

Write all task names and patterns to standard output.  Finding all
anonymous tasks is a computationally expensive process.  Do not
repeatedly call this function.

##Function lark.task.find

###Signature

name => (fn, match, pattern)

###Description

Return the task matching the given name.  If no name is given the
default task is returned.

Find first looks for named tasks with the given name.  If no explicitly
named task matches an anonymous task stored in a global variable of the
same name will be used.

When no named task matches a given name it will be tested against
pattern matching tasks.  The first pattern task to match the name will
be returned.  Pattern matching tasks will be tested in the order they
were defined.

###Parameters

**name** _(optional) string_

-- The name of a task that matches a task defined with task.create, task.name(), task.pattern().

**fn** _function_

-- The matching task function.  If nil all other return parameters
will be nil.

**match** _string_

-- The name of the matching task.

**pattern** _string_

-- The pattern which matched the task name, if a name was given and
no anonymous or explicitly named task could be matched.

##Function lark.task.get_name

###Signature

ctx => name

###Description

Retrieve the name of a (running) task from the task's context.

###Parameters

**context** _table_

-- Task context received as the first argument to a task function.

**name** _string_

-- The task's name explicity given to task.run() or derived for an
anonymous task.

##Function lark.task.get_param

###Signature

(ctx, name) => value

###Description

Retrieve the value of a task parameter (typically passed in through the
command line).

###Parameters

**context** _table_

-- Task context received as the first argument to a task function.

**name** _string_

-- The name of a task parameter.

**value** _string_

-- The value of the named parameter or nil if the task has no
parameter with the given name.  While task.run() does not restrict
the type of given parameter values all parameter values should be
treated as strings.

##Function lark.task.get_pattern

###Signature

ctx => patt

###Description

Retrieve the regular expression that matched a (running) task from the
task's context.  If the task was not matched using a pattern nil is
returned.

###Parameters

**context** _table_

-- Task context received as the first argument to a task function.

**patt** _string_

-- The pattern that matched the task name passed to task.run().

##Function lark.task.name

###Signature

name => fn => fn

###Description

Return a decorator that gives a task function an explicit name.
Explicitly named tasks are given the highest priority in matching
names given to find() and run().

###Parameters

**name** _string_

-- The task name.  A tasks may only consist of latin
alphanumerics and underscore '_'.

**fn** _function_

-- The task function.  The function may take one "context" argument
which allows runtime access to task metadeta and command line
parameters.

##Function lark.task.pattern

###Signature

patt => fn => fn

###Description

Returns a decorator that associates the given patten with a function.

###Parameters

**patt**

string -- A regular expression to match against task names

**fn**

function -- A task function which may take a context argument

##Function lark.task.run

###Signature

name => ()

###Description

Find and run the task with the given name.  See find() for more
information about task precedence.

###Parameters

**name** _string_

-- The name of the task to run.

