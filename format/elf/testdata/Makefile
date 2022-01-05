TARGETS=libbbb.o libbbb.so libbbb.a a.o a_dynamic a_stripped a_static

all: $(TARGETS)

all-plaforms:
	make build PLATFORM=linux/amd64 DIR=linux_amd64
	make build PLATFORM=linux/arm64 DIR=linux_arm64
	make build PLATFORM=linux/386 DIR=linux_386
	make build PLATFORM=linux/arm/v6 DIR=linux_arm_v6
	make build PLATFORM=linux/arm/v7 DIR=linux_arm_v7

clean:
	rm -f $(TARGETS)

build:
	docker run -ti --rm --platform $(PLATFORM) -v "$(PWD):$(PWD)" -w "$(PWD)" alpine:3.15.0 sh -c 'apk add build-base && make'
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
