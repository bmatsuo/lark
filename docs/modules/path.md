#Module path

##Description

The path module provides utilities for working with filesystem paths.

##Functions

**[base](#function-base)**

Returns the basename of the given path.

**[dir](#function-dir)**

Returns the directory containing the given path.

**[exists](#function-exists)**

Returns true if and only if path exists

**[ext](#function-ext)**

Returns the file extension of the given path.

**[glob](#function-glob)**

Returns an array of paths that match the given pattern.

**[is_dir](#function-is_dir)**

Returns true if and only if path exists and is a directory

**[join](#function-join)**

Joins the given paths using the filesystem path separator and returns
the result.

##Function path.base

###Signature

path => string

###Description

Returns the basename of the given path.

###Parameters

**path** _A file path that may not exist_

##Function path.dir

###Signature

path => string

###Description

Returns the directory containing the given path.

###Parameters

**path** _A file path that may not exist_

##Function path.exists

###Signature

path => bool

###Description

Returns true if and only if path exists

###Parameters

**path** _A file path that may not exist_

##Function path.ext

###Signature

path => string

###Description

Returns the file extension of the given path.

###Parameters

**path** _A file path that may not exist_

##Function path.glob

###Signature

patt => [string]

###Description

Returns an array of paths that match the given pattern.

###Parameters

**patt**

Pattern using star '*' as a wildcard.

##Function path.is_dir

###Signature

path => bool

###Description

Returns true if and only if path exists and is a directory

###Parameters

**path** _A file path that may not exist_

##Function path.join

###Signature

[path] => string

###Description

Joins the given paths using the filesystem path separator and returns the result.

###Parameters

**path** _A file path that may not exist_

