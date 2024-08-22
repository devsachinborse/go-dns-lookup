package main

import (
	"github.com/miekg/dns"
)

type DNSRecord struct {
	DomainType string
	DomainName string
	IPAddress  string
	TTL        uint32
	RecordType string
}

func LookupDNS(domain string) ([]DNSRecord, error) {
	var records []DNSRecord

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, err
	}

	for _, answer := range r.Answer {
		switch answer.(type) {
		case *dns.A:
			aRecord := answer.(*dns.A)
			records = append(records, DNSRecord{
				DomainType: "A",
				DomainName: domain,
				IPAddress:  aRecord.A.String(),
				TTL:        aRecord.Hdr.Ttl,
				RecordType: "A",
			})
		case *dns.AAAA:
			aaaaRecord := answer.(*dns.AAAA)
			records = append(records, DNSRecord{
				DomainType: "AAAA",
				DomainName: domain,
				IPAddress:  aaaaRecord.AAAA.String(),
				TTL:        aaaaRecord.Hdr.Ttl,
				RecordType: "AAAA",
			})
		case *dns.CNAME:
			cnameRecord := answer.(*dns.CNAME)
			records = append(records, DNSRecord{
				DomainType: "CNAME",
				DomainName: domain,
				IPAddress:  cnameRecord.Target,
				TTL:        cnameRecord.Hdr.Ttl,
				RecordType: "CNAME",
			})
		case *dns.NS:
			nsRecord := answer.(*dns.NS)
			records = append(records, DNSRecord{
				DomainType: "NS",
				DomainName: domain,
				IPAddress:  nsRecord.Ns,
				TTL:        nsRecord.Hdr.Ttl,
				RecordType: "NS",
			})
		case *dns.MX:
			mxRecord := answer.(*dns.MX)
			records = append(records, DNSRecord{
				DomainType: "MX",
				DomainName: domain,
				IPAddress:  mxRecord.Mx,
				TTL:        mxRecord.Hdr.Ttl,
				RecordType: "MX",
			})
		}
	}

	return records, nil
}
