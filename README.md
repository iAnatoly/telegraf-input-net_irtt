# IRTT Input Plugin

This plugin provides IRTT ([Isochronous Round-Trip Tester[(https://github.com/heistp/irtt)) statistics collection capabilities. 

In current version, this plugin collects RTT (round trip time) and IPDV (jitter) network stats. It can (and should be) expanded to collect more stats, but I would like to keep it as simple as possible for the moment. 

Couple of notes: 
- IRRT requires a server endpoint to talk to. This is trivial to set up on most systems. You can use one server endpoint for multiple clients.
- You will need to specify an HMAC pre-shared key on server side and plugin config, to ensure your IRTT infrastructure won't be abused. This is also trivial. 

## Configuration

TBD

Sample config:
```toml
#
```

## Installation

* Clone the repo

```
git clone 
```
* Build the "net_irtt" binary

```
$ go build -o net_irtt cmd/main.go
```
* You should be able to call this from telegraf now using execd
```
[[inputs.execd]]
  command = ["/path/to/net_irtt", "-poll_interval 1m"]
  signal = "none"
```
This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim)
