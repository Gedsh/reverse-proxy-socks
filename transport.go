package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

var socksPort = flag.Int("sockport", defaultSocksPort, "Socks5 proxy port")

func getTransport() *http.Transport {

	addr := localHost + ":" + strconv.Itoa(*socksPort)

	baseDialer := &net.Dialer{}

	dialer, err := proxy.SOCKS5("tcp", addr, nil, baseDialer)
	if err != nil {
		log.Fatalf("SOCKS5 dialer error: %v", err)
	}

	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

			if cd, ok := dialer.(proxy.ContextDialer); ok {
				return cd.DialContext(ctx, network, addr)
			}

			return dialer.Dial(network, addr)
		},
		DisableKeepAlives:   false,
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 30 * time.Second,
	}
}
