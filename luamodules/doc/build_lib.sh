#!/bin/bash

echo 'package doc' > doclib.go
echo '// DocLib contains Lua source code for the doc module.' >> doclib.go
echo -n 'var DocLib = `' >> doclib.go
cat doc.lua >> doclib.go
echo '`' >> doclib.go
goimports -w doclib.go
