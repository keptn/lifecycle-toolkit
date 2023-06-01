#!/bin/bash

# Keptn Lifecycle Toolkit Documentation generation
#
# This script support the release of the latest version of the documentation
# The final files are available unde the docs/release folder

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR=${SCRIPT_DIR/.github\/scripts/}
DOC_DIR="${ROOT_DIR}docs"
WORK_DIR=${DOC_DIR}/release

## Create the working folder
mkdir -p $WORK_DIR

cd $WORK_DIR

## Download the latest status of the main documentation
git clone https://github.com/keptn/lifecycle-toolkit.git
cd lifecycle-toolkit
git checkout page
git pull

## Rewrite with the current content
cp -r "${DOC_DIR}/content/en/docs" ./docs/content/en

