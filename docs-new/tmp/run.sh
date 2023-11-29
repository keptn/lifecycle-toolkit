#!/bin/zsh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

for i in $(find $1 -name "*.md" -type f); do
    [ -f "$i" ] || break
    go run ${DIR}/main.go $i
done