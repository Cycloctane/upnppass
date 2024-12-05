package main

import (
	"flag"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Cycloctane/upnppass/internal/upnp"
)

const defaultMaxAgeSeconds = 1800

func main() {
	locationStr := flag.String("u", "", "URL of upnp device's root desc xml")
	nicStr := flag.String("i", "", "Network interface for multicast")
	maxAge := flag.Int("t", defaultMaxAgeSeconds, "Max age of upnp notify in seconds")
	flag.Parse()
	if *maxAge < defaultMaxAgeSeconds {
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
	ads, err := upnp.SetupAdvertise(location.String(), desc, defaultMaxAgeSeconds)
	if err != nil {
		panic(err)
	}

	repeat := time.Tick(time.Duration(defaultMaxAgeSeconds) * time.Second)
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
