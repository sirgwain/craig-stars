package server

import (
	"net/http"

	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
)

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
				{HullComponent: cs.OmegaTorpedo.Name, HullSlotIndex: 5, Quantity: 3},
				{HullComponent: cs.Overthruster.Name, HullSlotIndex: 6, Quantity: 3},
				{HullComponent: cs.GorillaDelagator.Name, HullSlotIndex: 6, Quantity: 4},
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
				{HullComponent: cs.FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
			}))


	fleets := []*cs.Fleet{
		{
			MapObject: cs.MapObject{
				PlayerNum: player1.Num,
			},
			BaseName: "Battle Cruiser",
			Tokens: []cs.ShipToken{
				{
					DesignNum: player1.Designs[0].Num,
					Quantity:  2,
				},
			},
			FleetOrders: cs.FleetOrders{
				BattlePlanName: player1.BattlePlans[0].Name,
			},
		},
		// player2's teamster
		{
			MapObject: cs.MapObject{
				PlayerNum: player2.Num,
			},
			BaseName: "Teamster",
			Tokens: []cs.ShipToken{
				{
					Quantity:  5,
					DesignNum: player2.Designs[0].Num,
				},
				{
					Quantity:  1,
					DesignNum: player2.Designs[1].Num,
				},
			},
			FleetOrders: cs.FleetOrders{
				BattlePlanName: player2.BattlePlans[0].Name,
			},
		}}

	record := cs.RunTestBattle([]*cs.Player{player1, player2}, fleets)
	rest.RenderJSON(w, rest.JSON{"player": player1, "battle": record})
}
