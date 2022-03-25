#!/bin/sh

# run caddy in other terminal
# caddy run

gen_test() {
    fq -Rs -L . 'include "curltrace"; from_curl_trace.send' $1.trace >$1_client
    fq -Rs -L . 'include "curltrace"; from_curl_trace.recv' $1.trace >$1_server
    echo "\$ fq -d http dv $1_client" >$1_client.fqtest
    echo "\$ fq -d http dv $1_server" >$1_server.fqtest
}

echo reqbody | curl -s --trace sinple_request.trace -d @- http://0:8080/ok >/dev/null
gen_test sinple_request
rm -f sinple_request.trace

curl -s --trace multi_request.trace http://0:8080/aaa http://0:8080/bbb >/dev/null
gen_test multi_request
rm -f multi_request.trace

curl -s --trace multi_part_single.trace --form aaa_file='@static/aaa' http://0:8080/ok >/dev/null
gen_test multi_part_single
rm -f multi_part_single.trace

curl -s --trace multi_part_multi.trace --form aaa_file='@static/aaa' --form bbb_file='@static/bbb' http://0:8080/ok >/dev/null
gen_test multi_part_multi
rm -f multi_part_multi.trace

curl -s --trace gzip.trace --compressed http://0:8080/ccc >/dev/null
gen_test gzip
rm -f gzip.trace

curl -s --trace gzip_png.trace --compressed http://0:8080/4x4.png >/dev/null
gen_test gzip_png
rm -f gzip_png.trace
