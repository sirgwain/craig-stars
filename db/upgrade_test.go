package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ensureUpgrade(t *testing.T) {
	c, conn := connectTestDBWithConn()

	c.createTestFullGame(t.Context())

	version, err := c.getVersion(t.Context())
	if err != nil {
		t.Errorf("ensureUpgrade() failed to getVersion(): \n%v", err)
		return
	}

	// start at 0, run an upgrade
	assert.Equal(t, int64(0), version.Current)
	if err := conn.WrapInTransaction(func(c Client) error {
		return c.ensureUpgrade(context.Background())
	}); err != nil {
		t.Errorf("ensureUpgrade() failed: \n%v", err)
	}

	version, err = c.getVersion(t.Context())
	if err != nil {
		t.Errorf("ensureUpgrade() failed to getVersion() after upgrade: \n%v", err)
		return
	}

	assert.Equal(t, LATEST_VERSION, version.Current)
}
