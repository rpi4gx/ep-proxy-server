package main

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

func usage() {
	flag.Usage()
	os.Exit(2)
}

func main() {
	var listeningPort int
	var proxyType string
	var proxyMode string
	var proxyModeRotationalPoolSize int
	var rapidApiKey string

	flag.IntVar(&listeningPort, "port", 8090, "Listening port")
	flag.StringVar(&proxyType, "type", "datacenter", "Type of proxy <datacenter|residential>.")
	flag.StringVar(&proxyMode, "mode", "static", "Mode <static|rotational>.")
	flag.IntVar(&proxyModeRotationalPoolSize, "size", 50, "Pool size for rotational mode.")
	flag.StringVar(&rapidApiKey, "rapidApiKey", "", "Rapid API key.")
	// TODO: allow to select country
	flag.Parse()

	if proxyType != Datacenter && proxyType != Residential {
		usage()
	}

	if len(rapidApiKey) < 1 {
		usage()
	}

	if proxyMode == Static {
		proxyModeRotationalPoolSize = 1
	} else if proxyMode == Rotational {
		if proxyModeRotationalPoolSize < 0 {
			usage()
		}
	} else {
		usage()
	}

	proxyPool := newProxyPool(proxyType, proxyMode, proxyModeRotationalPoolSize, rapidApiKey)
	fmt.Println("Using Ephemeral Proxies API with", proxyType, "proxies in", proxyMode, "mode")
	serverStart(listeningPort, proxyPool)
}
