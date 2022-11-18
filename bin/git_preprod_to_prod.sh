#!/usr/bin/env bash

# Accepts a version string and prints it incremented by one.
# Usage: increment_version <version> [<position>] [<leftmost>]
increment_version() {
   local usage=" USAGE: $FUNCNAME [-l] [-t] <version> [<position>] [<leftmost>]
           -l : remove leading zeros
           -t : drop trailing zeros
    <version> : The version string.
   <position> : Optional. The position (starting with one) of the number 
                within <version> to increment.  If the position does not 
                exist, it will be created.  Defaults to last position.
   <leftmost> : The leftmost position that can be incremented.  If does not
                exist, position will be created.  This right-padding will
                occur even to right of <position>, unless passed the -t flag."

   # Get flags.
   local flag_remove_leading_zeros=0
   local flag_drop_trailing_zeros=0
   while [ "${1:0:1}" == "-" ]; do
      if [ "$1" == "--" ]; then shift; break
      elif [ "$1" == "-l" ]; then flag_remove_leading_zeros=1
      elif [ "$1" == "-t" ]; then flag_drop_trailing_zeros=1
      else echo -e "Invalid flag: ${1}\n$usage"; return 1; fi
      shift; done

   # Get arguments.
   if [ ${#@} -lt 1 ]; then echo "$usage"; return 1; fi
   local v="${1}"             # version string
   local targetPos=${2-last}  # target position
   local minPos=${3-${2-0}}   # minimum position

   # Split version string into array using its periods. 
   local IFSbak; IFSbak=IFS; IFS='.' # IFS restored at end of func to                     
   read -ra v <<< "$v"               #  avoid breaking other scripts.

   # Determine target position.
   if [ "${targetPos}" == "last" ]; then 
      if [ "${minPos}" == "last" ]; then minPos=0; fi
      targetPos=$((${#v[@]}>${minPos}?${#v[@]}:$minPos)); fi
   if [[ ! ${targetPos} -gt 0 ]]; then
      echo -e "Invalid position: '$targetPos'\n$usage"; return 1; fi
   (( targetPos--  )) || true # offset to match array index

   # Make sure minPosition exists.
   while [ ${#v[@]} -lt ${minPos} ]; do v+=("0"); done;

   # Increment target position.
   v[$targetPos]=`printf %0${#v[$targetPos]}d $((10#${v[$targetPos]}+1))`;

   # Remove leading zeros, if -l flag passed.
   if [ $flag_remove_leading_zeros == 1 ]; then
      for (( pos=0; $pos<${#v[@]}; pos++ )); do
         v[$pos]=$((${v[$pos]}*1)); done; fi

   # If targetPosition was not at end of array, reset following positions to
   #   zero (or remove them if -t flag was passed).
   if [[ ${flag_drop_trailing_zeros} -eq "1" ]]; then
        for (( p=$((${#v[@]}-1)); $p>$targetPos; p-- )); do unset v[$p]; done
   else for (( p=$((${#v[@]}-1)); $p>$targetPos; p-- )); do v[$p]=0; done; fi

   echo "${v[*]}"
   IFS=IFSbak
   return 0
}


# EXAMPLE     ------------->      # RESULT
# increment_version 1             # 2
# increment_version 1 2           # 1.1
# increment_version 1 3           # 1.0.1
# increment_version 1.0.0         # 1.0.1
# increment_version 1.2.3.9       # 1.2.3.10
# increment_version 00.00.001     # 00.00.002
# increment_version -l 00.001     # 0.2
# increment_version 1.1.1.1 2     # 1.2.0.0
# increment_version -t 1.1.1 2    # 1.2
# increment_version v1.1.3        # v1.1.4
# increment_version 1.2.9 2 4     # 1.3.0.0
# increment_version -t 1.2.9 2 4  # 1.3
# increment_version 1.2.9 last 4  # 1.2.9.1

printf "making sure no open pull requests exist\n"
prlist=$(gh pr list --json number,body)
if [ "$prlist" != "[]" ]; then
  printf "ERROR - there are open pull requests - close them first\n"
  gh pr list
  exit 1
fi

# check local and origin/remote latest commits on main/preprod and prod are the same
printf "making sure local and remote branches are aligned ..."
orimain=$(git ls-remote --head --exit-code origin main | cut -f 1)
oriprod=$(git ls-remote --head --exit-code origin prod | cut -f 1)
oripreprod=$(git ls-remote --head --exit-code origin preprod | cut -f 1)
main=$(git rev-parse main)
prod=$(git rev-parse prod)
preprod=$(git rev-parse preprod)
if [ "$orimain" != "$main" ]; then
  printf "\nmain branch not aligned with origin\n"
  exit 2
fi
if [ "$oripreprod" != "$preprod" ]; then
  printf "\npreprod branch not aligned with origin\n"
  exit 3
fi
if [ "$oriprod" != "$prod" ]; then
  printf "\nprod branch not aligned with origin\n"
  exit 4
fi
printf "looks OK\n"

# if preprod and prod are aligned -> cancel
if [ "$prod" == "$preprod" ]; then
  printf "\nPreprod commit is already in prod - nothing to do - cancelling\n"
  # Todo: exit 0
fi

# show info on screen
printf "\n\ncurrent latest commit on main branch\n"
git log -n 1 main
printf "\n\ncurrent latest commit on preprod branch\n"
git log -n 1 preprod
printf "\n\ncurrent latest commit on prod branch\n"
git log -n 1 prod

prodv=$(git tag -l --points-at prod | grep -E '^([^.]*\.[^.]*){1}$')
preprodv=$(git tag -l --points-at preprod | grep -E '^([^.]*\.[^.]*){2}$') 
printf "current preprodv: $preprodv, current prod version: $prodv\n\n"
version="v$(increment_version -t ${prodv:1} 2)"
message="release $(date +"%Y-%m-%d %T")"
printf "\n\nScript will\n"
printf "  - put preprod commit in prod\n"
printf "  - and tag it with '$version' with messsage '$message'\n\n"

# ask user if all seems OK 
while true
do
  read -p "OK (y/n)? " answer
  case $answer in
   [yY]* ) printf "deploying in prod ! ...\n"
           break;;

   [nN]* ) printf "Cancelling\n"
           exit;;
   
   * )     printf "enter Y or N\n\n";;
  esac
done

printf "aligning prod with preprod locally\n"
git branch -f prod preprod || { printf "can't align prod with preprod\n"; exit 5; }
printf "pushing to remote origin\n"
git push -f origin prod || { printf "can't force push prod to main repo\n"; exit 6; }
printf "tagging this commit\n"
git tag -a "$version" -m "$message"
git push origin "$version"