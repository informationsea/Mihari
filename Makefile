VERSION=0.0.1
NAME=mihari

SYSTEM=$(shell uname)
HARDWARE=$(shell uname -m)

DISTNAME=$(NAME)-$(VERSION)-$(SYSTEM)-$(HARDWARE)

all: bin/mihari

bin/mihari: mihari/*.go mihari/bindata.go
	mkdir -p bin
	(cd bin; go build ../mihari/)

mihari/bindata.go: assets/libmiharihook.so
	go-bindata -o mihari/bindata.go assets

assets/libmiharihook.so: libmiharihook/*c
	make -C libmiharihook
	mkdir -p assets
	cp libmiharihook/libmiharihook.so assets/libmiharihook.so

clean:
	-rm bin/mihari
	-rmdir bin
	-rm mihari/bindata.go
	-make -C libmiharihook clean
	-rm assets/libmiharihook.so
	-rmdir assets

dist: all
	-rm -rf dist/$(DISTNAME)
	mkdir -p dist/$(DISTNAME)
	cp bin/mihari dist/$(DISTNAME)/
	cp README.md dist/$(DISTNAME)/
	(cd dist; zip -r $(DISTNAME).zip $(DISTNAME)/)
	(cd dist; tar cjf $(DISTNAME).tar.gz $(DISTNAME)/)

.PHONY: clean all dist
