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
