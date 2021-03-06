package dhcpd

import (
	"fmt"
	"testing"
)

func TestHasStaticIP(t *testing.T) {

	nets := New()

	out, err := nets.SetAsAuto("wlp3s0")
	fmt.Println(out, err)
	out, err = nets.IsStaticIP("wlp3s0")
	fmt.Println(out, err)
	err = nets.SetStaticIP("wlp3s0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8")
	fmt.Println(out, err)

	//fmt.Println(nets.SetStaticIP(" wlp3s0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))
	//fmt.Println(nets.IsStaticIP("wlp3s0"))
	//fmt.Println(nets.SetAsAuto("wlp3s0"))
	//fmt.Println(nets.IsStaticIP("wlp3s0"))
	//fmt.Println(nets.SetStaticIP("eth1", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))

}
