package dhcpd

import (
	"fmt"
	"testing"
)

func TestHasStaticIP(t *testing.T) {

	nets := &DHCP{}

	fmt.Println(nets.HasStaticIP("eth0", true))
	fmt.Println(nets.SetStaticIP("eth0", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))
	fmt.Println(nets.SetStaticIP("eth1", "192.168.15.11/24", "192.168.15.1", "8.8.8.8"))

}
