#!/usr/bin/env bash

cd $GITDIR
hugo server --disableFastRender -e production >> log/serve.log 2>&1 &
disown
