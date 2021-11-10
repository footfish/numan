#! /bin/bash
# Adds bash autocomplete for 1st argument ie.'command'
#
# Usage
# make sure num command is working (may need server running)
# >source ../scripts/autocomplete.sh
# or place it your .bashrc file

num >/dev/null
if [[ $? -ne 0 ]]
then
	echo "didn't work - check num command is working"
    exit 1
fi

_numan_autocomplete() {
    case $COMP_CWORD in
        1) tab_arguments=`num | tail -n 2 |head -n 1|tr -d ,|sed "s,\x1B\[[0-9;]*[a-zA-Z],,g"`
            cur="${COMP_WORDS[COMP_CWORD]}"
	        COMPREPLY=( $(compgen -W "$tab_arguments" -- ${cur}) )
	        ;;
	    *) COMPREPLY=() ;;
    esac
    return 0
    }
complete -F _numan_autocomplete num 
echo "done - check tab completion is working"

