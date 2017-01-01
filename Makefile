PREFIX := /usr/local
BINDIR := $(PREFIX)/bin
DOTJSDIR := $(HOME)/.js

all: _build install-bin install-keypair install-plist enable-service

install-bin: _build/mdjsd
	install $< '$(BINDIR)'

install-keypair: _build/keypair
	mkdir -p '$(DOTJSDIR)/etc'
	cp _build/keypair/* '$(DOTJSDIR)/etc'

install-plist: _build/io.opts.modotjs.plist
	cp $< '$(HOME)/Library/LaunchAgents'

enable-service:
	launchctl load -w '$(HOME)/Library/LaunchAgents/io.opts.modotjs.plist'

clean:
	rm -r _build

_build:
	mkdir -p _build

_build/mdjsd: server/mdjsd.go
	go build -ldflags='-s -w' -o $@ $<

_build/keypair:
	mkdir -p _build/keypair
	cd _build/keypair && \
		openssl genrsa 2048 > server.key && \
		openssl req -new -batch -key server.key -subj '/O=modotjs/CN=localhost' > server.csr && \
		openssl x509 -req -days 3650 -signkey server.key < server.csr > server.crt

_build/io.opts.modotjs.plist: etc/io.opts.modotjs.plist
	sed \
		-e 's|{{BIN}}|$(BINDIR)/mdjsd|' \
		-e 's|{{DOTJSDIR}}|$(DOTJSDIR)|' \
		$< > $@
