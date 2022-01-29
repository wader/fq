#!/bin/sh

for i in $(cd $REPODIR/format && ls -1 | sort); do
    if [ ! -e $REPODIR/format/$i/doc.md ]; then
        continue
    fi

    echo "### $i"
    echo
    cat "$REPODIR/format/$i/doc.md"
done
