#!/usr/bin/env make -f

bindir=bin/$(shell uname -sm | awk '{ sub(/x86_64/, "amd64", $$2); print tolower($$1) "_" $$2; }')
templates_html=$(shell echo templates.html/{index,usepercent,tooltipable}.html)

.PHONY: all devel
all: $(bindir)/ostent
devel: $(shell echo src/ostential/{view/bindata.devel.go,assets/bindata.devel.go})

$(bindir)/amberpp: $(shell go list -f '\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}' amberp/amberpp | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	go build -o $@ amberp/amberpp

$(bindir)/ostent: $(shell \
go list -tags production -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostent | xargs \
go list -tags production -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	go build -tags production -o $@ ostent

$(bindir)/jsmakerule: $(shell \
go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostential/assets/jsmakerule | xargs \
go list -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	@echo '* Prerequisite: bin-jsmakerule'
	go build -o $@ ostential/assets/jsmakerule

tmp/jsassets.d: $(bindir)/jsmakerule
	@echo '* Prerequisite: tmp/jsassets.d'
	$^ assets/js/production/ugly/index.js >$@
include tmp/jsassets.d
assets/js/production/ugly/index.js:
	@echo    @uglifyjs -c -o $@ [devel-jsassets]
	@cat $^ | uglifyjs -c -o $@ -
#	uglifyjs -c -o $@ $^

assets/css/index.css: style/index.scss
	sass $< $@

templates.html/%.html: amber.templates/%.amber amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines amber.templates/defines.amber -output $@ $<
tmp/jscript.jsx: amber.templates/jscript.amber amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines amber.templates/defines.amber -j -output $@ $<

assets/js/devel/gen/jscript.js: tmp/jscript.jsx
	jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@

src/ostential/view/bindata.production.go: $(templates_html) # $(wildcard templates.html/*.html)
	go-bindata -pkg view -o $@ -tags production -prefix templates.html templates.html/...
src/ostential/view/bindata.devel.go: $(templates_html) # $(wildcard templates.html/*.html)
	go-bindata -pkg view -o $@ -tags '!production' -debug -prefix templates.html templates.html/...

src/ostential/assets/bindata.production.go: assets/css/index.css $(shell find assets -type f | grep -v assets/js/devel/) assets/js/production/ugly/index.js
	go-bindata -pkg assets -o $@ -tags production -prefix assets -ignore assets/js/devel/ assets/...
src/ostential/assets/bindata.devel.go: assets/css/index.css $(shell find assets -type f | grep -v assets/js/production/) assets/js/devel/gen/jscript.js
	go-bindata -pkg assets -o $@ -tags '!production' -debug -prefix assets -ignore assets/js/production/ assets/...
