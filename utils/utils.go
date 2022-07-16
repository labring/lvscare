package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/labring/lvscare/internal/glog"
	"github.com/labring/lvscare/internal/ipvs"
	"github.com/vishvananda/netlink"
)

//SplitServer is
func SplitServer(server string) (string, uint16) {
	glog.V(8).Infof("server %s", server)

	ip, port, err := net.SplitHostPort(server)
	if err != nil {
		glog.Errorf("SplitServer error: %v.", err)
		return "", 0
	}
	glog.V(8).Infof("SplitServer debug: TargetIP: %s, Port: %s", ip, port)
	p, err := strconv.Atoi(port)
	if err != nil {
		glog.Warningf("SplitServer error: %v", err)
		return "", 0
	}
	return ip, uint16(p)
}

//IsHTTPAPIHealth is check http error
func IsHTTPAPIHealth(ip, port, path, schem string) bool {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("%s://%s%s", schem, net.JoinHostPort(ip, port), path)
	resp, err := http.Get(url)
	if err != nil {
		glog.V(8).Infof("IsHTTPAPIHealth error: %v", err)
		return false
	}
	defer resp.Body.Close()

	_ = resp
	return true
}

func BuildVirtualServer(vip string) *ipvs.VirtualServer {
	ip, port := SplitServer(vip)
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

func BuildRealServer(real string) *ipvs.RealServer {
	ip, port := SplitServer(real)
	realServer := &ipvs.RealServer{
		Address: net.ParseIP(ip),
		Port:    port,
		Weight:  1,
	}
	return realServer
}

// IsIPv6 returns if netIP is IPv6.
func IsIPv6(netIP net.IP) bool {
	return netIP != nil && netIP.To4() == nil
}

func IsIpv4(ip string) bool {
	//matched, _ := regexp.MatchString("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}", ip)

	arr := strings.Split(ip, ".")
	if len(arr) != 4 {
		return false
	}
	for _, v := range arr {
		if v == "" {
			return false
		}
		if len(v) > 1 && v[0] == '0' {
			return false
		}
		num := 0
		for _, c := range v {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			} else {
				return false
			}
		}
		if num > 255 {
			return false
		}
	}
	return true
}

type LocalAddr struct {
	Addr      []net.Addr
	Interface net.Interface
}

func ListLocalHostAddrs() ([]LocalAddr, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var allAddrs []LocalAddr
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := netInterfaces[i].Addrs()
		if err != nil {
			fmt.Println("failed to get Addrs, ", err.Error())
		}
		var ipAddrs []net.Addr
		for j := 0; j < len(addrs); j++ {
			ipAddrs = append(ipAddrs, addrs[j])
		}
		allAddrs = append(allAddrs, LocalAddr{
			Addr:      ipAddrs,
			Interface: netInterfaces[i],
		})
	}
	return allAddrs, nil
}

func GetLocalDefaultRoute(ip string) (*LocalAddr, error) {
	locals, err := ListLocalHostAddrs()
	if err != nil {
		return nil, err
	}
	for _, local := range locals {
		for _, add := range local.Addr {
			glog.Infof("local net interface: %s, addr: %s", local.Interface.Name, add.String())
			if isContains, err := Contains(add.String(), ip); err != nil {
				return nil, err
			} else {
				if isContains {
					return &local, nil
				}
			}
		}
	}
	return nil, errors.New("not found contains ip")
}

func Contains(sub, s string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(sub)
	if err != nil {
		return false, err
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return false, fmt.Errorf("%s is not a valid TargetIP address", s)
	}
	return ipNet.Contains(ip), nil
}

func CreateDummyLinkIfNotExists(name string) (netlink.Link, error) {
	link, _ := netlink.LinkByName(name)
	if link != nil {
		return link, nil
	}
	link = &netlink.Dummy{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}
	return link, netlink.LinkAdd(link)
}

func AssignIPToLink(s string, link netlink.Link) error {
	if strings.Index(s, "/") < 0 {
		s = s + "/32"
	}
	addr, err := netlink.ParseAddr(s)
	if err != nil {
		return err
	}
	return netlink.AddrAdd(link, addr)
}

func DeleteLinkByName(name string) error {
	l, err := netlink.LinkByName(name)
	if err != nil {
		if err, ok := err.(netlink.LinkNotFoundError); !ok {
			return err
		}
		return nil
	}
	return netlink.LinkDel(l)
}
