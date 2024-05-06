package scan

import (
	"fmt"
	"net"
	"time"
)

type PortState struct {
	Port int
	Open state
}

// Result represent the san reuslts fo a single host
type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

type state bool

func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true

	return p
}

func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}
