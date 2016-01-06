#!/usr/bin/env gmake -f

PATH:=$(shell echo -n $$PATH:; echo $$GOPATH | sed 's,:\|$$,/bin:,g'):$$PWD/node_modules/.bin

# This repo clone location (final subdirectories) defines package name thus
# it should be */github.com/[ostrost]/ostent to make package=github.com/[ostrost]/ostent
package:=$(shell echo $$PWD | awk -F/ '{ OFS="/"; print $$(NF-2), $$(NF-1), $$NF }')
templateppackage=$(package)/cmd/ostent-templatepp

testpackage?=./...
singletestpackage=$(testpackage)
ifeq ($(testpackage), ./...)
singletestpackage=$(package)
endif

shareprefix=share
assets_devgo    = $(shareprefix)/assets/bindata.dev.go
assets_bingo    = $(shareprefix)/assets/bindata.bin.go
assets_nonego   = $(shareprefix)/assets/bindata.none.go
templates_devgo = $(shareprefix)/templates/bindata.dev.go
templates_bingo = $(shareprefix)/templates/bindata.bin.go

xargs=xargs
ifeq (Linux, $(shell uname -s))
xargs=xargs --no-run-if-empty
endif
go-bindata=go-bindata -ignore '.*\.go'# Go regexp syntax for -ignore
bingo_modtime=1400000000 # const mod time for bin bindata fileinfo
# Non-dev bindata mode templates identified by this value in templateutil.

.PHONY: all al init test covertest coverfunc coverhtml bindata bindata-dev bindata-bin check-update dev
.PHONY: all32 boot32
ifneq (init, $(MAKECMDGOALS))
# before init:
# - go list would fail (for *packagefiles)
# - go test fails without dependencies installed
# - go-bindata is not installed yet

cmdname=$(notdir $(PWD))
destbin:=$(shell echo $(GOPATH) | awk -F: '{ print $$1 "/bin" }')
# destbin:=$(shell go list -f '{{.Target}}' $(package) | $(xargs) dirname)

define golistfiles =
{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles}}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}
endef
packagefiles:=$(shell \
go list -tags   bin  -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(package) | $(xargs) \
go list -tags   bin  -f '$(golistfiles)' | sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort)
devpackagefiles:=$(shell \
go list -tags '!bin' -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(package) | $(xargs) \
go list -tags '!bin' -f '$(golistfiles)' | sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort)
templateppfiles:=$(shell \
go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(templateppackage) | $(xargs) \
go list -f '$(golistfiles)' | sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort)
templatepp=$(destbin)/$(notdir $(templateppackage))

all: $(destbin)/$(cmdname)
all32: $(destbin)/$(cmdname).32
endif
init:
	go get -u -v \
github.com/jteeuwen/go-bindata/go-bindata \
github.com/progrium/go-extpoints \
github.com/skelterjohn/rerun \
github.com/hashicorp/go-multierror \
github.com/spf13/cobra \
golang.org/x/net/html
# golang.org/x/net/html not in use yet
	go get -v $(package)
	go get -v -a -tags bin $(package)

check-update:
	npm outdated # upgrade with npm update --save-dev
	bower list # | grep latest\ is
	# update with bower install [] --save

%: %.sh # clear the implicit *.sh rule
# print-* rule for debugging. http://blog.jgc.org/2015/04/the-one-line-you-should-add-to-every.html :
print-%: ; @echo $*=$($*)

ifneq (init, $(MAKECMDGOALS))
test:
	go vet $(testpackage)
	go test -v $(testpackage)
covertest:           ; go test -coverprofile=coverage.out -covermode=count -v $(singletestpackage)
coverfunc: covertest ; go tool  cover  -func=coverage.out
coverhtml: covertest ; go tool  cover  -html=coverage.out

commands/extpoints/extpoints.go: commands/extpoints/interface.go ; cd $(dir $@) && go generate

al: $(packagefiles) $(devpackagefiles)
al: $(assets_nonego)
# al: like `all' but without final go build $(package). For when rerun does the build

$(templatepp): $(templateppfiles)
	go build -o $@ $(templateppackage)
$(destbin)/$(cmdname): $(packagefiles)
	go build -ldflags '-s -w' -a -tags bin -o $@ $(package)
$(destbin)/$(cmdname).32:
	CGO_ENABLED=1 GOARCH=386 \
	go build -ldflags '-s -w' -a -tags bin -o $@ $(package)
boot32:
	cd $(GOROOT)/src && \
	CGO_ENABLED=1 GOARCH=386 \
  ./make.bash --no-clean

dev: \
share/templates/index.html \
share/assets/js/src/bundle.js \
share/assets/css/index.css \
share/js/jsdefines.jsx

share/templates/index.html:
	type gulp >/dev/null || exit 0; gulp jade    --silent --input=./$< --output=$@

share/assets/css/index.css \
share/assets/js/src/bundle.js \
share/assets/js/min/bundle.min.js \
:
# the first prerequisite only is passed to gulp
	type gulp >/dev/null || exit 0; mkdir -p share/cache
	type gulp >/dev/null || exit 0; gulp webpack --silent --input=./$< --output=$@
# the rule above
share/assets/css/index.css:        share/style/index.scss
share/assets/js/src/bundle.js:     share/js/index.js share/js/jsdefines.jsx
share/assets/js/min/bundle.min.js: share/js/index.js share/js/jsdefines.jsx

share/templates/index.html: share/templatesorigin/index.jade
share/js/jsdefines.jsx: share/templates/index.html share/templatesorigin/jsdefines.jstmpl $(templatepp)
	$(templatepp) --output $@ --definesfrom share/templates/index.html --template share/templatesorigin/jsdefines.jstmpl

$(templates_bingo) $(templates_devgo): $(shell find share/templates/ -type f \! -name \*.go)

$(templates_bingo): ; cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin    -nomemcopy -mode 0644 -modtime $(bingo_modtime) ./...
$(templates_devgo): ; cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -dev ./...

$(assets_nonego):   ; cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags none   -ignore js/min/ -nocompress -mode 0644 -modtime $(bingo_modtime) ./...
$(assets_bingo):    ; cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin    -ignore js/src/ -nomemcopy  -mode 0644 -modtime $(bingo_modtime) ./...
$(assets_devgo):    ; cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -ignore js/min/ -dev ./...

$(assets_bingo):  $(shell find share/assets/ -type f \! -name '*.go' \! -path 'share/assets/js/src/*')
$(assets_devgo):  $(shell find share/assets/ -type f \! -name '*.go' \! -path 'share/assets/js/min/*')
$(assets_nonego): $(shell find share/assets/ -type f \! -name '*.go' \! -path 'share/assets/js/min/*')

# spare shortcuts
bindata-bin: $(assets_bingo) $(templates_bingo)
bindata-dev: $(assets_devgo) $(templates_devgo)
bindata: bindata-dev bindata-bin $(assets_nonego)

endif # END OF ifneq (init, $(MAKECMDGOALS))
