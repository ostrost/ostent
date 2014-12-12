// +build production

package assets

import (
	"sync"
	"time"
)

var uncompressedassets struct {
	cache map[string][]byte
	mutex sync.Mutex
}

func UncompressedAssetFunc(readfunc func(string) ([]byte, error)) func(string) ([]byte, error) {
	return func(name string) ([]byte, error) {
		return uncompressedasset(readfunc, name)
	}
}

func uncompressedasset(readfunc func(string) ([]byte, error), name string) ([]byte, error) {
	uncompressedassets.mutex.Lock()
	defer uncompressedassets.mutex.Unlock()
	if text, ok := uncompressedassets.cache[name]; ok {
		return text, nil
	}
	text, err := readfunc(name)
	if err != nil {
		uncompressedassets.cache[name] = text
	}
	return text, err
}

var STARTIME = time.Now()

func ModTime(string, string) (time.Time, error) {
	return STARTIME, nil
}
