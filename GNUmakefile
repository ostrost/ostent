#!/usr/bin/env make -f

templates_html           =$(shell echo share/templates.html/{index,usepercent,tooltipable}.html)
binassets_develgo        =src/ostential/assets/bindata.devel.go
binassets_productiongo   =src/ostential/assets/bindata.production.go
bintemplates_develgo     =src/ostential/view/bindata.devel.go
bintemplates_productiongo=src/ostential/view/bindata.production.go

bindir=bin/$(shell uname -sm | awk '{ sub(/x86_64/, "amd64", $$2); print tolower($$1) "_" $$2; }')

.PHONY: all devel
all: $(bindir)/ostent
devel: # $(shell echo src/ostential/{view,assets}/bindata.devel.go)
	go-bindata -pkg assets -o $(binassets_develgo) -tags '!production' -debug -prefix share/assets -ignore share/assets/js/production/ share/assets/...
	cd $(dir $(word 1, $(templates_html))) && go-bindata -pkg view -tags '!production' -debug -o ../$(bintemplates_develgo) $(notdir $(templates_html))

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

$(bindir)/jsmakerule: $(binassets_develgo) $(shell \
go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' ostential/assets/jsmakerule | xargs \
go list -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles     }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles    }}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr
#	@echo '* Sources:' $^
	@echo '* Prerequisite: bin-jsmakerule'
	go build -o $@ ostential/assets/jsmakerule

share/tmp/jsassets.d: # $(bindir)/jsmakerule
	@echo '* Prerequisite: share/tmp/jsassets.d'
#	$(MAKE) $(MFLAGS) $(bindir)/jsmakerule
	$(bindir)/jsmakerule share/assets/js/production/ugly/index.js >$@
#	$^ share/assets/js/production/ugly/index.js >$@
ifneq ($(MAKECMDGOALS), clean)
include share/tmp/jsassets.d
endif
share/assets/js/production/ugly/index.js:
	@echo    @uglifyjs -c -o $@ [devel-jsassets]
	@cat $^ | uglifyjs -c -o $@ -
#	uglifyjs -c -o $@ $^

share/assets/css/index.css: share/style/index.scss
	sass $< $@

share/templates.html/%.html: share/amber.templates/%.amber share/amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines share/amber.templates/defines.amber -output $@ $<
share/tmp/jscript.jsx: share/amber.templates/jscript.amber share/amber.templates/defines.amber $(bindir)/amberpp
	$(bindir)/amberpp -defines share/amber.templates/defines.amber -j -output $@ $<

share/assets/js/devel/milk/index.js: share/coffee/index.coffee
	coffee -p $^ >/dev/null && coffee -o $(@D)/ $^

share/assets/js/devel/gen/jscript.js: share/tmp/jscript.jsx
	jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@

$(bintemplates_productiongo): $(templates_html)
	cd $(<D) && go-bindata -ignore '.*\.go' -pkg view -tags production -o ../../$@ $(^F)
$(bintemplates_develgo): $(templates_html)
	cd $(<D) && go-bindata -ignore '.*\.go' -pkg view -tags '!production' -debug -o ../../$@ $(^F)

$(binassets_productiongo): share/assets/css/index.css $(shell find share/assets -type f | grep -v share/assets/js/devel/) share/assets/js/production/ugly/index.js
	go-bindata -ignore '.*\.go' -ignore jsmakerule -pkg assets -o $@ -tags production -prefix share/assets -ignore share/assets/js/devel/ share/assets/...
$(binassets_develgo): share/assets/css/index.css $(shell find share/assets -type f | grep -v share/assets/js/production/) share/assets/js/devel/gen/jscript.js
	go-bindata -ignore '.*\.go' -ignore jsmakerule -pkg assets -o $@ -tags '!production' -debug -prefix share/assets -ignore share/assets/js/production/ share/assets/...
