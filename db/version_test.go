package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateVersion(t *testing.T) {
	c := connectTestDB()
	defer func() { closeTestDB(c) }()

	version, err := c.getVersion()
	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	version.Current = 1
	if err := c.updateVersion(version); err != nil {
		t.Errorf("update version %s", err)
		return
	}

	updated, err := c.getVersion()

	if err != nil {
		t.Errorf("get version %s", err)
		return
	}

	assert.Equal(t, version.Current, updated.Current)
	assert.Less(t, version.UpdatedAt, updated.UpdatedAt)

}
