#!/usr/bin/env python3
import os
from dirsync import sync
import yaml
import argparse

parser = argparse.ArgumentParser(description="Keptn Documentation Updater")
parser.add_argument('--version', '-v', help='Keptn LT Version', default="development", required=True, dest='version')
parser.add_argument('--update-main', '-u', action='store_true', help='Update main version', dest='update_main')
parser.add_argument('--klt-repo', '-k', help='Keptn LT Repo Path', required=True, dest='klt_repo')
parser.add_argument('--klt-docs', '-d', help='Keptn LT Docs Repo Path', required=True, dest='klt_docs')
parser.add_argument('--klt-examples', '-e', help='Keptn LT Examples Repo Path', required=True, dest='klt_examples')

args = parser.parse_args()

klt_repo = args.klt_repo
klt_docs = args.klt_docs
klt_examples = args.klt_examples
version = args.version
update_main = args.update_main

if klt_docs == "" or klt_repo == "":
    print("Please provide the path to the Keptn LT and Keptn Docs Repos")
    exit(1)

# Sync the docs from the KLT repo to the docs folder, sync main-version docs to the root
sync(klt_repo + '/docs/content/docs', klt_docs + '/content/en/docs-' + version, 'sync', exclude=['^tmp', 'Makefile'], create=True)

# Update the version in the docs
with open(klt_docs + "/" + 'config.yaml', 'r') as f:
    config = f.read()
    data = yaml.safe_load(config)

if "versions" not in data['params']:
    data['params']['versions'] = []

version_exists = False
versions = data['params']['versions']
for v in versions:
    if v['version'] == version:
        version_exists = True

if not version_exists:
    versions.append({'version': version, 'url': '/docs-' + version + '/'})

versions.sort(key=lambda x: (x['version'][0].isdigit(), x['version']), reverse=True)

if update_main:
    sync(klt_docs + '/content/en/docs-' + version, klt_docs + '/content/en/docs', 'sync', exclude=['^tmp', 'Makefile'], create=True)
    data['params']['version'] = version
    sync(klt_repo + '/examples', klt_examples, 'sync', exclude=['^tmp'], create=True)


with open(klt_docs + "/" + 'config.yaml', 'w') as file:
    documents = yaml.dump(data, file)

