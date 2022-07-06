package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// compare two objects as json outputs
// if the comparison fails, this writes a want.json and got.json to the tmp folder
func CompareAsJSON(t *testing.T, got interface{}, want interface{}) bool {
	if got == nil && want == nil {
		return true
	} else if got == nil && want != nil || got != nil && want == nil {
		return false
	} else {
		gotJson, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Errorf("Failed to compare %s, error = %v", got, err)
		}
		wantJson, err := json.MarshalIndent(want, "", "  ")
		if err != nil {
			t.Errorf("Failed to compare %s, error = %v", want, err)
		}

		if string(gotJson) != string(wantJson) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Debug().Msgf("\n\ngot: %s\n", string(gotJson))
			log.Debug().Msgf("\n\nwant: %s\n", string(wantJson))

			_ = ioutil.WriteFile("../tmp/got.json", gotJson, 0644)
			_ = ioutil.WriteFile("../tmp/want.json", wantJson, 0644)
			return false
		} else {
			return true
		}
	}
}
