# Invoke make from fq root folder
# 
# e.g. make -f format/midi/testdata/Makefile midi

.PHONY: build

build:
	go fmt ./format/midi/...
	make -f Makefile fq

run: build
	go run . -d midi dv format/midi/testdata/test.mid

debug: build
# 	go run . -d midi dv format/midi/workdir/15ba2fe0cd4bf33a.mid
	go run . -d midi dv format/midi/testdata/unknown-chunks.mid

test: build
	go test ./format -run TestFormats/midi

lint: test
	make -f Makefile lint

regenerate: build
	go test ./format -run TestFormats/midi -update

query: build
	./fq -d midi '.. | select(.event=="track_name")? | "\(.name)"' format/midi/testdata/twinkle.mid 
	./fq -d midi '.. | select(.event=="tempo")?.tempo' format/midi/testdata/test.mid
	./fq -d midi '.. | select(.event=="key_signature")?.key' format/midi/testdata/test.mid
	./fq -d midi 'grep_by(.event=="note_on" or .event=="note_off") | "\(.time.tick) \(.note)"' format/midi/testdata/test.mid
