#!/usr/bin/env bash

# clean up branches remote first
cd $GITDIR

printf "git fetch and prune\n"
git fetch --all --prune;

branches=$(git branch -vv | awk '/: gone]/{print $1}')

if [ "$branches" == "" ]; then
  printf "nothing to delete\n"
  exit 0
else
  printf "\nwill also delete these local branches:\n\n"
  printf $branches
  printf "\n"
fi

# ask user if OK 
while true
do
  read -p "OK (y/n)? " answer
  case $answer in
   [yY]* ) printf "Deleting ...\n\n"
           break;;

   [nN]* ) exit;;
   
   * )     printf "enter Y or N\n\n";;
  esac
done
git branch -vv | awk '/: gone]/{print $1 }' | tr '\n' '\0' | xargs -0 git branch -D;