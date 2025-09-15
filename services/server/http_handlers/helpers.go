package httphandlers

import rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"

func BuildRules(hosts, ips, regexps []string) *rasp_rpc.NewRules {
	payload := &rasp_rpc.NewRules_Rules{
		Rules: &rasp_rpc.UpdatedSSRFRules{
			HostRules:   hosts,
			IPRules:     ips,
			RegexpRules: regexps,
		},
	}
	return &rasp_rpc.NewRules{
		Payload: payload,
	}
}
