#!/usr/bin/env make -f

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
templates_devgo = $(shareprefix)/templates/bindata.dev.go
templates_bingo = $(shareprefix)/templates/bindata.bin.go

xargs=xargs
ifeq (Linux, $(shell uname -s))
xargs=xargs --no-run-if-empty
endif
go-bindata=go-bindata -ignore '.*\.go'# Go regexp syntax for -ignore

.PHONY: all al init test covertest coverfunc coverhtml bindata bindata-dev bindata-bin
.PHONY: all32 boot32
ifneq (init, $(MAKECMDGOALS))
# before init:
# - go list would fail (for *packagefiles)
# - go test fails without dependencies installed
# - go-bindata is not installed yet

cmdname=$(notdir $(PWD))
destbin:=$(shell echo $(GOPATH) | awk -F: '{ print $$1 "/bin" }')
# destbin=$(shell go list -f '{{.Target}}' $(package) | $(xargs) dirname)

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
github.com/yosssi/ace \
golang.org/x/net/html
	go get -v $(package)
	go get -v -a -tags bin $(package)

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

share/assets/css/index.css \
share/assets/js/src/bundle.js \
share/assets/js/min/bundle.min.js \
:
# the first prerequisite only is passed to gulp
	type gulp  >/dev/null || exit 0; gulp wp --silent --input=./$< --output=$@

share/assets/css/index.css: share/style/index.scss # the above rule
share/js/jsdefines.js: share/tmp/jsdefines.jsx
	type babel >/dev/null || exit 0; babel --optional optimisation.react.constantElements --optional optimisation.react.inlineElements $^ -o $@

# "jsdefines.js" not passed to gulp/gulpfile.ls
share/assets/js/src/bundle.js:     share/js/index.js share/js/jsdefines.js # the above rule
share/assets/js/min/bundle.min.js: share/js/index.js share/js/jsdefines.js # the above rule

share/templates/index.html: share/ace.templates/index.ace share/ace.templates/defines.ace $(templatepp)
	$(templatepp) -defines share/ace.templates/defines.ace -output $@ $<
share/tmp/jsdefines.jsx: share/ace.templates/jsdefines.jstmpl share/ace.templates/defines.ace $(templatepp)
	$(templatepp) -defines share/ace.templates/defines.ace -output $@ -javascript $<

$(templates_bingo) $(templates_devgo): $(shell find share/templates/ -type f \! -name \*.go)

$(templates_bingo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin -mode 0644 -modtime 1400000000 -nomemcopy ./...
$(templates_devgo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -dev ./...

$(assets_bingo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin -mode 0644 -modtime 1400000000 -ignore js/src/ -nomemcopy ./...
$(assets_devgo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -dev -ignore js/min/ ./...

$(assets_bingo): $(shell find \
                           share/assets/ -type f \! -name '*.go' \! -path \
                          'share/assets/js/src/*')
$(assets_bingo): share/assets/css/index.css
$(assets_bingo): share/assets/js/min/bundle.min.js

$(assets_devgo): $(shell find \
                      share/assets/ -type f \! -name '*.go' \! -path \
                     'share/assets/js/min/*')
$(assets_devgo): share/assets/css/index.css
$(assets_devgo): share/assets/js/src/bundle.js

# spare shortcuts
bindata-bin: $(assets_bingo) $(templates_bingo)
bindata-dev: $(assets_devgo) $(templates_devgo)
bindata: bindata-dev bindata-bin

endif # END OF ifneq (init, $(MAKECMDGOALS))
