#Lark Modules

Lark allows scripts to make full use of the Lua scripting language including
its module system.  But Lark does not support shared libraries using
`LUA_PATH`.  Instead all modules must be installed in the directory
`lark_modules/`.

##Third-party modules

It is recommended that developers commit third-party modules to source control.
Although at the moment there is no tool to help manage this.

##Tips for writing modules

###Avoid global variables

Avoid using global variables when writing modules.  Instead, return a table at
the end of a module file.

**lark_modules/hello.lua:**
```
local hello = {}

hello.greet = function()
    local user = os.getenv('LOGNAME') or os.getenv('USER')
    local msg = string.format('hello, %s', user)
    lark.log{msg, color='green'}
end

return hello
```

**lark.lua**
```
local hello = require('hello')

lark.task{'hello', function()
    hello.greet()
end}
```
