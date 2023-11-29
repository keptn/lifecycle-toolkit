#!/bin/zsh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

for i in $(find $1 -name "_index.md" -type f); do
    [ -f "$i" ] || break
    mv $i ${i/_index/index}
done