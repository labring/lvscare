module github.com/labring/lvscare

go 1.13

require (
	github.com/lithammer/dedent v1.1.0
	github.com/moby/ipvs v1.0.1
	github.com/spf13/cobra v1.2.1
	github.com/vishvananda/netlink v1.1.0
)

replace (
	github.com/moby/ipvs => github.com/moby/ipvs v1.0.1
	github.com/vishvananda/netlink => github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns => github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae
)
