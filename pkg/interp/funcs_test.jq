include "assert";
include "funcs";

[
  ".",
  ".a",
  ".a[0]",
  ".a[123].bb",
  ".[123].a",
  ".[123][123].a",
  ".\"b b\"",
  ".\"a \\\\ b\"",
  ".\"a \\\" b\""
][] | assert("\(.) | expr_to_path | path_to_expr"; .; expr_to_path | path_to_expr)
