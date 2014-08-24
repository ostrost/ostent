getifaddrs [![godoc getifaddrs](https://godoc.org/github.com/ostrost/ostent/getifaddrs?status.svg)](https://godoc.org/github.com/ostrost/ostent/getifaddrs)
==========

getifaddrs is a Go package to call getifaddrs(3). cgo is implied. Linux, *BSD supported.

```go
import "github.com/ostrost/ostent/src/getifaddrs"

func main() {
ifdatas, err := getifaddrs.getifaddrs()
// ...
}
