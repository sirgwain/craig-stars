package server

import (
	"net/http"

	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
)

// send a test battle to the user. This is hardcoded for now, but eventually it will support custom battles designed by the player.
func (s *server) testBattle(w http.ResponseWriter, r *http.Request) {
	player1 := cs.NewPlayer(0, cs.NewRace()).WithNum(1)
	player2 := cs.NewPlayer(0, cs.NewRace().WithLRT(cs.RS)).WithNum(2)
	player1.Name = cs.AINames[0] + "s"
	player2.Name = cs.AINames[1] + "s"
	player1.Race.PluralName = cs.AINames[0] + "s"
	player2.Race.PluralName = cs.AINames[1] + "s"
	player2.Color = colors[1]
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
			WithName("Long Range Scout").
			WithHull(cs.Scout.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.CompletePhaseShield.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		cs.NewShipDesign(player1, 3).
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
		cs.NewShipDesign(player1, 4).
			WithName("Accelerator Platform").
			WithHull(cs.OrbitalFort.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.MassDriver5.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.Laser.Name, HullSlotIndex: 2, Quantity: 6},
				{HullComponent: cs.Laser.Name, HullSlotIndex: 5, Quantity: 6},
			}),
		cs.NewShipDesign(player1, 5).
			WithName("Frigate").
			WithHull(cs.Frigate.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.PhaserBazooka.Name, HullSlotIndex: 3, Quantity: 3},
				{HullComponent: cs.CowHideShield.Name, HullSlotIndex: 4, Quantity: 2},
			}),
		cs.NewShipDesign(player1, 6).
			WithName("Destroyer v1").
			WithHull(cs.Destroyer.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.TransGalacticDrive.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.ManeuveringJet.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: cs.BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
				{HullComponent: cs.Strobnium.Name, HullSlotIndex: 5, Quantity: 2},
				{HullComponent: cs.BetaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.BetaTorpedo.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: cs.BetaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
			}),
		cs.NewShipDesign(player1, 7).
			WithName("INT-1MF-SCP - Dragonfly").
			WithHull(cs.Destroyer.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.FuelTank.Name, HullSlotIndex: 6, Quantity: 1},
				{HullComponent: cs.BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
				{HullComponent: cs.PulsedSapper.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.BattleComputer.Name, HullSlotIndex: 4, Quantity: 1},
				{HullComponent: cs.PhaserBazooka.Name, HullSlotIndex: 3, Quantity: 1},
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
			WithName("FF 8").
			WithHull(cs.Frigate.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.PhaserBazooka.Name, HullSlotIndex: 3, Quantity: 3},
				{HullComponent: cs.CowHideShield.Name, HullSlotIndex: 4, Quantity: 2},
			}),
		cs.NewShipDesign(player2, 6).
			WithName("DD S 10").
			WithHull(cs.Destroyer.Name).
			WithSlots([]cs.ShipDesignSlot{
				{HullComponent: cs.FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
				{HullComponent: cs.ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 1},
				{HullComponent: cs.ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
				{HullComponent: cs.BearNeutrinoBarrier.Name, HullSlotIndex: 4, Quantity: 1},
			}),
		cs.NewShipDesign(player2, 7).
			WithName("Empty Spacedock").
			WithHull(cs.SpaceDock.Name).
			WithSlots([]cs.ShipDesignSlot{}),
	)

	p1BattleCruiser := player1.Designs[0]
	p1Scout := player1.Designs[1]
	p1StalwartDefender := player1.Designs[2]
	p1AcceleratorPlatform := player1.Designs[3]
	p1DestroyerV1 := player1.Designs[5]
	p1DragonFly := player1.Designs[6]
	p2Teamster := player2.Designs[0]
	p2Scout := player2.Designs[1]
	p2StalwartDefender := player2.Designs[2]
	p2Sapper := player2.Designs[3]
	p2Frigate := player2.Designs[4]
	p2DDS10 := player2.Designs[5]
	p2EmptySpaceDock := player2.Designs[6]
	_ = p1BattleCruiser
	_ = p1Scout
	_ = p1StalwartDefender
	_ = p1AcceleratorPlatform
	_ = p1DestroyerV1
	_ = p2Teamster
	_ = p2Scout
	_ = p2StalwartDefender
	_ = p2Sapper
	_ = p2Frigate
	_ = p2DDS10
	_ = p2EmptySpaceDock

	fleets := []*cs.Fleet{
		// player1's tokens
		{
			MapObject: cs.MapObject{
				PlayerNum: player1.Num,
			},
			BaseName: "Fleet1",
			Tokens: []cs.ShipToken{
				{
					DesignNum: p1DragonFly.Num,
					Quantity:  1,
				},
			},
		},
		// player2's tokens
		{
			MapObject: cs.MapObject{
				PlayerNum: player2.Num,
			},
			BaseName: "Fleet2",
			Tokens: []cs.ShipToken{
				{
					Quantity:  1,
					DesignNum: p2Scout.Num,
				},
			},
		}}

	record := cs.RunTestBattle([]*cs.Player{player1, player2}, fleets)
	rest.RenderJSON(w, rest.JSON{"player": player1, "battle": record, "fleets": fleets})
}
