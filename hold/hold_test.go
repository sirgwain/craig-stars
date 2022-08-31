package hold

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirgwain/craig-stars/config"
)

func TestDB_Connect(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Should connect"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DB{}
			dir, err := ioutil.TempDir("../tmp", "hold_testdb_connect")
			if err != nil {
				t.Error(err)
				return
			}
			cfg := &config.Config{}
			cfg.Database.Filename = fmt.Sprintf("%s/bolt.db", dir)
			db.Connect(cfg)
		})
	}
}
