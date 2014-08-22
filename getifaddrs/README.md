getifaddrs [![godoc getifaddrs](https://godoc.org/github.com/rzab/ostent/getifaddrs?status.svg)](https://godoc.org/github.com/rzab/ostent/getifaddrs)
==========

getifaddrs is a Go package to call getifaddrs(3). cgo is implied. Linux, *BSD supported.

```go
import "github.com/rzab/ostent/src/getifaddrs"

func main() {
ifdatas, err := getifaddrs.getifaddrs()
// ...
}
