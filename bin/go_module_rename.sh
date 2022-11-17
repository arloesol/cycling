#!/usr/bin/env bash
if [[ "$#" -ne 1 ]]; then
    echo "Usage: $(basename "$0") NEW_MODULE_NAME"
	exit 1
fi
if [[ ! -r go.mod ]]; then
    echo "Can't find or read go.mod file in this directory"
	exit 2    
fi
OLD_MODULE_NAME=$(grep "^module " go.mod | cut -c 8-)
NEW_MODULE_NAME=$1
printf "rename module from ${OLD_MODULE_NAME} to ${NEW_MODULE_NAME}\n"

# ask user if OK to go ahead 
while true
do
  read -p "OK (y/n)? " answer
  case $answer in
   [yY]* ) break;;

   [nN]* ) exit;;
   
   * )     printf "enter Y or N\n\n";;
  esac
done

go mod edit -module ${NEW_MODULE_NAME}
find . -type f -name '*.go' -exec sed -i -e "s,${OLD_MODULE_NAME},${NEW_MODULE_NAME},g" {} \;
