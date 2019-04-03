package service

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

//EndPoint  is
type EndPoint struct {
	IP   string
	Port string
}

//Lvser is
type Lvser interface {
	CreateInterface(name string, CIRD string) error
	CreateVirtualServer() error
	AddRealServer(ip, port string) error
	GetVirtualServer() (vs *EndPoint, rs *[]EndPoint)
	RemoveVirtualServer() error
	RemoveRealServer(ip, port string) error
}

type lvscare struct {
	vs EndPoint
	rs []EndPoint
}

func (l *lvscare) CreateInterface(name string, CIRD string) error {
	kubeIpvs1, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	fmt.Println("cird is: ", CIRD)
	addr, err := netlink.ParseAddr(CIRD)
	if err != nil {
		return err
	}
	netlink.AddrAdd(kubeIpvs1, addr)

	return err
}

func (l *lvscare) CreateVirtualServer() error {
	return nil
}

func (l *lvscare) AddRealServer(ip, port string) error {
	return nil
}

func (l *lvscare) GetVirtualServer() (vs *EndPoint, rs *[]EndPoint) {
	return nil, nil
}

func (l *lvscare) RemoveVirtualServer() error {
	return nil
}

func (l *lvscare) RemoveRealServer(ip, port string) error {
	return nil
}

//NewLvscare is
func NewLvscare(ip, port string) Lvser {
	return &lvscare{
		vs: EndPoint{IP: ip, Port: port},
	}
}
