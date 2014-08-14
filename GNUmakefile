#!/usr/bin/env make -f

binassets_develgo         = src/share/assets/bindata.devel.go
binassets_productiongo    = src/share/assets/bindata.production.go
bintemplates_develgo      = src/share/templates/bindata.devel.go
bintemplates_productiongo = src/share/templates/bindata.production.go
templates_dir             = src/share/templates/
templates_files           = index.html usepercent.html tooltipable.html
templates_html=$(addprefix $(templates_dir), $(templates_files))
bindir=bin/$(shell uname -sm | awk '{ sub(/x86_64/, "amd64", $$2); print tolower($$1) "_" $$2; }')
GOPATH=$(shell echo $$GOPATH):$(PWD)
PATH=$(shell printf %s: $$PATH; echo $$GOPATH | awk -F: 'BEGIN { OFS="/bin:"; } { print $$1,$$2,$$3,$$4,$$5,$$6,$$7,$$8,$$9 "/bin"}')

gobindata=go-bindata -ignore '.*\.go'

.PHONY: all test bootstrap bootstrap_develgo
all: $(bindir)/ostent
test:
	go test -v ./...
bootstrap:
	go get -v github.com/jteeuwen/go-bindata/go-bindata
	$(MAKE) $(MFLAGS) bootstrap_develgo
	go get -v github.com/skelterjohn/rerun ostent
	go get -v -tags production ostent
	rm -f bin/ostent
bootstrap_develgo: $(binassets_develgo) $(bintemplates_develgo)

%: %.sh # clear the implicit *.sh rule covering ./ostent.sh

$(bindir)/%:
	@echo '* Sources:' $^
	go build -o $@ $(patsubst src////%,%,$|)

$(bindir)/amberpp: | src////amberp/amberpp
$(bindir)/ostent:  | src////ostent

ifeq (, $(findstring bootstrap, $(MAKECMDGOALS)))
$(bindir)/amberpp: $(shell env GOPATH=$(GOPATH) go list -f '\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}' amberp/amberpp | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr

$(bindir)/ostent: $(shell \
env GOPATH=$(GOPATH) go list -tags production -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostent | xargs \
env GOPATH=$(GOPATH) go list -tags production -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	go build -tags production -o $@ ostent

$(bindir)/jsmakerule: $(binassets_develgo) $(shell \
env GOPATH=$(GOPATH) go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' share/assets/jsmakerule | xargs \
env GOPATH=$(GOPATH) go list -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	@echo '* Prerequisite: bin-jsmakerule'
	go build -o $@ share/assets/jsmakerule
endif

src/share/tmp/jsassets.d: # $(bindir)/jsmakerule
	@echo '* Prerequisite: src/share/tmp/jsassets.d'
#	$(MAKE) $(MFLAGS) $(bindir)/jsmakerule
	$(bindir)/jsmakerule src/share/assets/js/production/ugly/index.js >$@
#	$^ src/share/assets/js/production/ugly/index.js >$@
ifneq ($(MAKECMDGOALS), clean)
include src/share/tmp/jsassets.d
endif
src/share/assets/js/production/ugly/index.js:
	@echo @uglifyjs -c -o $@ ...
	@if type uglifyjs >/dev/null; then cat $^ | uglifyjs -c -o $@ -; fi
#	uglifyjs -c -o $@ $^

src/share/assets/css/index.css: src/share/style/index.scss
	if type sass >/dev/null; then sass $< $@; fi

src/share/assets/js/devel/milk/index.js: src/share/coffee/index.coffee
	if type coffee >/dev/null; then coffee -p $^ >/dev/null && coffee -o $(@D)/ $^; fi

src/share/assets/js/devel/gen/jscript.js: src/share/tmp/jscript.jsx
	if type jsx >/dev/null; then jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@; fi

src/share/templates/%.html: src/share/amber.templates/%.amber src/share/amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines src/share/amber.templates/defines.amber -output $@ $<
src/share/tmp/jscript.jsx: src/share/amber.templates/jscript.amber src/share/amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines src/share/amber.templates/defines.amber -j -output $@ $<

$(bintemplates_productiongo): $(templates_html)
	cd $(<D) && $(gobindata) -pkg templates -tags production -o $(@F) $(^F)
$(bintemplates_develgo): # $(templates_html)
#	$(templates_dir)   instead of $(<D)
#	$(templates_files) instead of $(^F)
	cd $(templates_dir) && $(gobindata) -pkg templates -tags '!production' -debug -o $(@F) $(templates_files)
# 	cd $(dir $(word 1, $(templates_html))) && $(gobindata) -pkg templates -tags '!production' -debug -o ../$(bintemplates_develgo) $(notdir $(templates_html))
ifeq (, $(findstring bootstrap, $(MAKECMDGOALS)))
$(bintemplates_develgo): $(templates_html)
endif

$(binassets_productiongo):
	$(gobindata) -ignore jsmakerule -pkg assets -o $@ -tags production -prefix src/share/assets -ignore src/share/assets/js/devel/ src/share/assets/...
$(binassets_develgo):
	$(gobindata) -ignore jsmakerule -pkg assets -o $@ -tags '!production' -debug -prefix src/share/assets -ignore src/share/assets/js/production/ src/share/assets/...

$(binassets_productiongo): $(shell find src/share/assets -type f \! -name '*.go' \! -path src/share/assets/js/devel/)
$(binassets_productiongo): src/share/assets/css/index.css
$(binassets_productiongo): src/share/assets/js/production/ugly/index.js

ifeq (, $(findstring bootstrap, $(MAKECMDGOALS)))
$(binassets_develgo): $(shell find src/share/assets -type f \! -name '*.go' \! -path src/share/assets/js/production/)
$(binassets_develgo): src/share/assets/css/index.css
$(binassets_develgo): src/share/assets/js/devel/gen/jscript.js
endif
