#!/usr/bin/env make -f

bindir=bin/$(shell uname -sm | awk '{ sub(/x86_64/, "amd64", $$2); print tolower($$1) "_" $$2; }')
templates_html=$(shell echo templates.html/{index,usepercent,tooltipable}.html)

.PHONY: all devel
all: $(bindir)/ostent
devel: # $(shell echo src/ostential/{view,assets}/bindata.devel.go)
	go-bindata -pkg assets -o src/ostential/assets/bindata.devel.go -tags '!production' -debug -prefix assets -ignore assets/js/production/ assets/...
	go-bindata -pkg view   -o src/ostential/view/bindata.devel.go   -tags '!production' -debug -prefix templates.html templates.html/...

%: %.sh # clear the implicit *.sh rule covering ./ostent.sh

$(bindir)/%:
	@echo '* Sources:' $^
	go build -o $@ $(patsubst src////%,%,$|)

$(bindir)/amberpp: | src////amberp/amberpp
$(bindir)/ostent:  | src////ostent

$(bindir)/amberpp: $(shell go list -f '\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}' amberp/amberpp | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr

$(bindir)/ostent: $(shell \
go list -tags production -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostent | xargs \
go list -tags production -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	go build -tags production -o $@ ostent

$(bindir)/jsmakerule: src/ostential/assets/bindata.devel.go $(shell \
go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostential/assets/jsmakerule | xargs \
go list -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	@echo '* Prerequisite: bin-jsmakerule'
	go build -o $@ ostential/assets/jsmakerule

tmp/jsassets.d: # $(bindir)/jsmakerule
	@echo '* Prerequisite: tmp/jsassets.d'
#	$(MAKE) $(MFLAGS) $(bindir)/jsmakerule
	$(bindir)/jsmakerule assets/js/production/ugly/index.js >$@
#	$^ assets/js/production/ugly/index.js >$@
ifneq ($(MAKECMDGOALS), clean)
include tmp/jsassets.d
endif
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

assets/js/devel/milk/index.js: coffee/index.coffee
	coffee -p $^ >/dev/null && coffee -o $(@D)/ $^

assets/js/devel/gen/jscript.js: tmp/jscript.jsx
	jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@

src/ostential/view/bindata.production.go: $(templates_html)
	cd $(<D) && go-bindata -pkg view -tags production -o ../$@ $(^F)
src/ostential/view/bindata.devel.go: $(templates_html)
	cd $(<D) && go-bindata -pkg view -tags '!production' -debug -o ../$@ $(^F)

src/ostential/assets/bindata.production.go: assets/css/index.css $(shell find assets -type f | grep -v assets/js/devel/) assets/js/production/ugly/index.js
	go-bindata -pkg assets -o $@ -tags production -prefix assets -ignore assets/js/devel/ assets/...
src/ostential/assets/bindata.devel.go: assets/css/index.css $(shell find assets -type f | grep -v assets/js/production/) assets/js/devel/gen/jscript.js
	go-bindata -pkg assets -o $@ -tags '!production' -debug -prefix assets -ignore assets/js/production/ assets/...
