package db

import (
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migrateDatabase_backups(t *testing.T) {
	datasource := filepath.Join(t.TempDir(), "data.db")
	backups := func() []string {
		matches, err := filepath.Glob(datasource + ".*.*")
		require.NoError(t, err)
		return matches
	}
	schema := fstest.MapFS{
		"schema/0001_first.up.sql": {Data: []byte("CREATE TABLE a (id INTEGER);")},
	}
	c := &dbConn{}

	// a new database has nothing worth backing up
	require.NoError(t, c.migrateDatabase(datasource, schema, "schema"))
	assert.Empty(t, backups())

	// neither does starting again with nothing to migrate
	require.NoError(t, c.migrateDatabase(datasource, schema, "schema"))
	assert.Empty(t, backups())

	// a pending migration backs up the database as it was
	schema["schema/0002_second.up.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE b (id INTEGER);")}
	require.NoError(t, c.migrateDatabase(datasource, schema, "schema"))
	require.Len(t, backups(), 1)
	assert.Contains(t, backups()[0], "data.db.1.")

	// a failed migration is backed up once, not every time the server tries to start
	schema["schema/0003_broken.up.sql"] = &fstest.MapFile{Data: []byte("CREATE TABLE a (id INTEGER);")}
	require.Error(t, c.migrateDatabase(datasource, schema, "schema"))
	require.Len(t, backups(), 2)
	err := c.migrateDatabase(datasource, schema, "schema")
	require.ErrorContains(t, err, "a previous migration to version 3 failed")
	assert.Len(t, backups(), 2)
}
