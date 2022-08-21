package dhcpd

import (
	"fmt"
	"testing"
)

var iface = "wlp0s20f3"

func Test_SetStaticIP(t *testing.T) {
	nets := New(&DHCP{FilePath: "dhcpcd.conf"})
	conf, err := nets.SetStaticIP(&SetStaticIP{
		Ip:                   "10.0.40.22",
		NetMask:              "255.255.0.0",
		IFaceName:            iface,
		GatewayIP:            "192.168.15.1",
		DnsIP:                "8.8.8.8",
		CheckInterfaceExists: false,
		SaveFile:             true,
	})
	fmt.Println(err)
	fmt.Println(conf)
}

func Test_Exists(*testing.T) {
	nets := New(&DHCP{FilePath: "dhcpcd.conf"})
	out, err := nets.Exists(iface)
	fmt.Println(out, err)
}

func Test_SetAsAuto(t *testing.T) {
	nets := New(&DHCP{FilePath: "dhcpcd.conf"})
	out, err := nets.SetAsAuto(iface)
	fmt.Println(out, err)
}
