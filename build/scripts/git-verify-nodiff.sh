#!/bin/bash

set -euo pipefail

if [[ ! -z $(git status --porcelain) ]]; then
    echo "Staging area is dirty, please add all files created by the build to .gitignore"
    git status -s
    exit 1
fi
