package care

import (
	"time"

	"github.com/fanux/LVScare/create"
	"github.com/fanux/LVScare/service"
)

//VsAndRsCare is
func VsAndRsCare(vs string, rs []string, beat int64, path string, schem string) error {
	lvs := service.BuildLvscare()
	t := time.NewTicker(time.Duration(beat) * time.Second)
	for {
		select {
		case <-t.C:
			//check virturl server
			service, _ := lvs.GetVirtualServer()
			if service == nil {
				create.VsAndRsCreate(vs, rs)
			}

			//check real server
			lvs.CheckRealServers(path, schem)
		}
	}

	return nil
}
