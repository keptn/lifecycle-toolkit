#!/bin/bash

srcDir=$1
dstDir=$2

rsync -av --exclude='charts/*.tgz' $srcDir $dstDir