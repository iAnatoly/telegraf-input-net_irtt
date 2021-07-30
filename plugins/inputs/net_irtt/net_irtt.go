package net_irtt

import (
	"context"
	"fmt"
	"time"
	// irtt imports:
	"github.com/heistp/irtt"
	// telegraf imports:
	"github.com/influxdata/telegraf"
)

const measurement = "irtt_stats"

type NetIrtt struct {
	RemoteAddress string
	PacketLength  int
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
	fmt.Printf("RTT min/mean/max: %s/%s/%s\n", r.RTTStats.Min, r.RTTStats.Mean(), r.RTTStats.Max)
	fmt.Printf("jitter min/mean/max: %s/%s/%s\n", r.RoundTripIPDVStats.Min, r.RoundTripIPDVStats.Mean(), r.RoundTripIPDVStats.Max)
	fmt.Printf("late packet count/percent: %d/%f\n", r.LatePackets, r.LatePacketsPercent)
	fmt.Printf("lost packet percent: %f\n", r.PacketLossPercent)

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
