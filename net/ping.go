package net

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/tatsushid/go-fastping"
)

// ErrAddrNotResponsive is an error returned when an address is not responsive.
// ErrTimeout is an error returned when a timeout occurs.
var (
	ErrAddrNotResponsive = errors.New("the address is not responsive")
	ErrTimeout           = errors.New("timeout")

	pingPool = &sync.Pool{
		New: func() any {
			return fastping.NewPinger()
		},
	}
)

// Pong represents the result of a ping operation, including the remote
// address, the round-trip time, and any error that occurred.
type Pong struct {
	Addr *net.IPAddr
	RTT  time.Duration
	Err  error
}

// ping sends an ICMP ping to the given address and returns a channel that receives the ping response.
// The ping is sent using the provided fastping.Pinger instance, and the ping is considered successful
// if a response is received within the given maxRTT and timeout durations.
// If the ping fails, the channel will receive a pong with an error.
func Ping(address string, maxRTT, timeout time.Duration) chan *Pong {
	p := pingPool.Get().(*fastping.Pinger)
	defer pingPool.Put(p)

	ch := make(chan *Pong, 1)
	go func() {

		ra, err := net.ResolveIPAddr("ip4:icmp", address)
		if err != nil {
			ra, err = net.ResolveIPAddr("ip6:ipv6-icmp", address)
			if err != nil {
				ch <- &Pong{
					Err: err,
				}
			}
		}
		p.AddIPAddr(ra)
		defer p.RemoveIPAddr(ra)

		p.MaxRTT = maxRTT

		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			ch <- &Pong{
				Addr: addr,
				RTT:  rtt,
			}
		}

		p.OnIdle = func() {
			if p.Err() != nil {
				ch <- &Pong{
					Addr: ra,
					Err:  err,
				}
			} else {
				ch <- &Pong{
					Addr: ra,
					Err:  ErrAddrNotResponsive,
				}
			}
		}

		p.RunLoop()

		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		for {
			select {
			case <-p.Done():
				if err := p.Err(); err != nil {
					ch <- &Pong{
						Addr: ra,
						Err:  err,
					}
				}
			case <-ticker.C:
			}
		}
	}()
	return ch
}
