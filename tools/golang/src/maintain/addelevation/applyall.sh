#!/usr/bin/env bash
cd $GITDIR/static/gpxfiles
for src in */*.withelevation; do
  mv "${src}" "${src%.withelevation}"
done