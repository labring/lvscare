package care

import (
	"github.com/wonderivan/logger"
	"time"

	"github.com/fanux/lvscare/create"
	"github.com/fanux/lvscare/service"
)

//VsAndRsCare is
func VsAndRsCare(vs string, rs []string, beat int64, path string, schem string) {
	lvs := service.BuildLvscare()
	logger.Debug("VsAndRsCare DeleteVirtualServer")
	err := lvs.DeleteVirtualServer(vs, false)
	print(err)
	t := time.NewTicker(time.Duration(beat) * time.Second)
	for {
		select {
		case <-t.C:
			//check virturl server
			isAvailable := lvs.IsVirtualServerAvailable(vs)
			if !isAvailable {
				create.VsAndRsCreate(vs, rs, lvs)
			}
			//check real server
			lvs.CheckRealServers(path, schem)
		}
	}

}
