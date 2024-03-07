package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestCidrContains(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		cidr        string
		expected    bool
		expectedErr string // Use this to expect specific error messages
	}{
		{"IPv4 in range", "192.168.1.1", "192.168.1.0/24", true, ""},
		{"IPv4 out of range", "192.168.2.1", "192.168.1.0/24", false, ""},
		{"IPv4, CIDR is IPv6", "192.168.1.1", "::1/128", false, "address is IPv4, but CIDR is IPv6"},
		{"IPv6 in range", "2001:db8::1", "2001:db8::/64", true, ""},
		{"IPv6, CIDR is IPv4", "2001:db8::1", "192.168.1.0/24", false, "address is IPv6, but CIDR is IPv4"},
		{"Invalid address format", "invalid", "192.168.1.0/24", false, "invalid address format"},
		{"Invalid CIDR format", "192.168.1.1", "invalid", false, "invalid CIDR format: invalid CIDR address: invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			contains, err := cidrContains(tc.address, tc.cidr)
			if err != nil && tc.expectedErr == "" {
				t.Errorf("Unexpected error: %v", err)
			} else if err == nil && tc.expectedErr != "" {
				t.Errorf("Expected error '%s', but got none", tc.expectedErr)
			} else if err != nil && tc.expectedErr != "" && err.Error() != tc.expectedErr {
				t.Errorf("Expected error '%s', but got '%s'", tc.expectedErr, err.Error())
			} else if contains != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, contains)
			}
		})
	}
}

func TestCIDRContainsFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		address     string
		cidr        string
		expectedErr *regexp.Regexp
		expected    bool
	}{
		{
			name:        "Address in range",
			address:     "\"192.168.1.1\"",
			cidr:        "\"192.168.1.0/24\"",
			expectedErr: nil,
			expected:    true,
		},
		{
			name:        "Address not in range",
			address:     "\"192.168.2.1\"",
			cidr:        "\"192.168.1.0/24\"",
			expectedErr: nil,
			expected:    false,
		},
		{
			name:        "IPv6 address in range",
			address:     "\"2001:db8::1\"",
			cidr:        "\"2001:db8::/32\"",
			expectedErr: nil,
			expected:    true,
		},
		{
			name:        "IPv6 address not in range",
			address:     "\"2001:db9::1\"",
			cidr:        "\"2001:db8::/32\"",
			expectedErr: nil,
			expected:    false,
		},
		{
			name:        "Exact match",
			address:     "\"192.168.1.0\"",
			cidr:        "\"192.168.1.0/32\"",
			expectedErr: nil,
			expected:    true,
		},
		{
			name:        "Single IP not in single IP CIDR",
			address:     "\"192.168.1.1\"",
			cidr:        "\"192.168.1.0/32\"",
			expectedErr: nil,
			expected:    false,
		},
		{
			name:        "Invalid address format",
			address:     "\"not.an.ip\"",
			cidr:        "\"192.168.1.0/24\"",
			expectedErr: regexp.MustCompile("invalid address format"),
			expected:    false,
		},
		{
			name:        "Invalid CIDR format",
			address:     "\"192.168.1.1\"",
			cidr:        "\"not.a.cidr\"",
			expectedErr: regexp.MustCompile("invalid CIDR format"),
			expected:    false,
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
								value = provider::tfutils::cidrcontains(%s, %s)
							}
							`, tt.cidr, tt.address),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("test", fmt.Sprintf("%t", tt.expected)),
						),
						ExpectError: tt.expectedErr,
					},
				},
			})
		})
	}
}
