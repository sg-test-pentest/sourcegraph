package providers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcegraph/sourcegraph/internal/extsvc"
	"github.com/sourcegraph/sourcegraph/schema"
)

type mockAuthProvider struct {
	configID ConfigID
	config   schema.AuthProviderCommon
}

func (m mockAuthProvider) ConfigID() ConfigID {
	return m.configID
}

func (m mockAuthProvider) Config() schema.AuthProviders {
	return schema.AuthProviders{
		Github: &schema.GitHubAuthProvider{
			Type: m.configID.Type,
		},
	}
}

func (m mockAuthProvider) CachedInfo() *Info {
	panic("should not be called")
}

func (m mockAuthProvider) Refresh(ctx context.Context) error {
	panic("should not be called")
}

func (m mockAuthProvider) ExternalAccountInfo(ctx context.Context, account extsvc.Account) (*extsvc.PublicAccountData, error) {
	panic("should not be called")
}

func TestSortedProviders(t *testing.T) {
	tests := []struct {
		name          string
		input         []Provider
		expectedOrder []int
	}{
		{
			name: "sort works as expected",
			input: []Provider{
				mockAuthProvider{configID: ConfigID{Type: "a", ID: "1"}, config: schema.AuthProviderCommon{Order: 2}},
				mockAuthProvider{configID: ConfigID{Type: "b", ID: "2"}},
				mockAuthProvider{configID: ConfigID{Type: "builtin", ID: "3"}},
				mockAuthProvider{configID: ConfigID{Type: "c", ID: "4"}, config: schema.AuthProviderCommon{Order: 1}},
				mockAuthProvider{configID: ConfigID{Type: "d", ID: "5"}, config: schema.AuthProviderCommon{Order: 1}},
				mockAuthProvider{configID: ConfigID{Type: "b", ID: "6"}, config: schema.AuthProviderCommon{Order: 1}},
			},
			expectedOrder: []int{3, 0, 2, 1, 5, 4},
		},
		{
			name:          "Behaves well for empty slice",
			input:         []Provider{},
			expectedOrder: []int{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			MockProviders = test.input
			t.Cleanup(func() {
				MockProviders = nil
			})

			sorted := SortedProviders()
			expected := make([]Provider, len(sorted))
			for i, order := range test.expectedOrder {
				expected[i] = test.input[order]
			}
			require.ElementsMatch(t, expected, sorted)
		})
	}
}