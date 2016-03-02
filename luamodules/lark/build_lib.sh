#!/bin/bash

echo 'package lark' > larklib.go
echo '// LarkLib contains Lua source code for the lark module.' >> larklib.go
echo -n 'var LarkLib = `' >> larklib.go
sed 's/`/` + "`" + `/g' lark.lua >> larklib.go
echo '`' >> larklib.go
goimports -w larklib.go
