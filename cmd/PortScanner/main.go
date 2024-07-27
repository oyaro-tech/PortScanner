package main

import (
	"flag"
	"log"
	"net/netip"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/oyaro-tech/PortScanner/internal/helper"
	"github.com/oyaro-tech/PortScanner/internal/scanner"
)

var (
	targetsFlag string
	portsFlag   string
	threadsFlag uint
	timeoutFlag uint
)

// TODO: Proxy flag - for usage tor or other services as proxy
// TODO: Add colors for logs

func main() {
	log.SetFlags(0)

	flag.StringVar(&targetsFlag, "target", "127.0.0.1", "comma splitted list of targets [IPv4, IPv6, CIDR, Domain]")
	flag.StringVar(&portsFlag, "port", "80,443", "comma splitted list of ports [supports ranges 22-1023]")
	flag.UintVar(&threadsFlag, "thread", 100, "number of threads")
	flag.UintVar(&timeoutFlag, "timeout", 1, "timeout for port probe in seconds")

	flag.Parse()

	if targetsFlag == "" {
		log.Fatalf("[!] Targets list cannot be empty")
	} else if portsFlag == "" {
		log.Fatalf("[!] Ports list cannot be empty")
	} else if threadsFlag <= 0 {
		log.Fatalf("[!] Threads cannot be less then 1")
	} else if timeoutFlag <= 0 {
		log.Fatalf("[!] Timeout cannot be less then 1 second")
	}

	var targets []string
	var ports []uint

	for _, port := range strings.Split(portsFlag, ",") {
		rangeRegexp, err := regexp.Compile("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])-([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])?$")
		if err != nil {
			log.Fatalf("[!] Regexp error: %s", err.Error())
		}
		portRegexp, err := regexp.Compile("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])?$")
		if err != nil {
			log.Fatalf("[!] Regexp error: %s", err.Error())
		}

		if rangeRegexp.MatchString(port) {
			_ports, err := helper.RangeToPortsList(port)
			if err != nil {
				log.Fatal(err.Error())
			}

			ports = append(ports, *_ports...)
			continue
		} else if portRegexp.MatchString(port) {
			_port, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				log.Fatalf("[!] Cannot convert port %s", port)
			}

			ports = append(ports, uint(_port))
			continue
		}

		log.Fatalf("[!] Invalid port %s", port)
	}

	for _, host := range strings.Split(targetsFlag, ",") {
		if _, err := netip.ParseAddr(host); err != nil {
			if helper.IsValidCIDR(host) {
				ips, err := helper.CIDRToIPList(host)
				if err != nil {
					log.Fatal(err.Error())
				}

				targets = append(targets, ips...)
				continue
			}

			ips, err := helper.DomainToIP(host)
			if err != nil {
				log.Fatal(err.Error())
			}

			targets = append(targets, ips...)
			continue
		}

		targets = append(targets, host)
	}

	slices.Sort(targets)
	slices.Sort(ports)

	targets = slices.Compact(targets)
	ports = slices.Compact(ports)

	s := scanner.NewScanner(threadsFlag, timeoutFlag, ports...)
	s.Scan(targets...)
}
