package care

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labring/lvscare/internal/glog"
	"github.com/labring/lvscare/internal/route"
	"github.com/labring/lvscare/service"
	"github.com/labring/lvscare/utils"
)

// VsAndRsCare sync route and set ipvs rules
func (care *LvsCare) VsAndRsCare() error {
	lvs := service.BuildLvscare()
	// set inner lvs
	care.lvs = lvs
	if care.Clean {
		glog.V(8).Info("lvscare deleteVirtualServer")
		err := lvs.DeleteVirtualServer(care.VirtualServer, false)
		if err != nil {
			glog.Infof("virtualServer is not exist skip...: %v", err)
		}
	}
	if err := care.test(); err != nil {
		return err
	}
	if err := care.createVsAndRs(); err != nil {
		return err
	}
	if care.RunOnce {
		return nil
	}
	t := time.NewTicker(time.Duration(care.Interval) * time.Second)
	defer t.Stop()
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
					return err
				}
			}
			// check real server
			lvs.CheckRealServers(care.HealthPath, care.HealthSchem)
		case signa := <-sig:
			glog.Infof("receive kill signal: %+v", signa)
			for _, fn := range care.cleanupFuncs {
				if err := fn(); err != nil {
					return err
				}
			}
			return nil
		}
	}
}

func (care *LvsCare) test() error {
	if care.Test {
		ips, err := ipsFromAddrs(care.RealServer...)
		if err != nil {
			return err
		}
		hasLocal, err := isAnyLocalHostAddr(ips...)
		if err != nil {
			return err
		}
		if !hasLocal {
			return nil
		}
		vip, _, err := net.SplitHostPort(care.VirtualServer)
		if err != nil {
			return err
		}
		glog.Infof("create dummy interface with name %s", dummyIfaceName)
		link, err := utils.CreateDummyLinkIfNotExists(dummyIfaceName)
		if err != nil {
			return err
		}
		care.cleanupFuncs = append(care.cleanupFuncs, func() error {
			glog.Infof("remove dummy interface %s", dummyIfaceName)
			return utils.DeleteLinkByName(dummyIfaceName)
		})
		glog.Infof("assign IP %s/32 to interface", vip)
		return utils.AssignIPToLink(vip+"/32", link)
	}
	return nil
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
		if err = LVS.Route.SetRoute(); err != nil {
			return err
		}
		care.cleanupFuncs = append(care.cleanupFuncs, LVS.Route.DelRoute)
		return nil
	}
	return errors.New("real server can't empty")
}

func SetTargetIP() error {
	if LVS.TargetIP == nil {
		LVS.TargetIP = net.ParseIP(os.Getenv("LVSCARE_NODE_IP"))
	}

	return nil
}

func ipsFromAddrs(s ...string) ([]string, error) {
	var ret []string
	for i := range s {
		host, _, err := net.SplitHostPort(s[i])
		if err != nil {
			return nil, err
		}
		ret = append(ret, host)
	}
	return ret, nil
}

func isAnyLocalHostAddr(s ...string) (bool, error) {
	locals, err := utils.ListLocalHostAddrs()
	if err != nil {
		return false, err
	}
	var ips []net.IP
	for i := range s {
		ip := net.ParseIP(s[i])
		if ip == nil {
			return false, fmt.Errorf("%s is not a valid IP address", s)
		}
		ips = append(ips, ip)
	}

	for _, local := range locals {
		for i := range local.Addr {
			if ipn, ok := local.Addr[i].(*net.IPNet); ok {
				for _, ip := range ips {
					if ipn.IP.Equal(ip) {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}
