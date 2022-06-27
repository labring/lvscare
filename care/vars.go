package care

import (
	"github.com/labring/lvscare/internal/route"
	"github.com/labring/lvscare/service"
	"net"
)

type LvsCare struct {
	HealthPath    string
	HealthSchem   string // http or https
	VirtualServer string
	RealServer    []string
	RunOnce       bool
	Clean         bool
	Interval      int32
	Logger        string
	TargetIP      net.IP
	//
	lvs   service.Lvser
	Route *route.Route
}

var LVS LvsCare
