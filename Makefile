# tide
# Copyright (C) 2021 Dakota Walsh
# BSD 3-Clause See LICENSE in this repo for details.
.POSIX:

include config.mk

all: clean build

build:
	go build

clean:
	rm -f tide

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f tide $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/tide

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/tide

.PHONY: all build clean install uninstall
