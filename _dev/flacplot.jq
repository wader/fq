#!/usr/bin/env fq -s

# {
#     title: "Name of plot",
#     size: [1000, 300],
#     data: [
#         {
#             title: "Line 1",
#             points: [
#                 null,
#                 null,
#                 [1,2],
#                 [2,2],
#                 [3,3],
#                 [4,4]
#             ]
#         },
#         {
#             title: "Line 2",
#             points: [
#                 [4,1],
#                 [6,6],
#                 null,
#                 [7,7],
#                 [8,6]
#             ]
#         }
#     ]
# }

def makeplot($data):
    {
        size: [1000, 300],
        data: [$data[] | {
            points: .
        }],
    };

def remove_dups($v):
    reduce .[] as $item ([];
        if . == [] or .[0] != $item then [$item] + . 
        else . end
    ) | reverse;
def ltrim($v):
    if .[0] == $v then .[1:] |ltrim($v)
    else . end;
def rtrim($v): reverse | ltrim($v) | reverse;
def trim($v): ltrim($v) | rtrim($v);
def vjoin($v):
    reduce .[] as $item ([];[ $item,$v] + .) | reverse | .[1:];

def gnuplot($plot):
    [
        "# gnuplot < file.plot > file.svg\n",
        "$data << EOD\n",
        (
            $plot.data[] | .points |
            (trim(null) | remove_dups(null)[] | (.[0], " " , .[1], "\n")), "\n","\n"
        ),
        "EOD\n",
        "\n",
        "set terminal svg ", [
            ($plot.size | if . then ["size ", join(",")] else empty end)
        ], "\n",
        ($plot.title | if . then ["set title '", ., "'\n"] else empty end),
        "plot ",
        (
            $plot.data | to_entries | map( 
                [
                    "$data index ", .key, " ",
                    (.value.title | if . then ["title '", ., "'"] else empty end),
                    "with linespoints"
                ]
            ) | vjoin(", ")
        ), "\n"
    ] | flatten | join("");

gnuplot(makeplot(
    [
        open($FILENAME) | probe | .frames[] | . as $f |
        [[$f.end_of_header.frame_number._value, ._size/8]] +
        [.subframes[] | [$f.end_of_header.frame_number._value, .rice_partitions._value]]
    ] | transpose
))
