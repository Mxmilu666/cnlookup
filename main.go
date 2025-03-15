package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

// 这谁写的代码啊，写的这么屎，建议删库跑路

// Custom DNS resolver
type customResolver struct {
	server string
}

func (r *customResolver) lookupCNAME(domain string) (string, error) {
	if r.server == "" {
		return net.LookupCNAME(domain)
	}

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", r.server)
		},
	}
	return resolver.LookupCNAME(context.Background(), domain)
}

func (r *customResolver) lookupIP(domain string) ([]net.IP, error) {
	if r.server == "" {
		addrs, err := net.LookupIP(domain)
		return addrs, err
	}

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", r.server)
		},
	}
	ipAddrs, err := resolver.LookupIPAddr(context.Background(), domain)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, len(ipAddrs))
	for i, addr := range ipAddrs {
		ips[i] = addr.IP
	}
	return ips, nil
}

func lookupCNAME(domain string, dnsServer string) {
	resolver := &customResolver{server: dnsServer}
	nestLevel := 0
	var domains []string // Track queried domains

	originalDomain := domain

	// Print header and DNS server info
	fmt.Printf("CNAME Lookup: %s\n", domain)
	if dnsServer != "" {
		fmt.Printf("Using DNS server: %s\n", dnsServer)
	}
	fmt.Println(strings.Repeat("-", 50))

	// Check if domain is Punycode
	unicodeDomain, err := idna.ToUnicode(domain)
	if err == nil && unicodeDomain != domain {
		fmt.Printf("Domain Unicode form: %s\n", unicodeDomain)
		fmt.Println(strings.Repeat("-", 50))
	}

	fmt.Println("CNAME resolution chain:")
	for {
		nestLevel++
		domains = append(domains, domain)

		cname, err := resolver.lookupCNAME(domain)
		if err != nil {
			fmt.Printf("Error: Unable to query CNAME record for %s: %v\n", domain, err)
			return
		}

		// Clean trailing dot from CNAME
		cleanCname := strings.TrimRight(cname, ".")
		cleanDomain := strings.TrimRight(domain, ".")

		// Check for Punycode
		unicodeCNAME, err := idna.ToUnicode(cleanCname)

		if cleanCname == cleanDomain {
			// Domain points to itself, no CNAME record
			fmt.Printf("  [%d] %s (No CNAME record)\n", nestLevel, domain)
			break
		} else {
			if err == nil && unicodeCNAME != cleanCname {
				fmt.Printf("  [%d] %s → %s (Unicode: %s)\n", nestLevel, domain, cleanCname, unicodeCNAME)
			} else {
				fmt.Printf("  [%d] %s → %s\n", nestLevel, domain, cleanCname)
			}
		}

		// Check for CNAME loops
		if contains(domains, cleanCname) {
			fmt.Println("Warning: CNAME circular reference detected")
			break
		}

		domain = cleanCname
	}

	fmt.Println(strings.Repeat("-", 50))

	// Query IP addresses for the final domain
	ips, err := resolver.lookupIP(domain)
	if err != nil {
		fmt.Printf("Error: Unable to resolve IP addresses for %s: %v\n", domain, err)
	} else {
		fmt.Println("IP addresses:")
		ipv4Count := 0
		ipv6Count := 0

		for _, ip := range ips {
			if ip.To4() != nil {
				fmt.Printf("  IPv4: %s\n", ip.String())
				ipv4Count++
			} else {
				fmt.Printf("  IPv6: %s\n", ip.String())
				ipv6Count++
			}
		}

		fmt.Printf("\nFound %d IPv4 address(es) and %d IPv6 address(es)\n", ipv4Count, ipv6Count)
	}

	// Display summary statistics
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Resolution summary: %s -> %s (%d level CNAME chain)\n",
		originalDomain, domain, nestLevel)
}

// Check if slice contains item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimRight(s, ".") == item {
			return true
		}
	}
	return false
}

func main() {
	var domain string
	var dnsServer string

	flag.StringVar(&domain, "d", "", "Domain to lookup")
	flag.StringVar(&dnsServer, "s", "", "DNS server (format: IP:PORT, e.g., 8.8.8.8:53)")
	flag.Parse()

	// Try to get domain
	if domain == "" && flag.NArg() > 0 {
		domain = flag.Arg(0)
	}

	if domain == "" {
		fmt.Println("Usage: cnlookup -d example.com [-s 8.8.8.8:53]")
		fmt.Println("   or: cnlookup example.com [-s 8.8.8.8:53]")
		os.Exit(1)
	}

	// Add DNS Server default port
	if dnsServer != "" && !strings.Contains(dnsServer, ":") {
		dnsServer = dnsServer + ":53"
	}

	lookupCNAME(domain, dnsServer)
}
