TARGETS=libbbb.o libbbb.so a.o a_dynamic a_stripped a_static
GENERATED_FILES=$(TARGETS) libbbb.a

FQ:=fq

all: $(TARGETS)

all-platforms:
	make build DIR=darwin_amd64 CFLAGS='-target x86_64-apple-macos10.12'
	make build CC=clang DIR=darwin_aarch64 CFLAGS='-target arm64-apple-macos11'
	make build_fat_targets DIR=darwin_fat DIR_X86=darwin_amd64 DIR_ARM=darwin_aarch64

clean:
	rm -f $(TARGETS)

build:
	make
	mkdir -p $(DIR)
	mv $(GENERATED_FILES) $(DIR)
	rm $(DIR)/*.o

build_fat_targets:
	mkdir -p $(DIR)
	make build_fat DIR=$(DIR) OUTPUT=a_dynamic TARGET_X86=$(DIR_X86)/a_dynamic TARGET_ARM=$(DIR_ARM)/a_dynamic
	make build_fat DIR=$(DIR) OUTPUT=a_static TARGET_X86=$(DIR_X86)/a_static TARGET_ARM=$(DIR_ARM)/a_static
	make build_fat DIR=$(DIR) OUTPUT=libbbb.so TARGET_X86=$(DIR_X86)/libbbb.so TARGET_ARM=$(DIR_ARM)/libbbb.so
	make build_fat DIR=$(DIR) OUTPUT=a_stripped TARGET_X86=$(DIR_X86)/a_stripped TARGET_ARM=$(DIR_ARM)/a_stripped

# build_fat creates a FAT binary from TARGET_X86 and TARGET_ARM
# creates FAT binary $(OUTPUT) under $(DIR)
build_fat:
	lipo -create -output $(DIR)/$(OUTPUT) $(TARGET_X86) $(TARGET_ARM)

libbbb.so: libbbb.o
	$(CC) $(CFLAGS) -shared -o $@ $+
libbbb.a: libbbb.o
	ar ru $@ $+
	ranlib $@
a_dynamic: a.o
	$(CC) $(CFLAGS) -o $@ $+ -L./ -lbbb
a_stripped: a_dynamic
	strip -o $@ $<
a_static: a.o libbbb.a
	$(CC) $(CFLAGS) -o $@ $+ libbbb.a

all-platforms-actual:
	make actual DIR=darwin_amd64
	make actual DIR=darwin_aarch64
	make actual DIR=darwin_fat

# generates or actualizes the test cases
actual:
	cd $(DIR) && echo $(TARGETS) | tr -s '[:blank:]' '\n' | grep -ivE '.*\.o$$' | xargs -I '{}' sh -c 'echo "$$ fq dv {}" > {}.fqtest && $(FQ) -d macho dv {} >> {}.fqtest'