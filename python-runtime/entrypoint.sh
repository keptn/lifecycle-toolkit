#!/bin/sh

regex='(https?|ftp|file)://[-[:alnum:]\+&@#/%?=~_|!:,.;]*[-[:alnum:]\+&@#/%=~_|]'

if [[ $SCRIPT =~ $regex ]]
then
    curl -s $SCRIPT | python3 $CMD_ARGS -
else
    python3 $CMD_ARGS $SCRIPT
fi
