# IRTT Input Plugin

## Summary

This plugin provides IRTT ([Isochronous Round-Trip Tester[(https://github.com/heistp/irtt)) statistics collection capabilities. 

In current version, this plugin collects RTT (round trip time) and IPDV (jitter) network stats. It can (and should be) expanded to collect more stats, but I would like to keep it as simple as possible for the moment. 

Couple of notes: 
- IRRT requires a server endpoint to talk to. This is trivial to set up on most systems. You can use one server endpoint for multiple clients.
- You will need to specify an HMAC pre-shared key on server side and plugin config, to ensure your IRTT infrastructure won't be abused. This is also trivial. 

## Configuration

Sample config (see plugin.conf in the repo):
```toml
[[inputs.net_irtt]]

  ## these ones you probably want to adjust.
  ## irtt server should be listening on remote_address, with the same hmac_key configured

  remote_addresses = [ "127.0.0.1:2112", "192.168.1.1:2112" ]
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

```

## Installation

* Clone the repo `git clone ...` 
* Build the "net_irtt" binary:
    * for AMD64: `$ go build -o net_irtt cmd/main.go`
    * for ARM64 (Raspberry Pi 4): `$ env GOOS=linux GOARCH=arm64 go build -o net_irtt.arm64 cmd/main.go`
    * for ARMv7l (Raspberry Pi 3b): `$ env GOOS=linux GOARCH=arm GOARM=7 go build -o net_irtt.armv7l cmd/main.go`
    * for ARMv6l (Raspberry Pi Zero): `$ env GOOS=linux GOARCH=arm GOARM=6 go build -o net_irtt.armv6l cmd/main.go`
* Edit the config: `vi plugin.config` 
* Copy the binary and the config to an appropriate location
```bash
$ sudo cp plugin.config /etc/telegraf/telegraf-irtt.config
$ sudo cp net_irtt /usr/lib/telegraf/plugins/
```
* You should be able to call this from telegraf now using execd:
```
[[inputs.execd]]
  command = ["/usr/lib/telegraf/plugins/net_irtt", "-config", "/etc/telegraf/telegraf-irtt.config" ]
  signal = "none"
```
## Credits
* This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim).
* This plugin is notitng but a wrapper for beautiful [IRTT](https://github.com/heistp/irtt). Credit for the ingenuity goes to the authors.
* The original idea to create a telegraf plugin for irtt came from [@nvitaly](https://github.com/nvitaly).
