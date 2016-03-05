#!/bin/bash

OUTPUT_PATH=$1
if [[ -z "$OUTPUT_PATH" ]]; then
    OUTPUT_PATH="gen_larkmeta.go"
fi

VERSION=$(head -n 1 ../VERSION | tr -d '\n')

echo -n > "$OUTPUT_PATH"
ln() {
    echo "$@" >> "$OUTPUT_PATH"
}

ln "package larkmeta"
ln
ln "import \"github.com/codegangsta/cli\""
ln
ln "// Version is the version of the lark distribution."
ln "var Version = \"$VERSION\""
ln
ln "// Authors is the list of lark authors."
ln "var Authors = []cli.Author{"
ln "	{Name: \"The Lark authors\"},"
ln "}"
