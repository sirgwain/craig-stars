package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureUpgrade(t *testing.T) {
	c := connectTestDB()

	c.createTestFullGame(t.Context())

	version, err := c.getVersion(t.Context())
	if err != nil {
		t.Errorf("EnsureUpgrade() failed to getVersion(): \n%v", err)
		return
	}

	// start at 0, run an upgrade
	assert.Equal(t, int64(0), version.Current)
	if err := c.ensureUpgrade(t.Context()); err != nil {
		t.Errorf("EnsureUpgrade() failed: \n%v", err)
		return
	}

	version, err = c.getVersion(t.Context())
	if err != nil {
		t.Errorf("EnsureUpgrade() failed to getVersion() after upgrade: \n%v", err)
		return
	}

	assert.Equal(t, LATEST_VERSION, version.Current)
}
