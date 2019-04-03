package create

import (
	"fmt"

	"github.com/fanux/LVScare/service"
	"github.com/fanux/LVScare/utils"
)

//VsAndRsCreate is
func VsAndRsCreate(vs string, rs []string) error {
	ip, port := utils.SplitServer(vs)
	if ip == "" || port == "" {
		return fmt.Errorf("ip and port is null")
	}
	lvs := service.NewLvscare(ip, port)

	err := lvs.CreateInterface("sealyun-ipvs", ip+"/32")
	if err != nil {
		return err
	}

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
