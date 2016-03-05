#Lark Scripting Reference

Lua reference documentation is available through the embedded REPL.

    $ lark lua
    > help()

The global `help()` function can be used to inspect the modules and functions
provided by lark.  Third-party modules may integrate their documentation with
this tool using the "doc" module.

The following REPL session loads the "doc" module, inspects its available
functions, and inspects a single function for a detailed description and
signature/schema information.

```
    > doc = require('doc')
    > help(doc)

    Functions

      help


      desc
          A decorator that describes an object.

      get


      sig
          A decorator that documents a function's signature.

      param
          A decorator that describes an function parameter.
    > help(doc.desc)

    A decorator that describes an object.

    s => fn => fn

      s
          String containing the object description
    >
```
