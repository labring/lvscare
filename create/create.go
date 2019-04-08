package create

import (
	"fmt"

	"github.com/fanux/LVScare/service"
	"github.com/fanux/LVScare/utils"
)

//VsAndRsCreate is
func VsAndRsCreate(vs string, rs []string) error {
	//ip, _ := utils.SplitServer(vs)
	lvs, err := service.BuildLvscare(vs, rs)
	if err != nil {
		return fmt.Errorf("build lvs failed: %s", err)
	}

	/*
		err = lvs.CreateInterface("sealyun-ipvs", ip+"/32")
		if err != nil {
			fmt.Println(err)
		}
	*/

	err = lvs.CreateVirtualServer()
	if err != nil {
		return err
	}

	for _, r := range rs {
		i, p := utils.SplitServer(r)
		if i == "" || p == "" {
			return fmt.Errorf("ip and port is null")
		}
		err = lvs.AddRealServer(i, p)
		if err != nil {
			return err
		}
	}

	return nil
}
