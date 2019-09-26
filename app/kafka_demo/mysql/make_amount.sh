#!/usr/bin/env bash

total_count=$1
target_file=$2

insert_once=1000
i=0

insert_head="INSERT user_balance (user_id, symbol, amount) VALUES"

echo ${insert_head} > ${target_file}
xx=${insert_head}
for ((j=0; j<${total_count}; ++j))
do
    let i++
    if [[ $(($i % ${insert_once})) -eq 0 || ${i} -eq ${total_count} ]]
    then
        xx=${xx}"  (${i}, 'A-B', 0);"
        #echo "  (${i}, A-B, 0);" >> ${target_file}
        echo ${xx} >> ${target_file}

        if [ ${i} -ne ${total_count} ]
        then
            xx=${insert_head}
            #echo ${insert_head} >> ${target_file}
        fi
    else
        xx=${xx}"  (${i}, 'A-B', 0),"
        #echo "  (${i}, A-B, 0)," >> ${target_file}
    fi
done
