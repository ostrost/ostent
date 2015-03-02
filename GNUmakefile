#!/usr/bin/env make -f

fqostent=github.com/ostrost/ostent

binassets_develgo         = share/assets/bindata.devel.go
binassets_productiongo    = share/assets/bindata.production.go
bintemplates_develgo      = share/templates/bindata.devel.go
bintemplates_productiongo = share/templates/bindata.production.go
templates_dir             = share/templates/
templates_files           = index.html usepercent.html tooltipable.html
templates_html=$(addprefix $(templates_dir), $(templates_files))

PATH=$(shell printf %s: $$PATH; echo $$GOPATH | awk -F: 'BEGIN { OFS="/bin:"; } { print $$1,$$2,$$3,$$4,$$5,$$6,$$7,$$8,$$9 "/bin"}')

xargs=xargs
sed-i=sed -i ''
ifeq (Linux, $(shell uname -s))
xargs=xargs --no-run-if-empty
sed-i=sed -i'' # GNU sed -i opt, not a flag
endif
sed-i-production-bindata=$(sed-i) -Ee 's/time\.Unix\([0-9]+,/time.Unix(1400000000,/g'
# -e '/^\/\/ AssetDir /,$$d' # applicable for devel assets and both templates \
# but then there're `imported and not used: "{path/filepath,io/ioutil}"`
go-bindata=go-bindata -ignore '.*\.go' # Go regexp syntax for -ignore

.PHONY: all al init test covertest coverfunc coverhtml bindata bindata-devel bindata-production
ifneq (init, $(MAKECMDGOALS))
# before init:
# - go list would fail => unknown $(destbin)
# - go test fails without dependencies installed
# - go-bindata is not installed yet

destbin=$(shell echo $(GOPATH) | awk -F: '{ print $$1 "/bin" }')
# destbin=$(abspath $(dir $(shell go list -f '{{.Target}}' $(fqostent))))
ostent_files=$(shell \
go list -tags production -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(fqostent) | xargs \
go list -tags production -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles}}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr

all: $(destbin)/ostent
endif
init:
	go get -u -v \
github.com/jteeuwen/go-bindata/go-bindata \
github.com/skelterjohn/rerun \
github.com/clipperhouse/gen
# TODO rm {src,pkg/*}/github.com/clipperhouse/slice{,.a}
	cd types && gen add github.com/rzab/slice
	git remote set-url origin https://$(fqostent) # travis & tip & https://code.google.com/p/go/issues/detail?id=8850
	go get -v -tags production $(fqostent)
	go list -f '{{.Target}}' $(fqostent) | $(xargs) rm # clean the library archive
	go get -v -a $(fqostent)
	go list -f '{{.Target}}' $(fqostent) | $(xargs) rm # clean the library archive

%: %.sh # clear the implicit *.sh rule covering ./ostent.sh

ifneq (init, $(MAKECMDGOALS))
test:
	go vet ./...
	go test -v ./...
covertest:
	go test -v -covermode=count -coverprofile=coverage.out $(fqostent)
coverfunc:
	go tool cover -func=coverage.out
coverhtml:
	go tool cover -html=coverage.out

$(PWD)/assetutil/%_slice.go: $(PWD)/assetutil/assetutil.go
	cd $(dir $@) && go generate
$(PWD)/types/%_slice.go: $(PWD)/types/types.go
	cd $(dir $@) && go generate

al: $(ostent_files)
# al: like `all' but without final go build ostent. For when rerun does the build

$(destbin)/ostent: $(ostent_files)
	go build -ldflags -w -a -tags production -o $@ $(fqostent)

$(destbin)/%:
	go build -o $@ $(fqostent)/$|
$(destbin)/amberpp:    | amberp/amberpp
$(destbin)/jsmakerule: | share/assets/jsmakerule

