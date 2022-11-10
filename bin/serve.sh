#!/usr/bin/env bash

cd $GITDIR
hugo server --disableFastRender -e production
