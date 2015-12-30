#Memoization

Lark does not currently have direct support for command memoization (dependency
checking) like other make replacements
([fabricate](https://github.com/SimonAlfie/fabricate) or
[memoize.py](https://github.com/kgaughan/memoize.py)).  However these projects
can be used as executables for external memoization.

```lua
lark.exec{'fabricate.py', 'cc', CC_OPTS, '-o', BIN, OBJECTS}
lark.exec{'memoize.py', 'cc', CC_OPTS, '-o', BIN, OBJECTS}
```
