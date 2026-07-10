package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// parsePorts turns a spec like "22,80,443,8000-8100" into a sorted, deduped slice.
func parsePorts(spec string) ([]int, error) {
	set := make(map[int]struct{})
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid port range %q: %w", part, err)
			}
			end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid port range %q: %w", part, err)
			}
			if start > end {
				start, end = end, start
			}
			for p := start; p <= end; p++ {
				if err := validatePort(p); err != nil {
					return nil, err
				}
				set[p] = struct{}{}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q: %w", part, err)
			}
			if err := validatePort(p); err != nil {
				return nil, err
			}
			set[p] = struct{}{}
		}
	}
	ports := make([]int, 0, len(set))
	for p := range set {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports, nil
}

func validatePort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("port %d out of range (1-65535)", p)
	}
	return nil
}
