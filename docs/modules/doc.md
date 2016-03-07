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

**desc**

A decorator that describes an object.

**get**

Retrieve a table containing documentation for obj.

**help**

Print the documentation for obj.

**param**

A decorator that describes a function parameter.

**sig**

A decorator that documents a function's signature.

**usage**

A decorator that documents the usage of an object.

**var**

A decorator that describes module variable (table field).

