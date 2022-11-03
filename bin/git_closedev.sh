#!/usr/bin/env bash

# script only does ff-only merges - any more complex merges need to be done manually before launching this script
cd $GITDIR
branchname=$(git branch --show-current)
nbrdesc=${branchname#*-}
issuenbr=${nbrdesc%%-*}
printf "on branch $branchname for issue $issuenbr\n\n"

# make sure all items are commited and pushed upstream
git add --all
git commit -m "close dev $branchname - closes #$issuenbr"
git push || exit 10

# make sure any changes in main are included in our branch
git checkout main
git pull --ff-only || (printf "can't pull main ! \n\n"; exit 1)
git push || (printf "can't push main ! \n\n"; exit 11)
git checkout $branchname
git pull --ff-only || (printf "can't pull branch \n\n"; exit 2)
git merge main --ff-only || (printf "merge issue - fix manually !\n\n"; exit 3)
git push || (printf "can't push merged branch\n\n"; exit 12)

# create pull request into main branch
gh pr create --base main --title "close $branchname" --body "closes #$issuenbr"