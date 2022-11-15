package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	ephemeralproxies "github.com/rpi4gx/ephemeral-proxies-go"
)

const (
	Datacenter  string = "datacenter"
	Residential string = "residential"
)

const (
	Static     string = "static"
	Rotational string = "rotational"
)

type proxyPool struct {
	proxyMode   string
	proxyType   ephemeralproxies.ProxyType
	rapidApiKey string
	size        int
	pool        []ephemeralproxies.Proxy
	mutex       sync.Mutex
	channel     chan string
}

func newProxyPool(proxyType string, proxyMode string, poolSize int, rapidApiKey string) *proxyPool {
	pp := proxyPool{
		proxyMode:   proxyMode,
		size:        poolSize,
		rapidApiKey: rapidApiKey,
		channel:     make(chan string),
	}
	switch proxyType {
	case Datacenter:
		pp.proxyType = ephemeralproxies.Datacenter
	case Residential:
		pp.proxyType = ephemeralproxies.Residential
	default:
		panic("invalid proxytype")
	}
	go pp.populatePool()
	pp.channel <- "go"
	return &pp
}

func (pp *proxyPool) getPoolLength() (length int) {
	pp.mutex.Lock()
	length = len(pp.pool)
	pp.mutex.Unlock()
	return
}

func (pp *proxyPool) populatePool() {
	for {
		<-pp.channel // on signal, fill the pool
		errors := 0
		for pp.getPoolLength() < pp.size {
			if errors >= 10 {
				fmt.Println("too many Ephemeral Proxies API failures, exiting ...")
				os.Exit(-1)
			}
			p, err := ephemeralproxies.NewProxy(pp.rapidApiKey, pp.proxyType)
			if err != nil {
				fmt.Println("failure from Ephemeral Proxies API:", err)
				time.Sleep(100 * time.Millisecond)
				errors++
				continue
			}
			pp.mutex.Lock()
			pp.pool = append(pp.pool, *p)
			pp.mutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (pp *proxyPool) getProxyFromPool() (string, int, error) {
	defer func() {
		go func(pp *proxyPool) {
			pp.channel <- "go"
		}(pp)
	}()
	if pp.getPoolLength() < 1 {
		return "", 0, errors.New("warning: proxy pool exahusted. Try to increase by using -size option")
	}
	pp.mutex.Lock()
	var host = pp.pool[0].Host
	var port = pp.pool[0].Port
	if pp.proxyMode == Rotational {
		pp.pool = pp.pool[1:]
	}
	pp.mutex.Unlock()
	return host, port, nil
}

func (pp *proxyPool) String() string {
	return "Pool of size: " + strconv.Itoa(pp.size)
}
