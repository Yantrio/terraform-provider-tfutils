package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = URLDecodeFunction{}
)

func NewURLDecodeFunction() function.Function {
	return URLDecodeFunction{}
}

type URLDecodeFunction struct{}

func (U URLDecodeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "urldecode"
}

func (U URLDecodeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "URLDecode function",
		MarkdownDescription: "Decodes a URL-encoded string",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "input",
				MarkdownDescription: "The URL-encoded string to decode",
			},
		},
		Return: function.StringReturn{},
	}
}

func (U URLDecodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string
	err := req.Arguments.Get(ctx, &data)

	resp.Error = function.ConcatFuncErrors(err)
	if resp.Error != nil {
		return
	}

	// Decode the URL-encoded string
	decoded, decodeErr := url.QueryUnescape(data)
	if decodeErr != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("failed to decode URL-encoded string: %s", decodeErr.Error()))
		return
	}

	err = resp.Result.Set(ctx, decoded)
	resp.Error = function.ConcatFuncErrors(err)
}
