package helper

import (
	"net"
	"strconv"
	"strings"

	"errors"
)

func DomainToIP(domains ...string) ([]string, error) {
	var result []string

	for _, domain := range domains {
		ips, err := net.LookupIP(domain)
		if err != nil {
			return nil, errors.New("[!] Cannot retrieve any IP address for " + domain)
		}

		for _, ip := range ips {
			result = append(result, ip.String())
		}
	}
	return result, nil
}

func RangeToPortsList(portsRange string) (*[]uint, error) {
	var ports []uint

	ranges := strings.Split(portsRange, "-")

	start, err := strconv.ParseUint(ranges[0], 10, 16)
	if err != nil {
		return nil, errors.New("[!] Cannot parse ports range")
	}

	end, err := strconv.ParseUint((ranges[len(ranges)-1]), 10, 16)
	if err != nil {
		return nil, errors.New("[!] Cannot parse ports range")
	}

	if start < end {
		for i := start; i <= end; i++ {
			ports = append(ports, uint(i))
		}
	} else if start > end {
		for i := end; i <= start; i++ {
			ports = append(ports, uint(i))
		}
	} else {
		ports = append(ports, uint(start))
	}

	return &ports, nil
}

func IsValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

func CIDRToIPList(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	// Remove network address and broadcast address
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}

	return ips, nil
}

// incrementIP increments the given IP address.
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
