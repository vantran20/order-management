package pg

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTx_Rollback(t *testing.T) {
	// Given:
	q := "SELECT * FROM test_transactions"
	dbConn, err := NewPool(
		os.Getenv("DB_URL"),
		3,
		1,
	)
	require.NoError(t, err)

	_, err = dbConn.Query(q)
	require.Error(t, err, "table 'test_transactions' should not exist at the beginning of this test")

	// When: executed parallely
	// each transaction should be able to create its own table and rollback -- no conflict
	for i := 0; i < 5; i++ {
		t.Run("thread"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			// When:
			require.Error(t, Tx(context.Background(), dbConn, func(tx ContextExecutor) error {
				// Create table
				_, err := tx.Exec("CREATE TABLE test_transactions (id SERIAL PRIMARY KEY)")
				require.NoError(t, err)

				// Query table
				_, err = tx.Query(q)
				require.NoError(t, err)

				return errors.New("some err")
			}))
		})
	}

	// Verify that table is not created (rolled back)
	_, err = dbConn.Query(q)
	require.Error(t, err, "table 'test_transactions' should not exist at the end of this test")
}
