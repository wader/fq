#!/bin/sh
# what formats has a doc.md file
DOC_FORMATS=$(echo $(ls -1 $REPODIR/format/*/doc.md | sed "s#$REPODIR/format/\(.*\)\/doc.md#\1#"))
./formats_table.jq --arg doc_formats "$DOC_FORMATS"
