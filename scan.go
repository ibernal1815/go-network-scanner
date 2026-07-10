package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// Result holds the outcome of scanning one host:port.
type Result struct {
	Host string
	Port int
	Open bool
}

// scanPort attempts a TCP connect to host:port within timeout.
func scanPort(host string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Scanner runs a concurrent TCP connect scan over hosts x ports.
type Scanner struct {
	Timeout     time.Duration
	Concurrency int
}

func (s *Scanner) Run(hosts []string, ports []int) []Result {
	type job struct {
		host string
		port int
	}

	jobs := make(chan job)
	resultsCh := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				open := scanPort(j.host, j.port, s.Timeout)
				resultsCh <- Result{Host: j.host, Port: j.port, Open: open}
			}
		}()
	}

	go func() {
		for _, h := range hosts {
			for _, p := range ports {
				jobs <- job{host: h, port: p}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var open []Result
	for r := range resultsCh {
		if r.Open {
			open = append(open, r)
		}
	}

	sort.Slice(open, func(i, j int) bool {
		if open[i].Host != open[j].Host {
			return open[i].Host < open[j].Host
		}
		return open[i].Port < open[j].Port
	})
	return open
}
