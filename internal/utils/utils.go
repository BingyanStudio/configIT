package utils

import "net"

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func IPInCIDR(ip, cidr string) bool {
	// Parse the IP address
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	// Parse the CIDR notation
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	// Check if the IP is contained in the CIDR range
	return ipNet.Contains(ipAddr)
}
