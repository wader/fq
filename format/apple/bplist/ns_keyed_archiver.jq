def from_ns_keyed_archiver(root):
  ( if _is_decode_value and format == "bplist" then _bplist_torepr end
  | ."$objects" as $objects
  | def _f($id; $seen_ids):
      ( def _r($id):
          if $seen_ids | has("\($id)") then "cycle-\($id)"
          else _f($id; $seen_ids | ."\($id)" = true)
          end;
        $objects[$id]
      | type as $type
      | if $type == "string" and . == "$null" then null
        elif $type |
          . == "number"
          or . == "boolean"
          or . == "null"
          or . == "string" then .
        elif $type == "array" then . # TODO: does this happen?
        elif $type == "object" then
          ( ."$class" as $class
          | if $class == null then # TODO: what case is this?
              with_entries(
                if .value | type == "object" and has("cfuid") then
                  .value |= _r(.cfuid)
                end
              )
            else
              ( $objects[$class.cfuid]."$classname" as $cname
              | if $cname == "NSDictionary"
                  or $cname == "NSMutableDictionary" then
                  # transform arrays [key_id1, key_id2,...] and [obj_id1, obj_id2,..] into {key: obj, ...}
                  ( [."NS.keys", ."NS.objects"]
                  | transpose
                  | map({key: _r(.[0].cfuid), value: _r(.[1].cfuid)})
                  | from_entries
                  )
                elif $cname == "NSArray"
                  or $cname == "NSMutableArray"
                  or $cname == "NSSet"
                  or $cname == "NSMutableSet" then
                  ( ."NS.objects"
                  | map(_r(.cfuid))
                  )
                elif $cname == "NSData" or $cname == "NSMutableData" then ."NS.Data"
                elif $cname == "NSDate" then "NS.time"
                elif $cname == "NSNull" then null
                elif $cname == "NSAttributedString"
                  or $cname == "NSMutableAttributedString" then
                  _r(.NSString.cfuid)
                elif $cname == "NSUUID" then ."NS.uuidbytes"
                else
                  # replace class ID with classname, and dereference all cfuid values.
                  ( ."$class" = $cname
                  | with_entries(
                      if .value | type == "object" and has("cfuid") then
                        .value |= _r(.cfuid)
                      end
                    )
                  )
                end
              )
            end
          )
        end
      );
    root as $root
  | _f($root; {"\($root)": true})
  );

def from_ns_keyed_archiver:
  from_ns_keyed_archiver(."$top"?.root?.cfuid // error("root node not found, must specify root ID"));
