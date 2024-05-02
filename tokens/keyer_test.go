package tokens_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/tokens"
)

func TestArgon2iKeyer(t *testing.T) {
	require := require.New(t)

	keyer := tokens.NewArgon2iKeyer(tokens.Argon2iKeyerInit{
		Time:    3,
		Memory:  32 * (1 << 10),
		Threads: 4,
		KeyLen:  32,
	})

	given := []byte("Royale with Cheese")
	h1, err := keyer.Key(given)
	require.NoError(err)

	err = keyer.Compare(given, h1)
	require.NoError(err)

	h2, err := keyer.Key([]byte("Le Big Mac"))
	require.NoError(err)

	err = keyer.Compare(given, h2)
	require.Error(err)
}
