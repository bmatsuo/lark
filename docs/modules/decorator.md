#Module decorator

##Functions

**annotator**

Return a copy/concat decorator that

**create**

Return a simple copy/concat decorator using the given decorating
function.

**metatable**

Create a new metatable for a basic decorator with concat/call syntax.

##Function decorator.annotator

###Signature

(tab, prepend) => annot => obj => obj

###Description

Return a copy/concat decorator that

###Parameters

**tab** _function_

-- The table in which annotations are stored.  The table may employ
weak references but this is not a requirement.

**prepend** _(optional) boolean_

-- When true multiple (chained) annotations on the same obj will be
prepended in an array instead of being overwritten.  Prepending
makes the apparent order equal to the insertion order (opposite
call resolution order).

**annot** _any_

-- An annotation, typically a string, that is associated with given
objects.

**obj** _any_

-- A value to be decorated.

##Function decorator.create

###Signature

dec => obj => obj

###Description

Return a simple copy/concat decorator using the given decorating
function.

###Parameters

**dec** _function_

-- The decorating function.  Typically dec will return the same
object it is given though it is free to wrap or transform the
value.

**obj** _any_

-- A value to be decorated.

##Function decorator.metatable

###Signature

call => mt

###Description

Create a new metatable for a basic decorator with concat/call
syntax.

###Parameters

**call** _function_

-- The value of __call in the returned metatable

**mt** _table_

-- A metamethod table with __call and __concat set to call.

