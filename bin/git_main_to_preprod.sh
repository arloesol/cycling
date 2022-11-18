#!/usr/bin/env bash
# NO checks !!! - should really check if local and remote are synched first !
prlist=$(gh pr list --json number,body)
if [ "$prlist" != "[]" ]; then
  printf "ERROR - there are open pull requests - close them first"
  gh pr list
  exit 5
fi
stash=$(git stash) || { printf "can't git stash"; exit 6; }
git fetch || { printf "can't git fetch"; exit 1; }
git checkout main || { printf "can't checkout main"; exit 7; }
git pull origin main --ff-only || { printf "can't git pull main"; exit 2; }
git branch -f preprod main || { printf "can't align preprod with main"; exit 3; }
git push -f origin preprod || { printf "can't force push preprod to main repo"; exit 4; }
if [ "$stash" != "No local changes to save" ]; then 
  echo "do a 'git stash pop' to get source back"
fi