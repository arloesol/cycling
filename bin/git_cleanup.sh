#!/usr/bin/env bash

# clean up branches all over the place
# Todo: before deleting show a list of barnches and ask user OK
cd $GITDIR
git fetch --all --prune; git branch -vv | awk '/: gone]/{print $1}' | xargs git branch -D;