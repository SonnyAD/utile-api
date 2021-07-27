package main

import (
	"context"
	"encoding/xml"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func DNSResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	ip, err := resolver.LookupHost(context.Background(), domain)

	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns DNSResolved
	dns.Addresses = ip

	output(w, r.Header["Accept"], dns, dns.Addresses[0])
}

type DNSResolved struct {
	XMLName   xml.Name `xml:"dnsresolution"`
	Addresses []string `json:"addresses" xml:"addresses" yaml:"addresses"`
}
