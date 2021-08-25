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

	records := make([]MXRecord, len(mx))

	for i, v := range mx {
		records[i].Host = v.Host
		records[i].Pref = v.Pref
	}

	dns.Records = records

	defaultOutput := fmt.Sprintf("%s %d", dns.Records[0].Host, dns.Records[0].Pref)

	output(w, r.Header["Accept"], dns, defaultOutput)
}

type MXResolved struct {
	XMLName xml.Name   `json:"-" xml:"mxresolution" yaml:"-"`
	Records []MXRecord `json:"records" xml:"record" yaml:"records"`
}

type MXRecord struct {
	Host string `json:"host" xml:"host" yaml:"host"`
	Pref uint16 `json:"pref" xml:"pref" yaml:"pref"`
}

//endregion

//region NS resolution
func NSResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	// Details: https://pkg.go.dev/net#Resolver.LookupNS
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	ns, err := resolver.LookupNS(context.Background(), domain)

	if err != nil || len(ns) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns NSResolved

	hosts := make([]string, len(ns))

	for i, v := range ns {
		hosts[i] = v.Host
	}

	dns.Hosts = hosts

	output(w, r.Header["Accept"], dns, dns.Hosts[0])
}

type NSResolved struct {
	XMLName xml.Name `json:"-" xml:"nsresolution" yaml:"-"`
	Hosts   []string `json:"hosts" xml:"host" yaml:"hosts"`
}

//endregion

//region TXT resolution
func TXTResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	// Details: https://pkg.go.dev/net#Resolver.LookupTXT
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	txt, err := resolver.LookupTXT(context.Background(), domain)

	if err != nil || len(txt) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns TXTResolved

	dns.Values = txt

	output(w, r.Header["Accept"], dns, dns.Values[0])
}

type TXTResolved struct {
	XMLName xml.Name `json:"-" xml:"nsresolution" yaml:"-"`
	Values  []string `json:"values" xml:"value" yaml:"values"`
}

//endregion

//region CNAME resolution
func CNAMEResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	// Details: https://pkg.go.dev/net#Resolver.LookupCNAME
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	cname, err := resolver.LookupCNAME(context.Background(), domain)

	if err != nil || len(cname) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var dns CNAMEResolved

	dns.Value = cname

	output(w, r.Header["Accept"], dns, dns.Value)
}

type CNAMEResolved struct {
	XMLName xml.Name `json:"-" xml:"nsresolution" yaml:"-"`
	Value   string   `json:"value" xml:"value" yaml:"value"`
}

//endregion
