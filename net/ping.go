package net

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/tatsushid/go-fastping"
)

var (
	ErrAddrNotResponsive = errors.New("the address is not responsive")
	ErrTimeout           = errors.New("timeout")
	pingPool             = &sync.Pool{
		New: func() any {
			return fastping.NewPinger()
		},
	}
)

// Pong holds the result of a ping operation.
type Pong struct {
	Addr *net.IPAddr
	RTT  time.Duration
	Err  error
}

// Ping sends an ICMP ping and returns a channel receiving the response.
// Optimization: Uses pooled Pinger instances to reduce allocations.
func Ping(address string, maxRTT, timeout time.Duration) chan *Pong {
	p := pingPool.Get().(*fastping.Pinger)
	defer pingPool.Put(p)
	ch := make(chan *Pong, 1)
	go func() {
		ra, err := net.ResolveIPAddr("ip4:icmp", address)
		if err != nil {
			ra, err = net.ResolveIPAddr("ip6:ipv6-icmp", address)
			if err != nil {
				ch <- &Pong{Err: err}
				return
			}
		}
		p.AddIPAddr(ra)
		defer p.RemoveIPAddr(ra)
		p.MaxRTT = maxRTT
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			ch <- &Pong{Addr: addr, RTT: rtt}
		}
		p.OnIdle = func() {
			if p.Err() != nil {
				ch <- &Pong{Addr: ra, Err: p.Err()}
			} else {
				ch <- &Pong{Addr: ra, Err: ErrAddrNotResponsive}
			}
		}
		p.RunLoop()
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		select {
		case <-p.Done():
			if err := p.Err(); err != nil {
				ch <- &Pong{Addr: ra, Err: err}
			}
		case <-ticker.C:
			ch <- &Pong{Addr: ra, Err: ErrTimeout}
		}
	}()
	return ch
}
