package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateVersion(t *testing.T) {
	c := connectTestDB()

	version, err := c.getVersion(t.Context())
	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	version.Current = 1
	if err := c.updateVersion(t.Context(), version); err != nil {
		t.Errorf("update version %s", err)
		return
	}

	updated, err := c.getVersion(t.Context())

	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	assert.Equal(t, version.Current, updated.Current)

}
