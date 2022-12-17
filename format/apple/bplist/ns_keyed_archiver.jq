def from_ns_keyed_archiver_root:
  (  . as {"$objects": $objs, "$top": {root: $root_uid, "$0": $zero}}
  | def _f($id):
      ( .
      | $objs[$id]
      | if type == "string" then .
        elif type == "number" then .
        elif type == "boolean" then .
        elif type == "null" then .
        elif type == "array" then .
        else
          (. as {"$class": $class}
          | if $class == null then . else
            $objs[$class]."$classname" as $cname
            | if $cname == "NSDictionary" or $cname == "NSMutableDictionary" then
                ( . as {"NS.keys": $ns_keys, "NS.objects": $ns_objects}
                | [$ns_keys, $ns_objects]
                | transpose
                | map
                    (
                    ( . as [$k, $o]
                    | {key: _f($k), value: _f($o)}
                    )
                    )
                | from_entries
                )
              elif ["NSArray", "NSMutableArray", "NSSet", "NSMutableSet"]
              | any(. == $cname) then
                ( . as {"NS.objects": $ns_objects}
                | $ns_objects
                | map(_f(.))
                )
              elif $cname == "NSData" or $cname == "NSMutableData" then ."NS.Data"
              elif $cname == "NSUUID" then ."NS.uuidbytes"
              else ."class"=$cname
              end
            end
          )
        end
      );
    _f($root_uid? // $zero)
  );


def from_ns_keyed_archiver:
  (  . as {
      "$objects": $objects,
      "$top": {root: $root}
      #"$top": {"796BFF22-6712-4486-A32C-A1C5DB3273BA": $root}
    }
  | def _f($id; $seen_ids):
      def _r($id):
        if $seen_ids | has("\($id)") then "cycle-\($id)"
        else _f($id; $seen_ids | ."\($id)" = true)
        end;
      ( $objects[$id]
      | type as $type |
        if $type == "string" and . == "$null" then null
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
                if (.value | type == "object") 
                    and (.value | has("cfuid")) 
                    then .value |= _r(.cfuid) end
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
                elif $cname == "NSData" or $cname == "NSMutableData" then ."NS.Data" # TODO: will be a json string?
                elif $cname == "NSDate" then "NS.time"
                elif $cname == "NSNull" then null
                elif $cname == "NSAttributedString"
                  or $cname == "NSMutableAttributedString" then
                  _r(.NSString)
                elif $cname == "NSUUID" then ."NS.uuidbytes" # TODO: will be a json string?
                else
                  # replace class ID with classname, and dereference all cfuid values.
                  ."$class" = $cname |
                  with_entries(
                    if (.value | type == "object") 
                      and (.value | has("cfuid")) 
                      then .value |= _r(.cfuid) end
                  )
                end
              )
            end
          )
        end
      );
    def _f($id): _f($id; {"\($id)": true});
    _f($root?.cfuid // 1)
  );

