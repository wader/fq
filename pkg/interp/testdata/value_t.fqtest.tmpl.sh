#!/bin/sh
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers/g' | sed 's/PROMPT/mp3/g' >value_array.fqtest
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers[0].header.flags.unsynchronisation/g' | sed 's/PROMPT/mp3/g' >value_boolean.fqtest
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers[0].padding/g' | sed 's/PROMPT/mp3/g' >value_null.fqtest
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers[0].header.version/g' | sed 's/PROMPT/mp3/g' >value_number.fqtest
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers[0].header.flags/g' | sed 's/PROMPT/mp3/g' >value_object.fqtest
sed 's/CMD/fq -i -d mp3 . test.mp3/g' <value_t.fqtest.tmpl | sed 's/EXPR/.headers[0].header.magic/g' | sed 's/PROMPT/mp3/g' >value_string.fqtest

sed "s/CMD/fq -i -n '\"[]\" | json'/g" <value_t.fqtest.tmpl | sed 's/EXPR/(.)/g' | sed 's/PROMPT/json/g' >value_json_array.fqtest
sed "s/CMD/fq -i -n '\"{}\" | json'/g" <value_t.fqtest.tmpl | sed 's/EXPR/(.)/g' | sed 's/PROMPT/json/g' >value_json_object.fqtest
