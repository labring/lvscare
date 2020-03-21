package service

import (
	"fmt"
	"github.com/fanux/LVScare/internal/ipvs"
	"github.com/fanux/LVScare/utils"
	"github.com/wonderivan/logger"
	"net"
	"strconv"
)

//EndPoint  is
type EndPoint struct {
	IP   string
	Port uint16
}

func (ep EndPoint) String() string {
	port := strconv.Itoa(int(ep.Port))
	return ep.IP + ":" + port
}

//Lvser is
type Lvser interface {
	CreateVirtualServer(vs string) error
	GetVirtualServer() (vs *EndPoint, rs []*EndPoint)
	//config is config to lvscare struct
	AddRealServer(rs string, config bool) error
	//config is config to lvscare struct
	RemoveRealServer(rs string, config bool) error

	CheckRealServers(path, schem string)
}

type lvscare struct {
	vs     *EndPoint
	rs     []*EndPoint
	handle ipvs.Interface
}

func (l *lvscare) buildVirtualServer(vip string) *ipvs.VirtualServer {
	ip, port := utils.SplitServer(vip)
	virServer := &ipvs.VirtualServer{
		Address:   net.ParseIP(ip),
		Protocol:  "TCP",
		Port:      port,
		Scheduler: "rr",
		Flags:     0,
		Timeout:   0,
	}
	return virServer
}

func (l *lvscare) buildRealServer(real string) *ipvs.RealServer {
	ip, port := utils.SplitServer(real)
	realServer := &ipvs.RealServer{
		Address: net.ParseIP(ip),
		Port:    port,
		Weight:  1,
	}
	return realServer
}

func (l *lvscare) CreateVirtualServer(vs string) error {
	virIp, virPort := utils.SplitServer(vs)
	if virIp == "" || virPort == 0 {
		logger.Error("CreateVirtualServer error: virtual server ip and port is empty")
		return fmt.Errorf("virtual server ip and port is empty")
	}
	l.vs = &EndPoint{IP: virIp, Port: virPort}
	vServer := l.buildVirtualServer(vs)
	err := l.handle.AddVirtualServer(vServer)
	if err != nil {
		logger.Error("CreateVirtualServer error: ", err)
		return fmt.Errorf("new virtual server failed: %s", err)
	}
	return nil
}

func (l *lvscare) AddRealServer(rs string, config bool) error {
	realIp, realPort := utils.SplitServer(rs)
	if realIp == "" || realPort == 0 {
		logger.Error("AddRealServer error: real server ip and port is empty")
		return fmt.Errorf("real server ip and port is empty")
	}
	rsEp := &EndPoint{IP: realIp, Port: realPort}
	if config {
		l.rs = append(l.rs, rsEp)
	}
	realServer := l.buildRealServer(rs)
	//vir server build
	if l.vs != nil {
		vServer := l.buildVirtualServer(l.vs.String())
		err := l.handle.AddRealServer(vServer, realServer)
		if err != nil {
			logger.Error("AddRealServer error: ", err)
			return fmt.Errorf("new real server failed: %s", err)
		}
		return nil
	} else {
		logger.Error("AddRealServer error: virtual server is empty.")
		return fmt.Errorf("virtual server is empty.")
	}

}
func (l *lvscare) GetVirtualServer() (vs *EndPoint, rs []*EndPoint) {
	virArray, err := l.handle.GetVirtualServers()
	if err != nil {
		logger.Error("GetVirtualServer error: vir servers is empty; ", err)
		return nil, nil
	}
	if l.vs != nil {
		resultVirServer := l.buildVirtualServer(l.vs.String())
		for _, vir := range virArray {
			logger.Debug("GetVirtualServer debug: check vir ip: %s, port %v ", vir.Address.String(), vir.Port)
			if vir.String() == resultVirServer.String() {
				return l.vs, l.rs
			}
		}
	} else {
		logger.Error("GetVirtualServer error: virtual server is empty.")
	}
	return
}

func (l *lvscare) healthCheck(ip, port, path, shem string) bool {
	return utils.IsHTTPAPIHealth(ip, port, path, shem)
}

func (l *lvscare) CheckRealServers(path, schem string) {
	//if realserver unavilable remove it, if recover add it back
	for _, realServer := range l.rs {
		ip := realServer.IP
		port := strconv.Itoa(int(realServer.Port))
		if !l.healthCheck(ip, port, path, schem) {
			err := l.RemoveRealServer(realServer.String(), false)
			if err != nil {
				logger.Warn("CheckRealServers error: %s;  %d; %v ", realServer.IP, realServer.Port, err)
			}
		} else {
			rs, weight := l.GetRealServer(realServer.String())
			if weight == 0 {
				err := l.RemoveRealServer(realServer.String(), false)
				logger.Debug("CheckRealServers debug: remove weight = 0 real server")
				if err != nil {
					logger.Warn("CheckRealServers error[remove weight = 0 real server failed]: %s;  %d; %v ", realServer.IP, realServer.Port, err)
				}
			}
			if rs == nil || weight == 0 {
				//add it back
				ip := realServer.IP
				port := strconv.Itoa(int(realServer.Port))
				err := l.AddRealServer(ip+":"+port, false)
				if err != nil {
					logger.Warn("CheckRealServers error[add real server failed]: %s;  %d; %v ", realServer.IP, realServer.Port, err)
				}
			}
		}
	}
}

func (l *lvscare) GetRealServer(rsHost string) (rs *EndPoint, weight int) {
	ip, port := utils.SplitServer(rsHost)
	vs := l.buildVirtualServer(l.vs.String())
	dstArray, err := l.handle.GetRealServers(vs)
	if err != nil {
		logger.Error("GetRealServer error[get real server failed]: %s;  %d; %v ", ip, port, err)
		return nil, 0
	}
	dip := net.ParseIP(ip)
	for _, dst := range dstArray {
		logger.Debug("GetRealServer debug[check real server ip]: %s;  %d; %v ", dst.Address.String(), dst.Port, err)
		if dst.Address.Equal(dip) && dst.Port == port {
			return &EndPoint{IP: ip, Port: port}, dst.Weight
		}
	}
	return nil, 0
}

//
func (l *lvscare) RemoveRealServer(rs string, config bool) error {
	realIp, realPort := utils.SplitServer(rs)
	if realIp == "" || realPort == 0 {
		logger.Error("RemoveRealServer error: real server ip and port is empty ")
		return fmt.Errorf("real server ip and port is null")
	}

	if l.vs == nil {
		logger.Error("RemoveRealServer error: virtual service is empty.")
		return fmt.Errorf("virtual service is empty.")
	}
	virServer := l.buildVirtualServer(l.vs.String())
	realServer := l.buildRealServer(rs)
	err := l.handle.DeleteRealServer(virServer, realServer)
	if err != nil {
		logger.Error("RemoveRealServer error[real server delete error]: ", err)
		return fmt.Errorf("real server delete error: %v", err)
	}
	if config {
		//clean delete data
		var resultRS []*EndPoint
		for _, r := range l.rs {
			if r.IP == realIp && r.Port == realPort {
				continue
			} else {
				resultRS = append(resultRS, &EndPoint{
					IP:   r.IP,
					Port: r.Port,
				})
			}
		}
		l.rs = resultRS
	}

	return nil
}
