// +build !production

package assets

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path"
	"path/filepath"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

// css_index_css reads file data from disk. It returns an error on failure.
func css_index_css() (*asset, error) {
	path := filepath.Join(rootDir, "css/index.css")
	name := "css/index.css"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// css_index_css_map reads file data from disk. It returns an error on failure.
func css_index_css_map() (*asset, error) {
	path := filepath.Join(rootDir, "css/index.css.map")
	name := "css/index.css.map"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// favicon_png reads file data from disk. It returns an error on failure.
func favicon_png() (*asset, error) {
	path := filepath.Join(rootDir, "favicon.png")
	name := "favicon.png"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_gen_jscript_js reads file data from disk. It returns an error on failure.
func js_devel_gen_jscript_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/gen/jscript.js")
	name := "js/devel/gen/jscript.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_milk_index_js reads file data from disk. It returns an error on failure.
func js_devel_milk_index_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/milk/index.js")
	name := "js/devel/milk/index.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_vendor_min_bootstrap_3_3_2_bootstrap_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_bootstrap_3_3_2_bootstrap_min_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/vendor/min/bootstrap/3.3.2/bootstrap.min.js")
	name := "js/devel/vendor/min/bootstrap/3.3.2/bootstrap.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_vendor_min_headroom_0_7_0_headroom_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_headroom_0_7_0_headroom_min_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/vendor/min/headroom/0.7.0/headroom.min.js")
	name := "js/devel/vendor/min/headroom/0.7.0/headroom.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_vendor_min_jquery_2_1_3_jquery_2_1_3_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_jquery_2_1_3_jquery_2_1_3_min_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/vendor/min/jquery/2.1.3/jquery-2.1.3.min.js")
	name := "js/devel/vendor/min/jquery/2.1.3/jquery-2.1.3.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// js_devel_vendor_min_react_0_12_2_react_min_js reads file data from disk. It returns an error on failure.
func js_devel_vendor_min_react_0_12_2_react_min_js() (*asset, error) {
	path := filepath.Join(rootDir, "js/devel/vendor/min/react/0.12.2/react.min.js")
	name := "js/devel/vendor/min/react/0.12.2/react.min.js"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// robots_txt reads file data from disk. It returns an error on failure.
func robots_txt() (*asset, error) {
	path := filepath.Join(rootDir, "robots.txt")
	name := "robots.txt"
	bytes, err := bindata_read(path, name)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}

	a := &asset{bytes: bytes, info: fi}
	return a, err
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"css/index.css": css_index_css,
	"css/index.css.map": css_index_css_map,
	"favicon.png": favicon_png,
	"js/devel/gen/jscript.js": js_devel_gen_jscript_js,
	"js/devel/milk/index.js": js_devel_milk_index_js,
	"js/devel/vendor/min/bootstrap/3.3.2/bootstrap.min.js": js_devel_vendor_min_bootstrap_3_3_2_bootstrap_min_js,
	"js/devel/vendor/min/headroom/0.7.0/headroom.min.js": js_devel_vendor_min_headroom_0_7_0_headroom_min_js,
	"js/devel/vendor/min/jquery/2.1.3/jquery-2.1.3.min.js": js_devel_vendor_min_jquery_2_1_3_jquery_2_1_3_min_js,
	"js/devel/vendor/min/react/0.12.2/react.min.js": js_devel_vendor_min_react_0_12_2_react_min_js,
	"robots.txt": robots_txt,
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
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"css": &_bintree_t{nil, map[string]*_bintree_t{
		"index.css": &_bintree_t{css_index_css, map[string]*_bintree_t{
		}},
		"index.css.map": &_bintree_t{css_index_css_map, map[string]*_bintree_t{
		}},
	}},
	"favicon.png": &_bintree_t{favicon_png, map[string]*_bintree_t{
	}},
	"js": &_bintree_t{nil, map[string]*_bintree_t{
		"devel": &_bintree_t{nil, map[string]*_bintree_t{
			"gen": &_bintree_t{nil, map[string]*_bintree_t{
				"jscript.js": &_bintree_t{js_devel_gen_jscript_js, map[string]*_bintree_t{
				}},
			}},
			"milk": &_bintree_t{nil, map[string]*_bintree_t{
				"index.js": &_bintree_t{js_devel_milk_index_js, map[string]*_bintree_t{
				}},
			}},
			"vendor": &_bintree_t{nil, map[string]*_bintree_t{
				"min": &_bintree_t{nil, map[string]*_bintree_t{
					"bootstrap": &_bintree_t{nil, map[string]*_bintree_t{
						"3.3.2": &_bintree_t{nil, map[string]*_bintree_t{
							"bootstrap.min.js": &_bintree_t{js_devel_vendor_min_bootstrap_3_3_2_bootstrap_min_js, map[string]*_bintree_t{
							}},
						}},
					}},
					"headroom": &_bintree_t{nil, map[string]*_bintree_t{
						"0.7.0": &_bintree_t{nil, map[string]*_bintree_t{
							"headroom.min.js": &_bintree_t{js_devel_vendor_min_headroom_0_7_0_headroom_min_js, map[string]*_bintree_t{
							}},
						}},
					}},
					"jquery": &_bintree_t{nil, map[string]*_bintree_t{
						"2.1.3": &_bintree_t{nil, map[string]*_bintree_t{
							"jquery-2.1.3.min.js": &_bintree_t{js_devel_vendor_min_jquery_2_1_3_jquery_2_1_3_min_js, map[string]*_bintree_t{
							}},
						}},
					}},
					"react": &_bintree_t{nil, map[string]*_bintree_t{
						"0.12.2": &_bintree_t{nil, map[string]*_bintree_t{
							"react.min.js": &_bintree_t{js_devel_vendor_min_react_0_12_2_react_min_js, map[string]*_bintree_t{
							}},
						}},
					}},
				}},
			}},
		}},
	}},
	"robots.txt": &_bintree_t{robots_txt, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
