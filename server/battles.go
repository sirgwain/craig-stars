package server

import (
	"net/http"

	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
)

// send a test battle to the user. This is hardcoded for now, but eventually it will support custom battles designed by the player.
func (s *server) testBattle(w http.ResponseWriter, r *http.Request) {
	player1 := cs.NewPlayer(0, cs.NewRace()).WithNum(1)
	player2 := cs.NewPlayer(0, cs.NewRace()).WithNum(2)
	player1.Name = cs.AINames[0] + "s"
	player2.Name = cs.AINames[1] + "s"
	player1.Race.PluralName = cs.AINames[0] + "s"
	player2.Race.PluralName = cs.AINames[1] + "s"
	player1.Relations = []cs.PlayerRelationship{{Relation: cs.PlayerRelationFriend}, {Relation: cs.PlayerRelationEnemy}}
	player2.Relations = []cs.PlayerRelationship{{Relation: cs.PlayerRelationEnemy}, {Relation: cs.PlayerRelationFriend}}
	player1.PlayerIntels.PlayerIntels = []cs.PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}
	player2.PlayerIntels.PlayerIntels = []cs.PlayerIntel{{Num: player1.Num}, {Num: player2.Num}}

	player1.Designs = append(player1.Designs,
		cs.NewShipDesign(player1, 1).
			WithName("Battle Cruiser").
			WithHull(cs.BattleCruiser.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 2, Quantity: 2},
				{HullComponent: cs.BattleSuperComputer.Name, HullSlotIndex: 3, Quantity: 2},
				{HullComponent: cs.ColloidalPhaser.Name, HullSlotIndex: 4, Quantity: 3},
				{HullComponent: cs.HeavyBlaster.Name, HullSlotIndex: 5, Quantity: 3},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 6, Quantity: 3},
				{HullComponent: cs.GorillaDelagator.Name, HullSlotIndex: 7, Quantity: 4},
			}),
		cs.NewShipDesign(player1, 2).
			WithName("Accelerator Platform").
			WithHull(cs.OrbitalFort.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.Laser.Name, HullSlotIndex: 2, Quantity: 6},
				{HullComponent: cs.Laser.Name, HullSlotIndex: 5, Quantity: 6},
			}),
	)

	player2.Designs = append(player2.Designs,
		cs.NewShipDesign(player2, 1).
			WithName("Teamster").
			WithHull(cs.SmallFreighter.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		cs.NewShipDesign(player2, 2).
			WithName("Long Range Scout").
			WithHull(cs.Scout.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.CompletePhaseShield.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		cs.NewShipDesign(player2, 3).
			WithName("Stalwart Defender").
			WithHull(cs.Destroyer.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.OmegaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: cs.Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: cs.BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
			}),
		cs.NewShipDesign(player2, 4).
			WithName("Stalwart Sapper").
			WithHull(cs.Destroyer.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.PulsedSapper.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.PulsedSapper.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: cs.Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 7, Quantity: 1},
			}),
		cs.NewShipDesign(player2, 5).
			WithName("Frigate").
			WithHull(cs.Frigate.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.BatScanner.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.Laser.Name, HullSlotIndex: 3, Quantity: 2},
				{HullComponent: cs.WolverineDiffuseShield.Name, HullSlotIndex: 4, Quantity: 2},
			}),
	)

	fleets := []*cs.Fleet{
		// player1's platform
		{
			MapObject: cs.MapObject{
				PlayerNum: player1.Num,
			},
			BaseName: "Accelerator Platform",
			Tokens: []cs.ShipToken{
				{
					DesignNum: player1.Designs[1].Num,
					Quantity:  1,
				},
			},
		},
		// player2's stalwart defenders
		{
			MapObject: cs.MapObject{
				PlayerNum: player2.Num,
			},
			BaseName: "Frigate",
			Tokens: []cs.ShipToken{
				{
					Quantity:  3,
					DesignNum: player2.Designs[4].Num,
				},
				{
					Quantity:  7,
					DesignNum: player2.Designs[4].Num,
				},
			},
		}}

	record := cs.RunTestBattle([]*cs.Player{player1, player2}, fleets)
	rest.RenderJSON(w, rest.JSON{"player": player1, "battle": record, "fleets": fleets})
}
