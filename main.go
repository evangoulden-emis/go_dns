package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	// Define a command-line flag for FQDNs
	fqdnsInput := flag.String("fqdns", "", "Comma-separated list of FQDNs to query")
	dnsServer := flag.String("server", "8.8.8.8:53", "DNS server to query")
	recursionDesired := flag.Bool("recursion", false, "Whether to enable recursion")
	flag.Parse()

	if *fqdnsInput == "" {
		fmt.Println("Please provide a list of FQDNs using -fqdns")
		return
	}

	fqdns := strings.Split(*fqdnsInput, ",")

	for _, fqdn := range fqdns {
		queryDNS(strings.TrimSpace(fqdn), *dnsServer, *recursionDesired)
	}
}

func queryDNS(fqdn, server string, recursionDesired bool) {
	c := new(dns.Client)
	m := new(dns.Msg)

	recordTypes := []uint16{dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX, dns.TypeNS, dns.TypeSOA, dns.TypeTXT}
	for _, recordType := range recordTypes {
		m.SetQuestion(dns.Fqdn(fqdn), recordType)
		m.RecursionDesired = recursionDesired
		r, _, err := c.Exchange(m, server)
		if err != nil {
			fmt.Printf("Error querying %s: %v\n", fqdn, err)
			return
		}

		if r.MsgHdr.RecursionAvailable {
			fmt.Printf("✅ Recursion is ENABLED on %s\n", server)
		} else {
			fmt.Printf("❌ Recursion is DISABLED on %s\n", server)
		}

		fmt.Printf("🔎 %s results for %s:\n", dns.TypeToString[recordType], fqdn)
		for _, ans := range r.Answer {
			fmt.Println(ans)
		}
		fmt.Println()
	}
}
