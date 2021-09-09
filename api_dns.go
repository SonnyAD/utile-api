package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

type DNSResolution struct {
	XMLName    xml.Name    `json:"-" xml:"dns" yaml:"-"`
	Type       string      `json:"type" xml:"type" yaml:"type"`
	Resolution interface{} `json:"resolution" xml:"resolution" yaml:"resolution"`
}

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

	var reply DNSResolution
	reply.Type = "dns"
	reply.Resolution = dns

	output(w, r.Header["Accept"], reply, dns.Addresses[0])
}

type DNSResolved struct {
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

	var reply DNSResolution
	reply.Type = "mx"
	reply.Resolution = dns

	defaultOutput := fmt.Sprintf("%s %d", dns.Records[0].Host, dns.Records[0].Pref)

	output(w, r.Header["Accept"], reply, defaultOutput)
}

type MXResolved struct {
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

	var reply DNSResolution
	reply.Type = "ns"
	reply.Resolution = dns

	output(w, r.Header["Accept"], reply, dns.Hosts[0])
}

type NSResolved struct {
	Hosts []string `json:"hosts" xml:"host" yaml:"hosts"`
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

	var reply DNSResolution
	reply.Type = "txt"
	reply.Resolution = dns

	output(w, r.Header["Accept"], reply, dns.Values[0])
}

type TXTResolved struct {
	Values []string `json:"values" xml:"value" yaml:"values"`
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

	var reply DNSResolution
	reply.Type = "cname"
	reply.Resolution = dns

	output(w, r.Header["Accept"], reply, dns.Value)
}

type CNAMEResolved struct {
	Value string `json:"value" xml:"value" yaml:"value"`
}

//endregion

//region CAA resolution
func CAAResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	// NOTE: Adding a dot at the end because the dns library is expecting a FQDN
	domain := mux.Vars(r)["domain"] + "."
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeCAA)
	result, _, err := c.Exchange(m, "1.1.1.1:53")
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}
	if result.Rcode != dns.RcodeSuccess || len(result.Answer) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var answer CAAResolved

	records := make([]CAARecord, len(result.Answer))

	for i, v := range result.Answer {
		if caa, ok := v.(*dns.CAA); ok {
			records[i].Flag = caa.Flag
			records[i].Tag = caa.Tag
			records[i].Value = caa.Value
		}
	}

	answer.Records = records

	var reply DNSResolution
	reply.Type = "caa"
	reply.Resolution = answer

	output(w, r.Header["Accept"], reply, strconv.Itoa((int)(answer.Records[0].Flag))+" "+answer.Records[0].Tag+" "+answer.Records[0].Value)
}

type CAAResolved struct {
	Records []CAARecord `json:"records" xml:"record" yaml:"records"`
}

type CAARecord struct {
	Flag  uint8  `json:"flag" xml:"flag" yaml:"flag"`
	Tag   string `json:"tag" xml:"tag" yaml:"tag"`
	Value string `json:"value" xml:"value" yaml:"value"`
}

//endregion

//region AAAA resolution
func AAAAResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	// NOTE: Adding a dot at the end because the dns library is expecting a FQDN
	domain := mux.Vars(r)["domain"] + "."
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeAAAA)
	result, _, err := c.Exchange(m, "1.1.1.1:53")
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}
	if result.Rcode != dns.RcodeSuccess || len(result.Answer) == 0 {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	var answer AAAAResolved

	hosts := make([]string, len(result.Answer))

	for i, v := range result.Answer {
		if aaaa, ok := v.(*dns.AAAA); ok {
			hosts[i] = aaaa.AAAA.String()
		}
	}

	answer.Hosts = hosts

	var reply DNSResolution
	reply.Type = "aaaa"
	reply.Resolution = answer

	output(w, r.Header["Accept"], reply, answer.Hosts[0])
}

type AAAAResolved struct {
	Hosts []string `json:"hosts" xml:"host" yaml:"hosts"`
}

//endregion

//region DMARC resolution
func DMARCResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := "_dmarc." + mux.Vars(r)["domain"]

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

	var dns DMARCResolved

	dns.Value = txt[0]

	var reply DNSResolution
	reply.Type = "dmarc"
	reply.Resolution = dns

	output(w, r.Header["Accept"], reply, dns.Value)
}

type DMARCResolved struct {
	Value string `json:"value" xml:"value" yaml:"value"`
}

//endregion

//region PTR resolution
func PTRResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	ip := mux.Vars(r)["ip"]

	// NOTE: First convert IP to ARPA Hostname
	arpa, err := dns.ReverseAddr(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// NOTE: Then lookup the ARPA domain PTR record
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(arpa, dns.TypePTR)
	result, _, err := c.Exchange(m, "1.1.1.1:53")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if result.Rcode != dns.RcodeSuccess || len(result.Answer) == 0 {
		http.Error(w, "No results found2", http.StatusNotFound)
		return
	}

	var answer PTRResolved

	domains := make([]string, len(result.Answer))

	for i, v := range result.Answer {
		if ptr, ok := v.(*dns.PTR); ok {
			domains[i] = ptr.Ptr
		}
	}

	answer.Domains = domains

	var reply DNSResolution
	reply.Type = "ptr"
	reply.Resolution = answer

	output(w, r.Header["Accept"], reply, answer.Domains[0])
}

type PTRResolved struct {
	Domains []string `json:"domains" xml:"domain" yaml:"domains"`
}

//endregion
