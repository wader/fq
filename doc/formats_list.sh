#!/bin/sh
# what formats has a .md file
DOC_FORMATS=$(echo $(ls -1 $REPODIR/format/*/*.md | sed "s#$REPODIR/format/.*/\(.*\).md#\1#"))
./formats_list.jq --arg doc_formats "$DOC_FORMATS"
