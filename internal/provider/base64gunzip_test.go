package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestBase64GunzipFunction(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		input         string
		wanted        string
		expectedError *regexp.Regexp
	}{
		{
			name:          "Valid: Hello world",
			input:         "\"H4sIANL/6WUA/wWAQQkAAAjEqhjHBhbQ3+DA/o/RB6nJswJWsRdKCwAAAA==\"",
			wanted:        "Hello World",
			expectedError: nil,
		},
		{
			name:          "Valid: Empty string",
			input:         "\"H4sIAOQA6mUA/wXAgQgAAAAAIH/rAwAAAAAAAAAA\"",
			wanted:        "",
			expectedError: nil,
		},
		{
			name:          "Invalid: Invalid base64",
			input:         "\"H4sIANL/6WUA/wWAQQkAAAjEqhjHBhbQ3+DA/o/RB6nJswJWsRdKCwAAAA\"",
			wanted:        "",
			expectedError: regexp.MustCompile("illegal base64"),
		},
		{
			name:          "Invalid: Invalid gzip",
			input:         "\"abc\"",
			wanted:        "",
			expectedError: regexp.MustCompile("illegal base64"),
		},
	}

	for _, tt := range cases {
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
								value = provider::tfutils::base64gunzip(%s)
							}
							`, tt.input),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("test", tt.wanted),
						),
						ExpectError: tt.expectedError,
					},
				},
			})
		})
	}
}
