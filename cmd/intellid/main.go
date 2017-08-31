package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	nats "github.com/nats-io/go-nats"

	"flag"

	"github.com/AutogrowSystems/go-intelli/device"
	"github.com/AutogrowSystems/go-intelli/util/tell"
)

const enumerationInterval = 10

// Version is the current version of the software: can be set by LDFLAGS at
// build time using: go build -ldflags "-X main.Version=1.0" ./cmd/natsgw
var Version = "version not set"

func main() {
	var delay int
	var natsHost string
	var debug bool
	var apiPort string
	var printVersion bool

	flag.StringVar(&natsHost, "nats", "localhost:4222", "the NATS URL to use")
	flag.StringVar(&apiPort, "p", ":9191", "the API port to serve on")
	flag.BoolVar(&debug, "debug", false, "Run gateway on debug mode")
	flag.BoolVar(&printVersion, "version", false, "print the version and exit")
	flag.IntVar(&delay, "delay", 15, "how often to poll the USB device")
	flag.Parse()

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if !debug {
		tell.Level = tell.INFO
	}

	nc, err := nats.Connect("nats://" + natsHost)
	if err != nil {
		tell.Fatalf("failed to connect to NATS: %s", err)
	}

	mgr := device.NewManager(enumerationInterval, delay)

	// send the shadow over NATS whenever the device shadow is updated
	mgr.OnDeviceUpdated(func(d device.Device) {
		tell.Debugf("device %s updated", d.SerialNumber)

		subj := fmt.Sprintf("intelli.%s", d.SerialNumber)
		data, err := json.Marshal(d.Shadow)
		if err != nil {
			tell.Errorf("failed to send device update over NATS: %s", err)
		}

		nc.Publish(subj, data)
	})

	// start discovering devices attached via USB (loops forever)
	go mgr.Discover()

	// attach and API to the manager to see whats going on
	r := gin.Default()
	mgr.AttachAPI(r)
	go r.Run(apiPort)

	// interrogate the readings from the discovered devices (loops forever)
	mgr.Interrogate()
}