$(destbin)/amberpp: $(shell go list -f '\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}' $(fqostent)/amberp/amberpp | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr

$(destbin)/jsmakerule: $(binassets_develgo)
$(destbin)/jsmakerule: $(shell \
go list -f '{{.ImportPath}}{{"\n"}}{{join .Deps "\n"}}' $(fqostent)/share/assets/jsmakerule | xargs \
go list -f '{{if and (not .Standard) (not .Goroot)}}\
{{$$dir := .Dir}}\
{{range .GoFiles }}{{$$dir}}/{{.}}{{"\n"}}{{end}}\
{{range .CgoFiles}}{{$$dir}}/{{.}}{{"\n"}}{{end}}{{end}}' | \
sed -n "s,^ *,,g; s,$(PWD)/,,p" | sort) # | tee /dev/stderr

share/tmp/jsassets.d: # $(destbin)/jsmakerule
#	$(MAKE) $(MFLAGS) $(destbin)/jsmakerule
	true && \
$(destbin)/jsmakerule share/assets/js/production/ugly/index.js >/dev/null && \
$(destbin)/jsmakerule share/assets/js/production/ugly/index.js >$@
#	$^ share/assets/js/production/ugly/index.js >$@
ifneq ($(MAKECMDGOALS), clean)
include share/tmp/jsassets.d
endif

# these four rules are actually independant of $(destbin) and could be set when the goal is `init', but we're keeping it simple
share/assets/js/production/ugly/index.js: # the prerequisites from included jsassets.d
	if type uglifyjs >/dev/null; then cat /dev/null $^ | uglifyjs -c -o $@ -; fi
share/assets/css/index.css: share/style/index.scss
	if type sass >/dev/null; then sass $< $@; fi
share/assets/js/devel/milk/index.js: share/coffee/index.coffee
	if type coffee >/dev/null; then coffee -p $^ >/dev/null && coffee -o $(@D)/ $^; fi
share/assets/js/devel/gen/jscript.js: share/tmp/jscript.jsx
	if type jsx >/dev/null; then jsx <$^ >/dev/null && jsx <$^ 2>/dev/null >$@; fi

share/templates/%.html: share/amber.templates/%.amber share/amber.templates/defines.amber $(destbin)/amberpp
	$(destbin)/amberpp -defines share/amber.templates/defines.amber -output $@ $<
share/tmp/jscript.jsx: share/amber.templates/jscript.amber share/amber.templates/defines.amber $(destbin)/amberpp
	$(destbin)/amberpp -defines share/amber.templates/defines.amber -j -output $@ $<

$(bintemplates_productiongo): $(templates_html)
	cd $(<D) && $(go-bindata) -pkg templates -tags production -o $(@F) $(^F) && $(sed-i-production-bindata) $(@F)
$(bintemplates_develgo): $(templates_html)
	cd $(templates_dir) && $(go-bindata) -pkg templates -tags '!production' -dev -o $(@F) $(templates_files)
#	# the target has no prerequisites e.g. $(templates_html):
#	# $(templates_dir)   instead of $(<D)
#	# $(templates_files) instead of $(^F)

$(binassets_productiongo):
	cd share/assets && $(go-bindata) -pkg assets -o $(@F) -tags production -ignore js/devel/ ./... && $(sed-i-production-bindata) $(@F)
$(binassets_develgo):
	cd share/assets && $(go-bindata) -pkg assets -o $(@F) -tags '!production' -dev -ignore js/production/ ./...

$(binassets_productiongo): $(shell find \
                           share/assets -type f \! -name '*.go' \! -path \
                          'share/assets/js/devel/*')
$(binassets_productiongo): share/assets/css/index.css
$(binassets_productiongo): share/assets/js/production/ugly/index.js

$(binassets_develgo): $(shell find \
                      share/assets -type f \! -name '*.go' \! -path \
                     'share/assets/js/production/*')
$(binassets_develgo): share/assets/css/index.css
$(binassets_develgo): share/assets/js/devel/gen/jscript.js

# spare shortcuts
bindata-production: $(binassets_productiongo) $(bintemplates_productiongo)
bindata-devel: $(binassets_develgo) $(bintemplates_develgo)
bindata: bindata-devel bindata-production

endif # END OF ifneq (init, $(MAKECMDGOALS))
