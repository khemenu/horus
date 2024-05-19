package alias_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/alias"
)

func TestNew(t *testing.T) {
	require.True(t, alias.Validate(alias.New()))
}

func TestValidate(t *testing.T) {
	tcs := []struct {
		desc  string
		given string
	}{
		{
			"single lower letter of alphabet",
			"a",
		},
		{
			"number in middle",
			"a0b",
		},
		{
			"end with number",
			"a0",
		},
		{
			"consecutive numbers",
			"a0012b",
		},
		{
			"split by dash",
			"a-b",
		},
		{
			"split by underscore",
			"a_b",
		},
		{
			"composite",
			"a_-_b-_-c",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)
			require.True(alias.Validate(tc.given))
		})
	}

	tcs = []struct {
		desc  string
		given string
	}{
		{
			"empty",
			"",
		},
		{
			"capital letter of alphabet",
			"A",
		},
		{
			"invalid character",
			"a귤b",
		},
		{
			"begin with number",
			"0a",
		},
		{
			"begin with dash",
			"-a",
		},
		{
			"begin with underscore",
			"_a",
		},
		{
			"end with dash",
			"a-",
		},
		{
			"end with underscore",
			"a_",
		},
		{
			"consecutive dashes",
			"a--b",
		},
		{
			"consecutive underscores",
			"a__b",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)
			require.False(alias.Validate(tc.given))
		})
	}
}
