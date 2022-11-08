#!/usr/bin/env bash
cd $GITDIR/data

declare -A prefixes
prefixes[flandersbybike]="com.flandersbybike"
prefixes[limburg]="be.visitlimburg"
prefixes[routen]="be.routen"
prefixes[vlaamsbrabant]="be.vlaamsbrabant"
prefixes[westtoer]="be.westtoer"

for dir in *; do
  prefix=${prefixes[$dir]} 
  printf "making a backup of $dir - $prefix\n"
  routedir=$GITDIR/content/route/$dir
  imgpfx=$GITDIR/assets/routes/gallery/$prefix.
  gpxdir=$GITDIR/static/gpxfiles/$dir
  tar cvfz $GITDIR/bak/$dir.$(date +%s).tgz $routedir $gpxdir $imgpfx*
done
