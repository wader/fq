TARGETS=libbbb.a libbbb.o libbbb.so a.o a_dynamic a_stripped a_static

FQ:=fq

all: $(TARGETS)

all-platforms:
	make build DIR=darwin_amd64

clean:
	rm -f $(TARGETS)

build:
	make
	mkdir -p $(DIR)
	mv $(TARGETS) $(DIR)
	rm $(DIR)/*.o

libbbb.so: libbbb.o
	$(CC) -shared -o $@ $+
libbbb.a: libbbb.o
	ar ru $@ $+
	ranlib $@
a_dynamic: a.o
	$(CC) -o $@ $+ -L./ -lbbb
a_stripped: a_dynamic
	strip -o $@ $<
a_static: a.o libbbb.a
	$(CC) -o $@ $+ libbbb.a

all-platforms-actual:
	make actual DIR=darwin_amd64

actual:
	cd $(DIR) && echo $(TARGETS) | tr -s '[:blank:]' '\n' | grep -ivE '.*\.o$$' | xargs -I '{}' sh -c 'echo "$$ fq -d macho v {}" > {}.fqtest && $(FQ) -d macho v {} >> {}.fqtest'