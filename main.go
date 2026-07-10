package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	target := flag.String("target", "", "target IP, hostname, or CIDR (e.g. 192.168.1.0/24)")
	portSpec := flag.String("ports", "1-1024", "ports to scan (e.g. 22,80,443 or 1-1024)")
	timeout := flag.Duration("timeout", 2*time.Second, "per-connection timeout")
	concurrency := flag.Int("concurrency", 200, "number of concurrent workers")
	flag.Parse()

	if *target == "" {
		fmt.Fprintln(os.Stderr, "error: -target is required")
		flag.Usage()
		os.Exit(1)
	}

	hosts, err := expandHosts(*target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing target: %v\n", err)
		os.Exit(1)
	}

	ports, err := parsePorts(*portSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing ports: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Scanning %d host(s) x %d port(s)...\n", len(hosts), len(ports))
	start := time.Now()

	scanner := &Scanner{Timeout: *timeout, Concurrency: *concurrency}
	results := scanner.Run(hosts, ports)

	if len(results) == 0 {
		fmt.Println("No open ports found.")
	} else {
		var currentHost string
		for _, r := range results {
			if r.Host != currentHost {
				currentHost = r.Host
				fmt.Printf("\n%s\n", currentHost)
			}
			fmt.Printf("  %d/tcp open  %s\n", r.Port, serviceName(r.Port))
		}
	}

	fmt.Printf("\nDone in %s\n", time.Since(start).Round(time.Millisecond))
}
