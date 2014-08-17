package cmd

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func init() {
	AddPlugin("DNS", "(?i)^\\.(dns|dig)$", MessageHandler(DnsLookup), false, false)
}

const DNSLOOKUP = "8.8.4.4:53"

func DnsLookup(msg *Message) {
	var lookup_type string
	var lookup_addr string
	if len(msg.Params) < 2 {
		msg.Return("Usage: .dns [A/AAAA/CNAME/PTR/TXT/SRV] [host]")
	} else if len(msg.Params) == 2 {
		lookup_type = "A"
		lookup_addr = msg.Params[1]
	} else {
		lookup_type = msg.Params[1]
		lookup_addr = msg.Params[2]
	}
	if strings.ToUpper(lookup_type) == "PTR" {
		if addr, err := dns.ReverseAddr(lookup_addr); err == nil {
			LookupHelper(msg, dns.TypePTR, addr)
		} else {
			msg.Return("Invalid PTR address")
			return
		}
	} else if _, isdomain := dns.IsDomainName(lookup_addr); isdomain {
		if querytype, ok := dns.StringToType[strings.ToUpper(lookup_type)]; ok {
			LookupHelper(msg, querytype, lookup_addr)
		}
	}
}

func LookupHelper(msg *Message, ltype uint16, laddr string) {
	m := new(dns.Msg)
	host := dns.Fqdn(laddr)
	m.SetQuestion(host, ltype)
	c := new(dns.Client)
	data, _, err := c.Exchange(m, DNSLOOKUP)
	if err != nil {
		msg.Return("Unspecified Error, sorry.")
		return
	}
	if data == nil {
		msg.Return("I didn't get a response for that query!")
		return
	}
	switch data.MsgHdr.Rcode {
	case dns.RcodeServerFailure:
		msg.Return("SERVFAIL - Blame Google DNS")
	case dns.RcodeNameError:
		msg.Return("NXDOMAIN - Doesn't exist")
	case dns.RcodeRefused:
		msg.Return("REFUSED - Sorry.")
	default:
		data := GetLookupString(data.Answer, ltype)
		out := fmt.Sprintf("%s got %v record(s): %s", laddr, len(data), strings.Join(data, ", "))
		msg.Return(out)
	}
}

func GetLookupString(answer []dns.RR, ltype uint16) []string {
	var output []string
	for _, data := range answer {
		resp := strings.Replace(data.String(), data.Header().String(), "", -1)
		switch ltype {
		case 1, 28, 16, 2, 5, 12:
			output = append(output, resp)
		case 15:
			data := strings.Split(resp, " ")
			output = append(output, fmt.Sprintf("%v [%v]", data[1], data[0]))
		case 6:
			data := strings.Split(resp, " ")
			out := fmt.Sprintf("Primary: %s, Admin: %s, Serial: %s, Refresh: %s, Retry: %s, Expiration: %s, Min: %s", data[0], data[1], data[2], data[3], data[4], data[5], data[6])
			output = append(output, out)
		case 33:
			data := strings.Split(resp, " ")
			out := fmt.Sprintf("Priority: %s, Weight: %s, Port: %s, Target: %s", data[0], data[1], data[2], data[3])
			output = append(output, out)
		}
	}
	return output
}
