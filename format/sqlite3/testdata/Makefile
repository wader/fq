all: $(foreach f,$(wildcard *.sql), $(f).db $(f).fqtest)
clean:
	rm -f *.db *.fqtest
%.sql.db: %.sql
	cat $< | sqlite3 $@
%.sql.fqtest: %.sql.db
	echo "$$ fq 'dv, torepr' $<" > $@
