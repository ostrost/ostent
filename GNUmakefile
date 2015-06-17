#!/usr/bin/env make -f

# This repo clone location (final subdirectories) defines package name thus
# it should be */github.com/[ostrost]/ostent to make package=github.com/[ostrost]/ostent
package=$(shell echo $$PWD | awk -F/ '{ OFS="/"; print $$(NF-2), $$(NF-1), $$NF }')
testpackage?=./...
singletestpackage=$(testpackage)
ifeq ($(testpackage), ./...)
singletestpackage=$(package)
endif

acepp.go=$(shell go list -f '{{.Dir}}' github.com/ostrost/ostent)/acepp/acepp.go

shareprefix=share
assets_devgo    = $(shareprefix)/assets/bindata.dev.go
assets_bingo    = $(shareprefix)/assets/bindata.bin.go
templates_devgo = $(shareprefix)/templates/bindata.dev.go
templates_bingo = $(shareprefix)/templates/bindata.bin.go

PATH=$(shell printf %s: $$PATH; echo $$GOPATH | awk -F: 'BEGIN { OFS="/bin:"; } { print $$1,$$2,$$3,$$4,$$5,$$6,$$7,$$8,$$9 "/bin"}')

xargs=xargs
ifeq (Linux, $(shell uname -s))
xargs=xargs --no-run-if-empty
endif
go-bindata=go-bindata -ignore '.*\.go'# Go regexp syntax for -ignore

.PHONY: all32 all al init test covertest coverfunc coverhtml bindata bindata-dev bindata-bin
ifneq (init, $(MAKECMDGOALS))
# before init:
# - go list would fail (for *packagefiles)
# - go test fails without dependencies installed
# - go-bindata is not installed yet

cmdname=$(notdir $(PWD))
destbin=$(shell echo $(GOPATH) | awk -F: '{ print $$1 "/bin" }')
# destbin=$(shell go list -f '{{.Target}}' $(package) | $(xargs) dirname)

define golistfiles =
{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles}}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}
endef
packagefiles=$(shell \
go list -tags   bin  -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(package) | $(xargs) \
go list -tags   bin  -f '$(golistfiles)' | sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort)
devpackagefiles=$(shell \
go list -tags '!bin' -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(package) | $(xargs) \
go list -tags '!bin' -f '$(golistfiles)' | sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort)

all: $(destbin)/$(cmdname)
all32: $(destbin)/$(cmdname).32
endif
init:
	go get -u -v \
github.com/jteeuwen/go-bindata/go-bindata \
github.com/skelterjohn/rerun \
github.com/campoy/jsonenums \
github.com/clipperhouse/gen \
code.google.com/p/go.net/html \
github.com/yosssi/ace
	cd system/operating && gen add github.com/rzab/slice
	git remote set-url origin https://$(package) # travis & tip & https://code.google.com/p/go/issues/detail?id=8850
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

system/operating/%_slice.go:     system/operating/operating.go ; cd $(dir $@) && go generate
client/enums/uint%_jsonenums.go: client/tabs.go                ; cd $(dir $@) && go generate

al: $(packagefiles) $(devpackagefiles)
# al: like `all' but without final go build $(package). For when rerun does the build

$(destbin)/$(cmdname): $(packagefiles)
	go build -ldflags -w -a -tags bin -o $@ $(package)
$(destbin)/$(cmdname).32:
	GOARCH=386 CGO_ENABLED=1 \
	go build -ldflags -w -a -tags bin -o $@ $(package)

share/assets/css/index.css: share/style/index.scss
	type sass   >/dev/null || exit 0; sass $< $@
share/assets/js/src/gen/jscript.js: share/tmp/jscript.jsx
	type jsx    >/dev/null || exit 0; jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@
share/assets/js/src/milk/index.js: share/coffee/index.coffee
	type coffee >/dev/null || exit 0; coffee -p $^ >/dev/null && coffee -o $(@D)/ $^
share/assets/js/min/index.min.js: $(shell find share/assets/js/src/ -type f)
	type r.js   >/dev/null || exit 0; cd share/assets/js/src/milk && r.js -o build.js

share/templates/index.html: share/ace.templates/index.ace share/ace.templates/defines.ace $(acepp.go)
	go run $(acepp.go) -defines share/ace.templates/defines.ace -output $@ $<
share/tmp/jscript.jsx: share/ace.templates/jscript.txt share/ace.templates/defines.ace $(acepp.go)
	go run $(acepp.go) -defines share/ace.templates/defines.ace -output $@ -javascript $<

$(templates_bingo) $(templates_devgo): $(shell find share/templates/ -type f \! -name \*.go)

$(templates_bingo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin -mode 0600 -modtime 1400000000 ./...
$(templates_devgo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -dev ./...

$(assets_bingo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags bin -mode 0600 -modtime 1400000000 -ignore js/src/ ./...
$(assets_devgo):
	cd $(@D) && $(go-bindata) -pkg $(notdir $(@D)) -o $(@F) -tags '!bin' -dev -ignore js/min/ ./...

$(assets_bingo): $(shell find \
                           share/assets/ -type f \! -name '*.go' \! -path \
                          'share/assets/js/src/*')
$(assets_bingo): share/assets/css/index.css
$(assets_bingo): share/assets/js/min/index.min.js

$(assets_devgo): $(shell find \
                      share/assets/ -type f \! -name '*.go' \! -path \
                     'share/assets/js/min/*')
$(assets_devgo): share/assets/css/index.css
$(assets_devgo): share/assets/js/src/gen/jscript.js

# spare shortcuts
bindata-bin: $(assets_bingo) $(templates_bingo)
bindata-dev: $(assets_devgo) $(templates_devgo)
bindata: bindata-dev bindata-bin

endif # END OF ifneq (init, $(MAKECMDGOALS))
