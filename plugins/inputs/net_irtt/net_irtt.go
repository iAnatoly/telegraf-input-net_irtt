package net_irtt

import (
	"context"
	"time"
	// irtt imports:
	"github.com/heistp/irtt"
	// telegraf imports:
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const measurement = "net_irtt"

type NetIrtt struct {
	RemoteAddress string
	PacketLength  int
}

func init() {
	inputs.Add("net_irtt", func() telegraf.Input {
		return &NetIrtt{
			PacketLength: 100,
		}
	})
}

func (s *NetIrtt) Description() string {
	return "TBD"
}

// SampleConfig returns sample configuration options.
func (s *NetIrtt) SampleConfig() string {
	return `
  ## TBD
`
}

// Gather is the interface for passing metrics to telegraf
func (n *NetIrtt) Gather(acc telegraf.Accumulator) error {

	cfg := irtt.NewClientConfig()

	cfg.LocalAddress = ":0"
	cfg.RemoteAddress = "127.0.0.1:2112"
	cfg.OpenTimeouts, _ = irtt.ParseDurations("1s")
	cfg.Duration, _ = time.ParseDuration("1s")
	cfg.Interval, _ = time.ParseDuration("20ms")
	cfg.Length = 100
	cfg.Clock = irtt.BothClocks
	cfg.IPVersion = irtt.IPVersionFromBooleans(true, true, irtt.DualStack)
	cfg.TTL = 64
	cfg.HMACKey = []byte("wazzup")

	c := irtt.NewClient(cfg)
	ctx := context.Background()
	r, err := c.Run(ctx)

	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"RTTMin":   r.RTTStats.Min,
		"RTTMean":  r.RTTStats.Mean(),
		"RTTMax":   r.RTTStats.Max,
		"IPDVMean": r.RoundTripIPDVStats.Mean(),
		"PLPerc":   r.LatePacketsPercent,
	}

	tags := map[string]string{"RemoteAddress": cfg.RemoteAddress}

	acc.AddFields(measurement, fields, tags)

	return nil
}
