#Module doc

    local doc = require('doc')

##Description

The doc module contains utilities for documenting Lua objects using
decorators.  Sections of documentation are declared separately using
small idiomatically named decorators.  Decorators are defined for
documenting (module) table descriptions, variables, and functions.  For
function decorators are defined to document signatures and parameter
values.

##Functions

**[desc](#function-docdesc)**

A decorator that describes an object.

**[get](#function-docget)**

Retrieve a table containing documentation for obj.

**[help](#function-dochelp)**

Print the documentation for obj.

**[param](#function-docparam)**

A decorator that describes a function parameter.

**[purge](#function-docpurge)**

Remove all declared documentation data.

**[sig](#function-docsig)**

A decorator that documents a function's signature.

**[usage](#function-docusage)**

A decorator that documents the usage of an object.

**[var](#function-docvar)**

A decorator that describes module variable (table field).

##Function doc.desc

###Signature

s => fn => fn

###Description

A decorator that describes an object.

###Parameters

**s**

string -- The object description.

##Function doc.get

###Signature

obj => table

###Description

Retrieve a table containing documentation for obj.

###Parameters

**obj**

table, function, or userdata -- The object to retrieve documentation for.

##Function doc.help

###Signature

obj => ()

###Description

Print the documentation for obj.

###Parameters

**obj**

table, function, or userdata -- The object to retrieve documentation for.

##Function doc.param

###Signature

s => fn => fn

###Description

A decorator that describes a function parameter.

###Parameters

**s**

string -- The parameter name and description separated by white space.

##Function doc.purge

###Signature

() => ()

###Description

Remove all declared documentation data.

###Parameters

****

##Function doc.sig

###Signature

s => fn => fn

###Description

A decorator that documents a function's signature.

###Parameters

**s**

string -- The function signature.

##Function doc.usage

###Signature

s => fn => fn

###Description

A decorator that documents the usage of an object.

###Parameters

**s**

string -- Text describing usage.

##Function doc.var

###Signature

s => fn => fn

###Description

A decorator that describes module variable (table field).

###Parameters

**s**

string -- The variable name and description separated by white space.

