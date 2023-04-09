package validator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		checkResult func(err error)
	}{
		{
			name:  "OK",
			value: "abc@ok.com",
			checkResult: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:  "Invalid Characters",
			value: "abc{}@ok.com",
			checkResult: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid email")
			},
		},
		{
			name:  "Invalid Email",
			value: "abc@ok",
			checkResult: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid email")
			},
		},
		{
			name:  "Missing @",
			value: "abc.ok",
			checkResult: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid email")
			},
		},
	}

	for index := range testCases {
		tc := testCases[index]

		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.value)
			tc.checkResult(err)
		})
	}
}
