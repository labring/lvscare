package test

import (
	"github.com/fanux/lvscare/internal/ipvs"
	"github.com/fanux/lvscare/utils"
	"testing"
)

func TestVirtualServer_Delete(t *testing.T) {
	handle := ipvs.New()
	vs := utils.BuildVirtualServer("10.103.97.12:6443")
	err := handle.DeleteVirtualServer(vs)
	t.Error(err)
}
