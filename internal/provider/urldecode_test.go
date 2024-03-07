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
			name:   "Standard URL to decode",
			input:  "\"https%3A%2F%2Fwww.example.com%2F%3Ffoo%3Dbar%26baz%3Dqux\"",
			wanted: "https://www.example.com/?foo=bar&baz=qux",
		},
		{
			name:   "empty string",
			input:  "\"\"",
			wanted: "",
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
