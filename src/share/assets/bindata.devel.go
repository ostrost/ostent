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
		"src/share/assets/css/index.css",
		"css/index.css",
	)
}

// css_index_css_map reads file data from disk. It returns an error on failure.
func css_index_css_map() ([]byte, error) {
	return bindata_read(
		"src/share/assets/css/index.css.map",
		"css/index.css.map",
	)
}

// favicon_png reads file data from disk. It returns an error on failure.
func favicon_png() ([]byte, error) {
	return bindata_read(
		"src/share/assets/favicon.png",
		"favicon.png",
	)
}

// js_devel_gen_jscript_js reads file data from disk. It returns an error on failure.
func js_devel_gen_jscript_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/gen/jscript.js",
		"js/devel/gen/jscript.js",
	)
}

// js_devel_milk_index_js reads file data from disk. It returns an error on failure.
func js_devel_milk_index_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/milk/index.js",
		"js/devel/milk/index.js",
	)
}

// js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js",
		"js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js",
	)
}

// js_devel_vendor_min_headroom_0_5_0_headroom_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_headroom_0_5_0_headroom_min_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/vendor/min/headroom/0.5.0/headroom.min.js",
		"js/devel/vendor/min/headroom/0.5.0/headroom.min.js",
	)
}

// js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js",
		"js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js",
	)
}

// js_devel_vendor_min_react_0_11_0_react_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_react_0_11_0_react_min_js() ([]byte, error) {
	return bindata_read(
		"src/share/assets/js/devel/vendor/min/react/0.11.0/react.min.js",
		"js/devel/vendor/min/react/0.11.0/react.min.js",
	)
}

// robots_txt reads file data from disk. It returns an error on failure.
func robots_txt() ([]byte, error) {
	return bindata_read(
		"src/share/assets/robots.txt",
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
	"css/index.css":                                        css_index_css,
	"css/index.css.map":                                    css_index_css_map,
	"favicon.png":                                          favicon_png,
	"js/devel/gen/jscript.js":                              js_devel_gen_jscript_js,
	"js/devel/milk/index.js":                               js_devel_milk_index_js,
	"js/devel/vendor/min/bootstrap/3.2.0/bootstrap.min.js": js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js,
	"js/devel/vendor/min/headroom/0.5.0/headroom.min.js":   js_devel_vendor_min_headroom_0_5_0_headroom_min_js,
	"js/devel/vendor/min/jquery/2.1.1/jquery-2.1.1.min.js": js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js,
	"js/devel/vendor/min/react/0.11.0/react.min.js":        js_devel_vendor_min_react_0_11_0_react_min_js,
	"robots.txt":                                           robots_txt,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"favicon.png": &_bintree_t{favicon_png, map[string]*_bintree_t{}},
	"js": &_bintree_t{nil, map[string]*_bintree_t{
		"devel": &_bintree_t{nil, map[string]*_bintree_t{
			"gen": &_bintree_t{nil, map[string]*_bintree_t{
				"jscript.js": &_bintree_t{js_devel_gen_jscript_js, map[string]*_bintree_t{}},
			}},
			"milk": &_bintree_t{nil, map[string]*_bintree_t{
				"index.js": &_bintree_t{js_devel_milk_index_js, map[string]*_bintree_t{}},
			}},
			"vendor": &_bintree_t{nil, map[string]*_bintree_t{
				"min": &_bintree_t{nil, map[string]*_bintree_t{
					"headroom": &_bintree_t{nil, map[string]*_bintree_t{
						"0.5.0": &_bintree_t{nil, map[string]*_bintree_t{
							"headroom.min.js": &_bintree_t{js_devel_vendor_min_headroom_0_5_0_headroom_min_js, map[string]*_bintree_t{}},
						}},
					}},
					"jquery": &_bintree_t{nil, map[string]*_bintree_t{
						"2.1.1": &_bintree_t{nil, map[string]*_bintree_t{
							"jquery-2.1.1.min.js": &_bintree_t{js_devel_vendor_min_jquery_2_1_1_jquery_2_1_1_min_js, map[string]*_bintree_t{}},
						}},
					}},
					"react": &_bintree_t{nil, map[string]*_bintree_t{
						"0.11.0": &_bintree_t{nil, map[string]*_bintree_t{
							"react.min.js": &_bintree_t{js_devel_vendor_min_react_0_11_0_react_min_js, map[string]*_bintree_t{}},
						}},
					}},
					"bootstrap": &_bintree_t{nil, map[string]*_bintree_t{
						"3.2.0": &_bintree_t{nil, map[string]*_bintree_t{
							"bootstrap.min.js": &_bintree_t{js_devel_vendor_min_bootstrap_3_2_0_bootstrap_min_js, map[string]*_bintree_t{}},
						}},
					}},
				}},
			}},
		}},
	}},
	"robots.txt": &_bintree_t{robots_txt, map[string]*_bintree_t{}},
	"css": &_bintree_t{nil, map[string]*_bintree_t{
		"index.css.map": &_bintree_t{css_index_css_map, map[string]*_bintree_t{}},
		"index.css":     &_bintree_t{css_index_css, map[string]*_bintree_t{}},
	}},
}}
