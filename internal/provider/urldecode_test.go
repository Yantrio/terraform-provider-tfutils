package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestURLDecodeFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		wanted string
	}{
		{
			name:   "URL with spaces encoded as plus signs",
			input:  "\"https%3A%2F%2Fwww.example.com%2Fsearch%3Fquery%3Dthe+quick+brown+fox\"",
			wanted: "https://www.example.com/search?query=the quick brown fox",
		},
		{
			name:   "URL with encoded special characters",
			input:  "\"https%3A%2F%2Fwww.example.com%2F%3FspecialChars%3D%25%26%2A%28%29%21\"",
			wanted: "https://www.example.com/?specialChars=%&*()!",
		},
		{
			name:   "URL with UTF-8 characters encoded",
			input:  "\"https%3A%2F%2Fwww.example.com%2F%3Futf8%3D%E2%9C%93\"",
			wanted: "https://www.example.com/?utf8=✓",
		},
		{
			name:   "URL with hexadecimal characters",
			input:  "\"https%3A%2F%2Fwww.example.com%2Fpath%3Fhex%3D0x123abc\"",
			wanted: "https://www.example.com/path?hex=0x123abc",
		},
		{
			name:   "URL with mixed encoding",
			input:  "\"https%3A%2F%2Fexample.com%2F%3Fmixed%3D%25AA%2B%2B%2F%2F\"",
			wanted: "https://example.com/?mixed=%AA++//",
		},
		{
			name:   "Incorrectly encoded URL - missing percent encoding",
			input:  "\"https://www.example.com/?foo=bar&baz=qux\"",
			wanted: "https://www.example.com/?foo=bar&baz=qux", // Assuming the function might handle or bypass such cases
		},
		{
			name:   "Encoded URL containing hash symbol",
			input:  "\"https%3A%2F%2Fwww.example.com%2Fpath%23anchor\"",
			wanted: "https://www.example.com/path#anchor",
		},
		{
			name:   "Encoded URL containing multiple query parameters",
			input:  "\"https%3A%2F%2Fwww.example.com%2Fsearch%3Fq%3Dopenai%26page%3D2%26sort%3Dnewest\"",
			wanted: "https://www.example.com/search?q=openai&page=2&sort=newest",
		},
		{
			name:   "URL with encoded slashes and colons",
			input:  "\"https%3A%2F%2Fwww.example.com%2Fpath%2Fto%2F%3Aresource\"",
			wanted: "https://www.example.com/path/to/:resource",
		},
		{
			name:   "Completely encoded URL",
			input:  "\"%68%74%74%70%73%3A%2F%2F%77%77%77%2E%65%78%61%6D%70%6C%65%2E%63%6F%6D\"",
			wanted: "https://www.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(version.Must(version.NewVersion("1.7.0"))), // Temporarily set to 1.7.0 as 1.8.0 is not released yet
				},
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							output "test" {
								value = provider::tfutils::urldecode(%s)
							}
							`, tt.input),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("test", tt.wanted),
						),
					},
				},
			})
		})
	}
}
