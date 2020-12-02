#!/usr/local/bin/fq -f

def chunk:
    def chunks:
        def _chunks:
            if End then empty
            else FieldStruct("chunk") | chunk, _chunks
            end;
        if End then empty
        else . as $c | FieldArray("chunk") | _chunks | $c
        end;
    FieldUTF8("chunk_id"; 4) as $chunk_id |
    FieldU32LE("chunk_size") as $chunk_size |
    if $chunk_id == "RIFF" then
        FieldUTF8("format"; 4),
        chunks
    elif $chunk_id == "fmt " then
        FieldU16LE("audio_format"),
	    FieldU16LE("num_channels"),
	    FieldU32LE("sample_rate"),
	    FieldU32LE("byte_rate"),
	    FieldU16LE("block_align"),
	    FieldU16LE("bits_per_sample")
    elif $chunk_id == "LIST" then
        FieldUTF8("list_type"; 4),
        chunks
    elif $chunk_id == "data" then
        FieldBitBufLen("samples"; $chunk_size*8)
    else
        FieldBitBufLen("data"; $chunk_size*8)
    end;

chunk

# FieldUTF8("magic"; 4),
# FieldStruct("tjo"),
# FieldU8("asd"),
# (FieldStruct("sdf") |
#     FieldU5("b"),
#     FieldU16("asd")),
# (FieldU5("count") as $count | foreach range($count) as $_ (FieldArray("items");.;.) |
#     FieldU8("item")
# ),
# (FieldStruct("sdf2") |
#     FieldU2("b"),
#     FieldU16("asd"))
