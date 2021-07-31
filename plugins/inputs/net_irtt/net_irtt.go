package net_irtt

import (
	"context"
	"fmt"
	"time"
	// irtt imports:
	"github.com/heistp/irtt"
	// telegraf imports:
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const measurement = "net_irtt"

// TODO: add all the irtt parameters
type NetIrtt struct {
	RemoteAddress string
	HmacKey       string
	Duration      config.Duration
	Interval      config.Duration
	PacketLength  int
	LocalAddress  string
	OpenTimeouts  []config.Duration
	Ipv4          bool
	Ipv6          bool
	Ttl           int
}

func init() {
	// TODO provid emeningful defaults
	inputs.Add("net_irtt", func() telegraf.Input {
		return &NetIrtt{
			PacketLength: 100,
		}
	})
}

func (s *NetIrtt) Description() string {
	return "Provide Isochronous Round-Trip Tester stats"
}

// SampleConfig returns sample configuration options.
func (s *NetIrtt) SampleConfig() string {
	// TODO: proivide an example
	return `
  ## net_irtt parameters
  
  # these ones you probably want to adjust

  remoteaddress = "remote_server:2112"
  hmackey = "verylongpassphrase"
  duration = "5s"

  # send packets every 20ms, 100b payload
  # very similar to RTP

  interval = "20ms"
  packetlength = 100

  # override if needed 

  localaddress = ":0"
  opentimeouts = "1s"
  ipv4 = true
  ipv6 = false
  ttl = 64

  ## uncomment to remove unneeded fields
  # fielddrop = [ "RTTMin", "IPDVMin" ]
`
}

func (n *NetIrtt) getClientConfig() *irtt.ClientConfig {
	cfg := irtt.NewClientConfig()

	cfg.LocalAddress = n.LocalAddress
	cfg.RemoteAddress = n.RemoteAddress
	cfg.OpenTimeouts = func(ts []config.Duration) []time.Duration {
		r := make([]time.Duration, len(ts))
		for i := range ts {
			r[i] = time.Duration(ts[i])
		}
		return r
	}(n.OpenTimeouts)
	cfg.Duration = time.Duration(n.Duration)
	cfg.Interval = time.Duration(n.Interval)
	cfg.Length = n.PacketLength
	cfg.Clock = irtt.BothClocks
	cfg.IPVersion = irtt.IPVersionFromBooleans(n.Ipv4, n.Ipv6, irtt.DualStack)
	cfg.TTL = n.Ttl
	cfg.HMACKey = []byte(n.HmacKey)

	return cfg
}

// Gather is the interface for passing metrics to telegraf
func (n *NetIrtt) Gather(acc telegraf.Accumulator) error {

	cfg := n.getClientConfig()
	c := irtt.NewClient(cfg)
	ctx := context.Background() // TODO: add signal handling
	r, err := c.Run(ctx)

	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"RTTMin":   r.RTTStats.Min.Microseconds(),
		"RTTMean":  r.RTTStats.Mean().Microseconds(),
		"RTTMax":   r.RTTStats.Max.Microseconds(),
		"IPDVMean": r.RoundTripIPDVStats.Mean().Microseconds(),
		"IPDVMin":  r.RoundTripIPDVStats.Min.Microseconds(),
		"IPDVMax":  r.RoundTripIPDVStats.Max.Microseconds(),
		"PLPerc":   r.LatePacketsPercent,
	}

	tags := map[string]string{"RemoteAddress": cfg.RemoteAddress}

	acc.AddFields(measurement, fields, tags)

	return nil
}
