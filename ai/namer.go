package ai

import (
	"github.com/golodash/godash/strings"
	"github.com/sirgwain/craig-stars/cs"
)

func (ai *aiPlayer) fleetName(purpose cs.FleetPurpose) string {

	// TODO: do some cool names
	return strings.StartCase(string(purpose))
}
