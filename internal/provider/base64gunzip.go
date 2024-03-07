package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = Base64GunzipFunction{}
)

func NewBase64GunzipFunction() function.Function {
	return Base64GunzipFunction{}
}

type Base64GunzipFunction struct{}

func (f Base64GunzipFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "base64gunzip"
}

func (f Base64GunzipFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Base64gunzip function",
		MarkdownDescription: "Decompresses a base64-encoded gzip string",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "str",
				MarkdownDescription: "The base64-encoded gzip string to decompress and decode",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f Base64GunzipFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var str string
	err := req.Arguments.Get(ctx, &str)

	resp.Error = function.ConcatFuncErrors(err)
	if resp.Error != nil {
		return
	}

	decoded, decodeErr := base64.StdEncoding.DecodeString(str)
	if decodeErr != nil {
		resp.Error = function.NewFuncError(decodeErr.Error())
		return
	}

	buffer := bytes.NewReader(decoded)
	gzipReader, readerErr := gzip.NewReader(buffer)
	if readerErr != nil {
		resp.Error = function.NewFuncError(readerErr.Error())
		return
	}
	defer gzipReader.Close()

	gunzipped, readErr := io.ReadAll(gzipReader)
	if err != nil {
		resp.Error = function.NewFuncError(readErr.Error())
		return
	}

	err = resp.Result.Set(ctx, string(gunzipped))
	resp.Error = function.ConcatFuncErrors(err)
}
