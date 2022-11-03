#!/usr/bin/env bash

#set -x
issuenbr=$1
# check issue exists and is OPEN - display label and ask confirmation
state=$(gh issue view $issuenbr --json state -q .state)
if [ "$state" != "OPEN" ]; then
    printf "!! ERROR !! - issue $issuenbr is not in OPEN state\n\n"
    gh issue view $issuenbr
    exit 1;
fi

# determine branchname
title=$(gh issue view $issuenbr --json title -q .title | tr " " "_")
labels=$(gh issue view $issuenbr --json labels -q .labels[].name)
if printf "$labels" | grep feature >/dev/null 2>&1; then
  issuetyp="feat"
elif printf "$labels" | grep content >/dev/null 2>&1; then
  issuetyp="cont"
elif printf "$labels" | grep bug >/dev/null 2>&1; then
  issuetyp="bf"
else
  printf "!! ERROR !! no feature, content nor bug label on issue\nlabels:\n$labels\n";
  exit 2;
fi
branchname="$issuetyp-$issuenbr-$title"

# show info on screen
printf "issue: \n\n"
gh issue view $issuenbr
printf "\n\nScript will create new branch $branchname\n\n";

# ask user if issue is the right one 
while true
do
  read -p "OK (y/n)? " answer
  case $answer in
   [yY]* ) printf "Creating branch $branchname\n\n"
           break;;

   [nN]* ) exit;;
   
   * )     printf "enter Y or N\n\n";;
  esac
done

# make sure main is up to date
cd $GITDIR
git checkout main || (printf "can't checkout main\n\n"; exit 1) 
git pull --ff-only || (printf "can't pull main\n\n"; exit 2) 
git push || (printf "can't push main\n\n"; exit 3)

# create branch and create it also upstream
git checkout -b $branchname
git push -u origin $branchname