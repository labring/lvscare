package version

import (
	"fmt"
	"runtime"
)

var (
	Version  = "v1.0.0"
	Build    = ""
	PrintStr = fmt.Sprintf("lvscare version %v, build %v %v", Version, Build, runtime.Version())
)
