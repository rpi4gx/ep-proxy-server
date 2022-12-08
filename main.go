package main

import (
	"os"

	"github.com/namsral/flag"
	"github.com/op/go-logging"
)

func usage() {
	flag.Usage()
	os.Exit(2)
}

var log = logging.MustGetLogger("")

func main() {
	var listeningPort int
	var proxyType string
	var proxyMode string
	var proxyModeRotationalPoolSize int
	var rapidApiKey string
	var debugMode bool
	var loggerStdout = logging.NewLogBackend(os.Stdout, "", 0)

	flag.IntVar(&listeningPort, "port", 8090, "Listening port")
	flag.StringVar(&proxyType, "type", "datacenter", "Type of proxy <datacenter|residential>.")
	flag.StringVar(&proxyMode, "mode", "static", "Mode <static|rotational>.")
	flag.IntVar(&proxyModeRotationalPoolSize, "size", 50, "Pool size for rotational mode.")
	flag.StringVar(&rapidApiKey, "rapidApiKey", "", "Rapid API key.")
	flag.BoolVar(&debugMode, "debug", false, "Set to true to enable debug mode")
	// TODO: allow to select country
	flag.Parse()

	if proxyType != Datacenter && proxyType != Residential {
		usage()
	}

	if len(rapidApiKey) < 1 {
		usage()
	}

	if debugMode {
		logging.AddModuleLevel(loggerStdout).SetLevel(logging.DEBUG, "")
	} else {
		logging.AddModuleLevel(loggerStdout).SetLevel(logging.INFO, "")
	}
	logging.SetBackend(loggerStdout)
	log.Infof("debug mode: %v\n", debugMode)
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
	log.Info("Using Ephemeral Proxies API with", proxyType, "proxies in", proxyMode, "mode")
	serverStart(listeningPort, proxyPool)
}
