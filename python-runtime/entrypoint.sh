#!/bin/bash

set -eu

regex='(https?|ftp|file)://[-[:alnum:]\+&@#/%?=~_|!:,.;]*[-[:alnum:]\+&@#/%=~_|]'

if [[ $SCRIPT =~ $regex ]]
then
    curl $SCRIPT | python3 $CMD_ARGS -
else
    python3 $CMD_ARGS $SCRIPT
fi