package tokens_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/tokens"
)

func TestArgon2iKeyer(t *testing.T) {
	require := require.New(t)

	keyer := tokens.NewArgon2i(&tokens.Argon2State{
		Parallelism: 4,
		TagLength:   32,
		MemorySize:  32 * (1 << 10),
		Iterations:  3,
	})

	given := []byte("Royale with Cheese")
	k1, err := keyer.Key(given)
	require.NoError(err)

	err = k1.Compare(given)
	require.NoError(err)

	k2, err := keyer.Key([]byte("Le Big Mac"))
	require.NoError(err)

	err = k2.Compare(given)
	require.Error(err)
}
