package test

import (
	"github.com/sealyun/lvscare/utils"
	"testing"
)

func TestVirtualServer_Delete(t *testing.T) {
	//handle := ipvs.New()
	_ = utils.BuildVirtualServer("10.103.97.12:6443")
	//err := handle.DeleteVirtualServer(vs)
	//t.Error(err)
}
