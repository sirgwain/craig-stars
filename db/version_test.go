package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateVersion(t *testing.T) {
	c := connectTestDB()

	version, err := c.getVersion()
	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	tx, err := c.db.Beginx()
	if err != nil {
		t.Errorf("begin transaction %s", err)
		return
	}

	version.Current = 1
	if err := c.updateVersion(version, tx); err != nil {
		t.Errorf("update version %s", err)
		return
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("commit transaction %s", err)
	}

	updated, err := c.getVersion()

	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	assert.Equal(t, version.Current, updated.Current)
	assert.Less(t, version.UpdatedAt, updated.UpdatedAt)

}

func TestEnsureUpgrade(t *testing.T) {
	c := connectTestDB()

	c.createTestFullGame()

	version, err := c.getVersion()
	if err != nil {
		t.Errorf("EnsureUpgrade() failed to getVersion(): %v", err)
		return
	}

	// start at 0, run an upgrade
	assert.Equal(t, 0, version.Current)
	if err := c.ensureUpgrade(); err != nil {
		t.Errorf("EnsureUpgrade() failed: %v", err)
		return
	}

	version, err = c.getVersion()
	if err != nil {
		t.Errorf("EnsureUpgrade() failed to getVersion() after upgrade: %v", err)
		return
	}

	assert.Equal(t, 1, version.Current)
}
