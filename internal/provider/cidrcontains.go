package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = CIDRContainsFunction{}
)

func NewCIDRContainsFunction() function.Function {
	return CIDRContainsFunction{}
}

type CIDRContainsFunction struct{}

func (f CIDRContainsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cidrcontains"
}

func (f CIDRContainsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "CIDRContains function",
		MarkdownDescription: "Determines if a CIDR prefix contains a given IP address",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "prefix",
				MarkdownDescription: "The CIDR prefix to check",
			},
			function.StringParameter{
				Name:                "address",
				MarkdownDescription: "The IP address to check",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f CIDRContainsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var address string
	var cidr string
	err := req.Arguments.Get(ctx, &cidr, &address)

	resp.Error = function.ConcatFuncErrors(err)
	if resp.Error != nil {
		return
	}

	contains, err := cidrContains(address, cidr)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	err = resp.Result.Set(ctx, contains)
	resp.Error = function.ConcatFuncErrors(err)
}

func cidrContains(address string, cidr string) (bool, *function.FuncError) {
	ip := net.ParseIP(address)
	if ip == nil {
		return false, function.NewArgumentFuncError(1, "invalid address format")
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, function.NewArgumentFuncError(0, fmt.Sprintf("invalid CIDR format: %s", err))
	}

	// Both need to be in the same family (IPv4 or IPv6)
	if ip.To4() != nil && ipnet.IP.To4() == nil {
		return false, function.NewArgumentFuncError(0, "address is IPv4, but CIDR is IPv6")
	} else if ip.To4() == nil && ipnet.IP.To4() != nil {
		return false, function.NewArgumentFuncError(0, "address is IPv6, but CIDR is IPv4")
	}

	return ipnet.Contains(ip), nil
}
