package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

// Each candidate must supply four fighters and two bombers on its own. Put an
// incomplete fleet first to catch ship counts leaking between candidates.
func Test_fleet_getFleetsMatchingMakeup(t *testing.T) {
	player := cs.NewPlayer(1, cs.NewRace())
	player.Num = 1
	player.Designs = []*cs.ShipDesign{
		cs.NewShipDesign(player.Num, 1).WithPurpose(cs.ShipDesignPurposeBeamFighter),
		cs.NewShipDesign(player.Num, 2).WithPurpose(cs.ShipDesignPurposeBeamFighter),
		cs.NewShipDesign(player.Num, 3).WithPurpose(cs.ShipDesignPurposeBomber),
	}
	ai := &aiPlayer{Player: player}
	makeup := fleet{
		purpose: cs.FleetPurposeBomber,
		ships: []fleetShip{
			{purpose: cs.ShipDesignPurposeBeamFighter, quantity: 4},
			{purpose: cs.ShipDesignPurposeBomber, quantity: 2},
		},
	}
	tests := []struct {
		name    string
		purpose cs.FleetPurpose
		tokens  []cs.ShipToken
		want    bool
	}{
		{
			name:    "incomplete fleets cannot combine their counts",
			purpose: cs.FleetPurposeBomber,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 2}, {DesignNum: 3, Quantity: 1}},
		},
		{
			name:    "complete fleet matches",
			purpose: cs.FleetPurposeBomber,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 4}, {DesignNum: 3, Quantity: 2}},
			want:    true,
		},
		{
			name:    "different fighter designs count toward the same purpose",
			purpose: cs.FleetPurposeBomber,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 2}, {DesignNum: 2, Quantity: 2}, {DesignNum: 3, Quantity: 2}},
			want:    true,
		},
		{
			name:    "extra fighters do not prevent a match",
			purpose: cs.FleetPurposeBomber,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 4}, {DesignNum: 2, Quantity: 4}, {DesignNum: 3, Quantity: 2}},
			want:    true,
		},
		{
			name:    "extra fighters cannot replace missing bombers",
			purpose: cs.FleetPurposeBomber,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 4}, {DesignNum: 2, Quantity: 4}},
		},
		{
			name:    "fleet assigned to another purpose is excluded",
			purpose: cs.FleetPurposeScout,
			tokens:  []cs.ShipToken{{DesignNum: 1, Quantity: 4}, {DesignNum: 3, Quantity: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incomplete := cs.NewFleet(player, 1, "Incomplete", []cs.Waypoint{{}})
			incomplete.SetTag(cs.TagPurpose, string(makeup.purpose))
			incomplete.Tokens = []cs.ShipToken{{DesignNum: 1, Quantity: 2}, {DesignNum: 3, Quantity: 1}}
			candidate := cs.NewFleet(player, 2, "Candidate", []cs.Waypoint{{}})
			candidate.SetTag(cs.TagPurpose, string(tt.purpose))
			candidate.Tokens = tt.tokens

			matches := makeup.getFleetsMatchingMakeup(ai, []*cs.Fleet{incomplete, candidate})
			if tt.want {
				assert.Equal(t, []*cs.Fleet{candidate}, matches)
			} else {
				assert.Empty(t, matches)
			}
		})
	}
}
