package care

import (
	"net"

	"github.com/labring/lvscare/internal/route"
	"github.com/labring/lvscare/service"
)

type LvsCare struct {
	HealthPath    string
	HealthSchem   string // http or https
	VirtualServer string
	RealServer    []string
	RunOnce       bool
	Clean         bool
	Test          bool
	Interval      int32
	Logger        string
	TargetIP      net.IP

	lvs          service.Lvser
	cleanupFuncs []func() error
	Route        *route.Route
}

var LVS LvsCare

const (
	dummyIfaceName = "lvscare"
)
