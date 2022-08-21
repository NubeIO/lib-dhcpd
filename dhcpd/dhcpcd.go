package dhcpd

import (
	"errors"
	"fmt"
	address "github.com/NubeIO/lib-networking/ip"
	"github.com/NubeIO/lib-networking/networking"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var filePath = "/etc/dhcpcd.conf" // is normally /etc/dhcpcd.conf

type DHCP struct {
	FilePath string
}

func New(opts *DHCP) *DHCP {
	if opts == nil {
		opts = &DHCP{}
	}
	if opts.FilePath != "" {
		filePath = opts.FilePath
	}
	return opts
}

type SetStaticIP struct {
	Ip                   string `json:"ip"`
	NetMask              string `json:"net_mask"`
	IFaceName            string `json:"i_face_name"`
	GatewayIP            string `json:"gateway_ip"`
	DnsIP                string `json:"dns_ip"`
	CheckInterfaceExists bool   `json:"check_interface_exists"`
	SaveFile             bool   `json:"save_file"`
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

// Exists check interface
func (inst *DHCP) Exists(iFaceName string) (bool, error) {
	if isLinux() {
		body, err := ioutil.ReadFile(filePath)
		if err != nil {
			return false, err
		}
		return hasStaticIPDhcpcdConf(string(body), iFaceName, false), nil
	}
	return false, fmt.Errorf("cannot check if IP is static: not supported on %s", runtime.GOOS)
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

//SetStaticIP Set a static IP for the specified network interface
func (inst *DHCP) SetStaticIP(body *SetStaticIP) (string, error) {
	if body == nil {
		return "", errors.New("body can not be empty")
	}
	if body.Ip == "" {
		return "", errors.New("ip can not be empty")
	}
	if body.NetMask == "" {
		return "", errors.New("NetMask can not be empty")
	}
	if body.GatewayIP == "" {
		return "", errors.New("GatewayIP can not be empty")
	}
	if body.IFaceName == "" {
		return "", errors.New("interface name can not be empty")
	}
	_, err := address.New().IsIPSubnet(body.NetMask)
	if err != nil {
		return "", err
	}
	ipAndSub, _, err := address.GetIPSubnet(body.Ip, body.NetMask)
	if err != nil {
		return "", err
	}
	_, err = inst.SetAsAuto(body.IFaceName) // remove if existing
	if err != nil {
		return "", err
	}
	if body.CheckInterfaceExists {
		_, err := nets.CheckInterfacesName(body.IFaceName)
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("network interface not found:%s", body.IFaceName))
		}
	}
	if isLinux() {
		return setStaticIPDHCPConf(body.IFaceName, ipAndSub, body.GatewayIP, body.DnsIP, body.SaveFile)
	}
	return "", fmt.Errorf("cannot set static IP on %s", runtime.GOOS)
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
func setStaticIPDHCPConf(iFaceName, ip, gatewayIP, dnsIP string, writeFile bool) (string, error) {
	//nets := networking.New()
	//if checkInterfaceExists {
	//	ipV4, err := nets.GetNetworkByIface(iFaceName)
	//	if err != nil {
	//		return "", err
	//	}
	//	ip = ipV4.IP
	//	gatewayIP, _ = nets.GetGatewayIP(iFaceName)
	//	if dnsIP == "" {
	//		dnsIP = ip
	//	}
	//}

	add := updateStaticIPDhcpcdConf(iFaceName, ip, gatewayIP, dnsIP)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	body = append(body, []byte(add)...)
	if writeFile {
		err = ioutil.WriteFile(filePath, body, 0755)
		if err != nil {
			return "", err
		}
	}
	return string(body), nil
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
