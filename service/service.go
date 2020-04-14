package service

import (
	"fmt"
	"net"
	"strconv"
	"syscall"

	"github.com/fanux/lvscare/pkg/nl"
	"github.com/fanux/lvscare/utils"
)

//EndPoint  is
type EndPoint struct {
	IP   string
	Port string
}

//Lvser is
type Lvser interface {
	//CreateInterface(name string, CIRD string) error
	CreateVirtualServer() error
	AddRealServer(ip, port string) error
	GetVirtualServer() (vs *EndPoint, rs *[]EndPoint)
	RemoveVirtualServer() error
	GetRealServer(ip, port string) (rs *EndPoint, weight int)
	RemoveRealServer(ip, port string) error

	CheckRealServers(path, schem string)
	Close()
}

type lvscare struct {
	vs           EndPoint
	rs           []EndPoint
	service      *Service
	destinations []*Destination
	handle       *Handle
}

func (l *lvscare) Close() {
	l.handle.Close()
}

/*
func (l *lvscare) CreateInterface(name string, CIRD string) error {
	interfa := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: name}}
	err := netlink.LinkAdd(interfa)
	if err != nil {
		return fmt.Errorf("create net interface failed: %s", err)
	}

	kubeIpvs, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("get interface failed: %s", err)
	}

	fmt.Println("cird is: ", CIRD)
	addr, err := netlink.ParseAddr(CIRD)
	if err != nil {
		return fmt.Errorf("CIRD failed: %s", err)
	}
	err = netlink.AddrAdd(kubeIpvs, addr)
	if err != nil {
		return fmt.Errorf("bind IP addr failed: %s", err)
	}

	return err
}
*/

func (l *lvscare) CreateVirtualServer() error {
	err := l.handle.NewService(l.service)
	if err != nil {
		return fmt.Errorf("New ipvs failed: %s", err)
	}

	return nil
}

func (l *lvscare) AddRealServer(ip, port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port is %s : %s", port, err)
	}

	if l.service == nil {
		return fmt.Errorf("service is nil: %s", err)
	}

	d := Destination{
		AddressFamily:   nl.FAMILY_V4,
		Address:         net.ParseIP(ip),
		Port:            uint16(p),
		Weight:          1,
		ConnectionFlags: ConnectionFlagMasq,
	}

	err = l.handle.NewDestination(l.service, &d)
	if err != nil {
		return fmt.Errorf("new destination failed: %s", err)
	}

	return nil
}

func (l *lvscare) GetVirtualServer() (vs *EndPoint, rs *[]EndPoint) {
	svcArray, err := l.handle.GetServices()
	if err != nil {
		fmt.Println("services is nil", err)
		return nil, nil
	}
	if len(svcArray) == 0 {
		fmt.Println("services list is empty", err)
		return nil, nil
	}

	for _, svc := range svcArray {
		if svc.Address.String() == l.service.Address.String() && svc.Port == l.service.Port {
			return &l.vs, &l.rs
		}
	}

	return nil, nil
}

func (l *lvscare) RemoveVirtualServer() error {
	return nil
}

func (l *lvscare) GetRealServer(ip, port string) (rs *EndPoint, weight int) {
	dip := net.ParseIP(ip)
	p, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("port is %s : %s", port, err)
		return nil, 0
	}
	dport := uint16(p)

	dstArray, err := l.handle.GetDestinations(l.service)
	if err != nil {
		fmt.Printf("get real servers failed %s : %s\n", ip, port)
		return nil, 0
	}

	for _, dst := range dstArray {
		if dst.Address.Equal(dip) && dst.Port == dport {
			return &EndPoint{IP: ip, Port: port}, dst.Weight
		}
	}

	return nil, 0
}

func (l *lvscare) RemoveRealServer(ip, port string) error {
	p, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port is %s : %s", port, err)
	}

	if l.service == nil {
		return fmt.Errorf("service is nil: %s", err)
	}

	d := &Destination{
		AddressFamily:   nl.FAMILY_V4,
		Address:         net.ParseIP(ip),
		Port:            uint16(p),
		Weight:          1,
		ConnectionFlags: ConnectionFlagMasq,
	}

	l.handle.DelDestination(l.service, d)
	return nil
}

func (l *lvscare) healthCheck(ip, port, path, shem string) bool {
	return utils.IsHTTPAPIHealth(ip, port, path, shem)
}

func (l *lvscare) CheckRealServers(path, schem string) {
	//if realserver unavilable remove it, if recover add it back
	for _, realServer := range l.rs {
		if !l.healthCheck(realServer.IP, realServer.Port, path, schem) {
			err := l.RemoveRealServer(realServer.IP, realServer.Port)
			if err != nil {
				fmt.Printf("remove real server failed %s:%s", realServer.IP, realServer.Port)
			}
			continue
		}
		rs, weight := l.GetRealServer(realServer.IP, realServer.Port)
		if weight == 0 {
			err := l.RemoveRealServer(realServer.IP, realServer.Port)
			fmt.Println("remove weight = 0 real server")
			if err != nil {
				fmt.Println("	Error remove weight = 0 real server failed", realServer.IP, realServer.Port)
			}
		}
		if rs == nil || weight == 0 {
			//add it back
			err := l.AddRealServer(realServer.IP, realServer.Port)
			if err != nil {
				fmt.Printf("add real server failed %s:%s", realServer.IP, realServer.Port)
			}
		}
	}
}

//BuildLvscare is
func BuildLvscare(vs string, rs []string) (Lvser, error) {
	ip, port := utils.SplitServer(vs)
	if ip == "" || port == "" {
		return nil, fmt.Errorf("virtual server ip and port is null")
	}

	l := &lvscare{
		vs: EndPoint{IP: ip, Port: port},
	}

	p, err := strconv.Atoi(l.vs.Port)
	if err != nil {
		return nil, fmt.Errorf("port is %s : %s", l.vs.Port, err)
	}

	l.service = &Service{
		AddressFamily: nl.FAMILY_V4,
		SchedName:     RoundRobin,
		Protocol:      syscall.IPPROTO_TCP,
		Port:          uint16(p),
		Address:       net.ParseIP(l.vs.IP),
		Netmask:       0xFFFFFFFF,
	}

	handle, err := New("")
	if err != nil {
		return nil, fmt.Errorf("New ipvs handle failed: %s", err)
	}
	l.handle = handle

	for _, r := range rs {
		i, p := utils.SplitServer(r)
		if i == "" || p == "" {
			return nil, fmt.Errorf("real server ip and port is null")
		}

		iport, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("real server port failed: %s", err)
		}

		l.rs = append(l.rs, EndPoint{IP: i, Port: p})
		l.destinations = append(l.destinations, &Destination{
			AddressFamily:   nl.FAMILY_V4,
			Address:         net.ParseIP(i),
			Port:            uint16(iport),
			Weight:          1,
			ConnectionFlags: ConnectionFlagMasq,
		})
	}

	return l, nil
}
