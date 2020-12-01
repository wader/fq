#!/usr/local/bin/fq -f jq

FieldStruct("tjo"),
FieldU8("asd"),
(FieldStruct("sdf") |
    FieldU5("b"),
    FieldU16("asd")),
(FieldU5("count") as $count | foreach range($count) as $_ (FieldArray("items");.;.) |
    FieldU8("item")
),
(FieldStruct("sdf2") |
    FieldU2("b"),
    FieldU16("asd"))
