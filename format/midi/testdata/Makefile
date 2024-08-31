# Invoke make from fq root folder
# 
# e.g. make -f format/midi/testdata/Makefile all

.PHONY: build
.PHONY: lint

build:
	go fmt ./format/midi/...
	make -f Makefile fq

regenerate: build
	go test ./format -run TestFormats/midi -update

run: build
	go run . -d midi dv format/midi/testdata/reference.mid

debug: build
	go run . -d midi dv format/midi/testdata/events/note-off.mid
	go run . -d midi dv format/midi/testdata/events/note-on.mid
	go run . -d midi dv format/midi/testdata/events/controller.mid
	go run . -d midi dv format/midi/testdata/events/polyphonic-pressure.mid
	go run . -d midi dv format/midi/testdata/events/program-change.mid
	go run . -d midi dv format/midi/testdata/events/channel-pressure.mid
	go run . -d midi dv format/midi/testdata/events/pitch-bend.mid

test: build
	go test ./format -run TestFormats/midi

lint: 
	make -f Makefile lint

fuzz: build
	make -f Makefile fuzz GROUP=midi

events: build
	go run . -d midi dv format/midi/testdata/events/sequence-number.mid
	go run . -d midi dv format/midi/testdata/events/text.mid
	go run . -d midi dv format/midi/testdata/events/copyright.mid
	go run . -d midi dv format/midi/testdata/events/track-name.mid
	go run . -d midi dv format/midi/testdata/events/instrument-name.mid
	go run . -d midi dv format/midi/testdata/events/lyric.mid
	go run . -d midi dv format/midi/testdata/events/marker.mid
	go run . -d midi dv format/midi/testdata/events/cuepoint.mid
	go run . -d midi dv format/midi/testdata/events/program-name.mid
	go run . -d midi dv format/midi/testdata/events/device-name.mid
	go run . -d midi dv format/midi/testdata/events/midi-channel-prefix.mid
	go run . -d midi dv format/midi/testdata/events/midi-port.mid
	go run . -d midi dv format/midi/testdata/events/tempo.mid
	go run . -d midi dv format/midi/testdata/events/smpte-offset.mid
	go run . -d midi dv format/midi/testdata/events/time-signature.mid
	go run . -d midi dv format/midi/testdata/events/key-signature.mid
	go run . -d midi dv format/midi/testdata/events/end-of-track.mid
	go run . -d midi dv format/midi/testdata/events/sequencer-specific-event.mid

	go run . -d midi dv format/midi/testdata/events/note-off.mid
	go run . -d midi dv format/midi/testdata/events/note-on.mid
	go run . -d midi dv format/midi/testdata/events/polyphonic-pressure.mid
	go run . -d midi dv format/midi/testdata/events/controller.mid
	go run . -d midi dv format/midi/testdata/events/program-change.mid
	go run . -d midi dv format/midi/testdata/events/channel-pressure.mid
	go run . -d midi dv format/midi/testdata/events/pitch-bend.mid

	go run . -d midi dv format/midi/testdata/events/sysex-message.mid
	go run . -d midi dv format/midi/testdata/events/sysex-continuation.mid
	go run . -d midi dv format/midi/testdata/events/sysex-escape.mid

fuzzes: build
	go run . -d midi dv format/midi/testdata/fuzz/15ba2fe0cd4bf33a.mid
	go run . -d midi dv format/midi/testdata/fuzz/8450de010e750ed5.mid

query: build
	./fq -d midi '.. | select(.event=="track_name")? | "\(.track_name)"' format/midi/testdata/midi/twinkle.mid 
	./fq -d midi '.. | select(.event=="tempo")?.tempo' format/midi/testdata/reference.mid
	./fq -d midi '.. | select(.event=="key_signature")?.key_signature' format/midi/testdata/reference.mid
	./fq -d midi 'grep_by(.event=="note_on") | "\(.time.tick) \(.note_on.note)"' format/midi/testdata/reference.mid

all: test lint
	make -f Makefile test
# 	go test ./format -run TestFormats/all/all.fqtest
# 	go run . -d midi dv format/midi/testdata/midi/unknown-chunks.mid
# 	go run . -d midi dv format/midi/testdata/midi/empty.mid
