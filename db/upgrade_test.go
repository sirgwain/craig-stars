package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureUpgrade(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

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

	assert.Equal(t, LATEST_VERSION, version.Current)
}
