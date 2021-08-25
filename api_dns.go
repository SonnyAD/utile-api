package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//region DNS resolution
func DNSResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	// Details: https://pkg.go.dev/net#Resolver.LookupHost
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

	if err != nil || len(ip) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns DNSResolved
	dns.Addresses = ip

	output(w, r.Header["Accept"], dns, dns.Addresses[0])
}

type DNSResolved struct {
	XMLName   xml.Name `json:"-" xml:"dnsresolution" yaml:"-"`
	Addresses []string `json:"addresses" xml:"address" yaml:"addresses"`
}

//endregion

//region MX resolution
func MXResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	// Details: https://pkg.go.dev/net#Resolver.LookupMX
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	mx, err := resolver.LookupMX(context.Background(), domain)

	if err != nil || len(mx) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns MXResolved

	records := make([]Record, len(mx))

	for i, v := range mx {
		records[i].Host = v.Host
		records[i].Pref = v.Pref
	}

	dns.Records = records

	defaultOutput := fmt.Sprintf("%s %d", dns.Records[0].Host, dns.Records[0].Pref)

	output(w, r.Header["Accept"], dns, defaultOutput)
}

type MXResolved struct {
	XMLName xml.Name `json:"-" xml:"mxresolution" yaml:"-"`
	Records []Record `json:"records" xml:"record" yaml:"records"`
}

type Record struct {
	Host string `json:"host" xml:"host" yaml:"host"`
	Pref uint16 `json:"pref" xml:"pref" yaml:"pref"`
}

//endregion
