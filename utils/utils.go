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
	netIP := net.ParseIP(ip)
	return netIP != nil && netIP.To4() != nil
}

type LocalAddr struct {
	Addr      []net.Addr
	Interface net.Interface
}

func IsLocalHostAddrs() ([]LocalAddr, error) {
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
	locals, err := IsLocalHostAddrs()
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
