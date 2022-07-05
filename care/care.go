package care

import (
	"errors"
	"github.com/labring/lvscare/internal/glog"
	"github.com/labring/lvscare/internal/route"
	"github.com/labring/lvscare/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
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
		case signa := <-sig:
			glog.Infof("receive kill signal: %+v", signa)
			_ = LVS.Route.DelRoute()
			return
		}
	}
}
func (care *LvsCare) SyncRouter() error {
	if len(LVS.VirtualServer) == 0 {
		return errors.New("virtual server can't empty")
	}
	if len(LVS.RealServer) > 0 {
		var ipv4 bool
		vIP, _, err := net.SplitHostPort(LVS.VirtualServer)
		if err != nil {
			return err
		}
		if utils.IsIPv6(LVS.TargetIP) {
			ipv4 = false
		} else {
			ipv4 = true
		}
		if !ipv4 {
			glog.Infof("tip: %s is not ipv4", LVS.TargetIP.String())
			return nil
		}
		glog.Infof("tip: %s,vip: %s", LVS.TargetIP.String(), vIP)
		LVS.Route = route.NewRoute(vIP, LVS.TargetIP.String())
		return LVS.Route.SetRoute()
	}
	return errors.New("real server can't empty")
}

func SetTargetIP() error {
	if LVS.TargetIP == nil {
		LVS.TargetIP = net.ParseIP(os.Getenv("LVSCARE_NODE_IP"))
	}

	return nil
}
