module github.com/fanux/lvscare/pkg/netlink

go 1.13

replace (
    github.com/vishvananda/netlink/nl => ./nl
    github.com/vishvananda/netns => ../netns
)
