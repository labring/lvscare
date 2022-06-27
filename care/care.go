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
		rIP, _, err := net.SplitHostPort(LVS.RealServer[0])
		if err != nil {
			return err
		}
		vIP, _, err := net.SplitHostPort(LVS.VirtualServer)
		if err != nil {
			return err
		}
		if utils.IsIpv4(rIP) {
			ipv4 = true
		} else {
			ipv4 = false
		}
		if !ipv4 {
			return nil
		}
		defaultInterface, err := utils.GetLocalDefaultRoute(rIP)
		if err != nil {
			return err
		}
		if len(defaultInterface.Addr) > 0 {
			for _, addr := range defaultInterface.Addr {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return err
				}
				if (!utils.IsIPv6(ip) && ipv4) || (utils.IsIPv6(ip) && !ipv4) {
					LVS.TargetIP = ip
				}
			}
		}
		LVS.Route = route.NewRoute(vIP, LVS.TargetIP.String())
		return LVS.Route.SetRoute()
	}
	return errors.New("real server can't empty")
}
