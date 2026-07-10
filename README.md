# go-network-scanner

A fast, concurrent TCP port scanner written in Go. Give it a host, hostname, or CIDR range, and it reports which ports are open.

## Features

- Concurrent scanning with a configurable worker pool
- Accepts a single host, a hostname, or a CIDR block (e.g. `192.168.1.0/24`)
- Flexible port specs: single ports, comma lists, and ranges (`22,80,443,8000-8100`)
- Configurable per-connection timeout
- Common service name lookup for well-known ports

## Installation

```bash
git clone https://github.com/yourusername/go-network-scanner.git
cd go-network-scanner
go build -o go-network-scanner .
```

Requires Go 1.22 or later.

## Usage

```bash
./go-network-scanner -target <host|CIDR> [flags]
```

### Flags

| Flag           | Default   | Description                                      |
|----------------|-----------|---------------------------------------------------|
| `-target`      | (required)| Target IP, hostname, or CIDR block               |
| `-ports`       | `1-1024`  | Ports to scan, e.g. `22,80,443` or `1-65535`      |
| `-timeout`     | `2s`      | Per-connection timeout                            |
| `-concurrency` | `200`     | Number of concurrent worker goroutines             |

### Examples

Scan common ports on a single host:

```bash
./go-network-scanner -target 192.168.1.10 -ports 22,80,443
```

Sweep a subnet for open ports in the well-known range:

```bash
./go-network-scanner -target 192.168.1.0/24 -ports 1-1024 -concurrency 500
```

Full port sweep with a tighter timeout on a fast local network:

```bash
./go-network-scanner -target 10.0.0.5 -ports 1-65535 -timeout 500ms -concurrency 1000
```

### Sample output

```
Scanning 1 host(s) x 4 port(s)...

192.168.1.10
  22/tcp open  ssh
  80/tcp open  http
  443/tcp open  https

Done in 812ms
```

## How it works

`go-network-scanner` performs a TCP connect scan: for each host/port pair, it attempts a full TCP handshake via `net.DialTimeout`. A successful connection means the port is open; a timeout or refusal means it isn't. Work is distributed across a fixed pool of goroutines via a job channel, and results are collected, filtered to open ports, and sorted before being printed.

This is a full-connect scan, not a stealth SYN scan — it's simpler to implement portably (no raw sockets or elevated privileges needed) at the cost of being more visible to the target and any logging in between.

## Project layout

```
main.go      CLI entry point and output formatting
scan.go      Worker pool and TCP connect logic
ports.go     Port spec parsing (lists and ranges)
hosts.go     Host/CIDR expansion
services.go  Well-known port -> service name lookup
```

## Legal notice

Only scan hosts and networks you own or have explicit permission to test. Unauthorized port scanning may violate computer misuse laws and the acceptable use policies of most networks and cloud providers. The author assumes no liability for misuse of this tool.

## License

MIT
