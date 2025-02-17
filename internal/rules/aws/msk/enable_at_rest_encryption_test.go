package msk

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/msk"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableAtRestEncryption(t *testing.T) {
	tests := []struct {
		name     string
		input    msk.MSK
		expected bool
	}{
		{
			name: "Cluster with at rest encryption enabled",
			input: msk.MSK{
				Clusters: []msk.Cluster{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						EncryptionAtRest: msk.EncryptionAtRest{
							Metadata:  defsecTypes.NewTestMetadata(),
							KMSKeyARN: defsecTypes.String("foo-bar-key", defsecTypes.NewTestMetadata()),
							Enabled:   defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "Cluster with at rest encryption disabled",
			input: msk.MSK{
				Clusters: []msk.Cluster{
					{
						Metadata: defsecTypes.NewTestMetadata(),
					},
				},
			},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.MSK = test.input
			results := CheckEnableAtRestEncryption.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableAtRestEncryption.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
