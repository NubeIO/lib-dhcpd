package dhcpd

import (
	"fmt"
	"testing"
)

func TestHasStaticIP(t *testing.T) {
	//Raspi currently using interface wlan0 (check with ifconfig)
	//Interface relevant for testing as raspi uses interface wlan0
	testInterface := "wlp0s20f3"

	nets := New()

	out, err := nets.SetAsAuto(testInterface)
	fmt.Println(out, err)
	out, err = nets.IsStaticIP(testInterface)
	fmt.Println(out, err)
	err = nets.SetStaticIP(testInterface, "192.168.15.222", "192.168.15.1", "1.1.1.1")
	//fmt.Println(out, err)

	//fmt.Println(nets.SetStaticIP(" wlp3s0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))
	//fmt.Println(nets.IsStaticIP("wlp3s0"))
	//fmt.Println(nets.SetAsAuto("wlp3s0"))
	//fmt.Println(nets.IsStaticIP("wlp3s0"))
	//fmt.Println(nets.SetStaticIP("eth1", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))

}
