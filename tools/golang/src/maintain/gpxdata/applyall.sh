#!/usr/bin/env bash
cd $GITDIR/content/routes
for src in */*.new; do
  mv "${src}" "${src%.new}"
done