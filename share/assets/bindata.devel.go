// +build !production

package assets

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// css_index_css reads file data from disk. It returns an error on failure.
func css_index_css() ([]byte, error) {
	return bindata_read(
		"share/assets/css/index.css",
		"css/index.css",
	)
}

// css_index_css_map reads file data from disk. It returns an error on failure.
func css_index_css_map() ([]byte, error) {
	return bindata_read(
		"share/assets/css/index.css.map",
		"css/index.css.map",
	)
}

// favicon_png reads file data from disk. It returns an error on failure.
func favicon_png() ([]byte, error) {
	return bindata_read(
		"share/assets/favicon.png",
		"favicon.png",
	)
}

// js_devel_gen_jscript_js reads file data from disk. It returns an error on failure.
func js_devel_gen_jscript_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/gen/jscript.js",
		"js/devel/gen/jscript.js",
	)
}

// js_devel_milk_index_js reads file data from disk. It returns an error on failure.
func js_devel_milk_index_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/milk/index.js",
		"js/devel/milk/index.js",
	)
}

// js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js",
		"js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js",
	)
}

// js_devel_vendor_min_headroom_0_7_0_headroom_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_headroom_0_7_0_headroom_min_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/vendor/min/headroom/0.7.0/headroom.min.js",
		"js/devel/vendor/min/headroom/0.7.0/headroom.min.js",
	)
}

// js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js",
		"js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js",
	)
}

// js_devel_vendor_min_react_0_11_2_react_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_react_0_11_2_react_min_js() ([]byte, error) {
	return bindata_read(
		"share/assets/js/devel/vendor/min/react/0.11.2/react.min.js",
		"js/devel/vendor/min/react/0.11.2/react.min.js",
	)
}

// robots_txt reads file data from disk. It returns an error on failure.
func robots_txt() ([]byte, error) {
	return bindata_read(
		"share/assets/robots.txt",
		"robots.txt",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"css/index.css": css_index_css,
	"css/index.css.map": css_index_css_map,
	"favicon.png": favicon_png,
	"js/devel/gen/jscript.js": js_devel_gen_jscript_js,
	"js/devel/milk/index.js": js_devel_milk_index_js,
	"js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js": js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js,
	"js/devel/vendor/min/headroom/0.7.0/headroom.min.js": js_devel_vendor_min_headroom_0_7_0_headroom_min_js,
	"js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js": js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js,
	"js/devel/vendor/min/react/0.11.2/react.min.js": js_devel_vendor_min_react_0_11_2_react_min_js,
	"robots.txt": robots_txt,
}
