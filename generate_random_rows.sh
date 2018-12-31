#!/bin/bash

rm -fr ./rand_data && mkdir rand_data

for (( j=1; j < 6; j++ )) {
    filecontent=""
    for (( i=0; i < 500000; i++ )) {
        [[ $( shuf -n 1 -i 1-3 ) -eq 3 ]] && {  filecontent="${filecontent}"$'\n'"soymatiasbarrios@gmail.com" ;  } || { filecontent="${filecontent}"$'\n'"$( dd if=/dev/urandom bs=1 count=$( shuf -n 1 -i 5-20) status=none   )"; }
    }
    echo "${filecontent}" > ./rand_data/file_$j.test
}
