package care

import (
	"github.com/labring/lvscare/internal/glog"
	"time"

	"github.com/labring/lvscare/service"
)

//VsAndRsCare is
func (care *LvsCare) VsAndRsCare() {
	lvs := service.BuildLvscare()
	//set inner lvs
	care.lvs = lvs
	if care.Clean {
		glog.V(8).Info("lvscare deleteVirtualServer")
		err := lvs.DeleteVirtualServer(care.VirtualServer, false)
		if err != nil {
			glog.Infof("virtualServer is not exist skip...: %v", err)
		}
	}
	care.createVsAndRs()
	if care.RunOnce {
		return
	}
	t := time.NewTicker(time.Duration(care.Interval) * time.Second)
	for {
		select {
		case <-t.C:
			// in some cases, virtual server maybe removed
			isAvailable := care.lvs.IsVirtualServerAvailable(care.VirtualServer)
			if !isAvailable {
				err := care.lvs.CreateVirtualServer(care.VirtualServer, true)
				//virtual server is exists
				if err != nil {
					glog.Errorf("failed to create virtual server: %v", err)

					return
				}
			}
			//check real server
			lvs.CheckRealServers(care.HealthPath, care.HealthSchem)
		}
	}
}
