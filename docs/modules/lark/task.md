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

**find**

Return the task matching the given name.

**create**

A decorator that creates an anonymous task from a function.

**pattern**

Returns a decorator that associates the given patten with a function.

**get_pattern**

Retrieve the regular expression that matched a (running) task from the
task's context.

**get_name**

Retrieve the name of a (running) task from the task's context.

**get_param**

Retrieve the value of a task parameter (typically passed in through the
command line).

**dump**

Write all task names and patterns to standard output.

**run**

Find and run the task with the given name.

**name**

Return a decorator that gives a task function an explicit name.
