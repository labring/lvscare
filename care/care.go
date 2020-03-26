package care

import (
	"github.com/wonderivan/logger"
	"time"

	"github.com/fanux/lvscare/service"
)

//VsAndRsCare is
func VsAndRsCare() {
	lvs := service.BuildLvscare()
	logger.Debug("VsAndRsCare DeleteVirtualServer")
	err := lvs.DeleteVirtualServer(VirtualServer, false)
	logger.Warn("VsAndRsCare DeleteVirtualServer:", err)
	if RunOnce {
		createVsAndRs(VirtualServer, RealServer, lvs)
		return
	}
	t := time.NewTicker(time.Duration(Interval) * time.Second)
	for {
		select {
		case <-t.C:
			createVsAndRs(VirtualServer, RealServer, lvs)
			//check real server
			lvs.CheckRealServers(HealthPath, HealthSchem)
		}
	}
}
