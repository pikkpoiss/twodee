.phony: build run package

APPDIR    = $(PROJECT)-linux
APPROOT   = build/$(APPDIR)
BUILDROOT = $(APPROOT)

$(BUILDROOT)/launch.sh: pkg/linux/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(BUILDROOT)/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(BUILDROOT)/resources/%: src/resources/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(APPDIR)-$(VERSION).zip: build
	cd build && zip -r $(notdir $@) $(APPDIR)

build: \
	$(BUILDROOT)/launch.sh \
	$(BUILDROOT)/$(PROJECT) \
	$(subst src/resources/,$(BUILDROOT)/resources/,$(ASSETS))

package: build/$(APPDIR)-$(VERSION).zip

run: build
	$(BUILDROOT)/launch.sh
