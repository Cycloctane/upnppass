package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Cycloctane/upnppass/internal/upnp"
)

const defaultMaxAgeSec = 1800

var version = "dev"

func main() {
	locationStr := flag.String("u", "", "URL of upnp device's root desc xml")
	nicStr := flag.String("i", "", "Network interface for multicast")
	maxAge := flag.Int("t", defaultMaxAgeSec, "Max age of upnp notify in seconds")
	showVersion := flag.Bool("v", false, "Show version")
	flag.Parse()
	if *showVersion {
		fmt.Println(version)
		return
	}
	if *maxAge < defaultMaxAgeSec {
		panic("Max-age should be greater than 1800s")
	}
	location, err := url.Parse(*locationStr)
	if err != nil || !location.IsAbs() {
		panic("Invalid root desc url")
	}

	if *nicStr != "" {
		nic, err := net.InterfaceByName(*nicStr)
		if err != nil {
			panic(err)
		}
		upnp.SetInterface(nic)
	}

	desc, err := upnp.GetDesc(location.String())
	if err != nil {
		panic(err)
	}
	ads, err := upnp.SetupAdvertise(location.String(), desc, defaultMaxAgeSec)
	if err != nil {
		panic(err)
	}

	repeat := time.Tick(time.Duration(defaultMaxAgeSec) * time.Second)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-repeat:
			if !upnp.IsAlive(location.String(), desc.Device.UDN) {
				break loop
			} else {
				if err := ads.NotifyAll(); err != nil {
					break loop
				}
			}
		}
	}

	ads.CloseAll()
}
