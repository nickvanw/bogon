package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/miekg/dns"
	"github.com/nickvanw/bogon/commands"
)

const (
	dnsServer = "8.8.4.4:53"
)

var dnsCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(dns|dig)$")
	return dnsTitle, out, dnsLookup, defaultOptions
}

func dnsLookup(msg commands.Message, ret commands.MessageFunc) string {
	var lookupType, lookupAddr string
	if len(msg.Params) < 2 {
		return "Usage: .dns [A/AAAA/CNAME/PTR/TXT/SRV] [host]"
	} else if len(msg.Params) == 2 {
		lookupType = "A"
		lookupAddr = msg.Params[1]
	} else {
		lookupType = msg.Params[1]
		lookupAddr = msg.Params[2]
	}
	if strings.ToUpper(lookupType) == "PTR" {
		if addr, err := dns.ReverseAddr(lookupAddr); err == nil {
			return lookupHelper(msg, ret, dns.TypePTR, addr)
		}
		return "Invalid PTR address"
	} else if _, isdomain := dns.IsDomainName(lookupAddr); isdomain {
		if querytype, ok := dns.StringToType[strings.ToUpper(lookupType)]; ok {
			return lookupHelper(msg, ret, querytype, lookupAddr)
		}
	}
	return ""
}

func lookupHelper(msg commands.Message, ret commands.MessageFunc, ltype uint16, laddr string) string {
	m := new(dns.Msg)
	host := dns.Fqdn(laddr)
	m.SetQuestion(host, ltype)
	c := new(dns.Client)
	data, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		return "Unspecified Error, sorry."
	}
	if data == nil {
		return "I didn't get a response for that query!"
	}
	switch data.MsgHdr.Rcode {
	case dns.RcodeServerFailure:
		return "SERVFAIL - Blame Google DNS"
	case dns.RcodeNameError:
		return "NXDOMAIN - Doesn't exist"
	case dns.RcodeRefused:
		return "REFUSED - Sorry."
	default:
		data := lookupString(data.Answer, ltype)
		return fmt.Sprintf("%s got %v record(s): %s", laddr, len(data), strings.Join(data, ", "))
	}
}

func lookupString(answer []dns.RR, ltype uint16) []string {
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
