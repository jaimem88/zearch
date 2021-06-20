package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadJSONFile(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		expectedOutput map[string]string
		expectedError  string
	}{
		{
			name:     "valid_json_file",
			filename: "testdata/valid.json",
			expectedOutput: map[string]string{
				"id": "1",
			},
		},
		{
			name:          "invalid_json_file",
			filename:      "testdata/invalid.json",
			expectedError: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out map[string]string
			err := ReadJSONFile(tt.filename, &out)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)

			for k, v := range tt.expectedOutput {
				assert.Equal(t, v, out[k], "value for key did not match")
			}
		})
	}
}
