// +build production

package assets

import (
	"os"
	"sync"
	"time"
)

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// ReadCache is a cache in AssetReadFunc.
	ReadCache struct {
		MU     sync.Mutex
		Byname map[string][]byte
	}
)

// ProductionAssetInfoFunc wraps bindata's AssetInfo func. ModTime is always BootTime.
func ProductionAssetInfoFunc(infofunc func(string) (os.FileInfo, error)) func(string) (TimeInfo, error) {
	return func(name string) (TimeInfo, error) {
		_, err := infofunc(name)
		if err != nil {
			return nil, err
		}
		return BootInfo{}, nil
	}
}

// ProductionAssetReadFunc wraps bindata's Asset func. Result is from cache or cached.
func ProductionAssetReadFunc(readfunc func(string) ([]byte, error)) func(string) ([]byte, error) {
	return func(name string) ([]byte, error) {
		return Read(readfunc, name)
	}
}

// BootInfo is a TimeInfo implementation.
type BootInfo struct{}

// ModTime returns BootTime.
func (bi BootInfo) ModTime() time.Time { return BootTime } // bi is unused

// Read returns cached readfunc result.
func Read(readfunc func(string) ([]byte, error), name string) ([]byte, error) {
	ReadCache.MU.Lock()
	defer ReadCache.MU.Unlock()
	if text, ok := ReadCache.Byname[name]; ok {
		return text, nil
	}
	text, err := readfunc(name)
	if err != nil {
		ReadCache.Byname[name] = text
	}
	return text, err
}
