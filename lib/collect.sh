#!/bin/bash

OUTPUT_PATH=$1
IMPORT_PATH=github.com/bmatsuo/lark

collect() {
    for pkg in $(find . -type d); do
        grep 'var Module ' "$pkg"/*.go > /dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            echo "$pkg" | sed "s!^.!$IMPORT_PATH/lib!"
        fi
    done
}

echo -n > "$OUTPUT_PATH"
ln() {
    echo $1 >> "$OUTPUT_PATH"
}

ln "// DO NOT EDIT"
ln "// THIS IS A GENERATED FILE"
ln
ln "package lib"
ln
ln "import ("
ln "	\"github.com/bmatsuo/lark/gluamodule\""
for pkg in $(collect); do
    ln "	\"$pkg\""
done
ln ")"
ln
ln "// Modules lists every module in the library."
ln "var Modules = []gluamodule.Module{"
for pkg in $(collect); do
    name=$(basename "$pkg" | sed s/^_//)
    ln "	$name.Module,"
done
ln "}"
ln
ln "// InteralModules modules that are not general purpose and should not be
imported by scripts."
ln "var InternalModules = []gluamodule.Module{"
for pkg in $(collect); do
    if basename "$pkg" | grep ^_ > /dev/null 2>&1
    then
        name=$(basename "$pkg" | sed s/^_//)
        ln "	$name.Module,"
    fi
done
ln "}"

gofmt -w "$OUTPUT_PATH"
