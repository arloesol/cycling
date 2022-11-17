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
  printf "Processing $dir - $prefix\n"
  printf "  making a backup\n"
  routedir=$GITDIR/content/route/$dir
  imgdir=$GITDIR/assets/routes/gallery
  imgpfx=$imgdir/$prefix.
  gpxdir=$GITDIR/static/gpxfiles/$dir
  #tar cvfz $GITDIR/bak/$dir.$(date +%s).tgz $routedir $gpxdir $imgpfx* >/dev/null || exit 1

  printf "  emptying folders\n"
  rm -fr $routedir
  rm -fr $gpxdir
  rm -fr $imgpfx*
  
  printf "  deploying folders\n"
  mkdir $gpxdir
  mkdir $routedir
  mv $dir/gpx/* $gpxdir
  mv $dir/img/gallery/* $imgdir
  mv $dir/route/* $routedir

  find . -type d -empty -delete
  find . -type f -empty -delete

done
