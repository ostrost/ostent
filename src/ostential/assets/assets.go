package assets
func BindataKeys() (keys []string) {
	for k := range _bindata {
		keys = append(keys, k)
	}
	return
}
