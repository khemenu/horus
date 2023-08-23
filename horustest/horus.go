package horustest

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"khepri.dev/horus"
	"khepri.dev/horus/service"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/enttest"
)

func WithHorus(conf *horus.Config, f func(require *require.Assertions, h horus.Horus)) func(t *testing.T) {
	return func(t *testing.T) {
		require := require.New(t)

		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", enttest.WithOptions(ent.Log(t.Log)))
		defer client.Close()

		stores, err := store.NewStores(client)
		require.NoError(err)

		horus, err := service.NewHorus(stores, conf)
		require.NoError(err)

		f(require, horus)
	}
}
