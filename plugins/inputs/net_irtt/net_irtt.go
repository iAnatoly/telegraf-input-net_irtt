package net_irtt

import (
	"context"
	"time"
	// irtt imports:
	"github.com/heistp/irtt"
	// telegraf imports:
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const measurement = "net_irtt"
const input = "net_irtt"

// TODO: add all the irtt parameters
type NetIrtt struct {
	RemoteAddresses []string `toml:"remote_addresses"`
	HmacKey         string   `toml:"hmac_key"`
	Duration        config.Duration
	Interval        config.Duration
	PacketLength    int               `toml:"packet_length"`
	LocalAddress    string            `toml:"local_address"`
	OpenTimeouts    []config.Duration `toml:"open_timeouts"`
	Ipv4            bool
	Ipv6            bool
	Ttl             int
}

func getDefaultConfig() *NetIrtt {
	return &NetIrtt{
		Duration:     config.Duration(time.Second * 5),
		Interval:     config.Duration(time.Millisecond * 20),
		PacketLength: 100,
		LocalAddress: ":0",
		OpenTimeouts: []config.Duration{config.Duration(time.Second * 1)},
		Ipv4:         true,
		Ipv6:         false,
		Ttl:          64,
	}
}

func init() {
	inputs.Add(input, func() telegraf.Input {
		return getDefaultConfig()
	})
}

func (s *NetIrtt) Description() string {
	return "Provide Isochronous Round-Trip Tester stats"
}

// SampleConfig returns sample configuration options.
func (s *NetIrtt) SampleConfig() string {
	return `
  ## these ones you probably want to adjust.
  ## irtt server should be listening on remote_address, with the same hmac_key configured

  remote_addresses = [ "127.0.0.1:2112", "192.168.1.2:2112" ]
  hmac_key = "wazzup"

  ## run the test for 5s
  duration = "5s"

  ## send packets every 20ms, 100b payload
  ## very similar to RTP

  interval = "20ms"
  # packet_length = 100

  ## override as needed

  local_address = ":0"
  open_timeouts = ["1s"]
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
	// everything but:
	// cfg.RemoteAddress = n.RemoteAddress
	// because we handle more than one remote address with the same set of parameters
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

	for _, server := range n.RemoteAddresses {
		cfg.RemoteAddress = server
		c := irtt.NewClient(cfg)
		ctx := context.Background()
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

		tags := map[string]string{"remote_address": server}

		acc.AddFields(measurement, fields, tags)
	}

	return nil
}
