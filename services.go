package main

var commonServices = map[int]string{
	21: "ftp", 22: "ssh", 23: "telnet", 25: "smtp", 53: "dns",
	80: "http", 110: "pop3", 143: "imap", 443: "https",
	3306: "mysql", 3389: "rdp", 5432: "postgresql", 6379: "redis",
	8080: "http-alt", 8443: "https-alt",
}

func serviceName(port int) string {
	if name, ok := commonServices[port]; ok {
		return name
	}
	return ""
}
