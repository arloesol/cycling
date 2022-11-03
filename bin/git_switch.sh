#!/usr/bin/env bash

cd $GITDIR
# show available local branches and let user select
PS3="Enter a number: "

select branch in $(git branch | cut -c 3- | grep ".*-.*-.*")
do
    git checkout $branch
    exit 0
done