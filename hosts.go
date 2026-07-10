package main

import "net"

// expandHosts accepts a single IP, hostname, or CIDR and returns target IPs.
// For a CIDR it drops the network and broadcast addresses.
func expandHosts(target string) ([]string, error) {
	if _, ipnet, err := net.ParseCIDR(target); err == nil {
		var ips []string
		ip := make(net.IP, len(ipnet.IP))
		copy(ip, ipnet.IP.Mask(ipnet.Mask))
		for ; ipnet.Contains(ip); incIP(ip) {
			ips = append(ips, ip.String())
		}
		if len(ips) > 2 {
			ips = ips[1 : len(ips)-1]
		}
		return ips, nil
	}
	return []string{target}, nil
}

func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
