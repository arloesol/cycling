#!/usr/bin/env bash
cd $GITDIR/scrapedata

declare -A prefixes
prefixes[flandersbybike]="com.flandersbybike"
prefixes[limburg]="be.visitlimburg"
prefixes[routen]="be.routen"
prefixes[vlaamsbrabant]="be.vlaamsbrabant"
prefixes[westtoer]="be.westtoer"

find . -type d -empty -delete
find . -type f -empty -delete

for dir in *; do
  prefix=${prefixes[$dir]}
  if [ "$prefix" = "" ]; then
    printf "skipping $dir\n"
    continue
  fi 
  routedir=$GITDIR/content/route/$dir

  printf "  emptying folders\n"
  rm -fr $routedir
  
  printf "  deploying folders\n"
  mkdir $routedir
  mv $dir/route/* $routedir

  find . -type d -empty -delete
  find . -type f -empty -delete

done
