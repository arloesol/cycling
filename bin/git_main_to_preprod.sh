#!/usr/bin/env bash
# NO checks !!! - should really check if local and remote are synched first !
git fetch || ( printf "can't git fetch"; exit 1 )
git pull origin main --ff-only || ( printf "can't git pull main"; exit 2 )
git branch -f preprod main || ( printf "can't align preprod with main"; exit 3 )
git push -f origin preprod || ( printf "can't force push preprod to main repo"; exit 4 )