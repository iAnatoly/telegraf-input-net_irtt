package net_irtt

import (
	"fmt"
	"github.com/heistp/irtt"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/testutil"
	"os"
	"testing"
	"time"
)

const addr string = "127.0.0.1:21112"
const dura time.Duration = 1 * time.Second
const hmac string = "wazzup"

func TestMain(m *testing.M) {
	cfg := irtt.NewServerConfig()

	cfg.Addrs = []string{addr}
	cfg.HMACKey = []byte(hmac)

	s := irtt.NewServer(cfg)

	go func() {
		// shut down in 2x duration
		time.Sleep(2 * dura)
		s.Shutdown()
	}()

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Printf("failed: %s\n", err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func TestGather(t *testing.T) {
	netIrtt := getDefaultConfig()

	netIrtt.Duration = config.Duration(dura)
	netIrtt.HmacKey = hmac
	netIrtt.RemoteAddresses = []string{addr}

	acc := new(testutil.Accumulator)
	err := acc.GatherError(netIrtt.Gather)
	if err != nil {
		t.Errorf("failed: %s\n", err)
	}
}
