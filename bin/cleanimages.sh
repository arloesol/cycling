#!/usr/bin/env bash

rm *.svg
rm en_blog*
for src in *; do
  mv $src ${src/*none_/}
done
