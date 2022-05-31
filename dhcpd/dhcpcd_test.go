package dhcpd

import (
	"fmt"
	"testing"
)

func TestHasStaticIP(t *testing.T) {

	nets := &DHCP{}

	fmt.Println(nets.SetAsAuto("wlp3s0"))
	fmt.Println(nets.IsStaticIP("wlp3s0"))
	fmt.Println(nets.SetStaticIP(" wlp3s0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))
	fmt.Println(nets.IsStaticIP("wlp3s0"))
	fmt.Println(nets.SetAsAuto("wlp3s0"))
	fmt.Println(nets.IsStaticIP("wlp3s0"))
	//fmt.Println(nets.SetStaticIP("eth1", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))

}
