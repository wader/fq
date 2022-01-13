#!/bin/sh

for i in $(cd "$REPODIR" && ls -1 format/*/*.md | sort -t / -k 3); do
    FORMAT=$(echo "$i" | sed 's#format/.*/\(.*\).md#\1#')
    echo "### $FORMAT"
    echo
    cat "$REPODIR/$i"
    echo
done
