
def times($n):
    def _times($ctx; $i):
        if ($i < $n) then $ctx, _times($ctx; $i+1)
        else empty
        end;
    _times(.; 0);

FieldU("count"; 5) as $count |
(
    FieldStruct("asdasd") |
    FieldArray("list") | times($count)) |(FieldStruct("asdasd") as $s|
        FieldU("asdsadcv"; 2),
        (FieldU("asd"; 2) |
            if . > 1 then $s | FieldU("true"; 1)
            else $s | FieldU("false"; 2)
            end
        ),
        FieldU("vxcv"; 2)
    )


