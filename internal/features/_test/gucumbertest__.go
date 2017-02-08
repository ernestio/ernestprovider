
package main

import (
	"github.com/gucumber/gucumber"
	_i0 "github.com/ernestio/ernestprovider/internal/features/real"
	
)

var (
	_ci0 = _i0.IMPORT_MARKER
	
)

func main() {
	
	gucumber.GlobalContext.Filters = []string{
	"@virtual_network",
	
	}
	
	gucumber.GlobalContext.RunDir("internal/features")
}
