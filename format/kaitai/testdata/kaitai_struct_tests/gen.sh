#!/bin/sh

for n in $(
    cd spec/ks
    echo *.kst
); do
    echo $n
    fq -r --arg n "$n" -f gen.jq spec/ks/"$n" >$n.fqtest
done
