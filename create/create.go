package create

import (
	"github.com/fanux/LVScare/service"
)

//VsAndRsCreate is
func VsAndRsCreate(vs string, rs []string) error {
	//ip, _ := utils.SplitServer(vs)
	lvs := service.BuildLvscare()

	err := lvs.CreateVirtualServer(vs)
	if err != nil {
		return err
	}
	for _, r := range rs {
		err = lvs.AddRealServer(r)
		if err != nil {
			return err
		}
	}
	return nil
}
