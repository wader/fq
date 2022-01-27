all: $(patsubst %.ksy,%.kst,$(wildcard *.ksy))

clean:
	rm -f *.kst

%.kst: %.ksy
	docker run -v "$(PWD):/share" -it --entrypoint=ksdump kaitai/ksv \
		-f json /dev/zero $< | \
		tail +5 | \
		fq -L . -r 'include "to_kst"; to_kst | to_yaml' \
		> $@
