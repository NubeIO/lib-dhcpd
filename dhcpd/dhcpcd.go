package dhcpd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/NubeIO/lib-networking/networking"
)

const filePath = "dhcpcd.conf" // is normally /etc/dhcpcd.conf

type DHCP struct {
}

func New() *DHCP {
	return &DHCP{}
}

var nets = networking.New()

type NetInterface struct {
	Name         string   // Network interface name
	MTU          int      // MTU
	HardwareAddr string   // Hardware address
	Addresses    []string // Array with the network interface addresses
	Subnets      []string // Array with CIDR addresses of this network interface
	Flags        string   // Network interface flags (up, broadcast, etc)
}

func removeLine(path string, lineNumber int) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	info, _ := os.Stat(path)
	mode := info.Mode()
	array := strings.Split(string(file), "\n")
	array = append(array[:lineNumber], array[lineNumber+1:]...)
	err = ioutil.WriteFile(path, []byte(strings.Join(array, "\n")), mode)
	return err
}

// SetAsAuto check to auto address
func (inst *DHCP) SetAsAuto(iFaceName string) (bool, error) {
	if isLinux() {
		body, err := ioutil.ReadFile(filePath)
		if err != nil {
			return false, err
		}
		return hasStaticIPDhcpcdConf(string(body), iFaceName, true), nil
	}
	return false, fmt.Errorf("cannot check if IP is static: not supported on %s", runtime.GOOS)
}

// IsStaticIP Check if network interface has a static IP configured
func (inst *DHCP) IsStaticIP(iFaceName string) (bool, error) {
	if isLinux() {
		body, err := ioutil.ReadFile(filePath)
		if err != nil {
			return false, err
		}
		return hasStaticIPDhcpcdConf(string(body), iFaceName, false), nil
	}
	return false, fmt.Errorf("cannot check if IP is static: not supported on %s", runtime.GOOS)
}

//SetStaticIP Set a static IP for the specified network interface
func (inst *DHCP) SetStaticIP(iFaceName, ip, gatewayIP, dnsIP string) error {
	iface, err := nets.CheckInterfacesName(iFaceName)
	if err != nil {
		return err
	}
	if !iface {
		return fmt.Errorf("network interface not found")
	}
	if isLinux() {
		return setStaticIPDHCPConf(iFaceName, ip, gatewayIP, dnsIP)
	}
	return fmt.Errorf("cannot set static IP on %s", runtime.GOOS)
}

// for dhcpcd.conf
func hasStaticIPDhcpcdConf(dhcpConf, iFaceName string, delete bool) bool {
	lines := strings.Split(dhcpConf, "\n")
	nameLine := fmt.Sprintf("interface %s", iFaceName)
	withinInterfaceCtx := false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if withinInterfaceCtx && len(line) == 0 {
			// an empty line resets our state
			withinInterfaceCtx = false
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		line = strings.TrimSpace(line)
		if !withinInterfaceCtx {
			matched, _ := regexp.MatchString(line, nameLine)
			if matched {
				// we found our interface
				withinInterfaceCtx = true
				if delete {
					for ii := 0; ii < 4; ii++ {
						fmt.Println("line number ", i, line, ii)
						err := removeLine(filePath, i)
						if err != nil {
							return false
						}
					}
				}
			}
		} else {
			if strings.HasPrefix(line, "interface ") {
				// we found another interface - reset our state
				withinInterfaceCtx = false
				continue
			}
			if strings.HasPrefix(line, "static ip_address=") {
				return true
			}
		}
	}
	return false
}

// setStaticIPDHCPConf - updates /etc/dhcpd.conf and sets the current IP address to be static
func setStaticIPDHCPConf(iFaceName, ip, gatewayIP, dnsIP string) error {

	add := updateStaticIPDhcpcdConf(iFaceName, ip, gatewayIP, dnsIP)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	body = append(body, []byte(add)...)
	err = ioutil.WriteFile(filePath, body, 0755)
	if err != nil {
		return err
	}
	return nil
}

// updates dhcpd.conf content -- sets static IP address there
// for dhcpcd.conf
func updateStaticIPDhcpcdConf(iFaceName, ip, gatewayIP, dnsIP string) string {
	var body []byte
	add := fmt.Sprintf("\ninterface %s\nstatic ip_address=%s\n",
		iFaceName, ip)

	body = append(body, []byte(add)...)

	if len(gatewayIP) != 0 {
		add = fmt.Sprintf("static routers=%s\n",
			gatewayIP)
		body = append(body, []byte(add)...)
	}
	add = fmt.Sprintf("static domain_name_servers=%s\n\n",
		dnsIP)
	body = append(body, []byte(add)...)
	return string(body)
}

// Gets a list of nameservers currently configured in the /etc/resolv.conf
func getEtcResolvConfServers() ([]string, error) {
	body, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("nameserver ([a-zA-Z0-9.:]+)")
	matches := re.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return nil, errors.New("found no DNS servers in /etc/resolv.conf")
	}
	addrs := make([]string, 0)
	for i := range matches {
		addrs = append(addrs, matches[i][1])
	}
	return addrs, nil
}

func isLinux() bool {
	if runtime.GOOS == "linux" {
		return true
	} else {
		return false
	}
}

// Restarting and flushing the services after successful IP file write operation
func restartServices(iFaceName string) {
	// Commands needed to run on raspi
	// sudo sytemctl daemon-reload
	// sudo sytemctl stop dhcpcd.service

	// sudo ip addr flush dev INTERFACE

	// sudo systemctl start dhcpcd.service
	//sudo systemctl restart networking.service

}
